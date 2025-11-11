# coding-agent-playground

Go multi-module repository with ConnectRPC HTTP server, structured logging (slog), and OpenTelemetry OTLP exporter.

## Project Structure

This is a Go multi-module repository with the following structure:

```
.
├── go.work                    # Go workspace file
├── cmd/
│   └── server/               # Main application
│       ├── main.go           # Entry point with errgroup and graceful shutdown
│       ├── env.go            # Environment variables and command-line flags
│       └── inject.go         # Dependency injection
├── pkg/
│   ├── proto/                # Protocol Buffer definitions
│   │   ├── api/v1/          # API v1 definitions
│   │   │   ├── service.proto
│   │   │   ├── service.pb.go
│   │   │   └── serviceconnect/
│   │   │       └── service.connect.go
│   │   ├── buf.yaml
│   │   └── buf.gen.yaml
│   └── service/              # Service implementations
│       └── handler.go        # ConnectRPC service handlers
```

## Features

- **Multi-module repository**: Uses Go workspaces for managing multiple modules
- **ConnectRPC**: HTTP server using ConnectRPC protocol
- **Structured logging**: Uses Go's standard `log/slog` package
- **OpenTelemetry**: OTLP exporter for distributed tracing
- **Configuration**: Command-line flags with environment variable fallbacks
- **Graceful shutdown**: Signal handling with graceful HTTP server shutdown
- **Dependency Injection**: Centralized DI in `inject.go`
- **GitHub Actions**: CI/CD pipelines for testing, linting, and releases
- **Docker**: Containerization with Docker and Docker Compose
- **Makefile**: Convenient commands for development tasks

## Configuration

The server can be configured via command-line flags or environment variables:

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `-addr` | `SERVER_ADDR` | `:8080` | Server address |
| `-otlp-endpoint` | `OTLP_ENDPOINT` | `localhost:4317` | OTLP endpoint |
| `-log-level` | `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `-service-name` | `SERVICE_NAME` | `greet-service` | Service name for tracing |
| `-shutdown-timeout` | `SHUTDOWN_TIMEOUT` | `30` | Shutdown timeout in seconds |

## Quick Start

### Using Make

```bash
# Run the server
make run

# Build the server
make build

# Run tests
make test

# Run linters
make lint
```

### Direct Go Commands

```bash
cd cmd/server
go run . -addr=:8080 -log-level=debug
```

Or with environment variables:

```bash
SERVER_ADDR=:8080 LOG_LEVEL=debug go run ./cmd/server
```

### Using Docker Compose

```bash
# Start server with Jaeger
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

The server will be available at `http://localhost:8080` and Jaeger UI at `http://localhost:16686`.

## API

The server exposes a simple Greet API:

### Greet

**Endpoint**: `/api.v1.GreetService/Greet`

**Request**:
```json
{
  "name": "World"
}
```

**Response**:
```json
{
  "message": "Hello, World!"
}
```

## Development

### Prerequisites

- Go 1.23 or later
- Docker and Docker Compose (optional, for containerized development)
- golangci-lint (optional, for linting)
- buf (optional, for Protocol Buffer code generation)

### Building

```bash
# Using Makefile
make build

# Or directly with Go
cd cmd/server
go build -o server .
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Linting

```bash
# Run linters for all modules
make lint
```

### Generating Protocol Buffers

```bash
# Requires buf to be installed
make proto-gen
```

### Testing the API

#### With curl

```bash
curl -X POST http://localhost:8080/api.v1.GreetService/Greet \
  -H "Content-Type: application/json" \
  -d '{"name": "World"}'
```

#### With grpcurl (for Connect protocol)

```bash
grpcurl -plaintext -d '{"name": "World"}' \
  localhost:8080 api.v1.GreetService/Greet
```

## CI/CD

### GitHub Actions

This project includes several GitHub Actions workflows:

- **CI** (`.github/workflows/ci.yml`): Runs on every push and pull request
  - Linting with golangci-lint
  - Testing on Go 1.23 and 1.24
  - Build verification
  - Code coverage reporting

- **Release** (`.github/workflows/release.yml`): Triggered on version tags
  - Builds binaries for multiple platforms (Linux, macOS, Windows)
  - Creates GitHub releases with built artifacts

- **Dependabot** (`.github/dependabot.yml`): Automatically updates dependencies
  - Go module dependencies
  - GitHub Actions versions

### Creating a Release

```bash
# Tag a new version
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions will automatically create a release with binaries
```

## Docker

### Building the Docker Image

```bash
make docker-build
```

### Running with Docker Compose

The `docker-compose.yml` includes:
- The server application
- Jaeger for distributed tracing

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f server

# Stop all services
docker-compose down
```

## Observability

### Distributed Tracing

The application exports traces to an OTLP endpoint. When using Docker Compose, traces are sent to Jaeger.

Access Jaeger UI at: `http://localhost:16686`

### Logging

Structured logging using `log/slog` with JSON output:

```json
{
  "time": "2024-01-01T12:00:00Z",
  "level": "INFO",
  "msg": "received greet request",
  "name": "World"
}
```