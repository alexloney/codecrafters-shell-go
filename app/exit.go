package main

import (
	"context"
	"os"
)

type Exit struct {
	StandardIO
}

func (e *Exit) Execute(ctx context.Context) error {
	e.ensureDefaults()

	os.Exit(0)
	return nil
}
func (e *Exit) GetType() string {
	return "exit is a shell builtin"
}
