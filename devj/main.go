package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ifo/dev.journal/entry"
	"github.com/ifo/dev.journal/filesystem"
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
	case "edit":
		pe := filesystem.Latest()
		if pe == "" {
			fmt.Println("no entry to edit")
			return
		}
		cmd := exec.Command("vim", pe) // TODO: allow the editor to be configured
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func MakeNewEntry() error {
	folder := filesystem.DateString(time.Now())
	file := fmt.Sprintf("%s.md", folder)

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	contents := entry.Default.Export()

	fname := filesystem.Latest()
	if fname != "" {
		bts, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		contents = string(bts)
	}

	f, err := os.OpenFile(filepath.Join(folder, file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(contents)
	if err != nil {
		return err
	}

	return nil
}
