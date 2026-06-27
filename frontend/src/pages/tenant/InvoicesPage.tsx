import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Receipt, CreditCard } from 'lucide-react';
import { useAuth } from '../../contexts/AuthContext';
import apiClient from '../../lib/apiClient';
import { formatRupiah } from '../../lib/formatCurrency';
import { StatusBadge } from '../../components/StatusBadge';
import { LoadingSpinner } from '../../components/LoadingSpinner';
import { EmptyState } from '../../components/EmptyState';

interface Invoice {
  id: number;
  period: string;
  amount: number;
  tax_portion: number;
  net_portion: number;
  status: 'unpaid' | 'paid' | 'issue';
  due_date: string;
  paid_at: string | null;
}

export function InvoicesPage() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchInvoices = async () => {
    try {
      const res = await apiClient.get('/invoices');
      if (res.data.success) setInvoices(res.data.data || []);
    } catch (e) { console.error(e); }
    finally { setLoading(false); }
  };

  useEffect(() => {
    fetchInvoices();
    const interval = setInterval(fetchInvoices, 10000); // Poll setiap 10 detik
    return () => clearInterval(interval);
  }, []);

  if (loading) return <LoadingSpinner />;

  return (
    <div className="animate-fade-in max-w-3xl mx-auto">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-white">Selamat Datang, {user?.name} 👋</h1>
        <p className="text-slate-400 mt-1">Berikut tagihan kos Anda</p>
      </div>

      {invoices.length === 0 ? (
        <EmptyState icon={Receipt} message="Belum ada tagihan untuk Anda" />
      ) : (
        <div className="space-y-4">
          {invoices.map((inv) => (
            <div key={inv.id} className="glass rounded-2xl p-6 hover:border-emerald-500/30 transition-all duration-300">
              <div className="flex items-start justify-between mb-4">
                <div>
                  <p className="text-sm text-slate-400">Periode</p>
                  <p className="text-lg font-bold text-white">{inv.period}</p>
                </div>
                <StatusBadge status={inv.status} />
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <p className="text-xs text-slate-500">Total Tagihan</p>
                  <p className="text-xl font-bold text-white">{formatRupiah(inv.amount)}</p>
                </div>
                <div>
                  <p className="text-xs text-slate-500">Jatuh Tempo</p>
                  <p className="text-sm text-slate-300">{new Date(inv.due_date).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })}</p>
                </div>
              </div>

              {inv.status === 'paid' && (
                <div className="bg-emerald-500/10 rounded-xl p-3 border border-emerald-500/20">
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-emerald-400">✅ Lunas pada {inv.paid_at ? new Date(inv.paid_at).toLocaleDateString('id-ID') : '-'}</span>
                  </div>
                  <div className="flex gap-4 mt-2 text-xs">
                    <span className="text-slate-400">Bersih: <span className="text-emerald-400">{formatRupiah(inv.net_portion)}</span></span>
                    <span className="text-slate-400">Pajak: <span className="text-amber-400">{formatRupiah(inv.tax_portion)}</span></span>
                  </div>
                </div>
              )}

              {inv.status === 'issue' && (
                <div className="bg-rose-500/10 rounded-xl p-3 border border-rose-500/20">
                  <p className="text-sm text-rose-400">⚠️ Ada kendala — admin telah diberitahu. Pengingat otomatis dihentikan.</p>
                </div>
              )}

              {inv.status === 'unpaid' && (
                <button
                  onClick={() => navigate(`/tenant/payment/${inv.id}`)}
                  className="w-full mt-2 flex items-center justify-center gap-2 px-6 py-3 bg-gradient-to-r from-emerald-600 to-emerald-500 hover:from-emerald-500 hover:to-emerald-400 text-white rounded-xl font-semibold transition-all duration-300 shadow-lg shadow-emerald-900/40 hover:shadow-emerald-800/60"
                >
                  <CreditCard size={18} />
                  Bayar Sekarang
                </button>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
