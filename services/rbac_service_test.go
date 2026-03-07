package services_test

import (
	"testing"

	"cbs-simulator/services"
)

func TestGetUserRoles(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	roles, err := services.GetUserRoles("CIF001")
	if err != nil {
		t.Fatalf("Failed to get user roles: %v", err)
	}

	if len(roles) < 2 {
		t.Errorf("CIF001 should have at least 2 roles (customer + admin), got %d", len(roles))
	}
}

func TestHasRole_Admin(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	hasAdmin, err := services.HasRole("CIF001", "admin")
	if err != nil {
		t.Fatalf("Failed to check role: %v", err)
	}
	if !hasAdmin {
		t.Error("CIF001 should have admin role")
	}
}

func TestHasRole_CustomerOnly(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	hasAdmin, _ := services.HasRole("CIF002", "admin")
	if hasAdmin {
		t.Error("CIF002 should NOT have admin role")
	}

	hasCustomer, _ := services.HasRole("CIF002", "customer")
	if !hasCustomer {
		t.Error("CIF002 should have customer role")
	}
}

func TestGetPrimaryRole_Admin(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	role, err := services.GetPrimaryRole("CIF001")
	if err != nil {
		t.Fatalf("Failed to get primary role: %v", err)
	}
	if role != "admin" {
		t.Errorf("CIF001 primary role should be 'admin', got '%s'", role)
	}
}

func TestGetPrimaryRole_Supervisor(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	role, err := services.GetPrimaryRole("CIF003")
	if err != nil {
		t.Fatalf("Failed to get primary role: %v", err)
	}
	if role != "supervisor" {
		t.Errorf("CIF003 primary role should be 'supervisor', got '%s'", role)
	}
}

func TestGetPrimaryRole_Customer(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	role, err := services.GetPrimaryRole("CIF002")
	if err != nil {
		t.Fatalf("Failed to get primary role: %v", err)
	}
	if role != "customer" {
		t.Errorf("CIF002 primary role should be 'customer', got '%s'", role)
	}
}

func TestAssignRole(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	// CIF002 should not have teller role initially
	hasTeller, _ := services.HasRole("CIF002", "teller")
	if hasTeller {
		t.Error("CIF002 should NOT have teller initially")
	}

	// Assign teller role
	err := services.AssignRole("CIF002", 2, "CIF001") // role_id 2 = teller
	if err != nil {
		t.Fatalf("Failed to assign role: %v", err)
	}

	// Now should have teller
	hasTeller, _ = services.HasRole("CIF002", "teller")
	if !hasTeller {
		t.Error("CIF002 should have teller role after assignment")
	}
}

func TestGetAllRoles(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	roles, err := services.GetAllRoles()
	if err != nil {
		t.Fatalf("Failed to get all roles: %v", err)
	}
	if len(roles) != 4 {
		t.Errorf("Should have 4 roles, got %d", len(roles))
	}
}

func TestGetRoleByName(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	role, err := services.GetRoleByName("admin")
	if err != nil {
		t.Fatalf("Failed to get role: %v", err)
	}
	if role.RoleName != "admin" {
		t.Errorf("Role name should be 'admin', got '%s'", role.RoleName)
	}
}

func TestGetRoleByName_NotFound(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	_, err := services.GetRoleByName("nonexistent")
	if err == nil {
		t.Error("Should fail for nonexistent role")
	}
}
