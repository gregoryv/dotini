package ingrid

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
)

func Map(mapping Mapfn, scanner *bufio.Scanner) error {
	var lineno int
	var section []byte
	for scanner.Scan() {
		lineno++
		buf := scanner.Bytes()
		buf = bytes.TrimSpace(buf)

		var key, value, comment []byte
		var lbrack, rbrack, equal, hash = indexElements(buf)

		section = grabSection(section, buf, lbrack, rbrack)

		key, value, err := grabKeyValue(buf, equal)
		if err != nil {
			return fmt.Errorf("%w: line %v: %s", err, lineno, string(buf))
		}
		comment = grabComment(buf, hash)

		switch {
		case !isEmpty(section, key, value, comment):
			mapping(
				string(section),
				string(key),
				string(value),
				string(comment),
			)

		case len(buf) > 0:
			return fmt.Errorf("syntax error: line %v: %s", lineno, string(buf))
		}

	}
	return nil
}

func grabComment(buf []byte, hash int) []byte {
	if hash == 0 {
		return buf[hash:]
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
			valstr, err := strconv.Unquote(string(value))
			if err != nil {
				return nil, nil, ErrSyntax
			}
			value = []byte(valstr)
		}
	}
	return
}

var ErrSyntax = fmt.Errorf("syntax error")

func grabSection(current, buf []byte, lbrack, rbrack int) []byte {
	if isSection(lbrack, rbrack) {
		section := buf[lbrack+1 : rbrack]
		section = bytes.TrimSpace(section)
		return section
	}
	return current
}

func indexElements(buf []byte) (lbrack, rbrack, equal, hash int) {
	lbrack, rbrack, equal, hash = -1, -1, -1, -1
	for i, b := range buf {
		setIndex(i, &lbrack, b, '[')
		setIndex(i, &rbrack, b, ']')
		setIndex(i, &equal, b, '=')

		if b == '#' {
			hash = i
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

type Mapfn func(section, key, value, comment string) error
