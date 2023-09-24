/*
Package ingrid parses .ini files.

# Features

The format is often used in .ini, .conf or .cfg files. This
implementation is focused on performance and comes with minor
limitations.

  - comments start with # or ;
  - values may be quoted using ", ` (backtick), or ' (single tick)
  - sections start with [ and end with ]
  - spaces before and after key, values are removed

# Limitations

Currently the limitations for this implementation are

  - no spaces in keys
  - no comments after a key value pair
  - no multiline values
*/
package ingrid

//go:generate go run ./internal -in doc.go -out README.md
