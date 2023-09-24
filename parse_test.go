package ingrid

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func Test_Map_incorrect(t *testing.T) {
	assert := func(label, input string) {
		t.Helper()
		t.Run(label, func(t *testing.T) {
			var full bytes.Buffer
			handler := newHandler(t, &full)
			err := Map(handler, bufio.NewScanner(strings.NewReader(input)))
			if err == nil {
				t.Log(full.String())
				t.Error("expected error")
			}
		})
	}
	assert("missing right bracket", `[test`)
	assert("missing equal sign", `key`)
	assert("space in key", `key 1 = value`)
	assert("missing quote", `key1 = "value`)
}

func Test_Map(t *testing.T) {
	example, _ := os.ReadFile("testdata/example.ini")
	var full bytes.Buffer
	handler := newHandler(t, &full)
	err := Map(handler, bufio.NewScanner(bytes.NewReader(example)))
	if err != nil {
		t.Log(full.String())
		t.Log(err)
	}
	golden.Assert(t, full.String())
}

func Benchmark_Map(b *testing.B) {
	example, _ := os.ReadFile("testdata/example.ini")
	handler := func(section, key, value, comment string) error { return nil }
	r := bytes.NewReader(example)
	scanner := bufio.NewScanner(r)
	for i := 0; i < b.N; i++ {
		err := Map(handler, scanner)
		if err != nil {
			b.Fatal(err)
		}
		r.Reset(example)
	}
}

func newHandler(t *testing.T, full *bytes.Buffer) Handler {
	return func(section, key, value, comment string) error {
		var buf bytes.Buffer
		if len(section) > 0 {
			fmt.Fprintf(&buf, "[%s]", section)
		}
		if len(key) > 0 {
			fmt.Fprint(&buf, key, "=")

			if len(value) > 0 {
				fmt.Fprintf(&buf, value)
			}
		}
		if len(comment) > 0 {
			fmt.Fprintf(&buf, comment)
		}
		full.Write(buf.Bytes())
		full.WriteString("\n")
		return nil
	}
}
