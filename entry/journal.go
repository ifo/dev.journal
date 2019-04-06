package entry

// Journal is a map of Entries where they key is the string of the date the Entry was written.
type Journal struct {
	Entries map[EntryName]Entry `json:"entries"`
}

// NewJournal creates a new empty journal with a non-nil Entries map.
func NewJournal() *Journal {
	return &Journal{Entries: map[EntryName]Entry{}}
}

func (j *Journal) Add(e Entry) {
	if existing, exists := j.Entries[e.Name]; exists {
		if !e.Equals(existing) {
			// TODO: a better updated version naming scheme.
			e.Name = e.Name + "1"
			j.Add(e)
		}
		return
	}

	j.Entries[e.Name] = e
}

// Contains is an expensive way of determining if a journal already has a specific entry.
func (j *Journal) Contains(e Entry) bool {
	for _, entry := range j.Entries {
		if entry.Equals(e) {
			return true
		}
	}
	return false
}
