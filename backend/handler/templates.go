package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/rhea/nas-dashboard/db"
	"github.com/rhea/nas-dashboard/docker"
)

type TemplateHandler struct {
	db     *db.DB
	docker *docker.Client
}

func NewTemplateHandler(database *db.DB, d *docker.Client) *TemplateHandler {
	return &TemplateHandler{db: database, docker: d}
}

func (h *TemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	templates, err := h.db.ListTemplates()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if templates == nil {
		templates = []db.Template{}
	}
	writeJSON(w, http.StatusOK, templates)
}

func (h *TemplateHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	tmpl, err := h.db.GetTemplate(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "template not found"})
		return
	}
	writeJSON(w, http.StatusOK, tmpl)
}

type deployRequest struct {
	TemplateID  int64             `json:"template_id,omitempty"`
	Compose     string            `json:"compose,omitempty"`
	ServiceName string            `json:"service_name"`
	Env         map[string]string `json:"env,omitempty"`
}

// Deploy creates a docker-compose file and runs docker compose up.
func (h *TemplateHandler) Deploy(w http.ResponseWriter, r *http.Request) {
	var req deployRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	if req.ServiceName == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "service_name required"})
		return
	}

	composeContent := req.Compose
	if req.TemplateID > 0 {
		tmpl, err := h.db.GetTemplate(req.TemplateID)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "template not found"})
			return
		}
		composeContent = tmpl.Compose
	}

	if composeContent == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "compose content required"})
		return
	}

	// Create service directory
	serviceDir := filepath.Join("./data/services", req.ServiceName)
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create service directory"})
		return
	}

	// Write compose file
	composePath := filepath.Join(serviceDir, "docker-compose.yml")
	if err := os.WriteFile(composePath, []byte(composeContent), 0644); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to write compose file"})
		return
	}

	// Write .env file if env vars provided
	if len(req.Env) > 0 {
		envContent := ""
		for k, v := range req.Env {
			envContent += fmt.Sprintf("%s=%s\n", k, v)
		}
		envPath := filepath.Join(serviceDir, ".env")
		if err := os.WriteFile(envPath, []byte(envContent), 0600); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to write env file"})
			return
		}
	}

	// Run docker compose up
	cmd := exec.CommandContext(r.Context(), "docker", "compose", "-f", composePath, "up", "-d")
	cmd.Dir = serviceDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":  "deploy failed",
			"output": string(output),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "deployed",
		"service": req.ServiceName,
		"output":  string(output),
	})
}
