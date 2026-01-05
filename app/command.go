package main

import (
	"context"
	"io"
)

type Commander interface {
	Execute(ctx context.Context) error
	GetType() string
	SetStdin(io.Reader)
	SetStdout(io.Writer)
	SetStderr(io.Writer)
}
