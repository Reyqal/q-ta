import { useState, useEffect, useCallback } from 'react';
import { Wallet, Landmark, DoorOpen, DoorClosed, AlertCircle } from 'lucide-react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend } from 'recharts';
import { StatCard } from '../../components/StatCard';
import { StatusBadge } from '../../components/StatusBadge';
import { LoadingSpinner } from '../../components/LoadingSpinner';
import { formatRupiah } from '../../lib/formatCurrency';
import apiClient from '../../lib/apiClient';

interface DashboardSummary {
  total_pendapatan_bulan_ini: number;
  total_pajak_bulan_ini: number;
  total_kamar: number;
  kamar_terisi: number;
  kamar_kosong: number;
  tagihan_stats: {
    lunas: number;
    belum_bayar: number;
    ada_kendala: number;
  };
  pendapatan_bulanan: Array<{
    bulan: string;
    pendapatan: number;
    pajak: number;
  }>;
  tagihan_terbaru: Array<{
    id: number;
    period: string;
    amount: number;
    tax_portion: number;
    net_portion: number;
    status: 'paid' | 'unpaid' | 'issue';
    due_date: string;
    tenant?: {
      user?: { name: string; };
      room?: { room_number: string; };
    };
  }>;
}

export function DashboardPage() {
  const [summary, setSummary] = useState<DashboardSummary | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const fetchDashboard = useCallback(async () => {
    try {
      const response = await apiClient.get('/dashboard/summary');
      if (response.data.success) {
        setSummary(response.data.data);
      }
    } catch {
      // Dashboard data unavailable
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchDashboard();
    const interval = setInterval(fetchDashboard, 30000);
    return () => clearInterval(interval);
  }, [fetchDashboard]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-32">
        <LoadingSpinner size="lg" text="Memuat dashboard..." />
      </div>
    );
  }

  if (!summary) {
    return (
      <div className="flex flex-col items-center justify-center py-32 text-slate-400">
        <AlertCircle className="w-12 h-12 mb-4 text-slate-600" />
        <p>Gagal memuat data dashboard</p>
        <button
          onClick={fetchDashboard}
          className="mt-4 px-4 py-2 rounded-xl glass hover:bg-white/10 text-sm transition-all"
        >
          Coba Lagi
        </button>
      </div>
    );
  }

  const pendapatanBersih = summary.total_pendapatan_bulan_ini - summary.total_pajak_bulan_ini;

  return (
    <div className="space-y-8 animate-fadeIn">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-white">Dashboard</h1>
        <p className="text-sm text-slate-400 mt-1">Ringkasan keuangan dan operasional kos Anda</p>
      </div>

      {/* Stat Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <StatCard
          title="Pendapatan Bersih"
          value={formatRupiah(pendapatanBersih)}
          subtitle="Bulan ini (setelah pajak)"
          icon={Wallet}
          iconColor="text-emerald-400"
          delay={0}
        />
        <StatCard
          title="Titipan Pajak"
          value={formatRupiah(summary.total_pajak_bulan_ini)}
          subtitle="Bulan ini"
          icon={Landmark}
          iconColor="text-blue-400"
          delay={1}
        />
        <StatCard
          title="Kamar Terisi"
          value={`${summary.kamar_terisi} / ${summary.total_kamar}`}
          subtitle={`${summary.total_kamar > 0 ? Math.round((summary.kamar_terisi / summary.total_kamar) * 100) : 0}% okupansi`}
          icon={DoorClosed}
          iconColor="text-amber-400"
          delay={2}
        />
        <StatCard
          title="Kamar Kosong"
          value={String(summary.kamar_kosong)}
          subtitle="Siap dihuni"
          icon={DoorOpen}
          iconColor="text-purple-400"
          delay={3}
        />
      </div>

      {/* Chart + Invoice Stats */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Chart */}
        <div className="lg:col-span-2 glass rounded-2xl p-6 animate-fadeIn" style={{ animationDelay: '0.2s' }}>
          <h3 className="text-lg font-semibold text-white mb-6">Pendapatan vs Pajak (6 Bulan)</h3>
          <div className="h-72">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart
                data={summary.pendapatan_bulanan}
                margin={{ top: 5, right: 10, left: 10, bottom: 5 }}
              >
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.05)" />
                <XAxis
                  dataKey="bulan"
                  tick={{ fill: '#94a3b8', fontSize: 12 }}
                  tickLine={false}
                  axisLine={{ stroke: 'rgba(255,255,255,0.1)' }}
                />
                <YAxis
                  tick={{ fill: '#94a3b8', fontSize: 12 }}
                  tickLine={false}
                  axisLine={{ stroke: 'rgba(255,255,255,0.1)' }}
                  tickFormatter={(value: number) =>
                    value >= 1_000_000 ? `${(value / 1_000_000).toFixed(0)}Jt` : `${(value / 1_000).toFixed(0)}Rb`
                  }
                />
                <Tooltip
                  contentStyle={{
                    backgroundColor: '#1e293b',
                    border: '1px solid rgba(255,255,255,0.1)',
                    borderRadius: '12px',
                    color: '#e2e8f0',
                    fontSize: '13px',
                  }}
                  formatter={(value: number, name: string) => [
                    formatRupiah(value),
                    name === 'pendapatan' ? 'Pendapatan' : 'Pajak',
                  ]}
                  cursor={{ fill: 'rgba(255,255,255,0.03)' }}
                />
                <Legend
                  formatter={(value: string) => (
                    <span className="text-xs text-slate-400">
                      {value === 'pendapatan' ? 'Pendapatan' : 'Pajak'}
                    </span>
                  )}
                />
                <Bar
                  dataKey="pendapatan"
                  fill="#10b981"
                  radius={[6, 6, 0, 0]}
                  maxBarSize={40}
                />
                <Bar
                  dataKey="pajak"
                  fill="#3b82f6"
                  radius={[6, 6, 0, 0]}
                  maxBarSize={40}
                />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Invoice Stats */}
        <div className="glass rounded-2xl p-6 animate-fadeIn" style={{ animationDelay: '0.3s' }}>
          <h3 className="text-lg font-semibold text-white mb-6">Status Tagihan</h3>
          <div className="space-y-4">
            <div className="flex items-center justify-between p-4 rounded-xl bg-emerald-500/5 border border-emerald-500/10">
              <div className="flex items-center gap-3">
                <div className="w-3 h-3 rounded-full bg-emerald-500" />
                <span className="text-sm text-slate-300">Lunas</span>
              </div>
              <span className="text-2xl font-bold text-emerald-400">{summary.tagihan_stats.lunas}</span>
            </div>
            <div className="flex items-center justify-between p-4 rounded-xl bg-amber-500/5 border border-amber-500/10">
              <div className="flex items-center gap-3">
                <div className="w-3 h-3 rounded-full bg-amber-500" />
                <span className="text-sm text-slate-300">Belum Bayar</span>
              </div>
              <span className="text-2xl font-bold text-amber-400">{summary.tagihan_stats.belum_bayar}</span>
            </div>
            <div className="flex items-center justify-between p-4 rounded-xl bg-rose-500/5 border border-rose-500/10">
              <div className="flex items-center gap-3">
                <div className="w-3 h-3 rounded-full bg-rose-500" />
                <span className="text-sm text-slate-300">Ada Kendala</span>
              </div>
              <span className="text-2xl font-bold text-rose-400">{summary.tagihan_stats.ada_kendala}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Recent Invoices Table */}
      <div className="glass rounded-2xl overflow-hidden animate-fadeIn" style={{ animationDelay: '0.4s' }}>
        <div className="p-6 border-b border-white/5">
          <h3 className="text-lg font-semibold text-white">Tagihan Terbaru</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-white/5">
                <th className="text-left text-xs font-medium text-slate-500 uppercase tracking-wider px-6 py-3">Penghuni</th>
                <th className="text-left text-xs font-medium text-slate-500 uppercase tracking-wider px-6 py-3">Kamar</th>
                <th className="text-left text-xs font-medium text-slate-500 uppercase tracking-wider px-6 py-3">Periode</th>
                <th className="text-left text-xs font-medium text-slate-500 uppercase tracking-wider px-6 py-3">Jumlah</th>
                <th className="text-left text-xs font-medium text-slate-500 uppercase tracking-wider px-6 py-3">Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {summary.tagihan_terbaru.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-6 py-10 text-center text-slate-500 text-sm">
                    Belum ada tagihan
                  </td>
                </tr>
              ) : (
                summary.tagihan_terbaru.map((invoice) => (
                  <tr key={invoice.id} className="hover:bg-white/[0.02] transition-colors">
                    <td className="px-6 py-4 text-sm text-slate-200">{invoice.tenant?.user?.name || '-'}</td>
                    <td className="px-6 py-4 text-sm text-slate-400">{invoice.tenant?.room?.room_number || '-'}</td>
                    <td className="px-6 py-4 text-sm text-slate-400">{invoice.period}</td>
                    <td className="px-6 py-4 text-sm text-white font-medium">{formatRupiah(invoice.amount)}</td>
                    <td className="px-6 py-4">
                      <StatusBadge status={invoice.status} />
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Disclaimer */}
      <div className="glass rounded-2xl p-4 border-amber-500/20 bg-amber-500/5">
        <p className="text-xs text-amber-400/80 leading-relaxed">
          <strong>Disclaimer:</strong> Pemisahan pajak ini hanya alat bantu pencatatan, bukan pengganti konsultan pajak.
          Pastikan untuk berkonsultasi dengan profesional pajak untuk pelaporan yang akurat.
        </p>
      </div>
    </div>
  );
}
