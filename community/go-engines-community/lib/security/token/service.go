package token

import (
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/golang-jwt/jwt/v4"
)

type Service interface {
	GenerateToken(id string, t time.Time) (string, error)
	ValidateToken(token string) (id string, err error)
}

func NewJwtService(
	secretKey []byte,
	issuer string,
	apiConfigProvider config.ApiConfigProvider,
) Service {
	return &jwtService{
		secretKey:         secretKey,
		issuer:            issuer,
		apiConfigProvider: apiConfigProvider,
	}
}

type jwtService struct {
	secretKey         []byte
	issuer            string
	apiConfigProvider config.ApiConfigProvider
}

type tokenClaims struct {
	ID string `json:"_id"`
	jwt.RegisteredClaims
}

func (s *jwtService) GenerateToken(id string, expiresAt time.Time) (string, error) {
	cfg := s.apiConfigProvider.Get()
	registeredClaims := jwt.RegisteredClaims{
		ID:       utils.NewID(),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Issuer:   s.issuer,
	}
	if !expiresAt.IsZero() {
		registeredClaims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	}
	claims := tokenClaims{
		ID:               id,
		RegisteredClaims: registeredClaims,
	}
	token := jwt.NewWithClaims(cfg.TokenSigningMethod, claims)

	t, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("cannot generate token: %w", err)
	}

	return t, nil
}

func (s *jwtService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		cfg := s.apiConfigProvider.Get()

		if token.Method.Alg() != cfg.TokenSigningMethod.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %q, expected %q", token.Method.Alg(), cfg.TokenSigningMethod.Alg())
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
