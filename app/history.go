package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type History struct {
	Args []string
}

func (h *History) GetHistoryFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/.simple_shell_history", nil
}

func (h *History) ReadHistoryLines() ([]string, error) {
	historyFilePath, err := h.GetHistoryFilePath()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(historyFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func (h *History) DisplayFullHistory() {
	lines, err := h.ReadHistoryLines()
	if err != nil {
		fmt.Println("Error reading history file:", err)
		return
	}
	for i, line := range lines {
		fmt.Printf("    %d  %s\n", i+1, line)
	}
}

func (h *History) DisplayLastNHistory(n int) {
	lines, err := h.ReadHistoryLines()
	if err != nil {
		fmt.Println("Error reading history file:", err)
		return
	}
	start := 0
	if n < len(lines) {
		start = len(lines) - n
	}
	for i := start; i < len(lines); i++ {
		fmt.Printf("    %d  %s\n", i+1, lines[i])
	}
}

func (h *History) Execute() {
	if len(h.Args) > 1 {
		fmt.Println("history: too many arguments")
		return
	} else if len(h.Args) == 0 {
		h.DisplayFullHistory()
		return
	} else if len(h.Args) == 1 {
		n, err := strconv.Atoi(h.Args[0])
		if err != nil || n <= 0 {
			fmt.Println("history: invalid argument:", h.Args[0])
			return
		}
		h.DisplayLastNHistory(n)
		return
	}
}
func (h *History) GetType() string {
	return "history is a shell builtin"
}
