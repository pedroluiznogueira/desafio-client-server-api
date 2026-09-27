package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type cotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("erro ao buscar cotacao do servidor: %v", err)
	}
	defer resp.Body.Close()

	var body cotacaoResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		log.Fatal(err)
	}

	content := "Dólar: " + body.Bid
	if err := os.WriteFile("cotacao.txt", []byte(content), 0644); err != nil {
		log.Fatal(err)
	}

	log.Println(content)
}
