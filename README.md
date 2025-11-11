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

## Configuration

The server can be configured via command-line flags or environment variables:

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `-addr` | `SERVER_ADDR` | `:8080` | Server address |
| `-otlp-endpoint` | `OTLP_ENDPOINT` | `localhost:4317` | OTLP endpoint |
| `-log-level` | `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `-service-name` | `SERVICE_NAME` | `greet-service` | Service name for tracing |
| `-shutdown-timeout` | `SHUTDOWN_TIMEOUT` | `30` | Shutdown timeout in seconds |

## Running the Server

```bash
cd cmd/server
go run . -addr=:8080 -log-level=debug
```

Or with environment variables:

```bash
SERVER_ADDR=:8080 LOG_LEVEL=debug go run ./cmd/server
```

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
- (Optional) buf for Protocol Buffer code generation

### Building

```bash
cd cmd/server
go build -o server .
```

### Testing with curl

```bash
curl -X POST http://localhost:8080/api.v1.GreetService/Greet \
  -H "Content-Type: application/json" \
  -d '{"name": "World"}'
```