package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ifo/dev.journal/entry"
)

func TestPostJournalHandler(t *testing.T) {
	defer resetFileSystem()

	// Overwrite the functions
	folderCreator = FakeCreateFolder
	fileWriter = FakeWriteFile

	recorder := httptest.NewRecorder()
	bts, _ := json.Marshal(entry.Journal{Entries: map[string]entry.Entry{"2019-03-19": entry.Default}})
	body := bytes.NewReader(bts)
	request, _ := http.NewRequest(http.MethodPost, "/", body)
	request = request.WithContext(context.WithValue(request.Context(), userKey, "user"))

	postJournalHandler(recorder, request)

	file := fileSystem["journals/user/2019-03-19/2019-03-19.md"]
	if file != entry.Default.Export() {
		t.Errorf("Got %s,\n\nexpected %s\n", file, entry.Default.Export())
	}
}

var fileSystem = map[string]string{}

func resetFileSystem() {
	fileSystem = map[string]string{}
}

func FakeCreateFolder(folder string) error {
	return nil
}

func FakeWriteFile(path string, s string) error {
	fileSystem[path] = s
	return nil
}
