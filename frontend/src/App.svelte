<script>
  import Router from 'svelte-spa-router';
  import { onMount } from 'svelte';
  import { user, loading } from './lib/stores.js';
  import { api } from './lib/api.js';
  import Navbar from './components/Navbar.svelte';
  import Login from './pages/Login.svelte';
  import Setup from './pages/Setup.svelte';
  import Dashboard from './pages/Dashboard.svelte';
  import Containers from './pages/Containers.svelte';
  import ContainerDetail from './pages/ContainerDetail.svelte';
  import Images from './pages/Images.svelte';
  import Templates from './pages/Templates.svelte';
  import Settings from './pages/Settings.svelte';

  const routes = {
    '/': Dashboard,
    '/login': Login,
    '/setup': Setup,
    '/containers': Containers,
    '/containers/:id': ContainerDetail,
    '/images': Images,
    '/templates': Templates,
    '/settings': Settings,
  };

  let isAuthenticated = $state(false);

  user.subscribe(value => {
    isAuthenticated = !!value;
  });

  onMount(async () => {
    try {
      const me = await api.me();
      user.set(me);
    } catch {
      try {
        const res = await fetch('/api/auth/setup', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({})
        });
        if (res.status === 400) {
          const body = await res.json();
          if (body.error === 'username and password required') {
            window.location.hash = '#/setup';
            loading.set(false);
            return;
          }
        }
      } catch { /* ignore */ }
      window.location.hash = '#/login';
    } finally {
      loading.set(false);
    }
  });
</script>

<div class="app">
  {#if $loading}
    <div class="loading">
      <div class="spinner"></div>
      <p>Loading...</p>
    </div>
  {:else}
    {#if isAuthenticated}
      <Navbar />
    {/if}
    <main class:with-nav={isAuthenticated}>
      <Router {routes} />
    </main>
  {/if}
</div>

<style>
  :global(*) {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  :global(body) {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #0f1117;
    color: #e1e4e8;
    min-height: 100vh;
  }

  .app {
    min-height: 100vh;
  }

  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100vh;
    gap: 1rem;
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 3px solid #30363d;
    border-top-color: #58a6ff;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  main {
    min-height: 100vh;
  }

  main.with-nav {
    margin-left: 240px;
    padding: 2rem;
  }
</style>
