package services_test

import (
	"context"
	"testing"

	"vaxen/api/internal/models"
	"vaxen/api/internal/providers"
	"vaxen/api/internal/providers/wallet"
	"vaxen/api/internal/services"
	"vaxen/api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupWalletService(t *testing.T) (*services.WalletService, models.Organization) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	reg := providers.NewRegistry()
	reg.SetWallet(wallet.NewAWSKMS(wallet.AWSKMSConfig{}))
	svc := services.NewWalletService(db, reg)
	org := testutil.SeedOrganization(t, db)
	return svc, org
}

func TestWalletService_CreateWeb3Wallet(t *testing.T) {
	svc, org := setupWalletService(t)

	w, err := svc.CreateWeb3Wallet(context.Background(), org.ID, "ethereum", "main wallet")
	require.NoError(t, err)
	assert.NotEmpty(t, w.ID)
	assert.Equal(t, org.ID, w.OrganizationID)
	assert.Equal(t, "ethereum", w.Network)
	assert.NotEmpty(t, w.Address)
	assert.True(t, w.IsActive)
}

func TestWalletService_CreateDefaultWallets(t *testing.T) {
	svc, org := setupWalletService(t)

	err := svc.CreateDefaultWallets(context.Background(), org.ID)
	require.NoError(t, err)

	// Check fiat wallets were created
	wallets, err := svc.GetWallets(org.ID)
	require.NoError(t, err)
	assert.Len(t, wallets, 3, "should have USD, EUR, GBP fiat wallets")

	currencies := make(map[string]bool)
	for _, w := range wallets {
		currencies[w.Currency] = true
	}
	assert.True(t, currencies["USD"])
	assert.True(t, currencies["EUR"])
	assert.True(t, currencies["GBP"])

	// Check web3 wallets were created
	web3, err := svc.GetWeb3Wallets(org.ID)
	require.NoError(t, err)
	assert.Len(t, web3, 2, "should have ethereum and polygon wallets")

	networks := make(map[string]bool)
	for _, w := range web3 {
		networks[w.Network] = true
	}
	assert.True(t, networks["ethereum"])
	assert.True(t, networks["polygon"])
}

func TestWalletService_GetWallet(t *testing.T) {
	svc, org := setupWalletService(t)
	svc.CreateDefaultWallets(context.Background(), org.ID)

	wallets, _ := svc.GetWallets(org.ID)
	require.NotEmpty(t, wallets)

	w, err := svc.GetWallet(org.ID, wallets[0].ID)
	require.NoError(t, err)
	assert.Equal(t, wallets[0].ID, w.ID)
}

func TestWalletService_GetWallet_NotFound(t *testing.T) {
	svc, org := setupWalletService(t)

	_, err := svc.GetWallet(org.ID, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestWalletService_GetWallet_WrongOrg(t *testing.T) {
	svc, org := setupWalletService(t)
	svc.CreateDefaultWallets(context.Background(), org.ID)

	wallets, _ := svc.GetWallets(org.ID)
	require.NotEmpty(t, wallets)

	_, err := svc.GetWallet("other-org-id", wallets[0].ID)
	assert.Error(t, err, "should not return wallet belonging to different org")
}

func TestWalletService_GetWeb3Balance(t *testing.T) {
	svc, org := setupWalletService(t)

	w, err := svc.CreateWeb3Wallet(context.Background(), org.ID, "ethereum", "test")
	require.NoError(t, err)

	balance, err := svc.GetWeb3Balance(context.Background(), w.ID)
	require.NoError(t, err)
	assert.Equal(t, "0", balance.Balance)
	assert.Equal(t, "ethereum", balance.Network)
}

func TestWalletService_NilProvider(t *testing.T) {
	db := testutil.SetupTestDB(t)
	reg := providers.NewRegistry()
	svc := services.NewWalletService(db, reg)

	_, err := svc.CreateWeb3Wallet(context.Background(), "org-1", "ethereum", "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not configured")
}
