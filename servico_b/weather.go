package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Exemplo de resposta simplificada do WeatherAPI
type WeatherApiResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func getTemperatureCelsius(ctx context.Context, city string) (float64, error) {
	tracer := otel.Tracer("Servico-B-Handlers")
	ctx, span := tracer.Start(ctx, "getTemperatureCelsius")
	defer span.End()

	span.SetAttributes(attribute.String("weather.city", city))

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
			return 0, errors.New("missing WEATHER_API_KEY")
	}

	// Codificar o nome da cidade para uso na URL
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, encodedCity)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
			return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
			return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return 0, fmt.Errorf("weather request failed: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
			return 0, err
	}

	var weatherResp WeatherApiResponse
	if err := json.Unmarshal(bodyBytes, &weatherResp); err != nil {
			return 0, err
	}

	return weatherResp.Current.TempC, nil
}
