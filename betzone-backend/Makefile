.PHONY: help build run dev test clean deps

help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make dev            - Run in development mode"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make deps           - Download dependencies"

build:
	@echo "Building application..."
	@go build -o bin/betzone-backend .

run: build
	@echo "Running application..."
	@./bin/betzone-backend

dev:
	@echo "Running in development mode..."
	@go run main.go

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@go clean

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

docker-build:
	@echo "Building Docker image..."
	@docker build -t betzone-backend:latest .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env betzone-backend:latest

lint:
	@echo "Running linter..."
	@go fmt ./...

