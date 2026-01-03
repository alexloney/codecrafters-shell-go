package main

import (
	"fmt"
	"os"
)

type Cat struct {
	Args []string
}

func (c *Cat) DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false // Path does not exist
	}
	// If there's no error, check if it's a directory
	return err == nil && info.IsDir()
}

func (c *Cat) Execute() {
	for _, fileName := range c.Args {
		data, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cat: "+fileName+": No such file or directory")
			continue
		}
		fmt.Print(string(data))
	}
}
func (c *Cat) GetType() string {
	return "cat is a shell builtin"
}
