<script>
  import { onMount } from 'svelte';
  import { api } from '../lib/api.js';
  import LogViewer from '../components/LogViewer.svelte';
  import StatsChart from '../components/StatsChart.svelte';

  let { params } = $props();

  let container = $state(null);
  let error = $state('');
  let activeTab = $state('stats');

  onMount(async () => {
    try {
      container = await api.container(params.id);
    } catch (e) {
      error = e.message;
    }
  });

  function getName() {
    return container?.Name?.replace(/^\//, '') || params.id.substring(0, 12);
  }

  function isRunning() {
    return container?.State?.Running;
  }

  function getPorts() {
    const ports = container?.NetworkSettings?.Ports;
    if (!ports) return [];
    const result = [];
    for (const [containerPort, bindings] of Object.entries(ports)) {
      if (bindings) {
        for (const b of bindings) {
          result.push(`${b.HostPort} -> ${containerPort}`);
        }
      }
    }
    return result;
  }

  function getEnv() {
    return container?.Config?.Env || [];
  }

  function getVolumes() {
    return container?.Mounts || [];
  }

  async function doAction(action) {
    try {
      if (action === 'start') await api.startContainer(params.id);
      else if (action === 'stop') await api.stopContainer(params.id);
      else if (action === 'restart') await api.restartContainer(params.id);
      container = await api.container(params.id);
    } catch (e) {
      alert(e.message);
    }
  }

  async function remove() {
    if (!confirm(`Remove container ${getName()}? This cannot be undone.`)) return;
    try {
      await api.removeContainer(params.id);
      window.location.hash = '#/containers';
    } catch (e) {
      alert(e.message);
    }
  }
</script>

<div class="detail-page">
  {#if error}
    <div class="error">{error}</div>
  {:else if !container}
    <p class="loading">Loading...</p>
  {:else}
    <div class="header">
      <div class="title">
        <div class="status-dot" class:green={isRunning()} class:red={!isRunning()}></div>
        <h1>{getName()}</h1>
      </div>
      <div class="actions">
        {#if isRunning()}
          <button class="btn stop" onclick={() => doAction('stop')}>Stop</button>
          <button class="btn restart" onclick={() => doAction('restart')}>Restart</button>
        {:else}
          <button class="btn start" onclick={() => doAction('start')}>Start</button>
        {/if}
        <button class="btn remove" onclick={remove}>Remove</button>
      </div>
    </div>

    <div class="info-grid">
      <div class="info-item">
        <span class="label">Image</span>
        <span class="value">{container.Config?.Image}</span>
      </div>
      <div class="info-item">
        <span class="label">Status</span>
        <span class="value">{container.State?.Status}</span>
      </div>
      <div class="info-item">
        <span class="label">Created</span>
        <span class="value">{new Date(container.Created).toLocaleString()}</span>
      </div>
      <div class="info-item">
        <span class="label">ID</span>
        <span class="value mono">{container.Id?.substring(0, 12)}</span>
      </div>
    </div>

    {#if getPorts().length > 0}
      <div class="section">
        <h3>Ports</h3>
        <div class="tags">
          {#each getPorts() as port}
            <span class="tag">{port}</span>
          {/each}
        </div>
      </div>
    {/if}

    <div class="tabs">
      <button class:active={activeTab === 'stats'} onclick={() => activeTab = 'stats'}>Stats</button>
      <button class:active={activeTab === 'logs'} onclick={() => activeTab = 'logs'}>Logs</button>
      <button class:active={activeTab === 'env'} onclick={() => activeTab = 'env'}>Environment</button>
      <button class:active={activeTab === 'volumes'} onclick={() => activeTab = 'volumes'}>Volumes</button>
    </div>

    <div class="tab-content">
      {#if activeTab === 'stats' && isRunning()}
        <StatsChart containerId={params.id} />
      {:else if activeTab === 'stats'}
        <p class="empty">Container is not running</p>
      {:else if activeTab === 'logs'}
        <LogViewer containerId={params.id} />
      {:else if activeTab === 'env'}
        <div class="env-list">
          {#each getEnv() as env}
            <pre class="env-line">{env}</pre>
          {/each}
        </div>
      {:else if activeTab === 'volumes'}
        <div class="volume-list">
          {#each getVolumes() as mount}
            <div class="volume-item">
              <span class="mono">{mount.Source}</span>
              <span class="arrow">-></span>
              <span class="mono">{mount.Destination}</span>
              <span class="tag">{mount.Mode || 'rw'}</span>
            </div>
          {/each}
          {#if getVolumes().length === 0}
            <p class="empty">No volumes mounted</p>
          {/if}
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .detail-page { max-width: 1000px; }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    flex-wrap: wrap;
    gap: 1rem;
  }

  .title {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  h1 { font-size: 1.5rem; }

  .status-dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
  }

  .status-dot.green { background: #3fb950; }
  .status-dot.red { background: #f85149; }

  .actions { display: flex; gap: 0.5rem; }

  .btn {
    padding: 0.5rem 1rem;
    border: 1px solid #30363d;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.85rem;
    background: #21262d;
    color: #e1e4e8;
  }

  .btn.start { border-color: #238636; color: #3fb950; }
  .btn.stop { border-color: #da3633; color: #f85149; }
  .btn.restart { border-color: #9e6a03; color: #d29922; }
  .btn.remove { border-color: #da3633; color: #f85149; }
  .btn:hover { background: #30363d; }

  .info-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 0.75rem;
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    padding: 1.25rem;
    margin-bottom: 1.5rem;
  }

  .info-item { display: flex; flex-direction: column; gap: 0.2rem; }
  .label { font-size: 0.75rem; color: #8b949e; text-transform: uppercase; }
  .value { font-size: 0.9rem; }
  .mono { font-family: monospace; font-size: 0.85rem; }

  .section { margin-bottom: 1.5rem; }
  .section h3 { font-size: 0.9rem; color: #8b949e; margin-bottom: 0.5rem; }

  .tags { display: flex; gap: 0.5rem; flex-wrap: wrap; }
  .tag {
    background: #21262d;
    border: 1px solid #30363d;
    border-radius: 4px;
    padding: 0.2rem 0.6rem;
    font-size: 0.8rem;
    font-family: monospace;
  }

  .tabs {
    display: flex;
    gap: 0;
    border-bottom: 1px solid #30363d;
    margin-bottom: 1rem;
  }

  .tabs button {
    background: none;
    border: none;
    color: #8b949e;
    padding: 0.75rem 1rem;
    cursor: pointer;
    font-size: 0.85rem;
    border-bottom: 2px solid transparent;
  }

  .tabs button.active {
    color: #58a6ff;
    border-bottom-color: #58a6ff;
  }

  .env-list {
    background: #0d1117;
    border: 1px solid #30363d;
    border-radius: 6px;
    padding: 1rem;
    max-height: 400px;
    overflow-y: auto;
  }

  .env-line {
    margin: 0;
    padding: 2px 0;
    font-size: 0.8rem;
    white-space: pre-wrap;
    word-break: break-all;
  }

  .volume-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 0;
    border-bottom: 1px solid #21262d;
    font-size: 0.85rem;
    flex-wrap: wrap;
  }

  .arrow { color: #8b949e; }
  .empty { color: #8b949e; padding: 2rem; text-align: center; }
  .loading { color: #8b949e; }
  .error { color: #f85149; padding: 1rem; background: #1c0d0d; border-radius: 6px; }
</style>
