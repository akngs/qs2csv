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

type QueryMap map[string]string

func ConvertLines(reader *bufio.Reader, writer *bufio.Writer, noHeader bool, columnNames []string, nilValue string) {
	csvWriter := csv.NewWriter(writer)

	// Print a header
	if !noHeader {
		csvWriter.Write(columnNames)
	}

	fields := make([]string, len(columnNames))
	valueMap := make(QueryMap)

	for {
		// Read a line
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		// Parse querystring
		err = parseQuery(valueMap, strings.TrimRight(line, "\n"))
		if err != nil {
			panic(err)
		}

		// Select columns
		for i, key := range columnNames {
			value, present := valueMap[key]
			if present {
				fields[i] = value
			} else {
				fields[i] = nilValue
			}

			// Reset valueMap in order to reuse it
			valueMap[key] = nilValue
		}

		// Print a row
		csvWriter.Write(fields)
	}
	csvWriter.Flush()
}

func parseQuery(m QueryMap, query string) (err error) {
	for query != "" {
		key := query
		if i := strings.Index(key, "&"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		m[key] = value
	}
	return err
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
