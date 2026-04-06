package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rhea/nas-dashboard/docker"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ContainerHandler struct {
	docker *docker.Client
}

func NewContainerHandler(d *docker.Client) *ContainerHandler {
	return &ContainerHandler{docker: d}
}

func (h *ContainerHandler) List(w http.ResponseWriter, r *http.Request) {
	containers, err := h.docker.ListContainers(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, containers)
}

func (h *ContainerHandler) Inspect(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	info, err := h.docker.InspectContainer(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func (h *ContainerHandler) Start(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.docker.StartContainer(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "started"})
}

func (h *ContainerHandler) Stop(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.docker.StopContainer(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}

func (h *ContainerHandler) Restart(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.docker.RestartContainer(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "restarted"})
}

func (h *ContainerHandler) Remove(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.docker.RemoveContainer(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "removed"})
}

// Logs streams container logs over WebSocket.
func (h *ContainerHandler) Logs(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	reader, err := h.docker.ContainerLogs(r.Context(), id)
	if err != nil {
		conn.WriteJSON(map[string]string{"error": err.Error()})
		return
	}
	defer reader.Close()

	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			// Docker log stream has 8-byte header per frame; strip it
			data := buf[:n]
			if len(data) > 8 {
				data = data[8:]
			}
			if writeErr := conn.WriteMessage(websocket.TextMessage, data); writeErr != nil {
				return
			}
		}
		if err != nil {
			if err != io.EOF {
				conn.WriteJSON(map[string]string{"error": err.Error()})
			}
			return
		}
	}
}

// Stats streams container stats over WebSocket.
func (h *ContainerHandler) Stats(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	reader, err := h.docker.ContainerStats(r.Context(), id)
	if err != nil {
		conn.WriteJSON(map[string]string{"error": err.Error()})
		return
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	for {
		var stats map[string]any
		if err := decoder.Decode(&stats); err != nil {
			if err != io.EOF {
				conn.WriteJSON(map[string]string{"error": err.Error()})
			}
			return
		}
		if err := conn.WriteJSON(stats); err != nil {
			return
		}
	}
}
