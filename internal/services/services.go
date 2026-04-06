package services

import (
	"vaxen/api/internal/config"
	"vaxen/api/internal/providers"

	"gorm.io/gorm"
)

// Container holds all service instances. Passed to handlers for dependency injection.
type Container struct {
	Auth     *AuthService
	KYC      *KYCService
	MFA      *MFAService
	Wallet   *WalletService
	Exchange *ExchangeService
	Payment  *PaymentService
	Approval *ApprovalService
	Admin    *AdminService
}

func NewContainer(db *gorm.DB, cfg *config.Config, reg *providers.Registry) *Container {
	approval := NewApprovalService(db)
	return &Container{
		Auth:     NewAuthService(db, cfg),
		KYC:      NewKYCService(db, reg),
		MFA:      NewMFAService(db, reg),
		Wallet:   NewWalletService(db, reg),
		Exchange: NewExchangeService(db, reg),
		Payment:  NewPaymentService(db, reg, approval),
		Approval: approval,
		Admin:    NewAdminService(db),
	}
}
