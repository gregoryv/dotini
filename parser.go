package dotini

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Parse parses a basic ini file, no sections.
func Parse(handle HandlerFunc, r io.Reader) error {
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
			_ = handle.UseIni(section, "", "", comment)
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
			err := handle.UseIni(section, "", "", comment)
			if err != nil {
				return err
			}
			continue
		}
		// text without equal sign, invalid
		i := strings.Index(line, "=")
		if i == -1 {
			return fmt.Errorf("ini syntax error %v", lineno)
		}
		// field line
		rawkey := line[:i]
		key := strings.TrimSpace(rawkey)
		rawvalue := line[i+1:]
		value := rawvalue
		comment := findComment(rawvalue)
		{
			i := strings.Index(rawvalue, "#")
			if i != -1 {
				value = strings.TrimSpace(rawvalue[:i])
			}
		}
		{
			err := handle.UseIni(section, key, value, comment)
			if err != nil {
				return err
			}
		}
		if err != nil && errors.Is(err, io.EOF) {
			break
		}
	}
	return nil
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

type HandlerFunc func(section, key, value, comment string) error

func (h HandlerFunc) UseIni(section, key, value, comment string) error {
	return h(section, key, value, comment)
}
