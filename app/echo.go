package main

import (
	"context"
	"fmt"
	"strings"
)

type Echo struct {
	StandardIO
	Args []string
}

func (e *Echo) Execute(ctx context.Context) error {
	e.ensureDefaults()

	fmt.Fprintln(e.stdout, strings.Join(e.Args, " "))
	return nil
}
func (e *Echo) GetType() string {
	return "echo is a shell builtin"
}
