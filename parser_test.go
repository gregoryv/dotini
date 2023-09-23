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
			Test:  "completely empty",
			Input: "",
		},
		{
			Test:  "almost empty",
			Input: "\n\n",
		},
		{
			Test: "only comments",
			Input: `
#
# comment`,
			ExpectLines:    2,
			ExpectComments: []string{"comment"},
		},
		{
			Test:         "empty value",
			Input:        "k1=",
			ExpectLines:  1,
			ExpectKeys:   []string{"k1"},
			ExpectValues: []string{""},
		},
		{
			Input:          "# a comment",
			ExpectLines:    1,
			ExpectComments: []string{"a comment"},
		},
		{
			Input:          "[hut] # green",
			ExpectLines:    1,
			ExpectComments: []string{"green"},
		},
		{
			Input: `# comment
field1=string value
field2 = 1`,
			ExpectLines:    3,
			ExpectKeys:     []string{"field1", "field2"},
			ExpectValues:   []string{"string value", "1"},
			ExpectComments: []string{"comment"},
		},
		{
			Input: `
k1=1
k2="hello \""
`,
			ExpectLines:  2,
			ExpectKeys:   []string{"k1", "k2"},
			ExpectValues: []string{"1", `hello "`},
		},
		{
			Input: `
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
			Input:       "key1",
			ExpectError: "[line 1] key1 syntax error",
			ExpectLines: 0,
		},
		{
			Test:        "incomplete section",
			Input:       "[incomplete-section",
			ExpectError: "[line 1] [incomplete-section syntax error",
			ExpectLines: 0,
		},
		{
			Test:        "open quote",
			Input:       `k1="hello`,
			ExpectError: `[line 1] k1="hello syntax error`,
		},
		{
			Input:       "[smurf]",
			ExpectError: "unknown section",
			ExpectLines: 1,
			HandleErr:   fmt.Errorf("unknown section"),
		},
		{
			Input:          "fx=233 # field comment",
			ExpectError:    ErrSyntax.Error(),
			ExpectLines:    1,
			ExpectKeys:     []string{"fx"},
			ExpectValues:   []string{"233"},
			ExpectComments: []string{"field comment"},
		},
		{
			Input:       "nosuch=abc",
			ExpectError: "unknown key",
			HandleErr:   fmt.Errorf("unknown key"),
			ExpectLines: 1,
		},
	}
	for _, c := range cases {
		t.Run(c.Test, func(t *testing.T) {
			r := strings.NewReader(c.Input)
			err := Parse(c.UseIni, r)

			if err != nil {
				if c.ExpectError == "" {
					t.Log(c.Input)
					t.Error(err)
				}
				if got := err.Error(); got != c.ExpectError {
					t.Log("got", got)
					t.Log("exp", c.ExpectError)
					t.Fail()
				}
			}
			c.Verify(t)
		})
	}
}

type IniCase struct {
	Test string

	Input       string
	HandleErr   error
	ExpectError string

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
		t.Log("input:", c.Input)
		t.Error("lines", c.lines)
	}
	if !reflect.DeepEqual(c.keys, c.ExpectKeys) {
		t.Log("input:", c.Input)
		t.Errorf("keys: %q", c.keys)
	}
	if !reflect.DeepEqual(c.values, c.ExpectValues) {
		t.Log("input:", c.Input)
		t.Errorf("values: %q", c.values)
	}
	if !reflect.DeepEqual(c.comments, c.ExpectComments) {
		t.Log("input:", c.Input)
		t.Errorf("comments: %q", c.comments)
	}
}
