package main

import (
	"fmt"
	"github.com/glxxyz/dedupe/param"
	"github.com/glxxyz/dedupe/repo"
	"log"
	"os"
)

var errLog = log.New(os.Stderr, "", 0)

func main() {
	if options, err := param.ParseParameters(); err == nil && options != nil {
		if options.Verbose() {
			fmt.Printf("options: %+v\n", options)
		}
		var matchRepo repo.MatchRepository
		scanForDuplicates(options, &matchRepo)
	} else if err != nil {
		errLog.Print(err)
		os.Exit(1)
	}
}
