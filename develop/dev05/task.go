/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

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

type Flags struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
	pattern    string
	fileName   string
}

func parseFlags() (*Flags, error) {
	flags := &Flags{}

	flag.IntVar(&flags.after, "A", 0, "Print N lines after match")
	flag.IntVar(&flags.before, "B", 0, "Print N lines before match")
	flag.IntVar(&flags.context, "C", 0, "Print N lines around the match")
	flag.BoolVar(&flags.count, "c", false, "Print the number of lines")
	flag.BoolVar(&flags.ignoreCase, "i", false, "Ignore case")
	flag.BoolVar(&flags.invert, "v", false, "Invert match")
	flag.BoolVar(&flags.fixed, "F", false, "Exact match with line, not a pattern")
	flag.BoolVar(&flags.lineNum, "n", false, "Print line number")

	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		return nil, fmt.Errorf("usage: grep [flags] pattern file")
	}

	flags.pattern = args[0]
	flags.fileName = args[1]

	return flags, nil
}

func grep(flags *Flags) error {
	file, err := os.Open(flags.fileName)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var (
		beforeLines  []string
		afterCounter int
		matchCount   int
		lineNumber   int
	)

	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		match := lineMatches(flags.pattern, line, flags.ignoreCase, flags.fixed, flags.invert)

		if flags.before > 0 || flags.context > 0 {
			if match {
				for _, bLine := range beforeLines {
					printLine(bLine, lineNumber-len(beforeLines), flags.lineNum)
				}
				beforeLines = nil
			} else if len(beforeLines) < flags.before+flags.context {
				beforeLines = append(beforeLines, line)
			} else {
				beforeLines = append(beforeLines[1:], line)
			}
		}

		if match {
			matchCount++
			if !flags.count {
				printLine(line, lineNumber, flags.lineNum)
			}
			if flags.after > 0 || flags.context > 0 {
				afterCounter = flags.after + flags.context
			}
		} else if afterCounter > 0 {
			afterCounter--
			printLine(line, lineNumber, flags.lineNum)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	if flags.count {
		fmt.Println(matchCount)
	}

	return nil
}

func lineMatches(pattern, line string, ignoreCase, fixed, invert bool) bool {
	if ignoreCase {
		line = strings.ToLower(line)
		pattern = strings.ToLower(pattern)
	}

	var match bool
	if fixed {
		match = line == pattern
	} else {
		match = strings.Contains(line, pattern)
	}

	if invert {
		return !match
	}
	return match
}

func printLine(line string, lineNumber int, printNum bool) {
	if printNum {
		fmt.Printf("%d: %s\n", lineNumber, line)
	} else {
		fmt.Println(line)
	}
}

func main() {
	flags, err := parseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = grep(flags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
