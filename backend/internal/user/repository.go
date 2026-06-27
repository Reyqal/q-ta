package user

import (
	"github.com/qta/backend/internal/models"
	"gorm.io/gorm"
)

// Repository menangani operasi database untuk User.
type Repository struct {
	db *gorm.DB
}

// NewRepository membuat instance baru UserRepository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindByID mencari user berdasarkan ID.
func (r *Repository) FindByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByPhoneNumber mencari user berdasarkan nomor telepon.
func (r *Repository) FindByPhoneNumber(phone string) (*models.User, error) {
	var user models.User
	result := r.db.Where("phone_number = ?", phone).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Create menyimpan user baru ke database.
func (r *Repository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update memperbarui data user di database.
func (r *Repository) Update(user *models.User) error {
	return r.db.Save(user).Error
}
