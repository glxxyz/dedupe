package main

import (
	"fmt"
	"github.com/glxxyz/dedupe/param"
	"github.com/glxxyz/dedupe/repo"
	"sync"
	"time"
)

func scanForDuplicates(options *param.Options, matchRepo *repo.MatchRepository) {
	var scanners sync.WaitGroup
	var matchers sync.WaitGroup
	var movers sync.WaitGroup

	var scans = make(chan string, options.ScanBuffer())
	var files = make(chan *repo.FileData, options.MatchBuffer())
	var moves = make(chan string, options.MoveBuffer())

	spawnScanners(options, &scanners, scans, files)
	spawnMatchers(options, matchRepo, &matchers, files, moves)
	spawnMovers(options, &movers, moves)
	seedScanners(options, scans)

	spawnChannelTicker(options, scans, files, moves)

	close(scans)
	scanners.Wait()
	close(files)
	matchers.Wait()
	close(moves)
	movers.Wait()
}

func spawnChannelTicker(options *param.Options, scans chan string, files chan *repo.FileData, moves chan string) {
	if options.Verbose() {
		ticker := time.NewTicker(time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					fmt.Printf("Channels:\tscans=%d/%d\tfiles=%d/%d\tmoves=%d/%d\n",
						len(scans), cap(scans), len(files), cap(files), len(moves), cap(moves))
				}
			}
		}()
	}
}

func seedScanners(options *param.Options, scans chan<- string) {
	for _, path := range options.Paths() {
		scans <- path
	}
}

func spawnScanners(options *param.Options, scanners *sync.WaitGroup, scans <-chan string, files chan<- *repo.FileData) {
	for i := 0; i < options.Scanners(); i++ {
		scanners.Add(1)
		go func(num int) {
			defer scanners.Done()
			scanWorker(num, options, scans, files)
		}(i)
	}
}

func scanWorker(num int, options *param.Options, scans <-chan string, files chan<- *repo.FileData) {
	if options.Verbose() {
		fmt.Printf("scanner %d starting\n", num)
	}
	for {
		path := <-scans
		if path != "" {
			Walk(options, path, files)
		} else {
			if options.Verbose() {
				fmt.Printf("scanner %d done\n", num)
			}
			break
		}
	}
}

func spawnMatchers(options *param.Options, matchRepo *repo.MatchRepository, matchers *sync.WaitGroup, files <-chan *repo.FileData, moves chan<- string) {
	for i := 0; i < options.Matchers(); i++ {
		matchers.Add(1)
		go func(num int) {
			defer matchers.Done()
			matchWorker(num, options, matchRepo, files, moves)
		}(i)
	}
}

func matchWorker(num int, options *param.Options, matchRepo *repo.MatchRepository, files <-chan *repo.FileData, moves chan<- string) {
	if options.Verbose() {
		fmt.Printf("matcher %d starting\n", num)
	}
	for {
		file := <-files
		if file != nil {
			if options.Verbose() {
				fmt.Printf("matcher %d working on file: %v\n", num, file)
			}
			if fileToMove, found := matchRepo.MatchFileToMove(options, file); found {
				moves <- fileToMove
			}
		} else {
			if options.Verbose() {
				fmt.Printf("matcher %d done\n", num)
			}
			break
		}
	}
}

func spawnMovers(options *param.Options, movers *sync.WaitGroup, moves <-chan string) {
	for i := 0; i < options.Matchers(); i++ {
		movers.Add(1)
		go func(num int) {
			defer movers.Done()
			moveWorker(num, options, moves)
		}(i)
	}
}

func moveWorker(num int, options *param.Options, moves <-chan string) {
	if options.Verbose() {
		fmt.Printf("mover %d starting\n", num)
	}
	for {
		filePath := <-moves
		if filePath != "" {
			Move(options, filePath)
		} else {
			if options.Verbose() {
				fmt.Printf("mover %d done\n", num)
			}
			break
		}
	}
}
