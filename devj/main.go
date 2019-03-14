package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ifo/dev.journal/entry"
	"github.com/ifo/dev.journal/filesystem"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("no command given")
		return
	}

	conf, err := ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(os.Args[1]) {
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
		cmd := exec.Command(conf.EditorCommand, pe)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

	case "export":
		jrn, err := conf.ImportJournal(".")
		if err != nil {
			log.Fatal(err)
		}

		url := os.Args[2]
		user := os.Args[3]
		pass := os.Args[4]
		if user == "" || pass == "" {
			log.Fatal("need both url, user and password")
		}
		if !strings.HasPrefix(url, "https://") {
			log.Fatal(`the url must use https (so must start with "https://")`)
		}

		body, err := json.Marshal(jrn)
		if err != nil {
			log.Fatal(err)
		}
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(user, pass)
		req.Header.Set("Content-Type", "text/plain")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("post finished with status %d\n", resp.StatusCode)

	case "viewconfig":
		fmt.Println(conf)

	default:
		fmt.Println("unknown command")
	}
}

func MakeNewEntry() error {
	folder := filesystem.DateString(time.Now())
	file := fmt.Sprintf("%s.md", folder)

	contents := entry.Default.Export()

	// Overwrite contents with the last journal, to give a better starting journal.
	if fname := filesystem.Latest() {
		bts, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		contents = string(bts)
	}

	if err := filesystem.EnsureFolderExists(folder); err != nil {
		return err
	}

	return filesystem.WriteFile(contents, filepath.Join(folder, file))
}
