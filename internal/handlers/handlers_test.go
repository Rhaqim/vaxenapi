package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"vaxen/api/internal/config"
	"vaxen/api/internal/handlers"
	"vaxen/api/internal/middleware"
	"vaxen/api/internal/models"
	"vaxen/api/internal/providers"
	"vaxen/api/internal/providers/exchange"
	"vaxen/api/internal/providers/kyc"
	"vaxen/api/internal/providers/mfa"
	"vaxen/api/internal/providers/payment"
	"vaxen/api/internal/providers/wallet"
	"vaxen/api/internal/services"
	"vaxen/api/internal/testutil"
	"vaxen/api/internal/utils"

	"github.com/shopspring/decimal"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type testEnv struct {
	router *gin.Engine
	h      *handlers.Handler
	svc    *services.Container
	db     *gorm.DB
	cfg    *config.Config
}

func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()
	db := testutil.SetupTestDB(t)
	db.AutoMigrate(&models.Beneficiary{}, &models.Payout{})

	cfg := testutil.TestConfig()
	reg := providers.NewRegistry()
	reg.SetKYC(kyc.NewSumSub(kyc.SumSubConfig{}))
	reg.SetMFA(mfa.NewTwilio(mfa.TwilioConfig{}))
	reg.SetWallet(wallet.NewAWSKMS(wallet.AWSKMSConfig{}))
	reg.SetExchange(exchange.NewInternal())
	reg.SetPayment(payment.NewCircle(payment.CircleConfig{}))

	svc := services.NewContainer(db, cfg, reg)
	h := handlers.NewHandler(svc)

	router := gin.New()
	return &testEnv{router: router, h: h, svc: svc, db: db, cfg: cfg}
}

// registerAndGetToken registers a user and returns the JWT token.
func registerAndGetToken(t *testing.T, env *testEnv) (string, string, string) {
	t.Helper()
	result, err := env.svc.Auth.Register(services.RegisterInput{
		Email:              "handler-test@corp.com",
		Password:           "testpassword1",
		FirstName:          "Handler",
		LastName:           "Test",
		CompanyName:        "Handler Corp",
		LegalName:          "Handler Corp Ltd",
		RegistrationNumber: "HC-001",
		Country:            "US",
	})
	require.NoError(t, err)
	return result.Token, result.User.ID, result.User.OrganizationID
}

func authedRequest(method, url string, body any, token string) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, url, &buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	return resp
}

// --- Auth Handler Tests ---

func TestHandler_Register(t *testing.T) {
	env := setupTestEnv(t)

	env.router.POST("/auth/register", handlers.Register(env.cfg, env.svc))

	body := map[string]any{
		"firstName":          "New",
		"lastName":           "User",
		"email":              "new@corp.com",
		"password":           "password123",
		"companyName":        "New Corp",
		"legalName":          "New Corp Ltd",
		"registrationNumber": "NC-001",
		"country":            "US",
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req := httptest.NewRequest("POST", "/auth/register", &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]any)
	assert.NotEmpty(t, data["token"])
	assert.True(t, data["requiresMfa"].(bool))
}

func TestHandler_Register_Honeypot(t *testing.T) {
	env := setupTestEnv(t)
	env.router.POST("/auth/register", handlers.Register(env.cfg, env.svc))

	body := map[string]any{
		"firstName": "Bot", "lastName": "User", "email": "bot@corp.com",
		"password": "password123", "companyName": "Bot Corp",
		"legalName": "Bot Corp Ltd", "registrationNumber": "BOT-1", "country": "US",
		"honeypot": "gotcha",
	}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req := httptest.NewRequest("POST", "/auth/register", &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Should succeed silently but not create any records
	var count int64
	env.db.Model(&models.User{}).Where("email = ?", "bot@corp.com").Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestHandler_Login(t *testing.T) {
	env := setupTestEnv(t)
	env.router.POST("/auth/login", handlers.Login(env.cfg, env.svc))

	// Register user first
	env.svc.Auth.Register(services.RegisterInput{
		Email: "login-handler@corp.com", Password: "mypassword12",
		FirstName: "Login", LastName: "Handler",
		CompanyName: "LH Corp", LegalName: "LH Corp Ltd",
		RegistrationNumber: "LH-001", Country: "US",
	})

	body := map[string]any{"email": "login-handler@corp.com", "password": "mypassword12"}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req := httptest.NewRequest("POST", "/auth/login", &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	resp := parseResponse(t, w)
	data := resp["data"].(map[string]any)
	assert.NotEmpty(t, data["token"])
}

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	env := setupTestEnv(t)
	env.router.POST("/auth/login", handlers.Login(env.cfg, env.svc))

	body := map[string]any{"email": "nobody@corp.com", "password": "wrongpass1"}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req := httptest.NewRequest("POST", "/auth/login", &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// --- MFA Handler Tests ---

func TestHandler_MFA_Enroll(t *testing.T) {
	env := setupTestEnv(t)
	token, _, _ := registerAndGetToken(t, env)

	protected := env.router.Group("/")
	protected.Use(middleware.AuthMiddleware(env.cfg))
	protected.POST("/mfa/enroll", env.h.EnrollMFA)

	req := authedRequest("POST", "/mfa/enroll", nil, token)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	data := resp["data"].(map[string]any)
	assert.NotEmpty(t, data["qrCodeUrl"])
	assert.NotEmpty(t, data["secret"])
}

func TestHandler_MFA_Confirm(t *testing.T) {
	env := setupTestEnv(t)
	token, _, _ := registerAndGetToken(t, env)

	protected := env.router.Group("/")
	protected.Use(middleware.AuthMiddleware(env.cfg))
	protected.POST("/mfa/enroll", env.h.EnrollMFA)
	protected.POST("/mfa/confirm", env.h.ConfirmMFA)

	// Enroll first
	req := authedRequest("POST", "/mfa/enroll", nil, token)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Confirm
	req = authedRequest("POST", "/mfa/confirm", map[string]string{"code": "123456"}, token)
	w = httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// --- KYB Handler Tests ---

func TestHandler_SubmitKYB(t *testing.T) {
	env := setupTestEnv(t)
	token, _, _ := registerAndGetToken(t, env)

	protected := env.router.Group("/")
	protected.Use(middleware.AuthMiddleware(env.cfg))
	protected.POST("/kyb/submit", env.h.SubmitKYB)

	body := map[string]any{
		"legalName":          "Handler Corp Ltd",
		"registrationNumber": "HC-001",
		"country":            "US",
		"address": map[string]string{
			"street": "123 Main St", "city": "New York", "country": "US",
		},
		"directors": []map[string]string{
			{"firstName": "John", "lastName": "Doe", "email": "john@test.com"},
		},
	}

	req := authedRequest("POST", "/kyb/submit", body, token)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse(t, w)
	data := resp["data"].(map[string]any)
	assert.NotEmpty(t, data["referenceId"])
	assert.Equal(t, "pending", data["status"])
}

func TestHandler_GetKYBStatus(t *testing.T) {
	env := setupTestEnv(t)
	token, _, _ := registerAndGetToken(t, env)

	protected := env.router.Group("/")
	protected.Use(middleware.AuthMiddleware(env.cfg))
	protected.GET("/kyb/status", env.h.GetKYBStatus)

	req := authedRequest("GET", "/kyb/status", nil, token)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- Wallet Handler Tests ---

func TestHandler_GetWallets(t *testing.T) {
	env := setupTestEnv(t)
	token, _, orgID := registerAndGetToken(t, env)

	// Create some wallets
	env.svc.Wallet.CreateDefaultWallets(t.Context(), orgID)

	protected := env.router.Group("/")
	protected.Use(middleware.AuthMiddleware(env.cfg))
	protected.GET("/wallets", env.h.GetWallets)

	req := authedRequest("GET", "/wallets", nil, token)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	data := resp["data"].([]any)
	assert.Len(t, data, 3) // USD, EUR, GBP
}

func TestHandler_GetWeb3Wallets(t *testing.T) {
	env := setupTestEnv(t)
	token, _, orgID := registerAndGetToken(t, env)

	env.svc.Wallet.CreateDefaultWallets(t.Context(), orgID)

	protected := env.router.Group("/")
	protected.Use(middleware.AuthMiddleware(env.cfg))
	protected.GET("/wallets/web3", env.h.GetWeb3Wallets)

	req := authedRequest("GET", "/wallets/web3", nil, token)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	data := resp["data"].([]any)
	assert.Len(t, data, 2) // ethereum, polygon
}

// --- Conversion/Exchange Handler Tests ---

func TestHandler_CreateQuote(t *testing.T) {
	env := setupTestEnv(t)
	token, _, _ := registerAndGetToken(t, env)

	protected := env.router.Group("/")
	protected.Use(middleware.AuthMiddleware(env.cfg))
	protected.POST("/quotes", env.h.CreateQuote)

	body := map[string]any{
		"fromCurrency": "USD", "toCurrency": "EUR",
		"amount": "1000", "side": "sell",
	}
	req := authedRequest("POST", "/quotes", body, token)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse(t, w)
	data := resp["data"].(map[string]any)
	assert.NotEmpty(t, data["QuoteID"])
}

func TestHandler_GetSupportedPairs(t *testing.T) {
	env := setupTestEnv(t)
	token, _, _ := registerAndGetToken(t, env)

	protected := env.router.Group("/")
	protected.Use(middleware.AuthMiddleware(env.cfg))
	protected.GET("/conversions/pairs", env.h.GetSupportedPairs)

	req := authedRequest("GET", "/conversions/pairs", nil, token)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	data := resp["data"].([]any)
	assert.NotEmpty(t, data)
}

func TestHandler_GetRate(t *testing.T) {
	env := setupTestEnv(t)
	token, _, _ := registerAndGetToken(t, env)

	protected := env.router.Group("/")
	protected.Use(middleware.AuthMiddleware(env.cfg))
	protected.GET("/conversions/rate", env.h.GetRate)

	req := authedRequest("GET", "/conversions/rate?from=USD&to=EUR", nil, token)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- Payout Handler Tests ---

func TestHandler_CreatePayout(t *testing.T) {
	env := setupTestEnv(t)
	token, _, orgID := registerAndGetToken(t, env)

	b := models.Beneficiary{OrganizationID: orgID, Name: "Test Ben"}
	env.db.Create(&b)

	protected := env.router.Group("/")
	protected.Use(middleware.AuthMiddleware(env.cfg))
	protected.POST("/payouts", env.h.CreatePayout)

	body := map[string]any{
		"beneficiaryId": b.ID, "amount": "500.00",
		"currency": "USD", "reference": "PAY-TEST",
	}
	req := authedRequest("POST", "/payouts", body, token)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// --- Approval Handler Tests ---

func TestHandler_ApprovalWorkflow(t *testing.T) {
	env := setupTestEnv(t)
	token, userID, orgID := registerAndGetToken(t, env)

	// Set policy to 2 approvals
	env.svc.Approval.UpdatePolicy(orgID, models.ApprovalActionPayout, 2)

	// Create a beneficiary and initiate payment (which will create an approval request)
	b := models.Beneficiary{OrganizationID: orgID, Name: "Approval Ben"}
	env.db.Create(&b)

	env.svc.Payment.InitiatePayment(t.Context(), services.SendPaymentInput{
		OrganizationID: orgID, InitiatedByID: userID,
		BeneficiaryID: b.ID, Amount: decimal.RequireFromString("1000"),
		Currency: "USD", Reference: "APR-001",
	})

	protected := env.router.Group("/")
	protected.Use(middleware.AuthMiddleware(env.cfg))
	protected.GET("/approvals/pending", env.h.GetPendingApprovals)

	req := authedRequest("GET", "/approvals/pending", nil, token)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse(t, w)
	data := resp["data"].([]any)
	assert.Len(t, data, 1)
}

// --- Admin Handler Tests ---

func TestHandler_Admin_GetSettings(t *testing.T) {
	env := setupTestEnv(t)
	token, _, _ := registerAndGetToken(t, env)

	// Make user an admin
	env.db.Model(&models.User{}).Where("email = ?", "handler-test@corp.com").Update("role", "admin")
	// Re-generate token with admin role
	token, _ = utils.GenerateToken(
		"", "", "handler-test@corp.com", "admin",
		env.cfg.JWTSecret, env.cfg.JWTExpiration,
	)

	// Seed a setting
	env.svc.Admin.UpsertSetting("test_key", "test_value", "general", "admin")

	admin := env.router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(env.cfg))
	admin.Use(middleware.AdminMiddleware())
	admin.GET("/settings", env.h.GetPlatformSettings)

	req := authedRequest("GET", "/admin/settings", nil, token)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Admin_FeatureFlags(t *testing.T) {
	env := setupTestEnv(t)

	token, _ := utils.GenerateToken("admin-1", "org-1", "admin@test.com", "admin",
		env.cfg.JWTSecret, env.cfg.JWTExpiration)

	admin := env.router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(env.cfg))
	admin.Use(middleware.AdminMiddleware())
	admin.PUT("/features", env.h.SetFeatureFlag)
	admin.GET("/features", env.h.GetFeatureFlags)

	// Set a flag
	body := map[string]any{"name": "dark_mode", "enabled": true, "description": "Dark mode"}
	req := authedRequest("PUT", "/admin/features", body, token)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Get flags
	req = authedRequest("GET", "/admin/features", nil, token)
	w = httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	resp := parseResponse(t, w)
	data := resp["data"].([]any)
	assert.Len(t, data, 1)
}

// --- Health Check ---

func TestHandler_HealthCheck(t *testing.T) {
	router := gin.New()
	router.GET("/health", handlers.HealthCheck)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
