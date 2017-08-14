package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func Latest() string {
	// TODO: handle the following cases:
	// - There has not been an entry for more than 7 days
	// - There are no files yet
	today := time.Now()
	for i := 0; i <= 7; i++ {
		folder := DateString(today.AddDate(0, 0, -i))
		file := filepath.Join(folder, fmt.Sprintf("%s.md", folder))
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return file
		}
	}
	return ""
}

func DateString(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", y, m, d)
}
