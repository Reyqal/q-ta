package models

import (
	"time"
)

// NotificationLog merepresentasikan log pengiriman notifikasi (WhatsApp, dll).
type NotificationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TenantID  *uint     `json:"tenant_id"`
	Channel   string    `gorm:"type:varchar(50);not null;default:'whatsapp'" json:"channel"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Status    string    `gorm:"type:varchar(30);not null;default:'simulated_sent'" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

// TableName mengatur nama tabel di database.
func (NotificationLog) TableName() string {
	return "notification_log"
}
