package main

import (
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type CEPRequest struct {
	CEP string `json:"cep"`
}

type ResponseBody struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func HandleTemperatura(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("Servico-B-Handler")
	ctx, span := tracer.Start(r.Context(), "HandleTemperatura")
	defer span.End()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Decodificar CEP
	var req CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid zipcode"})
		return
	}

	// Se não tiver 8 dígitos
	if len(req.CEP) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid zipcode"})
		return
	}

	// Adiciona atributo no span
	span.SetAttributes(attribute.String("cep", req.CEP))

	// Busca cidade via viaCep (ou fallback)
	city, err := getCityFromViaCep(ctx, req.CEP)
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "can not find zipcode"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	// Com a cidade, chamamos a Weather API e obtemos a temperatura em Celsius
	tempC, err := getTemperatureCelsius(ctx, city)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	// Converte para F e K
	tempF := (tempC * 1.8) + 32
	tempK := tempC + 273

	// Monta resposta
	response := ResponseBody{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
