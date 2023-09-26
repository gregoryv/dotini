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
			mapping := newHandler(t, &full)
			Map(mapping, bufio.NewScanner(strings.NewReader(input)))
			if got := full.String(); !strings.Contains(got, ErrSyntax.Error()) {
				t.Error(got)
			}
		})
	}
	assert("missing right bracket", `[test`)
	assert("missing equal sign", `key`)
	assert("space in key", `key 1 = value`)
	assert("missing quote", `key1 = "value`)
}

func Test_Map_allowed(t *testing.T) {
	assert := func(label, input string) {
		t.Helper()
		t.Run(label, func(t *testing.T) {
			var full bytes.Buffer
			mapping := newHandler(t, &full)
			Map(mapping, bufio.NewScanner(strings.NewReader(input)))
			if got := full.String(); strings.Contains(got, ErrSyntax.Error()) {
				t.Error(got)
			}
		})
	}

	assert("grub1",
		`GRUB_CMDLINE_LINUX_DEFAULT="quiet splash acpi_enf`+
			`orce_resources=lax snd-intel-dspcfg.dsp_driver=1"`,
	)
}

func Test_Map_cfg(t *testing.T) {
	example, _ := os.ReadFile("testdata/example.cfg")
	var full bytes.Buffer
	mapping := newHandler(t, &full)
	Map(mapping, bufio.NewScanner(bytes.NewReader(example)))
	golden.Assert(t, full.String())
}

func Test_Map_ini(t *testing.T) {
	example, _ := os.ReadFile("testdata/example.ini")
	var full bytes.Buffer
	mapping := newHandler(t, &full)
	Map(mapping, bufio.NewScanner(bytes.NewReader(example)))
	golden.Assert(t, full.String())
}

func Benchmark_Map(b *testing.B) {
	example, _ := os.ReadFile("testdata/example.ini")
	mapping := func(section, key, value, comment string, err error) {}
	r := bytes.NewReader(example)
	scanner := bufio.NewScanner(r)
	for i := 0; i < b.N; i++ {
		Map(mapping, scanner)
		r.Reset(example)
	}
}

func newHandler(t *testing.T, full *bytes.Buffer) Mapfn {
	return func(section, key, value, comment string, err error) {
		var buf bytes.Buffer
		if len(section) > 0 {
			fmt.Fprintf(&buf, "[%s]", section)
		}
		if len(key) > 0 {
			fmt.Fprint(&buf, key, "=")
			fmt.Fprintf(&buf, value)
		}
		if len(comment) > 0 {
			fmt.Fprintf(&buf, comment)
		}
		if err != nil {
			fmt.Fprint(&buf, err)
		}
		full.Write(buf.Bytes())
		full.WriteString("\n")

	}
}
