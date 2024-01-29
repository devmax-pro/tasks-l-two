package main

/*
=== Взаимодействие с ОС ===

# Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$: ")

		command, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		command = strings.TrimSpace(command)
		command = strings.TrimSuffix(command, "\n")

		if strings.Contains(command, "|") {
			if err := RunPipeline(command); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		} else {
			if err := RunCommand(command); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		}
	}
}

func RunCommand(command string) error {
	c := strings.Split(command, " ")
	out := os.Stdout
	switch c[0] {
	case "cd":
		if len(c) < 2 {
			dir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			return os.Chdir(dir)
		}
		return os.Chdir(c[1])
	case "pwd":
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Fprintln(out, pwd)
		return nil
	case "echo":
		fmt.Fprintln(out, strings.Join(c[1:], " "))
		return nil
	case "kill":
		if len(c) < 2 {
			return fmt.Errorf("kill: not enough arguments")
		}

		pid, err := strconv.Atoi(c[1])
		p, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
		return p.Kill()
	case "ps":
		cmd := exec.Command("ps", "aux")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	case "exec":
		if len(c) < 2 {
			return fmt.Errorf("exec: not enough arguments")
		}
		binary, err := exec.LookPath(c[1])
		if err != nil {
			return err
		}
		env := os.Environ()
		return syscall.Exec(binary, c[1:], env)

	case "quit":
		fmt.Fprint(out, "exiting from the shell\n")
		os.Exit(0)
	default:
		return fmt.Errorf("command not found: %s", c[0])
	}
	return nil
}

func RunPipeline(command string) error {
	c := strings.Split(command, " | ")
	if len(c) < 2 {
		return fmt.Errorf("pipe: not enough commands: '%v'", c)
	}

	var b bytes.Buffer
	for i := 0; i < len(c); i++ {
		com := exec.Command(c[i])
		commArgs := strings.Split(c[i], " ")
		if len(commArgs) > 1 {
			com = exec.Command(commArgs[0], commArgs[1:]...)
		}

		com.Stdin = bytes.NewReader(b.Bytes())
		b.Reset()
		com.Stdout = &b
		com.Stderr = os.Stderr
		err := com.Start()
		if err != nil {
			return err
		}
		err = com.Wait()
		if err != nil {
			return err
		}
	}

	fmt.Fprint(os.Stdout, b.String())

	return nil
}
