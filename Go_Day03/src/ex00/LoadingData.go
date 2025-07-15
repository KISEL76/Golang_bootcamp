package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/olivere/elastic/v7"
)

const mapping = `
{
  "properties": {
    "name": { 
		"type": "text" 
	},
    "address": { 
		"type": "text" 
	},
    "phone": {
		"type": "text" 
	},
    "location": { 
		"type": "geo_point" 
	}
  }
}`

type Restaurant struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location GeoPoint `json:"location"`
}

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func main() {
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		log.Fatalf("Ошибка подключения к Elasticsearch: %s", err)
	}

	for {
		fmt.Println("1 - создание индекса Elasticsearch")
		fmt.Println("2 - создание маппинга Elasticsearch")
		fmt.Println("3 - добавление данных Elasticsearch (Bulk API)")
		fmt.Println("4 - получение данных по ID")
		fmt.Println("5 - удаление индекса Elasticsearch")
		fmt.Println("6 - выход")
		fmt.Printf("Введите номер действия: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			createIndex(client)
		case "2":
			createMapping(client)
		case "3":
			addDataBulk(client)
		case "4":
			getDataByID(client)
		case "5":
			deleteIndex(client)
		case "6":
			return
		default:
			fmt.Printf("Некорректный ввод, попробуйте еще раз.\n")
		}
	}
}

func createIndex(client *elastic.Client) {
	fmt.Print("Введите название индекса: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	indexName := scanner.Text()

	exists, err := client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		log.Fatalf("Ошибка при проверке индекса: %s", err)
	}

	if exists {
		fmt.Println("Такой индекс уже существует.")
		return
	}

	_, err = client.CreateIndex(indexName).Do(context.Background())
	if err != nil {
		log.Fatalf("Ошибка при создании индекса: %s", err)
	}
	fmt.Printf("Индекс успешно создан!\n")
}

func createMapping(client *elastic.Client) {
	fmt.Print("Введите название индекса для маппинга: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	indexName := scanner.Text()

	_, err := client.PutMapping().Index(indexName).BodyString(mapping).Do(context.Background())
	if err != nil {
		log.Fatalf("Ошибка установки маппинга: %s", err)
	}
	fmt.Printf("Маппинг успешно установлен!\n")
}

func addDataBulk(client *elastic.Client) {
	fmt.Print("Введите название индекса для загрузки данных: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	indexName := scanner.Text()
	if indexName == "exit" {
		return
	}

	bulk := client.Bulk()

	for {
		fmt.Println("Ввод данных (напишите 'exit' в любое поле, чтобы выйти и загрузить все данные)")

		fmt.Print("Введите ID: ")
		scanner.Scan()
		idStr := scanner.Text()
		if idStr == "exit" {
			break
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			fmt.Println("Ошибка: ID должен быть числом!")
			continue
		}

		fmt.Print("Введите название: ")
		scanner.Scan()
		name := scanner.Text()
		if name == "exit" {
			break
		}

		fmt.Print("Введите адрес: ")
		scanner.Scan()
		address := scanner.Text()
		if address == "exit" {
			break
		}

		fmt.Print("Введите телефон: ")
		scanner.Scan()
		phone := scanner.Text()
		if phone == "exit" {
			break
		}

		fmt.Print("Введите широту (lat): ")
		scanner.Scan()
		latStr := scanner.Text()
		if latStr == "exit" {
			break
		}
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			fmt.Println("Ошибка: широта должна быть числом!")
			continue
		}

		fmt.Print("Введите долготу (lon): ")
		scanner.Scan()
		lonStr := scanner.Text()
		if lonStr == "exit" {
			break
		}
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			fmt.Println("Ошибка: долгота должна быть числом!")
			continue
		}

		restaurant := Restaurant{
			ID:      id,
			Name:    name,
			Address: address,
			Phone:   phone,
			Location: GeoPoint{
				Lat: lat,
				Lon: lon,
			},
		}

		req := elastic.NewBulkIndexRequest().
			Index(indexName).
			Id(strconv.FormatInt(id, 10)).
			Doc(restaurant)
		bulk = bulk.Add(req)

		fmt.Printf("Данные добавлены в очередь Bulk API (выполнится после 'exit')\n")
	}

	if bulk.NumberOfActions() > 0 {
		response, err := bulk.Do(context.Background())
		if err != nil {
			log.Fatalf("Ошибка при загрузке данных через Bulk API: %s", err)
		}

		for i, item := range response.Items {
			for action, result := range item {
				if result.Error != nil {
					fmt.Printf("Ошибка в документе %d: %s\n", i, result.Error.Reason)
				} else {
					fmt.Printf("Документ %d успешно загружен (операция %s)\n", i, action)
				}
			}
		}
	} else {
		fmt.Printf("Нет данных для загрузки.\n")
	}
}

func getDataByID(client *elastic.Client) {
	fmt.Print("Введите название индекса: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	indexName := scanner.Text()

	fmt.Print("Введите ID документа: ")
	scanner.Scan()
	docID := scanner.Text()

	result, err := client.Get().Index(indexName).Id(docID).Do(context.Background())
	if err != nil {
		fmt.Printf("Ошибка при получении документа: %s\n", err)
		return
	}

	if result.Found {
		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Printf("Ошибка при форматировании JSON: %s\n", err)
			return
		}
		fmt.Println(string(prettyJSON))
	} else {
		fmt.Println("Документ не найден.")
	}
}

func deleteIndex(client *elastic.Client) {
	fmt.Print("Введите название индекса для удаления: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	indexName := scanner.Text()

	_, err := client.DeleteIndex(indexName).Do(context.Background())
	if err != nil {
		fmt.Printf("Ошибка при удалении индекса: %s\n", err)
		return
	}

	fmt.Printf("Индекс успешно удалён!\n")
}
