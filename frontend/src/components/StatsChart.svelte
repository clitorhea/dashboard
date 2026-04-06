<script>
  import { onMount, onDestroy } from 'svelte';
  import { connectWebSocket } from '../lib/ws.js';

  let { containerId } = $props();

  let cpuPercent = $state(0);
  let memUsage = $state(0);
  let memLimit = $state(0);
  let netIn = $state(0);
  let netOut = $state(0);
  let ws;

  function formatBytes(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }

  function calcCPU(stats) {
    const cpuDelta = stats.cpu_stats.cpu_usage.total_usage - stats.precpu_stats.cpu_usage.total_usage;
    const systemDelta = stats.cpu_stats.system_cpu_usage - stats.precpu_stats.system_cpu_usage;
    const cpuCount = stats.cpu_stats.online_cpus || 1;
    if (systemDelta > 0) {
      return ((cpuDelta / systemDelta) * cpuCount * 100).toFixed(1);
    }
    return 0;
  }

  onMount(() => {
    ws = connectWebSocket(`/api/containers/${containerId}/stats`, (stats) => {
      if (stats.cpu_stats) {
        cpuPercent = calcCPU(stats);
      }
      if (stats.memory_stats) {
        memUsage = stats.memory_stats.usage || 0;
        memLimit = stats.memory_stats.limit || 0;
      }
      if (stats.networks) {
        let totalIn = 0, totalOut = 0;
        for (const iface of Object.values(stats.networks)) {
          totalIn += iface.rx_bytes || 0;
          totalOut += iface.tx_bytes || 0;
        }
        netIn = totalIn;
        netOut = totalOut;
      }
    });
  });

  onDestroy(() => {
    if (ws) ws.close();
  });
</script>

<div class="stats-grid">
  <div class="stat-card">
    <div class="stat-label">CPU</div>
    <div class="stat-value">{cpuPercent}%</div>
    <div class="stat-bar">
      <div class="stat-fill cpu" style="width: {Math.min(cpuPercent, 100)}%"></div>
    </div>
  </div>

  <div class="stat-card">
    <div class="stat-label">Memory</div>
    <div class="stat-value">{formatBytes(memUsage)} / {formatBytes(memLimit)}</div>
    <div class="stat-bar">
      <div class="stat-fill mem" style="width: {memLimit > 0 ? Math.min((memUsage / memLimit) * 100, 100) : 0}%"></div>
    </div>
  </div>

  <div class="stat-card">
    <div class="stat-label">Network I/O</div>
    <div class="stat-value">{formatBytes(netIn)} in / {formatBytes(netOut)} out</div>
  </div>
</div>

<style>
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
  }

  .stat-card {
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 6px;
    padding: 1rem;
  }

  .stat-label {
    font-size: 0.75rem;
    color: #8b949e;
    text-transform: uppercase;
    margin-bottom: 0.25rem;
  }

  .stat-value {
    font-size: 1.1rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
  }

  .stat-bar {
    height: 4px;
    background: #30363d;
    border-radius: 2px;
    overflow: hidden;
  }

  .stat-fill {
    height: 100%;
    border-radius: 2px;
    transition: width 0.5s ease;
  }

  .stat-fill.cpu { background: #58a6ff; }
  .stat-fill.mem { background: #3fb950; }
</style>
