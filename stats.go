package main

import (
	"fmt"
	"math"
	"time"
)

func featureStats(input <-chan *PageContainer, n int) {
	fmt.Println("Standard Deviation of features:")

	fmt.Println("# of Links:", calcStdDeviation(input, n, func(page *PageContainer) float64 {
		return float64(page.NumLinks)
	}))

	fmt.Println("Link Density:", calcStdDeviation(input, n, func(page *PageContainer) float64 {
		return page.LinkDensity
	}))

	fmt.Println("# of External Links:", calcStdDeviation(input, n, func(page *PageContainer) float64 {
		return float64(page.NumExternalLinks)
	}))

	fmt.Println("# of Headings:", calcStdDeviation(input, n, func(page *PageContainer) float64 {
		return float64(page.NumHeadings)
	}))

	fmt.Println("# of Refrences:", calcStdDeviation(input, n, func(page *PageContainer) float64 {
		return float64(page.NumRefs)
	}))

	fmt.Println("# of Categories:", calcStdDeviation(input, n, func(page *PageContainer) float64 {
		return float64(page.NumCategories)
	}))

	fmt.Println("# of Sentences:", calcStdDeviation(input, n, func(page *PageContainer) float64 {
		return float64(page.NumSentences)
	}))

	fmt.Println("Avg Sentence Length:", calcStdDeviation(input, n, func(page *PageContainer) float64 {
		return page.AvgSentenceLen
	}))

	fmt.Println("# of Words:", calcStdDeviation(input, n, func(page *PageContainer) float64 {
		return float64(page.NumWords)
	}))

	fmt.Println("Avg Word Length:", calcStdDeviation(input, n, func(page *PageContainer) float64 {
		return page.AvgWordLen
	}))

	fmt.Println("# of Unigrams:", calcStdDeviation(input, n, func(page *PageContainer) float64 {
		return float64(len(page.Unigrams))
	}))

	fmt.Println("# of Bigrams:", calcStdDeviation(input, n, func(page *PageContainer) float64 {
		return float64(len(page.Bigrams))
	}))

	fmt.Println("# of Trigrams:", calcStdDeviation(input, n, func(page *PageContainer) float64 {
		return float64(len(page.Trigrams))
	}))
}

func testTiming(iterations int, f func()) {
	t0 := time.Now()
	for i := 0; i < iterations; i++ {
		f()
	}
	t1 := time.Now()
	fmt.Printf("It took %v for %d iterations. (Avg. %fs)\n", t1.Sub(t0), iterations, t1.Sub(t0).Seconds()/float64(iterations))
	return
}

func calcStdDeviation(input <-chan *PageContainer, n int, f func(*PageContainer) float64) float64 {
	vals := make([]float64, n)
	mean := float64(0)
	for i := 0; i < len(vals); i++ {
		page := <-input
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
	return math.Sqrt(devSum / float64(len(vals)))
}
