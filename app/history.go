package main

import (
	"bufio"
	"fmt"
	"os"
)

type History struct {
	ICommand
	Args []string
}

func (h History) Execute() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error retrieving home directory:", err)
		return
	}

	f, err := os.Open(homeDir + "/.simple_shell_history")
	if err != nil {
		fmt.Println("Error opening history file:", err)
		return
	}
	// Ensure the file is closed when the function returns
	defer f.Close()

	// var lines []string
	scanner := bufio.NewScanner(f)
	// The scanner defaults to splitting by newlines (bufio.ScanLines)
	i := 1
	for scanner.Scan() {
		fmt.Printf("    %d  %s\n", i, scanner.Text())
		i++
		// lines = append(lines, scanner.Text())
	}

	// Check for errors during the scan process
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading history file:", err)
		return
	}
}
func (h History) GetType() string {
	return "history is a shell builtin"
}
