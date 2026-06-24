import { Copy, ExternalLink } from 'lucide-react';
import { type HistoryItem } from '../stores/walletStore';
import { weiToSuda, truncateAddr, formatAgo } from '../lib/format';
import BottomSheet from './BottomSheet';

interface Props {
  item: HistoryItem;
  myAddress: string;
  onClose: () => void;
}

function copyText(text: string) {
  void navigator.clipboard.writeText(text);
}

export default function TxDetails({ item, myAddress, onClose }: Props) {
  const isIn = item.to_address.toLowerCase() === myAddress.toLowerCase() || item.type === 'MINT' || item.type === 'PURCHASE';
  const amount = weiToSuda(item.amount_wei, 2);

  return (
    <BottomSheet onClose={onClose}>
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          gap: '8px',
          padding: '8px 0 20px',
          background: `linear-gradient(180deg, color-mix(in oklab, var(--primary), var(--background) 65%), transparent)`,
          borderRadius: '12px',
          marginBottom: '16px',
          textAlign: 'center',
        }}
      >
        <div
          style={{
            fontSize: '42px',
            fontWeight: 800,
            fontFamily: 'var(--ff-display)',
            color: isIn ? 'var(--success)' : 'var(--text-primary)',
          }}
        >
          {isIn ? '+' : '−'}{amount} <span style={{ fontSize: '16px', color: 'var(--text-muted)', fontWeight: 600 }}>SUDA</span>
        </div>
        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '6px',
            padding: '6px 14px',
            background: 'color-mix(in oklab, var(--success), transparent 85%)',
            border: '1px solid color-mix(in oklab, var(--success), transparent 60%)',
            borderRadius: '100px',
            fontSize: '13px',
            fontWeight: 700,
            color: 'var(--success)',
          }}
        >
          ✓ {item.status || 'Confirmed'}
        </div>
        {item.confirmed_at && (
          <div style={{ fontSize: '13px', color: 'var(--text-muted)' }}>
            {formatAgo(item.confirmed_at)}
          </div>
        )}
      </div>

      <div className="info-card">
        <Row label="Type" value={item.type.replace('_', ' ')} />
        <Row label="From" value={truncateAddr(item.from_address)} onCopy={() => copyText(item.from_address)} />
        <Row label="To" value={truncateAddr(item.to_address)} onCopy={() => copyText(item.to_address)} />
        {item.note && <Row label="Note" value={item.note} />}
        <Row label="Network fee" value="0.00 SUDA" sub="Paid by Suda" />
        <Row
          label="TX hash"
          value={`${item.tx_hash.slice(0, 10)}…${item.tx_hash.slice(-8)}`}
          mono
          onCopy={() => copyText(item.tx_hash)}
        />
        {item.block_number > 0 && (
          <Row label="Block" value={item.block_number.toLocaleString()} />
        )}
      </div>

      <button
        className="btn btn-ghost btn-md btn-full"
        style={{ marginTop: '12px', gap: '8px' }}
      >
        <ExternalLink size={16} />
        View on explorer
      </button>
    </BottomSheet>
  );
}

function Row({
  label,
  value,
  sub,
  mono,
  onCopy,
}: {
  label: string;
  value: string;
  sub?: string;
  mono?: boolean;
  onCopy?: () => void;
}) {
  return (
    <div className="detail-row">
      <span className="detail-label">{label}</span>
      <div className="detail-val" style={{ fontFamily: mono ? 'var(--ff-mono)' : undefined }}>
        <div>{value}</div>
        {sub && <div style={{ fontSize: '12px', color: 'var(--text-muted)', marginTop: '2px' }}>{sub}</div>}
      </div>
      {onCopy && (
        <button
          onClick={onCopy}
          style={{ background: 'transparent', border: 0, padding: '2px', flexShrink: 0, cursor: 'pointer' }}
        >
          <Copy size={14} color="var(--text-muted)" />
        </button>
      )}
    </div>
  );
}
