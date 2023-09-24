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
