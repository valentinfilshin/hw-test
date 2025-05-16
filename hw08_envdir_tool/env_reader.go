package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	fileSlice, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fmt.Println(fileSlice)

	environmentMap := make(Environment, len(dir))

	for _, file := range fileSlice {
		fileName := file.Name()

		if file.IsDir() || strings.Contains(fileName, "=") {
			continue
		}

		if envValue, err := readEnvFile(dir, fileName); err == nil {
			environmentMap[fileName] = envValue
		}
	}

	return environmentMap, nil
}

func readEnvFile(dir, fileName string) (EnvValue, error) {
	filePath := path.Join(dir, fileName)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return EnvValue{}, err
	}

	if !fileInfo.Mode().IsRegular() {
		return EnvValue{}, fmt.Errorf("%s is not a regular file", fileName)
	}

	fileSize := fileInfo.Size()

	if fileSize == 0 {
		return EnvValue{"", true}, nil
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return EnvValue{}, err
	}

	str := formatFileData(fileData)

	return EnvValue{str, false}, nil
}

func formatFileData(fileData []byte) string {
	idx := bytes.Index(fileData, []byte("\n"))

	if idx != -1 {
		fileData = fileData[:idx]
	}

	fileData = bytes.ReplaceAll(fileData, []byte("\000"), []byte("\n"))
	fileData = bytes.TrimRight(fileData, " \t")

	return string(fileData)
}
