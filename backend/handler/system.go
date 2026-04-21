package handler

import (
	"net/http"
	"time"

	"github.com/rhea/nas-dashboard/docker"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemHandler struct {
	docker    *docker.Client
	startTime time.Time
}

func NewSystemHandler(d *docker.Client) *SystemHandler {
	return &SystemHandler{docker: d, startTime: time.Now()}
}

func (h *SystemHandler) Info(w http.ResponseWriter, r *http.Request) {
	// Real host memory
	memStat, _ := mem.VirtualMemoryWithContext(r.Context())

	// Real CPU usage (sample over 200ms)
	cpuPcts, _ := cpu.PercentWithContext(r.Context(), 200*time.Millisecond, false)
	cpuPct := 0.0
	if len(cpuPcts) > 0 {
		cpuPct = cpuPcts[0]
	}

	// Disk usage of the root mount
	diskStat, _ := disk.UsageWithContext(r.Context(), "/")

	// Host uptime
	uptimeSecs, _ := host.UptimeWithContext(r.Context())
	hostInfo, _ := host.InfoWithContext(r.Context())

	hostname := ""
	platform := ""
	if hostInfo != nil {
		hostname = hostInfo.Hostname
		platform = hostInfo.Platform + " " + hostInfo.PlatformVersion
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"hostname":          hostname,
		"platform":          platform,
		"uptime_seconds":    uptimeSecs,
		"dashboard_uptime":  time.Since(h.startTime).Round(time.Second).String(),
		"cpu_percent":       cpuPct,
		"mem_total":         memStat.Total,
		"mem_used":          memStat.Used,
		"mem_percent":       memStat.UsedPercent,
		"disk_total":        diskStat.Total,
		"disk_used":         diskStat.Used,
		"disk_percent":      diskStat.UsedPercent,
	})
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
