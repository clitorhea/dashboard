<script>
  import { onMount } from 'svelte';
  import { api } from '../lib/api.js';
  import ContainerCard from '../components/ContainerCard.svelte';

  let containers = $state([]);
  let error = $state('');
  let filter = $state('all');

  async function load() {
    try {
      containers = await api.containers();
    } catch (e) {
      error = e.message;
    }
  }

  onMount(load);

  function filtered() {
    if (filter === 'running') return containers.filter(c => c.State === 'running');
    if (filter === 'stopped') return containers.filter(c => c.State !== 'running');
    return containers;
  }
</script>

<div class="containers-page">
  <div class="header">
    <h1>Containers</h1>
    <div class="filters">
      <button class:active={filter === 'all'} onclick={() => filter = 'all'}>
        All ({containers.length})
      </button>
      <button class:active={filter === 'running'} onclick={() => filter = 'running'}>
        Running ({containers.filter(c => c.State === 'running').length})
      </button>
      <button class:active={filter === 'stopped'} onclick={() => filter = 'stopped'}>
        Stopped ({containers.filter(c => c.State !== 'running').length})
      </button>
    </div>
  </div>

  {#if error}
    <div class="error">{error}</div>
  {:else if containers.length === 0}
    <p class="empty">No containers found</p>
  {:else}
    <div class="grid">
      {#each filtered() as container (container.Id)}
        <ContainerCard {container} onRefresh={load} />
      {/each}
    </div>
  {/if}
</div>

<style>
  .containers-page { max-width: 1200px; }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    flex-wrap: wrap;
    gap: 1rem;
  }

  h1 { font-size: 1.5rem; }

  .filters {
    display: flex;
    gap: 0.5rem;
  }

  .filters button {
    background: #21262d;
    border: 1px solid #30363d;
    color: #8b949e;
    padding: 0.4rem 0.8rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.8rem;
  }

  .filters button.active {
    background: #1c2128;
    border-color: #58a6ff;
    color: #58a6ff;
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 1rem;
  }

  .empty { color: #8b949e; text-align: center; padding: 3rem; }
  .error { color: #f85149; padding: 1rem; background: #1c0d0d; border-radius: 6px; }
</style>
