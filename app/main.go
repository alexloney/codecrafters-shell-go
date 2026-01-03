package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

func tokenize(input string) ([]string, string, bool, string, bool) {
	// Return values
	var args []string
	var stdout string
	var appendStdout bool
	var stderr string
	var appendStderr bool

	// Intermediate storage for the raw split
	var rawTokens []string

	// Use strings.Builder for efficient concatenation
	var currentToken strings.Builder

	inQuotes := false
	inBigQuotes := false

	// --- PHASE 1: Lexical Analysis (Split by spaces, handle quotes) ---
	for i := 0; i < len(input); i++ {
		char := input[i]

		if char == '\\' {
			if i+1 < len(input) {
				if inQuotes {
					currentToken.WriteByte(char)
				}
				if inBigQuotes {
					next := input[i+1]
					if next != '"' && next != '\\' && next != '`' && next != '$' && next != '\n' {
						currentToken.WriteByte(char)
					}
				}
				currentToken.WriteByte(input[i+1])
				i++
			} else {
				currentToken.WriteByte(char)
			}
			continue
		}

		if char == '"' && !inQuotes {
			inBigQuotes = !inBigQuotes
			continue
		} else if char == '\'' && !inBigQuotes {
			inQuotes = !inQuotes
			continue
		}

		if char == ' ' && !inQuotes && !inBigQuotes {
			if currentToken.Len() > 0 {
				rawTokens = append(rawTokens, currentToken.String())
				currentToken.Reset()
			}
		} else {
			currentToken.WriteByte(char)
		}
	}

	// Capture the final token if it exists
	if currentToken.Len() > 0 {
		rawTokens = append(rawTokens, currentToken.String())
	}

	// --- PHASE 2: Redirection Processing ---
	// Iterate through rawTokens. If we see a redirection symbol,
	// consume the NEXT token as the file, and do NOT add to args.
	for i := 0; i < len(rawTokens); i++ {
		token := rawTokens[i]

		switch token {
		case "1>", ">":
			if i+1 < len(rawTokens) {
				stdout = rawTokens[i+1]
				appendStdout = false
				i++ // Skip the filename so it isn't added to args
			}
		case ">>", "1>>":
			if i+1 < len(rawTokens) {
				stdout = rawTokens[i+1]
				appendStdout = true
				i++
			}
		case "2>":
			if i+1 < len(rawTokens) {
				stderr = rawTokens[i+1]
				appendStderr = false
				i++
			}
		case "2>>":
			if i+1 < len(rawTokens) {
				stderr = rawTokens[i+1]
				appendStderr = true
				i++
			}
		default:
			// If it's not a redirection operator or file, it's a command argument
			args = append(args, token)
		}
	}

	return args, stdout, appendStdout, stderr, appendStderr
}

func createCommand(tokens []string) Commander {
	if len(tokens) == 0 {
		return nil
	}

	switch tokens[0] {
	case "exit":
		return &Exit{}
	case "echo":
		return &Echo{Args: tokens[1:]}
	case "type":
		return &Type{Args: tokens[1:]}
	case "pwd":
		return &Pwd{Args: tokens[1:]}
	case "cd":
		return &Cd{Args: tokens[1:]}
	case "history":
		return &History{Args: tokens[1:]}
	default:
		return &Executable{Name: tokens[0], Args: tokens[1:]}
	}
}

// Given a tokenized input, handle processing the command.
// Returns false if the shell should exit, true otherwise.
func handleInput(input []string, stdout string, append_stdout bool, stderr string, append_stderr bool) {
	if len(input) == 0 {
		return
	}

	command := createCommand(input)

	stdout_to_use := os.Stdout
	stderr_to_use := os.Stderr
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
		stdout_to_use = f
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
		stderr_to_use = f
		defer f.Close()
		// command.SetStderr(f)
	}

	command.SetStdin(os.Stdin)
	command.SetStdout(stdout_to_use)
	command.SetStderr(stderr_to_use)

	command.Execute()
}

func main() {
	// Setup the history file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error retrieving home directory:", err)
		return
	}
	historyFilePath := homeDir + "/.simple_shell_history"

	os.Remove(historyFilePath) // Needed to pass the CodeCrafters tests

	// Initialize the Readline instance
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "$ ",
		HistoryFile: historyFilePath,
		// AutoComplete: completer, // Future enhancement for tab completion
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing readline:", err)
		return
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break // Exit on Ctrl+C with empty line
				} else {
					continue
				}
			} else if err == io.EOF {
				break // Exit on Ctrl+D
			}
			break
		}

		// Cleanup input
		command := strings.TrimSpace(line)

		// If empty, just continue and re-prompt
		if command == "" {
			continue
		}

		tokens, stdout, appendStdout, stderr, appendStderr := tokenize(command)
		handleInput(tokens, stdout, appendStdout, stderr, appendStderr)
	}
}
