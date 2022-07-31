package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for envName, envVariable := range env {
		if envVariable.NeedRemove {
			err := os.Unsetenv(envName)
			if err != nil {
				log.Fatal(fmt.Errorf("couldn't unset %s, %w", envName, err))
				continue
			}
		}

		err := os.Setenv(envName, envVariable.Value)
		if err != nil {
			log.Fatal(fmt.Errorf("couldn't set %s = %s, %w", envName, envVariable.Value, err))
		}
	}

	command := composeCommand(cmd)
	exitCode, err := runCommandAndWaitCompletion(command)
	if err != nil {
		log.Print(err)
	}

	return exitCode
}

func composeCommand(arguments []string) *exec.Cmd {
	commandName := arguments[0]
	args := arguments[1:]
	command := exec.Command(commandName, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	return command
}

func runCommandAndWaitCompletion(command *exec.Cmd) (int, error) {
	err := command.Start()
	if err != nil {
		return 0, fmt.Errorf("command wasn't started %w", err)
	}

	err = command.Wait()
	if err != nil {
		var execError *exec.ExitError
		ok := errors.As(err, &execError)
		if ok {
			status, sysCallOk := execError.Sys().(syscall.WaitStatus)
			if sysCallOk {
				return status.ExitStatus(), nil
			}
		}
	}

	return 0, nil
}
