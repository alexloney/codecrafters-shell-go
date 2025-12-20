package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Display the shell prompt, maybe in the future this
// could display user customizable prompts
func displayPrompt() {
	fmt.Print("$ ")
}

// Obtain the users input from stdin and trim any trailing newlines
func fetchTrimmedInput() string {
	command, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.TrimRight(command, "\r\n")
}

// First pass at tokenizing a string by splitting on spaces.
// This is a naive implementation and does not handle quotes or escaped spaces.
func tokenize(input string) []string {
	return strings.Split(input, " ")
}

// Given a tokenized input, handle processing the command.
// Returns false if the shell should exit, true otherwise.
func handleInput(input []string) bool {
	if len(input) == 0 {
		return true
	}

	command := input[0]
	args := input[1:]

	if command == "exit" {
		return false
	} else if command == "echo" {
		fmt.Println(strings.Join(args, " "))
		return true
	} else {
		fmt.Printf("%s: command not found\n", command)
	}

	return true
}

func main() {
	for true {
		displayPrompt()
		command := fetchTrimmedInput()
		tokens := tokenize(command)
		if !handleInput(tokens) {
			break
		}
	}
}
