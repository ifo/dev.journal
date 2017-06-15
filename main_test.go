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
		"default entry": {E: defaultEntry, Out: "General\n=======\n\n\n\nLearn\n=====\n\n\n"},
		"empty entry":   {E: Entry{Sections: nil}, Out: ""},
		"section with body": {
			E:   Entry{Sections: []Section{{Title: "Five", Body: "teve"}}},
			Out: "Five\n====\n\nteve\n"},
		"two sections with body": {
			E:   Entry{Sections: []Section{{Title: "Five", Body: "teve"}, {Title: "Four", Body: "too"}}},
			Out: "Five\n====\n\nteve\n\nFour\n====\n\ntoo\n"},
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
		"default entry": {In: "General\n=======\n\n\n\nLearn\n=====\n\n\n", E: defaultEntry, Err: nil},
		"empty entry":   {In: "", E: Entry{}, Err: fmt.Errorf("entry is empty")},
		"section with body": {
			In:  "Five\n====\n\nteve\n",
			E:   Entry{Sections: []Section{{Title: "Five", Body: "teve"}}},
			Err: nil},
		"two sections with body": {
			In:  "Five\n====\n\nteve\n\nFour\n====\n\ntoo\n",
			E:   Entry{Sections: []Section{{Title: "Five", Body: "teve"}, {Title: "Four", Body: "too"}}},
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
