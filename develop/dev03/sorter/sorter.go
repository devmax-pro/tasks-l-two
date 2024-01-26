package sorter

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type inputFlags struct {
	column     int
	numeric    bool
	reverse    bool
	unique     bool
	outputFile string
	inputFiles []string
}

func parseFlags() inputFlags {

	column := flag.Int("k", 0, "specify the column to be sorted")
	num := flag.Bool("n", false, "sort by numeric value")
	reverse := flag.Bool("r", false, "sort in reverse order")
	unique := flag.Bool("u", false, "do not print duplicate lines")
	outputFile := flag.String("o", "", "output file to write sorted lines")
	flag.Parse()

	flags := inputFlags{
		column:     *column,
		numeric:    *num,
		reverse:    *reverse,
		unique:     *unique,
		outputFile: *outputFile,
		inputFiles: flag.Args(),
	}
	return flags
}

func readFromFiles(inputFiles []string) ([]string, error) {

	if len(inputFiles) == 0 {
		return nil, errors.New("input file must be specified")
	}

	var data []string
	for _, inputFile := range inputFiles {
		if inputFile == "" {
		}

		lines, err := readLines(inputFile)
		if err != nil {
			return []string{}, err
		}
		data = append(data, lines...)
	}

	return data, nil
}

func SortFiles() {
	flags := parseFlags()
	lines, err := readFromFiles(flags.inputFiles)
	sortedLines := sortLines(lines, flags)

	if flags.outputFile == "" {
		writeToStdout(sortedLines)
		return
	}

	err = writeToFile(sortedLines, flags.outputFile)
	if err != nil {
		fmt.Println("Error writing:", err)
		os.Exit(1)
	}
}

func writeToStdout(data []string) {
	for _, v := range data {
		fmt.Fprintln(os.Stdout, v)
	}
}

func writeToFile(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := fmt.Fprintln(w, line)
		if err != nil {
			return err
		}
	}
	return w.Flush()
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func sortLines(lines []string, flags inputFlags) []string {

	if flags.unique {
		lines = removeDuplicates(lines)
	}

	sort.Slice(lines, func(i, j int) bool {
		if flags.numeric {
			return compareNumbers(lines[i], lines[j], flags.column)
		}

		if flags.column > 0 {
			return compareColumn(lines[i], lines[j], flags.column)
		}

		return lines[i] < lines[j]
	})

	if flags.reverse {
		reverseLines(lines)
	}
	return lines
}

func reverseLines(lines []string) {
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
}

func compareNumbers(line1, line2 string, col int) bool {
	num1, err1 := extractNumber(line1, col)
	num2, err2 := extractNumber(line2, col)

	if err1 != nil || err2 != nil {
		return line1 < line2
	}

	return num1 < num2
}

func compareColumn(line1, line2 string, col int) bool {

	cols1 := strings.Fields(line1)
	cols2 := strings.Fields(line2)

	if len(cols1) == 0 {
		return true
	}

	if len(cols2) == 0 {
		return false
	}

	if len(cols1) < col && len(cols2) >= col {
		return true
	}

	if len(cols1) >= col && len(cols2) < col {
		return false
	}

	if len(cols1) < col && len(cols2) < col {
		return cols1[0] < cols2[0]
	}

	return cols1[col-1] < cols2[col-1]
}

func extractNumber(line string, col int) (float64, error) {
	if col == 0 {
		return strconv.ParseFloat(line, 64)
	}

	nums := strings.Fields(line)
	if col >= len(nums) {
		return 0, fmt.Errorf("column index out of range")
	}
	return strconv.ParseFloat(nums[col-1], 64)
}

func removeDuplicates(lines []string) []string {
	seen := make(map[string]bool)
	var uniqueLines []string
	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			uniqueLines = append(uniqueLines, line)
		}
	}
	return uniqueLines
}
