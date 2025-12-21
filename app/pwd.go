package main

import (
	"fmt"
	"os"
)

type Pwd struct {
	ICommand
	Args []string
}

func (p Pwd) Execute() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error retrieving current directory:", err)
		return
	}
	fmt.Println(dir)
}
func (p Pwd) GetType() string {
	return "pwd is a shell builtin"
}
