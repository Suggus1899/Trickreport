# AGENTS.md

## Project Overview

Trickreport — Go backend + Astro frontend, orchestrated with Docker Compose.

## Stack

- **Backend:** Go 1.25 (chi router, pgx/v5, squirrel, golang-jwt/v5, go-ldap/v3, viper, zerolog, prometheus client, swaggo). Module: `github.com/trickreport/backend`.
- **Frontend:** Node + Astro 5 (TypeScript), ESLint, Prettier, Husky + lint-staged.
- **Database:** PostgreSQL 16 (via `postgres:16-alpine`).
- **Infra:** Docker Compose (root), Caddy reverse proxy, monitoring stack.

## Common Commands

### Frontend (run from `frontend/`)

```bash
cd frontend && npm install
cd frontend && npm run dev        # dev server (astro dev)
cd frontend && npm run build      # production build
cd frontend && npm run preview    # preview built site
cd frontend && npm run lint       # eslint
cd frontend && npm run format     # prettier --write .
```

### Backend (run from `backend/`)

```bash
cd backend && go run cmd/api/main.go          # run API
cd backend && go build -o dist/trickreport-api cmd/api/main.go
cd backend && go test ./... -race             # tests
cd backend && golangci-lint run ./...         # lint
cd backend && go fmt ./...                    # format
cd backend && go run cmd/api --migrate-up     # apply migrations
cd backend && go run cmd/api --migrate-down   # rollback last migration
```

### Root Makefile shortcuts

```bash
make dev            # backend + frontend concurrently (local, no Docker)
make test           # test-backend + test-frontend
make lint           # lint both
make build          # build both
make setup          # copy .env examples + npm install
```

### Docker

```bash
docker-compose up -d                 # start postgres + app
docker-compose logs -f app
docker-compose down
```

Required env vars before `up`: `POSTGRES_PASSWORD`, `JWT_SECRET`, `ADMIN_EMAIL`, `ADMIN_PASSWORD`.

## Conventions

- **Commits:** Conventional Commits (`feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, `test:`). No AI attribution, no "Co-Authored-By" lines.
- **Go:** `gofmt` + `golangci-lint`; table-driven tests; `go test ./... -race`.
- **Frontend:** ESLint + Prettier enforced via Husky/lint-staged on `*.{js,ts,astro}` and `*.{css,json,md}`.
- **Migrations:** via API flags (`--migrate-up` / `--migrate-down`), not raw SQL scripts.

## SDD (Spec-Driven Development)

Run `/sdd-init` with cwd here. Stack: Go 1.25 backend (chi + pgx + JWT + LDAP) + Astro 5 frontend (TypeScript) + PostgreSQL 16 + Docker Compose. Testing: `go test ./... -race` (backend), `npx vitest run` (frontend).
