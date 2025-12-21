package main

import (
	"fmt"
	"os"
)

type Cd struct {
	ICommand
	Args []string
}

func (c Cd) DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false // Path does not exist
	}
	// If there's no error, check if it's a directory
	return err == nil && info.IsDir()
}

func (c Cd) Execute() {
	// No path specified, return without changing directory
	if len(c.Args) == 0 {
		return
	}

	// Too many arguments, return without changing directory
	if len(c.Args) > 1 {
		return
	}

	current_dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error retrieving current directory:", err)
		return
	}

	// Special character, change to home directory
	if c.Args[0] == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error retrieving home directory:", err)
			return
		}
		err = os.Chdir(homeDir)
		if err != nil {
			fmt.Println("Error changing directory to home:", err)
		}
		os.Setenv("OLDPWD", current_dir)
		return
	}

	if c.Args[0] == "-" {
		prevDir := os.Getenv("OLDPWD")
		if prevDir == "" {
			fmt.Println("OLDPWD not set")
			return
		}
		err := os.Chdir(prevDir)
		if err != nil {
			fmt.Println("Error changing to previous directory:", err)
		}
		os.Setenv("OLDPWD", current_dir)
		return
	}

	if c.DirExists(c.Args[0]) {
		err := os.Chdir(c.Args[0])
		if err != nil {
			fmt.Println("Error changing directory:", err)
			return
		}
		os.Setenv("OLDPWD", current_dir)
	} else {
		fmt.Println("cd: " + c.Args[0] + ": No such file or directory")
	}

}
func (c Cd) GetType() string {
	return "cd is a shell builtin"
}
