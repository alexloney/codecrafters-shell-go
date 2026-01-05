package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
)

type Executable struct {
	StandardIO
	Name string
	Args []string
}

func (u *Executable) Execute(ctx context.Context) error {
	u.ensureDefaults()

	// LookPath checks if the command exists in PATH
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

	// If ctx is cancelled (Ctrl+C), this command will be killed automatically.
	cmd := exec.CommandContext(ctx, path, u.Args...)

	cmd.Stdin = u.stdin
	cmd.Stdout = u.stdout
	cmd.Stderr = u.stderr

	if err := cmd.Run(); err != nil {
		// Suppress "signal: killed" error if we caused it via Ctrl+C
		if ctx.Err() == context.Canceled {
			return nil
		}
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
