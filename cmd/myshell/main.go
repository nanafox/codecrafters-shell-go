package main

import (
	"fmt"
	"io"
	"os"
)

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

		lastExitCode = handleCommand(cmd)
	}
}
