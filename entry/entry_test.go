package entry

import (
	"fmt"
	"reflect"
	"testing"
)

const testFail = `Actual: "%+v" Expected: "%+v" Case: %q`

func TestImport(t *testing.T) {
	tests := map[string]struct {
		In  string
		E   Entry
		Err error
	}{
		// Pound titles.
		"default entry": {In: "# Do\n\n\n\n# Learn\n\n\n", E: Default, Err: nil},
		"empty entry":   {In: "", E: Entry{}, Err: fmt.Errorf("entry is empty")},
		"single rune":   {In: " ", E: Entry{}, Err: fmt.Errorf("entry is empty")},
		"one char line": {In: "# a\nb\nc",
			E:   Entry{Sections: []Section{{Title: "a", Body: "b\nc"}}},
			Err: nil},
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
		// Underline titles.
		"ul default": {In: "Do\n=\n\n\n\nLearn\n=\n\n\n", E: DefaultUnderline, Err: nil},
		"ul 1 rune line": {In: "a\n=\nb\nc",
			E:   Entry{Style: Underline, Sections: []Section{{Title: "a", Body: "b\nc"}}},
			Err: nil},
		"ul 3 lines": {In: "a\n=\na",
			E:   Entry{Style: Underline, Sections: []Section{{Title: "a", Body: "a"}}},
			Err: nil},
		"ul section with body": {
			In:  "Five\n=\n\nteve\n",
			E:   Entry{Style: Underline, Sections: []Section{{Title: "Five", Body: "teve"}}},
			Err: nil},
		"ul two sections with body": {
			In: "Five\n=\n\nteve\n\nFour\n=\n\ntoo\n",
			E: Entry{Style: Underline, Sections: []Section{{Title: "Five", Body: "teve"},
				{Title: "Four", Body: "too"}}},
			Err: nil},
		"ul repeated empty entries": {
			In:  "1\n=\n2\n=\n3\n=",
			E:   Entry{Style: Underline, Sections: []Section{{Title: "1"}, {Title: "2"}, {Title: "3"}}},
			Err: nil},
		"ul multiline body": {
			In:  "multi\n=\n\nmultiple\nlines",
			E:   Entry{Style: Underline, Sections: []Section{{Title: "multi", Body: "multiple\nlines"}}},
			Err: nil},
	}

	for id, test := range tests {
		ent, err := Import(test.In)
		if !reflect.DeepEqual(ent, test.E) {
			t.Errorf(testFail, ent, test.E, id)
		}
		if !errorEqual(err, test.Err) {
			t.Errorf(testFail, err, test.Err, id)
		}
	}
}

func TestEntry_Export(t *testing.T) {
	tests := map[string]struct {
		E   Entry
		Out string
	}{
		// Pound titles.
		"default entry": {E: Default, Out: "# Do\n\n\n\n# Learn\n\n\n"},
		"empty entry":   {E: Entry{Style: Pound, Sections: nil}, Out: ""},
		"section with body": {
			E:   Entry{Style: Pound, Sections: []Section{{Title: "Five", Body: "teve"}}},
			Out: "# Five\n\nteve\n"},
		"two sections with body": {
			E: Entry{Style: Pound,
				Sections: []Section{{Title: "Five", Body: "teve"}, {Title: "Four", Body: "too"}}},
			Out: "# Five\n\nteve\n\n# Four\n\ntoo\n"},
		// Underline titles.
		"ul default entry": {E: DefaultUnderline, Out: "Do\n==\n\n\n\nLearn\n=====\n\n\n"},
		"ul empty entry":   {E: Entry{Style: Underline, Sections: nil}, Out: ""},
		"ul section with body": {
			E:   Entry{Style: Underline, Sections: []Section{{Title: "Five", Body: "teve"}}},
			Out: "Five\n====\n\nteve\n"},
		"ul two sections with body": {
			E: Entry{Style: Underline,
				Sections: []Section{{Title: "Five", Body: "teve"}, {Title: "Four", Body: "too"}}},
			Out: "Five\n====\n\nteve\n\nFour\n====\n\ntoo\n"},
	}

	for id, test := range tests {
		exp := test.E.Export()
		if exp != test.Out {
			t.Errorf(testFail, exp, test.Out, id)
		}
	}
}

func TestEntry_Equals(t *testing.T) {
	entry1 := Entry{Name: "2019-01-01"}
	entry2 := Entry{Name: "2019-01-02"}
	entry2ul := Entry{Name: "2019-01-02", Style: Underline}

	tests := map[string]struct {
		E1     Entry
		E2     Entry
		Result bool
	}{
		"1": {E1: Default, E2: DefaultUnderline, Result: false},
		"2": {E1: entry1, E2: entry2, Result: true},
		"3": {E1: entry1, E2: entry2ul, Result: false},
		"4": {E1: entry2, E2: entry2ul, Result: false},
	}

	for id, test := range tests {
		result := test.E1.Equals(test.E2)
		if result != test.Result {
			t.Errorf(testFail, result, test.Result, id)
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
