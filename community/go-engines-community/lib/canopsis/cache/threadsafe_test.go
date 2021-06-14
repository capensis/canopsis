package cache_test

import (
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/cache"
	"git.canopsis.net/canopsis/go-engines/lib/testutils"
	. "github.com/smartystreets/goconvey/convey"
)

type TestCacheVal int

func (t TestCacheVal) CacheID() string {
	return "exists"
}

func NewTestCache() cache.Cache {
	kv := cache.NewKV()
	return cache.NewThreadSafeCache(kv)
}

func TestThreadSafeCache(t *testing.T) {
	Convey("Setup", t, func() {
		tsc := NewTestCache()

		v := TestCacheVal(42)
		err := tsc.Set(v)
		So(err, ShouldBeNil)

		var cv TestCacheVal

		exists := tsc.Get("exists", &cv)
		So(exists, ShouldBeTrue)
		So(int(cv), ShouldEqual, 42)

		exists = tsc.Get("noexists", &cv)
		So(exists, ShouldBeFalse)

	})
}

func TestCacheConcurrency(t *testing.T) {
	testutils.SkipLongIfSet(t)
	Convey("Setup", t, func() {
		tsc := NewTestCache()

		cacheWork := func() {
			run := 30000000

			tcv := TestCacheVal(1)
			tsc.Set(tcv)

			finished := make(chan bool, 2)

			cacheRead := func() {
				for i := 1; i < run; i++ {
					var val TestCacheVal
					exists := tsc.Get(tcv.CacheID(), &val)
					if int(val) == 0 || exists == false {
						t.Fatal("got 0 value")
					}
				}
				finished <- true
			}

			cacheWrite := func() {
				for i := 1; i < run; i++ {
					tsc.Set(TestCacheVal(i))
				}
				finished <- true
			}

			go cacheRead()
			go cacheWrite()
			<-finished
			<-finished
		}

		So(cacheWork, ShouldNotPanic)

	})
}
