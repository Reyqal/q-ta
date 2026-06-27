package dashboard

import (
	"github.com/qta/backend/internal/invoice"
	"gorm.io/gorm"
)

// DashboardSummary berisi ringkasan data untuk dashboard admin.
type DashboardSummary struct {
	TotalPendapatanBulanIni int64                  `json:"total_pendapatan_bulan_ini"`
	TotalPajakBulanIni      int64                  `json:"total_pajak_bulan_ini"`
	TotalKamar              int64                  `json:"total_kamar"`
	KamarTerisi             int64                  `json:"kamar_terisi"`
	KamarKosong             int64                  `json:"kamar_kosong"`
	TagihanStats            *invoice.StatusCounts  `json:"tagihan_stats"`
	PendapatanBulanan       []invoice.MonthlyStats `json:"pendapatan_bulanan"`
	TagihanTerbaru          interface{}            `json:"tagihan_terbaru"`
}

// Service mengelola logika bisnis untuk dashboard.
type Service struct {
	db          *gorm.DB
	invoiceRepo *invoice.Repository
}

// NewService membuat instance baru DashboardService.
func NewService(db *gorm.DB, invoiceRepo *invoice.Repository) *Service {
	return &Service{db: db, invoiceRepo: invoiceRepo}
}

// GetSummary mengambil ringkasan data untuk dashboard admin.
func (s *Service) GetSummary() (*DashboardSummary, error) {
	// Total pendapatan dan pajak bulan ini
	totals, err := s.invoiceRepo.GetCurrentMonthTotals()
	if err != nil {
		return nil, err
	}

	// Statistik kamar
	var totalKamar, kamarTerisi, kamarKosong int64
	s.db.Table("rooms").Count(&totalKamar)
	s.db.Table("rooms").Where("status = ?", "occupied").Count(&kamarTerisi)
	s.db.Table("rooms").Where("status = ?", "available").Count(&kamarKosong)

	// Statistik tagihan
	statusCounts, err := s.invoiceRepo.GetStatusCounts()
	if err != nil {
		return nil, err
	}

	// Pendapatan bulanan (6 bulan terakhir)
	monthlyStats, err := s.invoiceRepo.GetMonthlyStats()
	if err != nil {
		return nil, err
	}

	// Tagihan terbaru
	recentInvoices, err := s.invoiceRepo.FindRecent(10)
	if err != nil {
		return nil, err
	}

	return &DashboardSummary{
		TotalPendapatanBulanIni: totals.TotalPendapatan,
		TotalPajakBulanIni:      totals.TotalPajak,
		TotalKamar:              totalKamar,
		KamarTerisi:             kamarTerisi,
		KamarKosong:             kamarKosong,
		TagihanStats:            statusCounts,
		PendapatanBulanan:       monthlyStats,
		TagihanTerbaru:          recentInvoices,
	}, nil
}
