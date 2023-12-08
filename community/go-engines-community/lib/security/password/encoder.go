// password contains password encoders.
package password

//go:generate mockgen -destination=../../../mocks/lib/security/password/password.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password Encoder

import (
	"bytes"
	"crypto/sha1" //nolint:gosec
	"encoding/hex"
)

// Encoder is used to implement password encoder.
type Encoder interface {
	// EncodePassword encodes the raw password.
	EncodePassword(password []byte) []byte
	// IsValidPassword checks a raw password against an encoded password.
	IsValidPassword(encodedPassword, password []byte) bool
}

type sha1Encoder struct{}

// NewSha1Encoder creates new encoder.
func NewSha1Encoder() Encoder {
	return &sha1Encoder{}
}

func (e *sha1Encoder) EncodePassword(password []byte) []byte {
	h := sha1.New() //nolint:gosec
	_, err := h.Write(password)
	if err != nil {
		panic(err)
	}
	hash := h.Sum(nil)
	encodedPassword := make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(encodedPassword, hash)

	return encodedPassword
}

func (e *sha1Encoder) IsValidPassword(encodedPassword, password []byte) bool {
	return bytes.Equal(e.EncodePassword(password), encodedPassword)
}
