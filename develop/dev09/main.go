/*
=== Утилита wget ===

# Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
package main

import (
	"os"

	"module-wget/runner"
)

func main() {
	os.Exit(runner.CLI(os.Args[1:]))
}
