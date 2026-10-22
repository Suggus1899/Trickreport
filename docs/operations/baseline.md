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
failure, not evidence of clean coverage. GitNexus CLI and source/tests provide the
fallback. SQL, Astro scripts and dynamic calls require explicit source checks.

The stabilization plan is `docs/superpowers/plans/2026-10-21-estabilizacion-integral.md`.
Visual frontend work is deferred; security and API contract corrections remain
part of stabilization. Completion is recorded only after checks actually pass.
