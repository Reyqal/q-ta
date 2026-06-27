import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { CheckCircle, Loader2, ArrowLeft, Zap } from 'lucide-react';
import apiClient from '../../lib/apiClient';
import { formatRupiah } from '../../lib/formatCurrency';
import { LoadingSpinner } from '../../components/LoadingSpinner';

interface Invoice {
  id: number;
  period: string;
  amount: number;
  status: string;
}

interface QRISData {
  order_id: string;
  qr_string: string;
  qr_image_url: string;
  expires_at: string;
}

export function PaymentPage() {
  const { invoiceId } = useParams<{ invoiceId: string }>();
  const navigate = useNavigate();
  const [invoice, setInvoice] = useState<Invoice | null>(null);
  const [qris, setQris] = useState<QRISData | null>(null);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [isPaid, setIsPaid] = useState(false);
  const [simulating, setSimulating] = useState(false);

  // Fetch invoice detail
  useEffect(() => {
    const fetchInvoice = async () => {
      try {
        const res = await apiClient.get(`/invoices/${invoiceId}`);
        if (res.data.success) {
          setInvoice(res.data.data);
          if (res.data.data.status === 'paid') setIsPaid(true);
        }
      } catch (e) {
        console.error(e);
        navigate('/tenant/invoices');
      } finally { setLoading(false); }
    };
    fetchInvoice();
  }, [invoiceId, navigate]);

  // Create QRIS payment
  const createQRIS = async () => {
    setCreating(true);
    try {
      const res = await apiClient.post('/payments/create-qris', { invoice_id: Number(invoiceId) });
      if (res.data.success) setQris(res.data.data);
    } catch (e: any) {
      alert(e.response?.data?.message || 'Gagal membuat pembayaran QRIS');
    } finally { setCreating(false); }
  };

  // Poll payment status
  const checkStatus = useCallback(async () => {
    if (isPaid || !qris) return;
    try {
      const res = await apiClient.get(`/payments/${invoiceId}/status`);
      if (res.data.data?.invoice_status === 'paid') {
        setIsPaid(true);
      }
    } catch (e) { console.error(e); }
  }, [invoiceId, isPaid, qris]);

  useEffect(() => {
    if (!qris || isPaid) return;
    const interval = setInterval(checkStatus, 5000);
    return () => clearInterval(interval);
  }, [qris, isPaid, checkStatus]);

  // Simulate payment
  const simulatePayment = async () => {
    if (!qris) return;
    setSimulating(true);
    try {
      await apiClient.post('/webhooks/midtrans/simulate', { order_id: qris.order_id });
      // Status will update via polling
      setTimeout(checkStatus, 1000);
    } catch (e: any) {
      alert(e.response?.data?.message || 'Gagal simulasi');
    } finally { setSimulating(false); }
  };

  if (loading) return <LoadingSpinner />;

  // Success screen
  if (isPaid) {
    return (
      <div className="animate-fade-in max-w-md mx-auto text-center py-12">
        <div className="w-24 h-24 mx-auto mb-6 rounded-full bg-emerald-500/20 flex items-center justify-center animate-pulse">
          <CheckCircle size={48} className="text-emerald-400" />
        </div>
        <h1 className="text-3xl font-bold text-white mb-2">Pembayaran Berhasil! 🎉</h1>
        <p className="text-slate-400 mb-2">Tagihan periode {invoice?.period} telah lunas.</p>
        <p className="text-2xl font-bold text-emerald-400 mb-8">{formatRupiah(invoice?.amount || 0)}</p>
        <button onClick={() => navigate('/tenant/invoices')} className="px-6 py-3 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl font-medium transition-all">
          Kembali ke Daftar Tagihan
        </button>
      </div>
    );
  }

  return (
    <div className="animate-fade-in max-w-md mx-auto">
      <button onClick={() => navigate('/tenant/invoices')} className="flex items-center gap-2 text-slate-400 hover:text-white mb-6 transition-colors">
        <ArrowLeft size={18} /> Kembali
      </button>

      {/* Invoice Summary */}
      <div className="glass rounded-2xl p-6 mb-6">
        <h2 className="text-lg font-bold text-white mb-4">Rincian Pembayaran</h2>
        <div className="space-y-3">
          <div className="flex justify-between">
            <span className="text-slate-400">Periode</span>
            <span className="text-white font-medium">{invoice?.period}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-400">Total Tagihan</span>
            <span className="text-xl font-bold text-emerald-400">{formatRupiah(invoice?.amount || 0)}</span>
          </div>
        </div>
      </div>

      {/* QR Code Area */}
      {!qris ? (
        <button
          onClick={createQRIS}
          disabled={creating}
          className="w-full py-4 bg-gradient-to-r from-blue-600 to-blue-500 hover:from-blue-500 hover:to-blue-400 text-white rounded-2xl font-semibold text-lg transition-all duration-300 shadow-lg shadow-blue-900/40 disabled:opacity-50 flex items-center justify-center gap-2"
        >
          {creating ? <><Loader2 size={20} className="animate-spin" /> Membuat QRIS...</> : 'Buat Pembayaran QRIS'}
        </button>
      ) : (
        <div className="glass rounded-2xl p-6 text-center">
          <h3 className="text-sm font-semibold text-slate-300 mb-4">Scan QR Code untuk Membayar</h3>

          <div className="bg-white rounded-xl p-4 inline-block mb-4">
            <img
              src={qris.qr_image_url}
              alt="QR Code Pembayaran"
              className="w-56 h-56 object-contain"
              onError={(e) => {
                (e.target as HTMLImageElement).src = `https://api.qrserver.com/v1/create-qr-code/?size=224x224&data=${encodeURIComponent(qris.qr_string)}`;
              }}
            />
          </div>

          <div className="flex items-center justify-center gap-2 mb-4 text-amber-400">
            <Loader2 size={16} className="animate-spin" />
            <span className="text-sm">Menunggu pembayaran...</span>
          </div>

          <p className="text-xs text-slate-500 mb-4">
            Order ID: <code className="text-slate-400">{qris.order_id}</code>
          </p>

          {/* Dev Mode: Simulate Button */}
          <div className="border-t border-white/10 pt-4 mt-4">
            <p className="text-xs text-slate-500 mb-2">🧪 Mode Demo</p>
            <button
              onClick={simulatePayment}
              disabled={simulating}
              className="flex items-center justify-center gap-2 w-full py-2.5 bg-amber-600/20 hover:bg-amber-600/30 text-amber-400 rounded-xl text-sm font-medium transition-all border border-amber-500/20 disabled:opacity-50"
            >
              <Zap size={16} />
              {simulating ? 'Mensimulasikan...' : 'Simulasi Pembayaran Berhasil'}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
