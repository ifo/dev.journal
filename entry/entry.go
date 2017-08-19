package entry

import (
	"fmt"
	"strings"
)

type Style int

const (
	Pound Style = iota
	Underline
)

type Entry struct {
	Sections    []Section           `json:"sections"`
	Style       Style               `json:"style"`
	PublicFiles map[string][]byte   `json:"files"`
	FileNames   map[string]struct{} `json:"-"`
}

type Section struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

var Default = Entry{Style: Pound, Sections: []Section{{Title: "General"}, {Title: "Learn"}}}
var DefaultUnderline = Entry{Style: Underline,
	Sections: []Section{{Title: "General"}, {Title: "Learn"}}}

func (e Entry) Export() string {
	out := ""
	for i, s := range e.Sections {
		if i != 0 {
			out += "\n"
		}
		if e.Style == Pound {
			out += fmt.Sprintf("# %s\n\n%s\n", s.Title, s.Body)
		} else if e.Style == Underline {
			out += fmt.Sprintf("%s\n%s\n\n%s\n", s.Title, strings.Repeat("=", len(s.Title)), s.Body)
		}
	}
	return out
}

func Import(str string) (Entry, error) {
	if len(str) < 3 {
		return Entry{}, fmt.Errorf("entry is empty")
	}

	if str[:2] == "# " {
		return importPoundTitles(str)
	} else if lines := strings.SplitN(str, "\n", 3); len(lines) > 1 && areTitle(lines[0], lines[1]) {
		return importUnderlineTitles(str)
	}
	return Entry{}, fmt.Errorf("entries must start with a title")
}

func importPoundTitles(str string) (Entry, error) {
	lines := strings.Split(str, "\n")
	e := Entry{Style: Pound}
	s := Section{Title: lines[0][2:]} // Remove the starting "# " from the Title.

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

func importUnderlineTitles(str string) (Entry, error) {
	lines := strings.Split(str, "\n")
	e := Entry{Style: Underline}
	s := Section{Title: strings.TrimSpace(lines[0])}

	if len(lines) < 4 {
		s.Body = strings.TrimSpace(strings.Join(lines[2:], "\n"))
		e.Sections = append(e.Sections, s)
		return e, nil
	}

	past, curr, skip := "", lines[2], false
	for _, l := range lines[3:] {
		past = curr
		curr = l

		if skip {
			// We just passed a title, so move the window again.
			skip = false
			continue
		}

		switch {
		// The section is finished; start a new one.
		case areTitle(past, curr):
			s.Body = strings.TrimSpace(s.Body)
			e.Sections = append(e.Sections, s)
			s = Section{Title: strings.TrimSpace(past)}
			skip = true
		default:
			s.Body += "\n" + past
		}
	}

	if !skip {
		s.Body += "\n" + curr
	}
	s.Body = strings.TrimSpace(s.Body)
	e.Sections = append(e.Sections, s)

	return e, nil
}

// Two lines are a title if there is at least 1 non space rune on the first line
// and the 2nd line is more than 1 "=" sign, and entirely "=" signs.
func areTitle(line1, line2 string) bool {
	return len(strings.Replace(line1, " ", "", -1)) > 0 &&
		len(line2) > 0 &&
		len(strings.Replace(line2, "=", "", -1)) == 0
}

func (e *Entry) ImportFiles(pubSections map[string]struct{}, readFile func(string) ([]byte, error)) error {
	publicFiles := e.publicFileList(pubSections)
	fileMap := map[string][]byte{}
	for _, f := range publicFiles {
		data, err := readFile(f)
		if err != nil {
			return err
		}
		fileMap[f] = data
	}
	e.PublicFiles = fileMap
	return nil
}

func (e Entry) publicFileList(pubSections map[string]struct{}) []string {
	var expFileName []string
	for _, s := range e.Sections {
		if _, ok := pubSections[s.Title]; !ok {
			continue
		}
		for name, _ := range e.FileNames {
			if strings.Contains(s.Body, name) {
				expFileName = append(expFileName, name)
			}
		}
	}
	return expFileName
}
