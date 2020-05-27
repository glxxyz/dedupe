package main

import (
	"fmt"
	"github.com/glxxyz/dedupe/param"
	"github.com/glxxyz/dedupe/repo"
	"sync"
	"sync/atomic"
	"time"
)

func scanForDuplicates(options *param.Options, matchRepo *repo.MatchRepository) {
	var scanners sync.WaitGroup
	var matchers sync.WaitGroup
	var movers sync.WaitGroup

	var scans = make(chan string, options.ScanBuffer())
	var files = make(chan *repo.FileData, options.MatchBuffer())
	var moves = make(chan string, options.MoveBuffer())

	var scanCount uint32
	var fileCount uint32
	var moveCount uint32

	spawnScanners(options, &scanners, scans, files, &fileCount)
	spawnMatchers(options, matchRepo, &matchers, files, moves, &moveCount)
	spawnMovers(options, &movers, moves)
	seedScanners(options, scans, &scanCount)

	spawnChannelTicker(options, scans, files, moves, &scanCount, &fileCount, &moveCount)

	close(scans)
	scanners.Wait()
	close(files)
	matchers.Wait()
	close(moves)
	movers.Wait()
}

func spawnChannelTicker(
	options *param.Options,
	scans chan string, files chan *repo.FileData, moves chan string,
	scanCount *uint32, fileCount *uint32, moveCount *uint32) {
	if options.Verbose() {
		ticker := time.NewTicker(time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					fmt.Printf(
						"Channels:\tlen/cap/count\tscans=%d/%d/%d\tfiles=%d/%d/%d\tmoves=%d/%d/%d\n",
						len(scans), cap(scans), *scanCount, len(files), cap(files), *fileCount, len(moves), cap(moves), *moveCount)
				}
			}
		}()
	}
}

func seedScanners(options *param.Options, scans chan<- string, scanCount *uint32) {
	for _, path := range options.Paths() {
		scans <- path
		atomic.AddUint32(scanCount, 1)
	}
}

func spawnScanners(options *param.Options, scanners *sync.WaitGroup, scans <-chan string, files chan<- *repo.FileData, fileCount *uint32) {
	for i := 0; i < options.Scanners(); i++ {
		scanners.Add(1)
		go func(num int) {
			defer scanners.Done()
			scanWorker(num, options, scans, files, fileCount)
		}(i)
	}
}

func scanWorker(num int, options *param.Options, scans <-chan string, files chan<- *repo.FileData, fileCount *uint32) {
	if options.Verbose() {
		fmt.Printf("scanner %d starting\n", num)
	}
	for {
		path := <-scans
		if path != "" {
			Walk(options, path, files, fileCount)
		} else {
			if options.Verbose() {
				fmt.Printf("scanner %d done\n", num)
			}
			break
		}
	}
}

func spawnMatchers(options *param.Options, matchRepo *repo.MatchRepository, matchers *sync.WaitGroup, files <-chan *repo.FileData, moves chan<- string, moveCount *uint32) {
	for i := 0; i < options.Matchers(); i++ {
		matchers.Add(1)
		go func(num int) {
			defer matchers.Done()
			matchWorker(num, options, matchRepo, files, moves, moveCount)
		}(i)
	}
}

func matchWorker(num int, options *param.Options, matchRepo *repo.MatchRepository, files <-chan *repo.FileData, moves chan<- string, moveCount *uint32) {
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
				atomic.AddUint32(moveCount, 1)
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
