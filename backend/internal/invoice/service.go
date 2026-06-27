package invoice

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/qta/backend/internal/config"
	"github.com/qta/backend/internal/models"
	tenantPkg "github.com/qta/backend/internal/tenant"
)

// Service mengelola logika bisnis untuk Invoice.
type Service struct {
	repo       *Repository
	tenantRepo *tenantPkg.Repository
	cfg        *config.Config
}

// NewService membuat instance baru InvoiceService.
func NewService(repo *Repository, tenantRepo *tenantPkg.Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, tenantRepo: tenantRepo, cfg: cfg}
}

// GetAll mengambil semua invoice (untuk admin).
func (s *Service) GetAll() ([]models.Invoice, error) {
	return s.repo.FindAll()
}

// GetByTenantUserID mengambil invoice berdasarkan user_id penghuni.
func (s *Service) GetByTenantUserID(userID uint) ([]models.Invoice, error) {
	tenant, err := s.tenantRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("data penghuni tidak ditemukan")
	}
	return s.repo.FindByTenantID(tenant.ID)
}

// GetByID mengambil invoice berdasarkan ID.
func (s *Service) GetByID(id uint) (*models.Invoice, error) {
	invoice, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("tagihan tidak ditemukan")
	}
	return invoice, nil
}

// Create membuat invoice baru dengan perhitungan pajak otomatis.
func (s *Service) Create(tenantID uint, period string, amount int64, dueDate time.Time) (*models.Invoice, error) {
	if amount <= 0 {
		return nil, errors.New("jumlah tagihan harus lebih dari 0")
	}

	// Cek apakah sudah ada tagihan untuk tenant dan periode yang sama
	existing, _ := s.repo.FindByTenantAndPeriod(tenantID, period)
	if existing != nil {
		return nil, fmt.Errorf("tagihan untuk periode %s sudah ada", period)
	}

	// Hitung pemisahan pajak
	taxPortion := int64(math.Round(float64(amount) * s.cfg.TaxPercentage / 100))
	netPortion := amount - taxPortion

	invoice := &models.Invoice{
		TenantID:   tenantID,
		Period:     period,
		Amount:     amount,
		TaxPortion: taxPortion,
		NetPortion: netPortion,
		Status:     "unpaid",
		DueDate:    dueDate,
	}

	if err := s.repo.Create(invoice); err != nil {
		return nil, errors.New("gagal membuat tagihan")
	}

	return invoice, nil
}

// GenerateMonthly membuat tagihan bulanan untuk semua penghuni aktif.
func (s *Service) GenerateMonthly() (int, error) {
	tenants, err := s.tenantRepo.FindAllActive()
	if err != nil {
		return 0, errors.New("gagal mengambil data penghuni aktif")
	}

	now := time.Now()
	period := fmt.Sprintf("%d-%02d", now.Year(), now.Month())
	dueDate := time.Date(now.Year(), now.Month(), 10, 0, 0, 0, 0, time.Local)
	if dueDate.Before(now) {
		// Jika sudah lewat tanggal 10, set due date bulan depan
		dueDate = dueDate.AddDate(0, 1, 0)
	}

	count := 0
	for _, t := range tenants {
		// Skip jika sudah ada tagihan untuk periode ini
		existing, _ := s.repo.FindByTenantAndPeriod(t.ID, period)
		if existing != nil {
			continue
		}

		amount := t.Room.RentAmount
		taxPortion := int64(math.Round(float64(amount) * s.cfg.TaxPercentage / 100))
		netPortion := amount - taxPortion

		invoice := &models.Invoice{
			TenantID:   t.ID,
			Period:     period,
			Amount:     amount,
			TaxPortion: taxPortion,
			NetPortion: netPortion,
			Status:     "unpaid",
			DueDate:    dueDate,
		}

		if err := s.repo.Create(invoice); err == nil {
			count++
		}
	}

	return count, nil
}

// UpdateStatus mengubah status tagihan (admin: set ke issue atau unpaid).
func (s *Service) UpdateStatus(id uint, status string) (*models.Invoice, error) {
	if status != "issue" && status != "unpaid" {
		return nil, errors.New("status harus 'issue' atau 'unpaid'")
	}

	invoice, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("tagihan tidak ditemukan")
	}

	if invoice.Status == "paid" {
		return nil, errors.New("tidak bisa mengubah status tagihan yang sudah lunas")
	}

	invoice.Status = status
	if err := s.repo.Update(invoice); err != nil {
		return nil, errors.New("gagal memperbarui status tagihan")
	}

	return invoice, nil
}

// ConfirmManualPayment mengkonfirmasi pembayaran manual (cash/transfer).
func (s *Service) ConfirmManualPayment(id uint, paymentMethod string) (*models.Invoice, error) {
	if paymentMethod != "cash" && paymentMethod != "transfer" {
		return nil, errors.New("metode pembayaran harus 'cash' atau 'transfer'")
	}

	invoice, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("tagihan tidak ditemukan")
	}

	if invoice.Status == "paid" {
		return nil, errors.New("tagihan sudah lunas")
	}

	// Hitung ulang pemisahan pajak
	taxPortion := int64(math.Round(float64(invoice.Amount) * s.cfg.TaxPercentage / 100))
	netPortion := invoice.Amount - taxPortion

	now := time.Now()
	invoice.Status = "paid"
	invoice.PaidAt = &now
	invoice.TaxPortion = taxPortion
	invoice.NetPortion = netPortion

	if err := s.repo.Update(invoice); err != nil {
		return nil, errors.New("gagal mengkonfirmasi pembayaran")
	}

	return invoice, nil
}

// MarkAsPaid menandai invoice sebagai lunas (dipanggil dari webhook).
func (s *Service) MarkAsPaid(id uint) (*models.Invoice, error) {
	invoice, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("tagihan tidak ditemukan")
	}

	taxPortion := int64(math.Round(float64(invoice.Amount) * s.cfg.TaxPercentage / 100))
	netPortion := invoice.Amount - taxPortion

	now := time.Now()
	invoice.Status = "paid"
	invoice.PaidAt = &now
	invoice.TaxPortion = taxPortion
	invoice.NetPortion = netPortion

	if err := s.repo.Update(invoice); err != nil {
		return nil, errors.New("gagal memperbarui status pembayaran")
	}

	return invoice, nil
}
