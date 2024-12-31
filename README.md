# Sistema de Clima por CEP com OpenTelemetry e Zipkin

## Visão Geral

Este projeto contém dois serviços escritos em Go:

- **Serviço A**: Recebe um CEP via POST, valida se possui 8 dígitos e repassa ao Serviço B via HTTP.
- **Serviço B**: Recebe o CEP (já validado), consulta a API ViaCEP para extrair a cidade e chama a API de clima (WeatherAPI) para obter a temperatura em Celsius. Em seguida, converte para Fahrenheit e Kelvin, retornando uma resposta JSON.

Além disso, foi implementado **OpenTelemetry (OTEL)** para **Distributed Tracing** com **Zipkin**. Assim é possível visualizar spans e métricas de cada chamada.

## Pré-requisitos

- **Go** (>= 1.20) instalado localmente, caso queira rodar sem Docker.
- **Docker** e **Docker Compose** instalados, caso queira utilizar os containers.

## Rodando com Docker Compose

1. **Clonar o repositório e acessar a pasta principal.**
2. **Editar o arquivo `docker/docker-compose.yaml`:**

   - Inserir sua chave de API do WeatherAPI em `WEATHER_API_KEY`.
3. **Executar:**

   ```bash
   cd docker
   docker-compose up --build
   ```
4. **Aguardar a inicialização dos containers:**

   - `zipkin` na porta `9411`
   - `servicoA` na porta `8080`
   - `servicoB` na porta `9090`
5. **Testar com `curl` ou outro cliente:**

   ```bash
   curl -X POST -H "Content-Type: application/json" \
   -d '{"cep": "01001000"}' \
   http://localhost:8080/cep
   ```
6. **Acessar o Zipkin:** [http://localhost:9411](http://localhost:9411)

## Estrutura de Resposta

### Sucesso (200):

```bash
curl -s -w '\nHTTP Status: %{http_code}\n' -X POST -H "Content-Type: application/json" -d '{"cep": "01001000"}' http://localhost:8080/cep
```

```bash
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}

HTTP Status: 200
```

### CEP inválido (422):

```bash
curl -s -w '\nHTTP Status: %{http_code}\n' -X POST -H "Content-Type: application/json" -d '{"cep": "asdasd"}' http://localhost:8080/cep
```

```bash
{ "message": "invalid zipcode" }

HTTP Status: 422
```

### CEP não encontrado (404):

```bash
curl -s -w '\nHTTP Status: %{http_code}\n' -X POST -H "Content-Type: application/json" -d '{"cep": "01010101"}' http://localhost:8080/cep
```

```bash
{ "message": "can not find zipcode" }

HTTP Status: 404
```

---

## Conclusão

Esta estrutura e código-fonte ilustram:

1. **Serviço A**:

   - Recebe input `{ "cep": "29902555" }`
   - Valida e encaminha para o Serviço B.
   - Retorna eventuais erros (422) se o CEP for inválido.
2. **Serviço B**:

   - Recebe CEP válido.
   - Busca cidade via `ViaCEP`.
   - Busca clima via `WeatherAPI`.
   - Converte Celsius -> Fahrenheit -> Kelvin.
   - Retorna JSON com `{ city, temp_C, temp_F, temp_K }`.
   - Retorna 422 se CEP estiver no formato correto, mas for inválido na lógica do sistema.
   - Retorna 404 se o CEP não for encontrado pelo ViaCEP.
3. **Tracing distribuído** entre A e B, configurado para exportar spans ao Zipkin.
4. **Docker** e **docker-compose** para subir `servicoA`, `servicoB` e `zipkin` juntos.

## Referências

- [OpenTelemetry Go - Getting Started](https://opentelemetry.io/docs/instrumentation/go/getting-started/)
- [Zipkin.io](https://zipkin.io/)
- [ViaCEP](https://viacep.com.br/)
- [WeatherAPI](https://www.weatherapi.com/)
