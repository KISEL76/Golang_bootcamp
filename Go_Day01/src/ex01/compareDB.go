package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"go_day01/models"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

func loadXML(filePath string) (models.Recipes, error) {
	var recipes models.Recipes
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Ошибка чтения XML файла.")
		return recipes, err
	}
	err = xml.Unmarshal(file, &recipes)
	return recipes, err
}

func loadJSON(filePath string) (models.Recipes, error) {
	var recipes models.Recipes
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Ошибка чтения JSON файла")
		return recipes, err
	}
	err = json.Unmarshal(file, &recipes)
	return recipes, err
}

func compareRecipes(oldData, newData models.Recipes) {
	oldCakes := make(map[string]models.Cake)
	newCakes := make(map[string]models.Cake)

	for _, cake := range oldData.Cakes {
		oldCakes[cake.Name] = cake
	}
	for _, cake := range newData.Cakes {
		newCakes[cake.Name] = cake
	}

	for name := range newCakes {
		if _, exists := oldCakes[name]; !exists {
			fmt.Printf("ADDED cake \"%s\"\n", name)
		}
	}

	for name := range oldCakes {
		if _, exists := newCakes[name]; !exists {
			fmt.Printf("REMOVED cake \"%s\"\n", name)
		}
	}

	for name, oldCake := range oldCakes {
		if _, exists := newCakes[name]; !exists {
			continue
		}
		newCake := newCakes[name]

		if oldCake.StoveTime != newCake.StoveTime {
			fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n",
				name, newCake.StoveTime, oldCake.StoveTime)

		}
		compareIngredients(name, oldCake.Ingridients, newCake.Ingridients)
	}
}

func compareIngredients(cakeName string, oldIngredients, newIngredients []models.Ingredient) {
	oldIngr := make(map[string]models.Ingredient)
	newIngr := make(map[string]models.Ingredient)

	for _, ingr := range oldIngredients {
		oldIngr[ingr.Name] = ingr
	}

	for _, ingr := range newIngredients {
		newIngr[ingr.Name] = ingr
	}

	for name := range newIngr {
		if _, exists := oldIngr[name]; !exists {
			fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", name, cakeName)
		}
	}

	for name := range oldIngr {
		if _, exists := newIngr[name]; !exists {
			fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", name, cakeName)
		}
	}

	for name, oldIng := range oldIngr {
		newIng, exists := newIngr[name]
		if !exists {
			continue
		}

		if oldIng.Unit != newIng.Unit && newIng.Unit != "" {
			fmt.Printf("CHANGED unit for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n",
				name, cakeName, newIng.Unit, oldIng.Unit)
		} else if oldIng.Unit != "" {
			fmt.Printf("REMOVED unit \"%s\" for ingridient \"%s\" for cake \"%s\"\n", oldIng.Unit, name, cakeName)
		}

		if oldIng.Count != newIng.Count && oldIng.Unit == newIng.Unit {
			fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n",
				name, cakeName, newIng.Count, oldIng.Count)
		}
	}
}

func main() {
	firstFile := pflag.String("old", "", "Path to the input file (XML or JSON)")
	secondFile := pflag.String("new", "", "Path to the input file (XML or JSON)")
	pflag.Parse()

	if (*firstFile == "" || *secondFile == "") || (*firstFile == *secondFile) {
		fmt.Print("Пример использования: ./compareDB --old <old_file_path> --new <new_file_path>\n")
		os.Exit(1)
	}

	if strings.HasSuffix(*firstFile, ".xml") || strings.HasSuffix(*firstFile, ".json") && strings.HasSuffix(*secondFile, ".json") ||
		strings.HasSuffix(*secondFile, ".xml") {
		var oldData models.Recipes
		var err error

		if strings.HasSuffix(*firstFile, ".xml") {
			oldData, err = loadXML(*firstFile)
		} else {
			oldData, err = loadJSON(*firstFile)
		}
		if err != nil {
			log.Fatalf("Ошибка открытия %s файла.", *firstFile)
			os.Exit(1)
		}

		var newData models.Recipes
		var err1 error

		if strings.HasSuffix(*secondFile, ".json") {
			newData, err1 = loadJSON(*secondFile)
		} else {
			newData, err1 = loadXML(*secondFile)
		}
		if err1 != nil {
			log.Fatalf("Ошибка открытия %s файла.", *secondFile)
			os.Exit(1)
		}

		compareRecipes(oldData, newData)
	} else {
		log.Fatalf("Неподдержвиваемый формат. Только .json и .xml поддерживаются.")
	}
}
