package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	accessSecret   string
	refreshSecret  string
	accessExpires  time.Duration
	refreshExpires time.Duration
}

func NewJWTManager(accessSecret, refreshSecret string, accessExp, refreshExp time.Duration) *JWTManager {
	return &JWTManager{
		accessSecret:   accessSecret,
		refreshSecret:  refreshSecret,
		accessExpires:  accessExp,
		refreshExpires: refreshExp,
	}
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func (m *JWTManager) GenerateTokens(userID uuid.UUID, role string) (accessToken, refreshToken string, err error) {
	now := time.Now()

	// Access token
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessExpires)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})
	accessToken, err = access.SignedString([]byte(m.accessSecret))
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshExpires)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})
	refreshToken, err = refresh.SignedString([]byte(m.refreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (m *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	return m.validateToken(tokenString, m.accessSecret)
}

func (m *JWTManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return m.validateToken(tokenString, m.refreshSecret)
}

func (m *JWTManager) validateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
