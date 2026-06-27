package models

import (
	"time"
)

// Invoice merepresentasikan tagihan bulanan untuk penghuni kos.
type Invoice struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	TenantID   uint       `gorm:"not null" json:"tenant_id"`
	Period     string     `gorm:"type:varchar(7);not null" json:"period"`
	Amount     int64      `gorm:"not null" json:"amount"`
	TaxPortion int64      `gorm:"not null;default:0" json:"tax_portion"`
	NetPortion int64      `gorm:"not null;default:0" json:"net_portion"`
	Status     string     `gorm:"type:varchar(20);not null;default:'unpaid'" json:"status"`
	DueDate    time.Time  `gorm:"type:date;not null" json:"due_date"`
	PaidAt     *time.Time `json:"paid_at"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Tenant       Tenant        `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:InvoiceID" json:"transactions,omitempty"`
}

// TableName mengatur nama tabel di database.
func (Invoice) TableName() string {
	return "invoices"
}
