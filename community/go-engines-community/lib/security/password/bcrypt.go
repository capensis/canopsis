package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func NewBcryptEncoder() Encoder {
	return &bcryptEncoder{}
}

type bcryptEncoder struct{}

func (e *bcryptEncoder) EncodePassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func (e *bcryptEncoder) IsValidPassword(hashedPassword, password []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	return false, err
}
