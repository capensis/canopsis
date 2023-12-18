package types_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

func TestInterfaceToString(t *testing.T) {
	Convey("Given some values to convert", t, func() {
		Convey("I don't know how to convert a list", func() {
			vlstring := make([]string, 1)
			vlstring[0] = "bla"
			val, err := types.InterfaceToString(vlstring)
			So(err, ShouldNotBeNil)
			So(val, ShouldEqual, "")
		})

		Convey("I can convert string to string, as dumb as it sounds", func() {
			val, err := types.InterfaceToString("keukou")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "keukou")
		})

		Convey("I can convert a bool", func() {
			val, err := types.InterfaceToString(true)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "true")
		})

		Convey("I can convert a float64", func() {
			val, err := types.InterfaceToString(float64(1.1))
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "1.1")
		})

		Convey("I can convert an int", func() {
			val, err := types.InterfaceToString(int(1))
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "1")
		})

		Convey("I can convert an int64", func() {
			val, err := types.InterfaceToString(int64(1))
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "1")
		})

		Convey("I can convert an uint", func() {
			val, err := types.InterfaceToString(uint(1))
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "1")
		})

		Convey("I can convert an uint64", func() {
			val, err := types.InterfaceToString(uint64(1))
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "1")
		})

		Convey("I can convert a list of int", func() {
			values := make([]interface{}, 2)
			values[0] = 1
			values[1] = 2
			val, err := types.InterfaceToString(values)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "1,2")
		})
	})
}

func TestAsInteger(t *testing.T) {
	Convey("Given an int, AsInteger returns its value", t, func() {
		value, err := types.AsInteger(int(3))
		So(err, ShouldNotBeNil)
		So(value, ShouldEqual, 3)

		value, err = types.AsInteger(int(-728))
		So(err, ShouldNotBeNil)
		So(value, ShouldEqual, -728)
	})

	Convey("Given an uint, AsInteger returns its value", t, func() {
		value, err := types.AsInteger(uint(3))
		So(err, ShouldNotBeNil)
		So(value, ShouldEqual, 3)

		value, err = types.AsInteger(uint(728))
		So(err, ShouldNotBeNil)
		So(value, ShouldEqual, 728)
	})

	Convey("Given an int64, AsInteger returns its value", t, func() {
		value, err := types.AsInteger(int64(3))
		So(err, ShouldNotBeNil)
		So(value, ShouldEqual, 3)

		value, err = types.AsInteger(int64(-728))
		So(err, ShouldNotBeNil)
		So(value, ShouldEqual, -728)
	})

	Convey("Given an uint64, AsInteger returns its value", t, func() {
		value, err := types.AsInteger(uint64(3))
		So(err, ShouldNotBeNil)
		So(value, ShouldEqual, 3)

		value, err = types.AsInteger(uint64(728))
		So(err, ShouldNotBeNil)
		So(value, ShouldEqual, 728)
	})

	Convey("Given a CpsNumber, AsInteger returns its value", t, func() {
		value, err := types.AsInteger(types.CpsNumber(3))
		So(err, ShouldNotBeNil)
		So(value, ShouldEqual, 3)

		value, err = types.AsInteger(types.CpsNumber(-728))
		So(err, ShouldNotBeNil)
		So(value, ShouldEqual, -728)
	})

	Convey("Given a CpsTime, AsInteger returns its value", t, func() {
		value, err := types.AsInteger(datetime.CpsTime{Time: time.Unix(3, 0)})
		So(err, ShouldNotBeNil)
		So(value, ShouldEqual, 3)

		value, err = types.AsInteger(datetime.CpsTime{Time: time.Unix(-728, 0)})
		So(err, ShouldNotBeNil)
		So(value, ShouldEqual, -728)
	})
}

func TestCpsTime(t *testing.T) {
	Convey("Given a type struct with CpsTime", t, func() {
		type MyTime struct {
			TheTime datetime.CpsTime `json:"thetime" bson:"thetime"`
		}

		tnow := time.Now()
		stsnow := strconv.FormatInt(tnow.Unix(), 10)
		mt := MyTime{
			TheTime: datetime.CpsTime{Time: tnow},
		}

		Convey("I can marshal to JSON", func() {
			bjdoc, err := json.Marshal(mt)
			So(err, ShouldBeNil)
			So(string(bjdoc), ShouldEqual, `{"thetime":`+stsnow+"}")

			Convey("And get my struct back from Marshal-ed JSON", func() {
				var mtb MyTime
				err = json.Unmarshal(bjdoc, &mtb)
				So(err, ShouldBeNil)
				So(mtb.TheTime.Unix(), ShouldEqual, tnow.Unix())
			})
		})

		Convey("I can marshal to BSON", func() {
			bsdoc, err := bson.Marshal(mt)
			So(err, ShouldBeNil)

			Convey("And get my struct back from Marshal-ed BSON", func() {
				var mtb MyTime
				err := bson.Unmarshal(bsdoc, &mtb)
				So(err, ShouldBeNil)
				So(mtb.TheTime.Unix(), ShouldEqual, tnow.Unix())
			})
		})

		Convey("I can unmarshal a JSON float value", func() {
			bjdoc := fmt.Sprintf(`{"thetime":%g}`, 1.41421356237)

			var mtb MyTime
			err := json.Unmarshal([]byte(bjdoc), &mtb)
			So(err, ShouldBeNil)
			So(mtb.TheTime.Unix(), ShouldEqual, time.Unix(1, 0).Unix())
		})
	})
}

func TestCpsNumber(t *testing.T) {
	Convey("Setup", t, func() {
		type MyEvent struct {
			State types.CpsNumber `json:"state" bson:"state"`
		}

		me := MyEvent{
			State: 1.0,
		}

		Convey("For JSON documents", func() {
			Convey("Unmarshal", func() {
				floatJSON := `{"state": 1.1}`
				intJSON := `{"state": 1}`

				var floatEVT MyEvent
				var intEVT MyEvent

				err := json.Unmarshal([]byte(floatJSON), &floatEVT)
				So(err, ShouldBeNil)

				err = json.Unmarshal([]byte(intJSON), &intEVT)
				So(err, ShouldBeNil)

				So(floatEVT.State, ShouldEqual, types.CpsNumber(1.0))
				So(intEVT.State, ShouldEqual, types.CpsNumber(1.0))
				So(floatEVT.State, ShouldHaveSameTypeAs, types.CpsNumber(1))
				So(intEVT.State, ShouldHaveSameTypeAs, types.CpsNumber(1))
			})

			Convey("Marshal", func() {
				bjson, err := json.Marshal(me)
				So(err, ShouldBeNil)
				So(string(bjson), ShouldEqual, `{"state":1}`)
			})
		})

		Convey("For BSON documents", func() {
			bsdoc, err := bson.Marshal(me)
			So(err, ShouldBeNil)

			Convey("And get my struct back from Marshal-ed BSON", func() {
				var meb MyEvent
				err := bson.Unmarshal(bsdoc, &meb)
				So(err, ShouldBeNil)
				So(meb.State, ShouldEqual, 1)
			})
		})
	})
}

func TestBinaryCpsTime(t *testing.T) {
	Convey("Given time.Now()", t, func() {
		t := datetime.NewCpsTime()
		out := datetime.CpsTime{Time: time.Time{}}

		So(out.Equal(t.Time), ShouldBeFalse)

		Convey("I can encode", func() {
			b, err := t.MarshalBinary()
			So(err, ShouldBeNil)

			Convey("I can decode", func() {
				err := out.UnmarshalBinary(b)
				So(err, ShouldBeNil)
				So(out.Equal(t.Time), ShouldBeTrue)
			})
		})
	})
}
