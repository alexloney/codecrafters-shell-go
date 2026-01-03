package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type History struct {
	StandardIO
	Args    []string
	Manager *HistoryManager
}

func (h *History) DisplayFullHistory() {
	lines := h.Manager.GetLines()
	for i, line := range lines {
		fmt.Fprintf(h.stdout, "    %d  %s\n", i+1, line)
	}
}

func (h *History) DisplayLastNHistory(n int) {
	// This is not optimal as it reads the entire file into memory,
	// but works for this simple case. An improvement would be to
	// see the end of the file and read backward or use a ring buffer.
	lines := h.Manager.GetLines()
	start := 0
	if n < len(lines) {
		start = len(lines) - n
	}
	for i := start; i < len(lines); i++ {
		fmt.Fprintf(h.stdout, "    %d  %s\n", i+1, lines[i])
	}
}

func (h *History) AppendHistoryFile(filename string) error {
	in, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer in.Close()
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		h.Manager.Add(scanner.Text())
	}
	return nil
}

func (h *History) AppendHistoryToOutputFile(filename string) error {
	lines := h.Manager.GetLinesSinceLastAppend()

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines {
		fmt.Println("Appending line to file:", line)
		_, err := f.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *History) WriteHistoryFile(filename string) error {
	lines := h.Manager.GetLines()

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines {
		fmt.Println("Writing line to file:", line)
		_, err := f.WriteString(line + "\n")
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
	} else if len(h.Args) == 2 && h.Args[0] == "-a" {
		err := h.AppendHistoryToOutputFile(h.Args[1])
		if err != nil {
			fmt.Fprintln(h.stderr, "history: error reading file:", err)
		}
	}

	return nil
}
func (h *History) GetType() string {
	return "history is a shell builtin"
}
