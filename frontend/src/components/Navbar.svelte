<script>
  import { user, unreadCount } from '../lib/stores.js';
  import { api } from '../lib/api.js';

  let currentPath = $state(window.location.hash.replace('#', '') || '/');

  function updatePath() {
    currentPath = window.location.hash.replace('#', '') || '/';
  }

  $effect(() => {
    window.addEventListener('hashchange', updatePath);
    return () => window.removeEventListener('hashchange', updatePath);
  });

  async function logout() {
    await api.logout();
    user.set(null);
    window.location.hash = '#/login';
  }

  const links = [
    { path: '/', label: 'Dashboard', icon: '⊞' },
    { path: '/containers', label: 'Containers', icon: '▦' },
    { path: '/images', label: 'Images', icon: '◫' },
    { path: '/templates', label: 'Templates', icon: '⊕' },
    { path: '/settings', label: 'Settings', icon: '⚙' },
  ];
</script>

<nav class="sidebar">
  <div class="logo">
    <h2>NAS Dashboard</h2>
    {#if $unreadCount > 0}
      <span class="notif-badge" title="{$unreadCount} alert{$unreadCount > 1 ? 's' : ''}">{$unreadCount}</span>
    {/if}
  </div>

  <ul class="nav-links">
    {#each links as link}
      <li>
        <a href="#{link.path}" class:active={currentPath === link.path}>
          <span class="icon">{link.icon}</span>
          {link.label}
          {#if link.path === '/' && $unreadCount > 0}
            <span class="link-badge">{$unreadCount}</span>
          {/if}
        </a>
      </li>
    {/each}
  </ul>

  <div class="nav-footer">
    <span class="username">{$user?.username}</span>
    <button onclick={logout}>Logout</button>
  </div>
</nav>

<style>
  .sidebar {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    width: 240px;
    background: #161b22;
    border-right: 1px solid #30363d;
    display: flex;
    flex-direction: column;
    z-index: 100;
  }

  .logo {
    padding: 1.5rem;
    border-bottom: 1px solid #30363d;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .logo h2 {
    font-size: 1.1rem;
    color: #58a6ff;
    flex: 1;
  }

  .notif-badge {
    background: #da3633;
    color: #fff;
    font-size: 0.7rem;
    font-weight: 700;
    min-width: 18px;
    height: 18px;
    border-radius: 9px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 4px;
    animation: pulse 2s infinite;
  }

  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.6; }
  }

  .nav-links {
    list-style: none;
    padding: 0.5rem 0;
    flex: 1;
  }

  .nav-links a {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1.5rem;
    color: #8b949e;
    text-decoration: none;
    transition: all 0.15s;
  }

  .nav-links a:hover {
    color: #e1e4e8;
    background: #1c2128;
  }

  .nav-links a.active {
    color: #58a6ff;
    background: #1c2128;
    border-right: 2px solid #58a6ff;
  }

  .icon {
    font-size: 1.1rem;
    width: 1.5rem;
    text-align: center;
  }

  .link-badge {
    margin-left: auto;
    background: #da3633;
    color: #fff;
    font-size: 0.65rem;
    font-weight: 700;
    min-width: 16px;
    height: 16px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 3px;
  }

  .nav-footer {
    padding: 1rem 1.5rem;
    border-top: 1px solid #30363d;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .username {
    font-size: 0.85rem;
    color: #8b949e;
  }

  .nav-footer button {
    background: none;
    border: 1px solid #30363d;
    color: #8b949e;
    padding: 0.3rem 0.75rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.8rem;
  }

  .nav-footer button:hover {
    color: #e1e4e8;
    border-color: #8b949e;
  }
</style>
