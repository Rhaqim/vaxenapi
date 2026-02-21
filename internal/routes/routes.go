package routes

import (
	"vaxen/api/internal/config"
	"vaxen/api/internal/handlers"
	"vaxen/api/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
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
			public.POST("/auth/login", handlers.Login(cfg))
			public.POST("/auth/register", handlers.Register(cfg))
			public.POST("/auth/refresh", handlers.RefreshToken(cfg))
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
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
				kyb.POST("/submit", handlers.SubmitKYB)
				kyb.GET("/status", handlers.GetKYBStatus)
			}

			// Wallets
			wallets := protected.Group("/wallets")
			{
				wallets.GET("/", handlers.GetWallets)
				wallets.GET("/:id", handlers.GetWallet)
				wallets.POST("/", handlers.CreateWallet)
				wallets.GET("/:id/balance", handlers.GetWalletBalance)
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
				quotes.POST("/", handlers.CreateQuote)
				quotes.GET("/:id", handlers.GetQuote)
			}

			// Conversions
			conversions := protected.Group("/conversions")
			{
				conversions.POST("/", handlers.CreateConversion)
				conversions.GET("/", handlers.GetConversions)
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

			// Payouts
			payouts := protected.Group("/payouts")
			{
				payouts.GET("/", handlers.GetPayouts)
				payouts.GET("/:id", handlers.GetPayout)
				payouts.POST("/", handlers.CreatePayout)
			}

			// Crypto
			crypto := protected.Group("/crypto")
			{
				crypto.GET("/addresses", handlers.GetCryptoAddresses)
				crypto.POST("/addresses", handlers.CreateCryptoAddress)
				crypto.POST("/withdraw", handlers.CryptoWithdraw)
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

		// Admin routes (require admin role)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg))
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/users", handlers.GetAllUsers)
			admin.GET("/users/:id", handlers.GetUser)
			admin.PUT("/users/:id", handlers.UpdateUser)
			admin.DELETE("/users/:id", handlers.DeleteUser)
			admin.GET("/organizations", handlers.GetAllOrganizations)
			admin.POST("/organizations/:id/approve", handlers.ApproveOrganization)
			admin.POST("/organizations/:id/reject", handlers.RejectOrganization)
		}
	}
}
