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
	ConvertLines(toReader("a=1&b=2\na=3&b=4\n"), bufio.NewWriter(outBuf), false, []string{"a", "b"}, "")
	equals(t, "a,b\n1,2\n3,4\n", outBuf.String())
}

func TestConvertLinesWithoutHeading(t *testing.T) {
	outBuf := new(bytes.Buffer)
	ConvertLines(toReader("a=1&b=2\na=3&b=4\n"), bufio.NewWriter(outBuf), true, []string{"a", "b"}, "")
	equals(t, "1,2\n3,4\n", outBuf.String())
}

func TestConvertLinesWithNilColumn(t *testing.T) {
	outBuf := new(bytes.Buffer)
	ConvertLines(toReader("a=1&b=2\nb=3&c=4\n"), bufio.NewWriter(outBuf), false, []string{"a", "c"}, "NA")
	equals(t, "a,c\n1,NA\nNA,4\n", outBuf.String())
}

func equals(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Error(fmt.Sprintf("Expected %#v but %#v", expected, actual))
	}
}

func toReader(str string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(str))
}
