package models

import (
	"time"
)

// Transaction merepresentasikan transaksi pembayaran dari payment gateway.
type Transaction struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	InvoiceID        uint      `gorm:"not null" json:"invoice_id"`
	PaymentMethod    string    `gorm:"type:varchar(20);not null" json:"payment_method"`
	GatewayReference string    `gorm:"type:varchar(255)" json:"gateway_reference"`
	GatewayOrderID   string    `gorm:"type:varchar(255);uniqueIndex" json:"gateway_order_id"`
	Status           string    `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Amount           int64     `gorm:"not null" json:"amount"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Invoice Invoice `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
}

// TableName mengatur nama tabel di database.
func (Transaction) TableName() string {
	return "transactions"
}
