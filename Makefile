SHELL := /bin/bash

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
	go run main.go

.PHONY: build
build: deps test
	@echo "Building the application..."
	go build -o bin/boolean-ir-system main.go

.PHONY: run
run: deps test
	@echo "Running the application..."
	go run main.go
