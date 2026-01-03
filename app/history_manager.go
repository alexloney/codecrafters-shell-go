package main

import (
	"bufio"
	"os"
	"sync"
)

type HistoryManager struct {
	filePath       string
	lines          []string
	lastAppendLine int
	mu             sync.Mutex // Protects against concurrent access, if needed later
}

func (h *HistoryManager) GetLinesSinceLastAppend() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	lines := h.lines[h.lastAppendLine:]
	h.lastAppendLine = len(h.lines)
	return lines
}

func GetHistoryFilePath() (string, error) {
	// Use the default HISTFILE location
	historyFile := os.Getenv("HISTFILE")
	if historyFile != "" {
		return historyFile, nil
	}

	// If HISTFILE is not set, default to ~/.simple_shell_history
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/.simple_shell_history", nil
}

func NewHistoryManager() (*HistoryManager, error) {
	path, err := GetHistoryFilePath()
	if err != nil {
		return nil, err
	}

	hm := &HistoryManager{
		filePath: path,
		lines:    []string{},
	}

	// Pre-load existing history into memory
	if err := hm.load(); err != nil {
		// If file doesn't exist, that's fine, we start empty
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return hm, nil
}

// load reads the file into the in-memory slice
func (h *HistoryManager) load() error {
	f, err := os.Open(h.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		h.lines = append(h.lines, scanner.Text())
	}
	return nil
}

func (h *HistoryManager) clear() error {
	os.Remove(h.filePath)
	h.lines = []string{}
	h.lastAppendLine = 0
	return nil
}

// Add appends to memory AND writes to disk immediately
func (h *HistoryManager) Add(command string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 1. Update In-Memory
	h.lines = append(h.lines, command)

	// Not writing to disk because the Readline library already does this.
	// If we also write to disk here, we end up with duplicate entries.
}

// GetLines returns the current history (fast, from memory)
func (h *HistoryManager) GetLines() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	// Return a copy to be safe, or just the slice if you trust callers not to modify it
	return h.lines
}
