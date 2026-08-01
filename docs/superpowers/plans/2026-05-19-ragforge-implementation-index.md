# RAGForge Implementation Plan - Index

Implementation broken into 5 plan files by phase:

| File | Phase | Description |
|------|-------|-------------|
| `2026-05-19-ragforge-implementation-phase1.md` | Phase 1 | Project scaffold: go.mod, cmd, config, app.go |
| `2026-05-19-ragforge-implementation-phase2.md` | Phase 2 | Model entities (10 tables) + DAO layer (9 files) |
| `2026-05-19-ragforge-implementation-phase3.md` | Phase 3 | Engine layer: LLM provider + Embedding provider |
| `2026-05-19-ragforge-implementation-phase4.md` | Phase 4 | Knowledge Base module: DTOs, Service, Controller, Router |
| `2026-05-19-ragforge-implementation-phase5.md` | Phase 5 | Remaining 7 modules + Middleware + Integration |

**Execution order:** Phases must be executed sequentially (each depends on previous).
