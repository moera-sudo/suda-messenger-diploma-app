import { create } from 'zustand';

type AuthStatus = 'idle' | 'authing' | 'ready' | 'error';

interface AuthState {
  token: string | null;
  expiresAt: number;
  userId: string | null;
  username: string;
  avatarMediaId: string;
  initData: string | null;
  status: AuthStatus;
  errorMsg: string;
  setSession: (token: string, expiresInSec: number, userId: string, username: string, avatarMediaId?: string) => void;
  setInitData: (s: string) => void;
  setStatus: (s: AuthStatus, err?: string) => void;
  clear: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  token: null,
  expiresAt: 0,
  userId: null,
  username: '',
  avatarMediaId: '',
  initData: null,
  status: 'idle',
  errorMsg: '',
  setSession: (token, expiresInSec, userId, username, avatarMediaId = '') =>
    set({
      token,
      expiresAt: Date.now() + expiresInSec * 1000,
      userId,
      username,
      avatarMediaId,
    }),
  setInitData: (s) => set({ initData: s }),
  setStatus: (status, err = '') => set({ status, errorMsg: err }),
  clear: () => set({ token: null, userId: null, username: '', avatarMediaId: '', status: 'error', errorMsg: 'SESSION_EXPIRED' }),
}));
