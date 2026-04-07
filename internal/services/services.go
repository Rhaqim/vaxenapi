package services

import (
	"vaxen/api/internal/config"
	"vaxen/api/internal/providers"

	"gorm.io/gorm"
)

// Container holds all service instances. Passed to handlers for dependency injection.
type Container struct {
	Auth         *AuthService
	KYC          *KYCService
	MFA          *MFAService
	Wallet       *WalletService
	Exchange     *ExchangeService
	Payment      *PaymentService
	Approval     *ApprovalService
	Admin        *AdminService
	Email        *EmailService
	Beneficiary  *BeneficiaryService
	Organization *OrganizationService
	Report       *ReportService
	Account      *AccountService
	Webhook      *WebhookService
}

func NewContainer(db *gorm.DB, cfg *config.Config, reg *providers.Registry) *Container {
	approval := NewApprovalService(db)
	emailSvc := NewEmailService(cfg, reg)
	return &Container{
		Auth:         NewAuthService(db, cfg, emailSvc),
		KYC:          NewKYCService(db, reg),
		MFA:          NewMFAService(db, reg),
		Wallet:       NewWalletService(db, reg),
		Exchange:     NewExchangeService(db, reg),
		Payment:      NewPaymentService(db, reg, approval),
		Approval:     approval,
		Admin:        NewAdminService(db),
		Email:        emailSvc,
		Beneficiary:  NewBeneficiaryService(db, reg),
		Organization: NewOrganizationService(db),
		Report:       NewReportService(db),
		Account:      NewAccountService(db),
		Webhook:      NewWebhookService(db),
	}
}
