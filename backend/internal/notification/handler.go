package notification

import (
	"github.com/gofiber/fiber/v2"
)

// FiberHandler menangani HTTP request untuk notifikasi menggunakan Fiber.
type FiberHandler struct {
	repo *Repository
}

// NewFiberHandler membuat instance baru notification FiberHandler.
func NewFiberHandler(repo *Repository) *FiberHandler {
	return &FiberHandler{repo: repo}
}

// GetAll menangani GET /api/notifications — mengambil semua log notifikasi (admin only).
func (h *FiberHandler) GetAll(c *fiber.Ctx) error {
	logs, err := h.repo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data notifikasi",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil data notifikasi",
		"data":    logs,
	})
}
