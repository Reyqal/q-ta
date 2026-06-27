package payment

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/qta/backend/internal/middleware"
)

// Handler menangani HTTP request untuk pembayaran.
type Handler struct {
	service *Service
}

// NewHandler membuat instance baru PaymentHandler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateQRISRequest adalah request body untuk membuat pembayaran QRIS.
type CreateQRISRequest struct {
	InvoiceID uint `json:"invoice_id"`
}

// SimulateRequest adalah request body untuk simulasi pembayaran.
type SimulateRequest struct {
	OrderID string `json:"order_id"`
}

// CreateQRIS menangani POST /api/payments/create-qris — membuat pembayaran QRIS.
func (h *Handler) CreateQRIS(c *fiber.Ctx) error {
	var req CreateQRISRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tidak valid",
			"data":    nil,
		})
	}

	userID := middleware.ExtractUserID(c)

	qrisResp, err := h.service.CreateQRISPayment(req.InvoiceID, userID)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "akses ditolak: tagihan ini bukan milik Anda" {
			status = fiber.StatusForbidden
		} else if err.Error() == "tagihan sudah lunas" || err.Error() == "tagihan tidak ditemukan" {
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
		"message": "QRIS berhasil dibuat",
		"data":    qrisResp,
	})
}

// Webhook menangani POST /api/webhooks/midtrans — menerima callback dari Midtrans.
func (h *Handler) Webhook(c *fiber.Ctx) error {
	var payload WebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format webhook tidak valid",
		})
	}

	if err := h.service.ProcessWebhook(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Midtrans mengharapkan response 200 OK
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Webhook berhasil diproses",
	})
}

// SimulateWebhook menangani POST /api/webhooks/midtrans/simulate — simulasi pembayaran sukses.
func (h *Handler) SimulateWebhook(c *fiber.Ctx) error {
	var req SimulateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tidak valid",
			"data":    nil,
		})
	}

	if err := h.service.SimulatePayment(req.OrderID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Simulasi pembayaran berhasil",
		"data":    nil,
	})
}

// GetStatus menangani GET /api/payments/:invoice_id/status — cek status pembayaran.
func (h *Handler) GetStatus(c *fiber.Ctx) error {
	invoiceID, err := strconv.ParseUint(c.Params("invoice_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tagihan tidak valid",
			"data":    nil,
		})
	}

	invoiceStatus, transaction, err := h.service.GetPaymentStatus(uint(invoiceID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Status pembayaran",
		"data": fiber.Map{
			"invoice_status": invoiceStatus,
			"transaction":    transaction,
		},
	})
}
