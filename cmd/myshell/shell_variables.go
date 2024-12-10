package main

var lastExitCode int = 0

const EXIT_SUCCESS = 0

const INVALID_ARGUMENT = 1

const ILLEGAL_NUMBER = 2

const COMMAND_NOT_FOUND = 127

var shell_builtin_cmds = []string{
	"exit", "echo", "type", "pwd",
}
