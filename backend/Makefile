# ─── Variables ────────────────────────────────────────────────────────────────
BINARY     := trickreport-api
BUILD_DIR  := ./dist
CMD_PATH   := ./cmd/api
MIGRATIONS := ./migrations

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

## migrate: apply SQL migrations to the database
## Usage: make migrate DB_URL=postgres://user:pass@host:5432/dbname
migrate:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Error: DB_URL is required. Usage: make migrate DB_URL=..."; exit 1; \
	fi
	@for f in $(MIGRATIONS)/*.sql; do \
		echo "Applying $$f..."; \
		psql "$(DB_URL)" -f "$$f" || exit 1; \
	done
	@echo "All migrations applied."

## clean: remove build artifacts
clean:
	rm -rf $(BUILD_DIR)
