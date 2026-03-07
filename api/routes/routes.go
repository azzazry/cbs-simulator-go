package routes

import (
	"time"

	"cbs-simulator/api/handlers"
	"cbs-simulator/api/middleware"
	"cbs-simulator/config"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes and middleware
func SetupRoutes(router *gin.Engine) {
	// Global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())

	// Rate limiting (global)
	rateLimit := config.AppConfig.RateLimitPerMinute
	if rateLimit <= 0 {
		rateLimit = 60
	}
	router.Use(middleware.RateLimiterMiddleware(rateLimit, time.Minute))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "CBS Simulator"})
	})

	// API v1
	v1 := router.Group("/api/v1")

	// === Public routes (no auth required) ===
	auth := v1.Group("/auth")
	{
		auth.POST("/login", handlers.Login)
		auth.POST("/register", handlers.Register)

		// Self-service unlock flow (no auth needed - user is locked out)
		auth.POST("/otp/request", handlers.RequestOTP)
		auth.POST("/otp/verify", handlers.VerifyOTPHandler)
		auth.POST("/ekyc/verify", handlers.VerifyEKYC)
		auth.POST("/unlock", handlers.UnlockAccount)
		auth.POST("/reset-pin", handlers.ResetPINHandler)
	}

	// === Protected routes (JWT auth required) ===
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.AuditMiddleware())
	{
		// Auth (authenticated)
		protectedAuth := protected.Group("/auth")
		{
			protectedAuth.POST("/logout", handlers.Logout)
			protectedAuth.POST("/refresh", handlers.RefreshToken)
			protectedAuth.POST("/change-pin", handlers.ChangePIN)
			protectedAuth.GET("/profile", handlers.GetProfile)
		}

		// Customer management
		customer := protected.Group("/customers")
		{
			customer.GET("/:cif", handlers.GetAccountsByCIF)
			customer.GET("/:cif/accounts", handlers.GetAccountsByCIF)
		}

		// Account operations
		accounts := protected.Group("/accounts")
		{
			accounts.GET("/:account_number", handlers.GetAccountBalance)
			accounts.GET("/:account_number/transactions", handlers.GetAccountStatement)
			accounts.GET("/:account_number/balance", handlers.GetAccountBalance)
		}

		// Transfer operations
		transfers := protected.Group("/transfers")
		{
			transfers.POST("/intra", handlers.IntraBankTransfer)
			transfers.POST("/inter", handlers.InterBankTransfer)
			transfers.GET("/fees", handlers.GetTransferFees)
		}

		// Bill Payment
		bills := protected.Group("/bills")
		{
			bills.POST("/pay", handlers.PayBill)
			bills.GET("/history", handlers.GetAllBills)
		}

		// Card operations
		cards := protected.Group("/cards")
		{
			cards.GET("/:cif", handlers.GetCardsByCIF)
			cards.POST("/block", handlers.BlockCard)
			cards.POST("/unblock", handlers.UnblockCard)
		}

		// Loan operations
		loans := protected.Group("/loans")
		{
			loans.GET("/:cif", handlers.GetLoansByCIF)
			loans.GET("/detail/:loan_number", handlers.GetLoanDetails)
		}

		// Deposit operations
		deposits := protected.Group("/deposits")
		{
			deposits.GET("/:cif", handlers.GetDepositsByCIF)
			deposits.GET("/detail/:deposit_number", handlers.GetDepositDetails)
		}

		// Notification operations
		notifications := protected.Group("/notifications")
		{
			notifications.GET("/:cif", handlers.GetNotifications)
			notifications.POST("/read", handlers.MarkNotificationAsRead)
		}

		// Payment operations (QRIS, VA, E-Wallet, E-Money)
		payments := protected.Group("/payments")
		{
			payments.POST("/qris", handlers.ProcessQRISPayment)
			payments.POST("/va", handlers.ProcessVAPayment)
			payments.POST("/ewallet/topup", handlers.ProcessEWalletTopup)
			payments.POST("/emoney/topup", handlers.ProcessEMoneyTopup)
		}

		// === Phase 2: Core Banking Routes ===

		// General Ledger
		gl := protected.Group("/gl")
		{
			gl.GET("/chart-of-accounts", handlers.GetChartOfAccounts)
			gl.GET("/journal-entries", handlers.GetJournalEntries)
			gl.GET("/journal-entries/:id", handlers.GetJournalDetail)
			gl.GET("/trial-balance", handlers.GetTrialBalance)
			gl.GET("/account-balance/:code", handlers.GetGLAccountBalance)
		}

		// CIF Enhancement
		protected.GET("/customers/:cif/overview", handlers.GetCustomerOverview)
		protected.PUT("/customers/:cif/extended", handlers.UpdateCustomerExtended)
		protected.GET("/customers/search", handlers.SearchCustomers)

		// Interest
		interest := protected.Group("/interest")
		{
			interest.GET("/rates", handlers.GetInterestRates)
			interest.POST("/calculate", handlers.SimulateInterest)
		}

		// Standing Instructions
		si := protected.Group("/standing-instructions")
		{
			si.POST("", handlers.CreateStandingInstruction)
			si.GET("/:cif", handlers.GetStandingInstructions)
			si.PUT("/:si/pause", handlers.PauseStandingInstruction)
			si.DELETE("/:si", handlers.CancelStandingInstruction)
			si.GET("/:si/history", handlers.GetSIHistory)
		}

		// Account Management
		protected.POST("/accounts/open", handlers.OpenAccountHandler)
		protected.POST("/accounts/:account_number/close", handlers.CloseAccountHandler)
		protected.GET("/accounts/dormant", handlers.GetDormantAccounts)
		protected.POST("/accounts/:account_number/reactivate", handlers.ReactivateAccountHandler)

		// === Admin routes (require admin role) ===
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole("admin", "supervisor"))
		{
			admin.GET("/audit-logs", handlers.GetAuditLogs)
			admin.GET("/transaction-limits", handlers.GetTransactionLimits)
			admin.PUT("/transaction-limits", handlers.UpdateTransactionLimit)
			admin.GET("/roles", handlers.GetUserRolesHandler)
			admin.POST("/roles/assign", handlers.AssignRoleHandler)
			admin.POST("/unlock-account", handlers.AdminUnlockAccount)

			// Phase 2 Admin: EOD Processing
			admin.POST("/eod/run", handlers.RunEOD)
			admin.GET("/eod/status/:date", handlers.GetEODStatus)
			admin.GET("/eod/history", handlers.GetEODHistory)
		}
	}
}
