package main


import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"net/url"
)


func main() {
	// Parse command-line arguments
	columnNamesPtr := flag.String("c", "", "comma-separated list of column names")
	nilValuePtr := flag.String("n", "", "a string represents nil value")
	noHeaderPtr := flag.Bool("no-header", false, "do not print header")
	flag.Parse()

	columnsNames := strings.Split(*columnNamesPtr, ",")
	if len(columnsNames) == 0 {
		fmt.Println("Provide one or more columns using -f flag")
		os.Exit(1)
	}

	// Print a header
	if !*noHeaderPtr {
		fmt.Println(strings.Join(columnsNames, ","))
	}

	reader := bufio.NewReader(os.Stdin)
	writer := csv.NewWriter(os.Stdout)
	fields := make([]string, len(columnsNames))

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
		for i, key := range columnsNames {
			value, present := valueMap[key]
			if present {
				fields[i] = value[0]
			} else {
				fields[i] = *nilValuePtr
			}
		}

		// Print a row
		writer.Write(fields)
	}
	writer.Flush()
}
