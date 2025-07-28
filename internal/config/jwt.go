package config

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTConfig holds JWT-specific configuration
type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	AccessTokenTTL  int    `mapstructure:"access_token_ttl"`
	RefreshTokenTTL int    `mapstructure:"refresh_token_ttl"`
	Issuer          string `mapstructure:"issuer"`
}

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
	GetSugaredLogger().Info("🔄 Initializing JWT manager...")

	// Get JWT configuration
	jwtConfig := cfg.JWT

	// Validate JWT secret
	if jwtConfig.Secret == "" || jwtConfig.Secret == "your-secret-key" {
		if IsProduction() {
			GetSugaredLogger().Error("JWT secret must be set in production environment")
			return fmt.Errorf("JWT secret must be set in production environment")
		}
		GetSugaredLogger().Warn("⚠️  Warning: Using default JWT secret. Change this in production!")
	}

	jwtManager = &JWTManager{
		secretKey:       []byte(jwtConfig.Secret),
		accessTokenTTL:  time.Duration(jwtConfig.AccessTokenTTL) * time.Second,
		refreshTokenTTL: time.Duration(jwtConfig.RefreshTokenTTL) * time.Second,
		issuer:          jwtConfig.Issuer,
	}

	GetSugaredLogger().Infof("✅ JWT manager initialized with TTL: access=%v, refresh=%v",
		jwtManager.accessTokenTTL, jwtManager.refreshTokenTTL)

	return nil
}

// GetJWTManager returns the JWT manager
func GetJWTManager() *JWTManager {
	if jwtManager == nil {
		GetSugaredLogger().Fatal("JWT manager not initialized. Call InitJWT() first")
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
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		GetSugaredLogger().Errorf("Failed to generate access token for user %d: %v", userID, err)
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	GetSugaredLogger().Debugf("Generated access token for user %d (%s)", userID, username)
	return tokenString, nil
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
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		GetSugaredLogger().Errorf("Failed to generate refresh token for user %d: %v", userID, err)
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	GetSugaredLogger().Debugf("Generated refresh token for user %d", userID)
	return tokenString, nil
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
		GetSugaredLogger().Debugf("Failed to parse token: %v", err)
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		GetSugaredLogger().Debugf("Token validated for user %d (%s)", claims.UserID, claims.Username)
		return claims, nil
	}

	GetSugaredLogger().Debug("Invalid token provided")
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
