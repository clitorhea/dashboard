import { writable, derived } from 'svelte/store';

export const user = writable(null);
export const loading = writable(true);
export const notifications = writable([]);

// Derived count for the navbar badge
export const unreadCount = derived(notifications, $n => $n.length);
