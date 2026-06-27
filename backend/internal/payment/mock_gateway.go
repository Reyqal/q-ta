package payment

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"
)

// MockGateway adalah implementasi mock dari PaymentGatewayService.
// Digunakan untuk demo/development tanpa memerlukan API key payment gateway asli.
// Seluruh alur pembayaran bisa dijalankan dengan mock ini via endpoint /simulate.
type MockGateway struct{}

// NewMockGateway membuat instance baru MockGateway.
func NewMockGateway() *MockGateway {
	return &MockGateway{}
}

// CreateQRIS membuat respons QRIS dummy untuk simulasi.
// QR string berisi data mock yang bisa ditampilkan sebagai QR code di frontend.
func (m *MockGateway) CreateQRIS(orderID string, amount int64, description string) (*QRISResponse, error) {
	log.Printf("[MOCK PAYMENT] Membuat QRIS - OrderID: %s, Jumlah: Rp %d, Deskripsi: %s", orderID, amount, description)

	// Generate dummy QR string (format EMVCo-like)
	qrData := fmt.Sprintf("QTA-MOCK-%s-%d", orderID, amount)
	qrString := fmt.Sprintf(
		"00020101021226620014COM.QTA.MOCK01%02d%s0303QTA5204341453033605802ID5904QTA6007JAKARTA610540000540%d6304",
		len(orderID), orderID, amount,
	)

	// Generate dummy QR image URL (menggunakan public QR generator API untuk demo)
	qrImageURL := fmt.Sprintf(
		"https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=%s",
		base64.URLEncoding.EncodeToString([]byte(qrData)),
	)

	expiresAt := time.Now().Add(15 * time.Minute)

	return &QRISResponse{
		OrderID:    orderID,
		QRString:   qrString,
		QRImageURL: qrImageURL,
		ExpiresAt:  expiresAt,
	}, nil
}

// VerifyWebhookSignature selalu mengembalikan true untuk mock gateway.
// Di produksi, gunakan MidtransGateway yang memvalidasi SHA512 signature.
func (m *MockGateway) VerifyWebhookSignature(payload WebhookPayload) bool {
	log.Printf("[MOCK PAYMENT] Verifikasi webhook (selalu valid di mock) - OrderID: %s, Status: %s",
		payload.OrderID, payload.TransactionStatus)
	return true
}
