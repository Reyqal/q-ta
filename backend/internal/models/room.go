package models

import (
	"time"
)

// Room merepresentasikan kamar kos.
type Room struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RoomNumber  string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"room_number"`
	Status      string    `gorm:"type:varchar(20);not null;default:'available'" json:"status"`
	RentAmount  int64     `gorm:"not null" json:"rent_amount"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Tenant *Tenant `gorm:"foreignKey:RoomID" json:"tenant,omitempty"`
}

// TableName mengatur nama tabel di database.
func (Room) TableName() string {
	return "rooms"
}
