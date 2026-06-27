import { useState, useEffect } from 'react';
import { QrCode, BookOpen, FileText, MessageSquare, DoorOpen, Wifi, Building2 } from 'lucide-react';
import { Navbar } from '../components/Navbar';
import { StatusBadge } from '../components/StatusBadge';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { formatRupiah } from '../lib/formatCurrency';
import apiClient from '../lib/apiClient';

interface Room {
  id: number;
  room_number: string;
  rent_amount: number;
  status: 'available' | 'occupied';
  description?: string;
}

const features = [
  {
    icon: QrCode,
    title: 'Pembayaran QRIS',
    description: 'Terima pembayaran sewa melalui QRIS yang langsung tercatat otomatis di sistem.',
    color: 'from-emerald-500 to-teal-500',
  },
  {
    icon: BookOpen,
    title: 'Pencatatan Otomatis',
    description: 'Seluruh transaksi tercatat rapi secara digital. Tidak perlu lagi buku catatan manual.',
    color: 'from-blue-500 to-cyan-500',
  },
  {
    icon: FileText,
    title: 'Laporan Pajak',
    description: 'Pemisahan pendapatan bersih dan titipan pajak otomatis untuk pelaporan yang mudah.',
    color: 'from-purple-500 to-pink-500',
  },
  {
    icon: MessageSquare,
    title: 'Notifikasi WhatsApp',
    description: 'Kirim pengingat tagihan dan konfirmasi pembayaran langsung ke WhatsApp penghuni.',
    color: 'from-amber-500 to-orange-500',
  },
];

export function LandingPage() {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const response = await apiClient.get('/rooms');
        if (response.data.success) {
          setRooms(response.data.data);
        }
      } catch {
        // Rooms will be empty if backend isn't available
      } finally {
        setIsLoading(false);
      }
    };
    fetchRooms();
  }, []);

  return (
    <div className="min-h-screen bg-surface">
      <Navbar />

      {/* Hero Section */}
      <section className="relative pt-32 pb-20 px-4 overflow-hidden">
        {/* Background effects */}
        <div className="absolute inset-0 bg-grid opacity-40" />
        <div className="absolute top-20 left-1/4 w-96 h-96 bg-emerald-500/10 rounded-full blur-3xl" />
        <div className="absolute bottom-10 right-1/4 w-80 h-80 bg-blue-500/10 rounded-full blur-3xl" />

        <div className="relative max-w-5xl mx-auto text-center">
          <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full glass text-sm text-slate-300 mb-8 animate-slideDown">
            <Wifi className="w-3.5 h-3.5 text-primary" />
            Solusi Digital untuk Kos-Kosan Modern
          </div>

          <h1 className="text-4xl sm:text-5xl lg:text-7xl font-black text-white leading-tight mb-6 animate-fadeIn">
            <span className="text-gradient">Q-TA</span>
            <span className="block text-3xl sm:text-4xl lg:text-5xl font-bold mt-2 text-slate-200">
              Sistem Manajemen Kos Modern
            </span>
          </h1>

          <p className="text-lg sm:text-xl text-slate-400 max-w-2xl mx-auto mb-10 animate-fadeIn leading-relaxed" style={{ animationDelay: '0.1s' }}>
            Digitalisasi pencatatan, penagihan, dan pelaporan keuangan kos-kosan Anda.
            Dirancang khusus untuk UMKM Indonesia.
          </p>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4 animate-fadeIn" style={{ animationDelay: '0.2s' }}>
            <a
              href="/login"
              className="px-8 py-3.5 rounded-2xl bg-gradient-to-r from-emerald-500 to-emerald-600 text-white font-semibold text-base shadow-lg shadow-emerald-500/25 hover:shadow-emerald-500/40 hover:from-emerald-400 hover:to-emerald-500 transition-all"
            >
              Mulai Sekarang
            </a>
            <a
              href="#rooms"
              className="px-8 py-3.5 rounded-2xl glass hover:bg-white/10 text-slate-300 hover:text-white font-medium text-base transition-all"
            >
              Lihat Ketersediaan
            </a>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 px-4">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-14">
            <h2 className="text-3xl sm:text-4xl font-bold text-white mb-3 animate-fadeIn">
              Fitur Unggulan
            </h2>
            <p className="text-slate-400 max-w-lg mx-auto">
              Semua yang Anda butuhkan untuk mengelola kos-kosan dengan efisien
            </p>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {features.map((feature, index) => (
              <div
                key={feature.title}
                className="glass rounded-2xl p-6 hover:bg-white/[0.08] hover:border-white/15 transition-all duration-300 group animate-slideUp"
                style={{ animationDelay: `${index * 0.1}s` }}
              >
                <div className={`w-12 h-12 rounded-xl bg-gradient-to-br ${feature.color} flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300 shadow-lg`}>
                  <feature.icon className="w-6 h-6 text-white" />
                </div>
                <h3 className="text-lg font-semibold text-white mb-2">{feature.title}</h3>
                <p className="text-sm text-slate-400 leading-relaxed">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Room Availability */}
      <section id="rooms" className="py-20 px-4 bg-surface-light/30">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-14">
            <h2 className="text-3xl sm:text-4xl font-bold text-white mb-3 animate-fadeIn">
              Ketersediaan Kamar
            </h2>
            <p className="text-slate-400 max-w-lg mx-auto">
              Cek status dan harga kamar yang tersedia saat ini
            </p>
          </div>

          {isLoading ? (
            <div className="py-16">
              <LoadingSpinner size="lg" text="Memuat data kamar..." />
            </div>
          ) : rooms.length === 0 ? (
            <div className="text-center py-16">
              <DoorOpen className="w-12 h-12 text-slate-600 mx-auto mb-4" />
              <p className="text-slate-400">Belum ada data kamar tersedia</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
              {rooms.map((room, index) => (
                <div
                  key={room.id}
                  className="glass rounded-2xl overflow-hidden hover:bg-white/[0.08] hover:border-white/15 transition-all duration-300 group animate-slideUp"
                  style={{ animationDelay: `${index * 0.05}s` }}
                >
                  <div className="p-5">
                    <div className="flex items-start justify-between mb-3">
                      <div className="flex items-center gap-2">
                        <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center">
                          <DoorOpen className="w-4 h-4 text-primary" />
                        </div>
                        <span className="text-lg font-bold text-white">
                          Kamar {room.room_number}
                        </span>
                      </div>
                      <StatusBadge status={room.status} />
                    </div>

                    <div className="mt-4">
                      <p className="text-xs text-slate-500 uppercase tracking-wider mb-1">
                        Harga Sewa / Bulan
                      </p>
                      <p className="text-xl font-bold text-primary">
                        {formatRupiah(room.rent_amount)}
                      </p>
                    </div>

                    {room.description && (
                      <p className="text-xs text-slate-500 mt-3 line-clamp-2">
                        {room.description}
                      </p>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </section>

      {/* Footer */}
      <footer className="py-8 px-4 border-t border-white/5">
        <div className="max-w-6xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-slate-500">
          <div className="flex items-center gap-2">
            <div className="w-6 h-6 rounded-lg bg-gradient-to-br from-emerald-500 to-blue-500 flex items-center justify-center">
              <Building2 className="w-3.5 h-3.5 text-white" />
            </div>
            <span>Q-TA &copy; {new Date().getFullYear()}</span>
          </div>
          <p>Solusi manajemen kos-kosan untuk UMKM Indonesia</p>
        </div>
      </footer>
    </div>
  );
}
