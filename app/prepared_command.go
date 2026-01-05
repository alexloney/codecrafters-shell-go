package main

// Helper struct to track if the user manually redirected output
type PreparedCommand struct {
	Command          Commander
	StdoutRedirected bool // true if user used ">" or ">>"
	StderrRedirected bool // true if user used "2>" or "2>>"
}
