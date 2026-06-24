import { type ReactNode } from 'react';

interface Props {
  children: ReactNode;
}

export default function AppShell({ children }: Props) {
  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        minHeight: '100dvh',
        background: 'var(--background)',
        maxWidth: '480px',
        margin: '0 auto',
        position: 'relative',
      }}
    >
      {children}
    </div>
  );
}
