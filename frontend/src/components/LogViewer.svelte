<script>
  import { onMount, onDestroy } from 'svelte';
  import { connectWebSocket } from '../lib/ws.js';

  let { containerId } = $props();

  let logs = $state([]);
  let logContainer;
  let ws;

  onMount(() => {
    ws = connectWebSocket(`/api/containers/${containerId}/logs`, (data) => {
      logs = [...logs, typeof data === 'string' ? data : JSON.stringify(data)];
      // Auto-scroll
      if (logContainer) {
        requestAnimationFrame(() => {
          logContainer.scrollTop = logContainer.scrollHeight;
        });
      }
    });
  });

  onDestroy(() => {
    if (ws) ws.close();
  });
</script>

<div class="log-viewer" bind:this={logContainer}>
  {#if logs.length === 0}
    <p class="empty">Waiting for logs...</p>
  {:else}
    {#each logs as line}
      <pre class="log-line">{line}</pre>
    {/each}
  {/if}
</div>

<style>
  .log-viewer {
    background: #0d1117;
    border: 1px solid #30363d;
    border-radius: 6px;
    padding: 1rem;
    max-height: 500px;
    overflow-y: auto;
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    font-size: 0.8rem;
  }

  .log-line {
    margin: 0;
    padding: 1px 0;
    white-space: pre-wrap;
    word-break: break-all;
    color: #c9d1d9;
  }

  .empty {
    color: #8b949e;
    text-align: center;
    padding: 2rem;
  }
</style>
