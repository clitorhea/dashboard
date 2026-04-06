# NAS Dashboard

A self-hosted web dashboard for managing Docker services on your NAS. View containers, stream logs, monitor stats, deploy new services, and access it all remotely via Cloudflare Tunnel.

## Features

- **Container management** — start, stop, restart, remove containers from the browser
- **Live monitoring** — real-time log streaming and CPU/memory stats via WebSocket
- **Image management** — list, pull, and remove Docker images
- **Service templates** — deploy pre-configured services (Plex, Nextcloud, Pi-hole) with one click
- **Custom deploy** — paste any docker-compose YAML to deploy custom services
- **Remote access** — optional Cloudflare Tunnel integration for access outside your LAN
- **First-run setup** — creates your admin account on first visit, no config files needed
- **Dark theme** — designed for a clean, minimal look

## Quick Start

### Prerequisites

- Docker and Docker Compose installed on your NAS
- Git (to clone the repo)

### Deploy

```bash
git clone <your-repo-url> nas-dashboard
cd nas-dashboard
docker compose up -d --build
```

Open `http://<your-nas-ip>:3000` in your browser. On first visit you'll be prompted to create an admin account.

### With Cloudflare Tunnel (remote access)

1. Create a tunnel in the [Cloudflare Zero Trust dashboard](https://one.dash.cloudflare.com/)
2. Copy your tunnel token
3. Create a `.env` file:

```bash
CF_TUNNEL_TOKEN=your-token-here
```

4. Deploy with the tunnel profile:

```bash
docker compose --profile tunnel up -d --build
```

Your dashboard is now accessible at the domain you configured in Cloudflare.

## Development

### Prerequisites

- Go 1.23+
- Node.js 22+

### Backend

```bash
cd backend
go mod tidy
go run .
```

The API starts on `:3000`.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Vite dev server starts on `:5173` and proxies `/api` requests to the backend.

## Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md) for the full system design, API reference, and database schema.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Svelte 5 + Vite |
| Backend | Go (net/http + Docker SDK) |
| Database | SQLite (WAL mode) |
| Auth | bcrypt + session cookies |
| Remote Access | Cloudflare Tunnel |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | Dashboard port |
| `CF_TUNNEL_TOKEN` | — | Cloudflare Tunnel token (only needed with tunnel profile) |

## Project Structure

```
├── backend/          Go API server
│   ├── handler/      HTTP handlers (auth, containers, images, system, templates)
│   ├── middleware/    Auth middleware
│   ├── docker/       Docker SDK wrapper
│   └── db/           SQLite database + models
├── frontend/         Svelte SPA
│   └── src/
│       ├── pages/    Route pages
│       ├── components/  Reusable UI components
│       └── lib/      API client, stores, WebSocket helper
├── templates/        Built-in docker-compose templates
├── Dockerfile        Multi-stage build
└── docker-compose.yml
```

## License

MIT
