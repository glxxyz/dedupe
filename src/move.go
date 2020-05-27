package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type MoveOptions interface {
	DoMove() bool
	Trash() string
	Verbose() bool
}

func Move(options MoveOptions, filePath string) {
	if options.DoMove() {
		destPath := filepath.Join(options.Trash(), filePath)
		fmt.Printf(
			"Move:\t%v\t%v\n",
			strings.Replace(filePath, " ", "\\ ", -1),
			strings.Replace(destPath, " ", "\\ ", -1))
		folderPath := filepath.Dir(destPath)
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			errLog.Printf("error creating directory: %q: %v\n", folderPath, err)
		} else if err := os.Rename(filePath, destPath); err != nil {
			errLog.Printf("error moving file from: %q to: %q: %v\n", filePath, destPath, err)
		}
	} else {
		fmt.Printf("Move:\t%v\n", strings.Replace(filePath, " ", "\\ ", -1))
	}
}
