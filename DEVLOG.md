# NAS Docker Dashboard — Dev Log

## 2026-04-06 — Session 1: Project Planning

### What was done
- Discussed architecture options and tech stack trade-offs
- Chose tech stack: Go + Svelte + SQLite + Cloudflare Tunnel
- Created `ARCHITECTURE.md` with full system design
- Created `TASKLIST.md` and `DEVLOG.md`
- Initialized git repository

### Decisions made
- Go over Node.js, Svelte over React, SQLite over Postgres, Cloudflare Tunnel over VPN

---

## 2026-04-06 — Session 1 (continued): Full Implementation

### What was done

**Backend (Go)**
- `backend/main.go` — HTTP server with all routes, serves static frontend
- `backend/db/db.go` — SQLite connection, WAL mode, auto-migrations for all 4 tables
- `backend/db/models.go` — User, Session, Template models + all CRUD methods + template seeding
- `backend/docker/client.go` — Docker SDK wrapper (containers, images, system info, log/stats streaming)
- `backend/middleware/auth.go` — Session cookie validation middleware
- `backend/handler/auth.go` — Login, logout, setup (first-run), me endpoints
- `backend/handler/containers.go` — List, inspect, start/stop/restart/remove, WebSocket logs + stats
- `backend/handler/images.go` — List, remove, pull with streaming progress
- `backend/handler/system.go` — Host info + Docker engine info
- `backend/handler/templates.go` — Template CRUD + deploy via `docker compose up`

**Frontend (Svelte 5)**
- `src/lib/api.js` — Fetch wrapper with auto-redirect on 401
- `src/lib/stores.js` — User and loading stores
- `src/lib/ws.js` — WebSocket connection helper
- `src/App.svelte` — Root with router, auth check, first-run setup detection
- `src/components/Navbar.svelte` — Sidebar navigation with active state via hashchange
- `src/components/ContainerCard.svelte` — Container card with status dot and actions
- `src/components/LogViewer.svelte` — WebSocket-based real-time log viewer
- `src/components/StatsChart.svelte` — Live CPU/mem/network stats with progress bars
- `src/components/DeployModal.svelte` — Modal for deploying from template or custom compose
- `src/pages/Login.svelte` — Login form
- `src/pages/Setup.svelte` — First-run account creation
- `src/pages/Dashboard.svelte` — Overview with container counts + system info
- `src/pages/Containers.svelte` — Container grid with running/stopped filter
- `src/pages/ContainerDetail.svelte` — Detailed view with tabs (stats, logs, env, volumes)
- `src/pages/Images.svelte` — Image table with pull and remove
- `src/pages/Templates.svelte` — Template browser + custom deploy button
- `src/pages/Settings.svelte` — Account info + password change placeholder

**Infrastructure**
- `Dockerfile` — Multi-stage build (Node → Go → Alpine with docker-cli)
- `docker-compose.yml` — Dashboard + Cloudflare Tunnel (tunnel behind a profile)
- `templates/` — Built-in compose templates for Plex, Nextcloud, Pi-hole
- `vite.config.js` — Dev proxy to backend + build config
- Frontend builds successfully (77KB JS + 17KB CSS gzipped)

### Known issues / TODOs
- Go can't be compiled locally (not installed) — builds happen in Docker
- `go.sum` not yet generated (will be created on first `go mod tidy` in Docker build)
- Settings page password change needs a backend endpoint
- No `go.sum` file — the Dockerfile will need `go mod download` or `go mod tidy` first

### Next steps
- Generate `go.sum` via Docker or install Go locally
- Test Docker build end-to-end
- Write README with setup instructions
- Test on actual NAS hardware

---

## 2026-04-06 — Session 1 (part 3): Build Validation & README

### What was done
- Wrote `README.md` with quick start, Cloudflare Tunnel setup, dev instructions, env vars, project structure
- Downloaded Go 1.23.8 locally to validate the backend compiles
- Ran `go mod tidy` — resolved all dependencies, generated `go.sum`
- Fixed `golang.org/x/time` version conflict (v0.15.0 required go 1.25, pinned to v0.9.0 for go 1.23 compat)
- Successfully compiled backend: **15MB binary**, clean build
- Frontend builds clean: **77KB JS + 17KB CSS** (gzipped: 27KB + 3KB)
- Updated Dockerfile: pinned `golang:1.23.8-alpine`, set `GOTOOLCHAIN=local`, proper layer caching with separate go.mod/go.sum COPY

### Known issues / TODOs
- Settings page password change still needs a backend endpoint (`PUT /api/auth/password`)
- No Docker runtime on this dev machine — full `docker compose up` test pending on NAS
- Built-in templates (Plex, Nextcloud, Pi-hole) exist as YAML files but aren't seeded into SQLite yet at startup

### Next steps
- Test full Docker build + deploy on NAS
- Add template seeding on startup (load from `templates/` dir into SQLite)
- Add password change endpoint
