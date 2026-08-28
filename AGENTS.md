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
