package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error while reading directory: %w", err)
	}

	environmentResult := make(Environment, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.Contains(file.Name(), "=") {
			return nil, fmt.Errorf("file name contains invalid characters: %s", file.Name())
		}

		value, err := extractVariableFromFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		environmentResult[file.Name()] = *value
	}

	return environmentResult, nil
}

func extractVariableFromFile(filePath string) (*EnvValue, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error openning file: %s, %w", filePath, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("error reading file info: %s, %w", filePath, err)
	}

	if fileInfo.Size() == 0 {
		return &EnvValue{
			NeedRemove: true,
		}, nil
	}

	envValue, err := readEnvironmentalValue(file)
	if err != nil {
		return nil, err
	}

	return &EnvValue{
		Value: envValue,
	}, nil
}

func readEnvironmentalValue(file *os.File) (string, error) {
	reader := bufio.NewReader(file)

	byteLine, _, err := reader.ReadLine()
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("error reading line: %s, %w", file.Name(), err)
	}

	readLine := string(byteLine)
	readLine = strings.ReplaceAll(readLine, "\x00", "\n")
	readLine = strings.TrimRight(readLine, " \t")

	return readLine, nil
}
