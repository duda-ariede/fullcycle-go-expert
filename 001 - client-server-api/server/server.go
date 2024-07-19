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

const (
	dolarAPIURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	dbPath      = "./cotacoes.db"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	// Inicialização do servidor HTTP
	http.HandleFunc("/cotacao", cotacaoHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Configurando timeout de 200ms para a requisição da API de câmbio
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	// Requisição para obter a cotação atual do dólar
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dolarAPIURL, nil)
	if err != nil {
		log.Printf("Erro ao criar requisição HTTP: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao realizar requisição HTTP para a API de câmbio: %v", err)
		http.Error(w, "Erro ao obter cotação", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erro ao ler resposta da API de câmbio: %v", err)
		http.Error(w, "Erro ao ler resposta da API", http.StatusInternalServerError)
		return
	}

	var cotacao map[string]Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		log.Printf("Erro ao decodificar resposta JSON da API de câmbio: %v", err)
		http.Error(w, "Erro ao decodificar resposta da API", http.StatusInternalServerError)
		return
	}

	// Salvando a cotação no banco de dados SQLite
	err = salvaCotacao(ctx, cotacao)
	if err != nil {
		log.Printf("Erro ao salvar cotação no banco de dados: %v", err)
		// Não retornamos um erro HTTP específico aqui porque já enviamos a resposta JSON
	}

	//log.Println("Bid: ", cotacao["USDBRL"].Bid)
	// Retornando apenas o valor "bid" para o cliente
	response := struct {
		Valor string `json:"valor"`
	}{
		Valor: cotacao["USDBRL"].Bid,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Erro ao serializar resposta JSON: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func salvaCotacao(ctx context.Context, cotacao map[string]Cotacao) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão com banco de dados: %v", err)
	}
	defer db.Close()

	// Configurando timeout de 10ms para a persistência no banco de dados
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}

	stmt, err := tx.Prepare("INSERT INTO cotacoes (data, valor) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao preparar statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, time.Now().Format(time.RFC3339), cotacao["USDBRL"].Bid)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao executar insert: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("erro ao commitar transação: %v", err)
	}

	return nil
}
