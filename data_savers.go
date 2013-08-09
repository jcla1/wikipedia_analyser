package main

import (
	"compress/gzip"
	"encoding/gob"
	"io"
	"os"
	"github.com/jcla1/matrix"
)

func LoadMatrixFromFile(path string) *matrix.Matrix {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return LoadMatrix(file)
}

func LoadMatrix(r io.Reader) *matrix.Matrix {
	var m *matrix.Matrix
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&m)
	if err != nil {
		panic(err)
	}
	return m
}

func SaveMatrixToFile(m *matrix.Matrix, path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	SaveMatrix(m, file)
}

func SaveMatrix(m *matrix.Matrix, w io.Writer) {
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(m)
	if err != nil {
		panic(err)
	}
}

func BuildLMModule(in <-chan *PageContainer) <-chan *PageContainer {
	out := make(chan *PageContainer)

	ngramFile, err := os.Create(ngramFilePath)
	if err != nil {
		panic(err)
	}

	newLine := []byte("\n")

	go func() {
		var c *PageContainer
		var ok bool

		for {
			c, ok = <-in

			if !ok {
				close(out)
				return
			} else {
				if !c.IsRedirect {
					for k, v := range c.Unigrams {
						for i := 0; i < v; i++ {
							ngramFile.Write([]byte(k))
							ngramFile.Write(newLine)
						}
					}
					for k, v := range c.Bigrams {
						for i := 0; i < v; i++ {
							ngramFile.Write([]byte(k))
							ngramFile.Write(newLine)
						}
					}
					for k, v := range c.Trigrams {
						for i := 0; i < v; i++ {
							ngramFile.Write([]byte(k))
							ngramFile.Write(newLine)
						}
					}
				}

				out <- c
			}
		}
	}()

	return out
}

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

	featuredChannel, featuredWriter := ArticleWriter(featuredCompressed)
	normalChannel, normalWriter := ArticleWriter(normalCompressed)
	//featuredChannel, featuredWriter := ArticleWriter(featuredFile)
	//normalChannel, normalWriter := ArticleWriter(normalFile)

	DistrbuteArticles(input, featuredChannel, normalChannel)

	// Wait for all writers to finish
	<-featuredWriter
	<-normalWriter

	// Close all the gzip streams
	// I tried defering the close
	// call but kept getting some EOF errors
	featuredCompressed.Close()
	normalCompressed.Close()
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
	c := make(chan bool)

	go func() {
		var container *PageContainer
		var ok bool

		for {
			container, ok = <-input

			if !ok {
				// channel was closed!
				close(featuredChannel)
				close(normalChannel)
				c <- true
				return
			} else {

				if container.IsFeatured {
					featuredChannel <- container
				} else if !container.IsRedirect {
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
