FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./servico_b/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/servico_b main.go handlers.go via_cep.go weather.go
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/servico_b ./

CMD ["./servico_b"]
