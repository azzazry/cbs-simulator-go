package services

import (
	"cbs-simulator/database"
	"cbs-simulator/models"
	"fmt"
)

// GetUserRoles returns all roles assigned to a CIF
func GetUserRoles(cif string) ([]string, error) {
	query := `SELECT r.role_name FROM user_roles ur 
	          JOIN roles r ON ur.role_id = r.id 
	          WHERE ur.cif = ? AND r.is_active = 1`

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
		roles = append(roles, "customer") // default role
	}

	return roles, nil
}

// GetPrimaryRole returns the highest-priority role for a CIF
func GetPrimaryRole(cif string) (string, error) {
	roles, err := GetUserRoles(cif)
	if err != nil {
		return "customer", err
	}

	// Priority: admin > supervisor > teller > customer
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

// HasRole checks if a CIF has a specific role
func HasRole(cif, roleName string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_roles ur 
	          JOIN roles r ON ur.role_id = r.id 
	          WHERE ur.cif = ? AND r.role_name = ? AND r.is_active = 1`
	err := database.DB.QueryRow(query, cif, roleName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// AssignRole assigns a role to a CIF
func AssignRole(cif string, roleID int, assignedBy string) error {
	query := `INSERT OR IGNORE INTO user_roles (cif, role_id, assigned_by) VALUES (?, ?, ?)`
	_, err := database.DB.Exec(query, cif, roleID, assignedBy)
	return err
}

// RemoveRole removes a role from a CIF
func RemoveRole(cif string, roleID int) error {
	query := `DELETE FROM user_roles WHERE cif = ? AND role_id = ?`
	_, err := database.DB.Exec(query, cif, roleID)
	return err
}

// GetAllRoles returns all available roles
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

// GetUserRoleDetails returns detailed role info for a CIF
func GetUserRoleDetails(cif string) ([]models.Role, error) {
	query := `SELECT r.id, r.role_name, r.description, r.is_active, r.created_at 
	          FROM user_roles ur JOIN roles r ON ur.role_id = r.id 
	          WHERE ur.cif = ? AND r.is_active = 1`
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

// GetRoleByName returns a role by name
func GetRoleByName(roleName string) (*models.Role, error) {
	var role models.Role
	query := `SELECT id, role_name, description, is_active, created_at FROM roles WHERE role_name = ?`
	err := database.DB.QueryRow(query, roleName).Scan(&role.ID, &role.RoleName, &role.Description, &role.IsActive, &role.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("role not found: %s", roleName)
	}
	return &role, nil
}
