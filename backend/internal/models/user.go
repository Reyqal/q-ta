package models

import (
	"time"
)

// User merepresentasikan pengguna sistem (admin atau penghuni kos).
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	PhoneNumber  string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone_number"`
	Role         string    `gorm:"type:varchar(20);not null;default:'penghuni'" json:"role"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Tenant *Tenant `gorm:"foreignKey:UserID" json:"tenant,omitempty"`
}

// TableName mengatur nama tabel di database.
func (User) TableName() string {
	return "users"
}
