package activecheck_test

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/activecheck"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	. "github.com/smartystreets/goconvey/convey"
)

type redLockTest struct {
	activecheck.RedLocker

	active bool
}

func (rl *redLockTest) Lock(resource string, period int64) (string, error) {
	var err error

	if !rl.active {
		err = fmt.Errorf("acquire lock error")
	}
	return getRandStr(), err
}

func (rl *redLockTest) ExpireLock(_, _ string, _ int64) (int64, error) {
	return 0, nil
}

func (rl *redLockTest) setAcive(active bool) {
	rl.active = active
}

func NewRedLockTest() (*redLockTest, error) {
	return &redLockTest{
		active: true,
	}, nil
}

func getRandStr() string {
	b := make([]byte, 16)
	crand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func TestActiveCheckStart(t *testing.T) {
	Convey("Initialize", t, func() {
		ac := make([]struct {
			check  activecheck.ActiveChecker
			rl     *redLockTest
			active bool
		}, 2)

		ac[0].check = activecheck.NewActiveCheck("testLockKey")
		So(ac[0].check, ShouldNotBeNil)

		ac[1].check = activecheck.NewActiveCheck("testLockKey")
		So(ac[1].check, ShouldNotBeNil)

		Convey("Start the 1st node", func() {
			var err error
			checkPeriod := 6 * time.Second

			_, err = activecheck.NewRedLock([]string{""})
			So(err, ShouldNotBeNil)

			testAddr := activecheck.AddressListFlags{}
			testAddr.Set(os.Getenv(redis.EnvURL))

			rl, err := activecheck.NewRedLock(testAddr)
			So(err, ShouldBeNil)

			rl.SetRetryCount(5)
			rl.SetRetryDelay(rand.Intn(int((checkPeriod).Nanoseconds()) / 2 / 1e6))

			ac[0].rl, err = NewRedLockTest()

			ac[1].rl, err = NewRedLockTest()
			So(err, ShouldBeNil)

			dt := time.Now()
			ac[0].active, err = ac[0].check.Start(activecheck.RedLocker(ac[0].rl), checkPeriod)
			So(err, ShouldBeNil)

			if ac[0].active {

				Convey("keepalive 1st node", func() {
					err = ac[0].check.Keepalive()
					So(err, ShouldBeNil)

					Convey("Start the 2nd node", func() {
						ac[1].rl.setAcive(false)
						ac[1].active, err = ac[1].check.Start(activecheck.RedLocker(ac[1].rl), checkPeriod)

						So(ac[0].check.GetValue(), ShouldNotEqual, ac[1].check.GetValue())

						So(ac[1].active, ShouldBeFalse)

						So(err, ShouldBeNil)
						activeNode, passiveNode := 0, 0

						err = ac[0].check.Keepalive()
						if err != nil {
							ac[0].check.SetPassive()
							ac[0].active = ac[0].check.IsActive()
							So(time.Now(), ShouldNotHappenWithin, checkPeriod, dt)

							activeNode = 1
						} else {
							passiveNode = 1
						}

						So(ac[activeNode].active, ShouldBeTrue)
						So(ac[passiveNode].active, ShouldBeFalse)

						Convey("Switch active node", func() {
							ac[activeNode].check.SetPassive()
							ac[activeNode].active = ac[activeNode].check.IsActive()
							So(ac[activeNode].active, ShouldBeFalse)

							Convey("Wait active node switch", func() {
								// for real lock must sleep until expiration
								ac[passiveNode].rl.setAcive(true)
								ac[passiveNode].active, err = ac[passiveNode].check.Start(activecheck.RedLocker(ac[passiveNode].rl), checkPeriod)
								So(err, ShouldBeNil)

								So(ac[passiveNode].active, ShouldBeTrue)
							})
						})

					})
				})
			}
		})
	})
}
