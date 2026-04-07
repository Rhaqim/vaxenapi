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
	h := handlers.NewHandler(svc, cfg)

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
			public.POST("/auth/request-access", h.RequestAccess)
			public.POST("/auth/register", h.Register)
			public.POST("/auth/login", h.Login)
			public.POST("/auth/refresh", h.RefreshToken)
			public.POST("/auth/logout", h.Logout)
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// MFA
			mfa := protected.Group("/mfa")
			{
				mfa.POST("/enroll", h.EnrollMFA)
				mfa.POST("/confirm", h.ConfirmMFA)
				mfa.POST("/disable", h.DisableMFA)
			}

			// Organizations
			organizations := protected.Group("/organizations")
			{
				organizations.GET("/", h.GetOrganizations)
				organizations.GET("/:id", h.GetOrganization)
				organizations.POST("/", h.CreateOrganization)
				organizations.PUT("/:id", h.UpdateOrganization)
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
				accounts.GET("/", h.GetAccounts)
				accounts.GET("/:id", h.GetAccount)
				accounts.POST("/", h.CreateAccount)
			}

			// Quotes
			quotes := protected.Group("/quotes")
			{
				quotes.POST("/", h.CreateQuote)
				quotes.GET("/:id", h.GetQuote)
			}

			// Conversions / Swaps
			conversions := protected.Group("/conversions")
			{
				conversions.POST("/", h.CreateConversion)
				conversions.GET("/", h.GetConversions)
				conversions.GET("/pairs", h.GetSupportedPairs)
				conversions.GET("/rate", h.GetRate)
				conversions.GET("/:id", h.GetConversion)
			}

			// Orders
			orders := protected.Group("/orders")
			{
				orders.GET("/", h.GetOrders)
				orders.GET("/open", h.GetOpenOrders)
				orders.GET("/:id", h.GetOrder)
				orders.POST("/", h.CreateOrder)
				orders.PUT("/:id/cancel", h.CancelOrder)
			}

			// Beneficiaries
			beneficiaries := protected.Group("/beneficiaries")
			{
				beneficiaries.GET("/", h.GetBeneficiaries)
				beneficiaries.GET("/:id", h.GetBeneficiary)
				beneficiaries.POST("/", h.CreateBeneficiary)
				beneficiaries.PUT("/:id", h.UpdateBeneficiary)
				beneficiaries.DELETE("/:id", h.DeleteBeneficiary)
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
				crypto.GET("/addresses", h.GetCryptoAddresses)
				crypto.POST("/addresses", h.CreateCryptoAddress)
				crypto.POST("/withdraw", h.CryptoWithdraw)
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
				statements.GET("/", h.GetStatements)
				statements.GET("/:id", h.GetStatement)
			}

			// Reports
			reports := protected.Group("/reports")
			{
				reports.GET("/fx-pnl", h.GetFxPnlReport)
				reports.GET("/transactions", h.GetTransactionReport)
				reports.GET("/balances", h.GetBalanceReport)
			}

			// Audit
			audit := protected.Group("/audit")
			{
				audit.GET("/logs", h.GetAuditLogs)
			}

			// Providers
			providers := protected.Group("/providers")
			{
				providers.GET("/", h.GetProviders)
				providers.GET("/:id", h.GetProvider)
			}
		}

		// Webhooks (public but with provider authentication)
		webhooks := v1.Group("/webhooks")
		{
			webhooks.POST("/provider/:provider", h.ProcessWebhook)
		}

		// Platform Admin routes (require admin role)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg))
		admin.Use(middleware.AdminMiddleware())
		{
			// User management
			admin.GET("/users", h.GetAllUsers)
			admin.GET("/users/:id", h.GetUser)
			admin.PUT("/users/:id", h.UpdateUser)
			admin.DELETE("/users/:id", h.DeleteUser)

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

			// Access requests
			admin.GET("/access-requests", h.GetAccessRequests)
			admin.POST("/access-requests/:id/approve", h.ApproveAccessRequest)
			admin.POST("/access-requests/:id/reject", h.RejectAccessRequest)
		}
	}
}
