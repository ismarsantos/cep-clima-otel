package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"crypto/tls"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// ViaCepResponse exemplo simples do retorno JSON do ViaCep
type ViaCepResponse struct {
	Localidade string `json:"localidade"` // Nome da cidade
	Erro       bool   `json:"erro"`       // true se não encontrou
}

// getCityFromViaCep faz a consulta no ViaCep e retorna a localidade (cidade)
func getCityFromViaCep(ctx context.Context, cep string) (string, error) {
	tracer := otel.Tracer("Servico-B-Handlers")
	ctx, span := tracer.Start(ctx, "getCityFromViaCep")
	defer span.End()

	span.SetAttributes(attribute.String("viaCep.cep", cep))

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	fmt.Println("Calling URL:", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
			fmt.Println("Error creating request:", err)
			return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
			fmt.Println("Error making HTTP request:", err)
			return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Printf("HTTP response status: %d, body: %s\n", resp.StatusCode, string(bodyBytes))
			return "", errors.New("viaCep request failed")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
			fmt.Println("Error reading response body:", err)
			return "", err
	}

	var viaCepResp ViaCepResponse
	if err := json.Unmarshal(bodyBytes, &viaCepResp); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return "", err
	}

	if viaCepResp.Erro {
			fmt.Println("CEP not found")
			return "", errors.New("not_found")
	}

	fmt.Println("City found:", viaCepResp.Localidade)
	return viaCepResp.Localidade, nil
}

// getInsecureHttpClient cria um cliente HTTP que ignora validação TLS
func getInsecureHttpClient() *http.Client {
	return &http.Client{
			Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
	}
}
