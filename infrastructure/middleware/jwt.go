package middleware

import (
	"errors"
	"qlass-be/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// IJwtService defines the interface for mocking in tests
type JwtService interface {
	GenerateToken(uniID string, role string) (string, error)
	ValidateToken(tokenString string) (*JWTCustomClaims, error)
}

type jwtService struct {
	secretKey []byte
}

type JWTCustomClaims struct {
	UniversityID string `json:"university_id"`
	Role         string `json:"role"`
	jwt.RegisteredClaims
}

// NewJWTService works like ConnectDB or ConnectRedis
func NewJWTService(cfg *config.Config) JwtService {
	return &jwtService{
		secretKey: []byte(cfg.JWTSecret),
	}
}

func (s *jwtService) GenerateToken(uniID string, role string) (string, error) {
	claims := &JWTCustomClaims{
		uniID,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *jwtService) ValidateToken(tokenString string) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}