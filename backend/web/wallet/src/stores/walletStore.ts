import { create } from 'zustand';
import { api } from '../lib/api';
import { weiToSuda } from '../lib/format';

export interface HistoryItem {
  id: string;
  tx_hash: string;
  log_index: number;
  from_user_id: string;
  to_user_id: string;
  from_address: string;
  to_address: string;
  amount_wei: string;
  type: 'P2P_TRANSFER' | 'DONATION' | 'PURCHASE' | 'MINT';
  status: string;
  related_entity_type: string;
  related_entity_id: string;
  note: string;
  block_number: number;
  confirmed_at: string;
}

interface WalletState {
  address: string;
  balanceWei: bigint;
  syncedAt: number;
  history: HistoryItem[];
  historyLoading: boolean;
  historyHasMore: boolean;
  loading: boolean;
  error: string | null;
  // provisioning: wallet row not ready yet (/wallet/me 404). Not an error —
  // the welcome bonus is still being set up; UI shows a soft message and keeps polling.
  provisioning: boolean;

  fetchMe: () => Promise<void>;
  fetchHistory: (reset?: boolean, offset?: number) => Promise<void>;
  handleWsEvent: (evt: { type: string; payload: Record<string, unknown> }) => void;
  prependLiveHistory: (item: HistoryItem) => void;
}

const HISTORY_KEY = (a: string, b: number) => `${a}_${b}`;

export const useWalletStore = create<WalletState>((set, get) => ({
  address: '',
  balanceWei: 0n,
  syncedAt: 0,
  history: [],
  historyLoading: false,
  historyHasMore: true,
  loading: false,
  error: null,
  provisioning: false,

  fetchMe: async () => {
    set({ loading: true, error: null });
    let attempts = 0;
    while (attempts < 5) {
      try {
        const { data } = await api.get<{
          address: string;
          suda_balance_wei: string;
          decimals: number;
        }>('/tx/wallet/me');
        set({
          address: data.address,
          balanceWei: BigInt(data.suda_balance_wei),
          syncedAt: Date.now(),
          loading: false,
          error: null,
          provisioning: false,
        });
        return;
      } catch (err) {
        const axiosErr = err as { response?: { status?: number } };
        if (axiosErr.response?.status === 404) {
          // Wallet still being created (welcome bonus in flight). Retry a few
          // times, then hand off to the background poll in HomePage.
          attempts++;
          await new Promise((r) => setTimeout(r, 2000));
          continue;
        }
        // Real failure — surface it and stop the "Connecting…" spinner.
        set({ loading: false, error: 'Failed to load wallet', syncedAt: Date.now() });
        return;
      }
    }
    // 404 exhausted: wallet still provisioning. Mark synced (so the UI leaves
    // "Connecting…" and the welcome-bonus poll starts) without a hard error.
    set({ loading: false, error: null, provisioning: true, syncedAt: Date.now() });
  },

  fetchHistory: async (reset = false, offset = 0) => {
    if (get().historyLoading) return;
    set({ historyLoading: true });
    try {
      const { data } = await api.get<{ items: HistoryItem[] }>('/tx/wallet/history', {
        params: { limit: 50, offset },
      });
      const items = data.items ?? [];
      if (reset) {
        set({ history: items, historyHasMore: items.length === 50, historyLoading: false });
      } else {
        const existing = get().history;
        const existingKeys = new Set(existing.map((i) => HISTORY_KEY(i.tx_hash, i.log_index)));
        const newItems = items.filter((i) => !existingKeys.has(HISTORY_KEY(i.tx_hash, i.log_index)));
        set({
          history: [...existing, ...newItems],
          historyHasMore: items.length === 50,
          historyLoading: false,
        });
      }
    } catch {
      set({ historyLoading: false });
    }
  },

  handleWsEvent: (evt) => {
    switch (evt.type) {
      case 'SUDA_RECEIVED':
      case 'SUDA_SENT':
      case 'DONATION_SENT':
      case 'DONATION_RECEIVED':
      case 'PURCHASE_COMPLETED':
        // Cheap refresh — backend is authoritative; we don't try to merge
        // event payload into existing state to avoid drift / duplicates.
        void get().fetchMe();
        void get().fetchHistory(true);
        break;
      default:
        break;
    }
  },

  prependLiveHistory: (item) => {
    const existing = get().history;
    const key = HISTORY_KEY(item.tx_hash, item.log_index);
    if (existing.some((i) => HISTORY_KEY(i.tx_hash, i.log_index) === key)) return;
    set({ history: [item, ...existing] });
  },
}));

export function weiBalance(): string {
  return weiToSuda(useWalletStore.getState().balanceWei.toString(), 2);
}
