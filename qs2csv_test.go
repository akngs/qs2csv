package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestConvertLines(t *testing.T) {
	outBuf := new(bytes.Buffer)
	err := ConvertLines(toReader("a=1&b=2\na=3&b=4\n"), bufio.NewWriter(outBuf), false, []string{"a", "b"}, "")
	equals(t, nil, err)
	equals(t, "a,b\n1,2\n3,4\n", outBuf.String())
}

func TestConvertLinesWithoutHeading(t *testing.T) {
	outBuf := new(bytes.Buffer)
	err := ConvertLines(toReader("a=1&b=2\na=3&b=4\n"), bufio.NewWriter(outBuf), true, []string{"a", "b"}, "")
	equals(t, nil, err)
	equals(t, "1,2\n3,4\n", outBuf.String())
}

func TestConvertLinesWithNilColumn(t *testing.T) {
	outBuf := new(bytes.Buffer)
	err := ConvertLines(toReader("a=1&b=2\nb=3&c=4\n"), bufio.NewWriter(outBuf), false, []string{"a", "c"}, "NA")
	equals(t, nil, err)
	equals(t, "a,c\n1,NA\nNA,4\n", outBuf.String())
}

func TestConvertLinesWithMalformedInput(t *testing.T) {
	outBuf := new(bytes.Buffer)
	err := ConvertLines(toReader("a=1&b=2\n%XW=3&c=4\n"), bufio.NewWriter(outBuf), false, []string{"a", "c"}, "NA")
	equals(t, true, err != nil)
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
	outBuf := new(bytes.Buffer)
	ConvertLines(toReader(logString), bufio.NewWriter(outBuf), false, []string{"a", "c"}, "NA")
}

func equals(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Error(fmt.Sprintf("Expected %#v but %#v", expected, actual))
	}
}

func toReader(str string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(str))
}
