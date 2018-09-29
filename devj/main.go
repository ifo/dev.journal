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

type Config struct {
	PublicSections map[string]struct{} `json:"public_sections"`
	EditorCommand  string              `json:"editor_command"`
}

type lenientConfig struct {
	PublicSections map[string]interface{} `json:"public_sections"`
	EditorCommand  string                 `json:"editor_command"`
}

func ReadConfig() (*Config, error) {
	bts, err := ioutil.ReadFile(".devj")
	if err != nil {
		return nil, err
	}

	var conf *Config
	err = json.Unmarshal(bts, &conf)
	return conf, err
}

func (c *Config) UnmarshalJSON(buf []byte) error {
	lc := lenientConfig{}
	err := json.Unmarshal(buf, &lc)
	if err != nil {
		return err
	}
	c.EditorCommand = "vim"
	if lc.EditorCommand != "" {
		c.EditorCommand = lc.EditorCommand
	}
	c.PublicSections = map[string]struct{}{}
	for k, _ := range lc.PublicSections {
		c.PublicSections[k] = struct{}{}
	}
	return nil
}

func (c *Config) ImportJournal(basePath string) (*entry.Journal, error) {
	entries, err := filesystem.ListDirs(basePath)
	if err != nil {
		return nil, err
	}
	out := &entry.Journal{Entries: map[string]entry.Entry{}}
	for _, date := range entries {
		entryDir := filepath.Join(basePath, date)
		rawEntry, err := filesystem.ReadFile(filepath.Join(entryDir, fmt.Sprintf("%s.md", date)))
		if err != nil {
			return nil, err
		}

		e, err := entry.ImportPublic(string(rawEntry), c.PublicSections)
		if err != nil {
			return nil, err
		}

		files, err := filesystem.ListFiles(entryDir)
		if err != nil {
			return nil, err
		}
		e.FileNames = map[string]struct{}{}
		for _, name := range files {
			e.FileNames[name] = struct{}{}
		}

		err = e.ImportFiles(c.PublicSections, entryDir, filesystem.ReadFile)
		if err != nil {
			return nil, err
		}
		out.Entries[date] = e
	}
	return out, nil
}
