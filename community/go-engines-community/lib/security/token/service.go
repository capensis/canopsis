package token

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type Service interface {
	GenerateToken(id string) (string, time.Time, error)
	ValidateToken(token string) (id string, err error)
}

func NewJwtService(
	secretKey []byte,
	issuer string,
	expirationInterval time.Duration,
) Service {
	return &jwtService{
		secretKey:          secretKey,
		issuer:             issuer,
		expirationInterval: expirationInterval,
	}
}

type jwtService struct {
	secretKey          []byte
	issuer             string
	expirationInterval time.Duration
}

type tokenClaims struct {
	ID string `json:"_id"`
	jwt.StandardClaims
}

func (s *jwtService) GenerateToken(id string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.expirationInterval)
	claims := tokenClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    s.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("cannot generate token: %w", err)
	}

	return t, expiresAt, nil
}

func (s *jwtService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.secretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("cannot parse token : %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token is invalid")
	}

	if claims, ok := token.Claims.(*tokenClaims); ok {
		return claims.ID, nil
	}

	return "", fmt.Errorf("token claims are invalid")
}
