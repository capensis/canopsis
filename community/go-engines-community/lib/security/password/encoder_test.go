package password

import (
	"bytes"
	"testing"
)

func TestSha1Encoder_EncodePassword_GivenPassword_ShouldReturnSha1Hash(t *testing.T) {
	e := NewSha1Encoder()
	password := []byte("Test12345")
	encodedPassword := e.EncodePassword(password)
	expected := []byte("f988c245b3c789a608b34cd1b7c1b612542dbd09")

	if !bytes.Equal(encodedPassword, expected) {
		t.Errorf("expected encode: \"%v\", but got \"%v\"", string(expected), string(encodedPassword))
	}
}

func TestSha1Encoder_IsValidPassword_GivenEncodedPasswordAndPassword_ShouldReturnTrue(t *testing.T) {
	e := NewSha1Encoder()
	password := []byte("Test12345")
	encodedPassword := []byte("f988c245b3c789a608b34cd1b7c1b612542dbd09")

	if !e.IsValidPassword(encodedPassword, password) {
		t.Errorf("expected return true but got false")
	}
}

func TestSha1Encoder_IsValidPassword_GivenInvaludEncodedPasswordAndPassword_ShouldReturnFalse(t *testing.T) {
	e := NewSha1Encoder()
	password := []byte("Test12345")
	encodedPassword := []byte("invalidhash")

	if e.IsValidPassword(encodedPassword, password) {
		t.Errorf("expected return false but got true")
	}
}
