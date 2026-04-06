package services

import (
	"context"
	"errors"
	"time"

	"vaxen/api/internal/models"
	"vaxen/api/internal/providers"
	"vaxen/api/internal/providers/mfa"

	"gorm.io/gorm"
)

type MFAService struct {
	db  *gorm.DB
	reg *providers.Registry
}

func NewMFAService(db *gorm.DB, reg *providers.Registry) *MFAService {
	return &MFAService{db: db, reg: reg}
}

// EnrollMFA starts MFA enrollment for a user. Required for directors during signup.
func (s *MFAService) EnrollMFA(ctx context.Context, userID string, email string) (*mfa.EnrollResult, error) {
	provider, err := s.reg.RequireMFA()
	if err != nil {
		return nil, err
	}

	var existing models.MFAEnrollment
	if err := s.db.Where("user_id = ?", userID).First(&existing).Error; err == nil {
		if existing.IsConfirmed {
			return nil, errors.New("MFA already enrolled and confirmed")
		}
		s.db.Delete(&existing)
	}

	result, err := provider.Enroll(ctx, mfa.EnrollRequest{
		UserID: userID,
		Email:  email,
	})
	if err != nil {
		return nil, err
	}

	enrollment := models.MFAEnrollment{
		UserID:     userID,
		Provider:   provider.Name(),
		ProviderID: result.ProviderID,
		Secret:     result.Secret,
		BackupCode: result.BackupCode,
	}
	if err := s.db.Create(&enrollment).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// ConfirmMFA confirms MFA enrollment with the first valid code.
func (s *MFAService) ConfirmMFA(ctx context.Context, userID string, code string) error {
	provider, err := s.reg.RequireMFA()
	if err != nil {
		return err
	}

	var enrollment models.MFAEnrollment
	if err := s.db.Where("user_id = ?", userID).First(&enrollment).Error; err != nil {
		return errors.New("no MFA enrollment found")
	}

	if enrollment.IsConfirmed {
		return errors.New("MFA already confirmed")
	}

	if err := provider.ConfirmEnrollment(ctx, userID, code); err != nil {
		return err
	}

	now := time.Now()
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&enrollment).Updates(map[string]any{
			"is_confirmed": true,
			"confirmed_at": now,
		}).Error; err != nil {
			return err
		}
		return tx.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]any{
			"mfa_enabled":  true,
			"mfa_provider": provider.Name(),
		}).Error
	})
}

// SendChallenge sends an MFA challenge to the user (for login verification).
func (s *MFAService) SendChallenge(ctx context.Context, userID string, channel string) (*mfa.ChallengeResult, error) {
	provider, err := s.reg.RequireMFA()
	if err != nil {
		return nil, err
	}

	result, err := provider.SendChallenge(ctx, mfa.ChallengeRequest{
		UserID:  userID,
		Channel: channel,
	})
	if err != nil {
		return nil, err
	}

	challenge := models.MFAChallenge{
		UserID:      userID,
		ChallengeID: result.ChallengeID,
		Channel:     channel,
		ExpiresAt:   time.Now().Add(time.Duration(result.ExpiresInS) * time.Second),
	}
	s.db.Create(&challenge)

	return result, nil
}

// VerifyCode verifies an MFA code provided by the user.
func (s *MFAService) VerifyCode(ctx context.Context, userID string, code string) error {
	provider, err := s.reg.RequireMFA()
	if err != nil {
		return err
	}

	result, err := provider.VerifyCode(ctx, mfa.VerifyRequest{
		UserID: userID,
		Code:   code,
	})
	if err != nil {
		return err
	}

	if !result.Valid {
		return errors.New("invalid MFA code")
	}

	return nil
}

// UnenrollMFA removes MFA for a user (admin action).
func (s *MFAService) UnenrollMFA(ctx context.Context, userID string) error {
	provider, err := s.reg.RequireMFA()
	if err != nil {
		return err
	}

	if err := provider.Unenroll(ctx, userID); err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&models.MFAEnrollment{}).Error; err != nil {
			return err
		}
		return tx.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]any{
			"mfa_enabled":  false,
			"mfa_provider": "",
			"mfa_secret":   nil,
		}).Error
	})
}
