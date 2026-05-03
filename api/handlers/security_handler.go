package handlers

import (
	"net/http"
	"strconv"

	"cbs-simulator/services"

	"github.com/gin-gonic/gin"
)

// GetAuditLogs returns paginated audit logs (admin only)
func GetAuditLogs(c *gin.Context) {
	cif := c.Query("cif")
	action := c.Query("action")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := services.GetAuditLogs(cif, action, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"audit_logs": logs,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
		},
	})
}

// GetTransactionLimits returns all transaction limits (admin only)
// Query param ?role=customer untuk filter per role
func GetTransactionLimits(c *gin.Context) {
	roleName := c.Query("role")
	transactionType := c.Query("type")

	// Filter per role+type
	if roleName != "" && transactionType != "" {
		limit, err := services.GetTransactionLimit(roleName, transactionType)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   limit,
		})
		return
	}

	// Return semua limit
	limits, err := services.GetAllTransactionLimits()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   limits,
	})
}

// UpdateTransactionLimit updates a specific transaction limit (admin only)
// Body: { "role_name": "customer", "transaction_type": "transfer_intra", "daily_limit": 5000000, "per_transaction_limit": 1000000, "monthly_limit": 50000000 }
func UpdateTransactionLimit(c *gin.Context) {
	var req struct {
		RoleName            string  `json:"role_name" binding:"required"`
		TransactionType     string  `json:"transaction_type" binding:"required"`
		DailyLimit          float64 `json:"daily_limit"`
		PerTransactionLimit float64 `json:"per_transaction_limit"`
		MonthlyLimit        float64 `json:"monthly_limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "role_name and transaction_type are required",
		})
		return
	}

	if err := services.UpdateTransactionLimit(req.RoleName, req.TransactionType, req.DailyLimit, req.PerTransactionLimit, req.MonthlyLimit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Transaction limit updated for " + req.RoleName + " / " + req.TransactionType,
	})
}

// GetUserRolesHandler returns roles for a specific user or all roles (admin only)
func GetUserRolesHandler(c *gin.Context) {
	cif := c.Query("cif")

	if cif != "" {
		roles, err := services.GetUserRoleDetails(cif)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   gin.H{"cif": cif, "roles": roles},
		})
		return
	}

	roles, err := services.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   roles,
	})
}

// AssignRoleHandler assigns a role to a user (admin only)
func AssignRoleHandler(c *gin.Context) {
	var req struct {
		CIF      string `json:"cif" binding:"required"`
		RoleName string `json:"role_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "cif and role_name are required",
		})
		return
	}

	role, err := services.GetRoleByName(req.RoleName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	adminCIF, _ := c.Get("cif")

	if err := services.AssignRole(req.CIF, role.ID, adminCIF.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Role " + req.RoleName + " assigned to " + req.CIF,
	})
}

// AdminUnlockAccount force-unlocks an account (admin only, bypasses e-KYC)
func AdminUnlockAccount(c *gin.Context) {
	var req struct {
		CIF string `json:"cif" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "cif is required",
		})
		return
	}

	if err := services.UnlockAccount(req.CIF); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Account " + req.CIF + " unlocked by admin",
	})
}
