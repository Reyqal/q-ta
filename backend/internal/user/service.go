package user

import (
	"errors"

	"github.com/qta/backend/internal/models"
	"gorm.io/gorm"
)

// Service menangani business logic untuk User.
type Service struct {
	repo *Repository
}

// NewService membuat instance baru UserService.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetProfile mengambil profil user berdasarkan ID.
func (s *Service) GetProfile(userID uint) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengguna tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data pengguna")
	}
	return user, nil
}
