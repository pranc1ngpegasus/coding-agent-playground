# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /workspace

# Copy workspace file
COPY go.work go.work

# Copy all modules
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build the application
WORKDIR /workspace/cmd/server
RUN go build -o /server .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /server /app/server

EXPOSE 8080

ENTRYPOINT ["/app/server"]
