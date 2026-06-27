import { useState, useEffect } from 'react';
import { UserPlus, Users, Trash2, Copy } from 'lucide-react';
import apiClient from '../../lib/apiClient';
import { formatRupiah } from '../../lib/formatCurrency';
import { Modal } from '../../components/Modal';
import { LoadingSpinner } from '../../components/LoadingSpinner';
import { EmptyState } from '../../components/EmptyState';

interface Room { id: number; room_number: string; status: string; rent_amount: number; }
interface Tenant {
  id: number;
  user_id: number;
  room_id: number;
  join_date: string;
  is_active: boolean;
  user: { id: number; name: string; phone_number: string; };
  room: { id: number; room_number: string; rent_amount: number; };
}

export function TenantsPage() {
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [rooms, setRooms] = useState<Room[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [showResult, setShowResult] = useState(false);
  const [generatedPassword, setGeneratedPassword] = useState('');
  const [form, setForm] = useState({ name: '', phone_number: '', room_id: '' });
  const [submitting, setSubmitting] = useState(false);
  const [selectedRoomAmount, setSelectedRoomAmount] = useState(0);

  const fetchData = async () => {
    try {
      const [tenantRes, roomRes] = await Promise.all([
        apiClient.get('/tenants'),
        apiClient.get('/rooms'),
      ]);
      if (tenantRes.data.success) setTenants(tenantRes.data.data || []);
      if (roomRes.data.success) setRooms(roomRes.data.data || []);
    } catch (e) { console.error(e); }
    finally { setLoading(false); }
  };

  useEffect(() => { fetchData(); }, []);

  const availableRooms = rooms.filter(r => r.status === 'available');

  const openAdd = () => {
    setForm({ name: '', phone_number: '', room_id: '' });
    setSelectedRoomAmount(0);
    setShowModal(true);
  };

  const handleRoomChange = (roomId: string) => {
    setForm({ ...form, room_id: roomId });
    const room = rooms.find(r => r.id === Number(roomId));
    setSelectedRoomAmount(room?.rent_amount || 0);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      const res = await apiClient.post('/tenants', {
        name: form.name,
        phone_number: form.phone_number,
        room_id: Number(form.room_id),
      });
      if (res.data.success) {
        setGeneratedPassword(res.data.data.generated_password);
        setShowModal(false);
        setShowResult(true);
        fetchData();
      }
    } catch (e: any) {
      alert(e.response?.data?.message || 'Gagal menambahkan penghuni');
    } finally { setSubmitting(false); }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Yakin ingin menonaktifkan penghuni ini? Kamar akan kembali tersedia.')) return;
    try {
      await apiClient.delete(`/tenants/${id}`);
      fetchData();
    } catch (e: any) {
      alert(e.response?.data?.message || 'Gagal menonaktifkan');
    }
  };

  const copyPassword = () => {
    navigator.clipboard.writeText(generatedPassword);
    alert('Password berhasil disalin!');
  };

  if (loading) return <LoadingSpinner />;

  return (
    <div className="animate-fade-in">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-white">Kelola Penghuni</h1>
          <p className="text-slate-400 mt-1">Tambah dan kelola data penghuni kos</p>
        </div>
        <button onClick={openAdd} disabled={availableRooms.length === 0} className="flex items-center gap-2 px-4 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg transition-all duration-200 font-medium shadow-lg shadow-emerald-900/30 disabled:opacity-50 disabled:cursor-not-allowed">
          <UserPlus size={18} /> Tambah Penghuni
        </button>
      </div>

      {availableRooms.length === 0 && tenants.length > 0 && (
        <div className="glass rounded-xl p-4 mb-4 border-l-4 border-amber-500">
          <p className="text-amber-400 text-sm">⚠️ Semua kamar sudah terisi. Tambahkan kamar baru untuk menambah penghuni.</p>
        </div>
      )}

      {tenants.length === 0 ? (
        <EmptyState icon={Users} message="Belum ada data penghuni" />
      ) : (
        <div className="glass rounded-xl overflow-hidden">
          <table className="w-full">
            <thead>
              <tr className="border-b border-white/10">
                <th className="text-left px-6 py-4 text-sm font-semibold text-slate-300">Nama</th>
                <th className="text-left px-6 py-4 text-sm font-semibold text-slate-300">No. WhatsApp</th>
                <th className="text-left px-6 py-4 text-sm font-semibold text-slate-300">Kamar</th>
                <th className="text-left px-6 py-4 text-sm font-semibold text-slate-300">Sewa</th>
                <th className="text-left px-6 py-4 text-sm font-semibold text-slate-300">Tanggal Masuk</th>
                <th className="text-right px-6 py-4 text-sm font-semibold text-slate-300">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {tenants.map((t) => (
                <tr key={t.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                  <td className="px-6 py-4 font-semibold text-white">{t.user.name}</td>
                  <td className="px-6 py-4 text-slate-300">{t.user.phone_number}</td>
                  <td className="px-6 py-4 text-slate-300">{t.room.room_number}</td>
                  <td className="px-6 py-4 text-slate-300">{formatRupiah(t.room.rent_amount)}</td>
                  <td className="px-6 py-4 text-slate-400 text-sm">{new Date(t.join_date).toLocaleDateString('id-ID')}</td>
                  <td className="px-6 py-4 text-right">
                    <button onClick={() => handleDelete(t.id)} className="p-2 hover:bg-rose-500/20 text-rose-400 rounded-lg transition-colors"><Trash2 size={16} /></button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Modal Tambah Penghuni */}
      <Modal isOpen={showModal} onClose={() => setShowModal(false)} title="Tambah Penghuni Baru">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">Nama Lengkap</label>
            <input type="text" value={form.name} onChange={e => setForm({...form, name: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-lg text-white focus:ring-2 focus:ring-emerald-500 outline-none" required />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">Nomor WhatsApp</label>
            <input type="tel" value={form.phone_number} onChange={e => setForm({...form, phone_number: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-lg text-white focus:ring-2 focus:ring-emerald-500 outline-none" placeholder="08xxxxxxxxxx" required />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">Pilih Kamar</label>
            <select value={form.room_id} onChange={e => handleRoomChange(e.target.value)} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-lg text-white focus:ring-2 focus:ring-emerald-500 outline-none" required>
              <option value="">-- Pilih Kamar --</option>
              {availableRooms.map(r => (
                <option key={r.id} value={r.id}>Kamar {r.room_number} — {formatRupiah(r.rent_amount)}</option>
              ))}
            </select>
          </div>
          {selectedRoomAmount > 0 && (
            <div className="glass-light rounded-lg p-3">
              <p className="text-sm text-slate-400">Nominal Sewa: <span className="text-emerald-400 font-semibold">{formatRupiah(selectedRoomAmount)}</span>/bulan</p>
            </div>
          )}
          <p className="text-xs text-slate-500">⚙️ Akun login penghuni akan dibuat otomatis. Password akan ditampilkan setelah disimpan.</p>
          <div className="flex justify-end gap-3 pt-2">
            <button type="button" onClick={() => setShowModal(false)} className="px-4 py-2.5 text-slate-400 hover:text-white transition-colors">Batal</button>
            <button type="submit" disabled={submitting} className="px-6 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg font-medium transition-all disabled:opacity-50">
              {submitting ? 'Menyimpan...' : 'Simpan'}
            </button>
          </div>
        </form>
      </Modal>

      {/* Modal Hasil */}
      <Modal isOpen={showResult} onClose={() => setShowResult(false)} title="✅ Penghuni Berhasil Ditambahkan">
        <div className="space-y-4">
          <div className="glass-light rounded-lg p-4">
            <p className="text-sm text-slate-400 mb-2">Password yang di-generate untuk penghuni:</p>
            <div className="flex items-center gap-2">
              <code className="flex-1 px-3 py-2 bg-slate-800 rounded font-mono text-emerald-400 text-lg">{generatedPassword}</code>
              <button onClick={copyPassword} className="p-2 hover:bg-emerald-500/20 text-emerald-400 rounded-lg transition-colors"><Copy size={18} /></button>
            </div>
          </div>
          <div className="glass-light rounded-lg p-4 border-l-4 border-blue-500">
            <p className="text-sm text-blue-300">📋 Kredensial login telah dicatat di log notifikasi sebagai simulasi pengiriman via WhatsApp.</p>
          </div>
          <button onClick={() => setShowResult(false)} className="w-full py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg font-medium transition-all">
            Tutup
          </button>
        </div>
      </Modal>
    </div>
  );
}
