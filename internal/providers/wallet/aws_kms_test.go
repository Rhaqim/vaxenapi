package wallet_test

import (
	"context"
	"testing"

	"vaxen/api/internal/providers/wallet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAWSKMS() *wallet.AWSKMS {
	return wallet.NewAWSKMS(wallet.AWSKMSConfig{
		Region: "us-east-1",
	})
}

func TestAWSKMS_Name(t *testing.T) {
	assert.Equal(t, "aws_kms", newAWSKMS().Name())
}

func TestAWSKMS_ImplementsProvider(t *testing.T) {
	var _ wallet.Provider = newAWSKMS()
}

func TestAWSKMS_CreateWallet(t *testing.T) {
	kms := newAWSKMS()
	result, err := kms.CreateWallet(context.Background(), wallet.CreateWalletRequest{
		OwnerID:   "org-123",
		OwnerType: "organization",
		Network:   "ethereum",
		Label:     "main wallet",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.WalletID)
	assert.NotEmpty(t, result.Address)
	assert.Equal(t, "ethereum", result.Network)
	assert.Contains(t, result.KeyID, "arn:aws:kms")
}

func TestAWSKMS_GetAddress(t *testing.T) {
	kms := newAWSKMS()
	address, err := kms.GetAddress(context.Background(), "wallet-1")
	require.NoError(t, err)
	assert.NotEmpty(t, address)
	assert.Contains(t, address, "0x")
}

func TestAWSKMS_SignTransaction(t *testing.T) {
	kms := newAWSKMS()
	result, err := kms.SignTransaction(context.Background(), wallet.SignRequest{
		WalletID:    "wallet-1",
		Network:     "ethereum",
		Transaction: []byte("raw-tx-data"),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.SignedTransaction)
	assert.NotEmpty(t, result.TransactionHash)
}

func TestAWSKMS_GetBalance(t *testing.T) {
	kms := newAWSKMS()
	result, err := kms.GetBalance(context.Background(), "wallet-1", "ethereum")
	require.NoError(t, err)
	assert.Equal(t, "0", result.Balance)
	assert.Equal(t, "ethereum", result.Network)
}

func TestAWSKMS_DefaultRegion(t *testing.T) {
	kms := wallet.NewAWSKMS(wallet.AWSKMSConfig{})
	result, err := kms.CreateWallet(context.Background(), wallet.CreateWalletRequest{
		OwnerID:   "org-1",
		OwnerType: "organization",
		Network:   "ethereum",
	})
	require.NoError(t, err)
	assert.Contains(t, result.KeyID, "us-east-1")
}
