import { CheckCircle, Clock, AlertTriangle } from 'lucide-react';

type StatusType = 'paid' | 'unpaid' | 'issue' | 'available' | 'occupied' | 'sent' | 'failed' | 'pending';

interface StatusBadgeProps {
  status: StatusType;
  className?: string;
}

const statusConfig: Record<StatusType, { label: string; className: string; icon: React.ReactNode }> = {
  paid: {
    label: 'Lunas',
    className: 'bg-emerald-500/15 text-emerald-400 border-emerald-500/25',
    icon: <CheckCircle className="w-3.5 h-3.5" />,
  },
  unpaid: {
    label: 'Belum Bayar',
    className: 'bg-amber-500/15 text-amber-400 border-amber-500/25',
    icon: <Clock className="w-3.5 h-3.5" />,
  },
  issue: {
    label: 'Ada Kendala',
    className: 'bg-rose-500/15 text-rose-400 border-rose-500/25',
    icon: <AlertTriangle className="w-3.5 h-3.5" />,
  },
  available: {
    label: 'Tersedia',
    className: 'bg-emerald-500/15 text-emerald-400 border-emerald-500/25',
    icon: <CheckCircle className="w-3.5 h-3.5" />,
  },
  occupied: {
    label: 'Terisi',
    className: 'bg-blue-500/15 text-blue-400 border-blue-500/25',
    icon: <Clock className="w-3.5 h-3.5" />,
  },
  sent: {
    label: 'Terkirim',
    className: 'bg-emerald-500/15 text-emerald-400 border-emerald-500/25',
    icon: <CheckCircle className="w-3.5 h-3.5" />,
  },
  failed: {
    label: 'Gagal',
    className: 'bg-rose-500/15 text-rose-400 border-rose-500/25',
    icon: <AlertTriangle className="w-3.5 h-3.5" />,
  },
  pending: {
    label: 'Menunggu',
    className: 'bg-amber-500/15 text-amber-400 border-amber-500/25',
    icon: <Clock className="w-3.5 h-3.5" />,
  },
};

export function StatusBadge({ status, className = '' }: StatusBadgeProps) {
  const config = statusConfig[status] || statusConfig.pending;

  return (
    <span
      className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border transition-all ${config.className} ${className}`}
    >
      {config.icon}
      {config.label}
    </span>
  );
}
