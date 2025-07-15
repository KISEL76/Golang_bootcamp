package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"ex01/common"
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
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid JSON input")
		return
	}

	if order.CandyCount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid candy count")
		return
	}

	price, ok := prices[order.CandyType]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid candy type")
		return
	}

	totalPrice := price * order.CandyCount
	if order.Money < totalPrice {
		w.WriteHeader(402)
		fmt.Fprintf(w, "You need %d more money!", totalPrice-order.Money)
		return
	}

	if order.Money-totalPrice > 0 {
		w.WriteHeader(201)
		fmt.Fprintf(w, "Thank you! Your change is %d", order.Money-totalPrice)
	} else {
		w.WriteHeader(201)
		fmt.Fprint(w, "Thank you!")
	}
}

func main() {
	caCert, err := os.ReadFile("../cert/minica.pem")
	if err != nil {
		fmt.Printf("It's not possible to read CA: %s\n", err.Error())
		return
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		fmt.Println("It's not possible to add CA to pool")
		return
	}

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caCertPool,
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":3333",
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/buy_candy", buyCandyHandler)
	fmt.Println("HTTPS server with mTLS started at https://candy.tld:3333")

	err = server.ListenAndServeTLS("../cert/candy.tld/cert.pem", "../cert/candy.tld/key.pem")
	if err != nil {
		fmt.Printf("Mistake while launching the server: %s\n", err.Error())
	}
}
