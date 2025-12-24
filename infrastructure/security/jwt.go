package security

import (
	"fmt"
	"time"

	"portal-data-backend/infrastructure/config"
	"portal-data-backend/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager handles JWT token operations
type JWTManager struct {
	secret             []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	issuer             string
}

// Claims represents JWT claims
type Claims struct {
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
	RoleID         string `json:"role_id"`
	Email          string `json:"email"`
	jwt.RegisteredClaims
}

// TokenPair contains access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(cfg *config.JWTConfig) *JWTManager {
	return &JWTManager{
		secret:             []byte(cfg.Secret),
		accessTokenExpiry:  cfg.AccessTokenExpiry,
		refreshTokenExpiry: cfg.RefreshTokenExpiry,
		issuer:             cfg.Issuer,
	}
}

// GenerateTokenPair generates access and refresh tokens
func (j *JWTManager) GenerateTokenPair(userID, organizationID, roleID, email string) (*TokenPair, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}

	accessToken, err := j.generateAccessToken(userID, organizationID, roleID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := j.generateRefreshToken(userID, organizationID, roleID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(j.accessTokenExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// generateAccessToken generates an access token
func (j *JWTManager) generateAccessToken(userID, organizationID, roleID, email string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(j.accessTokenExpiry)

	claims := Claims{
		UserID:         userID,
		OrganizationID: organizationID,
		RoleID:         roleID,
		Email:          email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// generateRefreshToken generates a refresh token
func (j *JWTManager) generateRefreshToken(userID, organizationID, roleID, email string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(j.refreshTokenExpiry)

	claims := Claims{
		UserID:         userID,
		OrganizationID: organizationID,
		RoleID:         roleID,
		Email:          email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.ErrTokenExpired
		}
		return nil, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.ErrInvalidToken
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token from a valid refresh token
func (j *JWTManager) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	return j.generateAccessToken(claims.UserID, claims.OrganizationID, claims.RoleID, claims.Email)
}

// GetUserIDFromToken extracts user ID from token
func (j *JWTManager) GetUserIDFromToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
