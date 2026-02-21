package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type KybStatus string

const (
	KybStatusPending      KybStatus = "PENDING"
	KybStatusApproved     KybStatus = "APPROVED"
	KybStatusRejected     KybStatus = "REJECTED"
	KybStatusRequiresInfo KybStatus = "REQUIRES_INFO"
)

type UserRole string

const (
	UserRoleOwner   UserRole = "OWNER"
	UserRoleManager UserRole = "MANAGER"
	UserRoleFinance UserRole = "FINANCE"
	UserRoleViewer  UserRole = "VIEWER"
)

type WalletType string

const (
	WalletTypeFiat   WalletType = "FIAT"
	WalletTypeCrypto WalletType = "CRYPTO"
)

type TransactionType string

const (
	TransactionTypeDeposit     TransactionType = "DEPOSIT"
	TransactionTypeWithdrawal  TransactionType = "WITHDRAWAL"
	TransactionTypeTransferIn  TransactionType = "TRANSFER_IN"
	TransactionTypeTransferOut TransactionType = "TRANSFER_OUT"
	TransactionTypeConversion  TransactionType = "CONVERSION"
	TransactionTypeFee         TransactionType = "FEE"
)

type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "PENDING"
	TransactionStatusProcessing TransactionStatus = "PROCESSING"
	TransactionStatusCompleted  TransactionStatus = "COMPLETED"
	TransactionStatusFailed     TransactionStatus = "FAILED"
)

type AccountType string

const (
	AccountTypeIBAN  AccountType = "IBAN"
	AccountTypePIX   AccountType = "PIX"
	AccountTypeACH   AccountType = "ACH"
	AccountTypeSWIFT AccountType = "SWIFT"
)

type BeneficiaryType string

const (
	BeneficiaryTypeBank   BeneficiaryType = "BANK"
	BeneficiaryTypeCrypto BeneficiaryType = "CRYPTO"
)

type DepositType string

const (
	DepositTypeFiat   DepositType = "FIAT"
	DepositTypeCrypto DepositType = "CRYPTO"
)

type DepositStatus string

const (
	DepositStatusPending    DepositStatus = "PENDING"
	DepositStatusProcessing DepositStatus = "PROCESSING"
	DepositStatusCompleted  DepositStatus = "COMPLETED"
	DepositStatusFailed     DepositStatus = "FAILED"
)

type PayoutType string

const (
	PayoutTypeBank   PayoutType = "BANK"
	PayoutTypeCrypto PayoutType = "CRYPTO"
)

type PayoutStatus string

const (
	PayoutStatusPending    PayoutStatus = "PENDING"
	PayoutStatusProcessing PayoutStatus = "PROCESSING"
	PayoutStatusCompleted  PayoutStatus = "COMPLETED"
	PayoutStatusFailed     PayoutStatus = "FAILED"
	PayoutStatusCancelled  PayoutStatus = "CANCELLED"
)

type BatchStatus string

const (
	BatchStatusPending    BatchStatus = "PENDING"
	BatchStatusProcessing BatchStatus = "PROCESSING"
	BatchStatusCompleted  BatchStatus = "COMPLETED"
	BatchStatusFailed     BatchStatus = "FAILED"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusCompleted  OrderStatus = "COMPLETED"
	OrderStatusFailed     OrderStatus = "FAILED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
)

type ConversionType string

const (
	ConversionTypeMarket ConversionType = "MARKET"
	ConversionTypeLimit  ConversionType = "LIMIT"
)

type OrderType string

const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"
)

type JournalStatus string

const (
	JournalStatusPending  JournalStatus = "PENDING"
	JournalStatusPosted   JournalStatus = "POSTED"
	JournalStatusReversed JournalStatus = "REVERSED"
)

type ReconciliationStatus string

const (
	ReconciliationStatusPending   ReconciliationStatus = "PENDING"
	ReconciliationStatusCompleted ReconciliationStatus = "COMPLETED"
	ReconciliationStatusFailed    ReconciliationStatus = "FAILED"
)

type StatementType string

const (
	StatementTypePDF StatementType = "PDF"
	StatementTypeCSV StatementType = "CSV"
)

type StatementStatus string

const (
	StatementStatusGenerating StatementStatus = "GENERATING"
	StatementStatusReady      StatementStatus = "READY"
	StatementStatusFailed     StatementStatus = "FAILED"
)

type ComplianceType string

const (
	ComplianceTypeKYB ComplianceType = "KYB"
	ComplianceTypeKYT ComplianceType = "KYT"
	ComplianceTypeAML ComplianceType = "AML"
)

type ComplianceStatus string

const (
	ComplianceStatusPending      ComplianceStatus = "PENDING"
	ComplianceStatusApproved     ComplianceStatus = "APPROVED"
	ComplianceStatusRejected     ComplianceStatus = "REJECTED"
	ComplianceStatusRequiresInfo ComplianceStatus = "REQUIRES_INFO"
)

type ScreeningType string

const (
	ScreeningTypeTransaction ScreeningType = "TRANSACTION"
	ScreeningTypeAddress     ScreeningType = "ADDRESS"
	ScreeningTypeEntity      ScreeningType = "ENTITY"
)

type ScreeningStatus string

const (
	ScreeningStatusClean   ScreeningStatus = "CLEAN"
	ScreeningStatusFlagged ScreeningStatus = "FLAGGED"
	ScreeningStatusBlocked ScreeningStatus = "BLOCKED"
)

type WebhookStatus string

const (
	WebhookStatusPending   WebhookStatus = "PENDING"
	WebhookStatusProcessed WebhookStatus = "PROCESSED"
	WebhookStatusFailed    WebhookStatus = "FAILED"
)

// Organization represents a business organization
type Organization struct {
	ID                 string              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name               string              `gorm:"not null" json:"name"`
	LegalName          string              `gorm:"not null" json:"legalName"`
	RegistrationNumber string              `gorm:"not null" json:"registrationNumber"`
	TaxID              *string             `json:"taxId"`
	Country            string              `gorm:"type:char(2);not null" json:"country"`
	Address            datatypes.JSON      `gorm:"type:jsonb" json:"address"`
	KYBStatus          KybStatus           `gorm:"type:text;default:'PENDING'" json:"kybStatus"`
	KYBSubmittedAt     *time.Time          `json:"kybSubmittedAt"`
	KYBApprovedAt      *time.Time          `json:"kybApprovedAt"`
	Settings           datatypes.JSON      `gorm:"type:jsonb;default:'{}'" json:"settings"`
	CreatedAt          time.Time           `json:"createdAt"`
	UpdatedAt          time.Time           `json:"updatedAt"`
	Users              []User              `gorm:"foreignKey:OrganizationID" json:"-"`
	Wallets            []Wallet            `gorm:"foreignKey:OrganizationID" json:"-"`
	AccountNumbers     []AccountNumber     `gorm:"foreignKey:OrganizationID" json:"-"`
	Beneficiaries      []Beneficiary       `gorm:"foreignKey:OrganizationID" json:"-"`
	Deposits           []Deposit           `gorm:"foreignKey:OrganizationID" json:"-"`
	Payouts            []Payout            `gorm:"foreignKey:OrganizationID" json:"-"`
	PayoutBatches      []PayoutBatch       `gorm:"foreignKey:OrganizationID" json:"-"`
	ConversionOrders   []ConversionOrder   `gorm:"foreignKey:OrganizationID" json:"-"`
	LimitOrders        []LimitOrder        `gorm:"foreignKey:OrganizationID" json:"-"`
	Journals           []Journal           `gorm:"foreignKey:OrganizationID" json:"-"`
	ReconciliationRuns []ReconciliationRun `gorm:"foreignKey:OrganizationID" json:"-"`
	StatementFiles     []StatementFile     `gorm:"foreignKey:OrganizationID" json:"-"`
	ComplianceCases    []ComplianceCase    `gorm:"foreignKey:OrganizationID" json:"-"`
	ScreeningResults   []ScreeningResult   `gorm:"foreignKey:OrganizationID" json:"-"`
	AuditLogs          []AuditLog          `gorm:"foreignKey:OrganizationID" json:"-"`
	WebhookEvents      []WebhookEvent      `gorm:"foreignKey:OrganizationID" json:"-"`
}

func (Organization) TableName() string {
	return "organizations"
}

// User represents a user in the system
type User struct {
	ID             string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email          string         `gorm:"unique;not null" json:"email"`
	FirstName      string         `gorm:"not null" json:"firstName"`
	LastName       string         `gorm:"not null" json:"lastName"`
	PasswordHash   string         `gorm:"not null" json:"-"`
	Role           UserRole       `gorm:"type:text;default:'VIEWER'" json:"role"`
	OrganizationID string         `gorm:"type:uuid;not null" json:"organizationId"`
	MFAEnabled     bool           `gorm:"default:false" json:"mfaEnabled"`
	MFASecret      *string        `json:"mfaSecret"`
	LastLoginAt    *time.Time     `json:"lastLoginAt"`
	IsActive       bool           `gorm:"default:true" json:"isActive"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	Organization   Organization   `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	AuditLogs      []AuditLog     `gorm:"foreignKey:UserID" json:"-"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

// Wallet represents a currency wallet
type Wallet struct {
	ID               string              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID   string              `gorm:"type:uuid;not null" json:"organizationId"`
	Type             WalletType          `gorm:"type:text;not null" json:"type"`
	Currency         string              `gorm:"type:char(3);not null" json:"currency"`
	Balance          decimal.Decimal     `gorm:"type:decimal(20,8);default:0" json:"balance"`
	AvailableBalance decimal.Decimal     `gorm:"type:decimal(20,8);default:0" json:"availableBalance"`
	PendingBalance   decimal.Decimal     `gorm:"type:decimal(20,8);default:0" json:"pendingBalance"`
	IsActive         bool                `gorm:"default:true" json:"isActive"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
	Organization     Organization        `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	Transactions     []WalletTransaction `gorm:"foreignKey:WalletID" json:"-"`
}

func (Wallet) TableName() string {
	return "wallets"
}

// WalletTransaction represents a wallet transaction
type WalletTransaction struct {
	ID          string            `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WalletID    string            `gorm:"type:uuid;not null" json:"walletId"`
	Type        TransactionType   `gorm:"type:text;not null" json:"type"`
	Amount      decimal.Decimal   `gorm:"type:decimal(20,8);not null" json:"amount"`
	Currency    string            `gorm:"type:char(3);not null" json:"currency"`
	Description string            `gorm:"not null" json:"description"`
	Reference   *string           `json:"reference"`
	Status      TransactionStatus `gorm:"type:text;default:'PENDING'" json:"status"`
	Metadata    datatypes.JSON    `gorm:"type:jsonb" json:"metadata"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	Wallet      Wallet            `gorm:"foreignKey:WalletID;references:ID" json:"-"`
}

func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}

// AccountNumber represents a bank account number
type AccountNumber struct {
	ID             string       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string       `gorm:"type:uuid;not null" json:"organizationId"`
	Name           string       `gorm:"not null" json:"name"`
	Currency       string       `gorm:"type:char(3);not null" json:"currency"`
	Type           AccountType  `gorm:"type:text;not null" json:"type"`
	AccountNumber  string       `gorm:"not null" json:"accountNumber"`
	RoutingNumber  *string      `json:"routingNumber"`
	BankCode       *string      `json:"bankCode"`
	BankName       string       `gorm:"not null" json:"bankName"`
	BankCountry    string       `gorm:"type:char(2);not null" json:"bankCountry"`
	IsActive       bool         `gorm:"default:true" json:"isActive"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (AccountNumber) TableName() string {
	return "account_numbers"
}

// Beneficiary represents a payment beneficiary
type Beneficiary struct {
	ID             string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string          `gorm:"type:uuid;not null" json:"organizationId"`
	Name           string          `gorm:"not null" json:"name"`
	Type           BeneficiaryType `gorm:"type:text;not null" json:"type"`
	AccountNumber  *string         `json:"accountNumber"`
	RoutingNumber  *string         `json:"routingNumber"`
	BankCode       *string         `json:"bankCode"`
	BankName       *string         `json:"bankName"`
	BankCountry    *string         `gorm:"type:char(2)" json:"bankCountry"`
	Address        *string         `json:"address"`
	Currency       string          `gorm:"type:char(3);not null" json:"currency"`
	Network        *string         `json:"network"`
	IsActive       bool            `gorm:"default:true" json:"isActive"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	Organization   Organization    `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	Payouts        []Payout        `gorm:"foreignKey:BeneficiaryID" json:"-"`
}

func (Beneficiary) TableName() string {
	return "beneficiaries"
}

// Deposit represents a wallet deposit
type Deposit struct {
	ID                string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID    string          `gorm:"type:uuid;not null" json:"organizationId"`
	WalletID          string          `gorm:"type:uuid;not null" json:"walletId"`
	Type              DepositType     `gorm:"type:text;not null" json:"type"`
	Amount            decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"amount"`
	Currency          string          `gorm:"type:char(3);not null" json:"currency"`
	Reference         *string         `json:"reference"`
	Status            DepositStatus   `gorm:"type:text;default:'PENDING'" json:"status"`
	ProviderReference *string         `json:"providerReference"`
	ExecutedAt        *time.Time      `json:"executedAt"`
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
	Organization      Organization    `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (Deposit) TableName() string {
	return "deposits"
}

// Payout represents a payout transaction
type Payout struct {
	ID             string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string          `gorm:"type:uuid;not null" json:"organizationId"`
	Type           PayoutType      `gorm:"type:text;not null" json:"type"`
	Amount         decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"amount"`
	Currency       string          `gorm:"type:char(3);not null" json:"currency"`
	BeneficiaryID  string          `gorm:"type:uuid;not null" json:"beneficiaryId"`
	Reference      *string         `json:"reference"`
	Description    *string         `json:"description"`
	Status         PayoutStatus    `gorm:"type:text;default:'PENDING'" json:"status"`
	Fee            decimal.Decimal `gorm:"type:decimal(20,8);default:0" json:"fee"`
	ExecutedAt     *time.Time      `json:"executedAt"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	Organization   Organization    `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	Beneficiary    Beneficiary     `gorm:"foreignKey:BeneficiaryID;references:ID" json:"-"`
}

func (Payout) TableName() string {
	return "payouts"
}

// PayoutBatch represents a batch of payouts
type PayoutBatch struct {
	ID             string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string          `gorm:"type:uuid;not null" json:"organizationId"`
	Name           string          `gorm:"not null" json:"name"`
	TotalAmount    decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"totalAmount"`
	TotalCount     int             `gorm:"not null" json:"totalCount"`
	ProcessedCount int             `gorm:"default:0" json:"processedCount"`
	Status         BatchStatus     `gorm:"type:text;default:'PENDING'" json:"status"`
	FileURL        *string         `json:"fileUrl"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	Organization   Organization    `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (PayoutBatch) TableName() string {
	return "payout_batches"
}

// ConversionOrder represents a currency conversion order
type ConversionOrder struct {
	ID             string           `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string           `gorm:"type:uuid;not null" json:"organizationId"`
	FromCurrency   string           `gorm:"type:char(3);not null" json:"fromCurrency"`
	ToCurrency     string           `gorm:"type:char(3);not null" json:"toCurrency"`
	FromAmount     decimal.Decimal  `gorm:"type:decimal(20,8);not null" json:"fromAmount"`
	ToAmount       decimal.Decimal  `gorm:"type:decimal(20,8);not null" json:"toAmount"`
	Rate           decimal.Decimal  `gorm:"type:decimal(20,8);not null" json:"rate"`
	Fee            decimal.Decimal  `gorm:"type:decimal(20,8);default:0" json:"fee"`
	Status         OrderStatus      `gorm:"type:text;default:'PENDING'" json:"status"`
	Type           ConversionType   `gorm:"type:text;default:'MARKET'" json:"type"`
	LimitPrice     *decimal.Decimal `gorm:"type:decimal(20,8)" json:"limitPrice"`
	ExecutedAt     *time.Time       `json:"executedAt"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	Organization   Organization     `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (ConversionOrder) TableName() string {
	return "conversion_orders"
}

// LimitOrder represents a limit order
type LimitOrder struct {
	ID             string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string          `gorm:"type:uuid;not null" json:"organizationId"`
	FromCurrency   string          `gorm:"type:char(3);not null" json:"fromCurrency"`
	ToCurrency     string          `gorm:"type:char(3);not null" json:"toCurrency"`
	Amount         decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"amount"`
	LimitPrice     decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"limitPrice"`
	Type           OrderType       `gorm:"type:text;default:'LIMIT'" json:"type"`
	Status         OrderStatus     `gorm:"type:text;default:'PENDING'" json:"status"`
	ExecutedAt     *time.Time      `json:"executedAt"`
	CancelledAt    *time.Time      `json:"cancelledAt"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	Organization   Organization    `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (LimitOrder) TableName() string {
	return "limit_orders"
}

// Journal represents a journal entry
type Journal struct {
	ID             string        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string        `gorm:"type:uuid;not null" json:"organizationId"`
	Type           string        `gorm:"not null" json:"type"`
	Description    string        `gorm:"not null" json:"description"`
	Reference      *string       `json:"reference"`
	Status         JournalStatus `gorm:"type:text;default:'PENDING'" json:"status"`
	PostedAt       *time.Time    `json:"postedAt"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
	Organization   Organization  `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	Entries        []LedgerEntry `gorm:"foreignKey:JournalID" json:"-"`
}

func (Journal) TableName() string {
	return "journals"
}

// LedgerEntry represents a journal entry line
type LedgerEntry struct {
	ID             string           `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string           `gorm:"type:uuid;not null" json:"organizationId"`
	JournalID      string           `gorm:"type:uuid;not null" json:"journalId"`
	Account        string           `gorm:"not null" json:"account"`
	Debit          *decimal.Decimal `gorm:"type:decimal(20,8)" json:"debit"`
	Credit         *decimal.Decimal `gorm:"type:decimal(20,8)" json:"credit"`
	Description    string           `gorm:"not null" json:"description"`
	Reference      *string          `json:"reference"`
	Metadata       datatypes.JSON   `gorm:"type:jsonb" json:"metadata"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	Journal        Journal          `gorm:"foreignKey:JournalID;references:ID" json:"-"`
}

func (LedgerEntry) TableName() string {
	return "ledger_entries"
}

// ReconciliationRun represents a reconciliation run
type ReconciliationRun struct {
	ID             string               `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string               `gorm:"type:uuid;not null" json:"organizationId"`
	Date           time.Time            `json:"date"`
	Status         ReconciliationStatus `gorm:"type:text;default:'PENDING'" json:"status"`
	Discrepancies  datatypes.JSON       `gorm:"type:jsonb;default:'[]'" json:"discrepancies"`
	CompletedAt    *time.Time           `json:"completedAt"`
	CreatedAt      time.Time            `json:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt"`
	Organization   Organization         `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (ReconciliationRun) TableName() string {
	return "reconciliation_runs"
}

// StatementFile represents a statement file
type StatementFile struct {
	ID             string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string          `gorm:"type:uuid;not null" json:"organizationId"`
	Type           StatementType   `gorm:"type:text;not null" json:"type"`
	Period         datatypes.JSON  `gorm:"type:jsonb;not null" json:"period"`
	Currency       string          `gorm:"type:char(3);not null" json:"currency"`
	FileURL        string          `gorm:"not null" json:"fileUrl"`
	FileSize       int             `gorm:"not null" json:"fileSize"`
	Status         StatementStatus `gorm:"type:text;default:'GENERATING'" json:"status"`
	GeneratedAt    *time.Time      `json:"generatedAt"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	Organization   Organization    `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (StatementFile) TableName() string {
	return "statement_files"
}

// ComplianceCase represents a compliance case
type ComplianceCase struct {
	ID             string           `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string           `gorm:"type:uuid;not null" json:"organizationId"`
	Type           ComplianceType   `gorm:"type:text;not null" json:"type"`
	Status         ComplianceStatus `gorm:"type:text;default:'PENDING'" json:"status"`
	Provider       string           `gorm:"not null" json:"provider"`
	ProviderCaseID string           `gorm:"not null" json:"providerCaseId"`
	SubmittedAt    time.Time        `json:"submittedAt"`
	CompletedAt    *time.Time       `json:"completedAt"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	Organization   Organization     `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (ComplianceCase) TableName() string {
	return "compliance_cases"
}

// ScreeningResult represents a screening result
type ScreeningResult struct {
	ID             string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string          `gorm:"type:uuid;not null" json:"organizationId"`
	Type           ScreeningType   `gorm:"type:text;not null" json:"type"`
	Status         ScreeningStatus `gorm:"type:text;default:'CLEAN'" json:"status"`
	RiskScore      int             `gorm:"not null" json:"riskScore"`
	Provider       string          `gorm:"not null" json:"provider"`
	Details        datatypes.JSON  `gorm:"type:jsonb;not null" json:"details"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	Organization   Organization    `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (ScreeningResult) TableName() string {
	return "screening_results"
}

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID             string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string         `gorm:"type:uuid;not null" json:"organizationId"`
	UserID         *string        `gorm:"type:uuid" json:"userId"`
	Action         string         `gorm:"not null" json:"action"`
	Resource       string         `gorm:"not null" json:"resource"`
	ResourceID     *string        `json:"resourceId"`
	Details        datatypes.JSON `gorm:"type:jsonb" json:"details"`
	IPAddress      *string        `json:"ipAddress"`
	UserAgent      *string        `json:"userAgent"`
	CreatedAt      time.Time      `json:"createdAt"`
	Organization   Organization   `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	User           *User          `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// WebhookEvent represents a webhook event
type WebhookEvent struct {
	ID             string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string         `gorm:"type:uuid;not null" json:"organizationId"`
	Provider       string         `gorm:"not null" json:"provider"`
	EventType      string         `gorm:"not null" json:"eventType"`
	Payload        datatypes.JSON `gorm:"type:jsonb;not null" json:"payload"`
	Status         WebhookStatus  `gorm:"type:text;default:'PENDING'" json:"status"`
	ProcessedAt    *time.Time     `json:"processedAt"`
	RetryCount     int            `gorm:"default:0" json:"retryCount"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	Organization   Organization   `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (WebhookEvent) TableName() string {
	return "webhook_events"
}
