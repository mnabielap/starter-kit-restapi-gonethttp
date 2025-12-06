package utils

import (
	"fmt"
	"time"

	"starter-kit-restapi-gonethttp/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPayload struct {
	Sub  string `json:"sub"` // User ID
	Type string `json:"type"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT token
func GenerateToken(userID uuid.UUID, expires time.Duration, tokenType string, secret string) (string, time.Time, error) {
	expirationTime := time.Now().Add(expires)

	claims := &TokenPayload{
		Sub:  userID.String(),
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expirationTime, nil
}

// ValidateToken parses and verifies a JWT token
func ValidateToken(tokenString string, secret string) (*TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenPayload); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GenerateAuthTokens generates both Access and Refresh tokens
func GenerateAuthTokens(userID uuid.UUID, cfg *config.Config) (string, string, time.Time, time.Time, error) {
	accessTokenExpires := time.Duration(cfg.JWT.AccessExpirationMinutes) * time.Minute
	accessToken, accessExp, err := GenerateToken(userID, accessTokenExpires, "access", cfg.JWT.Secret)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	refreshTokenExpires := time.Duration(cfg.JWT.RefreshExpirationDays) * 24 * time.Hour
	refreshToken, refreshExp, err := GenerateToken(userID, refreshTokenExpires, "refresh", cfg.JWT.Secret)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	return accessToken, refreshToken, accessExp, refreshExp, nil
}