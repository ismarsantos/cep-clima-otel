package main

import (
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

)

func main() {
	zipkinEndpoint := os.Getenv("ZIPKIN_ENDPOINT")

	if zipkinEndpoint == "" {
		zipkinEndpoint = "http://localhost:9411/api/v2/spans"
	}
	exporter, err := zipkin.New(
		zipkinEndpoint,
	)
	if err != nil {
		log.Fatalf("erro ao criar exporter do Zipkin: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
    sdktrace.WithBatcher(exporter),
    sdktrace.WithResource(resource.NewWithAttributes(
        semconv.SchemaURL,
        attribute.String("service.name", "Servico-A"), // Altere para esta linha
    )),
	)
	defer func() {
		_ = tp.Shutdown(nil)
	}()
	otel.SetTracerProvider(tp)

	// Registrar rotas
	http.HandleFunc("/cep", HandleCEP)

	// Subir servidor
	log.Println("Serviço A rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
