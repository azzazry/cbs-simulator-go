package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"cbs-simulator/services"
)

// Login handles customer login
func Login(c *gin.Context) {
	var req services.LoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}
	
	response, err := services.Authenticate(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}

// GetProfile retrieves customer profile
func GetProfile(c *gin.Context) {
	cif := c.Param("cif")
	
	customer, err := services.GetCustomerByCIF(cif)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   customer,
	})
}

// ChangePIN handles PIN change request
func ChangePIN(c *gin.Context) {
	var req struct {
		CIF    string `json:"cif" binding:"required"`
		OldPIN string `json:"old_pin" binding:"required"`
		NewPIN string `json:"new_pin" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
		})
		return
	}
	
	if err := services.ChangePIN(req.CIF, req.OldPIN, req.NewPIN); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "PIN changed successfully",
	})
}

// Register handles customer registration
func Register(c *gin.Context) {
	var req services.RegisterRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}
	
	response, err := services.RegisterCustomer(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   response,
	})
}
