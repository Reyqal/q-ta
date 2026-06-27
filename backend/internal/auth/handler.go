package auth

import (
	"github.com/gofiber/fiber/v2"
)

// Handler menangani HTTP request untuk autentikasi.
type Handler struct {
	service *Service
}

// NewHandler membuat instance baru AuthHandler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// LoginRequest adalah request body untuk login.
type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

// Login menangani POST /api/auth/login.
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tidak valid",
			"data":    nil,
		})
	}

	// Validasi input
	if req.PhoneNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nomor telepon wajib diisi",
			"data":    nil,
		})
	}
	if req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Password wajib diisi",
			"data":    nil,
		})
	}

	result, err := h.service.Login(req.PhoneNumber, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data":    result,
	})
}
