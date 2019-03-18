# dev.journal

dev.journal is currently a small script that will setup a journal markdown file
for the day. The format is always `./<date>/<date>.md`.

It will copy the previous day's entry and uses the same format as above.

Things may change at any time.

### TODO

- Allow mixed parsing of `#` and `===` titles
- ? Allow mixed export of # and === titles
- Determine decent topic defaults (better that "General" and "Learn")

### Done

- Allow export of the `===` underlined title markdown
- Allow import of the `===` underlined title markdown
- Add parsing to allow for the manipulation of individual sections based on
  heading
- Add configuration
  - ? A .journal or journal.[json|toml|...] file
  - Ended up with `.devj` file, currently only json expected.
