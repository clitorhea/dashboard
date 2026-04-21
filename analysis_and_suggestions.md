# NAS Docker Dashboard — Project Analysis & Improvements

## Overview

The project is a well-scoped self-hosted Docker management dashboard built with Go + Svelte 5 + SQLite. The foundation is solid — clean separation of concerns, proper auth middleware, WebSocket streaming, and a multi-stage Docker build. Below is a thorough review with issues and prioritized improvement ideas.

---

## 🐛 Bugs & Broken Features

### 1. Password Change is Unimplemented (Settings page is dead)
**Files:** `backend/handler/auth.go`, `frontend/src/pages/Settings.svelte`

The Settings page has a full password-change form, but the handler is missing. The frontend just shows a hardcoded "not yet implemented" message.

**Fix:** Add `PUT /api/auth/password` that verifies the current password with bcrypt, then hashes and stores the new one.

---

### 2. Template Seeding is Never Called
**Files:** `backend/main.go`, `backend/db/models.go`, `templates/*.yaml`

`db.SeedTemplates()` exists and is ready, but `main.go` never calls it. The built-in YAML templates in `templates/` are also never loaded. The Templates page will always be empty on a fresh install.

**Fix:** In `main.go` after `db.Init()`, load the YAML files from `./templates/` and call `db.SeedTemplates()`.

---

### 3. Docker Log Frame Header Stripping is Fragile
**File:** `backend/handler/containers.go` (lines 100–107)

The 8-byte header is stripped with a naive slice, but `reader.Read()` can return a partial buffer that spans multiple frames, or a buffer that starts mid-frame. This causes garbled log output.

**Fix:** Use `stdcopy.StdCopy()` from the Docker SDK which correctly demultiplexes the stream.

---

### 4. No SPA Fallback for Deep Links
**File:** `backend/main.go` (line 71)

The frontend uses hash-based routing (`#/containers`), so deep links mostly work. However, `http.FileServer` will 404 on any path that doesn't exist as a file. If someone navigates to `/` and the `index.html` isn't in `./static`, it breaks silently.

**Fix:** Wrap the file server with a fallback handler that serves `index.html` for any non-`/api` request that results in a 404.

---

### 5. `SeedTemplates` Silently Fails on Parse Errors
**File:** `backend/db/models.go` (lines 155–172)

If any template YAML is malformed, the error is returned but likely logged nowhere if the caller doesn't check it. Since seeding isn't called from `main.go` yet (bug #2), this is latent.

---

## ⚠️ Security Issues

### 6. WebSocket CORS is Wide Open
**File:** `backend/handler/containers.go` (line 13)

```go
CheckOrigin: func(r *http.Request) bool { return true },
```

This allows any origin to open a WebSocket connection. Since the dashboard is self-hosted, it should at least validate the origin matches the request's `Host`.

---

### 7. Session Cookie Missing `Secure` Flag
**File:** `backend/handler/auth.go` (lines 72–79, 109–116)

The session cookie has `HttpOnly` but not `Secure`. When accessed over HTTPS (via Cloudflare Tunnel), the cookie should be marked `Secure` to prevent it being sent over plain HTTP.

**Fix:** Set `Secure: true` conditionally, or always (since Cloudflare always terminates TLS).

---

### 8. Deploy Endpoint Has No Path Traversal Protection
**File:** `backend/handler/templates.go` (line 89)

```go
serviceDir := filepath.Join("./data/services", req.ServiceName)
```

A malicious `service_name` like `../../etc/passwd` could write files outside the intended directory. `filepath.Join` normalizes the path, so this is partially mitigated, but an explicit check against the base dir is safer.

---

### 9. Passwords are Sent Back in API Response
**File:** `backend/handler/auth.go` (line 118)

`writeJSON(w, http.StatusOK, user)` — the `User` struct has `json:"-"` on the `Password` field, so this is safe *for now*. But it's a pattern that's easy to break if someone adds a field.

---

## 🔧 Code Quality & Architecture

### 10. `system/info` Returns Go Runtime Info, Not Real Host Stats
**File:** `backend/handler/system.go`

`systemInfo` returns `runtime.GOOS`, `runtime.NumCPU()`, etc. — these reflect the container's view, not the real NAS host. The dashboard claims "Dashboard Uptime" but the hostname and CPU count are the container's, not the host's.

To get real host CPU/memory/disk usage, integrate `github.com/shirou/gopsutil` (a popular Go sysinfo library) or read from `/proc/meminfo` and `/proc/stat` directly.

---

### 11. No Rate Limiting on Login
**File:** `backend/handler/auth.go`

There is no rate limiting on `POST /api/auth/login`. An attacker could brute-force the admin password with no throttling.

**Fix:** Add an in-memory rate limiter per IP (e.g., `golang.org/x/time/rate`) — this dependency is already in `go.mod`.

---

### 12. Expired Sessions Not Cleaned Up
**File:** `backend/db/models.go`

Sessions are validated with a `WHERE expires_at > datetime('now')` check, but expired sessions are never deleted from the DB. Over time this table will grow unboundedly.

**Fix:** Run a periodic cleanup goroutine or a `DELETE FROM sessions WHERE expires_at < datetime('now')` on startup.

---

### 13. No Confirmation on Destructive Container Actions
**File:** `frontend/src/components/ContainerCard.svelte`

The stop/restart/remove buttons should prompt for confirmation, especially `remove`. Currently they fire immediately.

---

### 14. `go.sum` Was Generated Locally But `ws` Package Not Shown
The backend imports `github.com/gorilla/websocket` — verify this is in `go.mod` and `go.sum`. The `ws/` subdirectory listed in the architecture doesn't exist in the actual file tree (the WebSocket logic is in `handler/containers.go`).

---

## 🚀 Improvement & New Feature Ideas

### Priority 1 — Fix Known TODOs (Low Effort, High Value)

| # | Feature | Effort |
|---|---------|--------|
| A | `PUT /api/auth/password` endpoint + wire up Settings page | S |
| B | Template seeding on startup (load YAMLs → SQLite) | S |
| C | Fix Docker log demux with `stdcopy.StdCopy` | S |
| D | Session cleanup goroutine (runs every hour) | S |
| E | Login rate limiting (5 attempts / 1 min per IP) | S |

---

### Priority 2 — Quality of Life (Medium Effort)

| # | Feature | Description |
|---|---------|-------------|
| F | **Real host stats** | Use `gopsutil` for actual CPU %, RAM used/total, disk usage. Show live progress bars on the Dashboard. |
| G | **Confirmation dialogs** | "Are you sure?" modal before stop/remove/restart actions. |
| H | **Network management** | View Docker networks, which containers are attached. API: `GET /api/networks`. |
| I | **Volume management** | List volumes, prune dangling volumes. API: `GET /api/volumes`. |
| J | **Container rename/update** | Rename containers; update restart policy from the UI. |
| K | **Log search/filter** | Add a search input to the LogViewer to filter displayed lines client-side. |
| L | **Image tagging display** | Show all tags per image in the Images list (currently just shows ID). |

---

### Priority 3 — New Functionality (Larger Features)

| # | Feature | Description |
|---|---------|-------------|
| M | **Compose file editor** | In-browser YAML editor for deployed stacks (stored in `./data/services/`). Edit and re-deploy without touching the CLI. |
| N | **Notification system** | Alert when a container stops unexpectedly (poll every 30s, compare to last known state). Store alerts in SQLite. Badge on Navbar. |
| O | **Container exec / terminal** | WebSocket-based TTY into a running container (`docker exec -it`). Opens a terminal pane in the UI. |
| P | **Update detection** | Compare running image digest to the latest on Docker Hub. Show an "update available" badge per container. |
| Q | **Metrics history** | Persist CPU/mem samples to SQLite every minute. Show a time-series chart (last 1h/6h/24h) on the ContainerDetail page. |
| R | **Multi-user support** | Add roles (admin/viewer). Viewer can see containers and logs but cannot start/stop/deploy. |
| S | **Backup / export** | One-click export of all deployed compose files and their `.env` files as a `.tar.gz` archive. |
| T | **More built-in templates** | Add Jellyfin, Portainer, Nginx Proxy Manager, Vaultwarden, Uptime Kuma, Immich, Gitea, etc. |
| U | **Dark/light theme toggle** | Currently hardcoded dark. Add a CSS variable-based theme system with a toggle stored in `localStorage`. |
| V | **PWA / mobile-friendly** | Add a manifest and service worker so the dashboard can be pinned to a phone's home screen. |

---

## Recommended Next Steps (Ordered)

1. **Fix template seeding** (B) — templates are already written, just need to be wired up.
2. **Fix log streaming** (C) — currently produces garbled output.
3. **Add password change endpoint** (A) — the UI is already there.
4. **Add real host stats** (F) — makes the Dashboard page actually useful.
5. **Add session cleanup + rate limiting** (D, E) — low-effort security wins.
6. **Notification system** (N) — high value for a "set it and forget it" NAS tool.
7. **Metrics history** (Q) — turns the stats chart from a sparkle into a real tool.
