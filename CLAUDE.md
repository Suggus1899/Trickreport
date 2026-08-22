# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Trickreport — multi-tenant help desk / ticketing platform. Go 1.25 backend (hexagonal/clean architecture) + Astro 5 SSR frontend (no UI framework, hand-rolled neumorphic CSS), PostgreSQL 16, orchestrated with Docker Compose. See [AGENTS.md](AGENTS.md) for the full command reference — the essentials are repeated below.

## Commands

```bash
# Backend (from backend/)
go run cmd/api/main.go            # run API (also runs embedded migrations + admin bootstrap)
go test ./... -race               # tests
golangci-lint run ./...           # lint
go run cmd/api --migrate-up       # apply migrations
go run cmd/api --migrate-down     # rollback last migration

# Frontend (from frontend/)
npm run dev                       # astro dev
npm run build                     # astro build
npx vitest run                    # tests (currently one suite: src/lib/api.test.ts)
npm run lint                      # eslint

# Root
make dev / make test / make lint / make build   # both sides at once
docker-compose up -d              # postgres + app (needs POSTGRES_PASSWORD, JWT_SECRET, ADMIN_EMAIL, ADMIN_PASSWORD)
```

Run a single Go test: `go test ./internal/application/ticket/... -run TestName -race`.

## Backend architecture

`backend/internal` follows strict hexagonal/clean layering. Dependencies point inward only — never skip a layer or import outward:

```
interfaces/http (chi handlers, middleware, DTOs)
        ↓ calls
application/<domain> (services — orchestration, no SQL, no HTTP)
        ↓ depends on interfaces defined by
domain/<domain> (entities, value objects, domain errors — no external deps)
        ↑ implemented by
infrastructure/postgres, infrastructure/auth, infrastructure/email, infrastructure/realtime, infrastructure/storage
```

- **Domains**: `ticket`, `user`, `auth`, `article` (knowledge base), `sla`, `automation`, `analytics` — each has its own `domain/<name>`, `application/<name>`, and (where needed) rows in `infrastructure/postgres`.
- **Wiring is manual and centralized** in [backend/internal/interfaces/http/wire.go](backend/internal/interfaces/http/wire.go): `NewRepos` → `NewServices` → `NewHandlers`. There is no DI framework/codegen (no wire/fx) — when adding a new repo/service/handler, wire it here by hand, in that order. Services are composed via setter methods after construction (e.g. `ticketSvc.SetEngine(engine)`, `authSvc.SetMFARepository(...)`) rather than large constructors — follow that pattern for new cross-service dependencies.
- **Multi-tenancy**: every row-owning entity carries a `TenantID`; `infrastructure/postgres/tenant_resolver.go` (+ `cached_tenant_resolver.go`) resolves tenant context. Always scope new repo queries by tenant.
- **Automation engine**: `application/automation` (`Engine`) evaluates SLA/automation rules and calls `infrastructure/postgres/automation_executor.go`, which runs raw SQL directly on the pool (bypassing the repo layer) for rule actions — this is a deliberate exception to the layering above.
- **Migrations are embedded and consolidated into a single file**: `internal/db/migrations/000001_init.{up,down}.sql` (prior history was squashed — see git log). Applied automatically by `db.RunMigrations` on every API startup, and manually via `--migrate-up`/`--migrate-down`.
- **Realtime**: `internal/realtime` (hub/client/message) is the low-level WebSocket hub; `infrastructure/realtime/adapter.go` bridges it to application services (e.g. ticket + notification services broadcast on mutation).
- **Logging**: single file `backend/trickreport.log` (zerolog), no `logs/` directory. Dev also mirrors to console via `io.MultiWriter`; production writes JSON to the file only — see [backend/cmd/api/main.go](backend/cmd/api/main.go).
- **Admin bootstrap**: `internal/bootstrap.EnsureAdmin` creates the initial admin user from `ADMIN_EMAIL`/`ADMIN_PASSWORD` on every startup (idempotent).
- API docs are swaggo-generated from annotations in `cmd/api/main.go` and handler files, served via `interfaces/http/handler/swagger.go`.

## Frontend architecture

Astro 5 in SSR mode (`@astrojs/node`, standalone adapter) — **no React/Vue/Svelte island framework**; interactivity is vanilla TypeScript in `<script>` blocks inside `.astro` files. There is no Tailwind despite what the README's tech-stack table implies — styling is a single hand-rolled neumorphic design system in [frontend/src/styles/global.css](frontend/src/styles/global.css) (~800 lines), applied through reusable `Neu*` components (`NeuButton`, `NeuCard`, `NeuModal`, etc. in `src/components/`).

- `src/pages/**/*.astro` are file-based routes (SSR pages call the backend directly server-side, then hydrate client state via inline scripts). `src/pages/admin/*` are admin-only views.
- `src/lib/api.ts` is the single fetch client wrapper for the Go API — extend it rather than calling `fetch` ad hoc from pages.
- `src/lib/offline.ts` + `src/scripts/register-sw.ts` + `pages/tickets/new-offline.astro` / `sync-status.astro` implement PWA offline ticket creation with background sync — read `offline.ts` before touching ticket creation flows, since online and offline paths must stay consistent.
- `src/lib/keyboard.ts` drives the global keyboard-shortcut system (paired with `KeyboardHelp.astro`).

## Conventions (also in AGENTS.md)

- Conventional Commits, no AI attribution (no `Co-Authored-By`, no "Generated with").
- Go: `gofmt` + `golangci-lint`, table-driven tests.
- Migrations go through the API's `--migrate-up`/`--migrate-down` flags, not raw SQL scripts.
