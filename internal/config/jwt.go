package config

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTManager struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	issuer          string
}

var jwtManager *JWTManager

// Claims represents JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// InitJWT initializes the JWT manager
func InitJWT() error {
	log.Println("🔄 Initializing JWT manager...")

	// Get JWT configuration
	jwtConfig := cfg.JWT

	// Validate JWT secret
	if jwtConfig.Secret == "" || jwtConfig.Secret == "your-secret-key" {
		if IsProduction() {
			return fmt.Errorf("JWT secret must be set in production environment")
		}
		log.Println("⚠️  Warning: Using default JWT secret. Change this in production!")
	}

	jwtManager = &JWTManager{
		secretKey:       []byte(jwtConfig.Secret),
		accessTokenTTL:  time.Duration(jwtConfig.AccessTokenTTL) * time.Second,
		refreshTokenTTL: time.Duration(jwtConfig.RefreshTokenTTL) * time.Second,
		issuer:          jwtConfig.Issuer,
	}

	log.Printf("✅ JWT manager initialized with TTL: access=%v, refresh=%v",
		jwtManager.accessTokenTTL, jwtManager.refreshTokenTTL)

	return nil
}

// GetJWTManager returns the JWT manager
func GetJWTManager() *JWTManager {
	if jwtManager == nil {
		log.Fatal("JWT manager not initialized. Call InitJWT() first")
	}
	return jwtManager
}

// GenerateAccessToken generates a new access token
func (j *JWTManager) GenerateAccessToken(userID uint, username, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// GenerateRefreshToken generates a new refresh token
func (j *JWTManager) GenerateRefreshToken(userID uint) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    j.issuer,
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshTokenTTL)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ValidateToken validates and parses a JWT token
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetTokenTTL returns the access token TTL
func (j *JWTManager) GetTokenTTL() time.Duration {
	return j.accessTokenTTL
}

// GetRefreshTokenTTL returns the refresh token TTL
func (j *JWTManager) GetRefreshTokenTTL() time.Duration {
	return j.refreshTokenTTL
}
