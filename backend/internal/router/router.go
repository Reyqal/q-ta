package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/qta/backend/internal/auth"
	"github.com/qta/backend/internal/dashboard"
	"github.com/qta/backend/internal/invoice"
	"github.com/qta/backend/internal/middleware"
	"github.com/qta/backend/internal/notification"
	"github.com/qta/backend/internal/payment"
	"github.com/qta/backend/internal/room"
	"github.com/qta/backend/internal/tenant"
	"github.com/qta/backend/internal/user"
)

// Handlers mengelompokkan semua handler yang diperlukan untuk routing.
type Handlers struct {
	Auth         *auth.Handler
	User         *user.Handler
	Room         *room.Handler
	Tenant       *tenant.Handler
	Invoice      *invoice.Handler
	Payment      *payment.Handler
	Dashboard    *dashboard.Handler
	Notification *notification.FiberHandler
}

// Setup mendaftarkan semua route ke Fiber app.
func Setup(app *fiber.App, h *Handlers, jwtSecret string) {
	api := app.Group("/api")

	// ==========================================
	// Public Routes (tanpa autentikasi)
	// ==========================================
	api.Get("/rooms", h.Room.GetAll)
	api.Post("/auth/login", h.Auth.Login)

	// Webhook endpoint (tidak perlu JWT, validasi via signature)
	api.Post("/webhooks/midtrans", h.Payment.Webhook)
	api.Post("/webhooks/midtrans/simulate", h.Payment.SimulateWebhook)

	// ==========================================
	// Authenticated Routes (perlu JWT)
	// ==========================================
	authenticated := api.Group("", middleware.JWTProtected(jwtSecret))

	// User profile
	authenticated.Get("/users/me", h.User.GetMe)

	// Invoices (role-aware: admin lihat semua, penghuni lihat miliknya)
	authenticated.Get("/invoices", h.Invoice.GetAll)
	authenticated.Get("/invoices/:id", h.Invoice.GetByID)

	// Payment (penghuni)
	authenticated.Post("/payments/create-qris", h.Payment.CreateQRIS)
	authenticated.Get("/payments/:invoice_id/status", h.Payment.GetStatus)

	// ==========================================
	// Admin Only Routes
	// ==========================================
	admin := authenticated.Group("", middleware.RequireRole("admin"))

	// Dashboard
	admin.Get("/dashboard/summary", h.Dashboard.GetSummary)

	// Room management
	admin.Post("/rooms", h.Room.Create)
	admin.Put("/rooms/:id", h.Room.Update)
	admin.Delete("/rooms/:id", h.Room.Delete)

	// Tenant management
	admin.Get("/tenants", h.Tenant.GetAll)
	admin.Post("/tenants", h.Tenant.Create)
	admin.Delete("/tenants/:id", h.Tenant.Delete)

	// Invoice management (admin)
	admin.Post("/invoices", h.Invoice.Create)
	admin.Post("/invoices/generate-monthly", h.Invoice.GenerateMonthly)
	admin.Put("/invoices/:id/status", h.Invoice.UpdateStatus)
	admin.Put("/invoices/:id/confirm-manual", h.Invoice.ConfirmManual)

	// Notifications
	admin.Get("/notifications", h.Notification.GetAll)
}
