package jwt

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/theotruvelot/catchook/internal/config"
)

type Manager interface {
	GenerateAccessToken(userID int, role string) (string, error)
	GenerateRefreshToken(userID int) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	ParseRefreshToken(refreshToken string) (*Claims, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type Claims struct {
	UserID    int    `json:"user_id"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"` // "access" ou "refresh"
	jwt.RegisteredClaims
}

type jwtManager struct {
	config config.JWTConfig
}

func NewManager(config config.JWTConfig) Manager {
	return &jwtManager{
		config: config,
	}
}

func (j *jwtManager) GenerateAccessToken(userID int, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			Subject:   strconv.Itoa(userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.config.AccessTokenDuration)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// Debug log pour voir les claims
	fmt.Printf("Generating token with claims: %+v\n", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	// Debug log pour voir le token généré
	fmt.Printf("Generated token: %s\n", tokenString)

	return tokenString, nil
}

func (j *jwtManager) GenerateRefreshToken(userID int) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			Subject:   strconv.Itoa(userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.config.RefreshTokenDuration)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

func (j *jwtManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (j *jwtManager) ParseRefreshToken(refreshToken string) (*Claims, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	return claims, nil
}

func (j *jwtManager) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := j.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	accessToken, err := j.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	return accessToken, nil
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}
