import axios from 'axios';
import { useAuthStore } from '../stores/authStore';
import { tryReAuth } from './webapp';

export const api = axios.create({ baseURL: '/api/v1' });

api.interceptors.request.use((cfg) => {
  const tok = useAuthStore.getState().token;
  if (tok && cfg.headers) {
    cfg.headers.Authorization = `Bearer ${tok}`;
  }
  return cfg;
});

let isRefreshing = false;

api.interceptors.response.use(undefined, async (err) => {
  if (err.response?.status === 401 && !isRefreshing) {
    isRefreshing = true;
    try {
      const ok = await tryReAuth();
      if (ok) {
        isRefreshing = false;
        return api.request(err.config);
      }
    } catch {
      // ignore
    }
    isRefreshing = false;
    useAuthStore.getState().clear();
  }
  throw err;
});
