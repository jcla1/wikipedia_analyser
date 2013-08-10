package main

import (
	"fmt"
	"math"
	"time"
)

func featureStats(input <-chan *PageContainer, n int) {

	containers := make([]*PageContainer, 0, n)

	for i := 0; i < n; i++ {
		containers = append(containers, <-input)
	}

	fmt.Println("Standard Deviation of features:")

	var avg float64
	var stdDev float64

	stdDev, avg = calcStats(containers, func(page *PageContainer) float64 {
		return float64(page.NumLinks)
	})

	fmt.Printf("# of Links:          %.3f (avg %.3f)\n", stdDev, avg)

	stdDev, avg = calcStats(containers, func(page *PageContainer) float64 {
		return page.LinkDensity
	})

	fmt.Printf("Link Density:        %.3f (avg %.3f)\n", stdDev, avg)

	stdDev, avg = calcStats(containers, func(page *PageContainer) float64 {
		return float64(page.NumExternalLinks)
	})

	fmt.Printf("# of External Links: %.3f (avg %.3f)\n", stdDev, avg)

	stdDev, avg = calcStats(containers, func(page *PageContainer) float64 {
		return float64(page.NumHeadings)
	})

	fmt.Printf("# of Headings:       %.3f (avg %.3f)\n", stdDev, avg)

	stdDev, avg = calcStats(containers, func(page *PageContainer) float64 {
		return float64(page.NumRefs)
	})

	fmt.Printf("# of Refrences:      %.3f (avg %.3f)\n", stdDev, avg)

	stdDev, avg = calcStats(containers, func(page *PageContainer) float64 {
		return float64(page.NumCategories)
	})

	fmt.Printf("# of Categories:     %.3f (avg %.3f)\n", stdDev, avg)

	stdDev, avg = calcStats(containers, func(page *PageContainer) float64 {
		return float64(page.NumSentences)
	})

	fmt.Printf("# of Sentences:      %.3f (avg %.3f)\n", stdDev, avg)

	stdDev, avg = calcStats(containers, func(page *PageContainer) float64 {
		return page.AvgSentenceLen
	})

	fmt.Printf("Avg Sentence Length: %.3f (avg %.3f)\n", stdDev, avg)

	stdDev, avg = calcStats(containers, func(page *PageContainer) float64 {
		return float64(page.NumWords)
	})

	fmt.Printf("# of Words:          %.3f (avg %.3f)\n", stdDev, avg)

	stdDev, avg = calcStats(containers, func(page *PageContainer) float64 {
		return page.AvgWordLen
	})

	fmt.Printf("Avg Word Length:     %.3f (avg %.3f)\n", stdDev, avg)

	stdDev, avg = calcStats(containers, func(page *PageContainer) float64 {
		return float64(len(page.Unigrams))
	})

	fmt.Printf("# of Unigrams:       %.3f (avg %.3f)\n", stdDev, avg)

	stdDev, avg = calcStats(containers, func(page *PageContainer) float64 {
		return float64(len(page.Bigrams))
	})

	fmt.Printf("# of Bigrams:        %.3f (avg %.3f)\n", stdDev, avg)

	stdDev, avg = calcStats(containers, func(page *PageContainer) float64 {
		return float64(len(page.Trigrams))
	})

	fmt.Printf("# of Trigrams:       %.3f (avg %.3f)\n", stdDev, avg)
}

func testTiming(iterations int, f func()) {
	t0 := time.Now()
	for i := 0; i < iterations; i++ {
		f()
	}
	t1 := time.Now()
	fmt.Printf("It took %v for %d iterations. (Avg. %.3fs)\n", t1.Sub(t0), iterations, t1.Sub(t0).Seconds()/float64(iterations))
	return
}

// Returns stdDeviation and Avg
func calcStats(containers []*PageContainer, f func(*PageContainer) float64) (float64, float64) {
	vals := make([]float64, len(containers))
	mean := float64(0)
	for i, page := range containers {
		if page.IsRedirect {
			i -= 1
			continue
		} else {
			vals[i] = f(page)
			mean += f(page)
		}
	}
	mean = mean / float64(len(vals))
	devSum := float64(0)
	for _, v := range vals {
		devSum += (float64(v) - mean) * (float64(v) - mean)
	}

	avg := float64(0)
	for _, v := range containers {
		avg += float64(f(v))
	}

	return math.Sqrt(devSum / float64(len(vals))), avg / float64(len(containers))
}

func calcStdDev(nums []float64) (float64, float64) {
	mean := float64(0)
	for _, v := range nums {
		mean += v
	}
	mean = mean / float64(len(nums))
	devSum := float64(0)
	for _, v := range nums {
		devSum += (float64(v) - mean) * (float64(v) - mean)
	}

	return math.Sqrt(devSum / float64(len(nums))), mean
}
