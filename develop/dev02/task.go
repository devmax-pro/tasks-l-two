/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
  - "a4bc2d5e" => "aaaabccddddde"
  - "abcd" => "abcd"
  - "45" => "" (некорректная строка)
  - "" => ""

Дополнительное задание: поддержка escape - последовательностей
  - qwe\4\5 => qwe45 (*)
  - qwe\45 => qwe44444 (*)
  - qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	testStrings := []string{"a4bc2d5e", "abcd", "45", "", `qwe\4\5`, `qwe\45`, `qwe\\5`}
	for _, str := range testStrings {
		unpacked, err := unpackString(str)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Unpacked: %s\n", unpacked)
		}
	}
}

func unpackString(input string) (string, error) {
	var result strings.Builder
	var escapeMode bool
	var prevRune rune

	for _, r := range input {
		switch {
		case r == '\\' && !escapeMode:
			if escapeMode {
				return "", fmt.Errorf("incorrect string")
			}
			escapeMode = true // Включаем режим escape
		case unicode.IsDigit(r) && !escapeMode:
			if prevRune == 0 {
				return "", fmt.Errorf("incorrect string")
			}
			count, _ := strconv.Atoi(string(r))
			result.WriteString(strings.Repeat(string(prevRune), count-1))
		case unicode.IsLetter(r) || escapeMode:
			result.WriteRune(r)
			prevRune = r
			escapeMode = false
		default:
			return "", fmt.Errorf("incorrect string")
		}
	}
	return result.String(), nil
}
