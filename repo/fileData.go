package repo

import (
	"os"
	"time"
)

type FileData struct {
	FilePath string
	Name     string
	size     int64
	modTime  time.Time
}

func NewFile(options MatchOptions, filePath string, info os.FileInfo) *FileData {
	var name string
	var modTime time.Time
	var size int64
	if options.Name() {
		name = info.Name()
	}
	if options.ModTime() {
		modTime = info.ModTime()
	}
	if options.Size() {
		size = info.Size()
	}
	return &FileData{
		FilePath: filePath,
		Name:     name,
		size:     size,
		modTime:  modTime,
	}
}
