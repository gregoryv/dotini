package dotini_test

import (
	"fmt"
	"strings"

	"github.com/gregoryv/dotini"
)

func Example_parse() {
	input := `

# generic things
debug = false
defaultBind = localhost:80 # used set for all servers

[example]
hostname = "example.com"

[github]
hostname=github.com
bind=localhost:443
`
	r := strings.NewReader(input)

	handler := func(section, key, value, comment string) error {
		switch key {
		case "hostname":
			fmt.Printf("%s.%s = %s\n", section, key, value)
		}
		return nil
	}
	dotini.Parse(handler, r)
	// output:
	// example.hostname =  "example.com"
	// github.hostname = github.com
}
