package tenant

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler menangani HTTP request untuk Tenant.
type Handler struct {
	service *Service
}

// NewHandler membuat instance baru TenantHandler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateTenantRequest adalah request body untuk membuat tenant.
type CreateTenantRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	RoomID      uint   `json:"room_id"`
}

// GetAll menangani GET /api/tenants - mengambil semua tenant aktif.
func (h *Handler) GetAll(c *fiber.Ctx) error {
	tenants, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data penghuni",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil data penghuni",
		"data":    tenants,
	})
}

// Create menangani POST /api/tenants - membuat tenant baru.
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tidak valid",
			"data":    nil,
		})
	}

	result, err := h.service.Create(req.Name, req.PhoneNumber, req.RoomID)
	if err != nil {
		status := fiber.StatusInternalServerError
		switch err.Error() {
		case "nama penghuni wajib diisi",
			"nomor telepon wajib diisi",
			"kamar wajib dipilih":
			status = fiber.StatusBadRequest
		case "kamar tidak ditemukan":
			status = fiber.StatusNotFound
		case "kamar sudah ditempati",
			"nomor telepon sudah terdaftar":
			status = fiber.StatusConflict
		}
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Penghuni berhasil ditambahkan",
		"data":    result,
	})
}

// Delete menangani DELETE /api/tenants/:id - menonaktifkan tenant.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID penghuni tidak valid",
			"data":    nil,
		})
	}

	if err := h.service.Deactivate(uint(id)); err != nil {
		status := fiber.StatusInternalServerError
		switch err.Error() {
		case "penghuni tidak ditemukan":
			status = fiber.StatusNotFound
		case "penghuni sudah tidak aktif":
			status = fiber.StatusConflict
		}
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Penghuni berhasil dinonaktifkan",
		"data":    nil,
	})
}
