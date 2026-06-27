package tenant

import (
	"github.com/qta/backend/internal/models"
	"gorm.io/gorm"
)

// Repository menangani operasi database untuk Tenant.
type Repository struct {
	db *gorm.DB
}

// NewRepository membuat instance baru TenantRepository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindAll mengambil semua tenant yang aktif beserta data user dan room.
func (r *Repository) FindAll() ([]models.Tenant, error) {
	var tenants []models.Tenant
	result := r.db.Preload("User").Preload("Room").
		Where("is_active = ?", true).
		Order("id DESC").
		Find(&tenants)
	return tenants, result.Error
}

// FindByID mencari tenant berdasarkan ID beserta data user dan room.
func (r *Repository) FindByID(id uint) (*models.Tenant, error) {
	var tenant models.Tenant
	result := r.db.Preload("User").Preload("Room").First(&tenant, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tenant, nil
}

// FindByUserID mencari tenant berdasarkan User ID.
func (r *Repository) FindByUserID(userID uint) (*models.Tenant, error) {
	var tenant models.Tenant
	result := r.db.Preload("User").Preload("Room").
		Where("user_id = ? AND is_active = ?", userID, true).
		First(&tenant)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tenant, nil
}

// Create menyimpan tenant baru ke database.
func (r *Repository) Create(tenant *models.Tenant) error {
	return r.db.Create(tenant).Error
}

// Update memperbarui data tenant di database.
func (r *Repository) Update(tenant *models.Tenant) error {
	return r.db.Save(tenant).Error
}

// FindAllActive mengambil semua tenant yang aktif beserta data room (untuk generate tagihan bulanan).
func (r *Repository) FindAllActive() ([]models.Tenant, error) {
	var tenants []models.Tenant
	result := r.db.Preload("User").Preload("Room").
		Where("is_active = ?", true).
		Find(&tenants)
	return tenants, result.Error
}

// GetDB mengembalikan instance database untuk transaction support.
func (r *Repository) GetDB() *gorm.DB {
	return r.db
}
