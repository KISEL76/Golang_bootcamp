package main

import (
	"fmt"
	"log"
	"net/http"

	"ex01_04/db"
	"ex01_04/handlers"
	"ex01_04/middleware"
)

func main() {
	store, err := db.NewElasticStore("http://localhost:9200", "places")
	if err != nil {
		log.Fatalf("Ошибка подключения к базе: %s", err)
	}

	http.HandleFunc("/api/get_token", middleware.GetTokenHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderPage(w, r, store)
	})
	http.HandleFunc("/api/places", func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderJSON(w, r, store)
	})
	http.HandleFunc("/api/recommend", middleware.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderRecommendations(w, r, store)
	}))

	fmt.Println("Сервер запущен на http://127.0.0.1:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
