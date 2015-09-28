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

type queryMap map[string]string

// ConvertLines takes lines of URL-encoded querystring and converts each line
// into CSV form.
func ConvertLines(in *bufio.Reader, out *bufio.Writer, noHeader bool, colNames []string, nilVal string) error {
	csvWriter := csv.NewWriter(out)

	// Print a header
	if !noHeader {
		csvWriter.Write(colNames)
	}

	fields := make([]string, len(colNames))

	// Fill query map to default nil values to skip present-checking
	queryMap := make(queryMap)
	for _, key := range colNames {
		queryMap[key] = nilVal
	}

	for {
		// Read a line
		line, err := in.ReadString('\n')
		if err == io.EOF {
			break
		}

		// Parse querystring
		err = parseQuery(queryMap, strings.TrimRight(line, "\n"))
		if err != nil {
			return err
		}

		// Select columns
		for i, key := range colNames {
			fields[i] = queryMap[key]

			// Reset query map to reuse it
			queryMap[key] = nilVal
		}

		// Print a row
		csvWriter.Write(fields)
	}
	csvWriter.Flush()

	return nil
}

func parseQuery(m queryMap, query string) (err error) {
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
	pColNames := flag.String("c", "", "comma-separated list of column names")
	pNilVal := flag.String("n", "", "a string represents nil value")
	pNoHeader := flag.Bool("no-header", false, "do not print header")
	flag.Parse()

	colNames := strings.Split(*pColNames, ",")
	if len(colNames) == 0 {
		fmt.Println("Provide one or more columns using -c flag")
		os.Exit(1)
	}

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	ConvertLines(in, out, *pNoHeader, colNames, *pNilVal)
}
