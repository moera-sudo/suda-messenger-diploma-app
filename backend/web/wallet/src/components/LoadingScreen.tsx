export default function LoadingScreen() {
  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '100dvh',
        gap: '16px',
        background: 'var(--background)',
      }}
    >
      <div className="spinner" />
      <p style={{ color: 'var(--text-muted)', fontSize: '14px', margin: 0 }}>
        Connecting wallet…
      </p>
    </div>
  );
}
