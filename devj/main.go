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
	if len(os.Args) > 1 && os.Args[1] == "new" {
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

		contents := journal.DefaultEntry.Export()

		fname := journal.PreviousEntry()
		fmt.Println(fname, root, file)
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
}
