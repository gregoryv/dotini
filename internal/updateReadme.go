package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var inFile = flag.String("in", "doc.go", "input file")
var exFile = flag.String("example", "map_test.go", "example input file")
var outFile = flag.String("out", "README.md", "output file")

func main() {
	var buf bytes.Buffer

	fmt.Fprintln(&buf,
		"<!-- GENERATED, DO NOT EDIT! See internal/updateReadme.go -->\n",
		`<img src="./internal/banner.png">`,
	)

	appendDoc(&buf, *inFile)
	appendExample(&buf, *exFile)
	appendBenchmark(&buf)

	if err := os.WriteFile(*outFile, buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
}
func appendExample(buf *bytes.Buffer, filename string) {
	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	fmt.Fprintln(buf, "## Example")
	fmt.Fprintln(buf)
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintln(buf, "   ", line)
	}
	fmt.Fprintln(buf)
}

func appendBenchmark(buf *bytes.Buffer) {
	cmd := exec.Command("go", "test", "-benchmem", "-bench", ".")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(buf, "## Benchmark")
	fmt.Fprintln(buf)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if skipAny(line, "PASS", "ok ") {
			continue
		}
		line = strings.ReplaceAll(line, "\t", "")
		line = strings.ReplaceAll(line, "  ", " ")
		fmt.Fprintln(buf, "    ", line)
	}
}

func appendDoc(buf *bytes.Buffer, filename string) {
	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		if skipAny(line, "/*", "*/", "//go:generate", "package ingrid") {
			continue
		}
		if strings.HasPrefix(line, "# ") {
			// add level
			fmt.Fprintf(buf, "#%s\n", line)
		} else {
			fmt.Fprintln(buf, line)
		}
	}
}

func skipAny(line string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(line, p) {
			return true
		}
	}
	return false
}
