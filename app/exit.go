package main

import (
	"io"
	"os"
)

type Exit struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (c *Exit) SetStdin(stdin io.Reader) {
	c.stdin = stdin
}

func (c *Exit) SetStdout(stdout io.Writer) {
	c.stdout = stdout
}

func (c *Exit) SetStderr(stderr io.Writer) {
	c.stderr = stderr
}

func (e *Exit) Execute() error {
	os.Exit(0)
	return nil
}
func (e *Exit) GetType() string {
	return "exit is a shell builtin"
}
