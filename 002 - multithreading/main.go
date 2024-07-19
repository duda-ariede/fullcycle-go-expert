package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Address struct {
	Source  string `json:"source"`
	CEP     string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro  string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF      string `json:"uf"`
}

type BrasilAPIResponse struct {
    CEP      string `json:"cep"`
    State    string `json:"state"`
    City     string `json:"city"`
    Neighborhood string `json:"neighborhood"`
    Street   string `json:"street"`
}

type ViaCEPResponse struct {
    CEP          string `json:"cep"`
    Logradouro  string `json:"logradouro"`
    Bairro      string `json:"bairro"`
    Localidade string `json:"localidade"`
    UF         string `json:"uf"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Por favor, informe um CEP válido como argumento")
		return
	}

	cep := os.Args[1]

	if len(cep) < 8 {
		fmt.Println("CEP inválido. O CEP deve ter pelo menos 8 caracteres")
		return
	}

	ch := make(chan Address, 2)

	go fetchFromBrasilAPI(cep, ch)
	go fetchFromViaCEP(cep, ch)

	select {
	case result := <-ch:
		printResult(result)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout: Não houve resposta dentro do tempo de 1 segundo")
	}
}

func fetchFromBrasilAPI(cep string, ch chan<- Address) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Erro ao acessar a API BrasilAPI:", err)
		ch <- Address{}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro lendo a resposta da API BrasilAPI:", err)
		ch <- Address{}
		return
	}

	var response BrasilAPIResponse
  err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Erro no parsing do JSON da API BrasilAPI:", err)
		ch <- Address{}
		return
	}

	address := Address{
    Source:  "BrasilAPI",
		CEP:     response.CEP,
		Logradouro: response.Street,
		Bairro:  response.Neighborhood,
		Localidade: response.City,
		UF:      response.State,
  }
  ch <- address
}

func fetchFromViaCEP(cep string, ch chan<- Address) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Erro ao acessar a API ViaCEP:", err)
		ch <- Address{}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro lendo a resposta da API ViaCEP:", err)
		ch <- Address{}
		return
	}

	var response ViaCEPResponse
  err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Erro no parsing do JSON da API ViaCEP:", err)
		ch <- Address{}
		return
	}

	address := Address{
		Source:  "ViaCEP",
		CEP:     response.CEP,
		Logradouro: response.Logradouro,
		Bairro:  response.Bairro,
		Localidade: response.Localidade,
		UF:      response.UF,
	}
	ch <- address
}

func printResult(address Address) {
	if address.CEP != "" {
		fmt.Printf("API Fonte: %s\n", address.Source)
		fmt.Printf("CEP: %s\n", address.CEP)
		fmt.Printf("Logradouro: %s\n", address.Logradouro)
		fmt.Printf("Bairro: %s\n", address.Bairro)
		fmt.Printf("Localidade: %s\n", address.Localidade)
		fmt.Printf("UF: %s\n", address.UF)
	} else {
		fmt.Println("Nenhum endereço válido encontrado")
	}
}
