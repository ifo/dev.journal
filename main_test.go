package main

import "testing"

const testFail = "Actual: %q Expected: %q Case: %q"

func TestEntryExport(t *testing.T) {
	cases := map[string]struct {
		E   Entry
		Out string
	}{
		"1": {E: defaultEntry, Out: "General\n=======\n\n\n\nLearn\n=====\n\n\n"},
		"2": {E: Entry{Sections: nil}, Out: ""},
	}

	for id, c := range cases {
		exp := c.E.Export()
		if exp != c.Out {
			t.Errorf(testFail, exp, c.Out, id)
		}
	}
}
