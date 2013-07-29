package main

import (
	"github.com/jcla1/matrix"
	"github.com/jcla1/nn"
	"github.com/jcla1/minimize"
)

const (
	bottleNeckFactor float64 = 0.3
)

type NN struct {
	LayerSizes []int
	Thetas     nn.Parameters
	Lambda     float64
}

func SetupCostGradFunc(n *NN, data []nn.TrainingExample) minimize.CostGradientFunc {

	f := func(t *matrix.Matrix) (float64, *matrix.Matrix) {
		vals := t.Values()
		thetas := make(nn.Parameters, 0, len(n.LayerSizes)-1)

		offset := 0
		for i := 0; i < len(n.LayerSizes)-1; i++ {
			thetas = append(thetas, matrix.FromSlice(vals[offset:offset+(n.LayerSizes[i+1]*n.LayerSizes[i]+1)], n.LayerSizes[i+1], n.LayerSizes[i]+1))
			offset += n.LayerSizes[i+1]*n.LayerSizes[i] + 1
		}

		cost := nn.CostFunction(data, thetas, n.Lambda)
		gradients := nn.BackProp(thetas, data, n.Lambda)

		unrolled := make([]float64, 0)

		for _, grad := range gradients {
			unrolled = append(unrolled, grad.Values()...)
		}

		return cost, matrix.FromSlice(unrolled, 1, len(unrolled))
	}

	return f
}

func FeatureVectors(input <-chan *PageContainer) <-chan *matrix.Matrix {
	outputChannel := make(chan *matrix.Matrix)

	go func() {
		var features []float64
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

					/*float64(len(container.Unigrams)),
					  float64(len(container.Bigrams)),
					  float64(len(container.Trigrams)),*/
				}

				outputChannel <- matrix.FromSlice(features, len(features), 1)
			}
		}
	}()

	return outputChannel
}
