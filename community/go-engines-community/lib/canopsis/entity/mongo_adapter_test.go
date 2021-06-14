package entity_test

import (
	"fmt"
	"strconv"
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/errt"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/testutils"
	. "github.com/smartystreets/goconvey/convey"
)

func testNewAdapter() (entity.Adapter, error) {
	mongo, err := mongo.NewSession(mongo.Timeout)
	if err != nil {
		return nil, err
	}

	collection := entity.DefaultCollection(mongo)
	_, err = collection.RemoveAll()
	if err != nil {
		return nil, err
	}

	return entity.NewAdapter(collection), nil
}

func TestEntityDBAdapter(t *testing.T) {
	Convey("I can get an adapter", t, func() {
		ad, err := testNewAdapter()
		So(err, ShouldBeNil)

		Convey("And check for an nonexistent entity", func() {
			_, err := ad.GetEntityByID("this is not this entity you are searching for")
			_, ok := err.(errt.NotFound)
			So(ok, ShouldBeTrue)
		})

		Convey("Set an entity in database", func() {
			ent := types.Entity{
				Enabled:       true,
				ID:            "squeeky/squeeke",
				Type:          "watcher",
				Depends:       []string{},
				EnableHistory: make([]types.CpsTime, 0),
				Impacts:       []string{},
			}

			So(ad.Insert(ent), ShouldBeNil)

			Convey("Get it back", func() {
				dbentity, err := ad.GetEntityByID("squeeky/squeeke")

				So(err, ShouldBeNil)
				So(dbentity, ShouldNotBeNil)
				So(dbentity.ID, ShouldEqual, "squeeky/squeeke")
				So(dbentity.Type, ShouldEqual, "watcher")
				So(dbentity.Enabled, ShouldEqual, true)
			})
		})
	})
}

// Verify that concurrency does not affect bulk operations
func TestBulkConcurrency(t *testing.T) {
	t.SkipNow()
	testutils.SkipLongIfSet(t)
	Convey("Given some bulk inserts in two goroutines", t, func() {
		ed, err := testNewAdapter()
		So(err, ShouldBeNil)

		id := ""
		name := "Spamalot"

		finished := make(chan bool, 2)
		start := 0
		mid := 40000
		high := 80000

		writeEven := func() {
			entity := types.NewEntity(
				id,
				name,
				types.EntityTypeComponent,
				nil, nil, nil,
			)
			for index := start; index < mid; index++ {
				entity.ID = strconv.Itoa(index)
				entity.Name = name + "_LOW_" + strconv.Itoa(index)
				ed.BulkInsert(entity)
			}
			finished <- true
		}
		writeOdd := func() {
			entity := types.NewEntity(
				id,
				name,
				types.EntityTypeComponent,
				nil, nil, nil,
			)
			for index := mid; index < high; index++ {
				entity.ID = strconv.Itoa(index)
				entity.Name = name + "_HIGH_" + strconv.Itoa(index)
				ed.BulkInsert(entity)
			}
			finished <- true
		}

		test := func() {
			go writeOdd()
			go writeEven()
			<-finished
			<-finished
			err := ed.FlushBulk()
			if err != nil {
				panic(fmt.Errorf("test flushbulk on entity db adapter: %v", err))
			}
		}

		Convey("Then concurrent access on bulk insert should work", func() {
			So(test, ShouldNotPanic)
			count, err := ed.Count()
			So(err, ShouldBeNil)
			So(count, ShouldEqual, high)
		})
	})
}
