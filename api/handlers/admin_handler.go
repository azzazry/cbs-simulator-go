package handlers

import (
	"net/http"

	"cbs-simulator/services"

	"github.com/gin-gonic/gin"
)

// GetAllBanks retrieves all active banks
// @Summary Get all banks
// @Description Retrieve list of all supported banks for transfers
// @Tags Admin - Banks
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /admin/banks [get]
func GetAllBanks(c *gin.Context) {
	banks, err := services.GetAllBanks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   banks,
	})
}

// GetTransferFees retrieves all transfer fees
// @Summary Get all transfer fees
// @Description Retrieve all interbank transfer fee configurations
// @Tags Admin - Fees
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /admin/fees/transfer [get]
func GetTransferFees(c *gin.Context) {
	fees, err := services.GetAllTransferFees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   fees,
	})
}

// UpdateTransferFeeRequest represents request to update transfer fee
type UpdateTransferFeeRequest struct {
	DestinationBankCode string  `json:"destination_bank_code" binding:"required"`
	FeeAmount           float64 `json:"fee_amount" binding:"required,min=0"`
	FeeType             string  `json:"fee_type" binding:"required,oneof=flat percentage"`
}

// UpdateTransferFee updates transfer fee configuration for outbound transfer to specific bank
// @Summary Update transfer fee
// @Description Update fee for outbound interbank transfer (our bank → destination bank)
// @Tags Admin - Fees
// @Accept json
// @Produce json
// @Param request body UpdateTransferFeeRequest true "Update fee request"
// @Success 200 {object} map[string]interface{}
// @Router /admin/fees/transfer [put]
func UpdateTransferFee(c *gin.Context) {
	var req UpdateTransferFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.UpdateTransferFee(req.DestinationBankCode, req.FeeAmount, req.FeeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Transfer fee updated successfully",
		"data": gin.H{
			"destination_bank_code": req.DestinationBankCode,
			"fee_amount":            req.FeeAmount,
			"fee_type":              req.FeeType,
		},
	})
}

// GetServiceFees retrieves all service fees
// @Summary Get all service fees
// @Description Retrieve all service fee configurations (e-wallet, e-money, VA, QRIS, etc)
// @Tags Admin - Fees
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /admin/fees/services [get]
func GetServiceFees(c *gin.Context) {
	serviceType := c.Query("type") // Optional filter by service type

	fees, err := services.GetAllServiceFees(serviceType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   fees,
	})
}

// UpdateServiceFeeRequest represents request to update service fee
type UpdateServiceFeeRequest struct {
	ServiceCode   string  `json:"service_code" binding:"required"`
	FeeAmount     float64 `json:"fee_amount"`
	FeePercentage float64 `json:"fee_percentage"`
	FeeType       string  `json:"fee_type" binding:"required,oneof=flat percentage"`
}

// UpdateServiceFee updates service fee configuration
// @Summary Update service fee
// @Description Update fee for service (e-wallet, e-money, VA payment, QRIS, etc)
// @Tags Admin - Fees
// @Accept json
// @Produce json
// @Param request body UpdateServiceFeeRequest true "Update service fee request"
// @Success 200 {object} map[string]interface{}
// @Router /admin/fees/services [put]
func UpdateServiceFee(c *gin.Context) {
	var req UpdateServiceFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.UpdateServiceFee(req.ServiceCode, req.FeeAmount, req.FeePercentage, req.FeeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Service fee updated successfully",
		"data": gin.H{
			"service_code":   req.ServiceCode,
			"fee_amount":     req.FeeAmount,
			"fee_percentage": req.FeePercentage,
			"fee_type":       req.FeeType,
		},
	})
}

// GetTransferFeeRequest represents request to get transfer fee calculation
type GetTransferFeeRequest struct {
	DestinationBankCode string  `json:"destination_bank_code" binding:"required"`
	Amount              float64 `json:"amount" binding:"required,min=1"`
}

// CalculateTransferFee calculates fee for given transfer parameters
// @Summary Calculate transfer fee
// @Description Calculate the fee that will be charged for outbound interbank transfer
// @Tags Admin - Fees
// @Accept json
// @Produce json
// @Param request body GetTransferFeeRequest true "Fee calculation request"
// @Success 200 {object} map[string]interface{}
// @Router /admin/fees/transfer/calculate [post]
func CalculateTransferFeeHandler(c *gin.Context) {
	var req GetTransferFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fee, err := services.CalculateTransferFee(req.DestinationBankCode, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalAmount := req.Amount + fee

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"destination_bank_code": req.DestinationBankCode,
			"transfer_amount":       req.Amount,
			"fee":                   fee,
			"total_amount":          totalAmount,
		},
	})
}

// GetServiceFeeRequest represents request to get service fee calculation
type GetServiceFeeRequest struct {
	ServiceCode string  `json:"service_code" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,min=1"`
}

// CalculateServiceFeeHandler calculates fee for given service parameters
// @Summary Calculate service fee
// @Description Calculate the fee that will be charged for service (e-wallet, VA, QRIS, etc)
// @Tags Admin - Fees
// @Accept json
// @Produce json
// @Param request body GetServiceFeeRequest true "Fee calculation request"
// @Success 200 {object} map[string]interface{}
// @Router /admin/fees/services/calculate [post]
func CalculateServiceFeeHandler(c *gin.Context) {
	var req GetServiceFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fee, err := services.CalculateServiceFee(req.ServiceCode, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalAmount := req.Amount + fee

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"service_code":   req.ServiceCode,
			"service_amount": req.Amount,
			"fee":            fee,
			"total_amount":   totalAmount,
		},
	})
}

// FeeStatisticsRequest represents request for fee statistics
type FeeStatisticsRequest struct {
	FeeType     string `json:"fee_type" binding:"oneof=transfer service"` // transfer or service
	ServiceType string `json:"service_type"`                              // optional for service fees
}

// GetFeeStatistics retrieves fee statistics
// @Summary Get fee statistics
// @Description Get fee configuration statistics
// @Tags Admin - Fees
// @Produce json
// @Param type query string false "Fee type: transfer or service"
// @Param service_type query string false "Service type (for service fees)"
// @Success 200 {object} map[string]interface{}
// @Router /admin/fees/statistics [get]
func GetFeeStatistics(c *gin.Context) {
	feeType := c.Query("type")
	_ = c.Query("service_type")

	if feeType == "service" {
		fees, err := services.GetAllServiceFees("")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Group by service type
		serviceTypeMap := make(map[string]int)
		for _, fee := range fees {
			serviceTypeMap[fee.ServiceType]++
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"total_services":  len(fees),
				"by_service_type": serviceTypeMap,
				"services":        fees,
			},
		})
	} else {
		fees, err := services.GetAllTransferFees()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"total_transfer_routes": len(fees),
				"transfers":             fees,
			},
		})
	}
}
