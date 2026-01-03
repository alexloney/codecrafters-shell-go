package main

import "io"

type StandardIO struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
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
