package main

import (
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	// Inicia o Tracer para Zipkin
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
			semconv.ServiceNameKey.String("Servico-B"),
		)),
	)
	defer func() {
		_ = tp.Shutdown(nil)
	}()
	otel.SetTracerProvider(tp)

	// Registrar rota
	http.HandleFunc("/temperatura", HandleTemperatura)

	// Subir servidor
	log.Println("Serviço B rodando na porta 9090...")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
