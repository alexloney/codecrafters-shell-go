package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Unknown struct {
	StandardIO
	Name string
	Args []string
}

func (u *Unknown) isExecutable(path string) bool {
	path, err := exec.LookPath(path)
	if err != nil {
		return false
	}
	return true
}

func (u *Unknown) fileExists(path string) bool {
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
func (u *Unknown) GetPath() string {
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
func (u *Unknown) Execute() error {
	// 1. LookPath checks if the command exists in PATH
	path, err := exec.LookPath(u.Name)
	if err != nil {
		if errors.Is(err, exec.ErrDot) {
			path = u.Name
		} else {
			fmt.Fprintf(u.stderr, "%s: command not found\n", u.Name)
			return nil // Return nil so the shell doesn't crash, just logs error
		}
	}

	// 2. Create the command using the found path
	cmd := exec.Command(path, u.Args...)

	// 3. Connect the I/O
	cmd.Stdin = u.stdin
	cmd.Stdout = u.stdout
	cmd.Stderr = u.stderr

	// 4. Run it
	if err := cmd.Run(); err != nil {
		// cmd.Run() waits for the command to finish automatically
		// If the command exits with non-zero, it returns an error.
		// We usually just print it.
		fmt.Fprintln(u.stderr, "Error executing command:", err)
	}
	return nil
}
func (u *Unknown) GetType() string {
	path, err := exec.LookPath(u.Name)
	if err != nil {
		return fmt.Sprintf("%s: not found", u.Name)
	}
	return fmt.Sprintf("%s is %s", u.Name, path)
}
