# Makefile for Course Enrollment Service

# Variables
APP_NAME=course-enrollment-service
BINARY_NAME=server
DOCKER_IMAGE=sonic-labs/course-enrollment-service
DOCKER_TAG=latest

# Go related variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build the application
.PHONY: build
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) cmd/server/main.go

# Run the application
.PHONY: run
run:
	$(GOBUILD) -o bin/$(BINARY_NAME) cmd/server/main.go && ./bin/$(BINARY_NAME)

# Run with live reload (requires air: go install github.com/cosmtrek/air@latest)
.PHONY: dev
dev:
	air

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Generate Swagger documentation
.PHONY: swagger
swagger:
	swag init -g cmd/server/main.go -o docs

# Format code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	golangci-lint run

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run with Docker Compose
.PHONY: docker-up
docker-up:
	docker-compose up --build

# Stop Docker Compose
.PHONY: docker-down
docker-down:
	docker-compose down

# Database migration (requires migrate tool)
.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database "postgresql://postgres:password@localhost:5432/course_enrollment?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database "postgresql://postgres:password@localhost:5432/course_enrollment?sslmode=disable" down

# Create new migration
.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Install development tools
.PHONY: install-tools
install-tools:
	$(GOGET) github.com/swaggo/swag/cmd/swag@latest
	$(GOGET) github.com/cosmtrek/air@latest

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run with live reload (requires air)"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  swagger       - Generate Swagger documentation"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code (requires golangci-lint)"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-up     - Run with Docker Compose"
	@echo "  docker-down   - Stop Docker Compose"
	@echo "  install-tools - Install development tools"
	@echo "  help          - Show this help message"
