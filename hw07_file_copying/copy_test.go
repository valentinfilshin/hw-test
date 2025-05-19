package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	resultFilePath := "output.txt"

	tests := []struct {
		name         string // Название теста
		offset       int64  // Смещение
		limit        int64  // Лимит чтения
		expectedFile string // Ожидаемый файл результата
		originalFile string // Ожидаемый файл результата
	}{
		{"case offset0 limit0", 0, 0, "testdata/out_offset0_limit0.txt", "testdata/input.txt"},
		{"case offset0 limit10", 0, 10, "testdata/out_offset0_limit10.txt", "testdata/input.txt"},
		{"case offset0 limit1000", 0, 1000, "testdata/out_offset0_limit1000.txt", "testdata/input.txt"},
		{"case offset0 limit10000", 0, 10000, "testdata/out_offset0_limit10000.txt", "testdata/input.txt"},
		{"case offset100 limit1000", 100, 1000, "testdata/out_offset100_limit1000.txt", "testdata/input.txt"},
		{"case offset6000 limit1000", 6000, 1000, "testdata/out_offset6000_limit1000.txt", "testdata/input.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.originalFile, resultFilePath, tt.offset, tt.limit)
			if err != nil {
				t.Fatalf("Ошибка копирования файла: %v", err)
			}

			result := parallelCompareFiles(tt.expectedFile, resultFilePath)

			assert.True(t, result, fmt.Sprintf("Файлы '%s' и '%s' не совпадают", tt.expectedFile, resultFilePath))

			err = os.Remove(resultFilePath)
			if err != nil {
				t.Fatalf("Ошибка удаления файла: %v", err)
			}
		})
	}

	t.Run("big file copy", func(t *testing.T) {
		originalFilePath := "testdata/bigInput.txt"

		limit := int64(0)
		offset := int64(0)

		err := Copy(originalFilePath, resultFilePath, offset, limit)
		if err != nil {
			t.Fatalf("Ошибка копирования файла: %v", err)
		}

		result := parallelCompareFiles(originalFilePath, resultFilePath)

		assert.True(t, result, fmt.Sprintf("Файлы '%s' и '%s' не совпадают", originalFilePath, resultFilePath))

		err = os.Remove(resultFilePath)
		if err != nil {
			t.Fatalf("Ошибка удаления файла: %v", err)
		}
	})

	originalFilePath := "testdata/input.txt"
	t.Run("offset bigger then file size", func(t *testing.T) {
		limit := int64(0)
		offset := int64(10000)

		err := Copy(originalFilePath, resultFilePath, offset, limit)

		require.Error(t, err, "Ошибки нет")
		require.Equal(t, err, ErrOffsetExceedsFileSize, "Ошибка не верная")
	})

	t.Run("negative limit", func(t *testing.T) {
		limit := int64(-10)
		offset := int64(10000)

		err := Copy(originalFilePath, resultFilePath, offset, limit)

		require.Error(t, err, "Ошибки нет")
		require.Equal(t, err, ErrWrongLimit, "Ошибка не верная")
	})

	t.Run("negative offset", func(t *testing.T) {
		limit := int64(0)
		offset := int64(-10000)

		err := Copy(originalFilePath, resultFilePath, offset, limit)

		require.Error(t, err, "Ошибки нет")
		require.Equal(t, err, ErrWrongOffset, "Ошибка не верная")
	})

	t.Run("same file", func(t *testing.T) {
		samePath := "testdata/../testdata/input.txt"

		limit := int64(0)
		offset := int64(0)

		err := Copy(originalFilePath, samePath, offset, limit)

		require.Error(t, err, "Ошибки нет")
		require.Equal(t, err, ErrWrongPaths, "Ошибка не верная")
	})

	t.Run("same paths", func(t *testing.T) {
		limit := int64(0)
		offset := int64(0)

		err := Copy(originalFilePath, originalFilePath, offset, limit)

		require.Error(t, err, "Ошибки нет")
		require.Equal(t, err, ErrWrongPaths, "Ошибка не верная")
	})

	t.Run("cant work with /dev/null", func(t *testing.T) {
		originalFilePath := "/dev/null"

		limit := int64(0)
		offset := int64(0)

		err := Copy(originalFilePath, resultFilePath, offset, limit)

		require.Error(t, err, "Ошибки нет")
		require.Equal(t, err, ErrUnsupportedFile, "Ошибка не верная")
	})
}

func parallelCompareFiles(filePath1 string, filePath2 string) bool {
	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan []byte, 2)

	go func() {
		defer wg.Done()
		hasher := sha256.New()
		file1, _ := os.Open(filePath1)
		defer file1.Close()
		_, err := io.Copy(hasher, file1)
		if err != nil {
			fmt.Println(err)
		}

		ch <- hasher.Sum(nil)
	}()

	go func() {
		defer wg.Done()
		hasher := sha256.New()
		file2, _ := os.Open(filePath2)
		defer file2.Close()
		_, err := io.Copy(hasher, file2)
		if err != nil {
			fmt.Println(err)
		}

		ch <- hasher.Sum(nil)
	}()

	wg.Wait()
	close(ch)

	hashes := make([][]byte, 2)
	hashIndex := 0
	for hash := range ch {
		hashes[hashIndex] = hash
		hashIndex++
	}

	return bytes.Equal(hashes[0], hashes[1])
}
