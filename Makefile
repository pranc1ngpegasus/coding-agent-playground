.PHONY: help build run test lint clean docker-build docker-up docker-down proto-gen

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the server binary
	cd cmd/server && go build -o ../../bin/server .

run: ## Run the server
	cd cmd/server && go run . -log-level=debug

test: ## Run tests
	cd cmd/server && go test -v -race ./...
	cd pkg/service && go test -v -race ./...

test-coverage: ## Run tests with coverage
	cd cmd/server && go test -v -race -coverprofile=coverage.out ./...
	cd pkg/service && go test -v -race -coverprofile=coverage.out ./...

lint: ## Run linters
	cd cmd/server && golangci-lint run
	cd pkg/service && golangci-lint run
	cd pkg/proto && golangci-lint run

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf dist/
	rm -f cmd/server/coverage.out
	rm -f pkg/service/coverage.out

docker-build: ## Build Docker image
	docker build -t greet-service:latest .

docker-up: ## Start services with Docker Compose
	docker-compose up -d

docker-down: ## Stop services with Docker Compose
	docker-compose down

docker-logs: ## View Docker Compose logs
	docker-compose logs -f

proto-gen: ## Generate protobuf code (requires buf)
	cd pkg/proto && buf generate

tidy: ## Run go mod tidy for all modules
	cd cmd/server && go mod tidy
	cd pkg/service && go mod tidy
	cd pkg/proto && go mod tidy
