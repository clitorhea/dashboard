# NAS Docker Dashboard — Task List

## Phase 1: Project Setup
- [x] Initialize Go backend module and directory structure
- [x] Initialize Svelte frontend with Vite
- [x] Create Dockerfile and docker-compose.yml
- [x] Set up .gitignore

## Phase 2: Backend Core
- [x] Set up SQLite database connection and migrations
- [x] Define database models (User, Session, Template, Setting)
- [x] Implement Docker client wrapper (connect via socket)
- [x] Set up HTTP router and server entrypoint

## Phase 3: Auth
- [x] Implement user registration (initial setup / first-run)
- [x] Implement login endpoint (bcrypt verify + session create)
- [x] Implement logout endpoint (session destroy)
- [x] Implement auth middleware (session validation)
- [x] Implement `/api/auth/me` endpoint

## Phase 4: Container Management API
- [x] GET `/api/containers` — list all containers
- [x] GET `/api/containers/:id` — inspect container
- [x] POST `/api/containers/:id/start`
- [x] POST `/api/containers/:id/stop`
- [x] POST `/api/containers/:id/restart`
- [x] DELETE `/api/containers/:id` — remove container
- [x] WebSocket `/api/containers/:id/logs` — stream logs
- [x] WebSocket `/api/containers/:id/stats` — stream CPU/mem

## Phase 5: Image & System API
- [x] GET `/api/images` — list images
- [x] DELETE `/api/images/:id` — remove image
- [x] POST `/api/images/pull` — pull image
- [x] GET `/api/system/info` — host CPU, RAM, disk
- [x] GET `/api/system/docker` — Docker engine info

## Phase 6: Service Templates API
- [x] GET `/api/templates` — list templates
- [x] GET `/api/templates/:id` — template detail
- [x] POST `/api/services/deploy` — deploy from template
- [x] Create built-in templates (Plex, Nextcloud, Pi-hole, etc.)

## Phase 7: Frontend — Layout & Auth
- [x] Set up Svelte routing (svelte-spa-router or similar)
- [x] Create app layout (Navbar, sidebar, main content area)
- [x] Build Login page
- [x] Implement auth store and API client with session handling
- [x] Add route guards (redirect to login if unauthenticated)

## Phase 8: Frontend — Dashboard & Containers
- [x] Build Dashboard page (overview cards: total containers, running, stopped, system stats)
- [x] Build Containers list page with status indicators
- [x] Build ContainerDetail page (info, env, ports, volumes)
- [x] Add start/stop/restart/remove actions with confirmation
- [x] Build LogViewer component (WebSocket streaming)
- [x] Build StatsChart component (live CPU/mem graphs)

## Phase 9: Frontend — Images, Templates & Settings
- [x] Build Images page (list, remove, pull)
- [x] Build Templates page (browse available services)
- [x] Build DeployModal (configure and deploy from template)
- [x] Build Settings page (change password, dashboard config)

## Phase 10: Deployment & Remote Access
- [x] Finalize multi-stage Dockerfile (Go build + Svelte build → minimal image)
- [x] Configure Cloudflare Tunnel in docker-compose
- [x] Add environment variable / secrets documentation
- [ ] Test full deployment on target NAS
- [x] Write user-facing README with setup instructions

## Phase 11: Bug Fixes & Improvements
- [x] Add `PUT /api/auth/password` endpoint
- [x] Wire Settings page to real password change endpoint
- [x] Seed built-in templates from `templates/` dir on startup
- [x] Fix Docker log demux using `stdcopy.StdCopy` (was garbled before)
- [x] Add session cleanup goroutine (runs every hour)
- [x] Add login rate limiter (5 failures / min per IP)
- [x] Add `Secure` flag to session cookies
- [x] Add WebSocket origin validation
- [x] Fix path traversal in deploy endpoint
- [x] Real host CPU/RAM/disk stats via `gopsutil/v3`
- [x] Confirmation dialogs on Stop/Restart container actions
- [x] Notification system — background polling for container health
- [x] Notification badge in Navbar
- [x] SPA 404 fallback for deep links
- [x] Add 6 more service templates (Jellyfin, Vaultwarden, NPM, Immich, Gitea, Uptime Kuma)
