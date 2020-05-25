package repo

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var errLog = log.New(os.Stderr, "", 0)

type MatchOptions interface {
	ModTime() bool
	Name() bool
	Size() bool
	Hash() bool
	Contents() bool
	MinBytes() int64
	SymLinks() bool
	Verbose() bool
	Paths() []string
}

type primaryKey struct {
	name    string
	modTime time.Time
	size    int64
}

var primaryMap sync.Map // primaryKey -> *matchAttributes

func MatchFileToMove(options MatchOptions, file *FileData) (string, bool) {
	if primary, found := findPrimaryMatch(file); found {
		if options.Verbose() {
			fmt.Printf("attributes match found for: %q\n", file.FilePath)
		}
		if shortHash, found := primary.findHeadMatch(options, file.FilePath); found {
			if options.Verbose() {
				fmt.Printf("head hash match found for: %q\n", file.FilePath)
			}
			if fullHash, found := shortHash.findFullMatch(options, file.FilePath); found {
				if options.Verbose() {
					fmt.Printf("full hash match found for: %q\n", file.FilePath)
				}
				return fullHash.lowestPriorityMatch(options, file.FilePath)
			}
		}
	}
	return "", false
}

func findPrimaryMatch(file *FileData) (*matchAttributes, bool) {
	actual, loaded := primaryMap.LoadOrStore(
		primaryKey{name: file.Name, modTime: file.modTime, size: file.size},
		&matchAttributes{singlePath: file.FilePath})
	return actual.(*matchAttributes), loaded
}
