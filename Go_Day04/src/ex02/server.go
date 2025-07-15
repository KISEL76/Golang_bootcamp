package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"ex02/common"
	"fmt"
	"net/http"
	"os"
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
		http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
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
		return
	}

	price, ok := prices[order.CandyType]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.ErrorResponse{Error: "Invalid candy type"})
		return
	}

	total := price * order.CandyCount
	if order.Money < total {
		w.WriteHeader(402)
		json.NewEncoder(w).Encode(common.ErrorResponse{
			Error: fmt.Sprintf("You need %d more money!", total-order.Money),
		})
		return
	}

	resp := common.SuccessResponse{
		Thanks: common.Cowify("Thank you!"),
		Change: order.Money - total,
	}
	data, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(data))
}

func main() {
	caCert, err := os.ReadFile("cert/minica.pem")
	if err != nil {
		fmt.Printf("Failed to read CA: %s\n", err.Error())
		return
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		fmt.Println("Failed to append CA cert")
		return
	}

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caPool,
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      "candy.tld:3333",
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/buy_candy", buyCandyHandler)

	fmt.Println("Server running at https://candy.tld:3333")
	err = server.ListenAndServeTLS("cert/candy.tld/cert.pem", "cert/candy.tld/key.pem")
	if err != nil {
		fmt.Printf("Server error: %s\n", err.Error())
	}
}
