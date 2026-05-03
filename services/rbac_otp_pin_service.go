package services

import (
	"cbs-simulator/config"
	"cbs-simulator/database"
	"cbs-simulator/models"
	"fmt"
	"math/rand"
	"time"
)

// =================== RBAC ===================

func GetUserRoles(cif string) ([]string, error) {
	query := `SELECT r.role_name FROM user_roles ur
	          JOIN roles r ON ur.role_id = r.id
	          WHERE ur.cif = $1 AND r.is_active = TRUE`

	rows, err := database.DB.Query(query, cif)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var roleName string
		if err := rows.Scan(&roleName); err != nil {
			return nil, err
		}
		roles = append(roles, roleName)
	}

	if len(roles) == 0 {
		roles = append(roles, "customer")
	}

	return roles, nil
}

func GetPrimaryRole(cif string) (string, error) {
	roles, err := GetUserRoles(cif)
	if err != nil {
		return "customer", err
	}

	priority := map[string]int{"admin": 4, "supervisor": 3, "teller": 2, "customer": 1}
	bestRole := "customer"
	bestPriority := 0

	for _, role := range roles {
		if p, ok := priority[role]; ok && p > bestPriority {
			bestPriority = p
			bestRole = role
		}
	}

	return bestRole, nil
}

func HasRole(cif, roleName string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_roles ur
	          JOIN roles r ON ur.role_id = r.id
	          WHERE ur.cif = $1 AND r.role_name = $2 AND r.is_active = TRUE`
	err := database.DB.QueryRow(query, cif, roleName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// AssignRole uses ON CONFLICT to safely ignore duplicates
func AssignRole(cif string, roleID int, assignedBy string) error {
	query := `INSERT INTO user_roles (cif, role_id, assigned_by) VALUES ($1, $2, $3)
	          ON CONFLICT (cif, role_id) DO NOTHING`
	_, err := database.DB.Exec(query, cif, roleID, assignedBy)
	return err
}

func RemoveRole(cif string, roleID int) error {
	query := `DELETE FROM user_roles WHERE cif = $1 AND role_id = $2`
	_, err := database.DB.Exec(query, cif, roleID)
	return err
}

func GetAllRoles() ([]models.Role, error) {
	query := `SELECT id, role_name, description, is_active, created_at FROM roles`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.RoleName, &role.Description, &role.IsActive, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func GetUserRoleDetails(cif string) ([]models.Role, error) {
	query := `SELECT r.id, r.role_name, r.description, r.is_active, r.created_at
	          FROM user_roles ur JOIN roles r ON ur.role_id = r.id
	          WHERE ur.cif = $1 AND r.is_active = TRUE`
	rows, err := database.DB.Query(query, cif)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.RoleName, &role.Description, &role.IsActive, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func GetRoleByName(roleName string) (*models.Role, error) {
	var role models.Role
	query := `SELECT id, role_name, description, is_active, created_at FROM roles WHERE role_name = $1`
	err := database.DB.QueryRow(query, roleName).Scan(&role.ID, &role.RoleName, &role.Description, &role.IsActive, &role.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("role not found: %s", roleName)
	}
	return &role, nil
}

// =================== OTP ===================

func GenerateOTP(cif, otpType, channel string) (string, error) {
	InvalidateOTP(cif, otpType)

	otpLength := config.AppConfig.OTPLength
	otp := generateRandomOTP(otpLength)

	expiryMinutes := config.AppConfig.OTPExpiry
	expiresAt := time.Now().Add(time.Duration(expiryMinutes) * time.Minute)

	query := `INSERT INTO otp_codes (cif, otp_code, otp_type, channel, expires_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := database.DB.Exec(query, cif, otp, otpType, channel, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to store OTP: %v", err)
	}

	return otp, nil
}

func VerifyOTP(cif, otp, otpType string) error {
	var id int
	var expiresAt time.Time
	var isUsed bool

	query := `SELECT id, expires_at, is_used FROM otp_codes
	          WHERE cif = $1 AND otp_code = $2 AND otp_type = $3
	          ORDER BY created_at DESC LIMIT 1`
	err := database.DB.QueryRow(query, cif, otp, otpType).Scan(&id, &expiresAt, &isUsed)
	if err != nil {
		return fmt.Errorf("invalid OTP code")
	}

	if isUsed {
		return fmt.Errorf("OTP has already been used")
	}

	if time.Now().After(expiresAt) {
		return fmt.Errorf("OTP has expired")
	}

	_, err = database.DB.Exec(`UPDATE otp_codes SET is_used = TRUE WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to mark OTP as used: %v", err)
	}

	return nil
}

func InvalidateOTP(cif, otpType string) {
	database.DB.Exec(`UPDATE otp_codes SET is_used = TRUE WHERE cif = $1 AND otp_type = $2 AND is_used = FALSE`, cif, otpType)
}

func generateRandomOTP(length int) string {
	digits := "0123456789"
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = digits[rand.Intn(len(digits))]
	}
	return string(otp)
}

// =================== PIN POLICY ===================

func ValidatePINPolicy(pin string) error {
	minLen := config.AppConfig.PINMinLength
	maxLen := config.AppConfig.PINMaxLength

	if len(pin) < minLen || len(pin) > maxLen {
		return fmt.Errorf("PIN must be %d digits", minLen)
	}

	for _, c := range pin {
		if c < '0' || c > '9' {
			return fmt.Errorf("PIN must contain only digits")
		}
	}

	isAscending := true
	isDescending := true
	for i := 1; i < len(pin); i++ {
		if pin[i] != pin[i-1]+1 {
			isAscending = false
		}
		if pin[i] != pin[i-1]-1 {
			isDescending = false
		}
	}
	if isAscending || isDescending {
		return fmt.Errorf("PIN cannot be sequential numbers")
	}

	allSame := true
	for i := 1; i < len(pin); i++ {
		if pin[i] != pin[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return fmt.Errorf("PIN cannot be all the same digit")
	}

	return nil
}

func RecordLoginAttempt(cif, ip string, success bool) error {
	query := `INSERT INTO login_attempts (cif, ip_address, attempt_type, is_success) VALUES ($1, $2, 'pin', $3)`
	_, err := database.DB.Exec(query, cif, ip, success)
	return err
}

func IsAccountLocked(cif string) (bool, error) {
	maxAttempts := config.AppConfig.MaxLoginAttempts
	lockoutMinutes := config.AppConfig.LockoutDurationMinutes

	windowStart := time.Now().UTC().Add(-time.Duration(lockoutMinutes) * time.Minute)

	query := `SELECT COUNT(*) FROM login_attempts
	          WHERE cif = $1 AND is_success = FALSE AND attempted_at > $2
	          AND attempted_at > COALESCE(
	              (SELECT MAX(attempted_at) FROM login_attempts WHERE cif = $3 AND is_success = TRUE),
	              '2000-01-01'::TIMESTAMPTZ
	          )`
	var failedCount int
	err := database.DB.QueryRow(query, cif, windowStart, cif).Scan(&failedCount)
	if err != nil {
		return false, err
	}

	return failedCount >= maxAttempts, nil
}

func GetFailedAttemptCount(cif string) (int, error) {
	lockoutMinutes := config.AppConfig.LockoutDurationMinutes
	windowStart := time.Now().UTC().Add(-time.Duration(lockoutMinutes) * time.Minute)

	query := `SELECT COUNT(*) FROM login_attempts
	          WHERE cif = $1 AND is_success = FALSE AND attempted_at > $2
	          AND attempted_at > COALESCE(
	              (SELECT MAX(attempted_at) FROM login_attempts WHERE cif = $3 AND is_success = TRUE),
	              '2000-01-01'::TIMESTAMPTZ
	          )`
	var count int
	err := database.DB.QueryRow(query, cif, windowStart, cif).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func UnlockAccount(cif string) error {
	query := `INSERT INTO login_attempts (cif, ip_address, attempt_type, is_success) VALUES ($1, 'system', 'unlock', TRUE)`
	_, err := database.DB.Exec(query, cif)
	return err
}
