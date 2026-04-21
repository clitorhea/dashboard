<script>
  import { api } from '../lib/api.js';

  let { container, onRefresh } = $props();

  let actionLoading = $state(false);

  function getName() {
    if (container.Names && container.Names.length > 0) {
      return container.Names[0].replace(/^\//, '');
    }
    return container.Id.substring(0, 12);
  }

  function getImage() {
    return container.Image || 'unknown';
  }

  function getStatus() {
    return container.State || 'unknown';
  }

  function isRunning() {
    return container.State === 'running';
  }

  async function doAction(action) {
    if (action === 'stop' && !confirm(`Stop "${getName()}"?`)) return;
    if (action === 'restart' && !confirm(`Restart "${getName()}"?`)) return;

    actionLoading = true;
    try {
      if (action === 'start') await api.startContainer(container.Id);
      else if (action === 'stop') await api.stopContainer(container.Id);
      else if (action === 'restart') await api.restartContainer(container.Id);
      if (onRefresh) onRefresh();
    } catch (e) {
      alert(`Action failed: ${e.message}`);
    } finally {
      actionLoading = false;
    }
  }
</script>

<div class="card" class:running={isRunning()} class:stopped={!isRunning()}>
  <div class="card-header">
    <div class="status-dot" class:green={isRunning()} class:red={!isRunning()}></div>
    <a href="#/containers/{container.Id}" class="name">{getName()}</a>
  </div>

  <div class="card-body">
    <div class="info-row">
      <span class="label">Image</span>
      <span class="value">{getImage()}</span>
    </div>
    <div class="info-row">
      <span class="label">Status</span>
      <span class="value">{container.Status || getStatus()}</span>
    </div>
  </div>

  <div class="card-actions">
    {#if isRunning()}
      <button class="btn stop" onclick={() => doAction('stop')} disabled={actionLoading}>Stop</button>
      <button class="btn restart" onclick={() => doAction('restart')} disabled={actionLoading}>Restart</button>
    {:else}
      <button class="btn start" onclick={() => doAction('start')} disabled={actionLoading}>Start</button>
    {/if}
  </div>
</div>

<style>
  .card {
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    padding: 1rem;
    transition: border-color 0.15s;
  }

  .card:hover {
    border-color: #58a6ff;
  }

  .card-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.75rem;
  }

  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
  }

  .status-dot.green { background: #3fb950; }
  .status-dot.red { background: #f85149; }

  .name {
    font-weight: 600;
    color: #58a6ff;
    text-decoration: none;
    font-size: 0.95rem;
  }

  .name:hover { text-decoration: underline; }

  .card-body {
    margin-bottom: 0.75rem;
  }

  .info-row {
    display: flex;
    justify-content: space-between;
    font-size: 0.8rem;
    padding: 0.2rem 0;
  }

  .label { color: #8b949e; }
  .value {
    color: #e1e4e8;
    max-width: 60%;
    text-align: right;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .card-actions {
    display: flex;
    gap: 0.5rem;
  }

  .btn {
    flex: 1;
    padding: 0.4rem;
    border: 1px solid #30363d;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.8rem;
    background: #21262d;
    color: #e1e4e8;
    transition: all 0.15s;
  }

  .btn:hover:not(:disabled) { background: #30363d; }
  .btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .btn.start { border-color: #238636; color: #3fb950; }
  .btn.stop { border-color: #da3633; color: #f85149; }
  .btn.restart { border-color: #9e6a03; color: #d29922; }
</style>
