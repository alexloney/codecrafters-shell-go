package main

import (
	"fmt"
	"io"
	"strings"
)

type Echo struct {
	Args   []string
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (c *Echo) SetStdin(stdin io.Reader) {
	c.stdin = stdin
}

func (c *Echo) SetStdout(stdout io.Writer) {
	c.stdout = stdout
}

func (c *Echo) SetStderr(stderr io.Writer) {
	c.stderr = stderr
}

func (e *Echo) Execute() error {
	fmt.Fprintln(e.stdout, strings.Join(e.Args, " "))
	return nil
}
func (e *Echo) GetType() string {
	return "echo is a shell builtin"
}
