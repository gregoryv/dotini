Package ingrid parses .ini files.

## Features

The format is often used in .ini, .conf or .cfg files. This
implementation is focused on performance and comes with minor
limitations.

  - lines starting with # are comments
  - values may be quoted using " or ` (backtick)
  - sections [ and end with ]

## Limitations

Currently the limitations for this implementation are

  - keys cannot contain spaces
  - no comments are allowed after a key value pair.
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
    
    [github]
    hostname=github.com
    bind=localhost:443
    `
    	mapping := func(section, key, value, comment string) error {
    		switch key {
    		case "hostname", "text":
    			fmt.Printf("%s.%s = %s\n", section, key, value)
    		}
    		return nil
    	}
    	ingrid.Map(mapping, bufio.NewScanner(strings.NewReader(input)))
    	// output:
    	// example.text = escaped "
    	// example.hostname = example.com
    	// github.hostname = github.com
    }
## Benchmark

     goos: linux
     goarch: amd64
     pkg: github.com/gregoryv/ingrid
     cpu: Intel(R) Xeon(R) E-2288G CPU @ 3.70GHz
     Benchmark_Map-16    	144609124	         7.636 ns/op	       0 B/op	       0 allocs/op
