package main

import "testing"

const testFail = "Actual: %q Expected: %q Case: %q"

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
