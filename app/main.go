package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	for true {
		fmt.Print("$ ")

		// Captures the user's command in the "command" variable
		command, _ := bufio.NewReader(os.Stdin).ReadString('\n')

		trimmed_command := strings.TrimRight(command, "\r\n")
		if trimmed_command == "exit" {
			break
		}

		fmt.Print(trimmed_command)
		fmt.Println(": command not found")
	}
}
