package token

import (
	"time"
)

type Generator interface {
	Generate(id string, t time.Time) (string, error)
	Validate(token string) (id string, err error)
}
