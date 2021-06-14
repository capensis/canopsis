package utils_test

import (
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	. "github.com/smartystreets/goconvey/convey"
)

type B struct {
	IntField       int
	IntPtrField    *int
	StringField    string
	StringPtrField *string
	MapField       map[string]interface{}
	MapPtrField    *map[string]interface{}
	NilPtrField    *int
}

type A struct {
	IntField       int
	IntPtrField    *int
	StringField    string
	StringPtrField *string
	MapField       map[string]interface{}
	MapPtrField    *map[string]interface{}
	NilPtrField    *int
	BField         B
	BPtrField      *B
}

func TestGetField(t *testing.T) {
	Convey("Given an instance of a struct, GetField returns the correct values for:", t, func() {
		int1 := 0
		int2 := 1
		string1 := "testptr"
		string2 := "abcde"
		map1 := make(map[string]interface{})
		map2 := make(map[string]interface{})

		map1["int"] = 5
		map1["string"] = "klmno"
		map1["nil"] = nil

		map2["int"] = 6
		map2["string"] = "pqrst"
		map2["nil"] = nil

		b := B{
			IntField:       2,
			IntPtrField:    &int1,
			StringField:    "fghij",
			StringPtrField: &string1,
			MapField:       map2,
			MapPtrField:    &map2,
			NilPtrField:    nil,
		}

		a := A{
			IntField:       3,
			IntPtrField:    &int2,
			StringField:    "test",
			StringPtrField: &string2,
			MapField:       map1,
			MapPtrField:    &map1,
			NilPtrField:    nil,
			BField:         b,
			BPtrField:      &b,
		}

		Convey("The fields of this object", func() {
			value, err := utils.GetField(a, "IntField")
			So(err, ShouldBeNil)
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 3)
			value, err = utils.GetField(a, "IntPtrField")
			So(err, ShouldBeNil)
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 1)
			value, err = utils.GetField(a, "StringField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "test")
			value, err = utils.GetField(a, "StringPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "abcde")
			value, err = utils.GetField(a, "MapField")
			So(err, ShouldBeNil)
			So(value, ShouldResemble, map1)
			value, err = utils.GetField(a, "MapPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldResemble, map1)
			value, err = utils.GetField(a, "NilPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
			value, err = utils.GetField(a, "BField")
			So(err, ShouldBeNil)
			So(value, ShouldResemble, b)
			value, err = utils.GetField(a, "BPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldResemble, b)
		})

		Convey("The values of a map in this object", func() {
			value, err := utils.GetField(a, "MapField.int")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 5)
			value, err = utils.GetField(a, "MapField.string")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "klmno")
			value, err = utils.GetField(a, "MapField.nil")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The values of a reference to a map in this object", func() {
			value, err := utils.GetField(a, "MapPtrField.int")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 5)
			value, err = utils.GetField(a, "MapPtrField.string")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "klmno")
			value, err = utils.GetField(a, "MapPtrField.nil")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The fields of an object in this object", func() {
			value, err := utils.GetField(a, "BField.IntField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 2)
			value, err = utils.GetField(a, "BField.IntPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 0)
			value, err = utils.GetField(a, "BField.StringField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "fghij")
			value, err = utils.GetField(a, "BField.StringPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "testptr")
			value, err = utils.GetField(a, "BField.NilPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The values of a map in an object in this object", func() {
			value, err := utils.GetField(a, "BField.MapField.int")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 6)
			value, err = utils.GetField(a, "BField.MapField.string")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "pqrst")
			value, err = utils.GetField(a, "BField.MapField.nil")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The values of a reference to a map in an object in this object", func() {
			value, err := utils.GetField(a, "BField.MapPtrField.int")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 6)
			value, err = utils.GetField(a, "BField.MapPtrField.string")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "pqrst")
			value, err = utils.GetField(a, "BField.MapPtrField.nil")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The fields of a reference to an object in this object", func() {
			value, err := utils.GetField(a, "BPtrField.IntField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 2)
			value, err = utils.GetField(a, "BPtrField.IntPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 0)
			value, err = utils.GetField(a, "BPtrField.StringField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "fghij")
			value, err = utils.GetField(a, "BPtrField.StringPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "testptr")
			value, err = utils.GetField(a, "BPtrField.NilPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The values of a map in a reference to an object in this object", func() {
			value, err := utils.GetField(a, "BPtrField.MapField.int")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 6)
			value, err = utils.GetField(a, "BPtrField.MapField.string")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "pqrst")
			value, err = utils.GetField(a, "BPtrField.MapField.nil")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The values of a reference to a map in a reference to an object in this object", func() {
			value, err := utils.GetField(a, "BPtrField.MapPtrField.int")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 6)
			value, err = utils.GetField(a, "BPtrField.MapPtrField.string")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "pqrst")
			value, err = utils.GetField(a, "BPtrField.MapPtrField.nil")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The fields of a reference to this object", func() {
			value, err := utils.GetField(&a, "IntField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 3)
			value, err = utils.GetField(&a, "IntPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 1)
			value, err = utils.GetField(&a, "StringField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "test")
			value, err = utils.GetField(&a, "StringPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "abcde")
			value, err = utils.GetField(&a, "MapField")
			So(err, ShouldBeNil)
			So(value, ShouldResemble, map1)
			value, err = utils.GetField(&a, "MapPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldResemble, map1)
			value, err = utils.GetField(&a, "NilPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
			value, err = utils.GetField(&a, "BField")
			So(err, ShouldBeNil)
			So(value, ShouldResemble, b)
			value, err = utils.GetField(&a, "BPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldResemble, b)
		})

		Convey("The values of a map in a reference to this object", func() {
			value, err := utils.GetField(&a, "MapField.int")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 5)
			value, err = utils.GetField(&a, "MapField.string")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "klmno")
			value, err = utils.GetField(&a, "MapField.nil")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The values of a reference to a map in a reference to this object", func() {
			value, err := utils.GetField(&a, "MapPtrField.int")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 5)
			value, err = utils.GetField(&a, "MapPtrField.string")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "klmno")
			value, err = utils.GetField(&a, "MapPtrField.nil")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The fields of an object in a reference to this object", func() {
			value, err := utils.GetField(&a, "BField.IntField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 2)
			value, err = utils.GetField(&a, "BField.IntPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 0)
			value, err = utils.GetField(&a, "BField.StringField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "fghij")
			value, err = utils.GetField(&a, "BField.StringPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "testptr")
			value, err = utils.GetField(&a, "BField.NilPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The values of a map in an object in a reference to this object", func() {
			value, err := utils.GetField(&a, "BField.MapField.int")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 6)
			value, err = utils.GetField(&a, "BField.MapField.string")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "pqrst")
			value, err = utils.GetField(&a, "BField.MapField.nil")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The values of a reference to a map in an object in a reference to this object", func() {
			value, err := utils.GetField(&a, "BField.MapPtrField.int")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 6)
			value, err = utils.GetField(&a, "BField.MapPtrField.string")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "pqrst")
			value, err = utils.GetField(&a, "BField.MapPtrField.nil")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The fields of a reference to an object in a reference to this object", func() {
			value, err := utils.GetField(&a, "BPtrField.IntField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 2)
			value, err = utils.GetField(&a, "BPtrField.IntPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 0)
			value, err = utils.GetField(&a, "BPtrField.StringField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "fghij")
			value, err = utils.GetField(&a, "BPtrField.StringPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "testptr")
			value, err = utils.GetField(&a, "BPtrField.NilPtrField")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The values of a map in a reference to an object in a reference to this object", func() {
			value, err := utils.GetField(&a, "BPtrField.MapField.int")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 6)
			value, err = utils.GetField(&a, "BPtrField.MapField.string")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "pqrst")
			value, err = utils.GetField(&a, "BPtrField.MapField.nil")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("The values of a reference to a map in a reference to an object in a reference to this object", func() {
			value, err := utils.GetField(&a, "BPtrField.MapPtrField.int")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, 6)
			value, err = utils.GetField(&a, "BPtrField.MapPtrField.string")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "pqrst")
			value, err = utils.GetField(&a, "BPtrField.MapPtrField.nil")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, nil)
		})

		Convey("Undefined fields", func() {
			_, err := utils.GetField(a, "undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(a, "IntPtrField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(a, "MapField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(a, "MapPtrField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(a, "BField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(a, "BField.MapField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(a, "BField.MapPtrField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(a, "BPtrField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(a, "BPtrField.MapField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(a, "BPtrField.MapPtrField.undefined")
			So(err, ShouldNotBeNil)

			_, err = utils.GetField(&a, "undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(&a, "IntPtrField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(&a, "MapField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(&a, "MapPtrField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(&a, "BField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(&a, "BField.MapField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(&a, "BField.MapPtrField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(&a, "BPtrField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(&a, "BPtrField.MapField.undefined")
			So(err, ShouldNotBeNil)
			_, err = utils.GetField(&a, "BPtrField.MapPtrField.undefined")
			So(err, ShouldNotBeNil)
		})
	})
}
