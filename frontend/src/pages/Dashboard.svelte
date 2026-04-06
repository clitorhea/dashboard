<script>
  import { onMount } from 'svelte';
  import { api } from '../lib/api.js';

  let dockerInfo = $state(null);
  let systemInfo = $state(null);
  let error = $state('');

  onMount(async () => {
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
  });

  function formatBytes(bytes) {
    if (!bytes) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }
</script>

<div class="dashboard">
  <h1>Dashboard</h1>

  {#if error}
    <div class="error">{error}</div>
  {:else if !dockerInfo}
    <p class="loading">Loading...</p>
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
          <span class="info-label">Memory</span>
          <span class="info-value">{formatBytes(dockerInfo.memory_total)}</span>
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
          <span class="info-value">{systemInfo?.uptime || '-'}</span>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .dashboard { max-width: 1000px; }

  h1 { margin-bottom: 1.5rem; font-size: 1.5rem; }
  h2 { margin-bottom: 1rem; font-size: 1.1rem; color: #8b949e; }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 1rem;
    margin-bottom: 2rem;
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
