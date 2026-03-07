package handlers

import (
	"cbs-simulator/models"
	"cbs-simulator/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// === GENERAL LEDGER ===

func GetChartOfAccounts(c *gin.Context) {
	accountType := c.Query("type")
	accounts, err := services.GetChartOfAccounts(accountType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": accounts})
}

func GetJournalEntries(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	refType := c.Query("type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	entries, total, err := services.GetJournalEntries(dateFrom, dateTo, refType, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{
		"journal_entries": entries, "total": total, "page": page, "page_size": pageSize,
	}})
}

func GetJournalDetail(c *gin.Context) {
	journalID, _ := strconv.Atoi(c.Param("id"))
	lines, err := services.GetJournalLines(journalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": lines})
}

func GetTrialBalance(c *gin.Context) {
	asOfDate := c.Query("date")
	report, err := services.GetTrialBalance(asOfDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var totalDebit, totalCredit float64
	for _, tb := range report {
		totalDebit += tb.DebitBalance
		totalCredit += tb.CreditBalance
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{
		"as_of_date": asOfDate, "accounts": report,
		"total_debit": totalDebit, "total_credit": totalCredit,
		"is_balanced": totalDebit == totalCredit,
	}})
}

func GetGLAccountBalance(c *gin.Context) {
	code := c.Param("code")
	balance, err := services.GetGLAccountBalance(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"account_code": code, "balance": balance}})
}

// === CIF ENHANCEMENT ===

func GetCustomerOverview(c *gin.Context) {
	cif := c.Param("cif")
	view, err := services.GetSingleCustomerView(cif)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": view})
}

func UpdateCustomerExtended(c *gin.Context) {
	cif := c.Param("cif")
	var ext models.CustomerExtended
	if err := c.ShouldBindJSON(&ext); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	ext.CIF = cif
	if err := services.UpdateCustomerExtended(cif, ext); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Customer extended data updated"})
}

func SearchCustomers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "search query 'q' is required"})
		return
	}
	customers, err := services.SearchCustomers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": customers})
}

// === INTEREST ===

func GetInterestRates(c *gin.Context) {
	productType := c.Query("product_type")
	rates, err := services.GetInterestRates(productType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": rates})
}

func SimulateInterest(c *gin.Context) {
	var req struct {
		ProductType string  `json:"product_type" binding:"required"`
		Principal   float64 `json:"principal" binding:"required"`
		TenorMonths int     `json:"tenor_months"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	sim, err := services.SimulateInterest(req.ProductType, req.Principal, req.TenorMonths)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": sim})
}

// === STANDING INSTRUCTIONS ===

func CreateStandingInstruction(c *gin.Context) {
	var req services.CreateSIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	si, err := services.CreateStandingInstruction(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": si})
}

func GetStandingInstructions(c *gin.Context) {
	cif := c.Param("cif")
	instructions, err := services.GetSIByCIF(cif)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": instructions})
}

func PauseStandingInstruction(c *gin.Context) {
	si := c.Param("id")
	if err := services.PauseSI(si); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Standing instruction paused"})
}

func CancelStandingInstruction(c *gin.Context) {
	si := c.Param("id")
	if err := services.CancelSI(si); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Standing instruction cancelled"})
}

func GetSIHistory(c *gin.Context) {
	si := c.Param("id")
	history, err := services.GetSIExecutionHistory(si)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": history})
}

// === EOD ===

func RunEOD(c *gin.Context) {
	var req struct {
		ProcessDate string `json:"process_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	result, err := services.RunEOD(req.ProcessDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}

func GetEODStatus(c *gin.Context) {
	date := c.Param("date")
	logs, err := services.GetEODStatus(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": logs})
}

func GetEODHistory(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	logs, total, err := services.GetEODHistory(dateFrom, dateTo, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"logs": logs, "total": total, "page": page}})
}

// === ACCOUNT MANAGEMENT ===

func OpenAccountHandler(c *gin.Context) {
	var req services.OpenAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	resp, err := services.OpenAccount(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": resp})
}

func CloseAccountHandler(c *gin.Context) {
	accountNumber := c.Param("account_number")
	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)
	if err := services.CloseAccount(accountNumber, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Account closed"})
}

func GetDormantAccounts(c *gin.Context) {
	accounts, err := services.GetDormantAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": accounts})
}

func ReactivateAccountHandler(c *gin.Context) {
	accountNumber := c.Param("account_number")
	if err := services.ReactivateAccount(accountNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Account reactivated"})
}
