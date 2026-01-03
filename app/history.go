package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type History struct {
	StandardIO
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
		fmt.Fprintln(h.stderr, "Error reading history file:", err)
		return
	}
	for i, line := range lines {
		fmt.Fprintf(h.stdout, "    %d  %s\n", i+1, line)
	}
}

func (h *History) DisplayLastNHistory(n int) {
	// This is not optimal as it reads the entire file into memory,
	// but works for this simple case. An improvement would be to
	// see the end of the file and read backward or use a ring buffer.
	lines, err := h.ReadHistoryLines()
	if err != nil {
		fmt.Fprintln(h.stderr, "Error reading history file:", err)
		return
	}
	start := 0
	if n < len(lines) {
		start = len(lines) - n
	}
	for i := start; i < len(lines); i++ {
		fmt.Fprintf(h.stdout, "    %d  %s\n", i+1, lines[i])
	}
}

func (h *History) AppendHistoryFile(filename string) error {
	historyFilePath, err := h.GetHistoryFilePath()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(historyFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	in, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer in.Close()
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		_, err := f.WriteString(scanner.Text() + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *History) WriteHistoryFile(filename string) error {
	historyFilePath, err := h.GetHistoryFilePath()
	if err != nil {
		return err
	}
	in, err := os.OpenFile(historyFilePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer in.Close()

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		_, err := f.WriteString(scanner.Text() + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *History) Execute() error {
	if len(h.Args) == 0 {
		h.DisplayFullHistory()
		return nil
	} else if len(h.Args) == 1 {
		n, err := strconv.Atoi(h.Args[0])
		if err != nil || n <= 0 {
			fmt.Fprintln(h.stderr, "history: invalid argument:", h.Args[0])
			return nil
		}
		h.DisplayLastNHistory(n)
		return nil
	} else if len(h.Args) == 2 && h.Args[0] == "-r" {
		err := h.AppendHistoryFile(h.Args[1])
		if err != nil {
			fmt.Fprintln(h.stderr, "history: error reading file:", err)
		}
	} else if len(h.Args) == 2 && h.Args[0] == "-w" {
		err := h.WriteHistoryFile(h.Args[1])
		if err != nil {
			fmt.Fprintln(h.stderr, "history: error reading file:", err)
		}
	}

	return nil
}
func (h *History) GetType() string {
	return "history is a shell builtin"
}
