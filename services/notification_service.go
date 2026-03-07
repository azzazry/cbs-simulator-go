package services

import (
	"database/sql"
	"log"

	"cbs-simulator/database"
	"cbs-simulator/models"
)

// SaveNotification saves notification to database
func SaveNotification(cif, notificationType, title, message, transactionID string) error {
	query := `INSERT INTO notifications (cif, notification_type, title, message, transaction_id) 
	          VALUES (?, ?, ?, ?, ?)`

	_, err := database.DB.Exec(query, cif, notificationType, title, message, transactionID)
	if err != nil {
		log.Printf("Error saving notification: %v", err)
	}
	return err
}

// GetNotifications retrieves notification history for a customer with pagination
func GetNotifications(cif string, limit int, offset int) ([]models.Notification, error) {
	var notifications []models.Notification

	query := `SELECT id, cif, notification_type, title, message, transaction_id, is_read, created_at 
	          FROM notifications 
	          WHERE cif = ? 
	          ORDER BY created_at DESC 
	          LIMIT ? OFFSET ?`

	rows, err := database.DB.Query(query, cif, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notif models.Notification
		err := rows.Scan(&notif.ID, &notif.CIF, &notif.NotificationType,
			&notif.Title, &notif.Message, &notif.TransactionID,
			&notif.IsRead, &notif.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}

	return notifications, nil
}

// MarkAsRead marks a notification as read
func MarkAsRead(notificationID int) error {
	_, err := database.DB.Exec(`UPDATE notifications SET is_read = 1 WHERE id = ?`, notificationID)
	return err
}

// RegisterFCMToken saves or updates device FCM token
func RegisterFCMToken(cif, deviceToken, deviceType, deviceName string) error {
	query := `INSERT OR REPLACE INTO fcm_tokens (cif, device_token, device_type, device_name, is_active) 
	          VALUES (?, ?, ?, ?, 1)`

	_, err := database.DB.Exec(query, cif, deviceToken, deviceType, deviceName)
	return err
}

// GetFCMTokens gets all active device tokens for a customer
func GetFCMTokens(cif string) ([]string, error) {
	var tokens []string

	query := `SELECT device_token FROM fcm_tokens WHERE cif = ? AND is_active = 1`
	rows, err := database.DB.Query(query, cif)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// SendPushNotification simulates sending push notification to FCM
// In production, integrate with Firebase Admin SDK
func SendPushNotification(cif, title, message string) error {
	tokens, err := GetFCMTokens(cif)
	if err != nil || len(tokens) == 0 {
		log.Printf("No device tokens found for CIF: %s", cif)
		return nil // Not critical error, just log
	}

	// TODO: Integrate dengan Firebase Admin SDK
	// For now, just log the push notifications
	for _, token := range tokens {
		log.Printf("PUSH NOTIFICATION - Token: %s, Title: %s, Message: %s", token, title, message)
	}

	return nil
}

// GetNotificationCount returns count of unread notifications
func GetNotificationCount(cif string) (int, error) {
	var count int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM notifications WHERE cif = ? AND is_read = 0`, cif).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return count, nil
}

// UpdateNotificationPreferences updates user notification settings
func UpdateNotificationPreferences(cif string, prefs models.NotificationPreference) error {
	// Check if preference already exists
	var id int
	err := database.DB.QueryRow(`SELECT id FROM notification_preferences WHERE cif = ?`, cif).Scan(&id)

	if err == sql.ErrNoRows {
		// Insert new preference
		query := `INSERT INTO notification_preferences (cif, transfer_notification, payment_notification, deposit_notification, loan_notification, promotion_notification) 
		          VALUES (?, ?, ?, ?, ?, ?)`
		_, err = database.DB.Exec(query, cif, prefs.TransferNotification, prefs.PaymentNotification,
			prefs.DepositNotification, prefs.LoanNotification, prefs.PromotionNotification)
		return err
	}

	// Update existing preference
	query := `UPDATE notification_preferences 
	          SET transfer_notification = ?, payment_notification = ?, deposit_notification = ?, 
	              loan_notification = ?, promotion_notification = ?, updated_at = CURRENT_TIMESTAMP 
	          WHERE cif = ?`
	_, err = database.DB.Exec(query, prefs.TransferNotification, prefs.PaymentNotification,
		prefs.DepositNotification, prefs.LoanNotification, prefs.PromotionNotification, cif)
	return err
}

// GetNotificationPreferences retrieves user notification settings
func GetNotificationPreferences(cif string) (*models.NotificationPreference, error) {
	var prefs models.NotificationPreference

	query := `SELECT id, cif, transfer_notification, payment_notification, deposit_notification, 
	          loan_notification, promotion_notification, updated_at 
	          FROM notification_preferences WHERE cif = ?`

	err := database.DB.QueryRow(query, cif).Scan(&prefs.ID, &prefs.CIF, &prefs.TransferNotification,
		&prefs.PaymentNotification, &prefs.DepositNotification,
		&prefs.LoanNotification, &prefs.PromotionNotification, &prefs.UpdatedAt)

	if err == sql.ErrNoRows {
		// Return default preferences
		return &models.NotificationPreference{
			CIF:                   cif,
			TransferNotification:  1,
			PaymentNotification:   1,
			DepositNotification:   1,
			LoanNotification:      1,
			PromotionNotification: 1,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &prefs, nil
}
