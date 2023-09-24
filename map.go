package ingrid

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
)

// Map maps each scanned line to mapping until EOF is reached.
// Returns ErrSyntax if line is badly formatted.
func Map(mapping Mapfn, scanner *bufio.Scanner) error {
	var lineno int
	// last section
	var section []byte

	for scanner.Scan() {
		lineno++
		buf := scanner.Bytes()
		buf = bytes.TrimSpace(buf)
		lbrack, rbrack, equal, semihash := indexElements(buf)

		// grab section, key, value and comment
		section = grabSection(buf, section, lbrack, rbrack)
		key, value, err := grabKeyValue(buf, equal)
		if err != nil {
			return fmt.Errorf("%w: line %v: %s", err, lineno, string(buf))
		}
		comment := grabComment(buf, semihash)

		if !isEmpty(section, key, value, comment) {
			mapping(
				string(section),
				string(key),
				string(value),
				string(comment),
			)
			continue
		}
		if len(buf) > 0 {
			return fmt.Errorf("%w: line %v: %s", ErrSyntax, lineno, string(buf))
		}
	}
	return nil
}

// grabComment returns entire buf if semihash is 0, nil otherwise,
// ie. comments are only allowed on separate lines.
func grabComment(buf []byte, semihash int) []byte {
	if semihash == 0 {
		return buf[semihash:]
	}
	return nil
}

// grabKeyValue returns key and value from buf. Quoted values are
// unquoted. Returns ErrSyntax if incorrectly formated.
func grabKeyValue(buf []byte, equal int) (key, value []byte, err error) {
	if equal == -1 {
		return
	}
	key = bytes.TrimSpace(buf[:equal])
	if bytes.ContainsAny(key, " ") {
		return nil, nil, ErrSyntax
	}

	value = buf[equal+1:]
	value = bytes.TrimSpace(value)
	if isQuoted(value) {
		normalizeQuotes(value)
		valstr, err := strconv.Unquote(string(value))
		if err != nil {
			return nil, nil, ErrSyntax
		}
		value = []byte(valstr)
	}
	return
}

var singleQuote byte = '\''
var ErrSyntax = fmt.Errorf("syntax error")

// normalizeQuotes replaces single tick quotes with `
func normalizeQuotes(value []byte) {
	last := len(value) - 1
	if value[0] == singleQuote && value[last] == singleQuote {
		value[0] = '`'
		value[last] = '`'
	}
}

// grabSection returns new section if buf contains one, otherwise
// current is returned.
func grabSection(buf, current []byte, lbrack, rbrack int) []byte {
	if isSection(lbrack, rbrack) {
		section := buf[lbrack+1 : rbrack]
		section = bytes.TrimSpace(section)
		return section
	}
	return current
}

// indexElements indexes first occurence of [, ], = and # or ; in buf
func indexElements(buf []byte) (lbrack, rbrack, equal, semihash int) {
	lbrack, rbrack, equal, semihash = -1, -1, -1, -1
	for i, b := range buf {
		if b == '#' || b == ';' {
			semihash = i
			break
		}
		setIndex(i, &lbrack, b, '[')
		setIndex(i, &rbrack, b, ']')
		setIndex(i, &equal, b, '=')
	}
	return
}

// setIndex updates dst with i if a == b
func setIndex(i int, dst *int, a, b byte) {
	if a != b {
		return
	}
	*dst = i
}

// isEmpty returns true if all arguments are empty
func isEmpty(section, key, value, comment []byte) bool {
	return len(section)+len(key)+len(value)+len(comment) == 0
}

// isSection returns true if lbrack somes before rbrack
func isSection(lbrack, rbrack int) bool {
	return rbrack > lbrack && lbrack >= 0 && rbrack >= 0
}

// isQuoted returns true if the first character of value looks like
// quote char, value cannot be empty
func isQuoted(value []byte) bool {
	const quoteChars = "\"'`"
	return bytes.ContainsAny(value[:1], quoteChars)
}

// Mapfn is called for each non empty line. section is always the
// current section. At least one of the arguments is not empty.
type Mapfn func(section, key, value, comment string) error
