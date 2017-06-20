package journal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Entry struct {
	Sections []Section
}

type Section struct {
	Title string
	Body  string
}

var DefaultEntry = Entry{Sections: []Section{{Title: "General"}, {Title: "Learn"}}}

func (e Entry) Export() string {
	out := ""
	for i, s := range e.Sections {
		if i != 0 {
			out += "\n"
		}
		out += fmt.Sprintf("# %s\n\n%s\n", s.Title, s.Body)
	}
	return out
}

func Import(str string) (Entry, error) {
	if str == "" {
		return Entry{}, fmt.Errorf("entry is empty")
	}
	lines := strings.Split(str, "\n")
	if len(strings.Replace(lines[0], " ", "", -1)) < 2 || lines[0][:2] != "# " {
		return Entry{}, fmt.Errorf("entries must start with a title")
	}

	e := Entry{}
	s := Section{Title: lines[0][2:]} // Remove the starting "# " from the Title

	for _, l := range lines[1:] {
		switch {
		// The section is finished; start a new one.
		case len(strings.Replace(l, " ", "", -1)) >= 2 && l[:2] == "# ":
			s.Body = strings.TrimSpace(s.Body)
			e.Sections = append(e.Sections, s)
			s = Section{Title: l[2:]}
		default:
			s.Body += "\n" + l
		}
	}

	s.Body = strings.TrimSpace(s.Body)
	e.Sections = append(e.Sections, s)

	return e, nil
}

func PreviousEntry() string {
	today := time.Now()
	for i := 1; i <= 7; i++ {
		y, m, d := today.AddDate(0, 0, -i).Date()
		root := fmt.Sprintf("%d-%02d-%02d", y, m, d)
		file := filepath.Join(root, fmt.Sprintf("%s.md", root))
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return file
		}
	}
	return ""
}
