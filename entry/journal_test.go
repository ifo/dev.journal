package entry

import (
	"reflect"
	"testing"
)

func TestJournal_Add(t *testing.T) {
	entry1 := Entry{Name: "2019-01-01"}
	entry2 := Entry{Name: "2019-01-02"}
	entry2ul := Entry{Name: "2019-01-02", Style: Underline}
	entry2ulAdded := Entry{Name: "2019-01-021", Style: Underline}
	entry2diff := Entry{Name: "2019-01-02", Sections: []Section{{Title: "a"}}}
	entry2diffAdded := Entry{Name: "2019-01-0211", Sections: []Section{{Title: "a"}}}
	outJournal1 := NewJournal()
	outJournal1.Entries[entry1.Name] = entry1
	outJournal2 := NewJournal()
	outJournal2.Entries[entry1.Name] = entry1
	outJournal2.Entries[entry2.Name] = entry2
	outJournal3 := NewJournal()
	outJournal3.Entries[entry2.Name] = entry2
	outJournal3.Entries[entry2ulAdded.Name] = entry2ulAdded
	outJournal3.Entries[entry2diffAdded.Name] = entry2diffAdded

	tests := map[string]struct {
		ToAdd      []Entry
		OutJournal *Journal
	}{
		"1": {ToAdd: []Entry{entry1}, OutJournal: outJournal1},
		"2": {ToAdd: []Entry{entry1, entry1}, OutJournal: outJournal1},
		"3": {ToAdd: []Entry{entry1, entry2}, OutJournal: outJournal2},
		"4": {ToAdd: []Entry{entry1, entry2, entry1}, OutJournal: outJournal2},
		"5": {ToAdd: []Entry{entry2, entry2ul, entry2diff}, OutJournal: outJournal3},
	}

	for id, test := range tests {
		journal := NewJournal()
		for _, entry := range test.ToAdd {
			journal.Add(entry)
		}

		if !reflect.DeepEqual(journal, test.OutJournal) {
			t.Errorf(testFail, journal, test.OutJournal, id)
		}
	}
}

func TestJournal_Contains(t *testing.T) {
	sections1 := []Section{{Title: "title1"}}
	sections2 := []Section{{Title: "title2"}}
	entry1 := Entry{Name: "2019-01-01", Sections: sections1}
	entry2 := Entry{Name: "2019-01-02", Sections: sections2}
	entry2ul := Entry{Name: "2019-01-02", Sections: sections2, Style: Underline}
	entry2ulAdded := Entry{Name: "2019-01-021", Sections: sections2, Style: Underline}
	journal1 := NewJournal()
	journal1.Entries[entry1.Name] = entry1
	journal2 := NewJournal()
	journal2.Entries[entry1.Name] = entry1
	journal2.Entries[entry2.Name] = entry2
	journal3 := NewJournal()
	journal3.Entries[entry2.Name] = entry2
	journal3.Entries[entry2ulAdded.Name] = entry2ulAdded

	tests := map[string]struct {
		Journal *Journal
		Entry   Entry
		Out     bool
	}{
		"1": {Journal: journal1, Entry: entry1, Out: true},
		"2": {Journal: journal1, Entry: entry2, Out: false},
		"3": {Journal: journal2, Entry: entry1, Out: true},
		"4": {Journal: journal2, Entry: entry2, Out: true},
		"5": {Journal: journal2, Entry: entry2ul, Out: false},
		"6": {Journal: journal3, Entry: entry1, Out: false},
		"7": {Journal: journal3, Entry: entry2, Out: true},
		"8": {Journal: journal3, Entry: entry2ul, Out: true},
	}

	for id, test := range tests {
		out := test.Journal.Contains(test.Entry)
		if out != test.Out {
			t.Errorf(testFail, out, test.Out, id)
		}
	}
}
