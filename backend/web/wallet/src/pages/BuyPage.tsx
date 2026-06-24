import { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, ArrowRight, Check, Lock } from 'lucide-react';
import { api } from '../lib/api';
import { extractApiError } from '../lib/errors';
import AppShell from '../components/AppShell';
import Toast from '../components/Toast';

type Stage = 'pick' | 'card' | 'processing' | 'success';

interface Package {
  code: string;
  title: string;
  suda_amount: number;
  suda_amount_wei: string;
  fiat_price: string;
  fiat_currency: string;
}

function ProcessingStep({ text, done, loading }: { text: string; done?: boolean; loading?: boolean }) {
  return (
    <div className="ps-row">
      <span className={`ps-dot ${done ? 'done' : loading ? 'loading' : ''}`}>
        {done && <Check size={12} color="#0F1020" strokeWidth={3} />}
      </span>
      <span>{text}</span>
    </div>
  );
}

export default function BuyPage() {
  const navigate = useNavigate();
  const [stage, setStage] = useState<Stage>('pick');
  const [packages, setPackages] = useState<Package[]>([]);
  const [selectedPkg, setSelectedPkg] = useState<Package | null>(null);
  const [purchaseId, setPurchaseId] = useState('');
  const [card, setCard] = useState('');
  const [cvv, setCvv] = useState('');
  const [processingStep, setProcessingStep] = useState(0);
  const [toast, setToast] = useState<string | null>(null);
  const [loadingPackages, setLoadingPackages] = useState(true);
  const [packagesError, setPackagesError] = useState(false);

  const loadPackages = useCallback(async () => {
    setLoadingPackages(true);
    setPackagesError(false);
    try {
      const { data } = await api.get<{ packages: Package[] }>('/tx/purchase/packages');
      setPackages(data.packages);
    } catch (err) {
      setPackagesError(true);
      setToast(extractApiError(err, 'Failed to load packages'));
    } finally {
      setLoadingPackages(false);
    }
  }, []);

  useEffect(() => {
    void loadPackages();
  }, [loadPackages]);

  const initiate = async () => {
    if (!selectedPkg) return;
    setStage('card');
    try {
      const { data } = await api.post<{ purchase_id: string }>('/tx/purchase/initiate', {
        package_code: selectedPkg.code,
        payment_method: 'CARD',
      });
      setPurchaseId(data.purchase_id);
    } catch (err) {
      setToast(extractApiError(err, 'Failed to initiate purchase'));
      setStage('pick');
    }
  };

  const confirm = async () => {
    setStage('processing');
    setProcessingStep(0);
    const t1 = setTimeout(() => setProcessingStep(1), 700);
    const t2 = setTimeout(() => setProcessingStep(2), 1400);
    try {
      await api.post(`/tx/purchase/${purchaseId}/confirm`, {
        card_number: card.replace(/\s/g, ''),
        cvv,
      });
      clearTimeout(t1);
      clearTimeout(t2);
      setProcessingStep(3);
      setTimeout(() => setStage('success'), 400);
    } catch (err) {
      clearTimeout(t1);
      clearTimeout(t2);
      setToast(extractApiError(err, 'Payment failed'));
      setStage('card');
    }
  };

  const formatCard = (v: string) =>
    v.replace(/\D/g, '').slice(0, 16).replace(/(.{4})/g, '$1 ').trim();

  return (
    <AppShell>
      <div className="wv-app-bar">
        <button
          className="btn btn-ghost btn-sm"
          style={{ padding: '8px', borderRadius: '50%', border: 0 }}
          onClick={() => {
            if (stage === 'pick') navigate(-1);
            else if (stage === 'card') setStage('pick');
          }}
        >
          <ArrowLeft size={20} />
        </button>
        <h3>
          {stage === 'pick' ? 'Buy SUDA' : stage === 'card' ? 'Payment' : stage === 'processing' ? 'Processing' : 'Done'}
        </h3>
        <span />
      </div>

      {stage === 'pick' && (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px', padding: '16px 14px' }}>
          <div style={{ textAlign: 'center', padding: '8px 4px 0' }}>
            <h4 style={{ margin: '0 0 6px', fontSize: '20px', fontWeight: 700 }}>Top up your wallet</h4>
            <p style={{ margin: 0, color: 'var(--text-muted)', fontSize: '14px', lineHeight: 1.6 }}>
              SUDA is the in-app token. Use it to subscribe to channels, send tips, and donate.
            </p>
          </div>

          {loadingPackages && packages.length === 0 && (
            <div style={{ display: 'flex', justifyContent: 'center', padding: '32px 0' }}>
              <div className="spinner" />
            </div>
          )}

          {!loadingPackages && packagesError && packages.length === 0 && (
            <div style={{ textAlign: 'center', padding: '24px 0', display: 'flex', flexDirection: 'column', gap: '12px' }}>
              <span style={{ color: 'var(--text-muted)', fontSize: '14px' }}>
                Couldn&apos;t load packages.
              </span>
              <button className="btn btn-ghost btn-sm" style={{ alignSelf: 'center' }} onClick={() => void loadPackages()}>
                Retry
              </button>
            </div>
          )}

          <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
            {packages.map((pkg) => (
              <button
                key={pkg.code}
                className={`package-card ${selectedPkg?.code === pkg.code ? 'active' : ''}`}
                onClick={() => setSelectedPkg(pkg)}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                  <div
                    style={{
                      width: '36px',
                      height: '36px',
                      borderRadius: '50%',
                      background: 'linear-gradient(135deg, var(--primary), var(--accent))',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      fontWeight: 800,
                      fontSize: '14px',
                      color: '#0F1020',
                      flexShrink: 0,
                    }}
                  >
                    S
                  </div>
                  <div style={{ textAlign: 'left' }}>
                    <div style={{ fontSize: '18px', fontWeight: 700, fontFamily: 'var(--ff-display)' }}>
                      {pkg.suda_amount.toLocaleString()} SUDA
                    </div>
                    <div style={{ fontSize: '12px', color: 'var(--text-muted)' }}>{pkg.title}</div>
                  </div>
                </div>
                <div style={{ textAlign: 'right' }}>
                  <div style={{ fontSize: '18px', fontWeight: 700, fontFamily: 'var(--ff-display)' }}>
                    ${pkg.fiat_price}
                  </div>
                  <div style={{ fontSize: '11px', color: 'var(--text-muted)' }}>{pkg.fiat_currency}</div>
                </div>
              </button>
            ))}
          </div>

          <button
            className="btn btn-primary btn-lg btn-full"
            disabled={!selectedPkg}
            onClick={() => void initiate()}
          >
            Continue {selectedPkg ? `· $${selectedPkg.fiat_price}` : ''}
            <ArrowRight size={18} />
          </button>
          <div style={{ textAlign: 'center', fontSize: '12px', color: 'var(--text-muted)' }}>
            Demo — no real card is charged.
          </div>
        </div>
      )}

      {stage === 'card' && selectedPkg && (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px', padding: '16px 14px' }}>
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              padding: '12px 14px',
              background: 'var(--surface)',
              border: '1px solid var(--border)',
              borderRadius: '14px',
            }}
          >
            <div>
              <div style={{ fontWeight: 600 }}>{selectedPkg.suda_amount.toLocaleString()} SUDA</div>
              <div style={{ fontSize: '12px', color: 'var(--text-muted)' }}>paid in USD</div>
            </div>
            <div style={{ fontSize: '18px', fontWeight: 700, fontFamily: 'var(--ff-display)' }}>
              ${selectedPkg.fiat_price}
            </div>
          </div>

          <div className="credit-card">
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '14px', fontSize: '16px', fontWeight: 700 }}>
              <span>Suda Pay</span>
              <span>✦</span>
            </div>
            <div className="cc-number">
              {formatCard(card).padEnd(19, '•').slice(0, 19)}
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '12px', gap: '20px' }}>
              <div>
                <small style={{ fontSize: '9px', opacity: 0.7, textTransform: 'uppercase', letterSpacing: '0.12em', display: 'block' }}>HOLDER</small>
                <div style={{ fontSize: '15px', fontWeight: 700 }}>YOU</div>
              </div>
              <div>
                <small style={{ fontSize: '9px', opacity: 0.7, textTransform: 'uppercase', letterSpacing: '0.12em', display: 'block' }}>CVV</small>
                <div style={{ fontSize: '15px', fontWeight: 700 }}>{cvv.replace(/./g, '•') || '•••'}</div>
              </div>
            </div>
          </div>

          <div className="suda-field">
            <div className="suda-field-label">CARD NUMBER</div>
            <div className="suda-field-wrap">
              <input
                type="text"
                inputMode="numeric"
                placeholder="1234 5678 9012 3456"
                value={formatCard(card)}
                onChange={(e) => setCard(e.target.value.replace(/\s/g, '').slice(0, 16))}
              />
            </div>
          </div>

          <div style={{ display: 'flex', gap: '12px' }}>
            <div className="suda-field" style={{ flex: 1 }}>
              <div className="suda-field-label">CVV</div>
              <div className="suda-field-wrap">
                <input
                  type="text"
                  inputMode="numeric"
                  placeholder="123"
                  value={cvv}
                  onChange={(e) => setCvv(e.target.value.replace(/\D/g, '').slice(0, 3))}
                />
              </div>
            </div>
          </div>

          <button
            className="btn btn-primary btn-lg btn-full"
            disabled={card.length < 12}
            onClick={() => void confirm()}
          >
            <Lock size={16} />
            Pay ${selectedPkg.fiat_price}
          </button>
          <div style={{ textAlign: 'center', fontSize: '12px', color: 'var(--text-muted)' }}>
            Secured by Suda Pay · Demo only
          </div>
        </div>
      )}

      {stage === 'processing' && (
        <div className="state-center">
          <div className="spinner" />
          <h4>Processing payment…</h4>
          <p>Crediting your wallet on-chain.</p>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '10px', marginTop: '8px', width: '100%', maxWidth: '260px' }}>
            <ProcessingStep text="Card authorized" done={processingStep >= 1} loading={processingStep === 0} />
            <ProcessingStep text="Minting tokens" done={processingStep >= 2} loading={processingStep === 1} />
            <ProcessingStep text="Sending to your wallet" done={processingStep >= 3} loading={processingStep === 2} />
          </div>
        </div>
      )}

      {stage === 'success' && selectedPkg && (
        <div className="state-center">
          <div className="check-circle">
            <Check size={28} color="#0F1020" strokeWidth={3} />
          </div>
          <h4>+{selectedPkg.suda_amount.toLocaleString()} SUDA</h4>
          <p>Added to your wallet.</p>
          <div
            style={{
              background: 'var(--surface)',
              border: '1px solid var(--border)',
              borderRadius: '14px',
              padding: '14px',
              width: '100%',
              display: 'flex',
              flexDirection: 'column',
              gap: '10px',
            }}
          >
            <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '14px' }}>
              <span style={{ color: 'var(--text-muted)' }}>Package</span>
              <strong>{selectedPkg.suda_amount.toLocaleString()} SUDA</strong>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '14px' }}>
              <span style={{ color: 'var(--text-muted)' }}>Paid</span>
              <strong>${selectedPkg.fiat_price}</strong>
            </div>
          </div>
          <button
            className="btn btn-primary btn-lg btn-full"
            onClick={() => navigate('/')}
          >
            Back to wallet
          </button>
        </div>
      )}

      {toast && <Toast message={toast} variant="error" onDismiss={() => setToast(null)} />}
    </AppShell>
  );
}
