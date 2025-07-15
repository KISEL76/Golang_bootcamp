package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"go_day01/models"
	"io"
	"log"
	"os"
	"strings"
)

func handleXML(filePath string) {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Ошибка открытия XML файла: %v", err)
	}
	defer xmlFile.Close()

	xmlData, err := io.ReadAll(xmlFile)
	if err != nil {
		log.Fatalf("Ошибка чтения XML файла: %v", err)
	}

	var recipes models.Recipes
	err = xml.Unmarshal(xmlData, &recipes)
	if err != nil {
		log.Fatalf("Ошибка обработки: %v", err)
	}

	jsonData, err := json.MarshalIndent(recipes, "", "    ")
	if err != nil {
		log.Fatalf("Ошибка кодирования в формат JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}

func handleJSON(filePath string) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Ошибка открытия JSON файла: %v", err)
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Ошибка чтения JSON файла: %v", err)
	}

	var recipes models.Recipes
	err = json.Unmarshal(jsonData, &recipes)
	if err != nil {
		log.Fatalf("Ошибка обработки JSON: %v", err)
	}

	xmlData, err := xml.MarshalIndent(recipes, "", "    ")
	if err != nil {
		log.Fatalf("Ошибка кодирования в формат XML: %v", err)
	}

	fmt.Println(string(xmlData))
}

func main() {
	filePath := flag.String("f", "", "Path to the input file (XML or JSON)")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Пример использования: ./readDB -f <file_path>")
		os.Exit(1)
	}

	if strings.HasSuffix(*filePath, ".xml") {
		handleXML(*filePath)
	} else if strings.HasSuffix(*filePath, ".json") {
		handleJSON(*filePath)
	} else {
		log.Fatalf("Неподдерживаемый формат: %s. Только .xml и .json поддерживаются.", *filePath)
	}
}
