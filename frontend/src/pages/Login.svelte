<script>
  import { user } from '../lib/stores.js';
  import { api } from '../lib/api.js';
  import { push } from 'svelte-spa-router';

  let username = $state('');
  let password = $state('');
  let error = $state('');
  let submitting = $state(false);

  async function login(e) {
    e.preventDefault();
    error = '';
    submitting = true;
    try {
      const u = await api.login(username, password);
      user.set(u);
      push('/');
    } catch (err) {
      error = err.message || 'Login failed';
    } finally {
      submitting = false;
    }
  }
</script>

<div class="login-page">
  <form class="login-form" onsubmit={login}>
    <h1>NAS Dashboard</h1>
    <p class="subtitle">Sign in to manage your services</p>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    <label>
      Username
      <input type="text" bind:value={username} required autocomplete="username" />
    </label>

    <label>
      Password
      <input type="password" bind:value={password} required autocomplete="current-password" />
    </label>

    <button type="submit" disabled={submitting}>
      {submitting ? 'Signing in...' : 'Sign In'}
    </button>
  </form>
</div>

<style>
  .login-page {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background: #0f1117;
  }

  .login-form {
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

  h1 {
    text-align: center;
    font-size: 1.3rem;
    color: #58a6ff;
  }

  .subtitle {
    text-align: center;
    font-size: 0.85rem;
    color: #8b949e;
  }

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

  input:focus {
    outline: none;
    border-color: #58a6ff;
  }

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
