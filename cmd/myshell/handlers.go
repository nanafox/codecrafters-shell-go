package main

import (
	"fmt"
	"os"
	"os/exec"
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
	case "pwd":
		return handleCwd()
	case "cd":
		return handleChdir(splitCommand)
	default:
		return runCommand(splitCommand)
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

// handlCmdNotFound handles the case where the command is not found.
func handleCmdNotFound(command string) (code int) {
	fmt.Fprintf(os.Stderr, "%s: command not found\n", command)
	return COMMAND_NOT_FOUND
}

// handleCwd prints the current working directory and returns.
func handleCwd() (code int) {
	if cwd, err := os.Getwd(); err == nil {
		fmt.Println(cwd)
		code = EXIT_SUCCESS
	} else {
		fmt.Fprintf(
			os.Stderr, "pwd: an error occurred while retrieving the current path",
		)
		code = 1
	}

	return
}

// handleChdir handles the `cd` command.
func handleChdir(splitCommand []string) (code int) {
	var err error
	var path string

	oldPwd := os.Getenv("PWD")

	if len(splitCommand) > 1 {
		path = strings.Trim(splitCommand[1], " ")
	}
	if len(splitCommand) == 1 || path == "~" || path == "" {
		if err = os.Chdir(os.Getenv("HOME")); err == nil {
			updateCwd(oldPwd, os.Getenv("HOME"))
			return EXIT_SUCCESS
		}
	}

	if err = os.Chdir(path); err == nil {
		updateCwd(oldPwd, path)
		return EXIT_SUCCESS
	}

	if path == "-" {
		return handleOldPwd(oldPwd)
	}

	cdError := strings.Replace(err.Error(), "chdir", "cd", 1)
	fmt.Fprintln(os.Stderr, cdError)
	return EXIT_FAILURE
}

// handleOldPwd handles the `cd -` command and argument. This ensures that
// changing back and forth between current and the immediate past path is
// possible.
func handleOldPwd(oldPwd string) (code int) {
	newPwd := os.Getenv("OLDPWD")

	if newPwd == "" {
		newPwd = oldPwd
	}

	if err := os.Chdir(newPwd); err == nil {
		updateCwd(oldPwd, newPwd)

		fmt.Println(newPwd)
	} else {
		fmt.Fprintln(os.Stderr, err)
		return EXIT_FAILURE

	}

	return EXIT_SUCCESS
}
