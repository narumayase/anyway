# anyway

Esta es una API en Go que recibe mensajes a través de HTTP POST y los reenvía a un tema de Kafka. Está diseñada para ser un productor de mensajes ligero y eficiente para Kafka.

## Características

*   **API HTTP:** Expone un endpoint RESTful para recibir mensajes.
*   **Integración con Kafka:** Produce mensajes de forma transparente a un tema de Kafka configurable.

### Prerrequisitos

*   Go: Versión 1.18 o superior.
*   Docker. 
*   Kafka.

### Configuración

La aplicación se puede configurar utilizando las siguientes variables de entorno:

*   `PORT`: El puerto en el que el servidor HTTP escuchará. (Por defecto: `8080`)
*   `KAFKA_BROKER`: La dirección del broker de Kafka (ej. `localhost:9092`). (Por defecto: `localhost:9092`)
*   `KAFKA_TOPIC`: El tema de Kafka al que se producirán los mensajes. (Por defecto: `anyway-topic`)
*   `LOG_LEVEL`: El nivel de registro (ej. `debug`, `info`, `warn`, `error`). (Por defecto: `info`)

Puedes crear un archivo `.env` en la raíz del proyecto para establecer estas variables, por ejemplo:

```
PORT=8080
KAFKA_BROKER=localhost:9092
KAFKA_TOPIC=mis-mensajes
LOG_LEVEL=debug
```

### Running the Application

1. Install dependencies:

```bash
go mod tidy
```

2. Configure environment variables:

```bash
cp env.example .env
# Edit .env with the values described below.
```

3. Run the application:

```bash
go run main.go
```

## Endpoints de la API

### `POST /api/v1/send`

Recibe un mensaje JSON y lo produce al tema de Kafka configurado.

**Ejemplo de Cuerpo de Solicitud:**

```json
{
    "key": "clave-mensaje-123",
    "headers": {
        "contentType": "application/json",
        "source": "mi-app"
    },
    "content": "SG9sYSBNdW5kbyBLYWZrYSE="
}
```
*   `key` (string, opcional): Una clave para el mensaje de Kafka.
*   `headers` (objeto, opcional): Un mapa de pares clave-valor de tipo string para los encabezados del mensaje de Kafka.
*   `content` (string, requerido): El mensaje, se espera que sea una cadena codificada en base64.

**Respuesta:**

*   `200 OK`: Mensaje enviado exitosamente a Kafka.
*   `400 Bad Request`: Formato de solicitud inválido.
*   `500 Internal Server Error`: Error al procesar o enviar el mensaje a Kafka.

### `GET /health`

Proporciona una simple verificación de estado para la API.

**Ejemplo de Respuesta:**

```json
{
    "status": "OK",
    "message": "La API anyway está en ejecución"
}
```

## Ejemplo de Uso

Para enviar un mensaje usando `curl`:

```bash
curl -X POST http://localhost:8080/api/v1/send \
     -H "Content-Type: application/json" \
     -d '{
           "key": "mi-clave-unica",
           "headers": {
             "source": "ejemplo-cli"
           },
           "content": "SG9sYSBkZXNkZSBjdXJsIQ=="
         }'
```

## Ejecución de las pruebas

Para ejecutar las pruebas unitarias:

```bash
go test ./...
```

### Test Coverage

To check test coverage (excluding mocks):

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage report in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# View coverage excluding mocks
go test -coverprofile=coverage.out ./... && \
go tool cover -func=coverage.out | grep -v "mocks"
```

## 🎗️ Architecture

This project follows Clean Architecture principles:

- **Domain**: Entities, repository interfaces, and use cases
- **Application**: Implementation of use cases
- **Infrastructure**: Kafka repository implementations
- **Interfaces**: HTTP controllers and routers

## 📁 Project Structure

```
anyompt/
├── cmd/                  # Application entry points
│   └── server/           # Main server
├── config/               # Configuration
├── internal/             # Project-specific code
│   ├── infrastructure/   # Repository implementations
│   └── interfaces/       # HTTP controllers
│       ├── http/         # Handler controller
│       └── middleware/   # Middlewares
│   ├── domain/           # Domain entities and interfaces
│   └── application/      # Use cases
├── main.go               # Main entry point
├── go.mod                # Go dependencies
├── README_ES.md          # README in spanish
└── README.md             # This file
```