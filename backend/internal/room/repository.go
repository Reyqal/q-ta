package room

import (
	"github.com/qta/backend/internal/models"
	"gorm.io/gorm"
)

// Repository menangani operasi database untuk Room.
type Repository struct {
	db *gorm.DB
}

// NewRepository membuat instance baru RoomRepository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindAll mengambil semua kamar.
func (r *Repository) FindAll() ([]models.Room, error) {
	var rooms []models.Room
	result := r.db.Order("room_number ASC").Find(&rooms)
	return rooms, result.Error
}

// FindByID mencari kamar berdasarkan ID.
func (r *Repository) FindByID(id uint) (*models.Room, error) {
	var room models.Room
	result := r.db.First(&room, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &room, nil
}

// Create menyimpan kamar baru ke database.
func (r *Repository) Create(room *models.Room) error {
	return r.db.Create(room).Error
}

// Update memperbarui data kamar di database.
func (r *Repository) Update(room *models.Room) error {
	return r.db.Save(room).Error
}

// Delete menghapus kamar dari database.
func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&models.Room{}, id).Error
}

// UpdateStatus memperbarui status kamar.
func (r *Repository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Room{}).Where("id = ?", id).Update("status", status).Error
}
