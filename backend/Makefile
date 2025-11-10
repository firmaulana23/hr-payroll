SHELL := /bin/bash
include .env

.PHONY: all build run fmt docs test clean db-create

BINARY := hr-payroll
BUILD_DIR := ./bin

all: build

build: ## Build the backend binary (linux/amd64, static)
	@echo "==> Building $(BINARY)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY) ./cmd

run: ## Run the app locally (uses local Go toolchain)
	@echo "==> Running app"
	@mkdir -p .tmp
	go run ./cmd

fmt: ## Run go fmt on project
	@echo "==> Formatting"
	go fmt ./...

docs: ## Generate swagger docs (requires swag installed)
	@echo "==> Generating Swagger docs"
	swag init -g ./cmd/main.go -o ./docs

test: ## Run unit tests
	@echo "==> Running tests"
	go test ./...

db-create: ## Create PostgreSQL database if it doesn't exist
	@echo "==> Ensuring database $(DB_NAME) exists"
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -p $(DB_PORT) -tc "SELECT 1 FROM pg_database WHERE datname='$(DB_NAME)'" \
		| grep -q 1 && echo "==> Database $(DB_NAME) already exists" || \
		(PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -p $(DB_PORT) -c "CREATE DATABASE $(DB_NAME);" && echo "==> Created database $(DB_NAME)")

clean: ## Remove built artifacts and generated docs
	@echo "==> Cleaning"
	rm -rf $(BUILD_DIR) ./docs
