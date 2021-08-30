package fixtures

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"github.com/brianvoe/gofakeit/v6"
)

type Faker struct {
	*gofakeit.Faker
	passwordEncoder password.Encoder
}

func (Faker) NowUnix() interface{} {
	return time.Now().Unix()
}

func (f Faker) DateUnix() interface{} {
	return time.Date(f.Number(1970, time.Now().Year()), time.Month(f.Month()), f.Day(), f.Hour(), f.Minute(), f.Second(), 0, time.UTC).Unix()
}

func (f Faker) Password(password string) string {
	return string(f.passwordEncoder.EncodePassword([]byte(password)))
}
