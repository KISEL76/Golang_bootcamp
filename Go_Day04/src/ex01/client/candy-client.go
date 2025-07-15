package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"ex01/common"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	candyType := flag.String("k", "", "Candy type")
	count := flag.Int("c", 1, "Number of candies")
	money := flag.Int("m", 0, "Amount of money given")
	flag.Parse()

	if *candyType == "" || *count <= 0 || *money <= 0 {
		fmt.Println("Usage: ./candy-client -k AA -c 2 -m 50")
		os.Exit(1)
	}

	cert, err := tls.LoadX509KeyPair("../cert/client/cert.pem", "../cert/client/key.pem")
	if err != nil {
		fmt.Printf("Didn't manage to download the client CA: %s", err.Error())
		return
	}

	caCert, err := os.ReadFile("../cert/minica.pem")
	if err != nil {
		fmt.Printf("Didn't manage to read CA: %s", err.Error())
		return
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		fmt.Println("CA wasn't added in the pool")
		return
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	order := common.CandyOrder{
		Money:      *money,
		CandyType:  *candyType,
		CandyCount: *count,
	}

	body, err := json.Marshal(order)
	if err != nil {
		fmt.Printf("Serialization error: %s", err.Error())
		return
	}

	resp, err := client.Post("https://candy.tld:3333/buy_candy", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Request error: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read the response error: %s", err.Error())
		return
	}

	fmt.Println(string(responseBody))
}
