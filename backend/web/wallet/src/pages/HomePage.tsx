import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowUpRight, Plus, QrCode } from 'lucide-react';
import { useWalletStore, type HistoryItem } from '../stores/walletStore';
import { useAuthStore } from '../stores/authStore';
import { connectWs } from '../lib/ws';
import { weiToSuda } from '../lib/format';
import AppShell from '../components/AppShell';
import QuickAction from '../components/QuickAction';
import ActivityItem from '../components/ActivityItem';
import AddressChip from '../components/AddressChip';
import Avatar from '../components/Avatar';
import WelcomeBonusBanner from '../components/WelcomeBonusBanner';
import TxDetails from '../components/TxDetails';

export default function HomePage() {
  const navigate = useNavigate();
  const { address, balanceWei, syncedAt, history, error, provisioning, fetchMe, fetchHistory } = useWalletStore();
  const username = useAuthStore((s) => s.username);
  const avatarMediaId = useAuthStore((s) => s.avatarMediaId);
  const [selectedTx, setSelectedTx] = useState<HistoryItem | null>(null);
  const bonusPolls = useRef(0);

  useEffect(() => {
    void fetchMe();
    void fetchHistory(true);
    connectWs();
  }, [fetchMe, fetchHistory]);

  // Welcome bonus is broadcast on-chain and lands a few seconds after wallet
  // creation. /tx/wallet/me returns 200 with balance 0 in the meantime, so the
  // 404-retry in fetchMe doesn't cover it. Poll a few times while the balance is
  // still 0 and no MINT (welcome bonus) appears in history yet. The WS
  // SUDA_RECEIVED event also refreshes, this is the fallback when WS is late.
  useEffect(() => {
    if (!syncedAt) return;
    const hasMint = history.some(
      (i) => i.type === 'MINT' && /^0x0+$/i.test(i.from_address),
    );
    // Keep polling while the balance is still 0 (welcome bonus not landed) or the
    // wallet is still provisioning. Longer cap so a slow bonus still shows up.
    if ((balanceWei !== 0n && !provisioning) || hasMint || bonusPolls.current >= 10) return;
    const t = setTimeout(() => {
      bonusPolls.current += 1;
      void fetchMe();
      void fetchHistory(true);
    }, 2500);
    return () => clearTimeout(t);
  }, [syncedAt, balanceWei, history, provisioning, fetchMe, fetchHistory]);

  const balance = weiToSuda(balanceWei.toString(), 2);
  const syncedText = provisioning
    ? 'Setting up wallet…'
    : error
      ? error
      : syncedAt
        ? 'Synced · just now'
        : 'Connecting…';

  return (
    <AppShell>
      {/* Hero / Balance Card */}
      <div className="wallet-hero">
        <div className="wallet-hero-bg">
          <svg width="100%" height="100%" viewBox="0 0 400 280" preserveAspectRatio="none">
            <defs>
              <radialGradient id="wgrad" cx="50%" cy="50%" r="60%">
                <stop offset="0%" stopColor="#29DBF2" stopOpacity="0.35" />
                <stop offset="100%" stopColor="#0F1020" stopOpacity="0" />
              </radialGradient>
              <linearGradient id="ring" x1="0" y1="0" x2="1" y2="1">
                <stop offset="0%" stopColor="#644599" />
                <stop offset="100%" stopColor="#29DBF2" />
              </linearGradient>
            </defs>
            <circle cx="80" cy="60" r="120" fill="url(#wgrad)" />
            <circle cx="320" cy="220" r="160" fill="url(#wgrad)" opacity="0.55" />
            <g transform="translate(330 50)" opacity="0.6">
              <circle r="32" fill="none" stroke="url(#ring)" strokeWidth="1.4" />
              <circle r="22" fill="none" stroke="url(#ring)" strokeWidth="1" />
            </g>
          </svg>
        </div>

        <div className="wallet-hero-content">
          <div className="wallet-greet">
            <span className="dot-online" />
            <span>{syncedText}</span>
            {username && (
              <span
                style={{
                  marginLeft: 'auto',
                  display: 'inline-flex',
                  alignItems: 'center',
                  gap: '8px',
                  color: 'var(--text-accent)',
                  fontWeight: 600,
                }}
              >
                {username}
                <Avatar mediaId={avatarMediaId} username={username} size={28} />
              </span>
            )}
          </div>
          <div className="balance-display">
            <span className="balance-num">{balance}</span>
            <span className="balance-tic">SUDA</span>
          </div>
          <AddressChip address={address} onQrClick={() => navigate('/receive')} />
        </div>
      </div>

      {/* Quick Actions */}
      <div style={{ display: 'flex', gap: '10px', padding: '16px 14px 10px' }}>
        <QuickAction icon={ArrowUpRight} label="Send" onClick={() => navigate('/send')} />
        <QuickAction icon={Plus} label="Buy" onClick={() => navigate('/buy')} />
        <QuickAction icon={QrCode} label="Receive" onClick={() => navigate('/receive')} />
      </div>

      {/* Activity section */}
      <div style={{ padding: '8px 14px 0' }}>
        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            paddingBottom: '10px',
          }}
        >
          <span style={{ fontSize: '15px', fontWeight: 700 }}>Activity</span>
          <button
            className="btn-link"
            style={{ fontSize: '13px', padding: '4px 0' }}
            onClick={() => navigate('/history')}
          >
            See all →
          </button>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
          {history.slice(0, 4).map((item) => (
            <ActivityItem
              key={`${item.tx_hash}_${item.log_index}`}
              item={item}
              myAddress={address}
              onClick={() => setSelectedTx(item)}
            />
          ))}
          {history.length === 0 && (
            <div
              style={{
                textAlign: 'center',
                padding: '32px 0',
                color: 'var(--text-muted)',
                fontSize: '14px',
              }}
            >
              No transactions yet
            </div>
          )}
        </div>
      </div>

      {/* Welcome bonus banner */}
      <div style={{ marginTop: '8px' }}>
        <WelcomeBonusBanner history={history} onBuy={() => navigate('/buy')} />
      </div>

      {/* Tx details sheet */}
      {selectedTx && (
        <TxDetails
          item={selectedTx}
          myAddress={address}
          onClose={() => setSelectedTx(null)}
        />
      )}
    </AppShell>
  );
}
