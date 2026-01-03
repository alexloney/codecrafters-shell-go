package main

import (
	"fmt"
	"io"
	"os"
)

type Pwd struct {
	Args   []string
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (c *Pwd) SetStdin(stdin io.Reader) {
	c.stdin = stdin
}

func (c *Pwd) SetStdout(stdout io.Writer) {
	c.stdout = stdout
}

func (c *Pwd) SetStderr(stderr io.Writer) {
	c.stderr = stderr
}

func (p *Pwd) Execute() error {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(p.stderr, "Error retrieving current directory:", err)
		return err
	}
	fmt.Fprintln(p.stdout, dir)
	return nil
}
func (p *Pwd) GetType() string {
	return "pwd is a shell builtin"
}
