package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Inteface
type ICommand interface {
	Execute()
	GetType() string
}

type Exit struct {
	ICommand
}

func (e Exit) Execute() {
	os.Exit(0)
}
func (e Exit) GetType() string {
	return "exit is a shell builtin"
}

type Echo struct {
	ICommand
	Args []string
}

func (e Echo) Execute() {
	fmt.Println(strings.Join(e.Args, " "))
}
func (e Echo) GetType() string {
	return "echo is a shell builtin"
}

type Type struct {
	ICommand
	Args []string
}

func (t Type) Execute() {
	command := createCommand(t.Args)
	fmt.Println(command.GetType())
}
func (t Type) GetType() string {
	return "type is a shell builtin"
}

type Unknown struct {
	ICommand
	Name string
	Args []string
}

func (u Unknown) isExecutable(path string) bool {
	path, err := exec.LookPath(path)
	if err != nil {
		return false
	}
	return true
}

func (u Unknown) fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	// Check if the error is specifically because the file does not exist
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	// For other errors (like permission issues), you might want to handle them differently or return false
	return false
}
func (u Unknown) GetPath() string {
	path_separator := string(os.PathListSeparator)
	file_separator := string(os.PathSeparator)
	path := os.Getenv("PATH")

	directories := strings.Split(path, path_separator)
	for _, dir := range directories {
		full_path := dir + file_separator + u.Name
		if u.fileExists(full_path) {
			if u.isExecutable(full_path) {
				return filepath.Clean(full_path)
			}
		}
	}
	return ""
}
func (u Unknown) Execute() {
	path := u.GetPath()
	if path != "" {
		command := exec.Command(u.Name, u.Args...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err := command.Start()
		if err != nil {
			fmt.Println(err)
			return
		}
		command.Wait()
		return
	}
	fmt.Printf("%s: command not found\n", u.Name)
}
func (u Unknown) GetType() string {
	path := u.GetPath()
	if path != "" {
		return fmt.Sprintf("%s is %s", u.Name, path)
	} else {
		return fmt.Sprintf("%s: not found", u.Name)
	}

}

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
