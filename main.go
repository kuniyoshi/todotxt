package main

import (
	"fmt"
	"os"
)

func main() {
	initTodoFile()
	
	command, args := parseArgs()
	
	if err := executeCommand(command, args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Run 'todo help' for usage information.")
		os.Exit(1)
	}
}

