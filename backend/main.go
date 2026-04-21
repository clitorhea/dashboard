package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rhea/nas-dashboard/db"
	"github.com/rhea/nas-dashboard/docker"
	"github.com/rhea/nas-dashboard/handler"
	"github.com/rhea/nas-dashboard/middleware"
)

func main() {
	// Initialize database
	database, err := db.Init("./data/data.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Seed built-in templates from ./templates/ directory
	if err := seedTemplatesFromDir(database, "./templates"); err != nil {
		log.Printf("Warning: template seeding failed: %v", err)
	}

	// Start background session cleanup goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if n, err := database.CleanupExpiredSessions(); err != nil {
				log.Printf("Session cleanup error: %v", err)
			} else if n > 0 {
				log.Printf("Cleaned up %d expired sessions", n)
			}
		}
	}()

	// Initialize Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to Docker: %v", err)
	}

	// Create handlers
	authHandler := handler.NewAuthHandler(database)
	containerHandler := handler.NewContainerHandler(dockerClient)
	imageHandler := handler.NewImageHandler(dockerClient)
	systemHandler := handler.NewSystemHandler(dockerClient)
	templateHandler := handler.NewTemplateHandler(database, dockerClient)

	// Auth middleware
	authMW := middleware.NewAuthMiddleware(database)

	mux := http.NewServeMux()

	// Auth routes (no middleware)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)
	mux.HandleFunc("GET /api/auth/me", authMW.Require(authHandler.Me))
	mux.HandleFunc("POST /api/auth/setup", authHandler.Setup)
	mux.HandleFunc("PUT /api/auth/password", authMW.Require(authHandler.ChangePassword))

	// Container routes
	mux.HandleFunc("GET /api/containers", authMW.Require(containerHandler.List))
	mux.HandleFunc("GET /api/containers/{id}", authMW.Require(containerHandler.Inspect))
	mux.HandleFunc("POST /api/containers/{id}/start", authMW.Require(containerHandler.Start))
	mux.HandleFunc("POST /api/containers/{id}/stop", authMW.Require(containerHandler.Stop))
	mux.HandleFunc("POST /api/containers/{id}/restart", authMW.Require(containerHandler.Restart))
	mux.HandleFunc("DELETE /api/containers/{id}", authMW.Require(containerHandler.Remove))
	mux.HandleFunc("GET /api/containers/{id}/logs", authMW.Require(containerHandler.Logs))
	mux.HandleFunc("GET /api/containers/{id}/stats", authMW.Require(containerHandler.Stats))

	// Image routes
	mux.HandleFunc("GET /api/images", authMW.Require(imageHandler.List))
	mux.HandleFunc("DELETE /api/images/{id}", authMW.Require(imageHandler.Remove))
	mux.HandleFunc("POST /api/images/pull", authMW.Require(imageHandler.Pull))

	// System routes
	mux.HandleFunc("GET /api/system/info", authMW.Require(systemHandler.Info))
	mux.HandleFunc("GET /api/system/docker", authMW.Require(systemHandler.DockerInfo))

	// Template routes
	mux.HandleFunc("GET /api/templates", authMW.Require(templateHandler.List))
	mux.HandleFunc("GET /api/templates/{id}", authMW.Require(templateHandler.Get))
	mux.HandleFunc("POST /api/services/deploy", authMW.Require(templateHandler.Deploy))

	// Notifications route
	mux.HandleFunc("GET /api/notifications", authMW.Require(containerHandler.Notifications))
	mux.HandleFunc("POST /api/notifications/{id}/dismiss", authMW.Require(containerHandler.DismissNotification))

	// Serve frontend static files with SPA fallback
	staticFS := http.Dir("./static")
	fileServer := http.FileServer(staticFS)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Let API requests fall through (should never reach here due to prefix matching)
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		// Try to serve the file; if not found, serve index.html (SPA fallback)
		path := filepath.Join("./static", r.URL.Path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.ServeFile(w, r, "./static/index.html")
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Dashboard starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// seedTemplatesFromDir reads all .yaml files in dir and seeds them into the DB.
func seedTemplatesFromDir(database *db.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // templates dir is optional
		}
		return err
	}

	var templates []db.Template
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			log.Printf("Warning: could not read template %s: %v", name, err)
			continue
		}

		meta := templateMetaFromFilename(name)
		templates = append(templates, db.Template{
			Name:        meta.name,
			Description: meta.description,
			Category:    meta.category,
			Compose:     string(content),
			Icon:        meta.icon,
		})
	}

	if len(templates) == 0 {
		return nil
	}

	return database.SeedTemplates(templates)
}

type templateMeta struct {
	name        string
	description string
	category    string
	icon        string
}

func templateMetaFromFilename(filename string) templateMeta {
	base := strings.TrimSuffix(strings.TrimSuffix(filename, ".yaml"), ".yml")
	known := map[string]templateMeta{
		"plex":              {name: "Plex Media Server", description: "Stream movies and TV shows to any device.", category: "media", icon: "🎬"},
		"nextcloud":         {name: "Nextcloud", description: "Self-hosted cloud storage and collaboration suite.", category: "storage", icon: "☁"},
		"pihole":            {name: "Pi-hole", description: "Network-wide ad blocker and DNS server.", category: "networking", icon: "🛡"},
		"jellyfin":          {name: "Jellyfin", description: "Free software media server — stream your collection.", category: "media", icon: "📺"},
		"vaultwarden":       {name: "Vaultwarden", description: "Lightweight Bitwarden-compatible password manager.", category: "security", icon: "🔐"},
		"nginx-proxy-manager": {name: "Nginx Proxy Manager", description: "Manage reverse proxies and SSL certs via a web UI.", category: "networking", icon: "🔀"},
		"immich":            {name: "Immich", description: "Self-hosted photo and video backup.", category: "media", icon: "🖼"},
		"gitea":             {name: "Gitea", description: "Lightweight self-hosted Git service.", category: "development", icon: "🐙"},
		"uptime-kuma":       {name: "Uptime Kuma", description: "Self-hosted monitoring tool for services.", category: "monitoring", icon: "📊"},
	}
	if m, ok := known[base]; ok {
		return m
	}
	return templateMeta{
		name:        base,
		description: base + " service",
		category:    "other",
		icon:        "📦",
	}
}
