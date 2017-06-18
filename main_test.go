package main

import (
	"fmt"
	"reflect"
	"testing"
)

const testFail = `Actual: "%+v" Expected: "%+v" Case: %q`

func TestEntryExport(t *testing.T) {
	cases := map[string]struct {
		E   Entry
		Out string
	}{
		"default entry": {E: defaultEntry, Out: "# General\n\n\n\n# Learn\n\n\n"},
		"empty entry":   {E: Entry{Sections: nil}, Out: ""},
		"section with body": {
			E:   Entry{Sections: []Section{{Title: "Five", Body: "teve"}}},
			Out: "# Five\n\nteve\n"},
		"two sections with body": {
			E:   Entry{Sections: []Section{{Title: "Five", Body: "teve"}, {Title: "Four", Body: "too"}}},
			Out: "# Five\n\nteve\n\n# Four\n\ntoo\n"},
	}

	for id, c := range cases {
		exp := c.E.Export()
		if exp != c.Out {
			t.Errorf(testFail, exp, c.Out, id)
		}
	}
}

func TestEntryImport(t *testing.T) {
	cases := map[string]struct {
		In  string
		E   Entry
		Err error
	}{
		"default entry": {In: "# General\n\n\n\n# Learn\n\n\n", E: defaultEntry, Err: nil},
		"empty entry":   {In: "", E: Entry{}, Err: fmt.Errorf("entry is empty")},
		"no title": {
			In:  "not a title",
			E:   Entry{},
			Err: fmt.Errorf("entries must start with a title")},
		"section with body": {
			In:  "# Five\n\nteve\n",
			E:   Entry{Sections: []Section{{Title: "Five", Body: "teve"}}},
			Err: nil},
		"two sections with body": {
			In:  "# Five\n\nteve\n\n# Four\n\ntoo\n",
			E:   Entry{Sections: []Section{{Title: "Five", Body: "teve"}, {Title: "Four", Body: "too"}}},
			Err: nil},
		"repeated empty entries": {
			In:  "# 1\n# 2\n# 3",
			E:   Entry{Sections: []Section{{Title: "1"}, {Title: "2"}, {Title: "3"}}},
			Err: nil},
		"multiline body": {
			In:  "# multi\n\nmultiple\nlines",
			E:   Entry{Sections: []Section{{Title: "multi", Body: "multiple\nlines"}}},
			Err: nil},
	}

	for id, c := range cases {
		ent, err := Import(c.In)
		if !reflect.DeepEqual(ent, c.E) {
			t.Errorf(testFail, ent, c.E, id)
		}
		if !errorEqual(err, c.Err) {
			t.Errorf(testFail, err, c.Err, id)
		}
	}
}

func errorEqual(e1, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	} else if e1 == nil && e2 != nil {
		return false
	} else if e2 == nil && e1 != nil {
		return false
	}
	return e1.Error() == e2.Error()
}
