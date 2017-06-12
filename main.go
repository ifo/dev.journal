package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// TODO: Lookup the .journal file for config

	year, month, day := time.Now().Date()
	root := fmt.Sprintf("%d-%02d-%02d", year, month, day)
	file := fmt.Sprintf("%s.md", root)

	if _, err := os.Stat(root); os.IsNotExist(err) {
		err = os.Mkdir(root, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	f, err := os.OpenFile(filepath.Join(root, file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	contents := defaultEntry.Export()

	fname := previousEntry()
	if fname != "" {
		bts, err := ioutil.ReadFile(fname)
		if err != nil {
			log.Fatal(err)
		}
		contents = string(bts)
	}

	_, err = f.WriteString(contents)
	if err != nil {
		log.Fatal(err)
	}
}

type Entry struct {
	Sections []Section
}

type Section struct {
	Title string
	Body  string
}

var defaultEntry = Entry{Sections: []Section{{Title: "General"}, {Title: "Learn"}}}

func (e Entry) Export() string {
	out := ""
	for i, s := range e.Sections {
		if i != 0 {
			out += "\n"
		}
		out += fmt.Sprintf("%s\n%s\n\n%s\n", s.Title, strings.Repeat("=", len(s.Title)), s.Body)
	}
	return out
}

func Import() (Entry, error) {
	return Entry{}, nil
}

func previousEntry() string {
	today := time.Now()
	for i := 1; i <= 7; i++ {
		y, m, d := today.AddDate(0, 0, -i).Date()
		root := fmt.Sprintf("%d-%02d-%02d", y, m, d)
		file := filepath.Join(root, fmt.Sprintf("%s.md", root))
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return file
		}
	}
	return ""
}
