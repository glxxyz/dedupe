package main

import (
	"./repo"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func Walk(options repo.MatchOptions, root string, files chan<- *repo.FileData) {
	if err := filepath.Walk(root, walkFunc(options, files)); err != nil {
		errLog.Printf("error walking path %q: %v\n", root, err)
		panic(err)
	}
}

func walkFunc(options repo.MatchOptions, files chan<- *repo.FileData) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errLog.Printf("failed to access path %q: %v\n", path, err)
			return nil
		}

		if visited(path) {
			if options.Verbose() {
				fmt.Printf("already visited: %q\n", path)
			}
			return filepath.SkipDir
		} else if info.IsDir() {
			if options.Verbose() {
				fmt.Printf("visiting dir: %q\n", path)
			}
		} else if info.Mode() & os.ModeSymlink != 0{
			if options.Verbose() {
				fmt.Printf("ignoring symbolic link: %q\n", path)
			}
		} else if !options.IgnoreZero() && info.Size() == 0 {
			if options.Verbose() {
				fmt.Printf("ignoring zero size file: %q\n", path)
			}
		} else {
			if options.Verbose() {
				fmt.Printf("visiting file: %q\n", path)
			}
			files <- repo.NewFile(options, path, info)
		}
		return nil
	}
}

var visitedMap sync.Map

func visited(path string) bool {
	_, loaded := visitedMap.LoadOrStore(path, true)
	return loaded
}
