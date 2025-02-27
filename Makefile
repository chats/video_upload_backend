# Variables
APP_NAME=video-upload-backend
BUILD_DIR=./build
MAIN_FILE=./cmd/api/main.go
MIGRATE_FILE=./cmd/migrate/main.go

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOMOD=$(GOCMD) mod
GOGET=$(GOCMD) get

# Determine the operating system
ifeq ($(OS),Windows_NT)
    EXT=.exe
else
    EXT=
endif

.PHONY: all build clean run test lint fmt migrate docker-build docker-run dev help

all: clean build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)$(EXT) $(MAIN_FILE)
	@echo "Build complete"

# Clean build files
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	$(GORUN) $(MAIN_FILE)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run linter
lint:
	@echo "Running linter..."
	$(GOVET) ./...

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -s -w .

# Run migrations
migrate:
	@echo "Running database migrations..."
	$(GORUN) $(MIGRATE_FILE)

# Update dependencies
deps:
	@echo "Updating dependencies..."
	$(GOMOD) tidy

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME) .

# Run with Docker
docker-run:
	@echo "Running with Docker..."
	docker compose up

# Run with Docker in detached mode
docker-up:
	@echo "Starting in detached mode..."
	docker compose up -d 

# Stop Docker containers
docker-down:
	@echo "Stopping containers..."
	docker compose down

# Setup environment
setup:
	@echo "Running setup script..."
	./setup.sh

# Development mode
dev:
	@echo "Starting in development mode with live reload..."
	air

# Display help information
help:
	@echo "Available targets:"
	@echo "  all          - Clean and build the application"
	@echo "  build        - Build the application"
	@echo "  clean        - Clean build files"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  migrate      - Run database migrations"
	@echo "  deps         - Update dependencies"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker"
	@echo "  docker-up    - Run with Docker in detached mode"
	@echo "  docker-down  - Stop Docker containers"
	@echo "  setup        - Setup environment"
	@echo "  dev          - Run in development mode with live reload"
	@echo "  help         - Display this help message"