package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("successful command, exit code = 0", func(t *testing.T) {
		exitCode := RunCmd([]string{"go", "version"}, nil)
		require.Equal(t, 0, exitCode)
	})

	t.Run("errored command, exit code = 1", func(t *testing.T) {
		exitCode := RunCmd([]string{"netstat", "-fake"}, nil)
		require.Equal(t, 1, exitCode)
	})
}
