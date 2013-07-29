package main

import (
	"./ngram"
	"flag"
	"fmt"
	"os"
	"runtime"
)

var _ = fmt.Println

var dumpFileName string

var (
	featuredFilePath string
	normalFilePath   string
	redirectFilePath string
)

func init() {
	flag.StringVar(&dumpFileName, "dumpFile", "data/latest.xml.bz2", "the dump file to work with")

	flag.StringVar(&featuredFilePath, "featuredFile", "data/featured.gob.gzip", "place to store the featured pages")
	flag.StringVar(&normalFilePath, "normalFile", "data/normal.gob.gzip", "place to store the normal pages")
	flag.StringVar(&redirectFilePath, "redirectFile", "data/redirect.gob.gzip", "place to store the redirect map")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	file, err := os.Open(dumpFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	pages := distributorFromXML(file)

	pages = BuildPipeline(pages, pipelineFuncs)

	c := make(chan ngram.LanguageModel)

	go func() {
		for i := 0; i < 10; i++ {
			c <- (<-pages).Unigrams
		}
		close(c)
	}()

	lmChannel := ngram.BuildLanguageModel(c)
	fmt.Println(len(<-lmChannel))

	/*file, err := os.Open(path)
	  if err != nil {
	    panic(err)
	  }
	  defer file.Close()

	  featuredPage := openCompressedPages(file)
		for fa := range featuredPages {
			fmt.Println(fa)
		}

		//featureStats(featuredPages, 5)*/
}
