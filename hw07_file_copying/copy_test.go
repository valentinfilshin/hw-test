package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
			fmt.Println(err)
		}

		result, err := CompareFiles(expectedFilePath, resultFilePath)
		if err != nil {
			fmt.Println(err)
		}

		assert.True(t, result)
		_ = os.Remove(resultFilePath)
	})

	// TODO Сделай табличный тест
	// TODO Как обрабатывать ошибки в тесте?

	// Данный тест не проходит на Windows, потому что перенос строки занимает более 1 байта
	t.Run("test offset0 limit10", func(t *testing.T) {
		expectedFilePath := "testdata/out_offset0_limit10.txt"

		offset := int64(0)
		limit := int64(10)

		err := Copy(originalFilePath, resultFilePath, offset, limit)
		if err != nil {
			fmt.Println(err)
		}

		result, err := CompareFiles(expectedFilePath, resultFilePath)
		if err != nil {
			fmt.Println(err)
		}

		assert.True(t, result)
		_ = os.Remove(resultFilePath)
	})

	// TODO проверка ошибок
}
