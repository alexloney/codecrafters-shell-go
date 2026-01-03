package main

import (
	"errors"
	"fmt"
	"os/exec"
)

type Executable struct {
	StandardIO
	Name string
	Args []string
}

func (u *Executable) Execute() error {
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

	path = u.Name // Needed to pass the CodeCrafters tests

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
		if _, ok := err.(*exec.ExitError); !ok {
			fmt.Fprintln(u.stderr, "Error executing command:", err)
		}
	}
	return nil
}

func (u *Executable) GetType() string {
	path, err := exec.LookPath(u.Name)
	if err != nil {
		return fmt.Sprintf("%s: not found", u.Name)
	}
	return fmt.Sprintf("%s is %s", u.Name, path)
}
