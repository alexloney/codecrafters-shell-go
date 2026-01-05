package main

import (
	"io"
	"os"
)

type StandardIO struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (s *StandardIO) ensureDefaults() {
	if s.stdin == nil {
		s.stdin = os.Stdin
	}
	if s.stdout == nil {
		s.stdout = os.Stdout
	}
	if s.stderr == nil {
		s.stderr = os.Stderr
	}
}

func (s *StandardIO) SetStdin(stdin io.Reader) {
	s.stdin = stdin
}

func (s *StandardIO) SetStdout(stdout io.Writer) {
	s.stdout = stdout
}

func (s *StandardIO) SetStderr(stderr io.Writer) {
	s.stderr = stderr
}
