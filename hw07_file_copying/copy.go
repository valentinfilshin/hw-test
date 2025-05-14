package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrCreateFile            = errors.New("can't create file")
	ErrReadFile              = errors.New("can't read file")
	ErrFileInfoNil           = errors.New("file info is nil")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if fromFileInfo == nil {
		return ErrFileInfoNil
	}

	if !fromFileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fromFileSize := fromFileInfo.Size()

	if offset > fromFileSize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit > fromFileSize-offset {
		limit = fromFileSize - offset
	}

	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0644)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer fromFile.Close()

	if offset > 0 {
		if _, err := fromFile.Seek(offset, io.SeekStart); err != nil {
			return err
		}
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return ErrCreateFile
	}
	defer toFile.Close()

	fmt.Println(limit)
	_, err = io.CopyN(toFile, fromFile, limit)
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = nil
			return nil
		}
		return err
	}

	return nil
}

func CompareFiles(expected, actual string) (bool, error) {
	bytes1, err := os.ReadFile(expected)
	if err != nil {
		return false, ErrReadFile
	}

	bytes2, err := os.ReadFile(actual)
	if err != nil {
		return false, ErrReadFile
	}

	return bytes.Equal(bytes1, bytes2), nil
}
