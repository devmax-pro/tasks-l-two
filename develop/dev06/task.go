/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	fields    string
	delimiter string
	separated bool
}

func cut(row string, conf *Config) (string, error) {
	// Parse fields to select
	fieldIndexes, err := parseFields(conf.fields)
	if err != nil {
		return "", err
	}
	if conf.separated && !strings.Contains(row, conf.delimiter) {
		return "", nil
	}

	fields := strings.Split(row, conf.delimiter)
	selectedFields := selectFields(fields, fieldIndexes)

	return strings.Join(selectedFields, conf.delimiter), nil
}

// parseFields parses the fields argument and returns a slice of field indexes.
func parseFields(fieldsArg string) ([]int, error) {
	var fields []int
	for _, field := range strings.Split(fieldsArg, ",") {
		var index int
		_, err := fmt.Sscanf(field, "%d", &index)
		if err != nil {
			return nil, err
		}
		fields = append(fields, index-1) // Subtract 1 to convert to zero-based index
	}
	return fields, nil
}

// selectFields selects and returns the requested fields from a slice of fields.
func selectFields(allFields []string, fieldIndexes []int) []string {
	var result []string
	for _, index := range fieldIndexes {
		if index >= 0 && index < len(allFields) {
			result = append(result, allFields[index])
		}
	}
	return result
}

func main() {

	conf := &Config{}
	flag.StringVar(&conf.fields, "f", "", "Fields to select (comma-separated).")
	flag.StringVar(&conf.delimiter, "d", "\t", "delimiter to use for splitting the fields.")
	flag.BoolVar(&conf.separated, "s", false, "Only output lines that contain the delimiter.")
	flag.Parse()

	// Check if fields were provided
	if conf.fields == "" {
		fmt.Fprintln(os.Stderr, "Error: No fields specified.")
		os.Exit(1)
	}

	// Create a scanner to read from STDIN
	scanner := bufio.NewScanner(os.Stdin)

	// Process lines from STDIN
	for scanner.Scan() {
		line := scanner.Text()
		result, err := cut(line, conf)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error cutting line:", err)
			continue
		}

		if result != "" {
			fmt.Println(result)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from input:", err)
	}
}
