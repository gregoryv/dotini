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

// grabComment returns entire buf semihash is 0, nil otherwise,
// ie. comments are only allowed on separate lines.
func grabComment(buf []byte, semihash int) []byte {
	if semihash == 0 {
		return buf[semihash:]
	}
	return nil
}

func grabKeyValue(buf []byte, equal int) (key, value []byte, err error) {
	if equal >= 0 {
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
	}
	return
}

var singleQuote byte = '\''
var ErrSyntax = fmt.Errorf("syntax error")

func normalizeQuotes(value []byte) {
	last := len(value) - 1
	if value[0] == singleQuote && value[last] == singleQuote {
		value[0] = '`'
		value[last] = '`'
	}
}

func grabSection(buf, current []byte, lbrack, rbrack int) []byte {
	if isSection(lbrack, rbrack) {
		section := buf[lbrack+1 : rbrack]
		section = bytes.TrimSpace(section)
		return section
	}
	return current
}

func indexElements(buf []byte) (lbrack, rbrack, equal, semihash int) {
	lbrack, rbrack, equal, semihash = -1, -1, -1, -1
	for i, b := range buf {
		setIndex(i, &lbrack, b, '[')
		setIndex(i, &rbrack, b, ']')
		setIndex(i, &equal, b, '=')

		if b == '#' || b == ';' {
			semihash = i
			break
		}
	}
	return
}

func setIndex(i int, to *int, b, is byte) {
	if b == is {
		*to = i
	}
}

func isEmpty(section, key, value, comment []byte) bool {
	return len(section)+len(key)+len(value)+len(comment) == 0
}

func isSection(lbrack, rbrack int) bool {
	return rbrack > lbrack && lbrack >= 0 && rbrack >= 0
}

func isQuoted(value []byte) bool {
	return len(value) > 0 && bytes.ContainsAny(value[:1], "\"'`")
}

// Mapfn is called for each non empty line. section is always the
// current section. At least one of the arguments is not empty.
type Mapfn func(section, key, value, comment string) error
