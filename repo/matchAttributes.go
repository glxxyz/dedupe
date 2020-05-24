package repo

import (
	"sync"
)

type matchAttributes struct {
	lock       sync.Mutex
	singlePath string
	headMap    sync.Map // uint32 -> *matchHeadHash
}

func (attributes *matchAttributes) findHeadMatch(options MatchOptions, filePath string) (*matchHeadHash, bool) {
	attributes.ensureMapExists(options)
	hash, err := calculateHeadHash(options, filePath)
	if err != nil {
		return nil, false
	}
	actual, loaded := attributes.headMap.LoadOrStore(hash, &matchHeadHash{singlePath: filePath})
	return actual.(*matchHeadHash), loaded
}

func (attributes *matchAttributes) ensureMapExists(options MatchOptions) {
	if attributes.singlePath != "" {
		attributes.lock.Lock()
		defer attributes.lock.Unlock()
		if attributes.singlePath != "" {
			hash, err := calculateHeadHash(options, attributes.singlePath)
			if err == nil {
				attributes.headMap.Store(hash, &matchHeadHash{singlePath: attributes.singlePath})
				attributes.singlePath = ""
			}
		}
	}
}