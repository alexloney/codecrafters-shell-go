package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/chzyer/readline"
)

// Splits the input string by pipes, taking care to ignore pipes within quotes.
func splitByPipe(input string) []string {
	var commands []string
	var current strings.Builder
	inQuote := false
	inDoubleQuote := false

	for i := 0; i < len(input); i++ {
		char := input[i]
		if char == '\'' && !inDoubleQuote {
			inQuote = !inQuote
		} else if char == '"' && !inQuote {
			inDoubleQuote = !inDoubleQuote
		}

		if char == '|' && !inQuote && !inDoubleQuote {
			if current.Len() > 0 {
				commands = append(commands, strings.TrimSpace(current.String()))
				current.Reset()
			}
		} else {
			current.WriteByte(char)
		}
	}
	if current.Len() > 0 {
		commands = append(commands, strings.TrimSpace(current.String()))
	}
	return commands
}

// Returns a fully configured Commander, ready to run
func parseCommand(input string, historyMgr *HistoryManager) *PreparedCommand {
	// tokenize handles parsing ">", ">>", etc.
	tokens, stdoutPath, appendStdout, stderrPath, appendStderr := tokenize(input)

	cmd := createCommand(tokens, historyMgr)
	if cmd == nil {
		return nil
	}

	pc := &PreparedCommand{
		Command: cmd,
	}

	// Handle STDOUT Redirection
	if stdoutPath != "" {
		pc.StdoutRedirected = true
		flags := os.O_RDWR | os.O_CREATE
		if appendStdout {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}
		f, err := os.OpenFile(stdoutPath, flags, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening stdout file:", err)
			return nil
		}
		cmd.SetStdout(f)
	}

	// Handle STDERR Redirection
	if stderrPath != "" {
		pc.StderrRedirected = true
		flags := os.O_RDWR | os.O_CREATE
		if appendStderr {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}
		f, err := os.OpenFile(stderrPath, flags, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening stderr file:", err)
			return nil
		}
		cmd.SetStderr(f)
	}

	return pc
}

func executePipeline(rawInput string, historyMgr *HistoryManager) {
	rawCommands := splitByPipe(rawInput)
	if len(rawCommands) == 0 {
		return
	}

	// Parse all commands
	var preparedCmds []*PreparedCommand
	for _, rawCmd := range rawCommands {
		pc := parseCommand(rawCmd, historyMgr)
		if pc == nil {
			return
		}
		preparedCmds = append(preparedCmds, pc)
	}

	var wg sync.WaitGroup

	// 'nextStdin' carries the Read-end of the pipe to the NEXT command.
	// Initial input is os.Stdin.
	var nextStdin io.Reader = os.Stdin

	for i, pc := range preparedCmds {
		cmd := pc.Command

		// 1. SETUP STDIN
		// Always take from the chain (either os.Stdin or previous pipe)
		cmd.SetStdin(nextStdin)

		// 2. SETUP STDOUT
		var pipeWriter *os.File

		// If this is NOT the last command, we usually output to a pipe...
		if i < len(preparedCmds)-1 {
			r, w, err := os.Pipe()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Pipe error:", err)
				return
			}

			// ...BUT if user already redirected to a file, we respect that.
			if pc.StdoutRedirected {
				// We keep the file as stdout.
				// We must close the pipe writer immediately because nothing will write to it.
				// The next command (r) will read EOF immediately.
				w.Close()
			} else {
				// Normal case: pipe output
				cmd.SetStdout(w)
				pipeWriter = w // We must close this after execution
			}

			// The next command reads from this pipe
			nextStdin = r

		} else {
			// This IS the last command.
			// If no file redirection, output to console.
			if !pc.StdoutRedirected {
				cmd.SetStdout(os.Stdout)
			}
		}

		// 3. SETUP STDERR
		// Pipes usually don't carry stderr. Default to console if not redirected.
		if !pc.StderrRedirected {
			cmd.SetStderr(os.Stderr)
		}

		wg.Add(1)
		go func(c Commander, w *os.File) {
			defer wg.Done()
			c.Execute()

			// Close the pipe writer so the NEXT command stops waiting for input.
			if w != nil {
				w.Close()
			}
		}(cmd, pipeWriter)
	}

	wg.Wait()
}

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

func createCommand(tokens []string, historyManager *HistoryManager) Commander {
	if len(tokens) == 0 {
		return nil
	}

	switch tokens[0] {
	case "exit":
		return &Exit{}
	case "echo":
		return &Echo{Args: tokens[1:]}
	case "type":
		return &Type{Args: tokens[1:], Manager: historyManager}
	case "pwd":
		return &Pwd{Args: tokens[1:]}
	case "cd":
		return &Cd{Args: tokens[1:]}
	case "history":
		return &History{Args: tokens[1:], Manager: historyManager}
	default:
		return &Executable{Name: tokens[0], Args: tokens[1:]}
	}
}

type BellCompleter struct {
	Completer readline.AutoCompleter
}

func (b *BellCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	newLine, length = b.Completer.Do(line, pos)

	if len(newLine) == 0 {
		fmt.Print("\a")
	}

	return newLine, length
}

func getBinaries() []readline.PrefixCompleterInterface {
	var items []readline.PrefixCompleterInterface

	// Add Builtins
	for _, cmd := range []string{"cd", "echo", "pwd", "exit", "type", "history"} {
		items = append(items, readline.PcItem(cmd))
	}

	// Add System Binaries
	path := os.Getenv("PATH")
	for _, dir := range strings.Split(path, string(os.PathListSeparator)) {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if !file.IsDir() {
				items = append(items, readline.PcItem(file.Name()))
			}
		}
	}
	return items
}

func main() {
	historyManager, err := NewHistoryManager()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing history manager:", err)
		return
	}

	completer := readline.NewPrefixCompleter(getBinaries()...)

	finalCompleter := &DoubleTabCompleter{
		Completer: completer,
	}

	// Initialize the Readline instance
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		HistoryFile:  "", // Disable history to file saving
		AutoComplete: finalCompleter,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing readline:", err)
		return
	}
	defer rl.Close()

	for _, cmd := range historyManager.GetLines() {
		rl.SaveHistory(cmd)
	}

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
		historyManager.Add(command)

		executePipeline(command, historyManager)
		// tokens, stdout, appendStdout, stderr, appendStderr := tokenize(command)
		// handleInput(tokens, historyManager, stdout, appendStdout, stderr, appendStderr)
	}
}
