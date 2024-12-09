package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var lastExitCode int = 0

const COMMAND_NOT_FOUND = 127

func main() {
	for {
		// Wait for user input
		cmd, err := getCmd()
		if err != nil {
			if err == io.EOF {
				fmt.Println("exit")
				os.Exit(lastExitCode)
			}
		}

		_ = handleCommand(cmd)
	}
}

// getCmd waits for user input and returns
func getCmd() (cmd string, err error) {
	fmt.Fprint(os.Stdout, "$ ")
	return bufio.NewReader(os.Stdin).ReadString('\n')
}

// handleCommand takes a command as input and executes it
func handleCommand(cmd string) int {
	if cmd[len(cmd)-1] == '\n' {
		cmd = cmd[:len(cmd)-1] // string the newline character
	}

	// hard code a few commands to start with
	if cmd == "exit" {
		os.Exit(lastExitCode)
	}

	if cmd == "$?" || cmd == "$status" {
		fmt.Println(lastExitCode)
		lastExitCode = 0
		return 0
	}

	// handle empty commands
	if cmd == "" {
		return 0
	}

	fmt.Fprintf(os.Stderr, "%s: command not found\n", cmd)
	lastExitCode = COMMAND_NOT_FOUND
	return lastExitCode
}
