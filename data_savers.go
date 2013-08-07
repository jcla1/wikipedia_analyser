package main

import (
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func openCompressedPages(file io.Reader) <-chan *PageContainer {
	decompressedStream, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}

	return distributorFromGob(decompressedStream)
}

func saveSplitPages(input <-chan *PageContainer) {
	featuredFile, err := os.Create(featuredFilePath)
	if err != nil {
		panic(err)
	}
	featuredCompressed := gzip.NewWriter(featuredFile)
	defer featuredFile.Close()

	normalFile, err := os.Create(normalFilePath)
	if err != nil {
		panic(err)
	}
	normalCompressed := gzip.NewWriter(normalFile)
	defer normalFile.Close()

	/*redirectFile, err := os.Create(redirectFilePath)
	if err != nil {
		panic(err)
	}
	//redirectCompressed := gzip.NewWriter(redirectFile)
	defer redirectFile.Close()*/

	featuredChannel, featuredWriter := ArticleWriter(featuredCompressed)
	normalChannel, normalWriter := ArticleWriter(normalCompressed)
	//featuredChannel, featuredWriter := ArticleWriter(featuredFile)
	//normalChannel, normalWriter := ArticleWriter(normalFile)

	/*
		redirectChannel, mapChannel := BuildRedirectMap()
		mapWriter := WriteMap(mapChannel, redirectCompressed)
		//mapWriter := WriteMap(mapChannel, redirectFile)
	*/

	//redirectChannel := RedirectTitleWriter(redirectFile)

	//DistrbuteArticles(input, featuredChannel, redirectChannel, normalChannel)
	DistrbuteArticles(input, featuredChannel, normalChannel)

	// Wait for all writers to finish
	//<-mapWriter
	<-featuredWriter
	<-normalWriter

	// Close all the gzip streams
	// I tried defering the close
	// call but kept getting some EOF errors
	featuredCompressed.Close()
	normalCompressed.Close()
	//redirectCompressed.Close()
}

func RedirectTitleWriter(w io.Writer) chan<- *PageContainer {
	input := make(chan *PageContainer)

	go func() {
		var container *PageContainer
		var ok bool

		for {
			container, ok = <-input

			if !ok {
				return
			} else {
				w.Write([]byte(fmt.Sprintf("%s\n", container.Page.Title)))
			}
		}
	}()

	return input
}

func splitDumpFile(input <-chan *PageContainer) {
	c := 0
	chunkSize := 10000

	var page *PageContainer
	var err error
	var chunkFile *os.File
	var encoder *json.Encoder
	var ok bool

	for {
		c += 1

		chunkFile, err = os.Create(fmt.Sprintf("data/parts/dumpFile-%04d.json", c))
		//chunkFile = os.Stdout
		if err != nil {
			panic(err)
		}

		encoder = json.NewEncoder(chunkFile)

		for i := 0; i < chunkSize; i++ {
			page, ok = <-input
			if !ok {
				chunkFile.Close()
				return
			} else {
				err = encoder.Encode(page)
				if err != nil {
					fmt.Printf("%#v\n", page)
					panic(err)
				}
			}
		}

		chunkFile.Close()
	}
}

func saveLanguageModel(input <-chan map[string]int, lmFile io.Writer) {
	mapWriter := WriteMap(input, lmFile)
	<-mapWriter
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

func DistrbuteArticles(input <-chan *PageContainer, featuredChannel, normalChannel chan<- *PageContainer) <-chan bool {
//func DistrbuteArticles(input <-chan *PageContainer, featuredChannel, redirectChannel, normalChannel chan<- *PageContainer) <-chan bool {
	c := make(chan bool)

	go func() {
		var container *PageContainer
		var ok bool

		for {
			container, ok = <-input

			if !ok {
				// channel was closed!
				//close(redirectChannel)
				close(featuredChannel)
				close(normalChannel)
				c <- true
				return
			} else {

				if container.IsRedirect {
					//redirectChannel <- container
					//<- container
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
