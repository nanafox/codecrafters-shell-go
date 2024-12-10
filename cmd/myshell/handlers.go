package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

// handleCommand takes a command as input and executes it
func handleCommand(cmd string) int {
	if cmd[len(cmd)-1] == '\n' {
		cmd = cmd[:len(cmd)-1] // string the newline character
	}

	splitCommand := strings.Split(cmd, " ")
	majorCommand := splitCommand[0]

	switch majorCommand {
	case "exit":
		return handleExit(splitCommand)
	case "":
		return EXIT_SUCCESS
	case "echo":
		handleEcho(splitCommand)
	case "type":
		return handleType(splitCommand)
	default:
		return handleCmdNotFound(majorCommand)
	}

	return EXIT_SUCCESS
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

// handleEcho handles the echo command and prints the arguments to the terminal
// or the exit code if the argument is $? or $status
func handleEcho(splitCommand []string) {
	if len(splitCommand) == 2 {
		if splitCommand[1] == "$?" || splitCommand[1] == "$status" {
			fmt.Println(lastExitCode)
			return
		}
	}

	fmt.Println(strings.Join(splitCommand[1:], " "))
}

// handleType handles the type command and prints the type of the command.
func handleType(splitCommand []string) (code int) {
	if len(splitCommand) != 2 {
		fmt.Fprintf(os.Stderr, "type: invalid number of arguments\n")
		return INVALID_ARGUMENT
	}
	command := splitCommand[1]
	if isShellBuiltin(command) {
		fmt.Printf("%s is a shell builtin\n", command)
	} else if path, err := exec.LookPath(command); err == nil {
		fmt.Printf("%s is %s\n", command, path)
	} else {
		fmt.Printf("%s: not found\n", command)
		return COMMAND_NOT_FOUND
	}
	return EXIT_SUCCESS
}

// isShellBuiltin checks if the command is a shell builtin.
func isShellBuiltin(command string) bool {
	return slices.Contains(shell_builtin_cmds, command)
}

// handlCmdNotFound handles the case where the command is not found.
func handleCmdNotFound(command string) (code int) {
	fmt.Fprintf(os.Stderr, "%s: command not found\n", command)
	return COMMAND_NOT_FOUND
}
