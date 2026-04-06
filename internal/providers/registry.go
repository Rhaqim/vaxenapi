package providers

import (
	"fmt"
	"sync"

	"vaxen/api/internal/providers/exchange"
	"vaxen/api/internal/providers/kyc"
	"vaxen/api/internal/providers/mfa"
	"vaxen/api/internal/providers/payment"
	"vaxen/api/internal/providers/wallet"
)

// Registry holds the active provider for each domain.
// Providers can be swapped at runtime via the Set* methods.
type Registry struct {
	mu sync.RWMutex

	kyc      kyc.Provider
	mfa      mfa.Provider
	wallet   wallet.Provider
	exchange exchange.Provider
	payment  payment.Provider
}

func NewRegistry() *Registry {
	return &Registry{}
}

// --- KYC ---

func (r *Registry) KYC() kyc.Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.kyc
}

func (r *Registry) SetKYC(p kyc.Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.kyc = p
}

// --- MFA ---

func (r *Registry) MFA() mfa.Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.mfa
}

func (r *Registry) SetMFA(p mfa.Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.mfa = p
}

// --- Wallet ---

func (r *Registry) Wallet() wallet.Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.wallet
}

func (r *Registry) SetWallet(p wallet.Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.wallet = p
}

// --- Exchange ---

func (r *Registry) Exchange() exchange.Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.exchange
}

func (r *Registry) SetExchange(p exchange.Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.exchange = p
}

// --- Payment ---

func (r *Registry) Payment() payment.Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.payment
}

func (r *Registry) SetPayment(p payment.Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.payment = p
}

// Validate ensures all required providers are configured.
func (r *Registry) Validate() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.kyc == nil {
		return fmt.Errorf("KYC provider not configured")
	}
	if r.mfa == nil {
		return fmt.Errorf("MFA provider not configured")
	}
	if r.wallet == nil {
		return fmt.Errorf("wallet provider not configured")
	}
	if r.exchange == nil {
		return fmt.Errorf("exchange provider not configured")
	}
	if r.payment == nil {
		return fmt.Errorf("payment provider not configured")
	}
	return nil
}
