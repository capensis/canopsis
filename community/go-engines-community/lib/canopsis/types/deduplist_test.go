package types_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDedupList(t *testing.T) {
	Convey("Setup", t, func() {
		fl := types.NewDedupList("imp1", "imp2")
		So("imp1", ShouldBeIn, fl.List())
		So("imp2", ShouldBeIn, fl.List())
		So(len(fl.List()), ShouldEqual, 2)

		fl.Add("imp3", "imp4")
		So("imp1", ShouldBeIn, fl.List())
		So("imp2", ShouldBeIn, fl.List())
		So("imp3", ShouldBeIn, fl.List())
		So("imp4", ShouldBeIn, fl.List())
		So(len(fl.List()), ShouldEqual, 4)

		// test remove middle element
		fl.Del("imp3")
		So("imp1", ShouldBeIn, fl.List())
		So("imp2", ShouldBeIn, fl.List())
		So("imp4", ShouldBeIn, fl.List())
		So("imp3", ShouldNotBeIn, fl.List())
		So(len(fl.List()), ShouldEqual, 3)

		// test remove last element
		fl.Del("imp4")
		So("imp1", ShouldBeIn, fl.List())
		So("imp2", ShouldBeIn, fl.List())
		So(len(fl.List()), ShouldEqual, 2)

		fl.Del("imp2")
		fl.Del("imp1")
		So(len(fl.List()), ShouldEqual, 0)
	})
}

func TestDedupListRead(t *testing.T) {
	Convey("Setup", t, func() {
		d := types.NewDedupList()
		d.Add("feeder2_0/feeder2_0")
		d.Add("feeder2_1/feeder2_0")
		d.Add("feeder2_0/feeder2_1")
		d.Add("feeder2_0/feeder2_1")
		d.Add("feeder2_0/feeder2_1")
		d.Add("feeder2_0/feeder2_1")
		d.Add("feeder2_0/feeder2_1")
		d.Add("feeder2_1/feeder2_1")

		l := d.List()
		So(len(l), ShouldEqual, 4)
		So(l, ShouldContain, "feeder2_0/feeder2_0")
		So(l, ShouldContain, "feeder2_0/feeder2_1")
		So(l, ShouldContain, "feeder2_1/feeder2_0")
		So(l, ShouldContain, "feeder2_1/feeder2_1")
	})
}

func TestDedupListJSON(t *testing.T) {
	Convey("Setup", t, func() {
		fl := types.NewDedupList("imp1", "imp2")
		b, err := json.Marshal(fl)
		So(err, ShouldBeNil)

		var fl2 types.DedupList
		So(json.Unmarshal(b, &fl2), ShouldBeNil)

		So(fl2.List(), ShouldContain, "imp1")
		So(fl2.List(), ShouldContain, "imp2")
		So(len(fl2.List()), ShouldEqual, 2)
	})
}

func TestDedupListBSON(t *testing.T) {
	t.SkipNow() // test creates bson.M instead of []interface{}
	Convey("Setup", t, func() {
		fl := types.NewDedupList("imp1", "imp2")
		b, err := bson.Marshal(fl)
		So(err, ShouldBeNil)

		var fl2 types.DedupList
		So(bson.Unmarshal(b, &fl2), ShouldBeNil)

		So(fl2.List(), ShouldContain, "imp1")
		So(fl2.List(), ShouldContain, "imp2")
		So(len(fl2.List()), ShouldEqual, 2)
	})
}

func BenchmarkDedupListAdd(b *testing.B) {
	l := types.NewDedupList()
	for i := 0; i < b.N; i++ {
		l.Add(fmt.Sprintf("item%d", i))
	}
}

func BenchmarkDedupListDel(b *testing.B) {
	l := types.NewDedupList()
	for i := 0; i < b.N; i++ {
		l.Add(fmt.Sprintf("item%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Del(fmt.Sprintf("item%d", i))
	}
}
