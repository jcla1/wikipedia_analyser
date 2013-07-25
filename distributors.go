package main

import (
	"compress/bzip2"
	"encoding/gob"
	"io"
)

func distributorFromGob(r io.Reader) <-chan *PageContainer {
	decoder := gob.NewDecoder(r)
	channel := make(chan *PageContainer, 50)

	go func() {
		var v *PageContainer
		var err error

		for {
			err = decoder.Decode(&v)

			if err == io.EOF {
				close(channel)
				return
			} else if err != nil {
				panic(err)
			}

			channel <- v
		}
	}()

	return channel
}

func distributorFromXML(r io.Reader) <-chan *PageContainer {
	return containerDistributor(xmlPageDistributor(r))
}

// Wraps our pages in a struct that can hold features
func containerDistributor(input <-chan *Page) <-chan *PageContainer {
	output := make(chan *PageContainer, 50)
	go func() {
		var page *Page
		var ok bool
		for {
			page, ok = <-input

			if !ok {
				// channel was closed!
				close(output)
				return
			} else {
				output <- &PageContainer{Page: page}
			}
		}
	}()

	return output
}

// Takes an io.Reader to read pages from
func xmlPageDistributor(r io.Reader) <-chan *Page {
	decompressedStream := bzip2.NewReader(r)

	parser, err := NewParser(decompressedStream)
	if err != nil {
		panic(err)
	}

	channel := make(chan *Page)
	go func() {
		var page *Page
		var err error

		for {
			page, err = parser.Next()

			if err == io.EOF {
				close(channel)
				return
			} else if err != nil {
				panic(err)
			}

			channel <- page
		}
	}()

	return channel
}
