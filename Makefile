SHELL := /bin/bash
APP_NAME := birs
BIN_DIR := bin
ENTRY_PATH := .

.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go get github.com/stretchr/testify
	go get github.com/joho/godotenv

.PHONY: test
test: deps
	@echo "Running tests..."
	go test ./... -v

.PHONY: run
run: deps
	@echo "Running the application..."
	go run $(ENTRY_PATH)

.PHONY: build
build: clean deps fmt
	@echo "Building the application..."
	go build -o $(BIN_DIR)/$(APP_NAME) $(ENTRY_PATH)

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	gofmt -w .
