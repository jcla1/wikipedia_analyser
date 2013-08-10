package main

import (
	"encoding/gob"
	"fmt"
	"github.com/jcla1/matrix"
	"github.com/jcla1/minimize"
	"github.com/jcla1/nn"
	"io"
	"os"
)

const (
	numHidden   int     = 3
	numFeatures int     = 9
	numIter     int     = 3000
	lambda      float64 = 0.0
)

func ProcessRecords(input <-chan *matrix.Matrix, n *NN, thresholdMat *matrix.Matrix) <-chan *matrix.Matrix {
	out := make(chan *matrix.Matrix)

	evaluateChannelIn := make(chan *matrix.Matrix)

	evaluateChannelOut := EvaluateNN(n, evaluateChannelIn)

	go func() {
		for m := range input {
			evaluateChannelIn <- m
			evOut := <-evaluateChannelOut

			diff, _ := m.Sub(evOut)
			predictionErr := diff.Power(2.0)

			wellnessSlice := make([]float64, numFeatures)

			for i, _ := range wellnessSlice {
				if predictionErr.Vals[i] <= thresholdMat.Vals[i] {
					wellnessSlice[i] = 1.0
				} else {
					wellnessSlice[i] = 0.0
				}
			}

			out <- matrix.FromSlice(wellnessSlice, numFeatures, 1)
		}
		close(out)
	}()

	return out
}

// Returns min and range
func FindNormalizedVectors(in <-chan *matrix.Matrix) (*matrix.Matrix, *matrix.Matrix) {
	p := <-in
	max := make([]float64, numFeatures)
	min := make([]float64, numFeatures)

	copy(max, p.Vals)
	copy(min, p.Vals)

	for p = range in {
		for i, v := range p.Vals {
			if max[i] < v {
				max[i] = v
			}

			if min[i] > v {
				min[i] = v
			}
		}
	}

	maxVec := matrix.FromSlice(max, numFeatures, 1)
	minVec := matrix.FromSlice(min, numFeatures, 1)

	r, _ := maxVec.Sub(minVec)
	return minVec, r
}

func oneOver(i int, v float64) float64 {
	return 1.0 / v
}

func Normalizer(in <-chan *matrix.Matrix, min, r *matrix.Matrix) <-chan *matrix.Matrix {
	out := make(chan *matrix.Matrix)

	rOver := r.Copy()
	rOver.Apply(oneOver)

	var v *matrix.Matrix

	go func() {
		for vec := range in {
			v, _ = vec.Sub(min)
			out <- v.EWProd(rOver)
		}
		close(out)
	}()

	return out
}

// Returns min, max of the deviation of the predictionError
func FindPredictionErrDeviation(input <-chan *PageContainer, n *NN, min, r *matrix.Matrix) (*matrix.Matrix, *matrix.Matrix) {
	evaluateChannelIn := make(chan *PageContainer)
	vectorizerChannelIn := make(chan *PageContainer)

	evaluateChannelOut := EvaluateNN(n, Normalizer(Vectorizer(evaluateChannelIn), min, r))
	vectorizerChannelOut := Normalizer(Vectorizer(vectorizerChannelIn), min, r)

	/*devMax := make([]float64, numFeatures)
	for i, _ := range devMax {
		devMax[i] = 0.0
	}

	for p := range input {
		evaluateChannelIn <- p
		vectorizerChannelIn <- p
		out := <-evaluateChannelOut
		in := <-vectorizerChannelOut

		diff, _ := in.Sub(out)
		predictionErr := diff.Power(2.0)

		for i, v := range predictionErr.Vals {
			if devMax[i] < v {
				devMax[i] = v
			}
		}

	}

	close(evaluateChannelIn)
	close(vectorizerChannelIn)

	return matrix.FromSlice(devMax, numFeatures, 1)*/

	vals := make([][]float64, numFeatures)
	for i, _ := range vals {
		vals[i] = make([]float64, 0, 1000)
	}

	for p := range input {
		evaluateChannelIn <- p
		vectorizerChannelIn <- p
		out := <-evaluateChannelOut
		in := <-vectorizerChannelOut

		diff, _ := in.Sub(out)
		predictionErr := diff.Power(2.0)

		for i, v := range predictionErr.Vals {
			vals[i] = append(vals[i], v)
		}

	}

	devs := make([]float64, numFeatures)
	means := make([]float64, numFeatures)

	for i, _ := range vals {
		devs[i], means[i] = calcStdDev(vals[i])
	}

	return matrix.FromSlice(devs, numFeatures, 1), matrix.FromSlice(means, numFeatures, 1)
}

func SaveNNToFile(n *NN, path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	SaveNN(n, file)
}

func LoadNNFromFile(path string) *NN {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return LoadNN(file)
}

func LoadNN(r io.Reader) *NN {
	var n *NN
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&n)
	if err != nil {
		panic(err)
	}
	return n
}

func SaveNN(n *NN, w io.Writer) {
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(n)
	if err != nil {
		panic(err)
	}
}

func EvaluateNN(n *NN, input <-chan *matrix.Matrix) <-chan *matrix.Matrix {
	c := make(chan *matrix.Matrix)

	tes := MakeTrainingEx(input)

	go func() {
		var te nn.TrainingExample
		var ok bool

		for {
			te, ok = <-tes

			if !ok {
				close(c)
				return
			} else {
				c <- nn.Hypothesis(n.Thetas, te)
			}
		}

	}()

	return c
}

func BuildTrainNN(input <-chan *PageContainer, min, r *matrix.Matrix) *NN {
	d := FeatureSlicer(MakeTrainingEx(Normalizer(Vectorizer(input), min, r)))
	n := NewNN([]int{numFeatures, numHidden, numFeatures}, lambda)
	TrainNN(n, d, numIter)

	return n
}

func TrainNN(n *NN, data []nn.TrainingExample, iter int) {
	f := SetupCostGradFunc(n, data)
	fmt.Println("Cost before:", nn.CostFunction(data, n.Thetas, n.Lambda))
	n.Thetas = ReshapeParams(minimize.Fmincg(f, UnrollParams(n.Thetas), iter, true), n.LayerSizes)
	fmt.Println("Cost after:", nn.CostFunction(data, n.Thetas, n.Lambda))
}

func FeatureSlicer(input <-chan nn.TrainingExample) []nn.TrainingExample {
	features := make([]nn.TrainingExample, 0, 2000)
	for v := range input {
		features = append(features, v)
	}

	return features
}

type NN struct {
	LayerSizes []int
	Thetas     nn.Parameters
	Lambda     float64
}

func NewNN(layerSizes []int, lambda float64) *NN {
	n := new(NN)

	n.LayerSizes = layerSizes
	n.Lambda = lambda

	n.InitParameters()

	return n
}

func (n *NN) InitParameters() {
	n.Thetas = make(nn.Parameters, len(n.LayerSizes)-1)
	for i, _ := range n.Thetas {
		n.Thetas[i] = matrix.Rand(n.LayerSizes[i+1], n.LayerSizes[i]+1)
	}
}

func SetupCostGradFunc(n *NN, data []nn.TrainingExample) minimize.CostGradientFunc {

	f := func(t *matrix.Matrix) (float64, *matrix.Matrix) {

		thetas := ReshapeParams(t, n.LayerSizes)

		cost := nn.CostFunction(data, thetas, n.Lambda)
		gradients := nn.BackProp(thetas, data, n.Lambda)

		return cost, UnrollParams(gradients)
	}

	return f
}

func ReshapeParams(val *matrix.Matrix, layerSizes []int) []*matrix.Matrix {
	vals := val.Values()
	thetas := make(nn.Parameters, 0, len(layerSizes)-1)

	offset := 0
	for i := 0; i < len(layerSizes)-1; i++ {
		thetas = append(thetas, matrix.FromSlice(vals[offset:offset+(layerSizes[i+1]*layerSizes[i]+1)], layerSizes[i+1], layerSizes[i]+1))
		offset += layerSizes[i+1]*layerSizes[i] + 1
	}

	return thetas
}

func UnrollParams(vals []*matrix.Matrix) *matrix.Matrix {
	unrolled := make([]float64, 0)

	for _, v := range vals {
		unrolled = append(unrolled, v.Values()...)
	}

	return matrix.FromSlice(unrolled, 1, len(unrolled))
}

func MakeFeatureVector(container *PageContainer) *matrix.Matrix {
	features := []float64{
		container.AvgSentenceLen,
		container.LinkDensity,
		float64(container.NumLinks),
		float64(container.NumExternalLinks),
		float64(container.NumHeadings),
		float64(container.NumRefs),
		float64(container.NumCategories),
		float64(container.NumSentences),
		float64(container.NumWords),

		//float64(len(container.Unigrams)),
		//float64(len(container.Bigrams)),
		//float64(len(container.Trigrams)),
	}

	return matrix.FromSlice(features, len(features), 1)
}

func Vectorizer(input <-chan *PageContainer) <-chan *matrix.Matrix {
	outputChannel := make(chan *matrix.Matrix)

	go func() {
		var container *PageContainer
		var ok bool

		for {
			container, ok = <-input

			if !ok {
				// channel closed!
				close(outputChannel)
				return
			} else {
				outputChannel <- MakeFeatureVector(container)
			}
		}
	}()

	return outputChannel
}

func MakeTrainingEx(input <-chan *matrix.Matrix) <-chan nn.TrainingExample {
	outputChannel := make(chan nn.TrainingExample)

	go func() {
		var m *matrix.Matrix
		var ok bool

		for {
			m, ok = <-input

			if !ok {
				// channel closed!
				close(outputChannel)
				return
			} else {
				outputChannel <- nn.TrainingExample{m, m}
			}
		}
	}()

	return outputChannel
}
