// Command ingrid parses .ini files
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gregoryv/ingrid"
)

func main() {
	next := parseFiles(os.Args[1:])
	var err error
	for next != nil {
		next, err = next()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func parseFiles(files []string) parseFn {
	return func() (parseFn, error) {
		if len(files) == 0 {
			return nil, nil
		}

		fh, err := os.Open(files[0])
		if err != nil {
			return nil, fmt.Errorf("%s:%w", files[0], err)
		}
		defer fh.Close()

		err = ingrid.Map(
			printKeyValue, bufio.NewScanner(fh),
		)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", files[0], err)
		}

		return parseFiles(files[1:]), nil
	}
}

func printKeyValue(section, key, value, comment string) error {
	if key == "" {
		return nil
	}
	prefix := ""
	if section != "" {
		prefix = section + "."
	}
	fmt.Printf("%s%s = %s\n", prefix, key, value)
	return nil
}

type parseFn func() (parseFn, error)
