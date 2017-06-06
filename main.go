package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// TODO: Lookup the .journal file for config
	// Get the current date.
	year, month, day := time.Now().Date()
	root := fmt.Sprintf("%d-%02d-%02d", year, month, day)
	file := fmt.Sprintf("%s.md", root)
	// Check to see if a directory or file exists.
	if _, err := os.Stat(root); os.IsNotExist(err) {
		// Create a directory (or file if config'd)
		err = os.Mkdir(root, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Create a md file with hard coded title keys.
	f, err := os.OpenFile(filepath.Join(root, file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString("General\n=======\n\nLearn\n=====\n")
	if err != nil {
		log.Fatal(err)
	}
	// TODO: If a previous entry exists, initially copy the new file from that.
	// Look back at most a week.
}
