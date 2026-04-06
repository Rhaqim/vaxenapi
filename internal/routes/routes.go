package routes

import (
	"vaxen/api/internal/config"
	"vaxen/api/internal/handlers"
	"vaxen/api/internal/middleware"
	"vaxen/api/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, cfg *config.Config, svc *services.Container) {
	h := handlers.NewHandler(svc)

	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		public := v1.Group("/")
		{
			public.POST("/auth/login", handlers.Login(cfg, svc))
			public.POST("/auth/register", handlers.Register(cfg, svc))
			public.POST("/auth/refresh", handlers.RefreshToken(cfg))
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// MFA (must be accessible before MFA is set up)
			mfa := protected.Group("/mfa")
			{
				mfa.POST("/enroll", h.EnrollMFA)
				mfa.POST("/confirm", h.ConfirmMFA)
				mfa.POST("/disable", h.DisableMFA)
			}

			// Organizations
			organizations := protected.Group("/organizations")
			{
				organizations.GET("/", handlers.GetOrganizations)
				organizations.GET("/:id", handlers.GetOrganization)
				organizations.POST("/", handlers.CreateOrganization)
				organizations.PUT("/:id", handlers.UpdateOrganization)
			}

			// KYB (Know Your Business)
			kyb := protected.Group("/kyb")
			{
				kyb.POST("/submit", h.SubmitKYB)
				kyb.GET("/status", h.GetKYBStatus)
			}

			// Wallets
			wallets := protected.Group("/wallets")
			{
				wallets.GET("/", h.GetWallets)
				wallets.GET("/web3", h.GetWeb3Wallets)
				wallets.GET("/:id", h.GetWallet)
				wallets.POST("/", h.CreateWallet)
				wallets.GET("/:id/balance", h.GetWalletBalance)
			}

			// Accounts
			accounts := protected.Group("/accounts")
			{
				accounts.GET("/", handlers.GetAccounts)
				accounts.GET("/:id", handlers.GetAccount)
				accounts.POST("/", handlers.CreateAccount)
			}

			// Quotes
			quotes := protected.Group("/quotes")
			{
				quotes.POST("/", h.CreateQuote)
				quotes.GET("/:id", handlers.GetQuote)
			}

			// Conversions / Swaps
			conversions := protected.Group("/conversions")
			{
				conversions.POST("/", h.CreateConversion)
				conversions.GET("/", handlers.GetConversions)
				conversions.GET("/pairs", h.GetSupportedPairs)
				conversions.GET("/rate", h.GetRate)
				conversions.GET("/:id", handlers.GetConversion)
			}

			// Orders
			orders := protected.Group("/orders")
			{
				orders.GET("/", handlers.GetOrders)
				orders.GET("/open", handlers.GetOpenOrders)
				orders.GET("/:id", handlers.GetOrder)
				orders.POST("/", handlers.CreateOrder)
				orders.PUT("/:id/cancel", handlers.CancelOrder)
			}

			// Beneficiaries
			beneficiaries := protected.Group("/beneficiaries")
			{
				beneficiaries.GET("/", handlers.GetBeneficiaries)
				beneficiaries.GET("/:id", handlers.GetBeneficiary)
				beneficiaries.POST("/", handlers.CreateBeneficiary)
				beneficiaries.PUT("/:id", handlers.UpdateBeneficiary)
				beneficiaries.DELETE("/:id", handlers.DeleteBeneficiary)
			}

			// Payouts / Send Money
			payouts := protected.Group("/payouts")
			{
				payouts.GET("/", h.GetPayouts)
				payouts.GET("/fees", h.GetPayoutFees)
				payouts.GET("/:id", h.GetPayout)
				payouts.POST("/", h.CreatePayout)
			}

			// Crypto
			crypto := protected.Group("/crypto")
			{
				crypto.GET("/addresses", handlers.GetCryptoAddresses)
				crypto.POST("/addresses", handlers.CreateCryptoAddress)
				crypto.POST("/withdraw", handlers.CryptoWithdraw)
			}

			// Approvals (multi-director workflow)
			approvals := protected.Group("/approvals")
			{
				approvals.GET("/pending", h.GetPendingApprovals)
				approvals.GET("/:id", h.GetApprovalRequest)
				approvals.POST("/:id/vote", h.VoteOnApproval)
				approvals.PUT("/policy", h.UpdateApprovalPolicy)
			}

			// Statements
			statements := protected.Group("/statements")
			{
				statements.GET("/", handlers.GetStatements)
				statements.GET("/:id", handlers.GetStatement)
			}

			// Reports
			reports := protected.Group("/reports")
			{
				reports.GET("/fx-pnl", handlers.GetFxPnlReport)
				reports.GET("/transactions", handlers.GetTransactionReport)
				reports.GET("/balances", handlers.GetBalanceReport)
			}

			// Audit
			audit := protected.Group("/audit")
			{
				audit.GET("/logs", handlers.GetAuditLogs)
			}

			// Providers
			providers := protected.Group("/providers")
			{
				providers.GET("/", handlers.GetProviders)
				providers.GET("/:id", handlers.GetProvider)
			}
		}

		// Webhooks (public but with provider authentication)
		webhooks := v1.Group("/webhooks")
		{
			webhooks.POST("/provider/:provider", handlers.ProcessWebhook)
		}

		// Platform Admin routes (require admin role)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg))
		admin.Use(middleware.AdminMiddleware())
		{
			// User management
			admin.GET("/users", h.GetAllUsers)
			admin.GET("/users/:id", handlers.GetUser)
			admin.PUT("/users/:id", handlers.UpdateUser)
			admin.DELETE("/users/:id", handlers.DeleteUser)

			// Organization management
			admin.GET("/organizations", h.GetAllOrganizations)
			admin.POST("/organizations/:id/approve", h.ApproveOrganization)
			admin.POST("/organizations/:id/reject", h.RejectOrganization)

			// Platform settings
			admin.GET("/settings", h.GetPlatformSettings)
			admin.PUT("/settings", h.UpdatePlatformSetting)

			// Feature flags
			admin.GET("/features", h.GetFeatureFlags)
			admin.PUT("/features", h.SetFeatureFlag)

			// Exchange rates
			admin.GET("/exchange-rates", h.GetExchangeRates)
			admin.PUT("/exchange-rates", h.UpsertExchangeRate)
		}
	}
}
