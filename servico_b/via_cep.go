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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
			return "", err
	}

	// Substitua http.DefaultClient por getInsecureHttpClient()
	resp, err := getInsecureHttpClient().Do(req)
	if err != nil {
			return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
			return "", errors.New("viaCep request failed")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
			return "", err
	}

	var viaCepResp ViaCepResponse
	if err := json.Unmarshal(bodyBytes, &viaCepResp); err != nil {
			return "", err
	}

	if viaCepResp.Erro {
			// caso CEP não encontrado
			return "", errors.New("not_found")
	}

	return viaCepResp.Localidade, nil
}


// getInsecureHttpClient cria um cliente HTTP que ignora validação TLS
func getInsecureHttpClient() *http.Client {
	transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: transport}
}