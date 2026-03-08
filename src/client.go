package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CotacaoResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	// Define URL padrão como localhost, mas permite override via variável de ambiente
	URL := os.Getenv("URL")
	if URL == "" {
		URL = "http://localhost:8080/cotacao"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		log.Fatalf("Falha ao criar requisição: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("Requisição para %s excedeu timeout: %v", URL, ctx.Err())
		} else {
			log.Printf("Falha ao fazer requisição para %s: %v", URL, err)
		}
		return
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Falha ao ler resposta: %v", err)
	}

	var cotacao CotacaoResponse
	err = json.Unmarshal(data, &cotacao)
	if err != nil {
		log.Fatalf("Falha ao fazer parse do JSON: %v", err)
	}

	// Salvar em arquivo
	conteudo := fmt.Sprintf("Dolar: %s", cotacao.USDBRL.Bid)
	err = os.WriteFile("cotacao.txt", []byte(conteudo), 0644)
	if err != nil {
		log.Fatalf("Falha ao salvar arquivo: %v", err)
	}

	fmt.Printf("Cotação salva: %s", conteudo)
}
