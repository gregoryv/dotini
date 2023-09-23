package dotini

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Benchmark_Parse(b *testing.B) {
	example := `
# comment 1
k1=1
k2=string
k3=value
k4 = something # commented
k5=123.5
`
	handler := func(section, key, value, comment string) error { return nil }

	for i := 0; i < b.N; i++ {
		r := strings.NewReader(example)
		err := Parse(handler, r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Test_Parse(t *testing.T) {
	cases := []*IniCase{
		{
			Test:    "completely empty",
			Example: "",
		},
		{
			Test:    "almost empty",
			Example: "\n\n",
		},
		{
			Test: "only comments",
			Example: `
#
# comment`,
			ExpectLines:    2,
			ExpectComments: []string{"comment"},
		},
		{
			Test:         "empty value",
			Example:      "k1=",
			ExpectLines:  1,
			ExpectKeys:   []string{"k1"},
			ExpectValues: []string{""},
		},
		{
			Example:        "# a comment",
			ExpectLines:    1,
			ExpectComments: []string{"a comment"},
		},
		{
			Example:        "[hut] # green",
			ExpectLines:    1,
			ExpectComments: []string{"green"},
		},
		{
			Example: `# comment
field1=string value
field2 = 1`,
			ExpectLines:    3,
			ExpectKeys:     []string{"field1", "field2"},
			ExpectValues:   []string{"string value", "1"},
			ExpectComments: []string{"comment"},
		},
		{
			Example: `
k1=1
k2="hello \""
`,
			ExpectLines:  2,
			ExpectKeys:   []string{"k1", "k2"},
			ExpectValues: []string{"1", `hello "`},
		},
		{
			Example: `
[server 1]
hostname=example.com

[server 2]
hostname=github.com`,
			ExpectLines:  4,
			ExpectKeys:   []string{"hostname", "hostname"},
			ExpectValues: []string{"example.com", "github.com"},
		},

		// ----------------------------------------
		// error cases below
		// ----------------------------------------

		{
			Test:        "missing equal sign",
			Example:     "key1",
			ExpectError: true,
			ExpectLines: 0,
		},
		{
			Example:     "[incomplete-section",
			ExpectError: true,
			ExpectLines: 0,
		},
		{
			Test:        "open quote",
			Example:     `k1="hello`,
			ExpectError: true,
		},
		{
			Example:     "[smurf]",
			ExpectError: true,
			ExpectLines: 1,
			HandleErr:   fmt.Errorf("unknown section"),
		},
		{
			Example:        "fx=233 # field comment",
			ExpectError:    true,
			ExpectLines:    1,
			ExpectKeys:     []string{"fx"},
			ExpectValues:   []string{"233"},
			ExpectComments: []string{"field comment"},
		},
		{
			Example:     "nosuch=abc",
			ExpectError: true,
			HandleErr:   fmt.Errorf("handler failed"),
			ExpectLines: 1,
		},
	}
	for _, c := range cases {
		t.Run(c.Test, func(t *testing.T) {
			r := strings.NewReader(c.Example)
			err := Parse(c.UseIni, r)

			if err != nil && !c.ExpectError {
				t.Log(c.Example)
				t.Error(err)
			}
			c.Verify(t)
		})
	}
}

type IniCase struct {
	Test string

	Example     string
	HandleErr   error
	ExpectError bool

	ExpectLines int
	lines       int

	ExpectKeys []string
	keys       []string

	ExpectValues []string
	values       []string

	ExpectComments []string
	comments       []string
}

func (c *IniCase) UseIni(section, key, value, comment string) error {
	c.lines++
	if c.HandleErr != nil {
		return c.HandleErr
	}
	if key != "" {
		c.keys = append(c.keys, key)
		c.values = append(c.values, value)
	}
	if comment != "" {
		c.comments = append(c.comments, comment)
	}
	return nil
}

func (c *IniCase) Verify(t *testing.T) {
	if c.lines != c.ExpectLines {
		t.Log("example:", c.Example)
		t.Error("lines", c.lines)
	}
	if !reflect.DeepEqual(c.keys, c.ExpectKeys) {
		t.Log("example:", c.Example)
		t.Errorf("keys: %q", c.keys)
	}
	if !reflect.DeepEqual(c.values, c.ExpectValues) {
		t.Log("example:", c.Example)
		t.Errorf("values: %q", c.values)
	}
	if !reflect.DeepEqual(c.comments, c.ExpectComments) {
		t.Log("example:", c.Example)
		t.Errorf("comments: %q", c.comments)
	}
}
