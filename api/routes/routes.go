package routes

import (
	"cbs-simulator/api/handlers"
	"cbs-simulator/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine) {
	// Apply middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "CBS Simulator",
			"version": "1.0.0",
		})
	})
	
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/register", handlers.Register)
			auth.POST("/change-pin", handlers.ChangePIN)
		}
		
		// Customer routes
		customer := v1.Group("/customers")
		{
			customer.GET("/:cif", handlers.GetProfile)
			customer.GET("/:cif/accounts", handlers.GetAccountsByCIF)
			customer.GET("/:cif/cards", handlers.GetCardsByCIF)
			customer.GET("/:cif/loans", handlers.GetLoansByCIF)
			customer.GET("/:cif/deposits", handlers.GetDepositsByCIF)
		}
		
		// Account routes
		accounts := v1.Group("/accounts")
		{
			accounts.GET("/:account_number", handlers.GetAccountBalance)
			accounts.GET("/:account_number/statement", handlers.GetAccountStatement)
		}
		
		// Transfer routes
		transfers := v1.Group("/transfers")
		{
			transfers.POST("/intra", handlers.IntraBankTransfer)
			transfers.POST("/inter", handlers.InterBankTransfer)
			transfers.GET("/:transaction_id", handlers.GetTransaction)
		}
		
		// Bill payment routes
		bills := v1.Group("/bills")
		{
			bills.GET("/billers", handlers.GetBillerList)
			bills.GET("/inquiry", handlers.InquiryBill)
			bills.POST("/pay", handlers.PayBill)
			bills.GET("/all", handlers.GetAllBills) // For testing
		}
		
		// Card routes
		cards := v1.Group("/cards")
		{
			cards.GET("/:card_number", handlers.GetCardDetails)
			cards.POST("/block", handlers.BlockCard)
			cards.POST("/unblock", handlers.UnblockCard)
		}
		
		// Loan routes
		loans := v1.Group("/loans")
		{
			loans.GET("/:loan_number", handlers.GetLoanDetails)
		}
		
		// Deposit routes
		deposits := v1.Group("/deposits")
		{
			deposits.GET("/:deposit_number", handlers.GetDepositDetails)
		}
		
		// Notification routes
		notifications := v1.Group("/notifications")
		{
			notifications.GET("/:cif", handlers.GetNotifications)
			notifications.GET("/:cif/count", handlers.GetNotificationCount)
			notifications.POST("/read", handlers.MarkNotificationAsRead)
			notifications.POST("/fcm-token", handlers.RegisterFCMToken)
			notifications.GET("/:cif/preferences", handlers.GetNotificationPreferences)
			notifications.PUT("/:cif/preferences", handlers.UpdateNotificationPreferences)
		}
		
		// Payment routes (QRIS, VA, E-Wallet, E-Money)
		payments := v1.Group("/payments")
		{
			// QRIS Payment
			payments.POST("/qris", handlers.ProcessQRISPayment)
			
			// Virtual Account Payment
			payments.POST("/va", handlers.ProcessVAPayment)
			
			// E-Wallet Top-up
			payments.POST("/ewallet/topup", handlers.ProcessEWalletTopup)
			payments.GET("/ewallet/providers", handlers.GetEWalletProviders)
			
			// E-Money Top-up
			payments.POST("/emoney/topup", handlers.ProcessEMoneyTopup)
			payments.GET("/emoney/providers", handlers.GetEMoneyProviders)
			
			// VA Information
			payments.GET("/va/providers", handlers.GetVAProviders)
		}
	}
	
	// Admin routes (unprotected for testing; add auth middleware in production)
	admin := router.Group("/api/v1/admin")
	{
		// Bank management routes
		banks := admin.Group("/banks")
		{
			banks.GET("", handlers.GetAllBanks)
		}
		
		// Fee management routes
		fees := admin.Group("/fees")
		{
			// Transfer fees
			fees.GET("/transfer", handlers.GetTransferFees)
			fees.PUT("/transfer", handlers.UpdateTransferFee)
			fees.POST("/transfer/calculate", handlers.CalculateTransferFeeHandler)
			
			// Service fees (e-wallet, e-money, VA, QRIS, etc)
			fees.GET("/services", handlers.GetServiceFees)
			fees.PUT("/services", handlers.UpdateServiceFee)
			fees.POST("/services/calculate", handlers.CalculateServiceFeeHandler)
			
			// Statistics
			fees.GET("/statistics", handlers.GetFeeStatistics)
		}
	}
}
