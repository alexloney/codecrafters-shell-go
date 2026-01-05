package main

import (
	"fmt"
)

type Type struct {
	StandardIO
	Args    []string
	Manager *HistoryManager
}

func (t *Type) Execute() error {
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
