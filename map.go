package ingrid

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
)

// Map maps each scanned line to mapping until EOF is reached.
// Returns ErrSyntax if line is badly formatted.
func Map(mapping Mapfn, scanner *bufio.Scanner) {
	var lineno int
	// current section
	var current []byte

	for scanner.Scan() {
		lineno++
		buf := scanner.Bytes()
		buf = bytes.TrimSpace(buf)

		if len(buf) == 0 {
			continue
		}

		// grab section, key, value and comment
		section, key, value, comment, err := parse(buf, current)
		current = section

		if err != nil {
			err = fmt.Errorf(
				"%v %s %w: %v", lineno, string(buf), ErrSyntax, err,
			)
		}
		mapping(
			string(section),
			string(key),
			string(value),
			string(comment),
			err,
		)
	}
}

// parse finds one or more of the allowed parts.
func parse(buf, current []byte) (
	section, key, value, comment []byte, err error,
) {
	lbrack, rbrack, equal, semihash := indexElements(buf)
	switch {
	case lbrack == 0:
		section = grabSection(&err, buf, current, lbrack, rbrack)

	case semihash == 0:
		comment = buf[semihash:]

	default:
		key, value = grabKeyValue(&err, buf, equal)
	}
	if len(section) == 0 {
		section = current
	}
	return
}

// indexElements indexes first occurence of [, ], = and # or ; in buf
func indexElements(buf []byte) (lbrack, rbrack, equal, semihash int) {
	lbrack, rbrack, equal, semihash = -1, -1, -1, -1
	for i, b := range buf {
		isCommentChar := b == '#' || b == ';'
		if isCommentChar {
			semihash = i
			break
		}
		setIndex(i, &lbrack, b, '[')
		setIndex(i, &rbrack, b, ']')
		setIndex(i, &equal, b, '=')
	}
	return
}

// setIndex updates dst with i if a == b and dst == -1
func setIndex(i int, dst *int, a, b byte) {
	if *dst != -1 {
		return
	}
	if a != b {
		return
	}
	*dst = i
}

// grabSection returns new section if buf contains one, otherwise
// current is returned.
func grabSection(err *error, buf, current []byte, lbrack, rbrack int) []byte {
	if lbrack == 0 && rbrack == -1 {
		*err = fmt.Errorf("missing right bracket")
	}
	if isSection(lbrack, rbrack) {
		section := buf[lbrack+1 : rbrack]
		section = bytes.TrimSpace(section)
		return section
	}
	return current
}

// grabKeyValue returns key and value from buf. Quoted values are
// unquoted. Returns ErrSyntax if incorrectly formated.
func grabKeyValue(err *error, buf []byte, equal int) (key, value []byte) {
	if equal == -1 {
		*err = fmt.Errorf("missing equal sign")
		return
	}
	key = bytes.TrimSpace(buf[:equal])
	if bytes.ContainsAny(key, " ") {
		*err = fmt.Errorf("space not allowed in key")
	}
	value = grabValue(err, buf, equal)
	return
}

func grabValue(err *error, buf []byte, equal int) (value []byte) {
	value = buf[equal+1:]
	value = bytes.TrimSpace(value)
	if isQuoted(value) {
		normalizeQuotes(value)
		valstr, e := strconv.Unquote(string(value))
		if e != nil {
			*err = fmt.Errorf("missing end quote")
		}
		value = []byte(valstr)
	}
	return
}

var singleQuote byte = '\''
var ErrSyntax = fmt.Errorf("SYNTAX ERROR")

// normalizeQuotes replaces single tick quotes with `
func normalizeQuotes(value []byte) {
	last := len(value) - 1
	if value[0] == singleQuote && value[last] == singleQuote {
		value[0] = '`'
		value[last] = '`'
	}
}

// isSection returns true if lbrack somes before rbrack
func isSection(lbrack, rbrack int) bool {
	return rbrack > lbrack && lbrack >= 0 && rbrack >= 0
}

// isQuoted returns true if the first character of value looks like
// quote char, value cannot be empty
func isQuoted(value []byte) bool {
	if len(value) == 0 {
		return false
	}
	const quoteChars = "\"'`"
	return bytes.ContainsAny(value[:1], quoteChars)
}

// Mapfn is called for each non empty line. section is always the
// current section. At least one of the arguments is not empty.
type Mapfn func(section, key, value, comment string, err error)
