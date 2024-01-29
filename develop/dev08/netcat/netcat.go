package main

/*
=== Взаимодействие с ОС ===

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	isUdp := flag.Bool("u", false, "Use UDP instead of the default option of TCP.")
	flag.Parse()
	if len(flag.Args()) < 2 {
		fmt.Println("Hostname and port required")
		return
	}
	serverHost := flag.Arg(0)
	serverPort := flag.Arg(1)
	err := runClient(*isUdp, serverHost, serverPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func runClient(isUdp bool, host, port string) error {
	var connType string
	if isUdp {
		connType = "udp"
	} else {
		connType = "tcp"
	}

	conn, err := net.Dial(connType, net.JoinHostPort(host, port))
	if err != nil {
		return fmt.Errorf("Can't connect to server: %s\n", err)
	}
	defer conn.Close()

	go func() {
		io.Copy(os.Stdout, conn)
	}()

	if err != nil {
		return fmt.Errorf("Connection error: %s\n", err)
	}
	fmt.Printf("Connected to %s\n", conn.RemoteAddr().String())

	// Читаем stdin и отправляем данные в соединение
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			return fmt.Errorf("error sending data: %w", err)
		}
	}

	if scanner.Err() != nil {
		return fmt.Errorf("error reading from stdin: %w", scanner.Err())
	}

	return nil
}
