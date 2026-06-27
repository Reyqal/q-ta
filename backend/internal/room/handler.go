package room

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/qta/backend/internal/models"
)

// Handler menangani HTTP request untuk Room.
type Handler struct {
	service *Service
}

// NewHandler membuat instance baru RoomHandler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateRoomRequest adalah request body untuk membuat kamar.
type CreateRoomRequest struct {
	RoomNumber  string `json:"room_number"`
	RentAmount  int64  `json:"rent_amount"`
	Description string `json:"description"`
}

// GetAll menangani GET /api/rooms - mengambil semua kamar.
func (h *Handler) GetAll(c *fiber.Ctx) error {
	rooms, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data kamar",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil data kamar",
		"data":    rooms,
	})
}

// Create menangani POST /api/rooms - membuat kamar baru.
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tidak valid",
			"data":    nil,
		})
	}

	room := &models.Room{
		RoomNumber:  req.RoomNumber,
		RentAmount:  req.RentAmount,
		Description: req.Description,
	}

	if err := h.service.Create(room); err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "nomor kamar wajib diisi" || err.Error() == "harga sewa harus lebih dari 0" {
			status = fiber.StatusBadRequest
		}
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Kamar berhasil ditambahkan",
		"data":    room,
	})
}

// Update menangani PUT /api/rooms/:id - memperbarui kamar.
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID kamar tidak valid",
			"data":    nil,
		})
	}

	var req CreateRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tidak valid",
			"data":    nil,
		})
	}

	input := &models.Room{
		RoomNumber:  req.RoomNumber,
		RentAmount:  req.RentAmount,
		Description: req.Description,
	}

	room, err := h.service.Update(uint(id), input)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "kamar tidak ditemukan" {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Kamar berhasil diperbarui",
		"data":    room,
	})
}

// Delete menangani DELETE /api/rooms/:id - menghapus kamar.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID kamar tidak valid",
			"data":    nil,
		})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "kamar tidak ditemukan" {
			status = fiber.StatusNotFound
		} else if err.Error() == "tidak bisa menghapus kamar yang sedang ditempati" {
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
		"message": "Kamar berhasil dihapus",
		"data":    nil,
	})
}
