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
	var output []string

	inQuotes := false
	inBigQuotes := false
	currentToken := ""
	// var currentToken strings.Builder

	for i := 0; i < len(input); i++ {
		char := input[i]

		if char == '\\' {
			// If the next character exists, add it literally
			if i+1 < len(input) {
				currentToken += string(input[i+1])
				i++ // Skip the next character as we've already processed it
			} else {
				// If there's no next character, just add the backslash
				currentToken += string(char)
			}
			continue
		}

		// If we hit a double-quote, while not in a single-quote, then process it
		if char == '"' && !inQuotes {
			inBigQuotes = !inBigQuotes
			continue
		} else if char == '\'' && !inBigQuotes {
			inQuotes = !inQuotes
			continue
		}

		if char == ' ' && !inQuotes && !inBigQuotes {
			if len(currentToken) > 0 {
				output = append(output, currentToken)
				currentToken = ""
			}
		} else {
			currentToken += string(char)
		}
	}

	if len(currentToken) > 0 {
		output = append(output, currentToken)
	}

	return output
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
		// fmt.Println(tokens)
		// for _, token := range tokens {
		// 	fmt.Println("Token:'", token, "'")
		// }
		handleInput(tokens)
	}
}
