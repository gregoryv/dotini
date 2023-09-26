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
	mapping := func(section, key, value, comment string, err error) {
		if key != "" {
			var prefix string
			if len(section) > 0 {
				prefix = section + "."
			}
			fmt.Printf("%s%s = %s\n", prefix, key, value)
		}
		if err != nil {
			fmt.Println(err)
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
}
