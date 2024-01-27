/*
=== Базовая задача ===

Создать программу, печатающую точное время с использованием NTP библиотеки. Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу, печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и go lint.
*/
package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"time"
)

func ntpTime() error {
	currTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return err
	}

	fmt.Println("The current time is", currTime.UTC().Format(time.UnixDate))
	return nil
}

func main() {
	err := ntpTime()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
