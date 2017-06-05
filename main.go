package main

import (
	"fmt"
	"time"
)

func main() {
	// Get the current date.
	year, month, day := time.Now().Date()
	root := fmt.Sprintf("%d-%02d-%02d", year, month, day)
	file := fmt.Sprintf("%s.md", root)
	// Check to see if a directory or file exists.
	// TODO: Lookup the .journal file for config
	// Create a directory (or file if config'd)
	// Create a md file with hard coded title keys.
	// If a previous entry exists, initially copy the new file from that.
}
