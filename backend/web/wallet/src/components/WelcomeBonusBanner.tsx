import { type HistoryItem } from '../stores/walletStore';

interface Props {
  history: HistoryItem[];
  onBuy: () => void;
}

const ZERO_ADDR = '0x0000000000000000000000000000000000000000';

export default function WelcomeBonusBanner({ history, onBuy }: Props) {
  const hasMint = history.some(
    (h) => h.type === 'MINT' && h.from_address.toLowerCase() === ZERO_ADDR.toLowerCase(),
  );
  if (!hasMint) return null;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', padding: '0 14px 16px' }}>
      <div className="tip-card">
        <div style={{ fontSize: '24px', flexShrink: 0 }}>🎁</div>
        <div>
          <strong style={{ display: 'block', fontSize: '14px', fontWeight: 700, marginBottom: '3px' }}>
            Welcome bonus active
          </strong>
          <span style={{ color: 'var(--text-muted)' }}>You got 100 SUDA. Try sending a tip in a chat.</span>
        </div>
      </div>
      <div className="tip-card alt">
        <div style={{ fontSize: '24px', flexShrink: 0 }}>📈</div>
        <div style={{ flex: 1 }}>
          <strong style={{ display: 'block', fontSize: '14px', fontWeight: 700, marginBottom: '3px' }}>
            Top up for more
          </strong>
          <span style={{ color: 'var(--text-muted)' }}>Reach 100 SUDA to unlock gated channels.</span>
        </div>
        <button className="btn-link" style={{ fontSize: '13px', padding: '4px 0' }} onClick={onBuy}>
          Buy →
        </button>
      </div>
    </div>
  );
}
