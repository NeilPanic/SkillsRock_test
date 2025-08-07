include .env
export

MIGRATION_DIR    := ./migrations
LOCAL_BIN        := $(CURDIR)/bin
GOOSE            := $(LOCAL_BIN)/goose
GOLANGCI_LINT    := $(LOCAL_BIN)/golangci-lint

LOCAL_MIGRATION_DIR := $(MIGRATION_DIR)
LOCAL_MIGRATION_DSN := host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable

lint:
	$(GOLANGCI_LINT) run ./... --config .golangci.yml

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.24.1
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4

#migrate-create:
#	$(GOOSE) -dir $(LOCAL_MIGRATION_DIR) create $(NAME) sql

migrate-status:
	$(GOOSE) -dir $(LOCAL_MIGRATION_DIR) postgres "$(LOCAL_MIGRATION_DSN)" status -v

migrate-up:
	$(GOOSE) -dir $(LOCAL_MIGRATION_DIR) postgres "$(LOCAL_MIGRATION_DSN)" up -v

migrate-down:
	$(GOOSE) -dir $(LOCAL_MIGRATION_DIR) postgres "$(LOCAL_MIGRATION_DSN)" down -v

up:
	docker compose up -d db

run: up migrate-up
	go run ./cmd/api

run-local: migrate-up
	@echo "Starting API on :8080 using local Postgres"
	DATABASE_DSN=$(DATABASE_DSN) go run ./cmd/api

down:
	docker compose down

clean:
	rm -rf $(LOCAL_BIN)

.PHONY: lint install-deps migrate-status migrate-up migrate-down up down clean
