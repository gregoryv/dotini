package ingrid

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

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
