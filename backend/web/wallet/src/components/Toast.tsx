import { useEffect } from 'react';
import { CheckCircle, XCircle, Info } from 'lucide-react';

type ToastVariant = 'success' | 'error' | 'info';

interface Props {
  message: string;
  variant?: ToastVariant;
  onDismiss: () => void;
  duration?: number;
}

const ICONS: Record<ToastVariant, typeof CheckCircle> = {
  success: CheckCircle,
  error: XCircle,
  info: Info,
};

export default function Toast({ message, variant = 'info', onDismiss, duration = 3000 }: Props) {
  useEffect(() => {
    const t = setTimeout(onDismiss, duration);
    return () => clearTimeout(t);
  }, [onDismiss, duration]);

  const Icon = ICONS[variant];

  return (
    <div className={`toast-wrap toast-${variant}`}>
      <Icon size={18} />
      <span style={{ flex: 1, fontSize: '14px' }}>{message}</span>
    </div>
  );
}
