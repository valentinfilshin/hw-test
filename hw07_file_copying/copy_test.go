package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
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
		tt := tt // Захватываем переменную внутри итерации
		t.Run(tt.name, func(t *testing.T) {
			originalFilePath := "testdata/input.txt"
			resultFilePath := "output.txt"

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

	t.Run("offset bigger then file size", func(t *testing.T) {
		originalFilePath := "testdata/input.txt"
		resultFilePath := "output.txt"

		limit := int64(0)
		offset := int64(10000)

		err := Copy(originalFilePath, resultFilePath, offset, limit)

		require.Error(t, err, "Ошибки нет")
		require.Equal(t, err, ErrOffsetExceedsFileSize, "Ошибка не верная")
	})

	t.Run("cant work with /dev/null", func(t *testing.T) {
		originalFilePath := "/dev/null"
		resultFilePath := "output.txt"

		limit := int64(0)
		offset := int64(0)

		err := Copy(originalFilePath, resultFilePath, offset, limit)

		require.Error(t, err, "Ошибки нет")
		require.Equal(t, err, ErrUnsupportedFile, "Ошибка не верная")
	})
}
