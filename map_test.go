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
empty=

# invalid lines
color
my name = john
[trouble
text='...
`
	mapping := func(section, key, value, comment string, err error) {
		switch {
		case errors.Is(err, ingrid.ErrSyntax):
			fmt.Printf("input line:%v\n", err)

		case key != "":
			fmt.Printf("%s%s=%s\n", prefix(section), key, value)
		}
	}
	ingrid.Map(mapping, bufio.NewScanner(strings.NewReader(input)))
	// output:
	// debug=false
	// bind=localhost:80
	// example.text=escaped "
	// example.hostname=example.com
	// example.more=single "quoted" string
	// github.hostname=github.com
	// github.bind=localhost:443
	// github.empty=
	// input line:17 color SYNTAX ERROR: missing equal sign
	// input line:18 my name = john SYNTAX ERROR: space not allowed in key
	// input line:19 [trouble SYNTAX ERROR: missing right bracket
	// input line:20 text='... SYNTAX ERROR: missing end quote
}

// prefix returns section. if not empty.
func prefix(section string) string {
	if len(section) == 0 {
		return ""
	}
	return section + "."
}
