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

// capitalizeFirst capitalizes the first letter of a string.
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// formatChdirError transforms the default chdir error into the desired format.
func formatChdirError(err error, dir string) error {
	// Check if it's a "chdir" error (basic safeguard, can be expanded as needed)
	if strings.Contains(err.Error(), "chdir") {
		return fmt.Errorf(
			"cd: %s: %s",
			dir,
			capitalizeFirst(strings.Replace(err.Error(), "chdir "+dir+": ", "", 1)),
		)
	}
	// Return the original error if it's not what we expect
	return err
}

// parseSingleQuotes strips away all single quotes from the command.
func parseSingleQuotes(splitCommand []string) (cmdArray []string) {
	if len(splitCommand) == 0 {
		return splitCommand
	}

	cmdArray = make([]string, 0, len(splitCommand))
	cmdArray = append(cmdArray, splitCommand[0])
	currentPosition := 1

	for i := 1; i <= len(splitCommand[1:]); i++ {
		word := splitCommand[i]
		if strings.HasSuffix(word, "'") {
			cmdArray = append(
				cmdArray,
				strings.ReplaceAll(strings.Join(splitCommand[currentPosition:i+1], " "), "'", ""),
			)
			currentPosition = i + 1
		} else {
			if word == "" || strings.HasPrefix(word, "'") {
				continue // skip empty words and the start of single quotes
			}
			cmdArray = append(cmdArray, word)
		}
	}

	return
}
