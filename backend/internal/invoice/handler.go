package invoice

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/qta/backend/internal/middleware"
)

// Handler menangani HTTP request untuk Invoice.
type Handler struct {
	service *Service
}

// NewHandler membuat instance baru InvoiceHandler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateInvoiceRequest adalah request body untuk membuat tagihan.
type CreateInvoiceRequest struct {
	TenantID uint   `json:"tenant_id"`
	Period   string `json:"period"`
	Amount   int64  `json:"amount"`
	DueDate  string `json:"due_date"` // Format: "2006-01-02"
}

// UpdateStatusRequest adalah request body untuk mengubah status tagihan.
type UpdateStatusRequest struct {
	Status string `json:"status"` // "issue" atau "unpaid"
}

// ConfirmManualRequest adalah request body untuk konfirmasi pembayaran manual.
type ConfirmManualRequest struct {
	PaymentMethod string `json:"payment_method"` // "cash" atau "transfer"
}

// GetAll menangani GET /api/invoices — admin: semua, penghuni: miliknya saja.
func (h *Handler) GetAll(c *fiber.Ctx) error {
	userID, role := middleware.ExtractClaims(c)

	var invoices interface{}
	var err error

	if role == "admin" {
		invoices, err = h.service.GetAll()
	} else {
		invoices, err = h.service.GetByTenantUserID(userID)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil data tagihan",
		"data":    invoices,
	})
}

// GetByID menangani GET /api/invoices/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tagihan tidak valid",
			"data":    nil,
		})
	}

	userID, role := middleware.ExtractClaims(c)
	invoice, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	// Penghuni hanya bisa melihat tagihan miliknya sendiri
	if role == "penghuni" && invoice.Tenant.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Akses ditolak",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil detail tagihan",
		"data":    invoice,
	})
}

// Create menangani POST /api/invoices — membuat tagihan baru (admin only).
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tidak valid",
			"data":    nil,
		})
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format tanggal jatuh tempo tidak valid (gunakan YYYY-MM-DD)",
			"data":    nil,
		})
	}

	invoice, err := h.service.Create(req.TenantID, req.Period, req.Amount, dueDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Tagihan berhasil dibuat",
		"data":    invoice,
	})
}

// GenerateMonthly menangani POST /api/invoices/generate-monthly.
func (h *Handler) GenerateMonthly(c *fiber.Ctx) error {
	count, err := h.service.GenerateMonthly()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Tagihan bulanan berhasil di-generate",
		"data": fiber.Map{
			"generated_count": count,
		},
	})
}

// UpdateStatus menangani PUT /api/invoices/:id/status.
func (h *Handler) UpdateStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tagihan tidak valid",
			"data":    nil,
		})
	}

	var req UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tidak valid",
			"data":    nil,
		})
	}

	invoice, err := h.service.UpdateStatus(uint(id), req.Status)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Status tagihan berhasil diperbarui",
		"data":    invoice,
	})
}

// ConfirmManual menangani PUT /api/invoices/:id/confirm-manual.
func (h *Handler) ConfirmManual(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tagihan tidak valid",
			"data":    nil,
		})
	}

	var req ConfirmManualRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tidak valid",
			"data":    nil,
		})
	}

	invoice, err := h.service.ConfirmManualPayment(uint(id), req.PaymentMethod)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Pembayaran manual berhasil dikonfirmasi",
		"data":    invoice,
	})
}
