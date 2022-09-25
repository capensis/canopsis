package fixtures

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/brianvoe/gofakeit/v6"
)

type Faker struct {
	*gofakeit.Faker
	passwordEncoder password.Encoder

	usedNames           map[string]struct{}
	pbehaviorExceptions map[string]struct{}
}

func NewFaker(passwordEncoder password.Encoder) *Faker {
	return &Faker{
		Faker:               gofakeit.New(0),
		passwordEncoder:     passwordEncoder,
		usedNames:           make(map[string]struct{}),
		pbehaviorExceptions: make(map[string]struct{}),
	}
}

func (*Faker) NowUnix() interface{} {
	return time.Now().Unix()
}

func (*Faker) NowUnixAdd(dStr string) (interface{}, error) {
	d, err := time.ParseDuration(dStr)
	if err != nil {
		return nil, err
	}

	return types.CpsTime{Time: time.Now().Add(d)}, nil
}

func (*Faker) GenerateExdates(count int) interface{} {
	exdates := make([]types.Exdate, count)
	now := time.Now()

	leftBound := now.AddDate(-1, 0, 0).Unix()
	upperBound := now.AddDate(1, 0, 0).Unix()
	interval := upperBound - leftBound

	for idx := range exdates {
		begin := rand.Int63n(interval)
		exdates[idx].Begin = types.CpsTime{Time: time.Unix(leftBound+begin, 0)}
		end := rand.Int63n(interval - begin)
		exdates[idx].End = types.CpsTime{Time: time.Unix(exdates[idx].Begin.Unix()+end, 0)}
	}

	return exdates
}

func (f *Faker) GeneratePBehaviorExceptionID() string {
	id := utils.NewID()
	f.pbehaviorExceptions[id] = struct{}{}

	return id
}

func (f *Faker) LinkPBehaviorExceptions() []string {
	exceptions := make([]string, 0, len(f.pbehaviorExceptions))
	for k := range f.pbehaviorExceptions {
		exceptions = append(exceptions, k)
	}

	return exceptions
}

func (f *Faker) DateUnix() interface{} {
	return f.Number(1, int(time.Now().Unix()))
}

func (f *Faker) Password(password string) string {
	return string(f.passwordEncoder.EncodePassword([]byte(password)))
}

func (f *Faker) UniqueName() (string, error) {
	for nameLen := 5; nameLen < 11; nameLen++ {
		for try := 0; try < 3; try++ {
			v := f.Generate(strings.Repeat("?", nameLen))
			if _, ok := f.usedNames[v]; !ok {
				f.usedNames[v] = struct{}{}
				return v, nil
			}
		}
	}

	return "", errors.New("cannot generate unique name")
}

func (f *Faker) ResetUniqueName() {
	f.usedNames = make(map[string]struct{})
}

func (f *Faker) ResetPBehaviorExceptions() {
	f.pbehaviorExceptions = make(map[string]struct{})
}
