const BASE = '/api';

async function request(path, options = {}) {
  const res = await fetch(`${BASE}${path}`, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...options.headers },
    ...options,
  });

  if (res.status === 401) {
    window.location.hash = '#/login';
    throw new Error('Unauthorized');
  }

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `Request failed: ${res.status}`);
  }

  return res.json();
}

export const api = {
  // Auth
  login: (username, password) =>
    request('/auth/login', { method: 'POST', body: JSON.stringify({ username, password }) }),
  logout: () =>
    request('/auth/logout', { method: 'POST' }),
  me: () =>
    request('/auth/me'),
  setup: (username, password) =>
    request('/auth/setup', { method: 'POST', body: JSON.stringify({ username, password }) }),

  // Containers
  containers: () => request('/containers'),
  container: (id) => request(`/containers/${id}`),
  startContainer: (id) => request(`/containers/${id}/start`, { method: 'POST' }),
  stopContainer: (id) => request(`/containers/${id}/stop`, { method: 'POST' }),
  restartContainer: (id) => request(`/containers/${id}/restart`, { method: 'POST' }),
  removeContainer: (id) => request(`/containers/${id}`, { method: 'DELETE' }),

  // Images
  images: () => request('/images'),
  removeImage: (id) => request(`/images/${id}`, { method: 'DELETE' }),
  pullImage: (image) => request('/images/pull', { method: 'POST', body: JSON.stringify({ image }) }),

  // System
  systemInfo: () => request('/system/info'),
  dockerInfo: () => request('/system/docker'),

  // Templates
  templates: () => request('/templates'),
  template: (id) => request(`/templates/${id}`),
  deploy: (data) => request('/services/deploy', { method: 'POST', body: JSON.stringify(data) }),
};
