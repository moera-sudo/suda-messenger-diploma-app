import { type ReactNode } from 'react';

interface Props {
  onClose: () => void;
  children: ReactNode;
}

export default function BottomSheet({ onClose, children }: Props) {
  return (
    <div className="sheet-overlay" onClick={onClose}>
      <div className="sheet-dim" />
      <div className="sheet-body" onClick={(e) => e.stopPropagation()}>
        <div className="sheet-handle" />
        {children}
      </div>
    </div>
  );
}
