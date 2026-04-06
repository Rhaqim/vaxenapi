package services

import (
	"context"
	"errors"

	"vaxen/api/internal/models"
	"vaxen/api/internal/providers"
	"vaxen/api/internal/providers/wallet"

	"gorm.io/gorm"
)

type WalletService struct {
	db  *gorm.DB
	reg *providers.Registry
}

func NewWalletService(db *gorm.DB, reg *providers.Registry) *WalletService {
	return &WalletService{db: db, reg: reg}
}

// CreateWeb3Wallet creates a new blockchain wallet for an organization.
func (s *WalletService) CreateWeb3Wallet(ctx context.Context, orgID string, network string, label string) (*models.Web3Wallet, error) {
	provider, err := s.reg.RequireWallet()
	if err != nil {
		return nil, err
	}

	result, err := provider.CreateWallet(ctx, wallet.CreateWalletRequest{
		OwnerID:   orgID,
		OwnerType: "organization",
		Network:   network,
		Label:     label,
	})
	if err != nil {
		return nil, err
	}

	w := models.Web3Wallet{
		OrganizationID: orgID,
		Address:        result.Address,
		Network:        result.Network,
		KeyID:          result.KeyID,
		Label:          label,
		IsActive:       true,
	}
	if err := s.db.Create(&w).Error; err != nil {
		return nil, err
	}

	return &w, nil
}

// CreateDefaultWallets creates the standard set of wallets for a new org (fiat + crypto).
func (s *WalletService) CreateDefaultWallets(ctx context.Context, orgID string) error {
	fiatCurrencies := []string{"USD", "EUR", "GBP"}
	for _, currency := range fiatCurrencies {
		w := models.Wallet{
			OrganizationID: orgID,
			Type:           models.WalletTypeFiat,
			Currency:       currency,
			IsActive:       true,
		}
		if err := s.db.Create(&w).Error; err != nil {
			return err
		}
	}

	networks := []string{"ethereum", "polygon"}
	for _, network := range networks {
		if _, err := s.CreateWeb3Wallet(ctx, orgID, network, network+" default"); err != nil {
			return err
		}
	}

	return nil
}

// GetWallets returns all fiat wallets for an organization.
func (s *WalletService) GetWallets(orgID string) ([]models.Wallet, error) {
	var wallets []models.Wallet
	if err := s.db.Where("organization_id = ?", orgID).Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

// GetWallet returns a specific wallet.
func (s *WalletService) GetWallet(orgID string, walletID string) (*models.Wallet, error) {
	var w models.Wallet
	if err := s.db.Where("id = ? AND organization_id = ?", walletID, orgID).First(&w).Error; err != nil {
		return nil, errors.New("wallet not found")
	}
	return &w, nil
}

// GetWeb3Wallets returns all Web3 wallets for an organization.
func (s *WalletService) GetWeb3Wallets(orgID string) ([]models.Web3Wallet, error) {
	var wallets []models.Web3Wallet
	if err := s.db.Where("organization_id = ?", orgID).Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

// GetWeb3Balance queries the on-chain balance via the wallet provider.
func (s *WalletService) GetWeb3Balance(ctx context.Context, walletID string) (*wallet.BalanceResult, error) {
	provider, err := s.reg.RequireWallet()
	if err != nil {
		return nil, err
	}

	var w models.Web3Wallet
	if err := s.db.First(&w, "id = ?", walletID).Error; err != nil {
		return nil, errors.New("web3 wallet not found")
	}

	return provider.GetBalance(ctx, w.KeyID, w.Network)
}
