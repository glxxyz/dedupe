package repo

import (
	"fmt"
	"strings"
	"sync"
)

type matchFullHash struct {
	lock      sync.Mutex
	filePaths []string
}

func (fullHash *matchFullHash) lowestPriorityMatch(options MatchOptions, filePath string) (string, bool) {
	fullHash.lock.Lock()
	defer fullHash.lock.Unlock()
	for num, testPath := range fullHash.filePaths {
		match, _ := fullByteMatch(options, testPath, filePath)
		if match {
			if firstIsHigherPriority(options.Paths(), testPath, filePath) {
				fmt.Printf("Dupe:\t%q\t%q\n", testPath, filePath)
				return filePath, true
			} else {
				fullHash.filePaths[num] = filePath
				fmt.Printf("Dupe:\t%q\t%q\n", filePath, testPath)
				return testPath, true
			}
		}
	}
	// There was no match, this implies a hash collision or a problem comparing the file
	fullHash.filePaths = append(fullHash.filePaths, filePath)
	return "", false
}

func firstIsHigherPriority(priorityPaths []string, first string, second string) bool {
	for _, priority := range priorityPaths {
		firstTest := strings.Index(first, priority)
		secondTest := strings.Index(second, priority)
		if firstTest == 0 && secondTest != 0 {
			return true
		} else if firstTest != 0 && secondTest == 0 {
			return false
		}
	}
	return first < second
}
