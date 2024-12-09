package main

import (
	"bufio"
	"fmt"
	"os"
)

// getCmd waits for user input and returns
func getCmd() (cmd string, err error) {
	fmt.Fprint(os.Stdout, "$ ")
	return bufio.NewReader(os.Stdin).ReadString('\n')
}
