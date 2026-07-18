# ─── Variables ────────────────────────────────────────────────────────────────
BINARY     := trickreport-api
BUILD_DIR  := ./dist
CMD_PATH   := ./cmd/api

.PHONY: all run build test migrate clean

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
