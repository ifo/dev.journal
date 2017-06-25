# dev.journal

dev.journal is currently a small script that will setup a journal markdown file
for the day. The format is always `./<date>/<date>.md`.

It will copy any previous days entry that's been written in the past 7 days and
uses the same format as above.

Configuration will be coming soon! And things may change at any time.

### TODO

- Add configuration
  - ? A .journal or journal.[json|toml|...] file
- Allow mixed parsing of `#` and `===` titles
- ? Allow mixed export of # and === titles
- Determine decent topic defaults (better that "General" and "Learn")

### Done

- Allow export of the `===` underlined title markdown
- Allow import of the `===` underlined title markdown
- Add parsing to allow for the manipulation of individual sections based on
  heading
