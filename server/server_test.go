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

func TestCreateBaseContext(t *testing.T) {
	user := "user"
	pass := "pass"

	called := false
	contextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxUser := r.Context().Value(userKey).(string)
		if ctxUser != user {
			t.Errorf("got %v expected %s", ctxUser, user)
		}

		ctxPass := r.Context().Value(passKey).(string)
		if ctxPass != pass {
			t.Errorf("got %v expected %s", ctxPass, pass)
		}

		called = true
	})

	ts := httptest.NewServer(
		CreateBaseContext(user, pass)(contextHandler))
	defer ts.Close()

	// Make a request call the tests.
	req, _ := http.NewRequest(http.MethodPost, ts.URL, nil)
	http.DefaultClient.Do(req)

	// Ensure things actually ran.
	if !called {
		t.Errorf("context checking tests were not run")
	}
}

func TestAuth(t *testing.T) {
	tests := []struct {
		BasicAuth  bool
		User       string
		Pass       string
		StatusCode int
	}{
		{BasicAuth: true, User: "user", Pass: "notpass", StatusCode: http.StatusForbidden},
		{BasicAuth: true, User: "notuser", Pass: "pass", StatusCode: http.StatusForbidden},
		{BasicAuth: true, User: "notuser", Pass: "notpass", StatusCode: http.StatusForbidden},
		{BasicAuth: false, User: "", Pass: "", StatusCode: http.StatusUnauthorized},
		{BasicAuth: true, User: "user", Pass: "pass", StatusCode: http.StatusOK},
	}

	user := "user"
	pass := "pass"

	ts := httptest.NewServer(
		CreateBaseContext(user, pass)(Auth(GetEmptyHandler())))
	defer ts.Close()

	for _, test := range tests {
		req, _ := http.NewRequest(http.MethodPost, ts.URL, nil)
		if test.BasicAuth {
			req.SetBasicAuth(test.User, test.Pass)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != test.StatusCode {
			t.Errorf("got %d status, expected %d\n", resp.StatusCode, test.StatusCode)
		}
	}
}

func TestPostJournalHandler(t *testing.T) {
	defer resetFileSystem()

	// Overwrite the functions
	folderCreator = FakeCreateFolder
	fileWriter = FakeWriteFile

	recorder := httptest.NewRecorder()
	bts, _ := json.Marshal(entry.Journal{Entries: map[entry.EntryName]entry.Entry{"2019-03-19": entry.Default}})
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

func GetEmptyHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
