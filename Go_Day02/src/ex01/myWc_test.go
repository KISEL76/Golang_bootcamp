package main

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func createTempFile(t *testing.T, content string) string {
	tempFile, err := os.CreateTemp("", "testfile_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.WriteString(content)
	tempFile.Close()
	return tempFile.Name()
}

func TestCountWords1(t *testing.T) {
	fileName := createTempFile(t, "Hello world! This is a test file.")
	defer os.Remove(fileName)

	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Printf("Ожидается: %d\nВывод: ", 7)
	countWords(fileName, &wg)
	wg.Wait()
	fmt.Printf("\n")
}

func TestCountWords2(t *testing.T) {
	fileName := createTempFile(t, "Hello world! This is a test file.\n\n I like triathlon:)")
	defer os.Remove(fileName)

	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Printf("Ожидается: %d\nВывод: ", 10)
	countWords(fileName, &wg)
	wg.Wait()
	fmt.Printf("\n")
}

func TestCountChars1(t *testing.T) {
	fileName := createTempFile(t, "Hello world! This is a test file.")
	defer os.Remove(fileName)

	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Printf("Ожидается: %d\nВывод: ", 33)
	countChars(fileName, &wg)
	wg.Wait()
	fmt.Printf("\n")
}

func TestCountChars2(t *testing.T) {
	fileName := createTempFile(t, "Hello world! This is a test file.\n\n\t I like triathlon:)")
	defer os.Remove(fileName)

	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Printf("Ожидается: %d\nВывод: ", 53)
	countChars(fileName, &wg)
	wg.Wait()
	fmt.Printf("\n")
}

func TestCountLines(t *testing.T) {
	fileName := createTempFile(t, "Hello world!\nThis is a test file.\nAnother line.")
	defer os.Remove(fileName)

	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Printf("\nОжидается: %d\nВывод: ", 3)
	countLines(fileName, &wg)
	wg.Wait()
	fmt.Printf("\n")
}

func TestCountLines2(t *testing.T) {
	fileName := createTempFile(t, "Hello world!\nThis is a test file.\nAnother line.\n\n\nKa-chow")
	defer os.Remove(fileName)

	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Printf("\nОжидается: %d\nВывод: ", 6)
	countLines(fileName, &wg)
	wg.Wait()
	fmt.Printf("\n")
}
