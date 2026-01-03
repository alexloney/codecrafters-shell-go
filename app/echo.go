package main

import (
	"fmt"
	"strings"
)

type Echo struct {
	Args []string
}

func (e *Echo) Execute() {
	fmt.Println(strings.Join(e.Args, " "))
}
func (e *Echo) GetType() string {
	return "echo is a shell builtin"
}
