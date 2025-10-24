package password

import (
	"bytes"
	"crypto/sha1" //nolint:gosec
	"encoding/hex"
)

func NewSha1Encoder() Encoder {
	return &sha1Encoder{}
}

type sha1Encoder struct{}

func (e *sha1Encoder) EncodePassword(password []byte) ([]byte, error) {
	h := sha1.New() //nolint:gosec
	_, err := h.Write(password)
	if err != nil {
		return nil, err
	}

	hash := h.Sum(nil)
	encodedPassword := make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(encodedPassword, hash)

	return encodedPassword, nil
}

func (e *sha1Encoder) IsValidPassword(storedEncodedPassword, password []byte) (bool, error) {
	encodedPassword, err := e.EncodePassword(password)
	if err != nil {
		return false, err
	}

	return bytes.Equal(encodedPassword, storedEncodedPassword), nil
}
