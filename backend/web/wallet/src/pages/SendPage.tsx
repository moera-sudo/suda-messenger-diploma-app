import { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, ArrowRight, Check } from 'lucide-react';
import { useWalletStore } from '../stores/walletStore';
import { api } from '../lib/api';
import { extractApiError } from '../lib/errors';
import { useDebounce } from '../hooks/useDebounce';
import { weiToSuda, sudaToWei } from '../lib/format';
import AppShell from '../components/AppShell';
import Toast from '../components/Toast';

type Stage = 'form' | 'confirm' | 'done';

interface ResolvedUser {
  found: boolean;
  user_id?: string;
  display_name?: string;
  wallet_address?: string;
}

export default function SendPage() {
  const navigate = useNavigate();
  const { balanceWei } = useWalletStore();
  const [stage, setStage] = useState<Stage>('form');
  const [recipient, setRecipient] = useState('');
  const [amount, setAmount] = useState('');
  const [note, setNote] = useState('');
  const [resolved, setResolved] = useState<ResolvedUser | null>(null);
  const [resolving, setResolving] = useState(false);
  const [txHash, setTxHash] = useState('');
  const [toast, setToast] = useState<string | null>(null);

  const resolveUser = useCallback(async (username: string) => {
    if (username.trim().length < 2) { setResolved(null); return; }
    setResolving(true);
    try {
      const { data } = await api.post<ResolvedUser>('/tx/wallet/resolve', { username });
      setResolved(data);
    } catch (err) {
      // Request failed (network / auth / server) — do NOT report it as
      // "user not found", which masks the real problem. Clear the resolution
      // and surface the actual error so the user knows to retry.
      setResolved(null);
      setToast(extractApiError(err, "Couldn't check username, try again"));
    } finally {
      setResolving(false);
    }
  }, []);

  const debouncedResolve = useDebounce(resolveUser, 400);

  const handleRecipientChange = (v: string) => {
    setRecipient(v);
    setResolved(null);
    debouncedResolve(v);
  };

  const amountWei = amount ? sudaToWei(amount) : '0';
  const balanceSuda = parseFloat(weiToSuda(balanceWei.toString(), 2));
  const amountNum = parseFloat(amount) || 0;
  const canSend = resolved?.found && amountNum > 0 && amountNum <= balanceSuda;

  const submit = async () => {
    setStage('confirm');
    try {
      const { data } = await api.post<{ tx_hash: string }>('/tx/wallet/transfer', {
        to_username: recipient,
        amount_wei: amountWei,
        note: note || undefined,
      });
      setTxHash(data.tx_hash);
      setStage('done');
    } catch (err) {
      setToast(extractApiError(err, 'Transfer failed'));
      setStage('form');
    }
  };

  const quickAmounts = [10, 25, 50, 100];

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
        <h3>Send SUDA</h3>
        <span />
      </div>

      {stage === 'form' && (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px', padding: '16px 14px' }}>
          <div className="suda-field">
            <div className="suda-field-label">RECIPIENT</div>
            <div className={`suda-field-wrap ${resolved?.found === false ? 'error' : ''}`}>
              <input
                type="text"
                placeholder="@username"
                value={recipient}
                onChange={(e) => handleRecipientChange(e.target.value)}
                autoFocus
              />
              {resolving && <div className="spinner" style={{ width: '16px', height: '16px', borderWidth: '2px', flexShrink: 0 }} />}
            </div>
          </div>

          {resolved?.found && (
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: '12px',
                padding: '12px 14px',
                background: 'var(--surface)',
                border: '1px solid var(--border)',
                borderRadius: '14px',
              }}
            >
              <div
                style={{
                  width: '40px',
                  height: '40px',
                  borderRadius: '50%',
                  background: 'linear-gradient(135deg, var(--primary), var(--accent))',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontWeight: 700,
                  fontSize: '16px',
                  flexShrink: 0,
                }}
              >
                {resolved.display_name?.[0]?.toUpperCase() ?? '?'}
              </div>
              <div style={{ flex: 1 }}>
                <div style={{ fontWeight: 600 }}>{resolved.display_name}</div>
                <div className="addr-mono" style={{ fontSize: '12px', marginTop: '2px' }}>
                  {resolved.wallet_address ? `${resolved.wallet_address.slice(0, 10)}…${resolved.wallet_address.slice(-6)}` : ''}
                </div>
              </div>
              <Check size={18} color="var(--success)" />
            </div>
          )}

          {resolved?.found === false && recipient.length > 1 && (
            <div style={{ fontSize: '13px', color: 'var(--danger)', paddingLeft: '4px' }}>
              User not found
            </div>
          )}

          <div>
            <div className="amount-input-wrap">
              <input
                type="text"
                inputMode="decimal"
                placeholder="0.00"
                value={amount}
                onChange={(e) => setAmount(e.target.value.replace(/[^\d.]/g, ''))}
              />
              <span className="unit">SUDA</span>
            </div>
            <div style={{ fontSize: '12px', color: 'var(--text-muted)', marginTop: '6px', paddingLeft: '4px' }}>
              Available: {weiToSuda(balanceWei.toString(), 2)} SUDA
            </div>
          </div>

          <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap' }}>
            {quickAmounts.map((n) => (
              <button
                key={n}
                className="chip"
                onClick={() => setAmount(String(n))}
              >
                {n}
              </button>
            ))}
            <button
              className="chip"
              onClick={() => setAmount(weiToSuda(balanceWei.toString(), 2))}
            >
              MAX
            </button>
          </div>

          <div className="suda-field">
            <div className="suda-field-label">NOTE (optional)</div>
            <div className="suda-field-wrap">
              <input
                type="text"
                placeholder="For the help"
                value={note}
                onChange={(e) => setNote(e.target.value.slice(0, 200))}
              />
            </div>
          </div>

          <div
            style={{
              background: 'var(--surface)',
              border: '1px solid var(--border)',
              borderRadius: '14px',
              padding: '14px',
              display: 'flex',
              flexDirection: 'column',
              gap: '10px',
            }}
          >
            <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '14px' }}>
              <span style={{ color: 'var(--text-muted)' }}>Balance after</span>
              <strong>{Math.max(0, balanceSuda - amountNum).toFixed(2)} SUDA</strong>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '14px' }}>
              <span style={{ color: 'var(--text-muted)' }}>Network fee</span>
              <strong style={{ color: 'var(--success)' }}>Free</strong>
            </div>
          </div>

          <button
            className="btn btn-primary btn-lg btn-full"
            disabled={!canSend}
            onClick={() => void submit()}
          >
            {amount ? `Send ${parseFloat(amount).toFixed(2)} SUDA` : 'Send'}
            <ArrowRight size={18} />
          </button>
        </div>
      )}

      {stage === 'confirm' && (
        <div className="state-center">
          <div className="spinner" />
          <h4>Confirming on-chain…</h4>
          <p>Don&apos;t close this window.</p>
        </div>
      )}

      {stage === 'done' && (
        <div className="state-center">
          <div className="check-circle">
            <Check size={28} color="#0F1020" strokeWidth={3} />
          </div>
          <h4>Sent successfully</h4>
          <p>{parseFloat(amount).toFixed(2)} SUDA → {recipient}</p>
          {txHash && (
            <div className="addr-mono" style={{ fontSize: '12px' }}>
              {txHash.slice(0, 10)}…{txHash.slice(-8)}
            </div>
          )}
          <button
            className="btn btn-primary btn-lg btn-full"
            onClick={() => navigate('/')}
            style={{ marginTop: '8px' }}
          >
            Back to wallet
          </button>
        </div>
      )}

      {toast && (
        <Toast message={toast} variant="error" onDismiss={() => setToast(null)} />
      )}

</AppShell>
  );
}
