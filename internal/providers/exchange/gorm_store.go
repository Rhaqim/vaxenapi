package exchange

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// GormRateStore implements RateStore by reading from the exchange_rates table.
type GormRateStore struct {
	db *gorm.DB
}

func NewGormRateStore(db *gorm.DB) *GormRateStore {
	return &GormRateStore{db: db}
}

// exchangeRateRow mirrors the models.ExchangeRate fields we need.
// Defined here to avoid an import cycle with models.
type exchangeRateRow struct {
	FromCurrency string          `gorm:"column:from_currency"`
	ToCurrency   string          `gorm:"column:to_currency"`
	Rate         decimal.Decimal `gorm:"column:rate"`
	Spread       decimal.Decimal `gorm:"column:spread"`
	IsActive     bool            `gorm:"column:is_active"`
}

func (exchangeRateRow) TableName() string { return "exchange_rates" }

func (s *GormRateStore) GetRate(from, to string) (decimal.Decimal, decimal.Decimal, bool) {
	var row exchangeRateRow
	err := s.db.Where("from_currency = ? AND to_currency = ? AND is_active = ?", from, to, true).First(&row).Error
	if err != nil {
		return decimal.Zero, decimal.Zero, false
	}
	return row.Rate, row.Spread, true
}

func (s *GormRateStore) ListPairs() []CurrencyPair {
	var rows []exchangeRateRow
	s.db.Where("is_active = ?", true).Find(&rows)

	pairs := make([]CurrencyPair, 0, len(rows))
	for _, r := range rows {
		pairType := "fiat"
		if isCrypto(r.FromCurrency) || isCrypto(r.ToCurrency) {
			if isCrypto(r.FromCurrency) && isCrypto(r.ToCurrency) {
				pairType = "crypto"
			} else {
				pairType = "cross"
			}
		}
		pairs = append(pairs, CurrencyPair{
			From:      r.FromCurrency,
			To:        r.ToCurrency,
			MinAmount: decimal.NewFromInt(1),
			MaxAmount: decimal.NewFromInt(1000000),
			Type:      pairType,
		})
	}
	return pairs
}

func isCrypto(currency string) bool {
	switch currency {
	case "BTC", "ETH", "USDT", "USDC", "SOL", "MATIC", "AVAX", "DOT", "ADA", "XRP", "BNB", "LINK":
		return true
	}
	return false
}
