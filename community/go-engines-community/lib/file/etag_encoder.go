package file

import (
	"crypto/sha1"
	"encoding/hex"
)

type sha1Encoder struct{}

func NewEtagEncoder() EtagEncoder {
	return &sha1Encoder{}
}

func (e *sha1Encoder) Encode(data []byte) (string, error) {
	h := sha1.New()
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	hash := h.Sum(nil)

	return hex.EncodeToString(hash), nil
}
