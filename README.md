# anyway

This is a Go API that receives messages via HTTP POST and forwards them to a Kafka topic. It's designed to be a lightweight and efficient message producer for Kafka.

## Features

*   **HTTP API:** Exposes a RESTful endpoint to receive messages.
*   **Kafka Integration:** Seamlessly produces messages to a configurable Kafka topic.

### Prerequisites

*   **Go:** Version 1.18 or higher.
*   **Docker (Optional).
*   **Kafka.

### Configuration

The application can be configured using the following environment variables:

*   `PORT`: The port on which the HTTP server will listen. (Default: `8080`)
*   `KAFKA_BROKER`: The address of the Kafka broker (e.g., `localhost:9092`). (Default: `localhost:9092`)
*   `KAFKA_TOPIC`: The Kafka topic to which messages will be produced. (Default: `anyway-topic`)
*   `LOG_LEVEL`: The logging level (e.g., `debug`, `info`, `warn`, `error`). (Default: `info`)

You can create an `.env` file in the project root to set these variables, for example:

```
PORT=8080
KAFKA_BROKER=localhost:9092
KAFKA_TOPIC=my-messages
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

## API Endpoints

### `POST /api/v1/send`

Receives a JSON message and produces it to the configured Kafka topic.

**Request Body Example:**

```json
{
    "key": "message-key-123",
    "headers": {
        "contentType": "application/json",
        "source": "my-app"
    },
    "content": "SGVsbG8gS2Fma2Egd29ybGQh"
}
```
*   `key` (string, optional): A key for the Kafka message.
*   `headers` (object, optional): A map of string key-value pairs for Kafka message headers.
*   `content` (string, required): The message payload, expected to be a base64 encoded string.

**Response:**

*   `200 OK`: Message successfully sent to Kafka.
*   `400 Bad Request`: Invalid request format.
*   `500 Internal Server Error`: Error processing or sending the message to Kafka.

### `GET /health`

Provides a simple health check for the API.

**Response Example:**

```json
{
    "status": "OK",
    "message": "anyway API is running"
}
```

## Running the tests

To run the unit tests:

```bash
go test ./...
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
