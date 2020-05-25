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

type MatchRepository struct {
	primaryMap sync.Map // primaryKey -> *matchAttributes
}

func (matchRepo *MatchRepository) MatchFileToMove(options MatchOptions, file *FileData) (string, bool) {
	if primary, found := matchRepo.findPrimaryMatch(file); found {
		if options.Verbose() {
			fmt.Printf("attributes match found for: %q\n", file.filePath)
		}
		if shortHash, found := primary.findHeadMatch(options, file.filePath); found {
			if options.Verbose() {
				fmt.Printf("head hash match found for: %q\n", file.filePath)
			}
			if fullHash, found := shortHash.findFullMatch(options, file.filePath); found {
				if options.Verbose() {
					fmt.Printf("full hash match found for: %q\n", file.filePath)
				}
				return fullHash.lowestPriorityMatch(options, file.filePath)
			}
		}
	}
	return "", false
}

func (matchRepo *MatchRepository) findPrimaryMatch(file *FileData) (*matchAttributes, bool) {
	actual, loaded := matchRepo.primaryMap.LoadOrStore(
		primaryKey{name: file.name, modTime: file.modTime, size: file.size},
		&matchAttributes{singlePath: file.filePath})
	return actual.(*matchAttributes), loaded
}
