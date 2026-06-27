package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/qta/backend/internal/config"
	"github.com/qta/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginResponse adalah response body untuk login yang berhasil.
type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Service menangani business logic untuk autentikasi.
type Service struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewService membuat instance baru AuthService.
func NewService(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{db: db, cfg: cfg}
}

// Login memvalidasi kredensial pengguna dan mengembalikan JWT token.
func (s *Service) Login(phoneNumber, password string) (*LoginResponse, error) {
	var user models.User
	result := s.db.Where("phone_number = ?", phoneNumber).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("nomor telepon atau password salah")
		}
		return nil, errors.New("terjadi kesalahan pada server")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("nomor telepon atau password salah")
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": float64(user.ID),
		"role":    user.Role,
		"exp":     time.Now().Add(time.Duration(s.cfg.JWTExpiryHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	return &LoginResponse{
		Token: tokenString,
		User:  user,
	}, nil
}
