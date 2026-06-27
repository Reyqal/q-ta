package payment

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/qta/backend/internal/models"
	"gorm.io/gorm"
)

// Service mengelola logika bisnis untuk pembayaran.
type Service struct {
	db      *gorm.DB
	gateway PaymentGatewayService
}

// NewService membuat instance baru PaymentService.
func NewService(db *gorm.DB, gateway PaymentGatewayService) *Service {
	return &Service{db: db, gateway: gateway}
}

// CreateQRISPayment membuat pembayaran QRIS untuk sebuah invoice.
func (s *Service) CreateQRISPayment(invoiceID uint, userID uint) (*QRISResponse, error) {
	// Ambil invoice dengan validasi
	var invoice models.Invoice
	err := s.db.Preload("Tenant").First(&invoice, invoiceID).Error
	if err != nil {
		return nil, errors.New("tagihan tidak ditemukan")
	}

	// Validasi bahwa penghuni hanya bisa membayar tagihannya sendiri
	if invoice.Tenant.UserID != userID {
		return nil, errors.New("akses ditolak: tagihan ini bukan milik Anda")
	}

	if invoice.Status == "paid" {
		return nil, errors.New("tagihan sudah lunas")
	}

	// Generate unique order ID
	orderID := fmt.Sprintf("QTA-%d-%s", invoiceID, uuid.New().String()[:8])

	// Buat transaksi di database
	transaction := &models.Transaction{
		InvoiceID:      invoiceID,
		PaymentMethod:  "qris",
		GatewayOrderID: orderID,
		Status:         "pending",
		Amount:         invoice.Amount,
	}
	if err := s.db.Create(transaction).Error; err != nil {
		return nil, errors.New("gagal membuat transaksi")
	}

	// Buat QRIS via payment gateway
	description := fmt.Sprintf("Pembayaran Kos Periode %s", invoice.Period)
	qrisResp, err := s.gateway.CreateQRIS(orderID, invoice.Amount, description)
	if err != nil {
		// Rollback transaction record
		s.db.Delete(transaction)
		return nil, fmt.Errorf("gagal membuat QRIS: %w", err)
	}

	return qrisResp, nil
}

// ProcessWebhook memproses webhook callback dari payment gateway.
func (s *Service) ProcessWebhook(payload WebhookPayload) error {
	// Verifikasi signature
	if !s.gateway.VerifyWebhookSignature(payload) {
		return errors.New("signature webhook tidak valid")
	}

	// Cari transaksi berdasarkan order ID
	var transaction models.Transaction
	err := s.db.Where("gateway_order_id = ?", payload.OrderID).First(&transaction).Error
	if err != nil {
		return fmt.Errorf("transaksi tidak ditemukan: %s", payload.OrderID)
	}

	// Update status transaksi
	transaction.Status = payload.TransactionStatus
	transaction.GatewayReference = payload.TransactionID
	s.db.Save(&transaction)

	log.Printf("[WEBHOOK] OrderID: %s, Status: %s", payload.OrderID, payload.TransactionStatus)

	// Jika pembayaran berhasil (settlement), update invoice
	if payload.TransactionStatus == "settlement" {
		return s.settleInvoice(transaction.InvoiceID)
	}

	return nil
}

// SimulatePayment mensimulasikan pembayaran sukses (untuk demo/development).
func (s *Service) SimulatePayment(orderID string) error {
	var transaction models.Transaction
	err := s.db.Where("gateway_order_id = ?", orderID).First(&transaction).Error
	if err != nil {
		return fmt.Errorf("transaksi tidak ditemukan: %s", orderID)
	}

	if transaction.Status == "settlement" {
		return errors.New("transaksi sudah diselesaikan")
	}

	// Update status transaksi
	transaction.Status = "settlement"
	transaction.GatewayReference = fmt.Sprintf("SIM-%s", uuid.New().String()[:8])
	s.db.Save(&transaction)

	log.Printf("[SIMULASI] Pembayaran berhasil untuk OrderID: %s", orderID)

	return s.settleInvoice(transaction.InvoiceID)
}

// settleInvoice menandai invoice sebagai lunas dan menghitung split pajak.
func (s *Service) settleInvoice(invoiceID uint) error {
	var invoice models.Invoice
	if err := s.db.First(&invoice, invoiceID).Error; err != nil {
		return errors.New("tagihan tidak ditemukan")
	}

	if invoice.Status == "paid" {
		return nil // Sudah lunas, idempotent
	}

	// Hitung pemisahan pajak (default 0.5%)
	taxPercentage := 0.5
	taxPortion := int64(float64(invoice.Amount) * taxPercentage / 100)
	netPortion := invoice.Amount - taxPortion

	now := time.Now()
	invoice.Status = "paid"
	invoice.PaidAt = &now
	invoice.TaxPortion = taxPortion
	invoice.NetPortion = netPortion

	if err := s.db.Save(&invoice).Error; err != nil {
		return errors.New("gagal memperbarui status tagihan")
	}

	log.Printf("[SETTLEMENT] Invoice #%d lunas - Total: %d, Pajak: %d, Bersih: %d",
		invoiceID, invoice.Amount, taxPortion, netPortion)

	return nil
}

// GetPaymentStatus mengecek status pembayaran untuk sebuah invoice.
func (s *Service) GetPaymentStatus(invoiceID uint) (string, *models.Transaction, error) {
	var invoice models.Invoice
	if err := s.db.First(&invoice, invoiceID).Error; err != nil {
		return "", nil, errors.New("tagihan tidak ditemukan")
	}

	// Cari transaksi terbaru untuk invoice ini
	var transaction models.Transaction
	err := s.db.Where("invoice_id = ?", invoiceID).
		Order("created_at DESC").First(&transaction).Error

	if err != nil {
		return invoice.Status, nil, nil
	}

	return invoice.Status, &transaction, nil
}
