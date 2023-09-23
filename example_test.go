package dotini_test

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/gregoryv/dotini"
)

func Example_parse() {
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
	handler := func(section, key, value, comment string) error {
		switch key {
		case "hostname", "text":
			fmt.Printf("%s.%s = %s\n", section, key, value)
		}
		return nil
	}
	dotini.Parse(handler, bufio.NewScanner(strings.NewReader(input)))
	// output:
	// example.text = escaped "
	// example.hostname = example.com
	// github.hostname = github.com
}
