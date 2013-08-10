package main

import (
	"fmt"
	"os"
)

var _ = fmt.Println

func Part1() {
	file, err := os.Open(dumpFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	pages := distributorFromXML(file)

	pages = BuildPipeline(pages, pipelineFuncs)
	pages = BuildLMModule(pages)
	saveSplitPages(pages)
}

func Part2() {

	Part2Normalize()
	Part2NN()
	Part2MaxPredictionErr()
}

func Part2Normalize() {
	featuredFile, err := os.Open(featuredFilePath)
	if err != nil {
		panic(err)
	}
	defer featuredFile.Close()

	featuredPages := openCompressedPages(featuredFile)

	min, r := FindNormalizedVectors(Vectorizer(featuredPages))

	fmt.Println(min)
	fmt.Println(r)

	SaveMatrixToFile(min, minMatrixFilePath)
	SaveMatrixToFile(r, rangeMatrixFilePath)
}

func Part2NN() {
	// Learning NN
	featuredFile, err := os.Open(featuredFilePath)
	if err != nil {
		panic(err)
	}
	defer featuredFile.Close()

	featuredPages := openCompressedPages(featuredFile)

	min := LoadMatrixFromFile(minMatrixFilePath)
	r := LoadMatrixFromFile(rangeMatrixFilePath)

	n := BuildTrainNN(featuredPages, min, r)
	SaveNNToFile(n, nnFilePath)
}

func Part2MaxPredictionErr() {
	featuredFile, err := os.Open(featuredFilePath)
	if err != nil {
		panic(err)
	}
	defer featuredFile.Close()

	input := openCompressedPages(featuredFile)
	min := LoadMatrixFromFile(minMatrixFilePath)
	r := LoadMatrixFromFile(rangeMatrixFilePath)
	n := LoadNNFromFile(nnFilePath)

	devMax, means := FindPredictionErrDeviation(input, n, min, r)

	thresholdMat, _ := devMax.Add(means)

	SaveMatrixToFile(thresholdMat, thresholdMatFilePath)
}
