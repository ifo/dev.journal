package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
	case "edit":
		pe := journal.LatestEntry()
		fmt.Println(pe)
		if pe == "" {
			fmt.Println("no entry to edit")
			return
		}
		cmd := exec.Command("vim", pe)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func MakeNewEntry() error {
	folder := journal.DateString(time.Now())
	file := fmt.Sprintf("%s.md", folder)

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	contents := journal.DefaultEntry.Export()

	fname := journal.LatestEntry()
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
