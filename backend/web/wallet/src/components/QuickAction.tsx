import { type LucideIcon } from 'lucide-react';

interface Props {
  icon: LucideIcon;
  label: string;
  onClick: () => void;
}

export default function QuickAction({ icon: Icon, label, onClick }: Props) {
  return (
    <button className="wq-btn" onClick={onClick}>
      <span className="wq-ico">
        <Icon size={20} color="#0F1020" strokeWidth={2.5} />
      </span>
      <span>{label}</span>
    </button>
  );
}
