package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gorilla/websocket"
	"github.com/rhea/nas-dashboard/docker"
)

var upgrader = websocket.Upgrader{
	// Validate that the WebSocket origin matches the host to prevent CSRF.
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true // same-origin requests have no Origin header
		}
		host := r.Host
		return strings.Contains(origin, host)
	},
}

// --- Notification types ---

type Notification struct {
	ID          string    `json:"id"`
	ContainerID string    `json:"container_id"`
	Name        string    `json:"name"`
	Message     string    `json:"message"`
	Level       string    `json:"level"` // "error" | "warning" | "info"
	CreatedAt   time.Time `json:"created_at"`
}

// notificationStore holds in-memory notifications and the last-known container states.
var notificationStore = &notifStore{
	notifications: make(map[string]Notification),
	lastState:     make(map[string]string),
}

type notifStore struct {
	mu            sync.RWMutex
	notifications map[string]Notification // id → notification
	lastState     map[string]string       // containerID → last known state
}

func (s *notifStore) upsertState(id, name, state string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	prev, seen := s.lastState[id]
	s.lastState[id] = state

	if seen && prev == "running" && state != "running" {
		// Container transitioned from running → something else
		notifID := id[:12] + "-" + time.Now().Format("150405")
		s.notifications[notifID] = Notification{
			ID:          notifID,
			ContainerID: id,
			Name:        name,
			Message:     name + " stopped unexpectedly (was running, now " + state + ")",
			Level:       "error",
			CreatedAt:   time.Now(),
		}
	}
}

func (s *notifStore) list() []Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Notification, 0, len(s.notifications))
	for _, n := range s.notifications {
		result = append(result, n)
	}
	return result
}

func (s *notifStore) dismiss(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.notifications[id]; !ok {
		return false
	}
	delete(s.notifications, id)
	return true
}

// --- Handler ---

type ContainerHandler struct {
	docker *docker.Client
}

func NewContainerHandler(d *docker.Client) *ContainerHandler {
	h := &ContainerHandler{docker: d}
	// Start background state polling for notifications (every 30s)
	go h.pollContainerStates()
	return h
}

func (h *ContainerHandler) pollContainerStates() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Run immediately on startup too
	h.snapshotStates()

	for range ticker.C {
		h.snapshotStates()
	}
}

func (h *ContainerHandler) snapshotStates() {
	containers, err := h.docker.ListContainers(nil)
	if err != nil {
		return
	}
	for _, c := range containers {
		name := c.ID[:12]
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}
		notificationStore.upsertState(c.ID, name, c.State)
	}
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
	// Clean up from notification state
	notificationStore.mu.Lock()
	delete(notificationStore.lastState, id)
	notificationStore.mu.Unlock()

	writeJSON(w, http.StatusOK, map[string]string{"status": "removed"})
}

// Logs streams container logs over WebSocket using stdcopy to correctly demux stdout/stderr.
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

	// Use a pipe so stdcopy can write to it while we read and forward over WS
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		// stdcopy correctly strips the 8-byte Docker multiplexed header
		stdcopy.StdCopy(pw, pw, reader)
	}()

	buf := make([]byte, 4096)
	for {
		n, err := pr.Read(buf)
		if n > 0 {
			if writeErr := conn.WriteMessage(websocket.TextMessage, buf[:n]); writeErr != nil {
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

// Notifications returns the current list of container health notifications.
func (h *ContainerHandler) Notifications(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, notificationStore.list())
}

// DismissNotification removes a notification by ID.
func (h *ContainerHandler) DismissNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !notificationStore.dismiss(id) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "notification not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "dismissed"})
}
