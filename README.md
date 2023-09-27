<!-- GENERATED, DO NOT EDIT! See internal/updateReadme.go -->
Package ingrid parses .ini files.

## Quickstart

	$ go get -u github.com/gregoryv/ingrid

## Features

The format is often used in .ini, .conf or .cfg files. This
implementation is focused on performance and comes with minor
limitations.

  - comments start with # or ;
  - values may be quoted using ", ` (backtick), or ' (single tick)
  - sections start with [ and end with ]
  - spaces before and after key, values are removed

## Limitations

Currently the limitations for this implementation are

  - no spaces in keys
  - no comments after a key value pair
  - no multiline values

## Example

    package ingrid_test
    
    import (
    	"bufio"
    	"errors"
    	"fmt"
    	"strings"
    
    	"github.com/gregoryv/ingrid"
    )
    
    func Example() {
    	input := `# generic things
    debug = false
    # default for servers
    bind= localhost:80
    
    [example]
    text = "escaped \""
    hostname = "example.com"
    more = 'single "quoted" string'
    
    [github]
    hostname=github.com
    bind=localhost:443
    
    # invalid lines
    color
    my name = john
    [trouble
    text='...
    `
    	mapping := func(section, key, value, comment string, err error) {
    		if errors.Is(err, ingrid.ErrSyntax) {
    			fmt.Printf("input line:%v\n", err)
    			return
    		}
    		if key != "" {
    			var prefix string
    			if len(section) > 0 {
    				prefix = section + "."
    			}
    			fmt.Printf("%s%s = %s\n", prefix, key, value)
    		}
    	}
    	ingrid.Map(mapping, bufio.NewScanner(strings.NewReader(input)))
    	// output:
    	// debug = false
    	// bind = localhost:80
    	// example.text = escaped "
    	// example.hostname = example.com
    	// example.more = single "quoted" string
    	// github.hostname = github.com
    	// github.bind = localhost:443
    	// input line:16 color SYNTAX ERROR: missing equal sign
    	// input line:17 my name = john SYNTAX ERROR: space not allowed in key
    	// input line:18 [trouble SYNTAX ERROR: missing right bracket
    	// input line:19 text='... SYNTAX ERROR: missing end quote
    }

## Benchmark

     goos: linux
     goarch: amd64
     pkg: github.com/gregoryv/ingrid
     cpu: Intel(R) Xeon(R) E-2288G CPU @ 3.70GHz
     Benchmark_Map-16  155136238     7.670 ns/op    0 B/op    0 allocs/op
