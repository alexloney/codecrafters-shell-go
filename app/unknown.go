package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Unknown struct {
	Commander
	Name string
	Args []string
}

func (u *Unknown) isExecutable(path string) bool {
	path, err := exec.LookPath(path)
	if err != nil {
		return false
	}
	return true
}

func (u *Unknown) fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	// Check if the error is specifically because the file does not exist
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	// For other errors (like permission issues), you might want to handle them differently or return false
	return false
}
func (u *Unknown) GetPath() string {
	path_separator := string(os.PathListSeparator)
	file_separator := string(os.PathSeparator)
	path := os.Getenv("PATH")

	directories := strings.Split(path, path_separator)
	for _, dir := range directories {
		full_path := dir + file_separator + u.Name
		if u.fileExists(full_path) {
			if u.isExecutable(full_path) {
				return filepath.Clean(full_path)
			}
		}
	}
	return ""
}
func (u *Unknown) Execute() {
	path := u.GetPath()
	if path != "" {
		command := exec.Command(u.Name, u.Args...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err := command.Start()
		if err != nil {
			fmt.Println(err)
			return
		}
		command.Wait()
		return
	}
	fmt.Printf("%s: command not found\n", u.Name)
}
func (u *Unknown) GetType() string {
	path := u.GetPath()
	if path != "" {
		return fmt.Sprintf("%s is %s", u.Name, path)
	} else {
		return fmt.Sprintf("%s: not found", u.Name)
	}

}
