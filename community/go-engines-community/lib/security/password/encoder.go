// Package password contains password encoders.
package password

//go:generate mockgen -destination=../../../mocks/lib/security/password/password.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password Encoder

// Encoder is used to implement password encoder.
type Encoder interface {
	// EncodePassword encodes the raw password.
	EncodePassword(password []byte) ([]byte, error)
	// IsValidPassword checks a raw password against an encoded password.
	IsValidPassword(encodedPassword, password []byte) (bool, error)
}
