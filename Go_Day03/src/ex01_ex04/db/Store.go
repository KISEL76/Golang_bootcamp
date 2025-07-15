package db

import (
	"context"
	"encoding/json"
	"log"

	"github.com/olivere/elastic/v7"
)

type Store interface {
	GetPlaces(limit, offset int) ([]Place, int, error)
	GetElementsCount() int64
	GetClosestPlaces(lat, lon float64, limit int) ([]Place, error)
}

type Place struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type ElasticStore struct {
	client    *elastic.Client
	indexName string
}

func NewElasticStore(url, index string) (*ElasticStore, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		return nil, err
	}
	return &ElasticStore{client: client, indexName: index}, nil
}

func (es *ElasticStore) GetPlaces(limit, offset int) ([]Place, int, error) {
	searchResult, err := es.client.Search().
		Index(es.indexName).
		From(offset).Size(limit).
		Do(context.Background())

	if err != nil {
		log.Printf("Ошибка поиска в Elasticsearch: %v", err)
		return nil, 0, err
	}

	var places []Place
	for _, hit := range searchResult.Hits.Hits {
		var place Place
		if err := json.Unmarshal(hit.Source, &place); err == nil {
			places = append(places, place)
		} else {
			log.Printf("Ошибка парсинга JSON: %v", err)
		}
	}

	total := int(searchResult.TotalHits())

	log.Printf("Найдено записей: %d, Всего записей: %d", len(places), total)

	return places, total, nil
}

func (es *ElasticStore) GetElementsCount() int64 {
	count, err := es.client.Count("places").Do(context.Background())
	if err != nil {
		log.Fatalln("Ошибка подсчета количества элементов индекса")
	}

	return count
}

func (es *ElasticStore) GetClosestPlaces(lat, lon float64, limit int) ([]Place, error) {
	log.Printf("Поиск ближайших мест для координат: lat=%.6f, lon=%.6f", lat, lon)

	searchResult, err := es.client.Search().
		Index(es.indexName).
		SortBy(elastic.NewGeoDistanceSort("location").
			Point(lat, lon).
			Order(true).
			Unit("km").
			DistanceType("arc"),
		).
		Query(elastic.NewGeoDistanceQuery("location").
			Point(lat, lon).
			Distance("50km"),
		).
		Size(limit).
		Do(context.Background())

	if err != nil {
		log.Printf("Ошибка поиска ближайших мест: %v", err)
		return nil, err
	}

	var places []Place
	for _, hit := range searchResult.Hits.Hits {
		var place Place
		if err := json.Unmarshal(hit.Source, &place); err == nil {
			places = append(places, place)
		} else {
			log.Printf("Ошибка парсинга JSON: %v", err)
		}
	}

	return places, nil
}
