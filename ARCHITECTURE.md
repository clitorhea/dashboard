# NAS Docker Dashboard — Architecture Overview

## Overview

A self-hosted web dashboard for managing Docker services on an Ubuntu-based NAS server. Provides a browser UI to view, manage, and deploy Docker containers, accessible both locally and remotely via Cloudflare Tunnel.

## Tech Stack

| Layer          | Technology     | Rationale                                              |
| -------------- | -------------- | ------------------------------------------------------ |
| Frontend       | Svelte + Vite  | Smallest bundle, fast on low-power NAS hardware        |
| Backend        | Go (net/http)  | Single binary, ~10MB RAM, first-class Docker SDK       |
| Database       | SQLite         | Zero setup, single file, perfect for single-user       |
| Auth           | Session + bcrypt | Simple session cookies, bcrypt-hashed passwords      |
| Remote Access  | Cloudflare Tunnel | Free, no port forwarding, built-in auth layer       |
| Containerization | Docker Compose | Dashboard itself runs as a container                 |

## System Architecture

```
Internet
    │
    ▼
┌──────────────────┐
│ Cloudflare Tunnel │  ← cloudflared container
│ (nas.domain.com)  │
└────────┬─────────┘
         │ :443
         ▼
┌──────────────────┐
│  Dashboard App    │  ← single Docker container
│                   │
│  ┌─────────────┐  │     ┌──────────┐
│  │  Frontend    │  │     │ SQLite   │
│  │  (Svelte)   │  │     │ data.db  │
│  ├─────────────┤  ├────►│          │
│  │  Backend     │  │     └──────────┘
│  │  (Go API)   │  │
│  └──────┬──────┘  │
│         │         │
└─────────┼─────────┘
          │
          ▼
┌──────────────────────────────────────┐
│  /var/run/docker.sock                │
│                                      │
│  Docker Engine                       │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌─────┐ │
│  │ Plex │ │ *arr │ │ NGINX│ │ ... │ │
│  └──────┘ └──────┘ └──────┘ └─────┘ │
└──────────────────────────────────────┘
```

## Project Structure

```
dashboard/
├── backend/
│   ├── main.go                 # entrypoint
│   ├── go.mod
│   ├── handler/
│   │   ├── containers.go       # container CRUD + actions
│   │   ├── images.go           # image list/remove
│   │   ├── system.go           # host stats (CPU, RAM, disk)
│   │   ├── templates.go        # service template CRUD
│   │   └── auth.go             # login/logout/session
│   ├── middleware/
│   │   └── auth.go             # session validation middleware
│   ├── docker/
│   │   └── client.go           # Docker SDK wrapper
│   ├── db/
│   │   ├── db.go               # SQLite connection + migrations
│   │   └── models.go           # User, Template, Setting models
│   └── ws/
│       └── logs.go             # WebSocket log/stats streaming
├── frontend/
│   ├── package.json
│   ├── vite.config.js
│   ├── src/
│   │   ├── App.svelte
│   │   ├── main.js
│   │   ├── lib/
│   │   │   ├── api.js          # fetch wrapper for backend API
│   │   │   ├── stores.js       # Svelte stores (auth, containers)
│   │   │   └── ws.js           # WebSocket client
│   │   ├── pages/
│   │   │   ├── Login.svelte
│   │   │   ├── Dashboard.svelte
│   │   │   ├── Containers.svelte
│   │   │   ├── ContainerDetail.svelte
│   │   │   ├── Images.svelte
│   │   │   ├── Templates.svelte
│   │   │   └── Settings.svelte
│   │   └── components/
│   │       ├── Navbar.svelte
│   │       ├── ContainerCard.svelte
│   │       ├── LogViewer.svelte
│   │       ├── StatsChart.svelte
│   │       └── DeployModal.svelte
│   └── static/
│       └── favicon.png
├── templates/                  # built-in service templates
│   ├── plex.yaml
│   ├── nextcloud.yaml
│   ├── pihole.yaml
│   └── ...
├── data/                       # persistent volume (gitignored)
│   └── data.db
├── Dockerfile
├── docker-compose.yml
├── ARCHITECTURE.md
├── TASKLIST.md
└── DEVLOG.md
```

## API Endpoints

### Auth
```
POST   /api/auth/login          # login, returns session cookie
POST   /api/auth/logout         # destroy session
GET    /api/auth/me             # current user info
```

### Containers
```
GET    /api/containers          # list all containers
GET    /api/containers/:id      # container inspect
POST   /api/containers/:id/start
POST   /api/containers/:id/stop
POST   /api/containers/:id/restart
DELETE /api/containers/:id      # remove container
GET    /api/containers/:id/logs # WebSocket — stream logs
GET    /api/containers/:id/stats # WebSocket — stream CPU/mem
```

### Images
```
GET    /api/images              # list images
DELETE /api/images/:id          # remove image
POST   /api/images/pull         # pull image by name:tag
```

### Services / Templates
```
GET    /api/templates           # list available service templates
GET    /api/templates/:id       # template details
POST   /api/services/deploy     # deploy a service from template or custom compose
```

### System
```
GET    /api/system/info         # host CPU, RAM, disk usage
GET    /api/system/docker       # Docker engine info/version
```

## Database Schema (SQLite)

```sql
CREATE TABLE users (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    username   TEXT    UNIQUE NOT NULL,
    password   TEXT    NOT NULL,  -- bcrypt hash
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id         TEXT    PRIMARY KEY,  -- random token
    user_id    INTEGER NOT NULL REFERENCES users(id),
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE templates (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL,
    description TEXT,
    category    TEXT,              -- media, networking, storage, etc.
    compose     TEXT    NOT NULL,  -- docker-compose YAML content
    icon        TEXT,              -- URL or embedded SVG
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE settings (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
```

## Security

- Docker socket is mounted read/write but never exposed to the network
- All API routes behind auth middleware (except `/api/auth/login`)
- Passwords stored as bcrypt hashes
- Sessions expire after 24 hours
- Cloudflare Access as an additional external auth gate
- CORS restricted to dashboard origin only

## Deployment

```yaml
# docker-compose.yml
services:
  dashboard:
    build: .
    ports:
      - "3000:3000"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./data:/app/data
    environment:
      - DASHBOARD_SECRET=${DASHBOARD_SECRET}
    restart: unless-stopped

  cloudflared:
    image: cloudflare/cloudflared:latest
    command: tunnel run
    environment:
      - TUNNEL_TOKEN=${CF_TUNNEL_TOKEN}
    restart: unless-stopped
```

## Build & Run

```bash
# Development
cd backend && go run .          # API on :3000
cd frontend && npm run dev      # Vite dev server on :5173

# Production (Docker)
docker compose up -d --build
```
