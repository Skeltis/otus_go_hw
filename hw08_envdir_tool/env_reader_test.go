package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type testFixture struct {
	testName       string
	fileName       string
	expectedString string
	expectedRemove bool
}

func TestReadDir(t *testing.T) {
	t.Run("Run on unexisting directory returns error", func(t *testing.T) {
		res, err := ReadDir("this-directory-doesn't exist")

		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("Run on directory with = in dir name returns error", func(t *testing.T) {
		dir, _ := os.MkdirTemp("", "tst")
		ioutil.TempFile(dir, "t=st")
		defer os.RemoveAll(dir)
		res, err := ReadDir(dir)

		require.ErrorContains(t, err, "file name contains invalid characters")
		require.Nil(t, res)
	})

	t.Run("Run on empty dir returns empty environment", func(t *testing.T) {
		dir, _ := os.MkdirTemp("", "tst")
		defer os.RemoveAll(dir)

		res, err := ReadDir(dir)

		require.NoError(t, err)
		require.Empty(t, res)
	})

	t.Run("Run on directory with subdirectories process only files", func(t *testing.T) {
		dir, _ := os.MkdirTemp("", "tst")
		os.MkdirTemp(dir, "tstSubdir")
		file, _ := ioutil.TempFile(dir, "fileTst")
		defer os.RemoveAll(dir)
		defer file.Close()
		res, err := ReadDir(dir)

		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, "", res[file.Name()].Value)
	})

	t.Run("Run on directory with subdirectories process only files", func(t *testing.T) {
		dir, _ := os.MkdirTemp("", "tst")
		os.MkdirTemp(dir, "tstSubdir")
		file, _ := ioutil.TempFile(dir, "fileTst")
		defer os.RemoveAll(dir)
		defer file.Close()
		res, err := ReadDir(dir)

		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, "", res[file.Name()].Value)
	})

	t.Run("Running on test data set", func(t *testing.T) {
		dir := filepath.Join("testdata", "env")
		result, err := ReadDir(dir)
		expected := Environment{
			"BAR": {
				Value:      "bar",
				NeedRemove: false,
			},
			"EMPTY": {
				Value:      "",
				NeedRemove: false,
			},
			"FOO": {
				Value:      "   foo\nwith new line",
				NeedRemove: false,
			},
			"HELLO": {
				Value:      "\"hello\"",
				NeedRemove: false,
			},
			"UNSET": {
				Value:      "",
				NeedRemove: true,
			},
		}
		require.Equal(t, expected, result)
		require.NoError(t, err)
	})
}

func TestExtractVariableFromFile(t *testing.T) {
	successfulFixtures := []testFixture{
		{
			testName:       "extract value from multiline file",
			fileName:       "BAR",
			expectedString: "bar",
			expectedRemove: false,
		},
		{
			testName:       "extract value from empty file",
			fileName:       "EMPTY",
			expectedString: "",
			expectedRemove: false,
		},
		{
			testName:       "extract value with 0x00 replacement and trim",
			fileName:       "FOO",
			expectedString: "   foo\nwith new line",
			expectedRemove: false,
		},
		{
			testName:       "extract value up to EOF",
			fileName:       "HELLO",
			expectedString: "\"hello\"",
			expectedRemove: false,
		},
		{
			testName:       "test unset value",
			fileName:       "UNSET",
			expectedString: "",
			expectedRemove: true,
		},
	}

	for _, tstFixture := range successfulFixtures {
		t.Run(tstFixture.testName, func(t *testing.T) {
			fileName := filepath.Join("testdata", "env", tstFixture.fileName)
			value, err := extractVariableFromFile(fileName)
			require.Equal(t, tstFixture.expectedString, value.Value)
			require.Equal(t, tstFixture.expectedRemove, value.NeedRemove)
			require.NoError(t, err)
		})
	}

	t.Run("testing not existing file", func(t *testing.T) {
		fileName := filepath.Join("testdata", "env", "who-are-you")
		_, err := extractVariableFromFile(fileName)
		require.ErrorContains(t, err, "error openning file")
	})
}
