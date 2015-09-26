package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

func ConvertLines(reader *bufio.Reader, writer *bufio.Writer, noHeader bool, columnNames []string, nilValue string) {
	csvWriter := csv.NewWriter(writer)
	fields := make([]string, len(columnNames))

	// Print a header
	if !noHeader {
		csvWriter.Write(columnNames)
	}

	for {
		// Read a line
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		// Parse querystring
		valueMap, err := url.ParseQuery(strings.TrimRight(line, "\n"))
		if err != nil {
			panic(err)
		}

		// Select columns
		for i, key := range columnNames {
			value, present := valueMap[key]
			if present {
				fields[i] = value[0]
			} else {
				fields[i] = nilValue
			}
		}

		// Print a row
		csvWriter.Write(fields)
	}
	csvWriter.Flush()
}

func main() {
	columnNamesPtr := flag.String("c", "", "comma-separated list of column names")
	nilValuePtr := flag.String("n", "", "a string represents nil value")
	noHeaderPtr := flag.Bool("no-header", false, "do not print header")
	flag.Parse()

	columnNames := strings.Split(*columnNamesPtr, ",")
	if len(columnNames) == 0 {
		fmt.Println("Provide one or more columns using -f flag")
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	ConvertLines(reader, writer, *noHeaderPtr, columnNames, *nilValuePtr)
}
