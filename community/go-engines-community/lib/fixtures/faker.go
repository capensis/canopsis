package fixtures

import (
	"errors"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Faker struct {
	*gofakeit.Faker
	passwordEncoder password.Encoder

	usedNames map[string]struct{}
}

func NewFaker(passwordEncoder password.Encoder) *Faker {
	return &Faker{
		Faker:           gofakeit.New(0),
		passwordEncoder: passwordEncoder,
		usedNames:       make(map[string]struct{}),
	}
}

func (*Faker) NowUnix() any {
	return time.Now().Unix()
}

func (*Faker) NowUnixAdd(dStr string) (any, error) {
	d, err := time.ParseDuration(dStr)
	if err != nil {
		return nil, err
	}

	return time.Now().Add(d).Unix(), nil
}

func (f *Faker) DateUnix() any {
	return f.Number(1, int(time.Now().Unix()))
}

func (f *Faker) Password(password string) (string, error) {
	h, err := f.passwordEncoder.EncodePassword([]byte(password))
	if err != nil {
		return "", err
	}

	return string(h), nil
}

func (f *Faker) UniqueName() (string, error) {
	for nameLen := 5; nameLen < 11; nameLen++ {
		for try := 0; try < 3; try++ {
			v, err := f.Generate(strings.Repeat("?", nameLen))
			if err != nil {
				return "", err
			}

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

func (f *Faker) UUID() string {
	return utils.NewID()
}

func (f *Faker) JWT() (string, error) {
	registeredClaims := jwt.RegisteredClaims{
		ID:       utils.NewID(),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Issuer:   canopsis.AppName,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims).SignedString([]byte(""))
}

// GenerateExdates
//
//nolint:gosec
func (*Faker) GenerateExdates(count int) any {
	exdates := make([]types.Exdate, count)
	now := time.Now()

	leftBound := now.AddDate(-1, 0, 0).Unix()
	upperBound := now.AddDate(1, 0, 0).Unix()
	interval := upperBound - leftBound

	for idx := range exdates {
		begin := rand.Int63n(interval)
		exdates[idx].Begin = datetime.CpsTime{Time: time.Unix(leftBound+begin, 0)}
		end := rand.Int63n(interval - begin)
		exdates[idx].End = datetime.CpsTime{Time: time.Unix(exdates[idx].Begin.Unix()+end, 0)}
	}

	return exdates
}

func (*Faker) GenerateBookmarks(prefix string, count int) []string {
	bookmarks := make([]string, count)

	for idx := range bookmarks {
		bookmarks[idx] = prefix + strconv.Itoa(idx)
	}

	return bookmarks
}

// ToObjectID converts a string of hex digits to a MongoDB ObjectId
func (*Faker) ToObjectID(s string) any {
	s = utils.ToObjectIDHex(s)
	objectID, err := bson.ObjectIDFromHex(s)
	if err != nil {
		return bson.NilObjectID
	}
	return objectID
}

func (*Faker) ObjectID() any {
	return bson.NewObjectID()
}

func (*Faker) AbsPath(path string) (string, error) {
	return filepath.Abs(path)
}
