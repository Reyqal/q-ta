import { useState, useEffect } from 'react';
import { Bell } from 'lucide-react';
import apiClient from '../../lib/apiClient';
import { LoadingSpinner } from '../../components/LoadingSpinner';
import { EmptyState } from '../../components/EmptyState';

interface NotifLog {
  id: number;
  tenant_id: number | null;
  channel: string;
  message: string;
  status: string;
  created_at: string;
}

export function NotificationsPage() {
  const [logs, setLogs] = useState<NotifLog[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchLogs = async () => {
      try {
        const res = await apiClient.get('/notifications');
        if (res.data.success) setLogs(res.data.data || []);
      } catch (e) { console.error(e); }
      finally { setLoading(false); }
    };
    fetchLogs();
  }, []);

  if (loading) return <LoadingSpinner />;

  return (
    <div className="animate-fade-in">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-white">Log Notifikasi</h1>
        <p className="text-slate-400 mt-1">Riwayat pengiriman notifikasi ke penghuni (simulasi)</p>
      </div>

      {logs.length === 0 ? (
        <EmptyState icon={Bell} message="Belum ada notifikasi" />
      ) : (
        <div className="space-y-3">
          {logs.map((log) => (
            <div key={log.id} className="glass rounded-xl p-4 hover:bg-white/5 transition-colors">
              <div className="flex items-start justify-between gap-4">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="inline-flex px-2 py-0.5 rounded-full text-[10px] font-semibold bg-blue-500/20 text-blue-400 uppercase">
                      {log.channel}
                    </span>
                    <span className={`inline-flex px-2 py-0.5 rounded-full text-[10px] font-semibold ${log.status === 'simulated_sent' ? 'bg-amber-500/20 text-amber-400' : log.status === 'sent' ? 'bg-emerald-500/20 text-emerald-400' : 'bg-rose-500/20 text-rose-400'}`}>
                      {log.status === 'simulated_sent' ? 'Simulasi' : log.status === 'sent' ? 'Terkirim' : 'Gagal'}
                    </span>
                  </div>
                  <pre className="text-sm text-slate-300 whitespace-pre-wrap font-sans leading-relaxed">{log.message}</pre>
                </div>
                <span className="text-xs text-slate-500 whitespace-nowrap">
                  {new Date(log.created_at).toLocaleString('id-ID')}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
