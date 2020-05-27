package repo

import "sync"

type matchHeadHash struct {
	lock        sync.Mutex
	singlePath  string
	fullHashMap sync.Map // uint64 -> *matchFullHash
}

func (headHash *matchHeadHash) findFullMatch(options MatchOptions, filePath string) (*matchFullHash, bool) {
	headHash.ensureMapExists(options)
	hash, err := calculateFullHash(options, filePath)
	if err != nil {
		return nil, false
	}
	actual, loaded := headHash.fullHashMap.LoadOrStore(hash, &matchFullHash{filePaths: []string{filePath}})
	return actual.(*matchFullHash), loaded
}

func (headHash *matchHeadHash) ensureMapExists(options MatchOptions) {
	if headHash.singlePath != "" {
		headHash.lock.Lock()
		defer headHash.lock.Unlock()
		if headHash.singlePath != "" {
			hash, err := calculateFullHash(options, headHash.singlePath)
			if err == nil {
				headHash.fullHashMap.Store(hash, &matchFullHash{filePaths: []string{headHash.singlePath}})
				headHash.singlePath = ""
			}
		}
	}
}
