package token

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(token string) string {
	sum := sha256.Sum256([]byte(token))

	return hex.EncodeToString(sum[:])
}
