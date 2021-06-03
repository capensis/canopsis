package types_test

import (
	gogob "encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/gob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/msgpack"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFailOnError(t *testing.T) {
	panicfunc := func() {
		msg := "I want to panic"
		err := errors.New("because i'm crazy")
		utils.FailOnError(err, msg)
	}

	coolfunc := func() {
		var err error
		msg := "you don't see meeeee"
		utils.FailOnError(err, msg)
	}

	Convey("Must panic if err is not nil", t, func() {
		So(panicfunc, ShouldPanicWith, "I want to panic: because i'm crazy")
	})

	Convey("Must not panic if err is nil", t, func() {
		So(coolfunc, ShouldNotPanic)
	})
}

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
		value, err := types.AsInteger(types.CpsTime{Time: time.Unix(3, 0)})
		So(err, ShouldNotBeNil)
		So(value, ShouldEqual, 3)

		value, err = types.AsInteger(types.CpsTime{Time: time.Unix(-728, 0)})
		So(err, ShouldNotBeNil)
		So(value, ShouldEqual, -728)
	})
}

func TestCpsTime(t *testing.T) {
	Convey("Given a type struct with CpsTime", t, func() {
		type MyTime struct {
			TheTime types.CpsTime `json:"thetime" bson:"thetime"`
		}

		tnow := time.Now()
		stsnow := fmt.Sprint(tnow.Unix())
		mt := MyTime{
			TheTime: types.CpsTime{Time: tnow},
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

func TestCpsDuration(t *testing.T) {
	Convey("Given a type struct with CpsDuration", t, func() {
		type MyTime struct {
			TheDuration types.CpsDuration `json:"theduration" bson:"theduration"`
		}
		duration := time.Second*10 + time.Hour*2 + time.Minute*10
		mt := MyTime{
			TheDuration: types.CpsDuration(duration),
		}

		Convey("I can marshal to JSON", func() {
			bjdoc, err := json.Marshal(mt)
			So(err, ShouldBeNil)
			So(string(bjdoc), ShouldEqual, `{"theduration":"`+duration.String()+`"}`)

			Convey("And get my struct back from Marshal-ed JSON", func() {
				var mtb MyTime
				err = json.Unmarshal(bjdoc, &mtb)
				So(err, ShouldBeNil)
				So(mtb.TheDuration.Duration(), ShouldEqual, duration)
			})
		})

		Convey("I can marshal to BSON", func() {
			bsdoc, err := bson.Marshal(mt)
			So(err, ShouldBeNil)

			Convey("And get my struct back from Marshal-ed BSON", func() {
				var mtb MyTime
				err := bson.Unmarshal(bsdoc, &mtb)
				So(err, ShouldBeNil)
				So(mtb.TheDuration.Duration(), ShouldEqual, duration)
			})
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

func TestGOBCpsTime(t *testing.T) {
	Convey("Given time.Now()", t, func() {
		x := types.CpsTime{Time: time.Time{}}
		gogob.Register(x)
		now := types.CpsTime{Time: time.Now()}
		encoder := gob.NewEncoder()
		b, err := encoder.Encode(now)
		So(err, ShouldBeNil)

		Convey("Decode encoded CpsTime", func() {
			decoder := gob.NewDecoder()
			var oldnow types.CpsTime
			err := decoder.Decode(b, &oldnow)
			So(err, ShouldBeNil)
			So(now.Equal(oldnow.Time), ShouldBeTrue)
		})
	})
}

func TestMsgPackCpsTime(t *testing.T) {
	Convey("Given time.Now()", t, func() {
		now := types.CpsTime{Time: time.Now()}
		encoder := msgpack.NewEncoder()
		b, err := encoder.Encode(now)
		So(err, ShouldBeNil)

		Convey("Decode encoded CpsTime", func() {
			var oldnow types.CpsTime
			decoder := msgpack.NewDecoder()
			err := decoder.Decode(b, &oldnow)
			So(err, ShouldBeNil)
			So(now.Equal(oldnow.Time), ShouldBeTrue)
		})
	})
}

func TestMsgPackCpsNumber(t *testing.T) {
	Convey("Int", t, func() {
		i := types.CpsNumber(10)
		encoder := msgpack.NewEncoder()
		b, err := encoder.Encode(i)
		So(err, ShouldBeNil)

		Convey("decode", func() {
			var di types.CpsNumber
			dec := msgpack.NewDecoder()
			err := dec.Decode(b, &di)
			So(err, ShouldBeNil)
			So(di.Float64(), ShouldEqual, 10)
		})
	})
}

func TestMsgPackCpsDuration(t *testing.T) {
	Convey("Duration", t, func() {
		d := types.CpsDuration(time.Second * 10)
		enc := msgpack.NewEncoder()
		b, err := enc.Encode(d)
		So(err, ShouldBeNil)

		Convey("decode", func() {
			var dd types.CpsDuration
			dec := msgpack.NewDecoder()
			err := dec.Decode(b, &dd)
			So(err, ShouldBeNil)
			So(dd.Duration().Seconds(), ShouldEqual, 10)
		})
	})
}
func TestBinaryCpsTime(t *testing.T) {
	Convey("Given time.Now()", t, func() {
		t := types.CpsTime{Time: time.Now()}
		out := types.CpsTime{Time: time.Time{}}

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