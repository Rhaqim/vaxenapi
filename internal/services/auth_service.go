package services

import (
	"errors"
	"time"

	"vaxen/api/internal/config"
	"vaxen/api/internal/models"
	"vaxen/api/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

type RegisterInput struct {
	Email              string
	Password           string
	FirstName          string
	LastName           string
	CompanyName        string
	LegalName          string
	RegistrationNumber string
	Country            string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	Token        string      `json:"token"`
	User         models.User `json:"user"`
	RequiresMFA  bool        `json:"requiresMfa"`
	MFAChallenge string      `json:"mfaChallenge,omitempty"`
}

// Register creates a new organization and its first director (owner).
func (s *AuthService) Register(input RegisterInput) (*AuthResult, error) {
	// Check if email already exists
	var existing models.User
	if err := s.db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create org and user in a transaction
	var org models.Organization
	var user models.User

	err = s.db.Transaction(func(tx *gorm.DB) error {
		org = models.Organization{
			Name:               input.CompanyName,
			LegalName:          input.LegalName,
			RegistrationNumber: input.RegistrationNumber,
			Country:            input.Country,
			KYBStatus:          models.KybStatusPending,
		}
		if err := tx.Create(&org).Error; err != nil {
			return err
		}

		user = models.User{
			Email:          input.Email,
			FirstName:      input.FirstName,
			LastName:       input.LastName,
			PasswordHash:   string(hash),
			Role:           models.UserRoleOwner,
			OrganizationID: org.ID,
			IsDirector:     true,
			IsActive:       true,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// Create default approval policies for the org
		policies := []models.ApprovalPolicy{
			{OrganizationID: org.ID, ActionType: models.ApprovalActionPayout, RequiredApprovals: 1},
			{OrganizationID: org.ID, ActionType: models.ApprovalActionConversion, RequiredApprovals: 1},
			{OrganizationID: org.ID, ActionType: models.ApprovalActionTransfer, RequiredApprovals: 1},
		}
		for _, p := range policies {
			if err := tx.Create(&p).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Directors must set up MFA before they can operate — return token but flag MFA required
	token, err := utils.GenerateToken(user.ID, org.ID, user.Email, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpiration)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResult{
		Token:       token,
		User:        user,
		RequiresMFA: true, // Directors must enroll MFA
	}, nil
}

// Login authenticates a user. If MFA is enabled, returns a challenge instead of a full token.
func (s *AuthService) Login(input LoginInput) (*AuthResult, error) {
	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", input.Email, true).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// If MFA is enabled, don't issue a full token yet — the handler must verify MFA first
	if user.MFAEnabled {
		return &AuthResult{
			User:        user,
			RequiresMFA: true,
		}, nil
	}

	now := time.Now()
	s.db.Model(&user).Update("last_login_at", now)

	token, err := utils.GenerateToken(user.ID, user.OrganizationID, user.Email, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpiration)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResult{
		Token: token,
		User:  user,
	}, nil
}

// CompleteLogin issues a token after MFA verification succeeds.
func (s *AuthService) CompleteLogin(userID string) (*AuthResult, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	now := time.Now()
	s.db.Model(&user).Update("last_login_at", now)

	token, err := utils.GenerateToken(user.ID, user.OrganizationID, user.Email, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpiration)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResult{
		Token: token,
		User:  user,
	}, nil
}
