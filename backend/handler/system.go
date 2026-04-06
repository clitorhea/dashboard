package handler

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/rhea/nas-dashboard/docker"
)

type SystemHandler struct {
	docker    *docker.Client
	startTime time.Time
}

func NewSystemHandler(d *docker.Client) *SystemHandler {
	return &SystemHandler{docker: d, startTime: time.Now()}
}

type systemInfo struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	CPUs     int    `json:"cpus"`
	Uptime   string `json:"uptime"`
}

func (h *SystemHandler) Info(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	info := systemInfo{
		Hostname: hostname,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		CPUs:     runtime.NumCPU(),
		Uptime:   time.Since(h.startTime).Round(time.Second).String(),
	}
	writeJSON(w, http.StatusOK, info)
}

func (h *SystemHandler) DockerInfo(w http.ResponseWriter, r *http.Request) {
	info, err := h.docker.Info(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	version, err := h.docker.ServerVersion(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"containers_running": info.ContainersRunning,
		"containers_stopped": info.ContainersStopped,
		"containers_total":   info.Containers,
		"images":             info.Images,
		"docker_version":     version.Version,
		"api_version":        version.APIVersion,
		"os":                 info.OperatingSystem,
		"architecture":       info.Architecture,
		"memory_total":       info.MemTotal,
		"cpus":               info.NCPU,
	})
}
