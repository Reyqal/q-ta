package dashboard

import (
	"github.com/gofiber/fiber/v2"
)

// Handler menangani HTTP request untuk dashboard admin.
type Handler struct {
	service *Service
}

// NewHandler membuat instance baru DashboardHandler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetSummary menangani GET /api/dashboard/summary — ringkasan dashboard admin.
func (h *Handler) GetSummary(c *fiber.Ctx) error {
	summary, err := h.service.GetSummary()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data dashboard",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil data dashboard",
		"data":    summary,
	})
}
