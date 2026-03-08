package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Currency struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type ResponseAPI struct {
	USDBRL Currency `json:"USDBRL"`
}

// Escrever um JSON simples: {status: "error", message: "...", data: []}
func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": msg,
		"data":    nil,
	})
}

var db *sql.DB

func main() {
	var err error

	db, err = sql.Open("sqlite3", "/data/cotacao.db")
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	err = createTable()
	if err != nil {
		log.Fatalf("Falha ao criar tabela: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", handleCurrencyAPI)

	port := ":8080"
	fmt.Printf("Servidor rodando em http://localhost%s\n", port)
	http.ListenAndServe(port, mux)
}

func createTable() error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS cotacoes(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        code TEXT,
        codeIn TEXT,
        name TEXT,
        high TEXT,
        low TEXT,
        varBid TEXT,
        pctChange TEXT,
        bid TEXT,
        ask TEXT,
        timestamp TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err := db.Exec(createTableSQL)
	return err
}

func handleCurrencyAPI(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	URL := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	go func() {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				log.Printf("%s %s chamado com timeout: %v", r.URL.Path, r.RemoteAddr, ctx.Err())
			}
		}
	}()

	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "falha para buscar a cotação")
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			writeError(w, http.StatusGatewayTimeout, "Serviço de cotação demorou demais para responder!")
			return
		}
		writeError(w, http.StatusBadGateway, "Falha ao buscar a cotação!")
		return
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Servidor demorou para ler a resposta da cotação!")
		return
	}

	var resp ResponseAPI
	err = json.Unmarshal(data, &resp)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Falha ao processar a resposta da cotação!")
		return
	}

	err = add_on_db(db, resp.USDBRL)
	if err != nil {
		log.Printf("Erro ao adicionar cotação ao banco de dados: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func add_on_db(db *sql.DB, currency Currency) error {
	if db != nil {
		dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer dbCancel()
		if _, err := db.ExecContext(dbCtx, `INSERT INTO cotacoes(code, codeIn, name, high, low, varBid, pctChange, bid, ask, timestamp) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, currency.Code, currency.Codein, currency.Name, currency.High, currency.Low, currency.VarBid, currency.PctChange, currency.Bid, currency.Ask, currency.Timestamp); err != nil {
			if dbCtx.Err() == context.DeadlineExceeded {
				log.Printf("/cotacao: persistência excedeu timeout: %v", dbCtx.Err())
			} else {
				log.Printf("/cotacao: falha ao persistir cotação: %v", err)
			}
		}
	}
	return nil
}
