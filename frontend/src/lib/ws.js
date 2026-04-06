export function connectWebSocket(path, onMessage, onError) {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const ws = new WebSocket(`${protocol}//${window.location.host}${path}`);

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      onMessage(data);
    } catch {
      // Plain text message (logs)
      onMessage(event.data);
    }
  };

  ws.onerror = (event) => {
    if (onError) onError(event);
  };

  return ws;
}
