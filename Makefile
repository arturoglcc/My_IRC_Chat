# Makefile

# Define directories
GO_SRC_DIR := ./src/server
JS_SRC_DIR := ./src/client
GO_BUILD_DIR := ./build
JS_BUILD_DIR := ./build/client
GO_EXEC := my_irc_server
SERVER_TEST_DIR := ./tests/server_tests
CLIENT_TEST_FILE := ./tests/client_test.js

# Default target
all: test-server test-client build-server build-client

# Run Go unit tests for the server
test-server:
	@echo "Running Go server unit tests..."
	go test -v $(SERVER_TEST_DIR)

# Run JavaScript unit tests for the client (using Jest)
test-client:
	@echo "Running JavaScript client unit tests..."
	cd $(JS_SRC_DIR) && npm test ../$(CLIENT_TEST_FILE)

# Build Go server executable
build-server:
	@echo "Building Go server executable..."
	mkdir -p $(GO_BUILD_DIR)
	go build -o $(GO_BUILD_DIR)/$(GO_EXEC) $(GO_SRC_DIR)/main.go

# Build JavaScript client (e.g., using Webpack)
build-client:
	@echo "Building JavaScript client..."
	mkdir -p $(JS_BUILD_DIR)
	cd $(JS_SRC_DIR) && npm run build -- --output-path=$(JS_BUILD_DIR)

# Clean all build files and outputs
clean:
	@echo "Cleaning up build files..."
	rm -rf $(GO_BUILD_DIR)/*
	rm -rf $(JS_BUILD_DIR)/*

# Phony targets to prevent conflicts with files
.PHONY: all test-server test-client build-server build-client clean

