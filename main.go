package main

import (
	"compress/gzip"
	"encoding/gob"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"time"
	"io"
)

var dumpFileName = flag.String("dumpFile", "data/latest.xml.bz2", "the dump file to work with")
var N = flag.Int("N", 100, "number of pages to read for timing")

func testTiming(pages <-chan *PageContainer) {
	n := *N
	t0 := time.Now()
	for i := 0; i < n; i++ {
		<-pages
	}
	t1 := time.Now()
	fmt.Printf("It took %v for %d articles. (Avg. %fs)\n", t1.Sub(t0), n, t1.Sub(t0).Seconds()/float64(n))
	return
}

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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	file, err := os.Open(*dumpFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	pages := distributor(file)

	pages = BuildPipeline(pages, pipelineFuncs)

	featuredFile, err := os.Create("data/featured.gob")
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

	redirectChannel, mapChannel := BuildRedirectMap()
	mapWriter := WriteMap(mapChannel, redirectCompressed)

	c := make(chan *PageContainer, 10)

	for i := 0; i < 10; i++ {
		c <- <-pages
	}

	DistrbuteArticles(c, featuredChannel, redirectChannel, normalChannel)
	close(c)
	time.Sleep(10*time.Second)
	<-mapWriter
	<-featuredWriter
	<-normalWriter

/*
		file, err := os.Open("data/normal.gob")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		r, err := gzip.NewReader(file)
		if err != nil {
			panic(err)
		}

		var v *PageContainer
		decoder := gob.NewDecoder(r)
		err = decoder.Decode(&v)
		if err != nil {
			panic(err)
		}
		fmt.Println(v)
*/
	//testTiming(pages)
	/*
		N := 1000

		fmt.Println("AvgSentenceLen:", calcStdDeviation(pages, N, func(page *PageContainer) float64 {
			return page.AvgSentenceLen
		}))
	*/
}
