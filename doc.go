/*
Package ingrid parses .ini files.

# Features

The format is often used in .ini, .conf or .cfg files. This
implementation is focused on performance and comes with minor
limitations.

  - lines starting with # or ; are comments
  - values may be quoted using " or ` (backtick)
  - sections [ and end with ]

# Limitations

Currently the limitations for this implementation are

  - keys cannot contain spaces
  - no comments are allowed after a key value pair.
  - no multiline values
*/
package ingrid

//go:generate go run ./internal -in doc.go -out README.md
