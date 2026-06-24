import { QrCode, Copy, Check } from 'lucide-react';
import { useState } from 'react';
import { truncateAddr } from '../lib/format';

interface Props {
  address: string;
  onQrClick?: () => void;
}

export default function AddressChip({ address, onQrClick }: Props) {
  const [copied, setCopied] = useState(false);

  const copy = async (e: React.MouseEvent) => {
    e.stopPropagation();
    await navigator.clipboard.writeText(address);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  return (
    <button className="addr-chip" onClick={onQrClick}>
      <QrCode size={12} color="var(--text-muted)" />
      <span className="addr-mono">{truncateAddr(address)}</span>
      <span onClick={copy}>
        {copied ? <Check size={12} color="var(--success)" /> : <Copy size={12} color="var(--text-muted)" />}
      </span>
    </button>
  );
}
