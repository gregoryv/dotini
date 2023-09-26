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
    	"fmt"
    	"strings"
    
    	"github.com/gregoryv/ingrid"
    )
    
    func Example() {
    	input := `# generic things
    debug = false
    # used set for all servers
    defaultBind = localhost:80
    
    [example]
    text = "escaped \""
    hostname = "example.com"
    more = 'single "quoted" string'
    
    [github]
    hostname=github.com
    bind=localhost:443
    
    color
    my name = john
    [trouble
    text='...
    `
    	mapping := func(section, key, value, comment string, err error) {
    		if err != nil {
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
    	// defaultBind = localhost:80
    	// example.text = escaped "
    	// example.hostname = example.com
    	// example.more = single "quoted" string
    	// github.hostname = github.com
    	// github.bind = localhost:443
    	// input line:15 color missing equal sign: syntax error
    	// input line:16 my name = john space not allowed in key: syntax error
    	// input line:17 [trouble missing right bracket: syntax error
    	// input line:18 text='... missing end quote: syntax error
    }

## Benchmark

     goos: linux
     goarch: amd64
     pkg: github.com/gregoryv/ingrid
     cpu: Intel(R) Xeon(R) E-2288G CPU @ 3.70GHz
     Benchmark_Map-16  153884863     8.000 ns/op    0 B/op    0 allocs/op
