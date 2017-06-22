package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ifo/dev.journal"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("no command given")
		return
	}
	switch os.Args[1] {
	case "new":
		err := MakeNewEntry()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("new entry created")
	}
}

func MakeNewEntry() error {
	year, month, day := time.Now().Date()
	folder := fmt.Sprintf("%d-%02d-%02d", year, month, day)
	file := fmt.Sprintf("%s.md", folder)

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(filepath.Join(folder, file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	contents := journal.DefaultEntry.Export()

	fname := journal.PreviousEntry()
	if fname != "" {
		bts, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		contents = string(bts)
	}

	_, err = f.WriteString(contents)
	if err != nil {
		return err
	}

	return nil
}
