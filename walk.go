package main

import (
	"fmt"
	"github.com/glxxyz/dedupe/repo"
	"os"
	"path/filepath"
	"sync"
)

func Walk(options repo.MatchOptions, root string, files chan<- *repo.FileData) {
	if err := filepath.Walk(root, walkFunc(options, files)); err != nil {
		panic(fmt.Errorf("error walking path %q: %v\n", root, err))
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
		} else if info.Mode()&os.ModeSymlink != 0 {
			walkSymLink(options, path, files)
		} else if info.Size() < options.MinBytes() {
			if options.Verbose() {
				fmt.Printf("ignoring file size %v bytes: %q\n", info.Size(), path)
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

func walkSymLink(options repo.MatchOptions, path string, files chan<- *repo.FileData) {
	if options.SymLinks() {
		dest, err := filepath.EvalSymlinks(path)
		if err != nil {
			errLog.Printf("failed to evaluate symbolic link %q: %v\n", path, err)
		} else {
			if options.Verbose() {
				fmt.Printf("following symbolic link: %q to %q\n", path, dest)
			}
			Walk(options, dest, files)
		}
	} else if options.Verbose() {
		fmt.Printf("ignoring symbolic link: %q\n", path)
	}
}

var visitedMap sync.Map

func visited(path string) bool {
	_, loaded := visitedMap.LoadOrStore(path, true)
	return loaded
}
