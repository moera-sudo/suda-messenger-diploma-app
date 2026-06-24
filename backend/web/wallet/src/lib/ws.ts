import { useAuthStore } from '../stores/authStore';
import { useWalletStore } from '../stores/walletStore';

let ws: WebSocket | null = null;
let reconnectTimer: ReturnType<typeof setTimeout> | null = null;

export function connectWs(): void {
  if (ws && (ws.readyState === WebSocket.CONNECTING || ws.readyState === WebSocket.OPEN)) return;

  const token = useAuthStore.getState().token;
  if (!token) return;

  const proto = location.protocol === 'https:' ? 'wss' : 'ws';
  ws = new WebSocket(`${proto}://${location.host}/ws?token=${token}`);

  ws.onmessage = (e) => {
    try {
      const evt = JSON.parse(e.data as string) as { type: string; payload: Record<string, unknown> };
      useWalletStore.getState().handleWsEvent(evt);
    } catch {
      // ignore malformed messages
    }
  };

  ws.onclose = () => {
    ws = null;
    if (reconnectTimer) clearTimeout(reconnectTimer);
    reconnectTimer = setTimeout(() => connectWs(), 2000);
  };

  ws.onerror = () => {
    ws?.close();
  };
}

export function disconnectWs(): void {
  if (reconnectTimer) {
    clearTimeout(reconnectTimer);
    reconnectTimer = null;
  }
  ws?.close();
  ws = null;
}
