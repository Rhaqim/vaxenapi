package services

import (
	"errors"

	"vaxen/api/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ReportService struct {
	db *gorm.DB
}

func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

// --- Statements ---

func (s *ReportService) GetStatements(orgID string) ([]models.StatementFile, error) {
	var statements []models.StatementFile
	if err := s.db.Where("organization_id = ?", orgID).Order("created_at desc").Find(&statements).Error; err != nil {
		return nil, err
	}
	return statements, nil
}

func (s *ReportService) GetStatement(orgID, id string) (*models.StatementFile, error) {
	var stmt models.StatementFile
	if err := s.db.Where("id = ? AND organization_id = ?", id, orgID).First(&stmt).Error; err != nil {
		return nil, errors.New("statement not found")
	}
	return &stmt, nil
}

// --- Audit Logs ---

func (s *ReportService) GetAuditLogs(orgID string, page, limit int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64
	s.db.Model(&models.AuditLog{}).Where("organization_id = ?", orgID).Count(&total)
	if err := s.db.Where("organization_id = ?", orgID).
		Order("created_at desc").
		Offset((page - 1) * limit).Limit(limit).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// --- Transaction Report ---

type TransactionReportResult struct {
	TotalTransactions int64           `json:"totalTransactions"`
	TotalVolume       decimal.Decimal `json:"totalVolume"`
	Currency          string          `json:"currency"`
}

func (s *ReportService) GetTransactionReport(orgID string) (*TransactionReportResult, error) {
	var totalCount int64
	s.db.Model(&models.WalletTransaction{}).
		Joins("JOIN wallets ON wallets.id = wallet_transactions.wallet_id").
		Where("wallets.organization_id = ?", orgID).
		Count(&totalCount)

	var totalVolume decimal.Decimal
	row := s.db.Model(&models.WalletTransaction{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Joins("JOIN wallets ON wallets.id = wallet_transactions.wallet_id").
		Where("wallets.organization_id = ?", orgID).
		Row()
	row.Scan(&totalVolume)

	return &TransactionReportResult{
		TotalTransactions: totalCount,
		TotalVolume:       totalVolume,
	}, nil
}

// --- Balance Report ---

type BalanceEntry struct {
	Currency         string          `json:"currency"`
	Balance          decimal.Decimal `json:"balance"`
	AvailableBalance decimal.Decimal `json:"availableBalance"`
}

func (s *ReportService) GetBalanceReport(orgID string) ([]BalanceEntry, error) {
	var wallets []models.Wallet
	if err := s.db.Where("organization_id = ? AND is_active = ?", orgID, true).Find(&wallets).Error; err != nil {
		return nil, err
	}

	entries := make([]BalanceEntry, 0, len(wallets))
	for _, w := range wallets {
		entries = append(entries, BalanceEntry{
			Currency:         w.Currency,
			Balance:          w.Balance,
			AvailableBalance: w.AvailableBalance,
		})
	}
	return entries, nil
}

// --- FX P&L Report ---

type FxPnlResult struct {
	TotalConversions int64           `json:"totalConversions"`
	TotalFees        decimal.Decimal `json:"totalFees"`
}

func (s *ReportService) GetFxPnlReport(orgID string) (*FxPnlResult, error) {
	var count int64
	s.db.Model(&models.ConversionOrder{}).Where("organization_id = ?", orgID).Count(&count)

	var totalFees decimal.Decimal
	row := s.db.Model(&models.ConversionOrder{}).
		Select("COALESCE(SUM(fee), 0) as total").
		Where("organization_id = ?", orgID).
		Row()
	row.Scan(&totalFees)

	return &FxPnlResult{
		TotalConversions: count,
		TotalFees:        totalFees,
	}, nil
}
