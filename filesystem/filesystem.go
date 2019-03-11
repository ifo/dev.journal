package filesystem

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"time"
)

func Latest() string {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return ""
	}

	regex := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	lastDirMatch := ""
	for _, f := range files {
		if f.IsDir() && regex.MatchString(f.Name()) {
			lastDirMatch = f.Name()
		}
	}
	return filepath.Join(lastDirMatch, fmt.Sprintf("%s.md", lastDirMatch))
}

func DateString(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", y, m, d)
}

func ListFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fileList []string
	for _, f := range files {
		fileList = append(fileList, f.Name())
	}
	return fileList, nil
}

func ListDirs(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fileList []string
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		fileList = append(fileList, f.Name())
	}
	return fileList, nil
}

func ReadFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}
