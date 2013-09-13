package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
//	"github.com/jcla1/matrix"
)

var _ = fmt.Println
var _ = os.Create

var dumpFileName string

var (
	featuredFilePath     string
	normalFilePath       string
	redirectFilePath     string
	nnFilePath           string
	ngramFilePath        string
	minMatrixFilePath    string
	rangeMatrixFilePath  string
	thresholdMatFilePath string
)

func init() {
	flag.StringVar(&dumpFileName, "dumpFile", "data/latest.xml.bz2", "the dump file to work with")

	flag.StringVar(&featuredFilePath, "featuredFile", "data/featured.gob.gzip", "place to store the featured pages")
	flag.StringVar(&normalFilePath, "normalFile", "data/normal.gob.gzip", "place to store the normal pages")
	flag.StringVar(&redirectFilePath, "redirectFile", "data/redirect.gob.gzip", "place to store the redirect map")

	flag.StringVar(&nnFilePath, "nnFile", "data/nn.gob", "place to put the trained NN")
	flag.StringVar(&minMatrixFilePath, "minFile", "data/min.gob", "place to put the min vector")
	flag.StringVar(&rangeMatrixFilePath, "rangeFile", "data/range.gob", "place to put the range vector")
	flag.StringVar(&thresholdMatFilePath, "thresholdMatrix", "data/thresholdMatrix.gob", "place to put the thresholdMatrix")

	flag.StringVar(&ngramFilePath, "ngramFile", "data/ngrams.txt", "place to put the ngrams")
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GOMAXPROCS(4)
	flag.Parse()

	min := LoadMatrixFromFile(minMatrixFilePath)
	r := LoadMatrixFromFile(rangeMatrixFilePath)

	thresholdMat := LoadMatrixFromFile(thresholdMatFilePath)
	thresholdMat = thresholdMat.Scale(0.8)

	n := LoadNNFromFile(nnFilePath)

	/*normalFile, err := os.Open(normalFilePath)
	if err != nil {
		panic(err)
	}
	defer normalFile.Close()

	normalPages := openCompressedPages(normalFile)*/

	file, err := os.Open(dumpFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	pages := distributorFromXML(file)

	featureChanIn := make(chan *PageContainer)
	processed := ProcessRecords(Normalizer(Vectorizer(featureChanIn), min, r), n, thresholdMat)

	for i := 0; i < 30; i++ {
		p := <-pages
		featureChanIn <- p
		fmt.Println(p.Page.Title)
		fmt.Println((<-processed).Vals)
	}

	/*testTiming(1, func() {
		Part1()
	})*/

	/*testTiming(1, func() {
		Part2NN()
	})*/

	/*min := LoadMatrixFromFile(minMatrixFilePath)
	r := LoadMatrixFromFile(rangeMatrixFilePath)

	normalFile, err := os.Open(normalFilePath)
	if err != nil {
		panic(err)
	}
	defer normalFile.Close()

	normalPages := openCompressedPages(normalFile)

	pages := Normalizer(Vectorizer(normalPages), min, r)

	fmt.Println(<-pages)
	fmt.Println()
	fmt.Println(<-pages)*/

	//n := LoadNNFromFile(nnFilePath)
	//fmt.Println(n.Thetas[1])

	/*file, err := os.Open(dumpFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	pages := distributorFromXML(file)

	pages = BuildPipeline(pages, pipelineFuncs)

	featureStats(pages, 10)*/
}
