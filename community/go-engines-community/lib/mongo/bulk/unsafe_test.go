package bulk_test

import (
	"fmt"
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/testutils"
	"github.com/globalsign/mgo/bson"

	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/mongo/bulk"
	. "github.com/smartystreets/goconvey/convey"
)

const CollectionUnsafe = "test_bulk_unsafe"

type TestVal struct {
	Bla int `bson:"bla"`
}

func TestAutoBulkInsertUnsafe(t *testing.T) {
	Convey("Setup", t, func() {
		s, err := mongo.NewSession(time.Second * 200)
		So(err, ShouldBeNil)

		collection := s.DB(mongo.DB).C(CollectionUnsafe)
		err = testutils.MongoDrop(collection)
		So(err, ShouldBeNil)

		max := 10
		ab := bulk.New(collection, max)

		Convey(fmt.Sprintf("Add %d documents", max), func() {
			for i := 0; i < max; i++ {
				br, err := ab.AddInsert(bulk.NewInsert(TestVal{Bla: i}))
				So(err, ShouldBeNil)
				So(br, ShouldBeNil)
			}

			Convey("Flush insert operations", func() {
				br, err := ab.PerformInserts()
				So(err, ShouldBeNil)
				So(br, ShouldNotBeNil)

				Convey("Count inserted documents", func() {
					n, err := collection.Count()
					So(err, ShouldBeNil)
					So(n, ShouldEqual, max)
				})
			})
		})

		err = testutils.MongoDrop(collection)
		So(err, ShouldBeNil)

		max = bulk.BulkSizeMax
		mult := 100
		ab = bulk.New(collection, max)

		Convey(fmt.Sprintf("Insert %d documents", mult*max), func() {
			for i := 0; i < mult*max; i++ {
				_, err := ab.AddInsert(bulk.NewInsert(TestVal{Bla: i}))
				So(err, ShouldBeNil)
			}

			_, err := ab.PerformInserts()
			So(err, ShouldBeNil)

			n, err := collection.Count()
			So(err, ShouldBeNil)
			So(n, ShouldEqual, mult*max)
		})
	})
}

func TestBulkUniqUpdates(t *testing.T) {
	Convey("Setup", t, func() {
		s, err := mongo.NewSession(time.Second * 200)
		So(err, ShouldBeNil)

		collection := s.DB(mongo.DB).C(CollectionUnsafe)
		err = testutils.MongoDrop(collection)
		So(err, ShouldBeNil)

		max := 10
		ab := bulk.New(collection, max)

		Convey("Insert a document", func() {
			err := collection.Insert(bson.M{"_id": "bulk_uniq_update"})
			So(err, ShouldBeNil)
			Convey("Add the first update", func() {
				r, err := ab.AddUniqUpdate(
					"bulk_uniq_update",
					bulk.NewUpdate(
						bson.M{
							"_id": "bulk_uniq_update",
						},
						bson.M{
							"_id":    "bulk_uniq_update",
							"i_must": "disappear",
						},
					),
				)
				So(r, ShouldBeNil)
				So(err, ShouldBeNil)

				Convey("Add the second update", func() {
					r, err := ab.AddUniqUpdate(
						"bulk_uniq_update",
						bulk.NewUpdate(
							bson.M{
								"_id": "bulk_uniq_update",
							},
							bson.M{
								"_id":    "bulk_uniq_update",
								"i_must": "be there for you",
							},
						),
					)
					So(r, ShouldBeNil)
					So(err, ShouldBeNil)

					Convey("Then flush updates", func() {
						r, err := ab.PerformUniqUpdates()
						So(err, ShouldBeNil)
						So(r, ShouldNotBeNil)
						So(r.Modified, ShouldEqual, 1)
						So(r.Matched, ShouldEqual, 1)

						Convey("Get updated document", func() {
							var doc bson.M
							err := collection.Find(bson.M{"_id": "bulk_uniq_update"}).One(&doc)
							So(err, ShouldBeNil)
							So(doc["_id"], ShouldEqual, "bulk_uniq_update")
							So(doc["i_must"], ShouldEqual, "be there for you")
						})
					})
				})
			})
		})

	})
}
