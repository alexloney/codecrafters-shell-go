package main

import (
	"fmt"
	"io"
)

type Type struct {
	Args   []string
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (c *Type) SetStdin(stdin io.Reader) {
	c.stdin = stdin
}

func (c *Type) SetStdout(stdout io.Writer) {
	c.stdout = stdout
}

func (c *Type) SetStderr(stderr io.Writer) {
	c.stderr = stderr
}

func (t *Type) Execute() error {
	command := createCommand(t.Args)
	fmt.Fprintln(t.stdout, command.GetType())
	return nil
}
func (t *Type) GetType() string {
	return "type is a shell builtin"
}
