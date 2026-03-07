package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"cbs-simulator/services"
)

// ===== ACCOUNT HANDLERS =====

// GetAccountBalance retrieves account balance
func GetAccountBalance(c *gin.Context) {
	accountNumber := c.Param("account_number")
	
	account, err := services.GetAccountBalance(accountNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   account,
	})
}

// GetAccountsByCIF retrieves all accounts for a customer
func GetAccountsByCIF(c *gin.Context) {
	cif := c.Param("cif")
	
	accounts, err := services.GetAccountsByCIF(cif)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   accounts,
	})
}

// GetAccountStatement retrieves transaction history
func GetAccountStatement(c *gin.Context) {
	accountNumber := c.Param("account_number")
	
	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	transactions, err := services.GetAccountStatement(accountNumber, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   transactions,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(transactions),
		},
	})
}

// ===== TRANSFER HANDLERS =====

// IntraBankTransfer handles intrabank transfer
func IntraBankTransfer(c *gin.Context) {
	var req services.TransferRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}
	
	req.TransferType = "intra"
	
	transaction, err := services.ProcessIntraBankTransfer(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Transfer successful",
		"data":    transaction,
	})
}

// InterBankTransfer handles interbank transfer
func InterBankTransfer(c *gin.Context) {
	var req services.TransferRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
		})
		return
	}
	
	if req.TransferType == "" {
		req.TransferType = "inter"
	}
	
	transaction, err := services.ProcessInterBankTransfer(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Transfer successful",
		"data":    transaction,
	})
}

// GetTransaction retrieves transaction details
func GetTransaction(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	
	transaction, err := services.GetTransactionByTransactionID(transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   transaction,
	})
}

// ===== BILL PAYMENT HANDLERS =====

// GetBillerList returns list of billers
func GetBillerList(c *gin.Context) {
	billers := services.GetBillerList()
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   billers,
	})
}

// InquiryBill retrieves bill information
func InquiryBill(c *gin.Context) {
	billerCode := c.Query("biller_code")
	customerNumber := c.Query("customer_number")
	
	if billerCode == "" || customerNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "biller_code and customer_number are required",
		})
		return
	}
	
	bill, err := services.GetBillByCustomerNumber(billerCode, customerNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   bill,
	})
}

// PayBill processes bill payment
func PayBill(c *gin.Context) {
	var req struct {
		AccountNumber string `json:"account_number" binding:"required"`
		BillNumber    string `json:"bill_number" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
		})
		return
	}
	
	transaction, err := services.PayBill(req.AccountNumber, req.BillNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Bill payment successful",
		"data":    transaction,
	})
}

// GetAllBills retrieves all unpaid bills (for testing)
func GetAllBills(c *gin.Context) {
	bills, err := services.GetAllUnpaidBills()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   bills,
	})
}
