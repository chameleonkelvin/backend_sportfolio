.PHONY: run build clean test help

# Run the application in development mode
run:
	go run main.go

# Build the application
build:
	go build -o scoring_app.exe

# Build for production (with optimizations)
build-prod:
	go build -ldflags="-s -w" -o scoring_app.exe

# Clean build artifacts
clean:
	del /f scoring_app.exe

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod tidy
	go mod download

# Update dependencies
deps-update:
	go get -u ./...
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Create database
db-create:
	mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS scoring_app_db;"

# Help command
help:
	@echo Available targets:
	@echo   run         - Run application in development mode
	@echo   build       - Build application
	@echo   build-prod  - Build application for production
	@echo   clean       - Clean build artifacts
	@echo   test        - Run tests
	@echo   deps        - Install dependencies
	@echo   deps-update - Update dependencies
	@echo   fmt         - Format code
	@echo   lint        - Run linter
	@echo   db-create   - Create database
	@echo   help        - Show this help message
