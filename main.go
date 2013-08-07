package main

import (
	"flag"
	"fmt"
	"runtime"
	"os"
)

var _ = fmt.Println

var dumpFileName string

var (
	featuredFilePath string
	normalFilePath   string
	redirectFilePath string
	nnFilePath       string
	unigramFilePath  string
	bigramFilePath   string
	trigramFilePath  string
)

func init() {
	flag.StringVar(&dumpFileName, "dumpFile", "data/latest.xml.bz2", "the dump file to work with")

	flag.StringVar(&featuredFilePath, "featuredFile", "data/featured.gob.gzip", "place to store the featured pages")
	flag.StringVar(&normalFilePath, "normalFile", "data/normal.gob.gzip", "place to store the normal pages")
	//flag.StringVar(&redirectFilePath, "redirectFile", "data/redirect.txt", "place to store the redirect map")

	flag.StringVar(&nnFilePath, "nnFile", "data/nn.gob", "place to put the trained NN")

	flag.StringVar(&unigramFilePath, "unigramFile", "data/unigrams.gob", "place to put the Unigrams")
	flag.StringVar(&bigramFilePath, "bigramFile", "data/bigrams.gob", "place to put the Bigrams")
	flag.StringVar(&trigramFilePath, "trigramFile", "data/trigrams.gob", "place to put the Trigrams")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	//Part1()

	file, err := os.Open(dumpFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	pages := distributorFromXML(file)

	//pages = BuildPipeline(pages, pipelineFuncs)
	for p := range pages {
		fmt.Println(p)
	}
}
