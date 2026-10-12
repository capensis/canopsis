package token

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

const (
	prefix        = "cps_"
	maskedStars   = 10
	maskedVisible = 5
)

func Generate() (string, error) {
	b := make([]byte, 32) // 256 bits
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return prefix + base64.RawURLEncoding.EncodeToString(b), nil
}

func Mask(token string) string {
	if len(token) <= maskedVisible {
		return strings.Repeat("*", maskedStars)
	}

	return strings.Repeat("*", maskedStars) + token[len(token)-maskedVisible:]
}
