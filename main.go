package main

import (
	"flag"
	"fmt"
	"os"
)

var dumpFileName = flag.String("dumpFile", "data/latest.xml.bz2", "the dump file to work with")

func main() {
	flag.Parse()

	file, err := os.Open(*dumpFileName)
	if err != nil {
		panic(err)
	}

	pages := distributor(file)

	fmt.Println((<-pages).Page.Text())
}
