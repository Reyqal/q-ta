import { useState, useEffect } from 'react';
import { FileText, CalendarPlus, AlertTriangle, CheckCircle } from 'lucide-react';
import apiClient from '../../lib/apiClient';
import { formatRupiah } from '../../lib/formatCurrency';
import { StatusBadge } from '../../components/StatusBadge';
import { LoadingSpinner } from '../../components/LoadingSpinner';
import { EmptyState } from '../../components/EmptyState';
import { Modal } from '../../components/Modal';

interface Invoice {
  id: number;
  tenant_id: number;
  period: string;
  amount: number;
  tax_portion: number;
  net_portion: number;
  status: 'unpaid' | 'paid' | 'issue';
  due_date: string;
  paid_at: string | null;
  tenant: {
    user: { name: string; phone_number: string; };
    room: { room_number: string; };
  };
}

export function InvoicesPage() {
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState(false);
  const [showManualModal, setShowManualModal] = useState(false);
  const [selectedInvoice, setSelectedInvoice] = useState<Invoice | null>(null);
  const [paymentMethod, setPaymentMethod] = useState('cash');

  const fetchInvoices = async () => {
    try {
      const res = await apiClient.get('/invoices');
      if (res.data.success) setInvoices(res.data.data || []);
    } catch (e) { console.error(e); }
    finally { setLoading(false); }
  };

  useEffect(() => { fetchInvoices(); }, []);

  const generateMonthly = async () => {
    if (!confirm('Generate tagihan bulanan untuk semua penghuni aktif?')) return;
    setGenerating(true);
    try {
      const res = await apiClient.post('/invoices/generate-monthly');
      alert(`${res.data.data?.generated_count || 0} tagihan berhasil di-generate`);
      fetchInvoices();
    } catch (e: any) {
      alert(e.response?.data?.message || 'Gagal generate');
    } finally { setGenerating(false); }
  };

  const toggleIssue = async (inv: Invoice) => {
    const newStatus = inv.status === 'issue' ? 'unpaid' : 'issue';
    try {
      await apiClient.put(`/invoices/${inv.id}/status`, { status: newStatus });
      fetchInvoices();
    } catch (e: any) {
      alert(e.response?.data?.message || 'Gagal mengubah status');
    }
  };

  const openManualConfirm = (inv: Invoice) => {
    setSelectedInvoice(inv);
    setPaymentMethod('cash');
    setShowManualModal(true);
  };

  const confirmManual = async () => {
    if (!selectedInvoice) return;
    try {
      await apiClient.put(`/invoices/${selectedInvoice.id}/confirm-manual`, { payment_method: paymentMethod });
      setShowManualModal(false);
      fetchInvoices();
    } catch (e: any) {
      alert(e.response?.data?.message || 'Gagal konfirmasi');
    }
  };

  if (loading) return <LoadingSpinner />;

  return (
    <div className="animate-fade-in">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-white">Kelola Tagihan</h1>
          <p className="text-slate-400 mt-1">Kelola tagihan semua penghuni</p>
        </div>
        <button onClick={generateMonthly} disabled={generating} className="flex items-center gap-2 px-4 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-all font-medium shadow-lg shadow-blue-900/30 disabled:opacity-50">
          <CalendarPlus size={18} /> {generating ? 'Generating...' : 'Generate Tagihan Bulanan'}
        </button>
      </div>

      <div className="glass rounded-xl p-3 mb-4 flex items-center gap-2 text-xs text-slate-400">
        <span className="inline-block w-3 h-3 rounded-full bg-emerald-500"></span> Lunas
        <span className="inline-block w-3 h-3 rounded-full bg-amber-500 ml-3"></span> Belum Bayar
        <span className="inline-block w-3 h-3 rounded-full bg-rose-500 ml-3"></span> Ada Kendala
      </div>

      {invoices.length === 0 ? (
        <EmptyState icon={FileText} message="Belum ada tagihan" />
      ) : (
        <div className="glass rounded-xl overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-white/10">
                  <th className="text-left px-4 py-3 text-xs font-semibold text-slate-300">Penghuni</th>
                  <th className="text-left px-4 py-3 text-xs font-semibold text-slate-300">Kamar</th>
                  <th className="text-left px-4 py-3 text-xs font-semibold text-slate-300">Periode</th>
                  <th className="text-right px-4 py-3 text-xs font-semibold text-slate-300">Total</th>
                  <th className="text-right px-4 py-3 text-xs font-semibold text-slate-300">Pajak</th>
                  <th className="text-right px-4 py-3 text-xs font-semibold text-slate-300">Bersih</th>
                  <th className="text-center px-4 py-3 text-xs font-semibold text-slate-300">Status</th>
                  <th className="text-right px-4 py-3 text-xs font-semibold text-slate-300">Aksi</th>
                </tr>
              </thead>
              <tbody>
                {invoices.map((inv) => (
                  <tr key={inv.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                    <td className="px-4 py-3 text-sm text-white font-medium">{inv.tenant?.user?.name || '-'}</td>
                    <td className="px-4 py-3 text-sm text-slate-300">{inv.tenant?.room?.room_number || '-'}</td>
                    <td className="px-4 py-3 text-sm text-slate-300">{inv.period}</td>
                    <td className="px-4 py-3 text-sm text-white text-right font-medium">{formatRupiah(inv.amount)}</td>
                    <td className="px-4 py-3 text-sm text-amber-400 text-right">{formatRupiah(inv.tax_portion)}</td>
                    <td className="px-4 py-3 text-sm text-emerald-400 text-right">{formatRupiah(inv.net_portion)}</td>
                    <td className="px-4 py-3 text-center"><StatusBadge status={inv.status} /></td>
                    <td className="px-4 py-3 text-right">
                      {inv.status !== 'paid' && (
                        <div className="flex items-center justify-end gap-1">
                          <button onClick={() => toggleIssue(inv)} title={inv.status === 'issue' ? 'Hapus Kendala' : 'Set Ada Kendala'} className="p-1.5 hover:bg-amber-500/20 text-amber-400 rounded-lg transition-colors">
                            <AlertTriangle size={14} />
                          </button>
                          <button onClick={() => openManualConfirm(inv)} title="Konfirmasi Bayar Manual" className="p-1.5 hover:bg-emerald-500/20 text-emerald-400 rounded-lg transition-colors">
                            <CheckCircle size={14} />
                          </button>
                        </div>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      <div className="mt-4 glass-light rounded-lg p-3">
        <p className="text-xs text-slate-500">⚠️ Pemisahan pajak (0,5% PPh Final UMKM) ini hanya alat bantu pencatatan, bukan pengganti konsultan pajak atau aplikasi resmi pemerintah.</p>
      </div>

      {/* Modal Konfirmasi Manual */}
      <Modal isOpen={showManualModal} onClose={() => setShowManualModal(false)} title="Konfirmasi Pembayaran Manual">
        <div className="space-y-4">
          {selectedInvoice && (
            <div className="glass-light rounded-lg p-4">
              <p className="text-sm text-slate-400">Tagihan: <span className="text-white font-semibold">{formatRupiah(selectedInvoice.amount)}</span></p>
              <p className="text-sm text-slate-400">Penghuni: <span className="text-white">{selectedInvoice.tenant?.user?.name}</span></p>
              <p className="text-sm text-slate-400">Periode: <span className="text-white">{selectedInvoice.period}</span></p>
            </div>
          )}
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Metode Pembayaran</label>
            <div className="flex gap-3">
              <label className="flex items-center gap-2 cursor-pointer">
                <input type="radio" value="cash" checked={paymentMethod === 'cash'} onChange={() => setPaymentMethod('cash')} className="accent-emerald-500" />
                <span className="text-sm text-slate-300">Tunai (Cash)</span>
              </label>
              <label className="flex items-center gap-2 cursor-pointer">
                <input type="radio" value="transfer" checked={paymentMethod === 'transfer'} onChange={() => setPaymentMethod('transfer')} className="accent-emerald-500" />
                <span className="text-sm text-slate-300">Transfer Bank</span>
              </label>
            </div>
          </div>
          <div className="flex justify-end gap-3 pt-2">
            <button onClick={() => setShowManualModal(false)} className="px-4 py-2.5 text-slate-400 hover:text-white transition-colors">Batal</button>
            <button onClick={confirmManual} className="px-6 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg font-medium transition-all">Konfirmasi Lunas</button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
