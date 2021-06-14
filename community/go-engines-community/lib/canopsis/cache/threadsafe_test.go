package cache_test

import (
	"context"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/cache"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/testutils"
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
	ctx := context.Background()
	Convey("Setup", t, func() {
		tsc := NewTestCache()

		v := TestCacheVal(42)
		err := tsc.Set(ctx, v)
		So(err, ShouldBeNil)

		var cv TestCacheVal

		exists := tsc.Get(ctx, "exists", &cv)
		So(exists, ShouldBeTrue)
		So(int(cv), ShouldEqual, 42)

		exists = tsc.Get(ctx, "noexists", &cv)
		So(exists, ShouldBeFalse)

	})
}

func TestCacheConcurrency(t *testing.T) {
	ctx := context.Background()

	testutils.SkipLongIfSet(t)
	Convey("Setup", t, func() {
		tsc := NewTestCache()

		cacheWork := func() {
			run := 30000000

			tcv := TestCacheVal(1)
			tsc.Set(ctx, tcv)

			finished := make(chan bool, 2)

			cacheRead := func() {
				for i := 1; i < run; i++ {
					var val TestCacheVal
					exists := tsc.Get(ctx, tcv.CacheID(), &val)
					if int(val) == 0 || exists == false {
						t.Fatal("got 0 value")
					}
				}
				finished <- true
			}

			cacheWrite := func() {
				for i := 1; i < run; i++ {
					tsc.Set(ctx, TestCacheVal(i))
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
