import { useState, useEffect } from 'react';
import { Plus, Pencil, Trash2, DoorOpen } from 'lucide-react';
import apiClient from '../../lib/apiClient';
import { formatRupiah } from '../../lib/formatCurrency';
import { Modal } from '../../components/Modal';
import { LoadingSpinner } from '../../components/LoadingSpinner';
import { EmptyState } from '../../components/EmptyState';

interface Room {
  id: number;
  room_number: string;
  status: 'available' | 'occupied';
  rent_amount: number;
  description: string;
}

export function RoomsPage() {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editRoom, setEditRoom] = useState<Room | null>(null);
  const [form, setForm] = useState({ room_number: '', rent_amount: '', description: '' });
  const [submitting, setSubmitting] = useState(false);

  const fetchRooms = async () => {
    try {
      const res = await apiClient.get('/rooms');
      if (res.data.success) setRooms(res.data.data || []);
    } catch (e) { console.error(e); }
    finally { setLoading(false); }
  };

  useEffect(() => { fetchRooms(); }, []);

  const openAdd = () => {
    setEditRoom(null);
    setForm({ room_number: '', rent_amount: '', description: '' });
    setShowModal(true);
  };

  const openEdit = (room: Room) => {
    setEditRoom(room);
    setForm({ room_number: room.room_number, rent_amount: String(room.rent_amount), description: room.description || '' });
    setShowModal(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      const payload = { room_number: form.room_number, rent_amount: Number(form.rent_amount), description: form.description };
      if (editRoom) {
        await apiClient.put(`/rooms/${editRoom.id}`, payload);
      } else {
        await apiClient.post('/rooms', payload);
      }
      setShowModal(false);
      fetchRooms();
    } catch (e: any) {
      alert(e.response?.data?.message || 'Gagal menyimpan');
    } finally { setSubmitting(false); }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Yakin ingin menghapus kamar ini?')) return;
    try {
      await apiClient.delete(`/rooms/${id}`);
      fetchRooms();
    } catch (e: any) {
      alert(e.response?.data?.message || 'Gagal menghapus');
    }
  };

  if (loading) return <LoadingSpinner />;

  return (
    <div className="animate-fade-in">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-white">Kelola Kamar</h1>
          <p className="text-slate-400 mt-1">Tambah, edit, dan hapus data kamar kos</p>
        </div>
        <button onClick={openAdd} className="flex items-center gap-2 px-4 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg transition-all duration-200 font-medium shadow-lg shadow-emerald-900/30">
          <Plus size={18} /> Tambah Kamar
        </button>
      </div>

      {rooms.length === 0 ? (
        <EmptyState icon={DoorOpen} message="Belum ada data kamar" />
      ) : (
        <div className="glass rounded-xl overflow-hidden">
          <table className="w-full">
            <thead>
              <tr className="border-b border-white/10">
                <th className="text-left px-6 py-4 text-sm font-semibold text-slate-300">No. Kamar</th>
                <th className="text-left px-6 py-4 text-sm font-semibold text-slate-300">Harga Sewa</th>
                <th className="text-left px-6 py-4 text-sm font-semibold text-slate-300">Status</th>
                <th className="text-left px-6 py-4 text-sm font-semibold text-slate-300">Deskripsi</th>
                <th className="text-right px-6 py-4 text-sm font-semibold text-slate-300">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {rooms.map((room) => (
                <tr key={room.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                  <td className="px-6 py-4 font-semibold text-white">{room.room_number}</td>
                  <td className="px-6 py-4 text-slate-300">{formatRupiah(room.rent_amount)}</td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex px-3 py-1 rounded-full text-xs font-semibold ${room.status === 'available' ? 'bg-emerald-500/20 text-emerald-400' : 'bg-rose-500/20 text-rose-400'}`}>
                      {room.status === 'available' ? 'Tersedia' : 'Terisi'}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-slate-400 text-sm">{room.description || '-'}</td>
                  <td className="px-6 py-4 text-right">
                    <div className="flex items-center justify-end gap-2">
                      <button onClick={() => openEdit(room)} className="p-2 hover:bg-blue-500/20 text-blue-400 rounded-lg transition-colors"><Pencil size={16} /></button>
                      <button onClick={() => handleDelete(room.id)} className="p-2 hover:bg-rose-500/20 text-rose-400 rounded-lg transition-colors" disabled={room.status === 'occupied'}><Trash2 size={16} /></button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <Modal isOpen={showModal} onClose={() => setShowModal(false)} title={editRoom ? 'Edit Kamar' : 'Tambah Kamar'}>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">Nomor Kamar</label>
            <input type="text" value={form.room_number} onChange={e => setForm({...form, room_number: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-lg text-white focus:ring-2 focus:ring-emerald-500 focus:border-transparent outline-none" placeholder="Contoh: 101" required />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">Harga Sewa (Rp)</label>
            <input type="number" value={form.rent_amount} onChange={e => setForm({...form, rent_amount: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-lg text-white focus:ring-2 focus:ring-emerald-500 focus:border-transparent outline-none" placeholder="1500000" required />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">Deskripsi</label>
            <textarea value={form.description} onChange={e => setForm({...form, description: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-lg text-white focus:ring-2 focus:ring-emerald-500 focus:border-transparent outline-none" rows={3} placeholder="Deskripsi kamar (opsional)" />
          </div>
          <div className="flex justify-end gap-3 pt-2">
            <button type="button" onClick={() => setShowModal(false)} className="px-4 py-2.5 text-slate-400 hover:text-white transition-colors">Batal</button>
            <button type="submit" disabled={submitting} className="px-6 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg font-medium transition-all disabled:opacity-50">
              {submitting ? 'Menyimpan...' : 'Simpan'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
