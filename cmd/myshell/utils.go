package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

// getCmd waits for user input and returns
func getCmd() (cmd string, err error) {
	fmt.Fprint(os.Stdout, "$ ")
	return bufio.NewReader(os.Stdin).ReadString('\n')
}

// updateCwd updates the `OLPDWD` and `PWD` variables so they match the current
// of the location in the filesystem. This helps to ensure that there's a smooth
// experiences when retrieving your current working directory and when you want
// to move between immediate-past directories and the current working directory.
func updateCwd(oldPath string, newPath string) {
	os.Setenv("PWD", newPath)
	os.Setenv("OLDPWD", oldPath)
}

// runCommand runs the command with the given arguments.
//
// It checks if the command is in the PATH and runs it if it is found. An error
// is returned if the command is not found.
func runCommand(splitCommand []string) (code int) {
	cmd := exec.Command(splitCommand[0], splitCommand[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return handleCmdNotFound(splitCommand[0])
		}

		return cmd.ProcessState.ExitCode() // return the error code from the process
	}
	return EXIT_SUCCESS
}

// isShellBuiltin checks if the command is a shell builtin.
func isShellBuiltin(command string) bool {
	return slices.Contains(shell_builtin_cmds, command)
}
