package cron

import (
	"fmt"
	"log"

	"github.com/qta/backend/internal/invoice"
	"github.com/qta/backend/internal/notification"
	robfigCron "github.com/robfig/cron/v3"
)

// Scheduler mengelola penjadwalan tugas otomatis.
type Scheduler struct {
	cron        *robfigCron.Cron
	invoiceRepo *invoice.Repository
	whatsApp    notification.WhatsAppGatewayService
}

// NewScheduler membuat instance baru Scheduler.
func NewScheduler(invoiceRepo *invoice.Repository, whatsApp notification.WhatsAppGatewayService) *Scheduler {
	return &Scheduler{
		cron:        robfigCron.New(),
		invoiceRepo: invoiceRepo,
		whatsApp:    whatsApp,
	}
}

// Start memulai semua cron job yang terdaftar.
func (s *Scheduler) Start() {
	// Jalankan pengecekan tagihan jatuh tempo setiap hari jam 08:00
	_, err := s.cron.AddFunc("0 8 * * *", s.checkOverdueInvoices)
	if err != nil {
		log.Printf("❌ Gagal mendaftarkan cron job: %v", err)
		return
	}

	s.cron.Start()
	log.Println("✅ Cron scheduler dimulai (pengecekan tagihan setiap hari jam 08:00)")
}

// Stop menghentikan semua cron job.
func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("⏹ Cron scheduler dihentikan")
}

// checkOverdueInvoices mengecek tagihan jatuh tempo dan mengirim reminder.
func (s *Scheduler) checkOverdueInvoices() {
	log.Println("[CRON] Memulai pengecekan tagihan jatuh tempo...")

	overdueInvoices, err := s.invoiceRepo.FindOverdue()
	if err != nil {
		log.Printf("[CRON] Gagal mengambil tagihan jatuh tempo: %v", err)
		return
	}

	log.Printf("[CRON] Ditemukan %d tagihan jatuh tempo", len(overdueInvoices))

	for _, inv := range overdueInvoices {
		// Skip tagihan dengan status "issue" — admin sudah menangani
		if inv.Status == "issue" {
			log.Printf("[CRON] Skip tagihan #%d (status: ada kendala)", inv.ID)
			continue
		}

		// Kirim reminder via WhatsApp (mock)
		message := fmt.Sprintf(
			"🏠 Pengingat Pembayaran Kos\n\n"+
				"Hai %s,\n"+
				"Tagihan kos Anda untuk periode %s sebesar Rp %d sudah jatuh tempo.\n\n"+
				"Silakan segera lakukan pembayaran melalui:\n"+
				"1️⃣ Bayar Sekarang (QRIS) — Login ke portal penghuni\n"+
				"2️⃣ Ada Kendala? — Hubungi admin kos\n\n"+
				"Terima kasih. 🙏",
			inv.Tenant.User.Name,
			inv.Period,
			inv.Amount,
		)

		err := s.whatsApp.SendMessage(inv.Tenant.User.PhoneNumber, message, inv.TenantID)
		if err != nil {
			log.Printf("[CRON] Gagal mengirim reminder untuk tagihan #%d: %v", inv.ID, err)
		} else {
			log.Printf("[CRON] Reminder terkirim untuk tagihan #%d (penghuni: %s)", inv.ID, inv.Tenant.User.Name)
		}
	}

	log.Println("[CRON] Pengecekan tagihan selesai")
}
