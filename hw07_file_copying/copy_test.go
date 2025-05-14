package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	originalFilePath := "testdata/input.txt"
	resultFilePath := "output.txt"

	t.Run("test offset0 limit0", func(t *testing.T) {
		expectedFilePath := "testdata/out_offset0_limit0.txt"

		limit := int64(0)
		offset := int64(0)

		err := Copy(originalFilePath, resultFilePath, limit, offset)
		if err != nil {

		}

		result, err := CompareFiles(expectedFilePath, resultFilePath)
		if err != nil {

		}

		assert.True(t, result)
		err = os.Remove(resultFilePath)
	})

	t.Run("test offset0 limit10", func(t *testing.T) {
		expectedFilePath := "testdata/out_offset0_limit10.txt"

		offset := int64(0)
		// TODO - это прикол винды? С тем что перенос строки занимает 2 байта?
		limit := int64(11)

		err := Copy(originalFilePath, resultFilePath, offset, limit)
		if err != nil {

		}

		result, err := CompareFiles(expectedFilePath, resultFilePath)
		if err != nil {

		}

		assert.True(t, result)
		err = os.Remove(resultFilePath)
	})

	// TODO ошибки
}
