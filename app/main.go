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

func createCommand(tokens []string) ICommand {
	if len(tokens) == 0 {
		return nil
	}

	switch tokens[0] {
	case "exit":
		return Exit{}
	case "echo":
		return Echo{Args: tokens[1:]}
	case "type":
		return Type{Args: tokens[1:]}
	case "pwd":
		return Pwd{Args: tokens[1:]}
	case "cd":
		return Cd{Args: tokens[1:]}
	default:
		return Unknown{Name: tokens[0], Args: tokens[1:]}
	}
}

// Given a tokenized input, handle processing the command.
// Returns false if the shell should exit, true otherwise.
func handleInput(input []string) {
	if len(input) == 0 {
		return
	}

	command := createCommand(input)
	command.Execute()
}

func main() {
	for true {
		displayPrompt()
		command := fetchTrimmedInput()
		tokens := tokenize(command)
		handleInput(tokens)
	}
}
