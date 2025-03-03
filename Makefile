SHELL := /bin/bash

.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go get github.com/stretchr/testify
	go get github.com/joho/godotenv

.PHONY: test
test:
	@echo "Running tests..."
	go test ./... -v