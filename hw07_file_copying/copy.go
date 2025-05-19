//nolint:depguard
package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrCreateFile            = errors.New("can't create file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrWrongOffset           = errors.New("offset can't be negative")
	ErrWrongLimit            = errors.New("limit can't be negative")
	ErrWrongPaths            = errors.New("path to files can't be same")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrWrongOffset
	}

	if limit < 0 {
		return ErrWrongLimit
	}

	fromFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	toFileInfo, err := os.Stat(toPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	if os.SameFile(fromFileInfo, toFileInfo) {
		return ErrWrongPaths
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

	fromFile, err := os.Open(fromPath)
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

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(fromFile)

	_, err = io.CopyN(toFile, barReader, limit)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	bar.Finish()

	return nil
}
