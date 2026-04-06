<script>
  import { onMount } from 'svelte';
  import { api } from '../lib/api.js';
  import DeployModal from '../components/DeployModal.svelte';

  let templates = $state([]);
  let error = $state('');
  let selectedTemplate = $state(null);
  let showCustomDeploy = $state(false);

  onMount(async () => {
    try {
      templates = await api.templates();
    } catch (e) {
      error = e.message;
    }
  });

  function closeModal() {
    selectedTemplate = null;
    showCustomDeploy = false;
  }
</script>

<div class="templates-page">
  <div class="header">
    <h1>Service Templates</h1>
    <button class="btn-custom" onclick={() => showCustomDeploy = true}>Custom Deploy</button>
  </div>

  {#if error}
    <div class="error">{error}</div>
  {:else if templates.length === 0}
    <div class="empty">
      <p>No templates available yet.</p>
      <p class="hint">Use "Custom Deploy" to deploy with your own docker-compose YAML.</p>
    </div>
  {:else}
    <div class="grid">
      {#each templates as tmpl (tmpl.id)}
        <div class="template-card">
          <div class="template-header">
            <h3>{tmpl.name}</h3>
            {#if tmpl.category}
              <span class="category">{tmpl.category}</span>
            {/if}
          </div>
          <p class="description">{tmpl.description || 'No description'}</p>
          <button class="btn-deploy" onclick={() => selectedTemplate = tmpl}>Deploy</button>
        </div>
      {/each}
    </div>
  {/if}

  {#if selectedTemplate}
    <DeployModal template={selectedTemplate} onClose={closeModal} onDeployed={closeModal} />
  {/if}

  {#if showCustomDeploy}
    <DeployModal onClose={closeModal} onDeployed={closeModal} />
  {/if}
</div>

<style>
  .templates-page { max-width: 1000px; }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }

  h1 { font-size: 1.5rem; }

  .btn-custom {
    background: #238636;
    color: #fff;
    border: none;
    border-radius: 4px;
    padding: 0.5rem 1rem;
    cursor: pointer;
    font-size: 0.85rem;
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 1rem;
  }

  .template-card {
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    padding: 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .template-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .template-header h3 { font-size: 1rem; }

  .category {
    font-size: 0.7rem;
    text-transform: uppercase;
    color: #8b949e;
    background: #21262d;
    padding: 0.15rem 0.5rem;
    border-radius: 3px;
  }

  .description {
    font-size: 0.85rem;
    color: #8b949e;
    flex: 1;
  }

  .btn-deploy {
    background: #21262d;
    border: 1px solid #30363d;
    color: #58a6ff;
    padding: 0.5rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.85rem;
  }

  .btn-deploy:hover { background: #30363d; }

  .empty { text-align: center; padding: 3rem; color: #8b949e; }
  .hint { font-size: 0.85rem; margin-top: 0.5rem; }
  .error { color: #f85149; padding: 1rem; background: #1c0d0d; border-radius: 6px; }
</style>
