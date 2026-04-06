package main

import (
	"log"
	"net/http"
	"os"

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

	// Serve frontend static files
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Dashboard starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
