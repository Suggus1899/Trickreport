# ─── Variables ────────────────────────────────────────────────────────────────
BINARY     := trickreport-api
BUILD_DIR  := ./dist
CMD_PATH   := ./cmd/api

.PHONY: all run build test migrate clean lint fmt coverage docker-build docker-run migrate-up migrate-down generate install-tools security

## all: build the binary
all: build

## run: run the API in development mode
run:
	go run $(CMD_PATH)/main.go

## build: compile the binary
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) $(CMD_PATH)/main.go
	@echo "Built $(BUILD_DIR)/$(BINARY)"

## test: run all tests
test:
	go test ./... -v -race

## migrate: apply database migrations (runs the API, which applies them on startup)
migrate:
	go run $(CMD_PATH)/main.go

## clean: remove build artifacts
clean:
	rm -rf $(BUILD_DIR)

## lint: run golangci-lint
lint:
	golangci-lint run ./...

## fmt: format all Go sources
fmt:
	go fmt ./...

## coverage: run tests with coverage and print per-function coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

## docker-build: build the backend Docker image
docker-build:
	docker build -t trickreport-backend .

## docker-run: run the backend Docker image
docker-run:
	docker run -p 8080:8080 trickreport-backend

## migrate-up: apply migrations up (via the API --migrate-up flag)
migrate-up:
	go run $(CMD_PATH) --migrate-up

## migrate-down: roll back the last migration (via the API --migrate-down flag)
migrate-down:
	go run $(CMD_PATH) --migrate-down

## generate: run go generate across the module
generate:
	go generate ./...

## install-tools: install developer tooling (swag, golangci-lint, etc.)
install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest

## security: run govulncheck across the module
security:
	govulncheck ./...
