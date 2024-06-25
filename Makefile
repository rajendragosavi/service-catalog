# Define project name
PROJECT_NAME := service-catalog

# Define directories
BIN_DIR := bin
SRC_DIRS := ./cmd ./internal ./pkg 

# Define output binary
BINARY := $(BIN_DIR)/$(PROJECT_NAME)

# Define Go command
GO := go

.PHONY: all build test clean

# Default target: build the project
all: build

# Build the project
build: $(BINARY)

# Build binary
$(BINARY): $(SRC_DIRS)
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BINARY) ./cmd/app

run:
	./bin/service-catalog

# Run tests
test:
	$(GO) test -v ./...

setup:
	echo "DATABASE_HOST=localhost" > .env
	echo "DATABASE_PORT=5432" >> .env
	echo "DATABASE_USER=postgres" >> .env
	echo "DATABASE_PASSWORD=postgres" >> .env
	echo "DATABASE_NAME=service" >> .env
	echo "SERVER_PORT=80" >> .env

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR)

docker-build :
	docker build -t service-catalog --no-cache .

deploy :
	docker run -d -p 80:80 service-catalog:latest