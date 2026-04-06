package wallet

import (
	"context"
	"fmt"
)

// AWSKMS implements the Wallet Provider interface using AWS Key Management Service.
// Keys are created and stored in KMS; signing is done server-side via KMS API.
type AWSKMS struct {
	region    string
	keyPolicy string
}

type AWSKMSConfig struct {
	Region    string
	KeyPolicy string // Optional custom key policy ARN
}

func NewAWSKMS(cfg AWSKMSConfig) *AWSKMS {
	region := cfg.Region
	if region == "" {
		region = "us-east-1"
	}
	return &AWSKMS{
		region:    region,
		keyPolicy: cfg.KeyPolicy,
	}
}

func (a *AWSKMS) Name() string { return "aws_kms" }

func (a *AWSKMS) CreateWallet(ctx context.Context, req CreateWalletRequest) (*WalletResult, error) {
	// TODO: Implement AWS KMS key creation
	// 1. Create an asymmetric ECC_SECG_P256K1 key in KMS (secp256k1 for Ethereum)
	// 2. Derive the public address from the KMS public key
	// 3. Store the Key ARN mapping
	//
	// aws kms create-key --key-spec ECC_SECG_P256K1 --key-usage SIGN_VERIFY
	// aws kms get-public-key --key-id <keyId>
	// derive Ethereum address from public key

	keyID := fmt.Sprintf("arn:aws:kms:%s:123456789:key/%s-%s", a.region, req.OwnerType, req.OwnerID)

	return &WalletResult{
		WalletID: fmt.Sprintf("kms-wallet-%s", req.OwnerID),
		Address:  "0x" + fmt.Sprintf("%040x", 0), // placeholder
		Network:  req.Network,
		KeyID:    keyID,
	}, nil
}

func (a *AWSKMS) GetAddress(ctx context.Context, walletID string) (string, error) {
	// TODO: Look up the KMS key and derive the address
	return "0x0000000000000000000000000000000000000000", nil
}

func (a *AWSKMS) SignTransaction(ctx context.Context, req SignRequest) (*SignResult, error) {
	// TODO: Implement KMS signing
	// aws kms sign --key-id <keyId> --signing-algorithm ECDSA_SHA_256 --message <txHash>
	return &SignResult{
		SignedTransaction: req.Transaction,
		TransactionHash:   "0x0000000000000000000000000000000000000000000000000000000000000000",
	}, nil
}

func (a *AWSKMS) GetBalance(ctx context.Context, walletID string, network string) (*BalanceResult, error) {
	// TODO: Query blockchain node (e.g. via Alchemy/Infura) for wallet balance
	return &BalanceResult{
		Balance:  "0",
		Currency: "ETH",
		Network:  network,
	}, nil
}
