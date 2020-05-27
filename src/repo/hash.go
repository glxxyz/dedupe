package repo

import (
	"bytes"
	"hash/crc32"
	"hash/crc64"
	"io"
	"os"
)

type HashOptions interface {
	Hash() bool
	Contents() bool
	Verbose() bool
}

// The idea of hashing the first few bytes came from https://stackoverflow.com/questions/748675/finding-duplicate-files-and-removing-them
func calculateHeadHash(options HashOptions, path string) (uint32, error) {
	if !options.Hash() {
		return 0, nil
	}
	file, err := os.Open(path)
	if err != nil {
		errLog.Printf("unable to open file: %v\n", err)
		return 0, err
	}
	defer file.Close()
	data := make([]byte, 1024)
	_, err = file.Read(data)
	if err != nil && err != io.EOF {
		errLog.Printf("error reading from file: %v\n", err)
		return 0, err
	}
	return crc32.ChecksumIEEE(data), nil
}

func calculateFullHash(options HashOptions, path string) (uint64, error) {
	if !options.Hash() {
		return 0, nil
	}
	file, err := os.Open(path)
	if err != nil {
		errLog.Printf("unable to open file: %v\n", err)
		return 0, err
	}
	defer file.Close()
	data := make([]byte, 8*1024)
	table := crc64.MakeTable(crc64.ECMA)
	var crc uint64
	for {
		_, err := file.Read(data)
		if err == io.EOF {
			break
		} else if err != nil {
			errLog.Printf("error reading from file: %v\n", err)
			return 0, err
		}
		crc = crc64.Update(crc, table, data)
	}
	return crc, nil
}

func fullByteMatch(options MatchOptions, pathA string, pathB string) (bool, error) {
	if !options.Contents() {
		return true, nil
	}

	fileA, err := os.Open(pathA)
	if err != nil {
		errLog.Printf("error opening file: %v\n", err)
		return false, err
	}
	defer fileA.Close()

	fileB, err := os.Open(pathB)
	if err != nil {
		errLog.Printf("error opening file %v\n", err)
		return false, err
	}
	defer fileB.Close()

	dataA := make([]byte, 8*1024)
	dataB := make([]byte, 8*1024)

	for {
		bytesA, err := fileA.Read(dataA)
		if err != nil && err != io.EOF {
			errLog.Printf("error reading from file: %v\n", err)
			return false, err
		}

		bytesB, err := fileB.Read(dataB)
		if err != nil && err != io.EOF {
			errLog.Printf("error reading from file: %v\n", err)
			return false, err
		}

		if bytesA == 0 && bytesB == 0 {
			return true, nil
		}

		if bytesA != bytesB {
			return false, nil
		}

		if !bytes.Equal(dataA, dataB) {
			return false, nil
		}
	}
}
