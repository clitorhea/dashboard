<script>
  import { api } from '../lib/api.js';

  let { template = null, onClose, onDeployed } = $props();

  let serviceName = $state(template?.name?.toLowerCase().replace(/\s+/g, '-') || '');
  let compose = $state(template?.compose || '');
  let envVars = $state('');
  let deploying = $state(false);
  let error = $state('');
  let result = $state('');

  async function deploy() {
    if (!serviceName) { error = 'Service name is required'; return; }
    deploying = true;
    error = '';

    const env = {};
    if (envVars.trim()) {
      for (const line of envVars.split('\n')) {
        const idx = line.indexOf('=');
        if (idx > 0) {
          env[line.substring(0, idx).trim()] = line.substring(idx + 1).trim();
        }
      }
    }

    try {
      const data = {
        service_name: serviceName,
        env: Object.keys(env).length > 0 ? env : undefined,
      };
      if (template) {
        data.template_id = template.id;
      } else {
        data.compose = compose;
      }

      const res = await api.deploy(data);
      result = res.output || 'Deployed successfully';
      if (onDeployed) onDeployed();
    } catch (e) {
      error = e.message;
    } finally {
      deploying = false;
    }
  }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="overlay" onclick={onClose}>
  <div class="modal" onclick={(e) => e.stopPropagation()}>
    <div class="modal-header">
      <h3>{template ? `Deploy ${template.name}` : 'Custom Deploy'}</h3>
      <button class="close" onclick={onClose}>&times;</button>
    </div>

    <div class="modal-body">
      <label>
        Service Name
        <input type="text" bind:value={serviceName} placeholder="my-service" />
      </label>

      {#if !template}
        <label>
          Docker Compose YAML
          <textarea bind:value={compose} rows="10" placeholder="services:&#10;  app:&#10;    image: ..."></textarea>
        </label>
      {/if}

      <label>
        Environment Variables (KEY=VALUE, one per line)
        <textarea bind:value={envVars} rows="4" placeholder="PUID=1000&#10;PGID=1000"></textarea>
      </label>

      {#if error}
        <div class="error">{error}</div>
      {/if}

      {#if result}
        <pre class="result">{result}</pre>
      {/if}
    </div>

    <div class="modal-footer">
      <button class="btn cancel" onclick={onClose}>Cancel</button>
      <button class="btn deploy" onclick={deploy} disabled={deploying}>
        {deploying ? 'Deploying...' : 'Deploy'}
      </button>
    </div>
  </div>
</div>

<style>
  .overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 200;
  }

  .modal {
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    width: 90%;
    max-width: 560px;
    max-height: 80vh;
    overflow-y: auto;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    border-bottom: 1px solid #30363d;
  }

  .modal-header h3 { font-size: 1rem; }

  .close {
    background: none;
    border: none;
    color: #8b949e;
    font-size: 1.5rem;
    cursor: pointer;
  }

  .modal-body {
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  label {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    font-size: 0.85rem;
    color: #8b949e;
  }

  input, textarea {
    background: #0d1117;
    border: 1px solid #30363d;
    border-radius: 4px;
    padding: 0.5rem;
    color: #e1e4e8;
    font-family: inherit;
    font-size: 0.85rem;
  }

  textarea { resize: vertical; font-family: monospace; }

  .error {
    color: #f85149;
    font-size: 0.85rem;
    padding: 0.5rem;
    background: #1c0d0d;
    border-radius: 4px;
  }

  .result {
    background: #0d1117;
    border: 1px solid #30363d;
    border-radius: 4px;
    padding: 0.75rem;
    font-size: 0.8rem;
    max-height: 150px;
    overflow-y: auto;
    white-space: pre-wrap;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.5rem;
    padding: 1rem 1.5rem;
    border-top: 1px solid #30363d;
  }

  .btn {
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.85rem;
    border: 1px solid #30363d;
  }

  .btn.cancel { background: #21262d; color: #e1e4e8; }
  .btn.deploy { background: #238636; color: #fff; border-color: #238636; }
  .btn:disabled { opacity: 0.5; cursor: not-allowed; }
</style>
