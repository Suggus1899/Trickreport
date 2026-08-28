# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Trickreport — multi-tenant help desk / ticketing platform. Go 1.25 backend (hexagonal/clean architecture) + Astro 5 SSR frontend with React islands (Tailwind + shadcn/ui), PostgreSQL 16, orchestrated with Docker Compose. See [AGENTS.md](AGENTS.md) for the full command reference — the essentials are repeated below.

## Commands

```bash
# Backend (from backend/)
go run cmd/api/main.go            # run API (also runs embedded migrations + admin bootstrap)
go test ./... -race               # tests
golangci-lint run ./...           # lint
go run cmd/api --migrate-up       # apply migrations
go run cmd/api --migrate-down     # rollback last migration
go run cmd/api --seed             # seed demo data (users/tickets/articles/SLA/automation rule); refuses to run when ENV=production

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

Astro 5 in SSR mode (`@astrojs/node`, standalone adapter) with **React islands** (`@astrojs/react`, hydrated via `client:*` directives) for interactive components — plain `.astro` files still handle SSR data-fetching and page shells. Styling is Tailwind + shadcn/ui primitives in [frontend/src/components/ui](frontend/src/components/ui) (the older hand-rolled neumorphic `global.css`/`Neu*` system was removed).

- `src/pages/**/*.astro` are file-based routes (SSR pages call the backend directly server-side, then render React islands for interactive parts). `src/pages/admin/*` are admin-only views.
- `src/lib/api.ts` is the single fetch client wrapper for the Go API — extend it rather than calling `fetch` ad hoc from pages. SSR code uses `API_URL`; browser-side code (React islands) must use `PUBLIC_API_URL` instead, since Astro/Vite only inline `PUBLIC_`-prefixed env vars into the client bundle.
- `src/lib/offline/{db,sync}.ts` + `src/scripts/register-sw.ts` + `pages/tickets/new-offline.astro` / `sync-status.astro` implement PWA offline ticket creation with background sync — read `offline/` before touching ticket creation flows, since online and offline paths must stay consistent.
- `src/components/chrome/KeyboardHelp.tsx` drives the global keyboard-shortcut system.

## Conventions (also in AGENTS.md)

- Conventional Commits, no AI attribution (no `Co-Authored-By`, no "Generated with").
- Go: `gofmt` + `golangci-lint`, table-driven tests.
- Migrations go through the API's `--migrate-up`/`--migrate-down` flags, not raw SQL scripts.

<!-- gitnexus:start -->
# GitNexus — Code Intelligence

This project is indexed by GitNexus as **Trickreport** (3962 symbols, 12553 relationships, 306 execution flows).

> Index stale? Run `node .gitnexus/run.cjs analyze --index-only` from the project root — it auto-selects an available runner. No `.gitnexus/run.cjs` yet? Bootstrap with `npx`, `bunx`, or `pnpm dlx` — e.g. `bunx gitnexus@latest analyze` (npm 11 npx crash; #1939).

## Always Do

- **MUST run impact analysis before editing.** Use `impact({target: "symbolName", direction: "upstream"})` (MCP) or `node .gitnexus/run.cjs impact "symbolName" --direction upstream --repo .` (CLI fallback); report callers, processes, and risk. Never substitute grep for graph analysis.
- **MUST analyze graph changes before committing.** Use `detect_changes({scope: "all"})` (MCP) or `node .gitnexus/run.cjs detect-changes --scope all --repo .` (CLI fallback). `partial: true` or `truncated: true` is not a clean check — a zero means unseen, not unaffected; re-run it. For regression review: `detect_changes({scope: "compare", base_ref: "main"})` or `node .gitnexus/run.cjs detect-changes --scope compare --base-ref "main" --repo .`.
- **MUST warn the user** if impact analysis returns HIGH or CRITICAL risk before proceeding with edits.
- **MUST treat `risk: UNKNOWN` as unresolved, not as low.** An empty caller set is not evidence the symbol is unused — it can also mean the callers are not resolvable by the index (plain-object property access, dynamic dispatch, cross-language calls). `impact` pairs `UNKNOWN` with a `riskNote` saying so. Confirm with a text search before treating the symbol as safe to change or delete; do not proceed on the strength of a zero.
- When exploring unfamiliar code, use `query({search_query: "concept"})` to find execution flows instead of grepping. It returns process-grouped results ranked by relevance.
- When you need full context on a specific symbol — callers, callees, which execution flows it participates in — use `context({name: "symbolName"})`.
- For security review, `explain({target: "fileOrSymbol"})` lists taint findings (source→sink flows; needs `analyze --pdg`).

## Never Do

- NEVER edit a function, class, or method before MCP/CLI impact analysis.
- NEVER ignore HIGH or CRITICAL risk warnings from impact analysis, and never read `UNKNOWN` as an all-clear — it means the walk could not answer, which is the one verdict that requires confirming by other means.
- NEVER rename symbols with find-and-replace — use `rename` which understands the call graph.
- NEVER commit before MCP/CLI graph change analysis.

## Resources

| Resource | Use for |
| --- | --- |
| `gitnexus://repo/Trickreport/context` | Codebase overview, check index freshness |
| `gitnexus://repo/Trickreport/clusters` | All functional areas |
| `gitnexus://repo/Trickreport/processes` | All execution flows |
| `gitnexus://repo/Trickreport/process/{name}` | Step-by-step execution trace |

## CLI

| Task | Read this skill file |
| --- | --- |
| Understand architecture / "How does X work?" | `.claude/skills/gitnexus-exploring/SKILL.md` |
| Blast radius / "What breaks if I change X?" | `.claude/skills/gitnexus-impact-analysis/SKILL.md` |
| Trace bugs / "Why is X failing?" | `.claude/skills/gitnexus-debugging/SKILL.md` |
| Rename / extract / split / refactor | `.claude/skills/gitnexus-refactoring/SKILL.md` |
| Tools, resources, schema reference | `.claude/skills/gitnexus-guide/SKILL.md` |
| Index, status, clean, wiki CLI commands | `.claude/skills/gitnexus-cli/SKILL.md` |

<!-- gitnexus:end -->
