package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type BrApiResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCepApiResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {

	brasilApiCh := make(chan *BrApiResponse)
	viaCepApiCh := make(chan *ViaCepApiResponse)
	go func() {
		brasilApiCh <- getResponseBrasilApi()
	}()
	go func() {
		viaCepApiCh <- getResponseViaCepApi()
	}()

	select {
	case brasilApi := <-brasilApiCh:
		fmt.Printf("Brasil API: Endereço: CEP %s, Estado %s, Cidade %s, Bairro %s, Rua %s\n",
			brasilApi.Cep, brasilApi.State, brasilApi.City, brasilApi.Neighborhood, brasilApi.Street)
	case viaCepApi := <-viaCepApiCh:
		fmt.Printf("VIA CEP: Endereço: CEP %s, Estado %s, Cidade %s, Bairro %s, Rua %s\n",
			viaCepApi.Cep, viaCepApi.Uf, viaCepApi.Localidade, viaCepApi.Bairro, viaCepApi.Logradouro)
	case <-time.After(1 * time.Second):
		println("Timeout")

	}

}

func fetchApiResponse(url string, responseStruct interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar a requisição: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar a requisição: %w", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(responseStruct)
	if err != nil {
		return fmt.Errorf("erro ao decodificar a resposta: %w", err)
	}

	return nil
}

func getResponseBrasilApi() *BrApiResponse {
	var responseApi BrApiResponse
	err := fetchApiResponse("https://brasilapi.com.br/api/cep/v1/79021210", &responseApi)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao obter resposta da Brasil API: %v\n", err)
		panic(err)
	}
	return &responseApi
}

func getResponseViaCepApi() *ViaCepApiResponse {
	var responseApi ViaCepApiResponse
	err := fetchApiResponse("http://viacep.com.br/ws/79021210/json", &responseApi)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao obter resposta da ViaCEP API: %v\n", err)
		panic(err)
	}
	return &responseApi
}
