package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

func readFileIntoSet(filename string) (map[string]struct{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
		return nil, err
	}
	defer file.Close()

	set := make(map[string]struct{})
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		set[scanner.Text()] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return set, nil
}

func main() {
	oldFile := pflag.String("old", "", "Path to the old_snapshot")
	newFile := pflag.String("new", "", "Path to the new_snapshot")
	pflag.Parse()

	if (*oldFile == "" || *newFile == "") || (*oldFile == *newFile) {
		fmt.Print("Пример использования: ./compareFS --old <old_snapshot_path> --new <new_snapshot_path>\n")
		os.Exit(1)
	}

	if strings.HasSuffix(*oldFile, ".txt") && strings.HasSuffix(*newFile, ".txt") {
		oldSet, err := readFileIntoSet(*oldFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка чтения снапшота: %v\n", err)
			os.Exit(1)
		}

		newSet, err := readFileIntoSet(*newFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка чтения снапшота: %v\n", err)
			os.Exit(1)
		}

		for file := range newSet {
			if _, exists := oldSet[file]; !exists {
				fmt.Printf("ADDED %s\n", file)
			}
		}

		for file := range oldSet {
			if _, exists := newSet[file]; !exists {
				fmt.Printf("REMOVED %s\n", file)
			}
		}
	} else {
		log.Fatal("Неподдержвиваемый формат. Только .txt поддержтваетсяя.")
	}
}
