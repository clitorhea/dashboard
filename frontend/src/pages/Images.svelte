<script>
  import { onMount } from 'svelte';
  import { api } from '../lib/api.js';

  let images = $state([]);
  let error = $state('');
  let pullInput = $state('');
  let pulling = $state(false);

  async function load() {
    try {
      images = await api.images();
    } catch (e) {
      error = e.message;
    }
  }

  onMount(load);

  function formatBytes(bytes) {
    if (!bytes) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }

  function getTag(img) {
    if (img.RepoTags && img.RepoTags.length > 0) return img.RepoTags[0];
    return img.Id?.substring(7, 19) || 'untagged';
  }

  async function pull() {
    if (!pullInput.trim()) return;
    pulling = true;
    try {
      await api.pullImage(pullInput.trim());
      pullInput = '';
      await load();
    } catch (e) {
      alert(`Pull failed: ${e.message}`);
    } finally {
      pulling = false;
    }
  }

  async function remove(id) {
    const tag = images.find(i => i.Id === id);
    if (!confirm(`Remove image ${getTag(tag)}?`)) return;
    try {
      await api.removeImage(id);
      await load();
    } catch (e) {
      alert(`Remove failed: ${e.message}`);
    }
  }
</script>

<div class="images-page">
  <div class="header">
    <h1>Images</h1>
    <div class="pull-form">
      <input type="text" bind:value={pullInput} placeholder="nginx:latest" onkeydown={(e) => e.key === 'Enter' && pull()} />
      <button onclick={pull} disabled={pulling}>{pulling ? 'Pulling...' : 'Pull'}</button>
    </div>
  </div>

  {#if error}
    <div class="error">{error}</div>
  {:else}
    <table>
      <thead>
        <tr>
          <th>Repository:Tag</th>
          <th>ID</th>
          <th>Size</th>
          <th>Created</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {#each images as img (img.Id)}
          <tr>
            <td class="mono">{getTag(img)}</td>
            <td class="mono">{img.Id?.substring(7, 19)}</td>
            <td>{formatBytes(img.Size)}</td>
            <td>{new Date(img.Created * 1000).toLocaleDateString()}</td>
            <td>
              <button class="btn-remove" onclick={() => remove(img.Id)}>Remove</button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>

<style>
  .images-page { max-width: 1000px; }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    flex-wrap: wrap;
    gap: 1rem;
  }

  h1 { font-size: 1.5rem; }

  .pull-form {
    display: flex;
    gap: 0.5rem;
  }

  .pull-form input {
    background: #0d1117;
    border: 1px solid #30363d;
    border-radius: 4px;
    padding: 0.5rem 0.75rem;
    color: #e1e4e8;
    font-size: 0.85rem;
    width: 220px;
  }

  .pull-form button {
    background: #238636;
    color: #fff;
    border: none;
    border-radius: 4px;
    padding: 0.5rem 1rem;
    cursor: pointer;
    font-size: 0.85rem;
  }

  .pull-form button:disabled { opacity: 0.5; }

  table {
    width: 100%;
    border-collapse: collapse;
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    overflow: hidden;
  }

  th {
    text-align: left;
    padding: 0.75rem 1rem;
    font-size: 0.75rem;
    text-transform: uppercase;
    color: #8b949e;
    border-bottom: 1px solid #30363d;
  }

  td {
    padding: 0.6rem 1rem;
    font-size: 0.85rem;
    border-bottom: 1px solid #21262d;
  }

  .mono { font-family: monospace; }

  .btn-remove {
    background: none;
    border: 1px solid #da3633;
    color: #f85149;
    padding: 0.25rem 0.6rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.75rem;
  }

  .btn-remove:hover { background: #da3633; color: #fff; }

  .error { color: #f85149; padding: 1rem; background: #1c0d0d; border-radius: 6px; }
</style>
