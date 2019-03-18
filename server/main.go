package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ifo/dev.journal/entry"
	"github.com/ifo/dev.journal/filesystem"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
)

const journalDir = "journals"

func main() {
	// Setup
	// Get the server port.
	portStr := os.Getenv("DEVJ_PORT")
	portDefault := 3000
	var err error
	if portStr != "" {
		portDefault, err = strconv.Atoi(portStr)
	}
	if err != nil {
		log.Fatal(err)
	}
	port := flag.Int("port", portDefault, "Port to run the server on")
	user := flag.String("u", os.Getenv("DEVJ_USER"), "The user")
	pass := flag.String("p", os.Getenv("DEVJ_PASSWORD"), "The user's password")

	flag.Parse()

	// TODO: allow for more than one user
	if *user == "" || *pass == "" {
		log.Fatalln("Both user and password must be non-empty")
	}

	// Ensure logging directory exists.
	if err := filesystem.EnsureFolderExists("logs"); err != nil {
		log.Fatalf("error making directory: %v", err)
	}
	f, err := os.OpenFile("logs/server.logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	logger := logrus.New()
	logger.Out = f

	// Ensure the file database directory exists.
	if err := filesystem.EnsureFolderExists(journalDir); err != nil {
		log.Fatalf("error making directory: %v", err)
	}

	r := chi.NewRouter()

	r.Use(CreateBaseContext(*user, *pass))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(lg.RequestLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(Auth)

	r.Post("/", postJournalHandler)

	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}

/*
// Types and Constants
*/

type contextKey string

var userKey = contextKey("user")
var passKey = contextKey("pass")

/*
// Middleware
*/

// CreateBaseContext sets the base context for every request.
func CreateBaseContext(user, pass string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), passKey, pass)
			ctx = context.WithValue(ctx, userKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Auth ensures that a user is authed.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "no credentials provided", 401)
			return
		}
		ctxUser := r.Context().Value(userKey).(string)
		ctxPass := r.Context().Value(passKey).(string)
		if user != ctxUser || pass != ctxPass {
			http.Error(w, "not authorized", 403)
			return
		}
		next.ServeHTTP(w, r)
	})
}

/*
// Handlers
*/

func postJournalHandler(w http.ResponseWriter, r *http.Request) {
	userDir := r.Context().Value(userKey).(string)
	// Ensure userDir exists.
	if err := filesystem.EnsureFolderExists(filepath.Join(journalDir, userDir)); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Decode the entries.
	var journal entry.Journal
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&journal)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer r.Body.Close()

	for dir, e := range journal.Entries {
		newDir := filepath.Join(journalDir, userDir, dir)
		if err := filesystem.EnsureFolderExists(newDir); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		// Export the entry.
		err := ioutil.WriteFile(filepath.Join(newDir, fmt.Sprintf("%s.md", dir)), []byte(e.Export()), 0644)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		// Export any public files.
		for name, contents := range e.PublicFiles {
			err = ioutil.WriteFile(filepath.Join(newDir, name), contents, 0644)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
	}
	// Empty 200 response.
}
