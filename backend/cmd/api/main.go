package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/qta/backend/internal/auth"
	"github.com/qta/backend/internal/config"
	cronPkg "github.com/qta/backend/internal/cron"
	"github.com/qta/backend/internal/dashboard"
	"github.com/qta/backend/internal/database"
	"github.com/qta/backend/internal/invoice"
	"github.com/qta/backend/internal/notification"
	"github.com/qta/backend/internal/payment"
	"github.com/qta/backend/internal/room"
	"github.com/qta/backend/internal/router"
	"github.com/qta/backend/internal/tenant"
	"github.com/qta/backend/internal/user"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Gagal memuat konfigurasi: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Gagal koneksi ke database: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(cfg, "db/migrations"); err != nil {
		log.Printf("⚠️ Migration: %v", err)
	}

	// ==========================================
	// Initialize Payment Gateway
	// ==========================================
	var paymentGateway payment.PaymentGatewayService
	if cfg.MidtransServerKey != "" {
		log.Println("💳 Menggunakan Midtrans Payment Gateway")
		paymentGateway = payment.NewMidtransGateway(cfg)
	} else {
		log.Println("💳 Menggunakan Mock Payment Gateway (mode demo)")
		paymentGateway = payment.NewMockGateway()
	}

	// ==========================================
	// Initialize WhatsApp Gateway
	// ==========================================
	whatsAppSvc := notification.NewMockWhatsAppService(db)
	log.Println("📱 Menggunakan Mock WhatsApp Gateway (log ke database)")

	// ==========================================
	// Initialize Repositories
	// ==========================================
	userRepo := user.NewRepository(db)
	roomRepo := room.NewRepository(db)
	tenantRepo := tenant.NewRepository(db)
	invoiceRepo := invoice.NewRepository(db)
	notifRepo := notification.NewRepository(db)

	// ==========================================
	// Initialize Services
	// ==========================================
	authService := auth.NewService(db, cfg)
	userService := user.NewService(userRepo)
	roomService := room.NewService(roomRepo)
	tenantService := tenant.NewService(tenantRepo)
	invoiceService := invoice.NewService(invoiceRepo, tenantRepo, cfg)
	paymentService := payment.NewService(db, paymentGateway)
	dashboardService := dashboard.NewService(db, invoiceRepo)

	// ==========================================
	// Initialize Handlers
	// ==========================================
	handlers := &router.Handlers{
		Auth:         auth.NewHandler(authService),
		User:         user.NewHandler(userService),
		Room:         room.NewHandler(roomService),
		Tenant:       tenant.NewHandler(tenantService),
		Invoice:      invoice.NewHandler(invoiceService),
		Payment:      payment.NewHandler(paymentService),
		Dashboard:    dashboard.NewHandler(dashboardService),
		Notification: notification.NewFiberHandler(notifRepo),
	}

	// ==========================================
	// Initialize Cron Scheduler
	// ==========================================
	scheduler := cronPkg.NewScheduler(invoiceRepo, whatsAppSvc)
	scheduler.Start()
	defer scheduler.Stop()

	// ==========================================
	// Setup Fiber App
	// ==========================================
	app := fiber.New(fiber.Config{
		AppName: "Q-TA API v1.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Terjadi kesalahan pada server",
				"data":    nil,
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.FrontendURL,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	// Register routes
	router.Setup(app, handlers, cfg.JWTSecret)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "Q-TA API",
		})
	})

	// ==========================================
	// Start Server with Graceful Shutdown
	// ==========================================
	go func() {
		addr := fmt.Sprintf(":%s", cfg.AppPort)
		log.Printf("🚀 Q-TA API berjalan di http://localhost%s", addr)
		log.Printf("📖 Environment: %s", cfg.AppEnv)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("❌ Gagal memulai server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⏹ Mematikan server...")
	if err := app.Shutdown(); err != nil {
		log.Printf("❌ Gagal mematikan server: %v", err)
	}
	log.Println("✅ Server berhasil dimatikan")
}
