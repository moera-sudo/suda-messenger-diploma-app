import axios from 'axios';
import { useAuthStore } from '../stores/authStore';
import { extractApiErrorCode } from './errors';

function parseInitData(): string | null {
  const url = new URL(window.location.href);
  const fromSearch = url.searchParams.get('initData');
  if (fromSearch) return fromSearch;
  const hash = url.hash;
  const qIndex = hash.indexOf('?');
  if (qIndex !== -1) {
    return new URLSearchParams(hash.slice(qIndex + 1)).get('initData');
  }
  return null;
}

interface WebappSessionResponse {
  access_token: string;
  expires_in: number;
  user_id: string;
  username: string;
  avatar_media_id?: string;
}

async function exchangeInitData(initData: string): Promise<WebappSessionResponse> {
  const { data } = await axios.post<WebappSessionResponse>(
    '/api/v1/messenger/auth/webapp-session',
    { init_data: initData },
    // Timeout so a hung/unreachable backend (502, no network) surfaces as an
    // error instead of an endless "Connecting wallet…" spinner.
    { timeout: 15000 },
  );
  return data;
}

export async function initAuth(): Promise<void> {
  const store = useAuthStore.getState();
  store.setStatus('authing');

  const initData = parseInitData();
  if (!initData) {
    store.setStatus('error', 'NO_INIT_DATA');
    return;
  }
  store.setInitData(initData);

  try {
    const session = await exchangeInitData(initData);
    store.setSession(session.access_token, session.expires_in, session.user_id, session.username, session.avatar_media_id);
    store.setStatus('ready');
  } catch (err) {
    store.setStatus('error', extractApiErrorCode(err, 'AUTH_ERROR'));
  }
}

export async function tryReAuth(): Promise<boolean> {
  const store = useAuthStore.getState();
  const initData = store.initData;
  if (!initData) return false;
  try {
    const session = await exchangeInitData(initData);
    store.setSession(session.access_token, session.expires_in, session.user_id, session.username, session.avatar_media_id);
    return true;
  } catch {
    return false;
  }
}
