package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

func validateCommand(flag1, flag2, flag3 *bool) int {
	flagCount := 0

	if *flag1 {
		flagCount++
	}
	if *flag2 {
		flagCount++
	}
	if *flag3 {
		flagCount++
	}
	return flagCount
}

func countWords(fileName string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка открытия файла %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	wordCount := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		words := strings.FieldsFunc(line, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})
		wordCount += len(words)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка чтения файла %s: %v\n", fileName, err)
	}
	fmt.Printf("%s\t%d\n", fileName, wordCount)
}

func countChars(fileName string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка открытия файла %s: %v\n", fileName, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	countChars := 0

	for scanner.Scan() {
		line := scanner.Text()
		countChars += utf8.RuneCountInString(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка чтения файла %s: %v\n", fileName, err)
	}
	fmt.Printf("%s\t%d\n", fileName, countChars)
}

func countLines(fileName string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка открытия файла %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	linesCount := 0

	for scanner.Scan() {
		linesCount++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка чтения файла %s: %v\n", fileName, err)
		return
	}
	fmt.Printf("%s\t%d\n", fileName, linesCount)
}

func wFlagWorks() {
	var wg sync.WaitGroup

	for _, fileName := range flag.Args()[0:] {
		wg.Add(1)
		go countWords(fileName, &wg)
	}
	wg.Wait()
}

func cFlagWorks() {
	var wg sync.WaitGroup

	for _, fileName := range flag.Args()[0:] {
		wg.Add(1)
		go countChars(fileName, &wg)
	}
	wg.Wait()
}

func lFlagWorks() {
	var wg sync.WaitGroup

	for _, fileName := range flag.Args()[0:] {
		wg.Add(1)
		go countLines(fileName, &wg)
	}
	wg.Wait()
}

func main() {
	flagW := flag.Bool("w", false, "to count words")
	flagC := flag.Bool("c", false, "to count characters")
	flagL := flag.Bool("l", false, "to count lines")
	flag.Parse()

	flagCount := validateCommand(flagW, flagC, flagL)

	if flagCount == 0 {
		wFlagWorks()
	} else if flagCount == 1 {
		if *flagW {
			wFlagWorks()
		} else if *flagC {
			cFlagWorks()
		} else if *flagL {
			lFlagWorks()
		}
	} else {
		fmt.Fprintf(os.Stderr, "Ошибка: утилита используется только с одним флагом\n")
	}
}
