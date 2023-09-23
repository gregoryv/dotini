/*
Package dotini provides parser for .ini files.

# Sections

Surrounded by square brackets

	[{SECTION}]

# Fields

Space separated

	{KEY1} = {VALUE}
	{KEY2}=   {VALUE}
	{KEY3}={VALUE}
	{KEY1}=     # empty values are accepted

Quoted values

	key1 = "hello"
	key2 = 'h'
	key3 = `h`

Comments

	# {COMMENT1}

	[{SECTION1}] #{COMMENT2}
	{KEY1}={VALUE} # {COMMENT3}

Empty lines are ignored.
*/
package dotini

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

// Parse file, fails on first syntax error and returns nil error once
// io.EOF is reached.
func Parse(handle Handler, r io.Reader) error {
	buf := bufio.NewReader(r)
	var lineno int
	var section string
	for {
		lineno++
		rawline, err := buf.ReadString('\n')

		line := strings.TrimSpace(rawline)

		// empty line and end of file
		if len(line) == 0 && errors.Is(err, io.EOF) {
			break
		}
		// skip empty lines
		if len(line) == 0 {
			continue
		}
		// only comment
		if line[0] == '#' {
			comment := findComment(line)
			_ = handle(section, "", "", comment)
			// ignore if handler fails when handling a comment
			continue
		}
		// section
		if line[0] == '[' {
			to := strings.Index(line, "]")
			if to == -1 {
				return fmt.Errorf("ini syntax error %v", lineno)
			}
			section = line[1:to]
			comment := findComment(line)
			err := handle(section, "", "", comment)
			if err != nil {
				return err
			}
			continue
		}
		// text without equal sign, invalid
		i := strings.Index(line, "=")
		if i == -1 {
			return fmt.Errorf("[line %v] %s %w", lineno, rawline, ErrSyntax)
		}
		// field line
		rawkey := line[:i]
		key := strings.TrimSpace(rawkey)
		rawvalue := line[i+1:]
		comment := findComment(rawvalue)
		{
			i := strings.Index(rawvalue, "#")
			if i != -1 {
				rawvalue = rawvalue[:i]
			}
		}
		{
			rawvalue = strings.TrimSpace(rawvalue)
			value := rawvalue
			if isQuoted(rawvalue) {
				unquoted, err := strconv.Unquote(rawvalue)
				if err != nil {
					log.Printf("%q %v", unquoted, err)
					if unquoted == "" && errors.Is(err, strconv.ErrSyntax) {
						return fmt.Errorf("[line %v] %s %w", lineno, rawline, ErrSyntax)
					}
				} else {
					value = unquoted
				}
			}
			err = handle(section, key, value, comment)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isQuoted(v string) bool {
	if len(v) == 0 {
		return false
	}
	switch v[0] {
	case '"', '\'', '`':
		return true
	}
	return false
}

func findComment(line string) string {
	i := strings.Index(line, "#")
	if i == -1 {
		return ""
	}
	rawComment := line[i+1:]
	comment := strings.TrimSpace(rawComment)
	return comment
}

type Handler func(section, key, value, comment string) error

var ErrSyntax = fmt.Errorf("syntax error")
