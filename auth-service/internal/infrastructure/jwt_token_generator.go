package infrastructure

import (
	"auth-service/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTTokenGenerator implementa el puerto TokenGenerator usando JWT
type JWTTokenGenerator struct {
	secretKey     []byte
	expirationMin int
}

// NewJWTTokenGenerator crea una nueva instancia del generador de tokens
func NewJWTTokenGenerator(secretKey string, expirationMin int) *JWTTokenGenerator {
	return &JWTTokenGenerator{
		secretKey:     []byte(secretKey),
		expirationMin: expirationMin,
	}
}

// Claims personalizados para el JWT
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken crea un token JWT con el ID del usuario
func (g *JWTTokenGenerator) GenerateToken(userID string) (*domain.AuthToken, error) {
	expirationTime := time.Now().Add(time.Duration(g.expirationMin) * time.Minute)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(g.secretKey)
	if err != nil {
		return nil, err
	}

	return &domain.AuthToken{
		Token:     tokenString,
		UserID:    userID,
		ExpiresAt: expirationTime.Unix(),
	}, nil
}

// ValidateToken valida un token JWT y retorna el ID del usuario
func (g *JWTTokenGenerator) ValidateToken(tokenString string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return g.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}
