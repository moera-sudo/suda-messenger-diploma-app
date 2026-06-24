import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft } from 'lucide-react';
import { useWalletStore, type HistoryItem } from '../stores/walletStore';
import AppShell from '../components/AppShell';
import ActivityItem from '../components/ActivityItem';
import TxDetails from '../components/TxDetails';

type Filter = 'all' | 'P2P_TRANSFER' | 'DONATION' | 'PURCHASE' | 'MINT';

const FILTERS: { id: Filter; label: string }[] = [
  { id: 'all', label: 'All' },
  { id: 'P2P_TRANSFER', label: 'Sent/Received' },
  { id: 'DONATION', label: 'Donations' },
  { id: 'PURCHASE', label: 'Purchases' },
];

export default function HistoryPage() {
  const navigate = useNavigate();
  const { address, history, historyLoading, historyHasMore, fetchHistory } = useWalletStore();
  const [filter, setFilter] = useState<Filter>('all');
  const [selectedTx, setSelectedTx] = useState<HistoryItem | null>(null);
  const loaderRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    void fetchHistory(true, 0);
  }, [fetchHistory]);

  useEffect(() => {
    const el = loaderRef.current;
    if (!el) return;
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry?.isIntersecting && historyHasMore && !historyLoading) {
          void fetchHistory(false, history.length);
        }
      },
      { threshold: 0.1 },
    );
    observer.observe(el);
    return () => observer.disconnect();
  }, [historyHasMore, historyLoading, history.length, fetchHistory]);

  const filtered = filter === 'all' ? history : history.filter((h) => h.type === filter);

  return (
    <AppShell>
      <div className="wv-app-bar">
        <button
          className="btn btn-ghost btn-sm"
          style={{ padding: '8px', borderRadius: '50%', border: 0 }}
          onClick={() => navigate(-1)}
        >
          <ArrowLeft size={20} />
        </button>
        <h3>Activity</h3>
      </div>

      <div className="chip-row">
        {FILTERS.map((f) => (
          <button
            key={f.id}
            className={`chip ${filter === f.id ? 'active' : ''}`}
            onClick={() => setFilter(f.id)}
          >
            {f.label}
          </button>
        ))}
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: '4px', padding: '0 14px 24px' }}>
        {filtered.map((item) => (
          <ActivityItem
            key={`${item.tx_hash}_${item.log_index}`}
            item={item}
            myAddress={address}
            onClick={() => setSelectedTx(item)}
          />
        ))}

        {filtered.length === 0 && !historyLoading && (
          <div style={{ textAlign: 'center', padding: '48px 0', color: 'var(--text-muted)', fontSize: '14px' }}>
            No transactions
          </div>
        )}

        <div ref={loaderRef} style={{ height: '1px' }} />

        {historyLoading && (
          <div style={{ display: 'flex', justifyContent: 'center', padding: '16px' }}>
            <div className="spinner" style={{ width: '28px', height: '28px', borderWidth: '2px' }} />
          </div>
        )}
      </div>

      {selectedTx && (
        <TxDetails item={selectedTx} myAddress={address} onClose={() => setSelectedTx(null)} />
      )}
    </AppShell>
  );
}
