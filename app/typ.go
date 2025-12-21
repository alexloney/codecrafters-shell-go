package main

import (
	"fmt"
)

type Type struct {
	ICommand
	Args []string
}

func (t Type) Execute() {
	command := createCommand(t.Args)
	fmt.Println(command.GetType())
}
func (t Type) GetType() string {
	return "type is a shell builtin"
}
