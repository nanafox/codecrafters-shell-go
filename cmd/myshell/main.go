package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var lastExitCode int = 0

const EXIT_SUCCESS = 0

const INVALID_ARGUMENT = 1

const ILLEGAL_NUMBER = 2

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

	splitCommand := strings.Split(cmd, " ")
	majorCommand := splitCommand[0]

	// hard code a few commands to start with
	if majorCommand == "exit" {
		return handleExit(splitCommand)
	}

	if cmd == "$?" || cmd == "$status" {
		fmt.Println(lastExitCode)
		lastExitCode = EXIT_SUCCESS
		return lastExitCode
	}

	// handle empty commands
	if cmd == "" {
		return EXIT_SUCCESS
	}

	fmt.Fprintf(os.Stderr, "%s: command not found\n", cmd)
	lastExitCode = COMMAND_NOT_FOUND
	return lastExitCode
}

// handleExit handles the exit command and exits the shell when the command is
// called with the correct arguments
func handleExit(splitCommand []string) (code int) {
	if len(splitCommand) == 2 {
		code, err := strconv.Atoi(splitCommand[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "exit: %s: numeric argument required\n", splitCommand[1])
			lastExitCode = ILLEGAL_NUMBER
			return lastExitCode
		}
		os.Exit(code)
	}
	os.Exit(lastExitCode)
	return
}
