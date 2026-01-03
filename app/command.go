package main

import "io"

type Commander interface {
	Execute() error
	GetType() string
	SetStdin(io.Reader)
	SetStdout(io.Writer)
	SetStderr(io.Writer)
}
