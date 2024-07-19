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

const (
	serverURL   = "http://localhost:8080/cotacao"
	outputFile  = "./cotacao.txt"
	clientTimeout = 300 * time.Millisecond
)

type CotacaoResponse struct {
	Valor string `json:"valor"`
}

func main() {
	// Configurando um timeout para o cliente
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	// Criando requisição para obter a cotação do servidor
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL, nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição HTTP: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Erro ao realizar requisição HTTP: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erro ao ler resposta HTTP: %v", err)
	}

	var cotacao CotacaoResponse
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		log.Fatalf("Erro ao decodificar resposta JSON: %v", err)
	}

	// Escrevendo a cotação em um arquivo
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Erro ao criar arquivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s\n", cotacao.Valor))
	if err != nil {
		log.Fatalf("Erro ao escrever no arquivo: %v", err)
	}

	log.Println("Cotação salva em cotacao.txt com sucesso!")
}
