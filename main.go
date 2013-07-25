package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"os"
	"runtime"
)

var _ = gzip.NewReader
var _ = fmt.Println

var dumpFileName = flag.String("dumpFile", "data/latest.xml.bz2", "the dump file to work with")
var N = flag.Int("N", 100, "number of pages to read for timing")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	file, err := os.Open(*dumpFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	pages := distributorFromXML(file)

	pages = BuildPipeline(pages, pipelineFuncs)

	//testTiming(1, func() {
	//	featureStats(pages, 100)
	//})
	<-pages
	fmt.Println(<-pages)

	/*featuredFile, err := os.Create("data/featured.gob")
	if err != nil {
		panic(err)
	}
	featuredCompressed := gzip.NewWriter(featuredFile)
	defer featuredCompressed.Close()
	defer featuredFile.Close()

	normalFile, err := os.Create("data/normal.gob")
	if err != nil {
		panic(err)
	}
	normalCompressed := gzip.NewWriter(normalFile)
	defer normalCompressed.Close()
	defer normalFile.Close()

	redirectFile, err := os.Create("data/redirect.gob")
	if err != nil {
		panic(err)
	}
	redirectCompressed := gzip.NewWriter(redirectFile)
	defer redirectCompressed.Close()
	defer featuredFile.Close()

	featuredChannel, featuredWriter := ArticleWriter(featuredCompressed)
	normalChannel, normalWriter := ArticleWriter(normalCompressed)
	//featuredChannel, featuredWriter := ArticleWriter(featuredFile)
	//normalChannel, normalWriter := ArticleWriter(normalFile)

	redirectChannel, mapChannel := BuildRedirectMap()
	mapWriter := WriteMap(mapChannel, redirectCompressed)
	//mapWriter := WriteMap(mapChannel, redirectFile)

	c := make(chan *PageContainer, 10)

	go func() {
		for i := 0; i < 500; i++ {
			c <- <-pages
		}
		close(c)
	}()

	testTiming(1, func() {
		DistrbuteArticles(c, featuredChannel, redirectChannel, normalChannel)

		// Wait for all writers to finish
		<-mapWriter
		<-featuredWriter
		<-normalWriter

		// Close all the gzip streams
		// I tried defering the close
		// call but kept getting some EOF errors
		featuredCompressed.Close()
		normalCompressed.Close()
		redirectCompressed.Close()
	})*/

	/*func() {
		file, err := os.Open("data/featured.gob")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		decompressedStream, err := gzip.NewReader(file)
		if err != nil {
			panic(err)
		}

		featuredPages := distributorFromGob(decompressedStream)

		for fa := range featuredPages {
			fmt.Println(fa)
		}

		//featureStats(featuredPages, 5)
	}()*/

	/*testTiming(*N, func() {
		<-pages
	})*/
}
