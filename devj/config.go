package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ifo/dev.journal/entry"
	"github.com/ifo/dev.journal/filesystem"
)

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

	// Let's set some defaults.
	if conf.EditorCommand == "" {
		conf.EditorCommand = "vim"
	}

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
	out := &entry.Journal{Entries: map[entry.EntryName]entry.Entry{}}
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
		out.Entries[entry.EntryName(date)] = e
	}
	return out, nil
}
