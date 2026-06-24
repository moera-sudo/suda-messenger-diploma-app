import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Copy, Check, Share2 } from 'lucide-react';
import { useWalletStore } from '../stores/walletStore';
import { useAuthStore } from '../stores/authStore';
import AppShell from '../components/AppShell';

export default function ReceivePage() {
  const navigate = useNavigate();
  const address = useWalletStore((s) => s.address);
  const username = useAuthStore((s) => s.username);
  const [copied, setCopied] = useState(false);

  const copy = async () => {
    await navigator.clipboard.writeText(address);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const share = async () => {
    if (navigator.share) {
      await navigator.share({ title: 'My Suda Wallet', text: address });
    }
  };

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
        <h3>Receive SUDA</h3>
        {typeof navigator.share === 'function' && (
          <button
            className="btn btn-ghost btn-sm"
            style={{ padding: '8px', borderRadius: '50%', border: 0 }}
            onClick={() => void share()}
          >
            <Share2 size={20} />
          </button>
        )}
      </div>

      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          gap: '20px',
          padding: '20px 16px',
        }}
      >
        <div className="qr-card">
          {username && (
            <div style={{ fontSize: '17px', fontWeight: 700 }}>{username}</div>
          )}
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              justifyContent: 'center',
            }}
          >
            <span
              className="addr-mono"
              style={{ fontSize: '12px', wordBreak: 'break-all', textAlign: 'center' }}
            >
              {address}
            </span>
          </div>
          <button
            className="btn btn-primary btn-md btn-full"
            style={{ marginTop: '4px' }}
            onClick={() => void copy()}
          >
            {copied ? <Check size={16} /> : <Copy size={16} />}
            {copied ? 'Copied!' : 'Copy address'}
          </button>
        </div>

        <div style={{ width: '100%', maxWidth: '340px' }}>
          <div
            style={{
              fontSize: '11px',
              fontWeight: 700,
              textTransform: 'uppercase',
              letterSpacing: '0.1em',
              color: 'var(--text-muted)',
              marginBottom: '8px',
            }}
          >
            HOW IT WORKS
          </div>
          <ol
            style={{
              margin: 0,
              paddingLeft: '18px',
              display: 'flex',
              flexDirection: 'column',
              gap: '8px',
            }}
          >
            <li style={{ fontSize: '13px', color: 'var(--text-muted)', lineHeight: 1.5 }}>
              Copy your address and share it with anyone in Suda.
            </li>
            <li style={{ fontSize: '13px', color: 'var(--text-muted)', lineHeight: 1.5 }}>
              They send SUDA to it from Send → @username or address.
            </li>
            <li style={{ fontSize: '13px', color: 'var(--text-muted)', lineHeight: 1.5 }}>
              SUDA lands in your wallet instantly.
            </li>
          </ol>
        </div>
      </div>
    </AppShell>
  );
}
