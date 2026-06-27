package invoice

import (
	"gorm.io/gorm"
	"github.com/qta/backend/internal/models"
)

// Repository mengelola akses data Invoice ke database.
type Repository struct {
	db *gorm.DB
}

// NewRepository membuat instance baru InvoiceRepository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindAll mengambil semua invoice dengan relasi tenant, user, dan room.
func (r *Repository) FindAll() ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := r.db.Preload("Tenant").Preload("Tenant.User").Preload("Tenant.Room").
		Order("created_at DESC").Find(&invoices).Error
	return invoices, err
}

// FindByTenantID mengambil invoice berdasarkan tenant_id.
func (r *Repository) FindByTenantID(tenantID uint) ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := r.db.Where("tenant_id = ?", tenantID).
		Preload("Tenant").Preload("Tenant.User").Preload("Tenant.Room").
		Order("created_at DESC").Find(&invoices).Error
	return invoices, err
}

// FindByID mengambil invoice berdasarkan ID.
func (r *Repository) FindByID(id uint) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.Preload("Tenant").Preload("Tenant.User").Preload("Tenant.Room").
		Preload("Transactions").First(&invoice, id).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

// Create menyimpan invoice baru ke database.
func (r *Repository) Create(invoice *models.Invoice) error {
	return r.db.Create(invoice).Error
}

// Update memperbarui invoice di database.
func (r *Repository) Update(invoice *models.Invoice) error {
	return r.db.Save(invoice).Error
}

// FindByTenantAndPeriod mencari invoice berdasarkan tenant dan periode.
func (r *Repository) FindByTenantAndPeriod(tenantID uint, period string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.Where("tenant_id = ? AND period = ?", tenantID, period).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

// FindOverdue mengambil invoice yang sudah jatuh tempo dan belum dibayar.
func (r *Repository) FindOverdue() ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := r.db.Where("status = 'unpaid' AND due_date <= CURRENT_DATE").
		Preload("Tenant").Preload("Tenant.User").Preload("Tenant.Room").
		Find(&invoices).Error
	return invoices, err
}

// GetMonthlyStats mengambil statistik pendapatan bulanan (6 bulan terakhir).
type MonthlyStats struct {
	Bulan      string `json:"bulan"`
	Pendapatan int64  `json:"pendapatan"`
	Pajak      int64  `json:"pajak"`
}

func (r *Repository) GetMonthlyStats() ([]MonthlyStats, error) {
	var stats []MonthlyStats
	err := r.db.Raw(`
		SELECT 
			period AS bulan,
			COALESCE(SUM(net_portion), 0) AS pendapatan,
			COALESCE(SUM(tax_portion), 0) AS pajak
		FROM invoices 
		WHERE status = 'paid' 
			AND paid_at >= CURRENT_DATE - INTERVAL '6 months'
		GROUP BY period
		ORDER BY period ASC
	`).Scan(&stats).Error
	return stats, err
}

// GetCurrentMonthTotals mengambil total pendapatan dan pajak bulan ini.
type CurrentMonthTotals struct {
	TotalPendapatan int64 `json:"total_pendapatan"`
	TotalPajak      int64 `json:"total_pajak"`
}

func (r *Repository) GetCurrentMonthTotals() (*CurrentMonthTotals, error) {
	var totals CurrentMonthTotals
	err := r.db.Raw(`
		SELECT 
			COALESCE(SUM(net_portion), 0) AS total_pendapatan,
			COALESCE(SUM(tax_portion), 0) AS total_pajak
		FROM invoices 
		WHERE status = 'paid' 
			AND TO_CHAR(paid_at, 'YYYY-MM') = TO_CHAR(CURRENT_DATE, 'YYYY-MM')
	`).Scan(&totals).Error
	return &totals, err
}

// GetStatusCounts mengambil jumlah tagihan per status.
type StatusCounts struct {
	Lunas      int64 `json:"lunas"`
	BelumBayar int64 `json:"belum_bayar"`
	AdaKendala int64 `json:"ada_kendala"`
}

func (r *Repository) GetStatusCounts() (*StatusCounts, error) {
	var counts StatusCounts
	err := r.db.Raw(`
		SELECT 
			COUNT(CASE WHEN status = 'paid' THEN 1 END) AS lunas,
			COUNT(CASE WHEN status = 'unpaid' THEN 1 END) AS belum_bayar,
			COUNT(CASE WHEN status = 'issue' THEN 1 END) AS ada_kendala
		FROM invoices
	`).Scan(&counts).Error
	return &counts, err
}

// FindRecent mengambil N invoice terbaru.
func (r *Repository) FindRecent(limit int) ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := r.db.Preload("Tenant").Preload("Tenant.User").Preload("Tenant.Room").
		Order("created_at DESC").Limit(limit).Find(&invoices).Error
	return invoices, err
}
