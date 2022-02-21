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
	return f.Number(1, int(time.Now().Unix()))
}

func (f Faker) Password(password string) string {
	return string(f.passwordEncoder.EncodePassword([]byte(password)))
}
