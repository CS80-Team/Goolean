SHELL := /bin/bash
APP_NAME := goolean
BIN_DIR := bin
ENTRY_PATH := .

.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy

.PHONY: test
test: deps
	@echo "Running tests..."
	go test ./... -v

.PHONY: run
run: fmt deps
	@echo "Running the application..."
	go run $(ENTRY_PATH)

.PHONY: build
build: clean deps fmt test
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
