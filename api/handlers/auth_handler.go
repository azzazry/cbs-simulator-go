package handlers

import (
	"net/http"
	"time"

	"cbs-simulator/services"

	"github.com/gin-gonic/gin"
)

// Login handles user authentication and returns JWT tokens
func Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
		})
		return
	}

	response, err := services.Authenticate(req, c.ClientIP())
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

// GetProfile returns the authenticated user's profile
func GetProfile(c *gin.Context) {
	cif, exists := c.Get("cif")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Authentication required",
		})
		return
	}

	customer, err := services.GetCustomerByCIF(cif.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Get roles
	roles, _ := services.GetUserRoleDetails(cif.(string))

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"customer": customer,
			"roles":    roles,
		},
	})
}

// ChangePIN handles PIN change requests
func ChangePIN(c *gin.Context) {
	var req struct {
		CIF    string `json:"cif"`
		OldPIN string `json:"old_pin"`
		NewPIN string `json:"new_pin"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
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
			"message": "Invalid request body: " + err.Error(),
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

// Logout invalidates the current JWT token
func Logout(c *gin.Context) {
	jti, _ := c.Get("jti")
	cif, _ := c.Get("cif")
	expiresAt, _ := c.Get("token_expires_at")

	if jti == nil || cif == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid token",
		})
		return
	}

	expiry, ok := expiresAt.(time.Time)
	if !ok {
		expiry = time.Now().Add(24 * time.Hour) // default fallback
	}

	if err := services.LogoutUser(jti.(string), cif.(string), expiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logged out successfully",
	})
}

// RefreshToken generates new token pair from a valid refresh token
func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "refresh_token is required",
		})
		return
	}

	tokenPair, err := services.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   tokenPair,
	})
}

// RequestOTP generates and sends an OTP code
func RequestOTP(c *gin.Context) {
	var req struct {
		CIF     string `json:"cif" binding:"required"`
		OTPType string `json:"otp_type" binding:"required"` // unlock_account, reset_pin
		Channel string `json:"channel"`                     // sms, email
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "cif and otp_type are required",
		})
		return
	}

	if req.Channel == "" {
		req.Channel = "sms"
	}

	otp, err := services.GenerateOTP(req.CIF, req.OTPType, req.Channel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// In simulator mode: return OTP in response for testing
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "OTP sent via " + req.Channel,
		"data": gin.H{
			"otp_type": req.OTPType,
			"channel":  req.Channel,
			"otp":      otp, // NOTE: In production, this should NOT be returned
		},
	})
}

// VerifyOTPHandler verifies an OTP code
func VerifyOTPHandler(c *gin.Context) {
	var req struct {
		CIF     string `json:"cif" binding:"required"`
		OTP     string `json:"otp" binding:"required"`
		OTPType string `json:"otp_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "cif, otp, and otp_type are required",
		})
		return
	}

	if err := services.VerifyOTP(req.CIF, req.OTP, req.OTPType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "OTP verified successfully",
	})
}

// VerifyEKYC performs e-KYC verification
func VerifyEKYC(c *gin.Context) {
	var req services.EKYCVerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "cif and id_card_number are required",
		})
		return
	}

	result, err := services.VerifyEKYC(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// UnlockAccount handles self-service account unlock (e-KYC + OTP)
func UnlockAccount(c *gin.Context) {
	var req struct {
		CIF            string `json:"cif" binding:"required"`
		OTP            string `json:"otp" binding:"required"`
		VerificationID string `json:"verification_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "cif, otp, and verification_id are required",
		})
		return
	}

	if err := services.SelfServiceUnlockAccount(req.CIF, req.OTP, req.VerificationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Account unlocked successfully. You can now login with your PIN.",
	})
}

// ResetPINHandler handles PIN reset after e-KYC + OTP verification
func ResetPINHandler(c *gin.Context) {
	var req struct {
		CIF            string `json:"cif" binding:"required"`
		NewPIN         string `json:"new_pin" binding:"required"`
		VerificationID string `json:"verification_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "cif, new_pin, and verification_id are required",
		})
		return
	}

	if err := services.ResetPIN(req.CIF, req.NewPIN, req.VerificationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "PIN reset successfully. You can now login with your new PIN.",
	})
}
