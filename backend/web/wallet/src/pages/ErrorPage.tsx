interface Props {
  code?: string;
}

const ERROR_MESSAGES: Record<string, { title: string; body: string }> = {
  NO_INIT_DATA: {
    title: 'Open from Suda app',
    body: 'This wallet can only be opened from the Suda messenger app.',
  },
  INIT_DATA_EXPIRED: {
    title: 'Session expired',
    body: 'Your session has expired. Please reopen the wallet from the app.',
  },
  SESSION_EXPIRED: {
    title: 'Session expired',
    body: 'Your session has expired. Please reopen the wallet from the app.',
  },
  NOT_FOUND: {
    title: 'Page not found',
    body: "The page you're looking for doesn't exist.",
  },
  AUTH_ERROR: {
    title: 'Authentication failed',
    body: 'Could not authenticate. Check your connection and try again, or reopen the wallet from the app.',
  },
};

export default function ErrorPage({ code = 'AUTH_ERROR' }: Props) {
  const msg = ERROR_MESSAGES[code] ?? ERROR_MESSAGES['AUTH_ERROR']!;

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '100dvh',
        padding: '32px 24px',
        textAlign: 'center',
        background: 'var(--background)',
        gap: '16px',
      }}
    >
      <div
        style={{
          width: '72px',
          height: '72px',
          borderRadius: '22px',
          background: 'var(--surface)',
          border: '1px solid var(--border)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: '32px',
        }}
      >
        🔒
      </div>
      <h2 style={{ margin: 0, fontSize: '22px', fontWeight: 700 }}>{msg.title}</h2>
      <p style={{ margin: 0, color: 'var(--text-muted)', fontSize: '14px', lineHeight: 1.6 }}>
        {msg.body}
      </p>
      <div style={{ display: 'flex', gap: '10px', marginTop: '8px' }}>
        <button
          className="btn btn-primary btn-md"
          onClick={() => window.location.reload()}
        >
          Retry
        </button>
        <button className="btn btn-ghost btn-md" onClick={() => window.close()}>
          Close
        </button>
      </div>
    </div>
  );
}
