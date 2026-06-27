package payment

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/qta/backend/internal/config"
)

// MidtransGateway adalah implementasi PaymentGatewayService menggunakan Midtrans Core API.
// Aktif ketika MIDTRANS_SERVER_KEY diisi di .env
//
// Untuk menggunakan:
// 1. Daftar di https://dashboard.sandbox.midtrans.com
// 2. Ambil Server Key dari Settings → Access Keys
// 3. Isi MIDTRANS_SERVER_KEY di .env
// 4. Set MIDTRANS_IS_PRODUCTION=false untuk sandbox
//
// Catatan: Implementasi ini menggunakan HTTP langsung ke Midtrans Core API
// tanpa library resmi midtrans-go untuk mengurangi dependency.
// Untuk integrasi penuh, pertimbangkan menggunakan github.com/midtrans/midtrans-go
type MidtransGateway struct {
	serverKey    string
	isProduction bool
}

// NewMidtransGateway membuat instance baru MidtransGateway.
func NewMidtransGateway(cfg *config.Config) *MidtransGateway {
	return &MidtransGateway{
		serverKey:    cfg.MidtransServerKey,
		isProduction: cfg.MidtransIsProduction,
	}
}

// getBaseURL mengembalikan base URL Midtrans berdasarkan environment.
func (m *MidtransGateway) getBaseURL() string {
	if m.isProduction {
		return "https://api.midtrans.com"
	}
	return "https://api.sandbox.midtrans.com"
}

// CreateQRIS membuat transaksi QRIS via Midtrans Core API.
//
// Endpoint: POST {baseURL}/v2/charge
// Body: {"payment_type":"qris","transaction_details":{"order_id":"...","gross_amount":...}}
//
// TODO: Implementasi HTTP call ke Midtrans Core API.
// Untuk saat ini, menggunakan mock response dengan format yang sesuai Midtrans.
// Saat siap integrasi real, uncomment kode HTTP call di bawah dan hapus mock response.
func (m *MidtransGateway) CreateQRIS(orderID string, amount int64, description string) (*QRISResponse, error) {
	log.Printf("[MIDTRANS] Membuat QRIS - OrderID: %s, Jumlah: %d", orderID, amount)

	// ====================================================================
	// TODO: Uncomment blok berikut untuk integrasi real Midtrans Core API
	// ====================================================================
	//
	// requestBody := map[string]interface{}{
	//     "payment_type": "qris",
	//     "transaction_details": map[string]interface{}{
	//         "order_id":     orderID,
	//         "gross_amount": amount,
	//     },
	//     "qris": map[string]interface{}{
	//         "acquirer": "gopay",
	//     },
	// }
	//
	// jsonBody, _ := json.Marshal(requestBody)
	// req, _ := http.NewRequest("POST", m.getBaseURL()+"/v2/charge", bytes.NewBuffer(jsonBody))
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(m.serverKey+":")))
	//
	// client := &http.Client{Timeout: 30 * time.Second}
	// resp, err := client.Do(req)
	// if err != nil {
	//     return nil, fmt.Errorf("gagal menghubungi Midtrans: %w", err)
	// }
	// defer resp.Body.Close()
	//
	// var midtransResp struct {
	//     StatusCode  string `json:"status_code"`
	//     QRString    string `json:"qr_string"`
	//     Actions     []struct {
	//         Name   string `json:"name"`
	//         URL    string `json:"url"`
	//     } `json:"actions"`
	// }
	//
	// json.NewDecoder(resp.Body).Decode(&midtransResp)
	//
	// qrImageURL := ""
	// for _, action := range midtransResp.Actions {
	//     if action.Name == "generate-qr-code" {
	//         qrImageURL = action.URL
	//     }
	// }
	//
	// return &QRISResponse{
	//     OrderID:    orderID,
	//     QRString:   midtransResp.QRString,
	//     QRImageURL: qrImageURL,
	//     ExpiresAt:  time.Now().Add(15 * time.Minute),
	// }, nil
	// ====================================================================

	// Sementara gunakan mock response yang sesuai format Midtrans
	qrString := fmt.Sprintf("00020101021226620014COM.MIDTRANS01%02d%s5204341453033605802ID6007JAKARTA540%d6304",
		len(orderID), orderID, amount)

	qrImageURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=%s", orderID)

	return &QRISResponse{
		OrderID:    orderID,
		QRString:   qrString,
		QRImageURL: qrImageURL,
		ExpiresAt:  time.Now().Add(15 * time.Minute),
	}, nil
}

// VerifyWebhookSignature memvalidasi signature webhook dari Midtrans.
// Formula: SHA512(order_id + status_code + gross_amount + server_key)
func (m *MidtransGateway) VerifyWebhookSignature(payload WebhookPayload) bool {
	input := payload.OrderID + payload.StatusCode + payload.GrossAmount + m.serverKey
	hash := sha512.Sum512([]byte(input))
	calculated := hex.EncodeToString(hash[:])

	isValid := calculated == payload.SignatureKey
	if !isValid {
		log.Printf("[MIDTRANS] Signature tidak valid untuk OrderID: %s", payload.OrderID)
	}
	return isValid
}
