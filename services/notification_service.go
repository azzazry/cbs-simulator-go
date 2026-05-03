package services

import (
	"database/sql"
	"log"

	"cbs-simulator/database"
	"cbs-simulator/models"
)

func SaveNotification(cif, notificationType, title, message, transactionID string) error {
	query := `INSERT INTO notifications (cif, notification_type, title, message, transaction_id)
	          VALUES ($1, $2, $3, $4, $5)`

	_, err := database.DB.Exec(query, cif, notificationType, title, message, transactionID)
	if err != nil {
		log.Printf("Error saving notification: %v", err)
	}
	return err
}

func GetNotifications(cif string, limit int, offset int) ([]models.Notification, error) {
	var notifications []models.Notification

	query := `SELECT id, cif, notification_type, title, message, transaction_id, is_read, created_at
	          FROM notifications
	          WHERE cif = $1
	          ORDER BY created_at DESC
	          LIMIT $2 OFFSET $3`

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

func MarkAsRead(notificationID int) error {
	_, err := database.DB.Exec(`UPDATE notifications SET is_read = TRUE WHERE id = $1`, notificationID)
	return err
}

// RegisterFCMToken saves or updates device FCM token (upsert via ON CONFLICT)
func RegisterFCMToken(cif, deviceToken, deviceType, deviceName string) error {
	query := `INSERT INTO fcm_tokens (cif, device_token, device_type, device_name, is_active)
	          VALUES ($1, $2, $3, $4, TRUE)
	          ON CONFLICT (device_token) DO UPDATE
	          SET cif = EXCLUDED.cif,
	              device_type = EXCLUDED.device_type,
	              device_name = EXCLUDED.device_name,
	              is_active = TRUE,
	              updated_at = NOW()`

	_, err := database.DB.Exec(query, cif, deviceToken, deviceType, deviceName)
	return err
}

func GetFCMTokens(cif string) ([]string, error) {
	var tokens []string

	query := `SELECT device_token FROM fcm_tokens WHERE cif = $1 AND is_active = TRUE`
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

func SendPushNotification(cif, title, message string) error {
	tokens, err := GetFCMTokens(cif)
	if err != nil || len(tokens) == 0 {
		log.Printf("No device tokens found for CIF: %s", cif)
		return nil
	}

	// TODO: Integrate dengan Firebase Admin SDK
	for _, token := range tokens {
		log.Printf("PUSH NOTIFICATION - Token: %s, Title: %s, Message: %s", token, title, message)
	}

	return nil
}

func GetNotificationCount(cif string) (int, error) {
	var count int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM notifications WHERE cif = $1 AND is_read = FALSE`, cif).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return count, nil
}

func UpdateNotificationPreferences(cif string, prefs models.NotificationPreference) error {
	var id int
	err := database.DB.QueryRow(`SELECT id FROM notification_preferences WHERE cif = $1`, cif).Scan(&id)

	if err == sql.ErrNoRows {
		query := `INSERT INTO notification_preferences (cif, transfer_notification, payment_notification, deposit_notification, loan_notification, promotion_notification)
		          VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = database.DB.Exec(query, cif, prefs.TransferNotification, prefs.PaymentNotification,
			prefs.DepositNotification, prefs.LoanNotification, prefs.PromotionNotification)
		return err
	}

	query := `UPDATE notification_preferences
	          SET transfer_notification = $1, payment_notification = $2, deposit_notification = $3,
	              loan_notification = $4, promotion_notification = $5, updated_at = NOW()
	          WHERE cif = $6`
	_, err = database.DB.Exec(query, prefs.TransferNotification, prefs.PaymentNotification,
		prefs.DepositNotification, prefs.LoanNotification, prefs.PromotionNotification, cif)
	return err
}

func GetNotificationPreferences(cif string) (*models.NotificationPreference, error) {
	var prefs models.NotificationPreference

	query := `SELECT id, cif, transfer_notification, payment_notification, deposit_notification,
	          loan_notification, promotion_notification, updated_at
	          FROM notification_preferences WHERE cif = $1`

	err := database.DB.QueryRow(query, cif).Scan(&prefs.ID, &prefs.CIF, &prefs.TransferNotification,
		&prefs.PaymentNotification, &prefs.DepositNotification,
		&prefs.LoanNotification, &prefs.PromotionNotification, &prefs.UpdatedAt)

	if err == sql.ErrNoRows {
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
