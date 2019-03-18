package filesystem

import (
	"fmt"
	"io/ioutil"
	"os"
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

	if lastDirMatch != "" {
		return filepath.Join(lastDirMatch, fmt.Sprintf("%s.md", lastDirMatch))
	}
	return ""
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

func EnsureFolderExists(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func SafeWriteFile(path string, s string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("file at path: %s already exists", path)
	}

	return ioutil.WriteFile(path, []byte(s), 0644)
}
