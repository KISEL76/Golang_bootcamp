package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: ./myXargs <команда> [аргументы...]")
		os.Exit(1)
	}

	var args []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		args = append(args, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения stdin:", err)
		os.Exit(1)
	}

	cmd := exec.Command(os.Args[1], append(os.Args[2:], args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Ошибка выполнения команды:", err)
		os.Exit(1)
	}
}
