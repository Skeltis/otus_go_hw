package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	commandLineArguments := os.Args[1:]
	if len(commandLineArguments) < 2 {
		log.Fatalf("At least command and directory should be specified: %s", commandLineArguments)
	}

	directory, err := getDirectory(commandLineArguments)
	if err != nil {
		log.Fatalf("%v", err)
	}

	environment, err := ReadDir(directory)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot read dir = %w", err))
	}

	exitCode := RunCmd(commandLineArguments[1:], environment)
	os.Exit(exitCode)
}

func getDirectory(arguments []string) (string, error) {
	folderPath := arguments[0]

	if _, err := os.Stat(folderPath); err != nil {
		return "", fmt.Errorf("error while extracting folder argument: %w", err)
	}
	return folderPath, nil
}
