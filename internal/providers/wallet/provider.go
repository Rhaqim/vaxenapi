package wallet

import "context"

// Provider defines the interface for Web3 wallet creation and key management.
// Implementations can wrap AWS KMS, HashiCorp Vault, Fireblocks, etc.
type Provider interface {
	Name() string

	// CreateWallet generates a new blockchain wallet and stores the keys securely.
	CreateWallet(ctx context.Context, req CreateWalletRequest) (*WalletResult, error)

	// GetAddress returns the public address for a managed wallet.
	GetAddress(ctx context.Context, walletID string) (string, error)

	// SignTransaction signs a raw transaction with the managed private key.
	SignTransaction(ctx context.Context, req SignRequest) (*SignResult, error)

	// GetBalance queries the on-chain balance for a wallet.
	GetBalance(ctx context.Context, walletID string, network string) (*BalanceResult, error)
}

type CreateWalletRequest struct {
	OwnerID   string // Organization or user ID
	OwnerType string // "organization" or "user"
	Network   string // e.g. "ethereum", "polygon", "bitcoin"
	Label     string // Human-readable label
}

type WalletResult struct {
	WalletID  string // Internal reference
	Address   string // Public blockchain address
	Network   string
	KeyID     string // Provider's key identifier (e.g. AWS KMS Key ARN)
}

type SignRequest struct {
	WalletID    string
	Network     string
	Transaction []byte // Raw unsigned transaction
}

type SignResult struct {
	SignedTransaction []byte
	TransactionHash  string
}

type BalanceResult struct {
	Balance  string // Decimal string
	Currency string // e.g. "ETH", "MATIC", "BTC"
	Network  string
}
