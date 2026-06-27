package notification

import (
	"fmt"
	"log"

	"github.com/qta/backend/internal/models"
	"gorm.io/gorm"
)

// WhatsAppGatewayService adalah interface untuk abstraksi pengiriman pesan WhatsApp.
//
// Untuk integrasi real WhatsApp, implementasikan interface ini dengan:
// - Baileys (WhatsApp Web API, gratis tapi tidak resmi)
// - WhatsApp Business API (resmi, berbayar via BSP seperti Wablas, Fonnte, dll.)
// - WhatsApp Cloud API (Meta, gratis 1000 pesan/bulan)
type WhatsAppGatewayService interface {
	SendMessage(phoneNumber string, message string, tenantID uint) error
}

// MockWhatsAppService adalah implementasi mock yang mencatat pesan ke database.
type MockWhatsAppService struct {
	db *gorm.DB
}

// NewMockWhatsAppService membuat instance baru MockWhatsAppService.
func NewMockWhatsAppService(db *gorm.DB) *MockWhatsAppService {
	return &MockWhatsAppService{db: db}
}

// SendMessage mencatat pesan ke tabel notification_log sebagai simulasi pengiriman.
// TODO: Ganti implementasi ini dengan WhatsApp Business API atau Baileys untuk produksi.
func (m *MockWhatsAppService) SendMessage(phoneNumber string, message string, tenantID uint) error {
	log.Printf("[MOCK WHATSAPP] Mengirim ke %s: %s", phoneNumber, message)

	notifLog := &models.NotificationLog{
		TenantID: tenantID,
		Channel:  "whatsapp",
		Message:  fmt.Sprintf("[Ke: %s] %s", phoneNumber, message),
		Status:   "simulated_sent",
	}

	if err := m.db.Create(notifLog).Error; err != nil {
		log.Printf("[MOCK WHATSAPP] Gagal mencatat notifikasi: %v", err)
		return err
	}

	return nil
}

// Repository mengelola akses data NotificationLog.
type Repository struct {
	db *gorm.DB
}

// NewRepository membuat instance baru NotificationRepository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindAll mengambil semua log notifikasi.
func (r *Repository) FindAll() ([]models.NotificationLog, error) {
	var logs []models.NotificationLog
	err := r.db.Order("created_at DESC").Limit(100).Find(&logs).Error
	return logs, err
}
