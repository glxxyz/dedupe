package main

import (
	"./param"
	"fmt"
	"log"
	"os"
)

var errLog = log.New(os.Stderr, "", 0)

func main() {
	options := param.ParseParameters()
	if options.Verbose() {
		fmt.Printf("options: %+v\n", options)
	}
	scanForDuplicates(options)
}

