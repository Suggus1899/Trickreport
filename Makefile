# Trickreport — Root Makefile for local development (no Docker)
# Local development ALWAYS runs without Docker.

.PHONY: dev backend frontend build test lint fmt install setup clean

## dev: start backend + frontend concurrently (local, no Docker)
dev:
	@echo "Starting Trickreport dev environment (local, no Docker)..."
	@bash ./scripts/dev.sh 2>/dev/null || powershell -ExecutionPolicy Bypass -File ./scripts/dev.ps1

## backend: start only the backend
backend:
	cd backend && go run cmd/api/main.go

## frontend: start only the frontend
frontend:
	cd frontend && npm run dev

## setup: copy .env files from examples (edit credentials before running)
setup:
	@if [ ! -f backend/.env ]; then cp backend/.env.example backend/.env; echo "Created backend/.env — edit DATABASE_URL with your credentials"; fi
	@if [ ! -f frontend/.env.local ]; then cp frontend/.env.example frontend/.env.local; echo "Created frontend/.env.local"; fi
	@cd frontend && npm install

## build: build both backend and frontend
build: build-backend build-frontend

## build-backend: compile the backend binary
build-backend:
	cd backend && go build -o dist/trickreport-api cmd/api/main.go

## build-frontend: build the frontend for production
build-frontend:
	cd frontend && npm run build

## test: run all tests (backend + frontend)
test: test-backend test-frontend

## test-backend: run backend tests
test-backend:
	cd backend && go test ./... -race

## test-frontend: run frontend tests
test-frontend:
	cd frontend && npx vitest run

## lint: lint both backend and frontend
lint: lint-backend lint-frontend

## lint-backend: lint the backend with golangci-lint
lint-backend:
	cd backend && golangci-lint run ./...

## lint-frontend: lint the frontend with eslint
lint-frontend:
	cd frontend && npm run lint

## fmt: format both backend and frontend
fmt: fmt-backend fmt-frontend

## fmt-backend: format Go sources
fmt-backend:
	cd backend && go fmt ./...

## fmt-frontend: format frontend sources with prettier
fmt-frontend:
	cd frontend && npm run format

## install: install frontend dependencies
install:
	cd frontend && npm install

## clean: remove build artifacts
clean:
	rm -rf backend/dist frontend/dist
