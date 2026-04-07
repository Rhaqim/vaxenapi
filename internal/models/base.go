package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GenerateUUID sets the ID field to a new UUID if it's empty.
// Called as a BeforeCreate hook on all models with UUID primary keys.
func GenerateUUID(id *string) {
	if *id == "" {
		*id = uuid.New().String()
	}
}

// BeforeCreate hooks for all models with UUID primary keys.
// This replaces PostgreSQL-specific default values to keep
// the codebase database-agnostic and testable with SQLite.

func (o *Organization) BeforeCreate(tx *gorm.DB) error    { GenerateUUID(&o.ID); return nil }
func (u *User) BeforeCreate(tx *gorm.DB) error             { GenerateUUID(&u.ID); return nil }
func (w *Wallet) BeforeCreate(tx *gorm.DB) error            { GenerateUUID(&w.ID); return nil }
func (wt *WalletTransaction) BeforeCreate(tx *gorm.DB) error { GenerateUUID(&wt.ID); return nil }
func (a *AccountNumber) BeforeCreate(tx *gorm.DB) error     { GenerateUUID(&a.ID); return nil }
func (b *Beneficiary) BeforeCreate(tx *gorm.DB) error       { GenerateUUID(&b.ID); return nil }
func (d *Deposit) BeforeCreate(tx *gorm.DB) error           { GenerateUUID(&d.ID); return nil }
func (p *Payout) BeforeCreate(tx *gorm.DB) error            { GenerateUUID(&p.ID); return nil }
func (pb *PayoutBatch) BeforeCreate(tx *gorm.DB) error      { GenerateUUID(&pb.ID); return nil }
func (co *ConversionOrder) BeforeCreate(tx *gorm.DB) error  { GenerateUUID(&co.ID); return nil }
func (lo *LimitOrder) BeforeCreate(tx *gorm.DB) error       { GenerateUUID(&lo.ID); return nil }
func (j *Journal) BeforeCreate(tx *gorm.DB) error           { GenerateUUID(&j.ID); return nil }
func (le *LedgerEntry) BeforeCreate(tx *gorm.DB) error      { GenerateUUID(&le.ID); return nil }
func (rr *ReconciliationRun) BeforeCreate(tx *gorm.DB) error { GenerateUUID(&rr.ID); return nil }
func (sf *StatementFile) BeforeCreate(tx *gorm.DB) error    { GenerateUUID(&sf.ID); return nil }
func (cc *ComplianceCase) BeforeCreate(tx *gorm.DB) error   { GenerateUUID(&cc.ID); return nil }
func (sr *ScreeningResult) BeforeCreate(tx *gorm.DB) error  { GenerateUUID(&sr.ID); return nil }
func (al *AuditLog) BeforeCreate(tx *gorm.DB) error         { GenerateUUID(&al.ID); return nil }
func (we *WebhookEvent) BeforeCreate(tx *gorm.DB) error     { GenerateUUID(&we.ID); return nil }
func (q *Quote) BeforeCreate(tx *gorm.DB) error             { GenerateUUID(&q.ID); return nil }

// New models
func (ap *ApprovalPolicy) BeforeCreate(tx *gorm.DB) error   { GenerateUUID(&ap.ID); return nil }
func (ar *ApprovalRequest) BeforeCreate(tx *gorm.DB) error  { GenerateUUID(&ar.ID); return nil }
func (av *ApprovalVote) BeforeCreate(tx *gorm.DB) error     { GenerateUUID(&av.ID); return nil }
func (ps *PlatformSetting) BeforeCreate(tx *gorm.DB) error  { GenerateUUID(&ps.ID); return nil }
func (ff *FeatureFlag) BeforeCreate(tx *gorm.DB) error      { GenerateUUID(&ff.ID); return nil }
func (er *ExchangeRate) BeforeCreate(tx *gorm.DB) error     { GenerateUUID(&er.ID); return nil }
func (me *MFAEnrollment) BeforeCreate(tx *gorm.DB) error    { GenerateUUID(&me.ID); return nil }
func (mc *MFAChallenge) BeforeCreate(tx *gorm.DB) error     { GenerateUUID(&mc.ID); return nil }
func (ww *Web3Wallet) BeforeCreate(tx *gorm.DB) error       { GenerateUUID(&ww.ID); return nil }
func (ar2 *AccessRequest) BeforeCreate(tx *gorm.DB) error   { GenerateUUID(&ar2.ID); return nil }
