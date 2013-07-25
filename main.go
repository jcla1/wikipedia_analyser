package main

import (
	"compress/gzip"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
)

var _ = gzip.NewReader
var _ = fmt.Println

var dumpFileName = flag.String("dumpFile", "data/latest.xml.bz2", "the dump file to work with")
var N = flag.Int("N", 100, "number of pages to read for timing")

func WriteMap(mapChannel <-chan map[string]int, w io.Writer) <-chan bool {
	c := make(chan bool)
	encoder := gob.NewEncoder(w)

	go func() {
		m := <-mapChannel
		err := encoder.Encode(m)
		if err != nil {
			panic(err)
		}

		c <- true
	}()

	return c
}

func ReadMap(r io.Reader) (map[string]int, error) {
	var m map[string]int

	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func DistrbuteArticles(input <-chan *PageContainer, featuredChannel, redirectChannel, normalChannel chan<- *PageContainer) <-chan bool {
	c := make(chan bool)

	go func() {
		var container *PageContainer
		var ok bool

		for {
			container, ok = <-input

			if !ok {
				// channel was closed!
				close(redirectChannel)
				close(featuredChannel)
				close(normalChannel)
				c <- true
				return
			} else {

				if container.IsRedirect {
					redirectChannel <- container
				} else if container.IsFeatured {
					featuredChannel <- container
				} else {
					normalChannel <- container
				}

			}
		}
	}()

	return c
}

func ArticleWriter(w io.Writer) (chan<- *PageContainer, <-chan bool) {
	c := make(chan *PageContainer)
	finishChannel := make(chan bool)
	encoder := gob.NewEncoder(w)

	go func() {
		var container *PageContainer
		var ok bool
		var err error

		for {
			container, ok = <-c

			if !ok {
				// channel was closed!
				finishChannel <- true
				return
			} else {
				err = encoder.Encode(container)
				if err != nil {
					panic(err)
				}
			}
		}
	}()

	return c, finishChannel
}

func BuildRedirectMap() (chan<- *PageContainer, <-chan map[string]int) {
	mapChannel := make(chan map[string]int)
	input := make(chan *PageContainer)
	m := make(map[string]int, 1000)

	go func() {
		var page *PageContainer
		var ok bool

		for {
			page, ok = <-input

			if !ok {
				// channel was closed!
				mapChannel <- m
				return
			} else {
				m[page.Page.Redirect.Title] += 1
			}
		}
	}()

	return input, mapChannel
}

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
