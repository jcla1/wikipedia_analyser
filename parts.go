package main

import (
	"./ngram"
	"os"
)

func Part1() {
	file, err := os.Open(dumpFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	pages := distributorFromXML(file)

	pages = BuildPipeline(pages, pipelineFuncs)
	saveSplitPages(pages)
}

func Part2() {
	Part2NN()
	Part2LM()
}

func Part2NN() {
	// Learning NN
	featuredFile, err := os.Open(featuredFilePath)
	if err != nil {
		panic(err)
	}
	defer featuredFile.Close()

	featuredPages := openCompressedPages(featuredFile)

	n := BuildTrainNN(featuredPages)
	SaveNNToFile(n, nnFilePath)
}

func Part2LM() {
	// Building Language Models
	featuredFile, err := os.Open(featuredFilePath)
	if err != nil {
		panic(err)
	}
	featuredPages := openCompressedPages(featuredFile)

	// I had to use all of the pages otherwise some of
	// the probabilities would have been +Inf! May still
	// give good measure of likelyhood.

	normalFile, err := os.Open(normalFilePath)
	if err != nil {
		panic(err)
	}
	normalPages := openCompressedPages(normalFile)

	unigramFile, err := os.Create(unigramFilePath)
	if err != nil {
		panic(err)
	}
	defer unigramFile.Close()

	bigramFile, err := os.Create(bigramFilePath)
	if err != nil {
		panic(err)
	}
	defer bigramFile.Close()

	trigramFile, err := os.Create(trigramFilePath)
	if err != nil {
		panic(err)
	}
	defer trigramFile.Close()

	unigramChannel := make(chan map[string]int)
	unigramBackChannel := ngram.BuildLanguageModel(unigramChannel)

	bigramChannel := make(chan map[string]int)
	bigramBackChannel := ngram.BuildLanguageModel(bigramChannel)

	trigramChannel := make(chan map[string]int)
	trigramBackChannel := ngram.BuildLanguageModel(trigramChannel)

	for container := range featuredPages {
		unigramChannel <- container.Unigrams
		bigramChannel <- container.Bigrams
		trigramChannel <- container.Trigrams
	}

	for container := range normalPages {
		unigramChannel <- container.Unigrams
		bigramChannel <- container.Bigrams
		trigramChannel <- container.Trigrams
	}

	close(unigramChannel)
	close(bigramChannel)
	close(trigramChannel)

	saveLanguageModel(unigramBackChannel, unigramFile)
	saveLanguageModel(bigramBackChannel, bigramFile)
	saveLanguageModel(trigramBackChannel, trigramFile)
}
