package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/qta/backend/internal/middleware"
)

// Handler menangani HTTP request untuk User.
type Handler struct {
	service *Service
}

// NewHandler membuat instance baru UserHandler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetMe menangani GET /api/users/me - mengembalikan profil user yang sedang login.
func (h *Handler) GetMe(c *fiber.Ctx) error {
	userID := middleware.ExtractUserID(c)

	user, err := h.service.GetProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil data profil",
		"data":    user,
	})
}
