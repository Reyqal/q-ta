package room

import (
	"errors"

	"github.com/qta/backend/internal/models"
	"gorm.io/gorm"
)

// Service menangani business logic untuk Room.
type Service struct {
	repo *Repository
}

// NewService membuat instance baru RoomService.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetAll mengambil semua kamar.
func (s *Service) GetAll() ([]models.Room, error) {
	return s.repo.FindAll()
}

// GetByID mengambil detail kamar berdasarkan ID.
func (s *Service) GetByID(id uint) (*models.Room, error) {
	room, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kamar tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data kamar")
	}
	return room, nil
}

// Create membuat kamar baru.
func (s *Service) Create(room *models.Room) error {
	if room.RoomNumber == "" {
		return errors.New("nomor kamar wajib diisi")
	}
	if room.RentAmount <= 0 {
		return errors.New("harga sewa harus lebih dari 0")
	}
	room.Status = "available"
	return s.repo.Create(room)
}

// Update memperbarui data kamar.
func (s *Service) Update(id uint, input *models.Room) (*models.Room, error) {
	room, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kamar tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data kamar")
	}

	if input.RoomNumber != "" {
		room.RoomNumber = input.RoomNumber
	}
	if input.RentAmount > 0 {
		room.RentAmount = input.RentAmount
	}
	if input.Description != "" {
		room.Description = input.Description
	}

	if err := s.repo.Update(room); err != nil {
		return nil, errors.New("gagal memperbarui data kamar")
	}
	return room, nil
}

// Delete menghapus kamar.
func (s *Service) Delete(id uint) error {
	room, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("kamar tidak ditemukan")
		}
		return errors.New("gagal mengambil data kamar")
	}

	if room.Status == "occupied" {
		return errors.New("tidak bisa menghapus kamar yang sedang ditempati")
	}

	return s.repo.Delete(id)
}
