package main

import (
	"bytes"
	"encoding/json"
	"flag"
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
		if err := EditEntry(conf); err != nil {
			log.Fatal(err)
		}

	case "export":
		if err := ExportJournal(conf); err != nil {
			log.Fatal(err)
		}
		fmt.Println("journal export complete")

	case "viewconfig":
		fmt.Println(conf)

	default:
		fmt.Println("unknown command")
	}
}

func ExportJournal(conf *Config) error {
	jrn, err := conf.ImportJournal(".")
	if err != nil {
		return err
	}

	var url, user, pass string
	flag.StringVar(&url, "url", "", "url to send the journal to")
	flag.StringVar(&user, "user", "", "username")
	flag.StringVar(&pass, "pass", "", "password")
	flag.Parse()

	if user == "" || pass == "" {
		log.Fatal("need both url, user and password")
	}
	if !strings.HasPrefix(url, "https://") {
		log.Fatal(`the url must use https (so must start with "https://")`)
	}

	body, err := json.Marshal(jrn)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.SetBasicAuth(user, pass)
	req.Header.Set("Content-Type", "text/plain")
	if resp, err := http.DefaultClient.Do(req); err != nil {
		return err
	} else {
		fmt.Printf("post finished with status %d\n", resp.StatusCode)
	}

	return nil
}

func MakeNewEntry() error {
	folder := filesystem.DateString(time.Now())
	file := fmt.Sprintf("%s.md", folder)

	contents := entry.Default.Export()

	// Overwrite contents with the last journal, to give a better starting journal.
	if fname := filesystem.Latest(); fname != "" {
		bts, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		contents = string(bts)
	}

	if err := filesystem.EnsureFolderExists(folder); err != nil {
		return err
	}

	return filesystem.SafeWriteFile(filepath.Join(folder, file), []byte(contents))
}

func EditEntry(conf *Config) error {
	pe := filesystem.Latest()
	if pe == "" {
		return fmt.Errorf("no entry to edit")
	}
	cmd := exec.Command(conf.EditorCommand, pe)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
