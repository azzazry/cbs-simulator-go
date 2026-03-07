package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"cbs-simulator/services"
)

// ===== CARD HANDLERS =====

// GetCardsByCIF retrieves all cards for a customer
func GetCardsByCIF(c *gin.Context) {
	cif := c.Param("cif")
	
	cards, err := services.GetCardsByCIF(cif)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   cards,
	})
}

// GetCardDetails retrieves card details
func GetCardDetails(c *gin.Context) {
	cardNumber := c.Param("card_number")
	
	card, err := services.GetCardByNumber(cardNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   card,
	})
}

// BlockCard blocks a card
func BlockCard(c *gin.Context) {
	var req struct {
		CardNumber string `json:"card_number" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
		})
		return
	}
	
	if err := services.BlockCard(req.CardNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Card blocked successfully",
	})
}

// UnblockCard unblocks a card
func UnblockCard(c *gin.Context) {
	var req struct {
		CardNumber string `json:"card_number" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
		})
		return
	}
	
	if err := services.UnblockCard(req.CardNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Card unblocked successfully",
	})
}

// ===== LOAN HANDLERS =====

// GetLoansByCIF retrieves all loans for a customer
func GetLoansByCIF(c *gin.Context) {
	cif := c.Param("cif")
	
	loans, err := services.GetLoansByCIF(cif)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   loans,
	})
}

// GetLoanDetails retrieves loan details
func GetLoanDetails(c *gin.Context) {
	loanNumber := c.Param("loan_number")
	
	loan, err := services.GetLoanByNumber(loanNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   loan,
	})
}

// ===== DEPOSIT HANDLERS =====

// GetDepositsByCIF retrieves all deposits for a customer
func GetDepositsByCIF(c *gin.Context) {
	cif := c.Param("cif")
	
	deposits, err := services.GetDepositsByCIF(cif)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   deposits,
	})
}

// GetDepositDetails retrieves deposit details
func GetDepositDetails(c *gin.Context) {
	depositNumber := c.Param("deposit_number")
	
	deposit, err := services.GetDepositByNumber(depositNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   deposit,
	})
}
