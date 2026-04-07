package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"vaxen/api/internal/config"
	"vaxen/api/internal/models"
	"vaxen/api/internal/utils"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

// --- Step 1: Request Access (public form, no password) ---

type RequestAccessInput struct {
	FirstName    string   `json:"firstName" binding:"required"`
	LastName     string   `json:"lastName" binding:"required"`
	Email        string   `json:"email" binding:"required,email"`
	Company      string   `json:"company" binding:"required"`
	Role         string   `json:"role" binding:"required"`
	Country      string   `json:"country"`
	Markets      []string `json:"markets"`
	AnnualVolume string   `json:"annualVolume"`
	UseCase      string   `json:"useCase"`
	Website      string   `json:"website"`
	Notes        string   `json:"notes"`
	Honeypot     string   `json:"honeypot"`
}

// RequestAccess stores a business's application for platform access.
// An admin reviews it. No account or password is created at this stage.
func (s *AuthService) RequestAccess(input RequestAccessInput) (*models.AccessRequest, error) {
	// Check for duplicate pending request
	var existing models.AccessRequest
	if err := s.db.Where("email = ? AND status = ?", input.Email, models.AccessRequestPending).First(&existing).Error; err == nil {
		return nil, errors.New("an access request for this email is already pending")
	}

	marketsJSON, _ := json.Marshal(input.Markets)
	req := models.AccessRequest{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		Company:      input.Company,
		Role:         input.Role,
		Country:      input.Country,
		Markets:      datatypes.JSON(marketsJSON),
		AnnualVolume: input.AnnualVolume,
		UseCase:      input.UseCase,
		Website:      input.Website,
		Notes:        input.Notes,
		Status:       models.AccessRequestPending,
	}
	if err := s.db.Create(&req).Error; err != nil {
		return nil, err
	}

	return &req, nil
}

// ApproveAccessRequest is called by an admin. Generates an invite token
// so the applicant can register with a password.
func (s *AuthService) ApproveAccessRequest(requestID string, adminID string) (*models.AccessRequest, error) {
	var req models.AccessRequest
	if err := s.db.First(&req, "id = ?", requestID).Error; err != nil {
		return nil, errors.New("access request not found")
	}

	if req.Status != models.AccessRequestPending {
		return nil, errors.New("access request is not pending")
	}

	token := generateInviteToken()
	now := time.Now()

	if err := s.db.Model(&req).Updates(map[string]any{
		"status":       models.AccessRequestApproved,
		"reviewed_by":  adminID,
		"reviewed_at":  now,
		"invite_token": token,
	}).Error; err != nil {
		return nil, err
	}

	req.Status = models.AccessRequestApproved
	req.InviteToken = &token

	// TODO: Send invitation email with link containing the invite token

	return &req, nil
}

// RejectAccessRequest is called by an admin.
func (s *AuthService) RejectAccessRequest(requestID string, adminID string, reason string) error {
	var req models.AccessRequest
	if err := s.db.First(&req, "id = ?", requestID).Error; err != nil {
		return errors.New("access request not found")
	}

	if req.Status != models.AccessRequestPending {
		return errors.New("access request is not pending")
	}

	now := time.Now()
	return s.db.Model(&req).Updates(map[string]any{
		"status":      models.AccessRequestRejected,
		"reviewed_by": adminID,
		"reviewed_at": now,
		"review_note": reason,
	}).Error
}

// GetAccessRequests returns access requests, optionally filtered by status.
func (s *AuthService) GetAccessRequests(status string) ([]models.AccessRequest, error) {
	var requests []models.AccessRequest
	q := s.db
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("created_at desc").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// --- Step 2: Register (after admin approval, with invite token) ---

type RegisterInput struct {
	InviteToken string `json:"inviteToken" binding:"required"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
}

type AuthResult struct {
	Token        string      `json:"token"`
	User         models.User `json:"user"`
	RequiresMFA  bool        `json:"requiresMfa"`
	MFAChallenge string      `json:"mfaChallenge,omitempty"`
}

// Register creates the org and first director from an approved access request.
func (s *AuthService) Register(input RegisterInput) (*AuthResult, error) {
	// Validate invite token
	var accessReq models.AccessRequest
	if err := s.db.Where("invite_token = ? AND status = ?", input.InviteToken, models.AccessRequestApproved).First(&accessReq).Error; err != nil {
		return nil, errors.New("invalid or expired invite token")
	}

	if accessReq.Email != input.Email {
		return nil, errors.New("email does not match the invitation")
	}

	// Check if already registered
	var existing models.User
	if err := s.db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	var org models.Organization
	var user models.User

	err = s.db.Transaction(func(tx *gorm.DB) error {
		org = models.Organization{
			Name:    accessReq.Company,
			Country: accessReq.Country,
			// LegalName and RegistrationNumber will be collected during KYB
			LegalName:          accessReq.Company,
			RegistrationNumber: "PENDING",
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

		// Default approval policies
		for _, action := range []models.ApprovalActionType{
			models.ApprovalActionPayout,
			models.ApprovalActionConversion,
			models.ApprovalActionTransfer,
		} {
			if err := tx.Create(&models.ApprovalPolicy{
				OrganizationID:    org.ID,
				ActionType:        action,
				RequiredApprovals: 1,
			}).Error; err != nil {
				return err
			}
		}

		// Consume the invite token
		return tx.Model(&accessReq).Update("invite_token", nil).Error
	})
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, org.ID, user.Email, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpiration)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResult{
		Token:       token,
		User:        user,
		RequiresMFA: true,
	}, nil
}

// --- Login ---

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	MFACode  string `json:"mfaCode,omitempty"`
}

func (s *AuthService) Login(input LoginInput) (*AuthResult, error) {
	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", input.Email, true).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

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

func generateInviteToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
