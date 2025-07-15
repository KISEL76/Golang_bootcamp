package main

import (
	"encoding/json"
	"ex00/common"
	"fmt"
	"net/http"
)

var prices = map[string]int{
	"CE": 10,
	"AA": 15,
	"NT": 17,
	"DE": 21,
	"YR": 23,
}

func buyCandyHandler(w http.ResponseWriter, r *http.Request) {
	var order common.CandyOrder

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.ErrorResponse{Error: "Invalid JSON input"})
		return
	}

	if order.CandyCount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.ErrorResponse{Error: "Invalid candy count!"})
	}

	price, ok := prices[order.CandyType]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.ErrorResponse{Error: "Invalid candy type"})
		return
	}

	totalPrice := price * order.CandyCount
	if order.Money < totalPrice {
		w.WriteHeader(402)
		json.NewEncoder(w).Encode(common.ErrorResponse{Error: fmt.Sprintf("You need %d more money!", totalPrice-order.Money)})
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(common.SuccessResponse{
		Thanks: "Thank you!",
		Change: order.Money - totalPrice,
	})
}

func main() {
	http.HandleFunc("/buy_candy", buyCandyHandler)
	fmt.Println("Server started at http://localhost:3333")
	http.ListenAndServe(":3333", nil)
}
