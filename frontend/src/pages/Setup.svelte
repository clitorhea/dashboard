<script>
  import { user } from '../lib/stores.js';
  import { api } from '../lib/api.js';
  import { push } from 'svelte-spa-router';

  let username = $state('');
  let password = $state('');
  let confirmPassword = $state('');
  let error = $state('');
  let submitting = $state(false);

  async function setup(e) {
    e.preventDefault();
    error = '';

    if (password !== confirmPassword) {
      error = 'Passwords do not match';
      return;
    }

    if (password.length < 8) {
      error = 'Password must be at least 8 characters';
      return;
    }

    submitting = true;
    try {
      const res = await api.setup(username, password);
      user.set(res.user);
      push('/');
    } catch (err) {
      error = err.message || 'Setup failed';
    } finally {
      submitting = false;
    }
  }
</script>

<div class="setup-page">
  <form class="setup-form" onsubmit={setup}>
    <h1>Welcome</h1>
    <p class="subtitle">Create your admin account to get started</p>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    <label>
      Username
      <input type="text" bind:value={username} required autocomplete="username" />
    </label>

    <label>
      Password
      <input type="password" bind:value={password} required autocomplete="new-password" />
    </label>

    <label>
      Confirm Password
      <input type="password" bind:value={confirmPassword} required />
    </label>

    <button type="submit" disabled={submitting}>
      {submitting ? 'Creating...' : 'Create Account'}
    </button>
  </form>
</div>

<style>
  .setup-page {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background: #0f1117;
  }

  .setup-form {
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    padding: 2rem;
    width: 100%;
    max-width: 360px;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  h1 { text-align: center; font-size: 1.3rem; color: #58a6ff; }
  .subtitle { text-align: center; font-size: 0.85rem; color: #8b949e; }

  label {
    display: flex;
    flex-direction: column;
    gap: 0.3rem;
    font-size: 0.85rem;
    color: #8b949e;
  }

  input {
    background: #0d1117;
    border: 1px solid #30363d;
    border-radius: 4px;
    padding: 0.6rem;
    color: #e1e4e8;
    font-size: 0.9rem;
  }

  input:focus { outline: none; border-color: #58a6ff; }

  button {
    background: #238636;
    color: #fff;
    border: none;
    border-radius: 4px;
    padding: 0.7rem;
    font-size: 0.9rem;
    cursor: pointer;
    margin-top: 0.5rem;
  }

  button:hover:not(:disabled) { background: #2ea043; }
  button:disabled { opacity: 0.6; cursor: not-allowed; }

  .error {
    color: #f85149;
    font-size: 0.85rem;
    text-align: center;
    padding: 0.5rem;
    background: #1c0d0d;
    border-radius: 4px;
  }
</style>
