package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestConvertLines(t *testing.T) {
	out := new(bytes.Buffer)
	err := ConvertLines(toReader("a=1&b=2\na=3&b=4\n"), bufio.NewWriter(out), false, []string{"a", "b"}, "")
	equals(t, err, nil)
	equals(t, out.String(), "a,b\n1,2\n3,4\n")
}

func TestConvertLinesWithoutHeading(t *testing.T) {
	out := new(bytes.Buffer)
	err := ConvertLines(toReader("a=1&b=2\na=3&b=4\n"), bufio.NewWriter(out), true, []string{"a", "b"}, "")
	equals(t, err, nil)
	equals(t, out.String(), "1,2\n3,4\n")
}

func TestConvertLinesWithNilColumn(t *testing.T) {
	out := new(bytes.Buffer)
	err := ConvertLines(toReader("a=1&b=2\nb=3&c=4\n"), bufio.NewWriter(out), false, []string{"a", "c"}, "NA")
	equals(t, nil, err)
	equals(t, out.String(), "a,c\n1,NA\nNA,4\n")
}

func TestConvertLinesWithMalformedInput(t *testing.T) {
	out := new(bytes.Buffer)
	err := ConvertLines(toReader("a=1&b=2\n%XW=3&c=4\n"), bufio.NewWriter(out), false, []string{"a", "c"}, "NA")
	equals(t, err != nil, true)
}

func BenchmarkConvertLines(b *testing.B) {
	// Prepare logs
	logs := make([]string, 0, 100000)
	for i := 0; i < cap(logs); i++ {
		logs = append(logs, "a=1&b=2&c=3")
	}
	logString := strings.Join(logs, "\n")

	// Perform benchmark
	b.ResetTimer()
	out := new(bytes.Buffer)
	ConvertLines(toReader(logString), bufio.NewWriter(out), false, []string{"a", "c"}, "NA")
}

func equals(t *testing.T, actual interface{}, expected interface{}) {
	if expected != actual {
		t.Error(fmt.Sprintf("Got %#v want %#v", actual, expected))
	}
}

func toReader(str string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(str))
}
