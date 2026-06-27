package models

import (
	"time"
)

// Tenant merepresentasikan penyewa kamar kos (penghuni).
type Tenant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	RoomID    uint      `gorm:"not null;uniqueIndex" json:"room_id"`
	JoinDate  time.Time `gorm:"type:date;not null;default:CURRENT_DATE" json:"join_date"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Room     Room      `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	Invoices []Invoice `gorm:"foreignKey:TenantID" json:"invoices,omitempty"`
}

// TableName mengatur nama tabel di database.
func (Tenant) TableName() string {
	return "tenants"
}
