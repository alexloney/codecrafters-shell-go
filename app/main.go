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
	input := strings.TrimRight(command, "\r\n")

	// Log history
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error retrieving home directory:", err)
		return ""
	}
	file, err := os.OpenFile(homeDir+"/.simple_shell_history", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening history file:", err)
		return input
	}
	defer file.Close()
	if _, err := file.WriteString(input + "\n"); err != nil {
		fmt.Println("Error writing to history file:", err)
	}

	return input
}

// First pass at tokenizing a string by splitting on spaces.
// This is a naive implementation and does not handle quotes or escaped spaces.
func tokenize(input string) ([]string, string, bool, string, bool) {
	var output []string
	stdout := ""
	append_stdout := false
	stderr := ""
	append_stderr := false

	inQuotes := false
	inBigQuotes := false
	currentToken := ""
	// var currentToken strings.Builder

	for i := 0; i < len(input); i++ {
		char := input[i]

		if char == '\\' {
			// If the next character exists, add it literally
			if i+1 < len(input) {
				// Only add the backslash if we are inside quotes
				if inQuotes {
					currentToken += string(char)
				}

				if inBigQuotes {
					if input[i+1] != '"' && input[i+1] != '\\' && input[i+1] != '`' && input[i+1] != '$' && input[i+1] != '\n' {
						currentToken += string(char)
					}
				}

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

	for i := 0; i < len(output); i++ {
		token := output[i]

		if token == "1>" || token == ">" {
			if i+1 < len(output) {
				stdout = output[i+1]
				// Remove these two tokens from output
				output = append(output[:i], output[i+2:]...)
			}
		} else if token == ">>" || token == "1>>" {
			if i+1 < len(output) {
				stdout = output[i+1]
				append_stdout = true
				// Remove these two tokens from output
				output = append(output[:i], output[i+2:]...)
			}
		} else if token == "2>" {
			if i+1 < len(output) {
				stderr = output[i+1]
				// Remove these two tokens from output
				output = append(output[:i], output[i+2:]...)
			}
		} else if token == "2>>" {
			if i+1 < len(output) {
				stderr = output[i+1]
				append_stderr = true
				// Remove these two tokens  from output
				output = append(output[:i], output[i+2:]...)
			}
		}
	}

	return output, stdout, append_stdout, stderr, append_stderr
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
	case "history":
		return History{Args: tokens[1:]}
	// case "cat":
	// 	return Cat{Args: tokens[1:]}
	default:
		return Unknown{Name: tokens[0], Args: tokens[1:]}
	}
}

// Given a tokenized input, handle processing the command.
// Returns false if the shell should exit, true otherwise.
func handleInput(input []string, stdout string, append_stdout bool, stderr string, append_stderr bool) {
	if len(input) == 0 {
		return
	}

	command := createCommand(input)

	old_stdout := os.Stdout
	old_stderr := os.Stderr
	if stdout != "" {
		flags := os.O_RDWR | os.O_CREATE
		if append_stdout {
			flags |= os.O_APPEND
		}
		f, err := os.OpenFile(stdout, flags, 0644)
		if err != nil {
			fmt.Println("Error opening file for stdout redirection:", err)
			return
		}
		os.Stdout = f
		defer f.Close()
		// command.SetStdout(f)
	}

	if stderr != "" {
		flags := os.O_RDWR | os.O_CREATE
		if append_stderr {
			flags |= os.O_APPEND
		}
		f, err := os.OpenFile(stderr, flags, 0644)
		if err != nil {
			fmt.Println("Error opening file for stderr redirection:", err)
			return
		}
		os.Stderr = f
		defer f.Close()
		// command.SetStderr(f)
	}

	command.Execute()

	os.Stdout = old_stdout
	os.Stderr = old_stderr
}

func main() {
	h := History{}
	filepath, _ := h.GetHistoryFilePath()
	os.Remove(filepath)

	for true {
		displayPrompt()
		command := fetchTrimmedInput()
		tokens, stdout, append_stdout, stderr, append_stderr := tokenize(command)
		// fmt.Println(tokens)
		// for _, token := range tokens {
		// 	fmt.Println("Token:'", token, "'")
		// }
		handleInput(tokens, stdout, append_stdout, stderr, append_stderr)
	}
}
