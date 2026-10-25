# Stabilization baseline — 2026-10-22

Work proceeds on local `main` at the user's request. Starting product HEAD:
`382b6417393f011ff41a1c8cb919a2385055fbe6`. The uncommitted formatting change in
`frontend/src/components/articles/ArticleEditor.tsx` is preserved.

Remote history was rewritten. Different commit hashes mostly contain equivalent
patches; the three local authentication, reset and article-editor improvements
remain present. Local main was advanced from its verified ancestor to this HEAD.
No remote branch is overwritten or force-pushed. Remote integration must preserve
the local functional improvements and the corrected realtime test expectation.

Codebase Memory is registered but its session tools currently return `Transport
closed`, including project listing, status and coverage. This is an availability
failure, not evidence of clean coverage. Graphify and source/tests provide the
fallback. SQL, Astro scripts and dynamic calls require explicit source checks.

The stabilization plan is `docs/superpowers/plans/2026-10-21-estabilizacion-integral.md`.
Visual frontend work is deferred; security and API contract corrections remain
part of stabilization. Completion is recorded only after checks actually pass.

## Addendum — 2026-10-23, Stage 0 on top of the stabilization merge

`main` received the platform-stabilization merge (`5668f4c`) and a conflict
resolution (`ed1d631`) today. The conflict resolution resurrected code that the
merge had removed; this addendum records what was verified and decided on the
resulting HEAD. Runtimes used: Go 1.26.10, Node 24.22.1, PostgreSQL 16.15.

**Security fix (`6ce4708`).** The resurrected
`backend/internal/infrastructure/storage/file_validator.go` once again allowed
`image/svg+xml` uploads. The effective upload path
(`application/ticket/attachment_validation.go`) already rejected SVG, and the
storage validator had no callers, but the whitelist itself was a stored-XSS
hole waiting for the planned Storage wiring. SVG is removed from the storage
whitelist, with regression tests on both validation paths.

**Base verification (point: "base comprobable").** All of the following ran
against the current HEAD and passed:

- Backend: `go build ./...`, `go vet ./...`, `gofmt -l` clean (16 files were
  reformatted in `28142bc`; formatting only).
- Backend unit suite: `go test ./... -race` — 27 packages ok. On this 2-vCPU
  host the suite must run with limited parallelism; a fully parallel run
  starves the race-instrumented test builds and reports spurious
  `[build failed]` results that do not reproduce per-package.
- Backend integration and contract suites with `TEST_DATABASE_URL` set
  against a disposable PostgreSQL 16.15: `infrastructure/postgres`,
  `interfaces/http/handler` and `worker` all pass with **zero skips**
  (previously these silently skipped without a database).
- Workspaces (Next.js apps): `npm run check` (tsc), `npm run lint`,
  `npm test` (2/2) and `npm run build` for all three apps pass.
- Astro app (`frontend/`): `npm run check` went from 21 errors to **0**,
  `npm run lint` from 3 errors to **0**, `npm run test:ci` passes 71/71 and
  `npm run build` succeeds (commit `af20463`). 34 `astro check` hints remain —
  deprecations (`FormEvent`, recharts `Cell`), not defects.

**Identity wiring completed (`af20463`).** B1.5's identity-scoped offline API
was only half-wired: both Astro pages already passed an `identity` prop that
`SyncStatus`, `OfflineTicketForm` and `useOfflineTickets` never accepted, so
they called the scoped API with the old arities. The components now thread
`OfflineIdentity` through every list/save/sync call.

**Bootstrap password policy (`de911cd`).** `EnsureAdmin` now validates
`ADMIN_PASSWORD` with the domain complexity policy before touching the
database; deployed environments already refused empty or `change_me` values
in `config.Load`. No `changeme` default remains anywhere in the repository.

**Design decision — resurrected dead code removed.** The conflict resolution
also brought back the neumorphic design set (`frontend/src/components/Neu*.astro`,
5 components), the 798-line `frontend/src/styles/global.css`, and the legacy
`frontend/src/lib/offline.ts`. None of them is referenced by any page, layout
or component (verified by search), and `lib/offline.ts` carries the unsafe
pre-B1.5 semantics without identity scoping. All seven files are deleted
(commits `eebf068`…`063047e`). The live transitional design system is Tailwind
plus `frontend/src/components/ui/*`; per the master plan it is replaced by
Next.js + HeroUI in Stage 2, at which point the whole Astro app is removed.

**History note.** The history contains a duplicated patch pair
(`8e36f0b` = `aefadf9`). `main` history is not rewritten; the duplication is
recorded here only.

**Still pending after this addendum.** Playwright E2E was not executed in
this environment (B3.4 remains a CI gate, not a verified pass); Codebase
Memory is still unavailable as noted above; the remaining Stage-0 blocks of
the stabilization plan continue from this verified base.

## Correction addendum — 2026-10-23 (workspace lockfile and bundler)

The addendum above stated that the root workspace build passed; that was
only partially verified at the time. A from-scratch `npm ci` followed by
`npm run build` on Linux x64 actually failed, for two independent reasons
found and fixed here:

1. **Incomplete root `package-lock.json`.** The lockfile had no package
   entries for 42 platform-specific optional dependencies — the native
   bindings `@tailwindcss/oxide-*` (11), `@next/swc-*` (7) and
   `@img/sharp-*` / `@img/sharp-libvips-*` (24) — so a Linux `npm ci`
   installed everything except the binaries the build needs. The entries
   were completed from the npm registry metadata (version, `resolved`,
   `integrity`, `cpu`, `os`, `libc`); no pinned version changed.
2. **Turbopack cannot resolve the hoisted oxide binding.** With the
   binding installed and loadable by plain Node, `next build` (Turbopack)
   still failed evaluating `@tailwindcss/postcss` for `globals.css`:
   `Cannot find module '@tailwindcss/oxide-linux-x64-gnu'` from inside
   Turbopack's Node.js evaluation in this npm-workspaces layout. The same
   app builds cleanly with webpack. Until Turbopack resolves hoisted
   optional native dependencies in workspaces, the three apps use
   `next build --webpack` and `next dev --webpack` (one-line revert per
   app in `apps/*/package.json`).

Re-verified after the fix, sequentially, from a clean `npm ci`:
`npm run check`, `npm run lint` and `npm run test` (2/2) pass, and
`npm run build` builds `@trickreport/customer-service`, `@trickreport/it`
and `@trickreport/msp`.

## Addendum — 2026-10-25, closure of Stage E0 (local)

Scope: the E0 items of `docs/superpowers/plans/2026-10-25-cierre-produccion.md`,
done on local `main` only. Nothing was pushed; `origin/main` stays at
`9647012`. The E0 checkbox in the plan remains **unchecked**: it requires the
seven CI jobs to be green on one commit, which needs a push that has not been
authorized.

Environment: Go 1.26.10; golangci-lint 2.15.0; Node v25.8.2 (the apps declare
`>=24.22.1 <25`, so npm prints EBADENGINE); npm 11.11.1; PostgreSQL 18.4 as a
disposable local instance on `127.0.0.1:55432` (CI uses 16.15);
`TZ=America/Caracas`. No C compiler is installed, so `go test -race` cannot run
on this host. Docker/Compose, Caddy, Trivy, govulncheck and promtool are not
installed.

### Commits (local, not pushed)

| Commit | Item | Content |
| --- | --- | --- |
| `286f350` | plan | production closure plan |
| `bdeb5e8` | E0-a | provider webhooks accepted without a browser `Origin` (Caddyfile, `proxy.ts`, tests) |
| `3a70d33` | E0-b | client scope on MSP project tasks and remote sessions; integration test |
| `cb176fb` | E0-c | stored contract fields kept; one datetime conversion |
| `cfd1759` | E0-d | SLA deadlines for MSP and catalog tickets (11 files) |
| `dd45902` | E0-e | sidebar sign-out reachable on short viewports (CSS); E2E spec |
| `810e41c` | E0-f | CX conversation refreshed after actions; portal users linkable to contacts |
| `be55532` | E0-g | gofmt for CX domain structs |
| `d11e7a3` | E0-g | errcheck and ineffassign fixes in repositories |
| `c155e96` | E0-g | unused helpers removed; TODO for metrics |
| `ffd4e7b` | E0-g | unused test mocks and no-op assertions removed |
| `fb98721` | E0-g | PostCSS 8.5.31 (clears the high advisory) |
| `ef7498b` | test | outbox assertion pinned to its own message |
| `f5294f8` | docs | README and agent guides aligned with the E0 state |

### Verification on this HEAD

- Backend: `go build ./...`, `go vet ./...` and `gofmt -l` clean.
- Backend suite with `TEST_DATABASE_URL` set: `go test -count=1 -p 1 -v ./...`
  exits 0. 32 packages ok, 0 FAIL, 0 SKIP, 758 PASS entries. Four packages have
  no test files: `cmd/api`, `application/cx`, `application/msp`,
  `domain/automation`. Named checks that pass include the outbox restart and
  rollback test, the MSP contact scope test (4 subtests), MSP and catalog ticket
  deadline tests, the SLA worker concurrency and history-rollback tests, and
  the SLA breach test for MSP tickets.
- Workspace (`packages/workspace`, Vitest): 27/27 in 4 files (proxy 5, itsm 5,
  cx 7, msp 10). The MSP contract test passes in `America/Caracas` and `UTC`.
- Next apps: `npm run check`, `npm run lint` and `npm run build` pass after the
  PostCSS bump. `npm audit --audit-level=high` reports 0.
- Playwright E2E `apps/it/tests/integration/service.spec.ts`: 2 passed (PostgreSQL
  18.4, API on 8088, Next on 3000). It was red three times before the CSS fix.
  This run was made earlier in the stage, not on the final HEAD.
- golangci-lint with `--max-same-issues=0`: **143 issues**, exit 1 (errcheck 50,
  gocyclo 11, gosec 32, revive 50). `ineffassign`, `unused`, `staticcheck`,
  `govet` and `misspell` report 0. The plan says 132; the cause of the
  difference is not established.

### Status per item

- **E0-a:** proxy tests 5/5 pass. `graphify-out/verify-webhook-boundary.mjs`
  does not exist and was not run.
- **E0-b:** integration test passes (4 subtests).
- **E0-c:** the regression is a Vitest test (`msp.test.ts`, 10/10). There is no
  Go test for the partial update.
- **E0-d:** Go tests for deadlines and the worker pass. Known limits: no SLA
  creation history; ITSM approval is not atomic with deadline creation; no
  backfill for existing tickets; MSP tickets have no escalation; the ITSM
  policy validates without the 525600-minute cap.
- **E0-e:** the E2E failure was fixed by the sidebar CSS change. The Select
  label hypothesis was tested and discarded, so the label association was not
  changed.
- **E0-f:** no tests added, no E2E, no targeted Vitest. The backend validates
  `user_id` against the tenant only, not against end users. Inactive users stay
  selectable. A reload failure after a successful action is shown as an error
  (minor).
- **E0-g:** npm audit is clean. Lint is not closed (143). The five high Trivy
  findings for the image were not scanned. `application/msp` and
  `application/cx` still have no test files; `application/itsm` now has tests.

### Not verified

- Prettier: declared in `frontend/package.json` but not installed. The husky
  `pre-commit` hook runs `npx lint-staged` from the repository root, which has
  no config. Not changed here.
- The seven remote CI jobs (no push).
- `go test -race` (no C compiler on this host).
- E2E and a targeted Vitest for E0-f.
- Docker Compose stack, Caddy routing, Trivy, govulncheck, promtool and
  `restore.test.sh`.
- GitNexus index was not available; impact checks were done with grep and the
  test suites.
- Outbox test: the assertion is now pinned to the test's own message. The
  cross-package race was not reproduced in parallel runs.

### Follow-ups (outside E0)

- CX: validate contact `user_id` against end users; hide inactive users from
  the selector; separate reload failures from action failures.
- Go tests for `application/msp` and `application/cx`.
- SLA: record deadline creation history; make ITSM approval atomic with
  deadlines; backfill existing tickets; escalate MSP tickets; apply the
  525600-minute cap in the ITSM policy.
- Tests in different packages share one database (`TEST_DATABASE_URL`); give
  each package its own database or schema.
- `middleware_test`: the next handler does not assert the tenant.
- Environment parity: local Node 25 and PostgreSQL 18.4 versus CI Node 24 and
  PostgreSQL 16.15.

### Decisions for the user

1. Authorize a push to `origin/main` so the seven CI jobs can run. The E0
   checkbox can only be closed after that.
2. Lint: fix the 143 findings, or change `backend/.golangci.yml` (for example
   `max-same-issues`). Changing the config would hide findings, so it needs a
   decision.
3. Accept the PostCSS lockfile change: three nested postcss copies removed, and
   twelve `@tailwindcss/oxide-*` packages now marked dev because their parent is
   a dev dependency. No production dependency, version or integrity changed.
4. Render RPO: option A (accept and document the RPO; recommended to start) or
   option B (synchronous replica), as set out in the plan.
5. Vercel Functions request-body limit of 4.5 MB (413): decide how large
   request bodies reach the Go API.
6. SLA: history rule for deadline changes, backfill of existing tickets, and the
   525600-minute cap.
7. MSP contacts linked without `user_id`, and CX end-user validation.
8. Per-package test database isolation.
