package main

import (
	"os"
)

type Exit struct {
}

func (e *Exit) Execute() {
	os.Exit(0)
}
func (e *Exit) GetType() string {
	return "exit is a shell builtin"
}
