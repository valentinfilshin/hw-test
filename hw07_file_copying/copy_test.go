package main

import (
	"bytes"
	"fmt"
	"os"
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
	}{
		{"case offset0 limit0", 0, 0, "testdata/out_offset0_limit0.txt"},
		{"case offset0 limit10", 0, 10, "testdata/out_offset0_limit10.txt"},
		{"case offset0 limit1000", 0, 1000, "testdata/out_offset0_limit1000.txt"},
		{"case offset0 limit10000", 0, 10000, "testdata/out_offset0_limit10000.txt"},
		{"case offset100 limit1000", 100, 1000, "testdata/out_offset100_limit1000.txt"},
		{"case offset6000 limit1000", 6000, 1000, "testdata/out_offset6000_limit1000.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalFilePath := "testdata/input.txt"

			err := Copy(originalFilePath, resultFilePath, tt.offset, tt.limit)
			if err != nil {
				t.Fatalf("Ошибка копирования файла: %v", err)
			}

			result, err := CompareFiles(tt.expectedFile, resultFilePath)
			if err != nil {
				t.Fatalf("Ошибка сравнения файлов: %v", err)
			}

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

		result, err := CompareFiles(originalFilePath, resultFilePath)
		if err != nil {
			t.Fatalf("Ошибка сравнения файлов: %v", err)
		}

		assert.True(t, result, fmt.Sprintf("Файлы '%s' и '%s' не совпадают", originalFilePath, resultFilePath))

		err = os.Remove(resultFilePath)
		if err != nil {
			t.Fatalf("Ошибка удаления файла: %v", err)
		}
	})

	t.Run("offset bigger then file size", func(t *testing.T) {
		originalFilePath := "testdata/input.txt"

		limit := int64(0)
		offset := int64(10000)

		err := Copy(originalFilePath, resultFilePath, offset, limit)

		require.Error(t, err, "Ошибки нет")
		require.Equal(t, err, ErrOffsetExceedsFileSize, "Ошибка не верная")
	})

	t.Run("negative limit", func(t *testing.T) {
		originalFilePath := "testdata/input.txt"

		limit := int64(-10)
		offset := int64(10000)

		err := Copy(originalFilePath, resultFilePath, offset, limit)

		require.Error(t, err, "Ошибки нет")
		require.Equal(t, err, ErrWrongLimit, "Ошибка не верная")
	})

	t.Run("negative offset", func(t *testing.T) {
		originalFilePath := "testdata/input.txt"

		limit := int64(0)
		offset := int64(-10000)

		err := Copy(originalFilePath, resultFilePath, offset, limit)

		require.Error(t, err, "Ошибки нет")
		require.Equal(t, err, ErrWrongOffset, "Ошибка не верная")
	})

	t.Run("same paths", func(t *testing.T) {
		originalFilePath := "testdata/input.txt"

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
