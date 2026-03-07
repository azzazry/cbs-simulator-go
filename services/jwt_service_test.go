package services_test

import (
	"testing"
	"time"

	"cbs-simulator/services"
)

func TestGenerateTokenPair(t *testing.T) {
	ensureConfig()

	pair, err := services.GenerateTokenPair("CIF001", "admin")
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	if pair.AccessToken == "" {
		t.Error("Access token should not be empty")
	}
	if pair.RefreshToken == "" {
		t.Error("Refresh token should not be empty")
	}
	if pair.TokenType != "Bearer" {
		t.Errorf("Token type should be 'Bearer', got '%s'", pair.TokenType)
	}
	if pair.ExpiresIn <= 0 {
		t.Error("ExpiresIn should be positive")
	}
}

func TestValidateToken_Valid(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	pair, _ := services.GenerateTokenPair("CIF001", "admin")

	claims, err := services.ValidateToken(pair.AccessToken)
	if err != nil {
		t.Fatalf("Token should be valid: %v", err)
	}

	if claims.CIF != "CIF001" {
		t.Errorf("CIF should be 'CIF001', got '%s'", claims.CIF)
	}
	if claims.Role != "admin" {
		t.Errorf("Role should be 'admin', got '%s'", claims.Role)
	}
	if claims.TokenType != "access" {
		t.Errorf("Token type should be 'access', got '%s'", claims.TokenType)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	ensureConfig()

	_, err := services.ValidateToken("invalid-token-string")
	if err == nil {
		t.Error("Should fail with invalid token")
	}
}

func TestValidateToken_BlacklistedToken(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	pair, _ := services.GenerateTokenPair("CIF001", "admin")

	// Validate to get JTI
	claims, _ := services.ValidateToken(pair.AccessToken)

	// Blacklist it
	err := services.BlacklistToken(claims.ID, "CIF001", time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("Failed to blacklist token: %v", err)
	}

	// Should now fail
	_, err = services.ValidateToken(pair.AccessToken)
	if err == nil {
		t.Error("Blacklisted token should be rejected")
	}
}

func TestRefreshAccessToken(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	pair, _ := services.GenerateTokenPair("CIF001", "admin")

	newPair, err := services.RefreshAccessToken(pair.RefreshToken)
	if err != nil {
		t.Fatalf("Failed to refresh token: %v", err)
	}

	if newPair.AccessToken == "" {
		t.Error("New access token should not be empty")
	}
	if newPair.AccessToken == pair.AccessToken {
		t.Error("New access token should be different from old one")
	}
}

func TestRefreshAccessToken_WithAccessToken(t *testing.T) {
	ensureConfig()
	setupTestDB(t)
	defer teardownTestDB(t)

	pair, _ := services.GenerateTokenPair("CIF001", "admin")

	// Using access token instead of refresh token should fail
	_, err := services.RefreshAccessToken(pair.AccessToken)
	if err == nil {
		t.Error("Should fail when using access token to refresh")
	}
}
