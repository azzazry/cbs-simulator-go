package routes

import (
	"github.com/gin-gonic/gin"
	"cbs-simulator/api/handlers"
	"cbs-simulator/api/middleware"
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
	}
}
