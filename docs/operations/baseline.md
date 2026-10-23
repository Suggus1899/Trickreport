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
