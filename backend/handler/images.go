package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rhea/nas-dashboard/docker"
)

type ImageHandler struct {
	docker *docker.Client
}

func NewImageHandler(d *docker.Client) *ImageHandler {
	return &ImageHandler{docker: d}
}

func (h *ImageHandler) List(w http.ResponseWriter, r *http.Request) {
	images, err := h.docker.ListImages(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, images)
}

func (h *ImageHandler) Remove(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.docker.RemoveImage(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "removed"})
}

type pullRequest struct {
	Image string `json:"image"`
}

func (h *ImageHandler) Pull(w http.ResponseWriter, r *http.Request) {
	var req pullRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	reader, err := h.docker.PullImage(r.Context(), req.Image)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer reader.Close()

	// Stream pull progress to client
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.WriteHeader(http.StatusOK)
	flusher, _ := w.(http.Flusher)

	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
			if flusher != nil {
				flusher.Flush()
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
	}
}
