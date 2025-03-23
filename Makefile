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

.PHONY: proto
proto:
	@echo "Compiling proto files..."
	rm -rf internal/transport/*.pb.go
	protoc \
	--go_out=. \
	--go_opt=module=github.com/CS80-Team/Goolean \
	--go-grpc_out=. \
	--go-grpc_opt=module=github.com/CS80-Team/Goolean \
	api/document.proto \
	api/query.proto \
	api/load.proto


.PHONY: build
build: proto clean deps fmt test
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

.PHONY: run_service
run_service: build
	@echo "Starting the service..."
	$(BIN_DIR)/$(APP_NAME) -service $(IP):$(PORT)