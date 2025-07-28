# Makefile for the Go project

.PHONY: run build test

# Go parameters
GOBASE := $(shell pwd)
GOPACKAGES := $(shell go list ./... | grep -v /vendor/)
GOFILES := $(shell find . -name "*.go" -not -path "./vendor/*")

# Variables
APP_NAME=aiload
CMD_PATH=./cmd

# Default target
all: build

# Run the application using air for live reloading
run:
	@echo "Running the application with air..."
	@air

# Build the application
build:
	@echo "Building the application..."
	@go build -o $(APP_NAME) $(CMD_PATH)/main.go

# Run tests
test:
	@echo "Running tests..."
	@CGO_ENABLED=1 go test -v -cover $(GOPACKAGES)

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(APP_NAME)