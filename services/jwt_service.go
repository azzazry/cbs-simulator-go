package services

import (
	"fmt"
	"time"

	"cbs-simulator/config"
	"cbs-simulator/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims represents custom JWT claims
type JWTClaims struct {
	CIF       string `json:"cif"`
	Role      string `json:"role"`
	TokenType string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds until access token expires
	TokenType    string `json:"token_type"`
}

// GenerateTokenPair creates a new access + refresh token pair
func GenerateTokenPair(cif, role string) (*TokenPair, error) {
	secret := []byte(config.AppConfig.JWTSecret)

	// Access token
	accessJTI := uuid.New().String()
	accessExpiry := time.Now().Add(time.Duration(config.AppConfig.JWTAccessExpiry) * time.Minute)
	accessClaims := JWTClaims{
		CIF:       cif,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        accessJTI,
			Issuer:    "cbs-simulator",
			Subject:   cif,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %v", err)
	}

	// Refresh token
	refreshJTI := uuid.New().String()
	refreshExpiry := time.Now().Add(time.Duration(config.AppConfig.JWTRefreshExpiry) * time.Hour)
	refreshClaims := JWTClaims{
		CIF:       cif,
		Role:      role,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        refreshJTI,
			Issuer:    "cbs-simulator",
			Subject:   cif,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %v", err)
	}

	expiresIn := int64(config.AppConfig.JWTAccessExpiry * 60) // convert minutes to seconds

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken validates a JWT token and returns claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
	secret := []byte(config.AppConfig.JWTSecret)

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Check if token is blacklisted
	blacklisted, err := IsTokenBlacklisted(claims.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %v", err)
	}
	if blacklisted {
		return nil, fmt.Errorf("token has been revoked")
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token from a valid refresh token
func RefreshAccessToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := ValidateToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %v", err)
	}

	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	// Generate new token pair
	return GenerateTokenPair(claims.CIF, claims.Role)
}

// BlacklistToken adds a token to the blacklist (for logout)
func BlacklistToken(jti, cif string, expiresAt time.Time) error {
	query := `INSERT OR IGNORE INTO token_blacklist (token_jti, cif, expires_at) VALUES (?, ?, ?)`
	_, err := database.DB.Exec(query, jti, cif, expiresAt)
	return err
}

// IsTokenBlacklisted checks if a token JTI is in the blacklist
func IsTokenBlacklisted(jti string) (bool, error) {
	var count int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM token_blacklist WHERE token_jti = ?`, jti).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CleanupExpiredBlacklistTokens removes expired tokens from the blacklist
func CleanupExpiredBlacklistTokens() error {
	_, err := database.DB.Exec(`DELETE FROM token_blacklist WHERE expires_at < ?`, time.Now())
	return err
}
