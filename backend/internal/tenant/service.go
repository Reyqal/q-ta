package tenant

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/qta/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateTenantResponse adalah response setelah berhasil membuat tenant.
type CreateTenantResponse struct {
	Tenant   models.Tenant `json:"tenant"`
	Password string        `json:"generated_password"`
}

// Service menangani business logic untuk Tenant.
type Service struct {
	repo *Repository
}

// NewService membuat instance baru TenantService.
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetAll mengambil semua tenant aktif.
func (s *Service) GetAll() ([]models.Tenant, error) {
	return s.repo.FindAll()
}

// Create membuat tenant baru dalam satu transaksi database.
// Langkah: create user -> create tenant -> update room status -> log notifikasi.
func (s *Service) Create(name, phoneNumber string, roomID uint) (*CreateTenantResponse, error) {
	if name == "" {
		return nil, errors.New("nama penghuni wajib diisi")
	}
	if phoneNumber == "" {
		return nil, errors.New("nomor telepon wajib diisi")
	}
	if roomID == 0 {
		return nil, errors.New("kamar wajib dipilih")
	}

	// Generate random password
	password := generateRandomPassword(8)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal membuat password")
	}

	var tenant models.Tenant
	db := s.repo.GetDB()

	// Execute in transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		// Check if room exists and is available
		var room models.Room
		if err := tx.First(&room, roomID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("kamar tidak ditemukan")
			}
			return errors.New("gagal mengecek data kamar")
		}
		if room.Status == "occupied" {
			return errors.New("kamar sudah ditempati")
		}

		// Check if phone number already exists
		var existingUser models.User
		if err := tx.Where("phone_number = ?", phoneNumber).First(&existingUser).Error; err == nil {
			return errors.New("nomor telepon sudah terdaftar")
		}

		// Create user
		user := models.User{
			Name:         name,
			PhoneNumber:  phoneNumber,
			Role:         "penghuni",
			PasswordHash: string(hashedPassword),
		}
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("gagal membuat user: %w", err)
		}

		// Create tenant
		tenant = models.Tenant{
			UserID:   user.ID,
			RoomID:   roomID,
			JoinDate: time.Now(),
			IsActive: true,
		}
		if err := tx.Create(&tenant).Error; err != nil {
			return fmt.Errorf("gagal membuat tenant: %w", err)
		}

		// Update room status to occupied
		if err := tx.Model(&models.Room{}).Where("id = ?", roomID).Update("status", "occupied").Error; err != nil {
			return fmt.Errorf("gagal mengubah status kamar: %w", err)
		}

		// Log notification (simulate WhatsApp send)
		message := fmt.Sprintf(
			"Selamat datang di kos, %s!\nNomor Kamar: %s\nLogin:\nTelp: %s\nPassword: %s",
			name, room.RoomNumber, phoneNumber, password,
		)
		notifLog := models.NotificationLog{
			TenantID: &tenant.ID,
			Channel:  "whatsapp",
			Message:  message,
			Status:   "simulated_sent",
		}
		if err := tx.Create(&notifLog).Error; err != nil {
			return fmt.Errorf("gagal mencatat notifikasi: %w", err)
		}

		// Reload tenant with relations
		tenant.User = user
		tenant.Room = room

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &CreateTenantResponse{
		Tenant:   tenant,
		Password: password,
	}, nil
}

// Deactivate menonaktifkan tenant dan mengembalikan status kamar ke 'available'.
func (s *Service) Deactivate(id uint) error {
	db := s.repo.GetDB()

	return db.Transaction(func(tx *gorm.DB) error {
		var tenant models.Tenant
		if err := tx.First(&tenant, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("penghuni tidak ditemukan")
			}
			return errors.New("gagal mengambil data penghuni")
		}

		if !tenant.IsActive {
			return errors.New("penghuni sudah tidak aktif")
		}

		// Deactivate tenant
		tenant.IsActive = false
		if err := tx.Save(&tenant).Error; err != nil {
			return errors.New("gagal menonaktifkan penghuni")
		}

		// Set room back to available
		if err := tx.Model(&models.Room{}).Where("id = ?", tenant.RoomID).Update("status", "available").Error; err != nil {
			return errors.New("gagal mengubah status kamar")
		}

		return nil
	})
}

// generateRandomPassword membuat password acak sepanjang n karakter.
func generateRandomPassword(n int) string {
	const charset = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	result := make([]byte, n)
	for i := range result {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[idx.Int64()]
	}
	return string(result)
}
