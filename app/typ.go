package main

import (
	"context"
	"fmt"
)

type Type struct {
	StandardIO
	Args    []string
	Manager *HistoryManager
}

func (t *Type) Execute(ctx context.Context) error {
	t.ensureDefaults()

	command := createCommand(t.Args, t.Manager)
	if command != nil {
		fmt.Fprintln(t.stdout, command.GetType())
	}

	return nil
}
func (t *Type) GetType() string {
	return "type is a shell builtin"
}
