package main

import (
	"encoding/gob"
	"github.com/jcla1/matrix"
	"github.com/jcla1/minimize"
	"github.com/jcla1/nn"
	"io"
	"os"
)

const (
	numHidden   int = 3
	numFeatures int = 8
	numIter     int = 200
)

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

func EvaluateNN(n *NN, input <-chan *PageContainer) <-chan *PageContainer {
	c := make(chan *PageContainer)

	go func() {
		var features []float64
		var m *matrix.Matrix
		var container *PageContainer
		var ok bool
		var te nn.TrainingExample

		for {
			if !ok {
				close(c)
				return
			} else {
				features = []float64{
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
				m = matrix.FromSlice(features, len(features), 1)
				te = nn.TrainingExample{m, m}

				container.NNOutput = nn.Hypothesis(n.Thetas, te)
				c <- container
			}
		}

	}()

	return c
}

func BuildTrainNN(input <-chan *PageContainer) *NN {
	d := FeatureSlicer(FeatureVectors(input))
	n := NewNN([]int{numFeatures, numHidden, numFeatures}, 0)
	TrainNN(n, d, numIter)

	return n
}

func TrainNN(n *NN, data []nn.TrainingExample, iter int) {
	f := SetupCostGradFunc(n, data)
	n.Thetas = ReshapeParams(minimize.Fmincg(f, UnrollParams(n.Thetas), iter, false), n.LayerSizes)
}

func FeatureSlicer(input <-chan nn.TrainingExample) []nn.TrainingExample {
	features := make([]nn.TrainingExample, 0, 100)
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

func FeatureVectors(input <-chan *PageContainer) <-chan nn.TrainingExample {
	outputChannel := make(chan nn.TrainingExample)

	go func() {
		var features []float64
		var m *matrix.Matrix
		var container *PageContainer
		var ok bool

		for {
			container, ok = <-input

			if !ok {
				// channel closed!
				close(outputChannel)
				return
			} else {
				features = []float64{
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
				m = matrix.FromSlice(features, len(features), 1)
				outputChannel <- nn.TrainingExample{m, m}
			}
		}
	}()

	return outputChannel
}
