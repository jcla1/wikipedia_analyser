package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

//	"time"
)

var dumpFileName = flag.String("dumpFile", "data/latest.xml.bz2", "the dump file to work with")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	file, err := os.Open(*dumpFileName)
	if err != nil {
		panic(err)
	}

	pages := distributor(file)
	pages = BuildPipeline(pages, pipelineFuncs)

	/*t0 := time.Now()
	for i := 0; i < 100; i++ {
		<-pages
	}
	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))*/

	//fmt.Println((<-pages).Page.Text())
	<-pages
	fmt.Println((<-pages).PlainText)
}
