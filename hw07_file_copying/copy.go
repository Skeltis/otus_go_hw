package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	progressbar "github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Place your code here.
	return copyInternal(fromPath, toPath, offset, limit)
}

func copyInternal(sourceFilePath, targetFilePath string, offset, limit int64) error {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return fmt.Errorf("error while openning file: %w", err)
	}
	defer sourceFile.Close()

	actualLimit, err := checkSourceCountLimit(sourceFile, offset, limit)
	if err != nil {
		return err
	}

	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	// actual counted limit was 0, so we only create empty file
	if actualLimit == 0 {
		return nil
	}

	if offset != 0 {
		_, err = sourceFile.Seek(offset, 0)
		if err != nil {
			return fmt.Errorf("error while moving file cursor: %w", err)
		}
	}

	progressBar := progressbar.Default(actualLimit)

	_, err = io.CopyN(io.MultiWriter(targetFile, progressBar), sourceFile, actualLimit)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("error while copying %s to %s: %w", sourceFilePath, targetFilePath, err)
	}

	return nil
}

func checkSourceCountLimit(file *os.File, offset, limit int64) (int64, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("error reading file info: %w", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return 0, ErrUnsupportedFile
	}

	fileSize := fileInfo.Size()
	if offset > fileSize {
		return 0, ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit > fileSize-offset {
		return fileSize - offset, nil
	}

	return limit, nil
}
