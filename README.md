Package ingrid parses .ini files.

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
    `
    	mapping := func(section, key, value, comment string) error {
    		switch key {
    		case "hostname", "text", "more":
    			fmt.Printf("%s.%s = %s\n", section, key, value)
    		}
    		return nil
    	}
    	ingrid.Map(mapping, bufio.NewScanner(strings.NewReader(input)))
    	// output:
    	// example.text = escaped "
    	// example.hostname = example.com
    	// example.more = single "quoted" string
    	// github.hostname = github.com
    }
## Benchmark

     goos: linux
     goarch: amd64
     pkg: github.com/gregoryv/ingrid
     cpu: Intel(R) Xeon(R) E-2288G CPU @ 3.70GHz
     Benchmark_Map-16    	152860821	         7.873 ns/op	       0 B/op	       0 allocs/op
