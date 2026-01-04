package main

import (
	"fmt"

	"github.com/chzyer/readline"
)

type DoubleTabCompleter struct {
	// Embed the real completer (PrefixCompleter or whatever you use)
	Completer readline.AutoCompleter

	// State to track double-tabs
	lastLine string
}

func (d *DoubleTabCompleter) Do(line []rune, pos int) ([][]rune, int) {
	// 1. Get real candidates from the underlying completer
	candidates, length := d.Completer.Do(line, pos)

	// If 0 or 1 match, just behave normally (no beep, no double-tab logic)
	if len(candidates) <= 1 {
		d.lastLine = ""
		return candidates, length
	}

	// 2. Logic for Multiple Matches (>1)
	currentLine := string(line)

	// Case A: This is the SECOND tab press (Input hasn't changed since last tab)
	if currentLine == d.lastLine {
		// Reset state so a 3rd tab doesn't get weird (optional)
		d.lastLine = ""
		// Return ALL candidates so readline displays the list
		return candidates, length
	}

	// Case B: This is the FIRST tab press
	d.lastLine = currentLine
	fmt.Print("\a") // Ring Bell (0x07)

	// 3. "Fill to partial" logic
	// Even though we aren't showing the list, we should still auto-complete
	// as much as possible (the Longest Common Prefix).
	// E.g., if we have "git-a", "git-b", we should fill "git-"
	lcp := d.longestCommonPrefix(candidates)

	if len(lcp) > 0 {
		// Return the LCP as the *only* candidate.
		// Readline will see 1 candidate and auto-fill it, but won't show the menu.
		return [][]rune{lcp}, length
	}

	// If absolutely nothing in common, do nothing
	return nil, 0
}

// Helper to find the common runes shared by all candidates
func (d *DoubleTabCompleter) longestCommonPrefix(candidates [][]rune) []rune {
	if len(candidates) == 0 {
		return nil
	}
	// Start with the first candidate as the prefix
	prefix := candidates[0]

	for _, c := range candidates[1:] {
		// Truncate prefix to the length of the shorter string
		if len(c) < len(prefix) {
			prefix = prefix[:len(c)]
		}
		// Check character by character
		for i := 0; i < len(prefix); i++ {
			if c[i] != prefix[i] {
				// Mismatch found, cut the prefix here
				prefix = prefix[:i]
				break
			}
		}
	}
	return prefix
}
