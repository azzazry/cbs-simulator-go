package handlers

import (
	"net/http"
	"strconv"

	"cbs-simulator/models"
	"cbs-simulator/services"

	"github.com/gin-gonic/gin"
)

// GetNotifications retrieves notification history for a customer with pagination
func GetNotifications(c *gin.Context) {
	cif := c.Param("cif")
	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	notifications, err := services.GetNotifications(cif, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch notifications",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   notifications,
	})
}

// GetNotificationCount retrieves count of unread notifications
func GetNotificationCount(c *gin.Context) {
	cif := c.Param("cif")

	count, err := services.GetNotificationCount(cif)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get notification count",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"unread_count": count,
		},
	})
}

// MarkNotificationAsRead marks a notification as read
func MarkNotificationAsRead(c *gin.Context) {
	var req struct {
		NotificationID int `json:"notification_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request",
		})
		return
	}

	if err := services.MarkAsRead(req.NotificationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to mark notification as read",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Notification marked as read",
	})
}

// RegisterFCMToken registers device FCM token for push notifications
func RegisterFCMToken(c *gin.Context) {
	cif, _ := c.Get("cif")

	var req struct {
		DeviceToken string `json:"device_token" binding:"required"`
		DeviceType  string `json:"device_type"`
		DeviceName  string `json:"device_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request",
		})
		return
	}

	if err := services.RegisterFCMToken(cif.(string), req.DeviceToken, req.DeviceType, req.DeviceName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to register FCM token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "FCM token registered successfully",
	})
}

// GetNotificationPreferences retrieves user notification settings
func GetNotificationPreferences(c *gin.Context) {
	cif := c.Param("cif")

	prefs, err := services.GetNotificationPreferences(cif)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch notification preferences",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   prefs,
	})
}

// UpdateNotificationPreferences updates user notification settings
func UpdateNotificationPreferences(c *gin.Context) {
	cif := c.Param("cif")

	var prefs models.NotificationPreference
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
		})
		return
	}

	prefs.CIF = cif
	if err := services.UpdateNotificationPreferences(cif, prefs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update notification preferences",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Notification preferences updated successfully",
	})
}

func UnregisterFCMToken(c *gin.Context) {
	cif, _ := c.Get("cif")

	var req struct {
		DeviceToken string `json:"device_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "device_token is required",
		})
		return
	}

	if err := services.UnregisterFCMToken(cif.(string), req.DeviceToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to unregister FCM token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Device token unregistered successfully",
	})
}

func GetFCMDevices(c *gin.Context) {
	cif, _ := c.Get("cif")

	devices, err := services.GetFCMDevices(cif.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch devices",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   devices,
	})
}
