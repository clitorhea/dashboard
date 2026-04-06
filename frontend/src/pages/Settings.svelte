<script>
  import { user } from '../lib/stores.js';
  import { api } from '../lib/api.js';

  let currentPassword = $state('');
  let newPassword = $state('');
  let confirmPassword = $state('');
  let message = $state('');
  let error = $state('');

  async function changePassword(e) {
    e.preventDefault();
    error = '';
    message = '';

    if (newPassword !== confirmPassword) {
      error = 'New passwords do not match';
      return;
    }

    if (newPassword.length < 8) {
      error = 'Password must be at least 8 characters';
      return;
    }

    // For now, re-login with current password to verify, then we'd need a dedicated endpoint
    // This is a placeholder — the backend would need a PUT /api/auth/password endpoint
    message = 'Password change not yet implemented on the backend';
  }
</script>

<div class="settings-page">
  <h1>Settings</h1>

  <div class="section">
    <h2>Account</h2>
    <div class="info-row">
      <span class="label">Username</span>
      <span class="value">{$user?.username}</span>
    </div>
  </div>

  <div class="section">
    <h2>Change Password</h2>
    <form onsubmit={changePassword}>
      {#if error}
        <div class="error">{error}</div>
      {/if}
      {#if message}
        <div class="message">{message}</div>
      {/if}

      <label>
        Current Password
        <input type="password" bind:value={currentPassword} required />
      </label>
      <label>
        New Password
        <input type="password" bind:value={newPassword} required />
      </label>
      <label>
        Confirm New Password
        <input type="password" bind:value={confirmPassword} required />
      </label>
      <button type="submit">Update Password</button>
    </form>
  </div>

  <div class="section">
    <h2>About</h2>
    <div class="info-row">
      <span class="label">Version</span>
      <span class="value">0.1.0</span>
    </div>
  </div>
</div>

<style>
  .settings-page { max-width: 600px; }

  h1 { font-size: 1.5rem; margin-bottom: 1.5rem; }
  h2 { font-size: 1rem; color: #8b949e; margin-bottom: 1rem; }

  .section {
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    padding: 1.5rem;
    margin-bottom: 1.5rem;
  }

  .info-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .label { color: #8b949e; font-size: 0.85rem; }
  .value { font-size: 0.9rem; }

  form {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
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
    padding: 0.5rem;
    color: #e1e4e8;
    font-size: 0.85rem;
  }

  input:focus { outline: none; border-color: #58a6ff; }

  button {
    background: #238636;
    color: #fff;
    border: none;
    border-radius: 4px;
    padding: 0.6rem;
    cursor: pointer;
    font-size: 0.85rem;
    margin-top: 0.5rem;
  }

  .error {
    color: #f85149;
    font-size: 0.85rem;
    padding: 0.5rem;
    background: #1c0d0d;
    border-radius: 4px;
  }

  .message {
    color: #3fb950;
    font-size: 0.85rem;
    padding: 0.5rem;
    background: #0d1c0d;
    border-radius: 4px;
  }
</style>
