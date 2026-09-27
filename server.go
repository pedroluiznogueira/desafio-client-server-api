package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type exchangeResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type cotacaoResponse struct {
	Bid string `json:"bid"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", "cotacoes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		created_at DATETIME
	)`); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/cotacao", cotacaoHandler)

	log.Println("server ouvindo em :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	bid, err := fetchExchangeRate(r.Context())
	if err != nil {
		log.Printf("erro ao buscar cotacao: %v", err)
		http.Error(w, "erro ao buscar cotacao", http.StatusInternalServerError)
		return
	}

	if err := saveQuote(r.Context(), bid); err != nil {
		log.Printf("erro ao salvar cotacao no banco: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cotacaoResponse{Bid: bid})
}

func fetchExchangeRate(parent context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(parent, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body exchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}

	return body.USDBRL.Bid, nil
}

func saveQuote(parent context.Context, bid string) error {
	ctx, cancel := context.WithTimeout(parent, 10*time.Millisecond)
	defer cancel()

	_, err := db.ExecContext(ctx, "INSERT INTO cotacoes (bid, created_at) VALUES (?, ?)", bid, time.Now())
	return err
}
