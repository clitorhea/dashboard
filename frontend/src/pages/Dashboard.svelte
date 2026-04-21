<script>
  import { onMount, onDestroy } from 'svelte';
  import { api } from '../lib/api.js';
  import { notifications } from '../lib/stores.js';

  let dockerInfo = $state(null);
  let systemInfo = $state(null);
  let error = $state('');
  let interval;

  async function load() {
    try {
      const [docker, system] = await Promise.all([
        api.dockerInfo(),
        api.systemInfo(),
      ]);
      dockerInfo = docker;
      systemInfo = system;
    } catch (e) {
      error = e.message;
    }
  }

  async function loadNotifications() {
    try {
      const notifs = await api.notifications();
      notifications.set(notifs);
    } catch (_) {}
  }

  onMount(() => {
    load();
    loadNotifications();
    // Refresh stats every 10s
    interval = setInterval(() => {
      load();
      loadNotifications();
    }, 10_000);
  });

  onDestroy(() => clearInterval(interval));

  async function dismissNotif(id) {
    try {
      await api.dismissNotification(id);
      notifications.update(n => n.filter(x => x.id !== id));
    } catch (_) {}
  }

  function formatBytes(bytes) {
    if (!bytes) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }

  function formatUptime(seconds) {
    if (!seconds) return '-';
    const d = Math.floor(seconds / 86400);
    const h = Math.floor((seconds % 86400) / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    if (d > 0) return `${d}d ${h}h`;
    if (h > 0) return `${h}h ${m}m`;
    return `${m}m`;
  }
</script>

<div class="dashboard">
  <h1>Dashboard</h1>

  {#if $notifications.length > 0}
    <div class="alerts">
      {#each $notifications as notif (notif.id)}
        <div class="alert alert-{notif.level}">
          <span class="alert-icon">{notif.level === 'error' ? '⚠' : 'ℹ'}</span>
          <span class="alert-msg">{notif.message}</span>
          <button class="alert-dismiss" onclick={() => dismissNotif(notif.id)} aria-label="Dismiss">✕</button>
        </div>
      {/each}
    </div>
  {/if}

  {#if error}
    <div class="error">{error}</div>
  {:else if !dockerInfo}
    <p class="loading">Loading…</p>
  {:else}
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-number running">{dockerInfo.containers_running}</div>
        <div class="stat-label">Running</div>
      </div>
      <div class="stat-card">
        <div class="stat-number stopped">{dockerInfo.containers_stopped}</div>
        <div class="stat-label">Stopped</div>
      </div>
      <div class="stat-card">
        <div class="stat-number total">{dockerInfo.containers_total}</div>
        <div class="stat-label">Total Containers</div>
      </div>
      <div class="stat-card">
        <div class="stat-number images">{dockerInfo.images}</div>
        <div class="stat-label">Images</div>
      </div>
    </div>

    <!-- Real host resource bars -->
    {#if systemInfo}
      <div class="resource-section">
        <h2>Host Resources</h2>
        <div class="resource-grid">
          <div class="resource-item">
            <div class="resource-header">
              <span>CPU</span>
              <span class="resource-value">{systemInfo.cpu_percent?.toFixed(1)}%</span>
            </div>
            <div class="bar-track">
              <div class="bar-fill cpu" style="width: {Math.min(systemInfo.cpu_percent ?? 0, 100)}%"></div>
            </div>
          </div>
          <div class="resource-item">
            <div class="resource-header">
              <span>Memory</span>
              <span class="resource-value">{formatBytes(systemInfo.mem_used)} / {formatBytes(systemInfo.mem_total)}</span>
            </div>
            <div class="bar-track">
              <div class="bar-fill mem" style="width: {Math.min(systemInfo.mem_percent ?? 0, 100)}%"></div>
            </div>
          </div>
          <div class="resource-item">
            <div class="resource-header">
              <span>Disk (/)</span>
              <span class="resource-value">{formatBytes(systemInfo.disk_used)} / {formatBytes(systemInfo.disk_total)}</span>
            </div>
            <div class="bar-track">
              <div class="bar-fill disk" style="width: {Math.min(systemInfo.disk_percent ?? 0, 100)}%"></div>
            </div>
          </div>
        </div>
      </div>
    {/if}

    <div class="info-section">
      <h2>System</h2>
      <div class="info-grid">
        <div class="info-item">
          <span class="info-label">Hostname</span>
          <span class="info-value">{systemInfo?.hostname || '-'}</span>
        </div>
        <div class="info-item">
          <span class="info-label">OS</span>
          <span class="info-value">{dockerInfo.os}</span>
        </div>
        <div class="info-item">
          <span class="info-label">Architecture</span>
          <span class="info-value">{dockerInfo.architecture}</span>
        </div>
        <div class="info-item">
          <span class="info-label">CPUs</span>
          <span class="info-value">{dockerInfo.cpus}</span>
        </div>
        <div class="info-item">
          <span class="info-label">Host Uptime</span>
          <span class="info-value">{formatUptime(systemInfo?.uptime_seconds)}</span>
        </div>
        <div class="info-item">
          <span class="info-label">Docker Version</span>
          <span class="info-value">{dockerInfo.docker_version}</span>
        </div>
        <div class="info-item">
          <span class="info-label">API Version</span>
          <span class="info-value">{dockerInfo.api_version}</span>
        </div>
        <div class="info-item">
          <span class="info-label">Dashboard Uptime</span>
          <span class="info-value">{systemInfo?.dashboard_uptime || '-'}</span>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .dashboard { max-width: 1000px; }

  h1 { margin-bottom: 1.5rem; font-size: 1.5rem; }
  h2 { margin-bottom: 1rem; font-size: 1.1rem; color: #8b949e; }

  /* Notifications */
  .alerts { margin-bottom: 1.5rem; display: flex; flex-direction: column; gap: 0.5rem; }

  .alert {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    border-radius: 6px;
    font-size: 0.875rem;
  }

  .alert-error { background: #1c0d0d; border: 1px solid #da3633; color: #f85149; }
  .alert-warning { background: #1a150a; border: 1px solid #9e6a03; color: #d29922; }
  .alert-info { background: #0c1a2a; border: 1px solid #1f6feb; color: #58a6ff; }

  .alert-icon { font-size: 1rem; }
  .alert-msg { flex: 1; }
  .alert-dismiss {
    background: none;
    border: none;
    color: inherit;
    cursor: pointer;
    opacity: 0.6;
    font-size: 0.9rem;
    padding: 0.1rem 0.3rem;
  }
  .alert-dismiss:hover { opacity: 1; }

  /* Container stat cards */
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  .stat-card {
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    padding: 1.25rem;
    text-align: center;
  }

  .stat-number {
    font-size: 2rem;
    font-weight: 700;
    margin-bottom: 0.25rem;
  }

  .stat-number.running { color: #3fb950; }
  .stat-number.stopped { color: #f85149; }
  .stat-number.total { color: #58a6ff; }
  .stat-number.images { color: #d29922; }

  .stat-label {
    font-size: 0.8rem;
    color: #8b949e;
    text-transform: uppercase;
  }

  /* Resource bars */
  .resource-section {
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    padding: 1.5rem;
    margin-bottom: 1.5rem;
  }

  .resource-grid {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .resource-item { display: flex; flex-direction: column; gap: 0.4rem; }

  .resource-header {
    display: flex;
    justify-content: space-between;
    font-size: 0.85rem;
    color: #8b949e;
  }

  .resource-value { color: #e1e4e8; font-variant-numeric: tabular-nums; }

  .bar-track {
    width: 100%;
    height: 6px;
    background: #21262d;
    border-radius: 3px;
    overflow: hidden;
  }

  .bar-fill {
    height: 100%;
    border-radius: 3px;
    transition: width 0.5s ease;
  }

  .bar-fill.cpu { background: #58a6ff; }
  .bar-fill.mem { background: #3fb950; }
  .bar-fill.disk { background: #d29922; }

  /* System info */
  .info-section {
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    padding: 1.5rem;
  }

  .info-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 0.75rem;
  }

  .info-item {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
  }

  .info-label { font-size: 0.75rem; color: #8b949e; text-transform: uppercase; }
  .info-value { font-size: 0.95rem; }

  .loading { color: #8b949e; }
  .error { color: #f85149; padding: 1rem; background: #1c0d0d; border-radius: 6px; }
</style>
