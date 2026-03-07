package handlers

import (
	"net/http"

	"cbs-simulator/services"

	"github.com/gin-gonic/gin"
)

// ProcessQRISPayment handles QRIS payment
func ProcessQRISPayment(c *gin.Context) {
	var req services.QRISPaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	transaction, err := services.ProcessQRISPayment(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "QRIS Payment successful",
		"data":    transaction,
	})
}

// ProcessVAPayment handles Virtual Account payment
func ProcessVAPayment(c *gin.Context) {
	var req services.VAPaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	transaction, err := services.ProcessVAPayment(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Virtual Account Payment successful",
		"data":    transaction,
	})
}

// ProcessEWalletTopup handles e-wallet top-up
func ProcessEWalletTopup(c *gin.Context) {
	var req services.EWalletTopupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	transaction, err := services.ProcessEWalletTopup(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "E-Wallet Top-up successful",
		"data":    transaction,
	})
}

// ProcessEMoneyTopup handles e-money top-up
func ProcessEMoneyTopup(c *gin.Context) {
	var req services.EMoneyTopupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	transaction, err := services.ProcessEMoneyTopup(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "E-Money Top-up successful",
		"data":    transaction,
	})
}

// GetServiceFees retrieves all service fees (for admin/info)

// GetEWalletProviders returns list of supported e-wallet providers
func GetEWalletProviders(c *gin.Context) {
	providers := []map[string]interface{}{
		{
			"code":    "OVO",
			"name":    "OVO",
			"logo":    "ovo",
			"fee":     2500,
			"minimum": 10000,
			"maximum": 10000000,
		},
		{
			"code":    "DANA",
			"name":    "DANA",
			"logo":    "dana",
			"fee":     2500,
			"minimum": 10000,
			"maximum": 10000000,
		},
		{
			"code":    "GOPAY",
			"name":    "GoPay",
			"logo":    "gopay",
			"fee":     2500,
			"minimum": 10000,
			"maximum": 10000000,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   providers,
	})
}

// GetEMoneyProviders returns list of supported e-money providers
func GetEMoneyProviders(c *gin.Context) {
	providers := []map[string]interface{}{
		{
			"code":    "LINKAJA",
			"name":    "LinkAja",
			"logo":    "linkaja",
			"fee":     2500,
			"minimum": 10000,
			"maximum": 10000000,
		},
		{
			"code":    "MANDIRIEMONEY",
			"name":    "Mandiri e-Money",
			"logo":    "mandiri_emoney",
			"fee":     2500,
			"minimum": 10000,
			"maximum": 10000000,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   providers,
	})
}

// GetVAProviders returns list of supported VA providers
func GetVAProviders(c *gin.Context) {
	providers := []map[string]interface{}{
		{
			"code":    "MANDIRI",
			"name":    "Mandiri Virtual Account",
			"bank":    "Bank Mandiri",
			"fee":     0,
			"minimum": 1000,
			"maximum": 999999999,
		},
		{
			"code":    "BCA",
			"name":    "BCA Virtual Account",
			"bank":    "Bank BCA",
			"fee":     0,
			"minimum": 1000,
			"maximum": 999999999,
		},
		{
			"code":    "BRI",
			"name":    "BRI Virtual Account",
			"bank":    "Bank BRI",
			"fee":     0,
			"minimum": 1000,
			"maximum": 999999999,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   providers,
	})
}
