package payment

import (
	"time"
)

// QRISResponse merepresentasikan respons dari payment gateway setelah membuat transaksi QRIS.
type QRISResponse struct {
	OrderID    string    `json:"order_id"`
	QRString   string    `json:"qr_string"`
	QRImageURL string    `json:"qr_image_url"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// WebhookPayload merepresentasikan data yang diterima dari payment gateway via webhook.
type WebhookPayload struct {
	OrderID           string `json:"order_id"`
	TransactionID     string `json:"transaction_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	TransactionStatus string `json:"transaction_status"`
	PaymentType       string `json:"payment_type"`
}

// PaymentGatewayService adalah interface untuk abstraksi payment gateway.
// Implementasi default menggunakan mock gateway.
// Untuk produksi, gunakan MidtransGateway dengan mengisi MIDTRANS_SERVER_KEY di .env
type PaymentGatewayService interface {
	// CreateQRIS membuat transaksi QRIS baru dan mengembalikan data QR code.
	CreateQRIS(orderID string, amount int64, description string) (*QRISResponse, error)

	// VerifyWebhookSignature memvalidasi signature dari webhook callback.
	VerifyWebhookSignature(payload WebhookPayload) bool
}
