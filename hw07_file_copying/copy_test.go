package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

type limitOffsetTest struct {
	limit            int64
	offset           int64
	expectedFilePath string
}

func TestCopy(t *testing.T) {
	inputFile := path.Join("testdata", "input.txt")
	targetFile := path.Join("testdata", "output.txt")

	if runtime.GOOS != "windows" {
		t.Run("unsupported file", func(t *testing.T) {
			sourceFile := path.Join("/", "dev", "zero")
			err := Copy(sourceFile, targetFile, 0, 0)
			require.Error(t, err)
			require.ErrorIs(t, err, ErrUnsupportedFile)
		})
	}

	t.Run("offset is bigger than file", func(t *testing.T) {
		err := Copy(inputFile, targetFile, 10000, 0)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("offset leads to the end of file", func(t *testing.T) {
		fileBytes, err := os.ReadFile(inputFile)
		require.NoError(t, err)

		err = Copy(inputFile, targetFile, int64(len(fileBytes)), 0)
		require.NoError(t, err)

		copiedBytes, err := os.ReadFile(targetFile)
		require.NoError(t, err)
		require.Len(t, copiedBytes, 0)
	})

	validTestCases := []limitOffsetTest{
		{
			offset:           0,
			limit:            0,
			expectedFilePath: path.Join("testdata", "out_offset0_limit0.txt"),
		},
		{
			offset:           0,
			limit:            10,
			expectedFilePath: path.Join("testdata", "out_offset0_limit10.txt"),
		},
		{
			offset:           0,
			limit:            1000,
			expectedFilePath: path.Join("testdata", "out_offset0_limit1000.txt"),
		},
		{
			offset:           0,
			limit:            10000,
			expectedFilePath: path.Join("testdata", "out_offset0_limit10000.txt"),
		},
		{
			offset:           100,
			limit:            10000,
			expectedFilePath: path.Join("testdata", "out_offset100_limit10000.txt"),
		},
		{
			offset:           6000,
			limit:            1000,
			expectedFilePath: path.Join("testdata", "out_offset6000_limit1000.txt"),
		},
	}

	for _, currentCase := range validTestCases {
		t.Run(fmt.Sprintf("offset %d limit %d", currentCase.offset, currentCase.limit), func(t *testing.T) {
			err := Copy(inputFile, targetFile, currentCase.offset, currentCase.limit)
			require.NoError(t, err)

			copiedBytes, err := os.ReadFile(targetFile)
			require.NoError(t, err)

			var expectedBytes []byte
			if runtime.GOOS == "windows" {
				fileBytes, err := os.ReadFile(inputFile)
				require.NoError(t, err)

				endIndex := currentCase.offset + currentCase.limit
				if endIndex > int64(len(fileBytes)) || currentCase.limit == 0 {
					endIndex = int64(len(fileBytes))
				}

				expectedBytes = fileBytes[currentCase.offset:endIndex]
			} else {
				expectedBytes, err = os.ReadFile(currentCase.expectedFilePath)
			}

			require.NoError(t, err)
			require.Equal(t, expectedBytes, copiedBytes)

			os.Remove(targetFile)
		})
	}
}
