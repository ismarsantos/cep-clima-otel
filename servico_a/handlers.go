package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type CEPRequest struct {
	CEP string `json:"cep"`
}

func HandleCEP(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("Servico-A-Handler")
	ctx, span := tracer.Start(r.Context(), "HandleCEP")
	defer span.End()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid zipcode"})
		return
	}

	if len(req.CEP) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid zipcode"})
		return
	}

	span.SetAttributes(attribute.String("input.cep", req.CEP))

	// Chamar Serviço B
	bodyBytes, statusCode, err := callServiceB(ctx, req.CEP)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	// Encaminhar resposta do serviço B ao cliente final
	w.WriteHeader(statusCode)
	w.Write(bodyBytes)
}

func callServiceB(ctx context.Context, cep string) ([]byte, int, error) {
	svcBURL := os.Getenv("SERVICO_B_URL")
	if svcBURL == "" {
		svcBURL = "http://localhost:9090/temperatura"
	}

	// Monta o json para enviar ao Serviço B
	payload, err := json.Marshal(map[string]string{"cep": cep})
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, svcBURL, nil)
	if err != nil {
		return nil, 0, err
	}

	// Usar Body para enviar o CEP em JSON
	req.Body = io.NopCloser((io.Reader)(jsonBodyReader{payload}))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return respBytes, resp.StatusCode, nil
}

type jsonBodyReader struct {
	data []byte
}

func (jbr jsonBodyReader) Read(p []byte) (int, error) {
	n := copy(p, jbr.data)
	jbr.data = jbr.data[n:]
	if len(jbr.data) == 0 {
		return n, io.EOF
	}
	return n, nil
}
