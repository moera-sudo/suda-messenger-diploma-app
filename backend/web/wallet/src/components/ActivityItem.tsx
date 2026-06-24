import { ArrowDownLeft, ArrowUpRight, Heart, Plus, Coins } from 'lucide-react';
import { type HistoryItem } from '../stores/walletStore';
import { weiToSuda, formatAgo, truncateAddr } from '../lib/format';
import { useAuthStore } from '../stores/authStore';

interface Props {
  item: HistoryItem;
  myAddress: string;
  onClick?: () => void;
}

function txMeta(item: HistoryItem, myAddress: string) {
  const isMe = item.from_address.toLowerCase() === myAddress.toLowerCase();
  if (item.type === 'MINT') {
    return {
      icon: Coins,
      color: '#3DDC97',
      title: 'Welcome Bonus',
      subtitle: 'SUDA credited to your wallet',
      dir: 'in' as const,
    };
  }
  if (item.type === 'PURCHASE') {
    return {
      icon: Plus,
      color: '#FFB347',
      title: 'Bought SUDA',
      subtitle: item.note || 'Package purchase',
      dir: 'in' as const,
    };
  }
  if (item.type === 'DONATION') {
    if (isMe) {
      return {
        icon: Heart,
        color: '#FF8FB3',
        title: `Donation`,
        subtitle: item.note || 'Donation sent',
        dir: 'out' as const,
      };
    }
    return {
      icon: Heart,
      color: '#FF8FB3',
      title: `Donation received`,
      subtitle: item.note || 'Donation',
      dir: 'in' as const,
    };
  }
  if (isMe) {
    return {
      icon: ArrowUpRight,
      color: 'var(--text-accent)',
      title: `To ${truncateAddr(item.to_address)}`,
      subtitle: item.note || 'Sent SUDA',
      dir: 'out' as const,
    };
  }
  return {
    icon: ArrowDownLeft,
    color: '#3DDC97',
    title: `From ${truncateAddr(item.from_address)}`,
    subtitle: item.note || 'Received SUDA',
    dir: 'in' as const,
  };
}

export default function ActivityItem({ item, myAddress, onClick }: Props) {
  const meta = txMeta(item, myAddress);
  const Icon = meta.icon;
  const amount = weiToSuda(item.amount_wei, 2);
  const time = item.confirmed_at ? formatAgo(item.confirmed_at) : '';

  return (
    <button className="tx-row" onClick={onClick}>
      <span
        className="tx-ico"
        style={{ color: meta.color, background: `${meta.color}22` }}
      >
        <Icon size={18} color={meta.color} />
      </span>
      <div style={{ flex: 1, minWidth: 0 }}>
        <div style={{ fontSize: '15px', fontWeight: 600, color: 'var(--text-primary)', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
          {meta.title}
        </div>
        <div style={{ fontSize: '13px', color: 'var(--text-muted)', marginTop: '2px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
          {meta.subtitle}
          {time && ` · ${time}`}
        </div>
      </div>
      <div className={`tx-amount-col ${meta.dir}`}>
        <span>{meta.dir === 'in' ? '+' : '−'}{amount}</span>
        <small style={{ fontSize: '11px', color: 'var(--text-muted)', fontWeight: 600 }}>SUDA</small>
      </div>
    </button>
  );
}

export function useMyAddress() {
  return useAuthStore((s) => s.userId);
}
