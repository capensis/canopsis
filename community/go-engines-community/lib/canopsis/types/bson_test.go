package types_test

import (
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

type intWrapper struct {
	Value types.OptionalInt64 `bson:"value"`
}

func TestOptionalInt64MongoDriver(t *testing.T) {
	Convey("Given a BSON document containing an integer", t, func() {
		document := bson.M{
			"value": 12,
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("The document should be decoded without error", func() {
			var w intWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldBeNil)

			Convey("The value should be set", func() {
				So(w.Value.Set, ShouldBeTrue)
				So(w.Value.Value, ShouldEqual, 12)
			})
		})
	})
	Convey("Given a BSON document containing an double", t, func() {
		document := bson.M{
			"value": 3.0,
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("The document should be decoded without error", func() {
			var w intWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldBeNil)

			Convey("The value should be set", func() {
				So(w.Value.Set, ShouldBeTrue)
				So(w.Value.Value, ShouldEqual, 3)
			})
		})
	})

	Convey("Given a BSON document containing a string", t, func() {
		document := bson.M{
			"value": "a string",
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("Decoding the document should return an error", func() {
			var w intWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldNotBeNil)
		})
	})

	Convey("Given an empty BSON document", t, func() {
		document := bson.M{}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("The document should be decoded without error", func() {
			var w intWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldBeNil)

			Convey("The value should not be set", func() {
				So(w.Value.Set, ShouldBeFalse)
			})
		})
	})
}

type boolWrapper struct {
	Value types.OptionalBool `bson:"value"`
}

func TestOptionalBoolMongoDriver(t *testing.T) {
	Convey("Given a BSON document containing a true boolean", t, func() {
		document := bson.M{
			"value": true,
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("The document should be decoded without error", func() {
			var w boolWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldBeNil)

			Convey("The value should be set", func() {
				So(w.Value.Set, ShouldBeTrue)
				So(w.Value.Value, ShouldBeTrue)
			})
		})
	})

	Convey("Given a BSON document containing a false boolean", t, func() {
		document := bson.M{
			"value": false,
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("The document should be decoded without error", func() {
			var w boolWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldBeNil)

			Convey("The value should be set", func() {
				So(w.Value.Set, ShouldBeTrue)
				So(w.Value.Value, ShouldBeFalse)
			})
		})
	})

	Convey("Given a BSON document containing an integer", t, func() {
		document := bson.M{
			"value": 12,
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("Decoding the document should return an error", func() {
			var w boolWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON document containing a string", t, func() {
		document := bson.M{
			"value": "a string",
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("Decoding the document should return an error", func() {
			var w boolWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldNotBeNil)
		})
	})

	Convey("Given an empty BSON document", t, func() {
		document := bson.M{}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("The document should be decoded without error", func() {
			var w boolWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldBeNil)

			Convey("The value should not be set", func() {
				So(w.Value.Set, ShouldBeFalse)
			})
		})
	})
}

type stringWrapper struct {
	Value types.OptionalString `bson:"value"`
}

func TestOptionalStringMongoDriver(t *testing.T) {
	Convey("Given a BSON document containing a string", t, func() {
		document := bson.M{
			"value": "a string",
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("The document should be decoded without error", func() {
			var w stringWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldBeNil)

			Convey("The value should be set", func() {
				So(w.Value.Set, ShouldBeTrue)
				So(w.Value.Value, ShouldEqual, "a string")
			})
		})
	})

	Convey("Given a BSON document containing an integer", t, func() {
		document := bson.M{
			"value": 12,
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("Decoding the document should return an error", func() {
			var w stringWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldNotBeNil)
		})
	})

	Convey("Given an empty BSON document", t, func() {
		document := bson.M{}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("The document should be decoded without error", func() {
			var w stringWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldBeNil)

			Convey("The value should not be set", func() {
				So(w.Value.Set, ShouldBeFalse)
			})
		})
	})
}

type interfaceWrapper struct {
	Value types.OptionalInterface `bson:"value"`
}

func TestOptionalInterfaceMongoDriver(t *testing.T) {
	Convey("Given a BSON document containing an integer", t, func() {
		document := bson.M{
			"value": 12,
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("The document should be decoded without error", func() {
			var w interfaceWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldBeNil)

			Convey("The value should be set", func() {
				So(w.Value.Set, ShouldBeTrue)
				So(w.Value.Value, ShouldEqual, 12)
			})
		})
	})

	Convey("Given a BSON document containing a string", t, func() {
		document := bson.M{
			"value": "a string",
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("The document should be decoded without error", func() {
			var w interfaceWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldBeNil)

			Convey("The value should be set", func() {
				So(w.Value.Set, ShouldBeTrue)
				So(w.Value.Value, ShouldEqual, "a string")
			})
		})
	})

	Convey("Given an empty BSON document", t, func() {
		document := bson.M{}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("The document should be decoded without error", func() {
			var w interfaceWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldBeNil)

			Convey("The value should not be set", func() {
				So(w.Value.Set, ShouldBeFalse)
			})
		})
	})
}

type regexpWrapper struct {
	Value types.OptionalRegexp `bson:"value"`
}

func TestOptionalRegexpMongoDriver(t *testing.T) {
	Convey("Given a BSON document containing a valid regular expression", t, func() {
		document := bson.M{
			"value": "abc-(.*)-def",
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("The document should be decoded without error", func() {
			var w regexpWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldBeNil)

			Convey("The value should be set", func() {
				So(w.Value.Set, ShouldBeTrue)
				So(w.Value.Value.Match([]byte("abc-test-def")), ShouldBeTrue)
				So(w.Value.Value.Match([]byte("abc-ok-!-def")), ShouldBeTrue)
				So(w.Value.Value.Match([]byte("not-a-match")), ShouldBeFalse)
			})
		})
	})

	Convey("Given a BSON document containing a valid negative lookahead regular expression", t, func() {
		document := bson.M{
			"value": "^(?!resource_CPU).*$",
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("The document should be decoded without error", func() {
			var w regexpWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldBeNil)

			Convey("The value should be set", func() {
				So(w.Value.Set, ShouldBeTrue)
				So(w.Value.Value.Match([]byte("resource_CPU_overloaded")), ShouldBeFalse)
				So(w.Value.Value.Match([]byte("server_resource_CPU_drop")), ShouldBeTrue)
				So(w.Value.Value.Match([]byte("resource_CPU")), ShouldBeFalse)
				So(w.Value.Value.Match([]byte("just a string")), ShouldBeTrue)
			})
		})
	})

	Convey("Given a BSON document containing an invalid regular expression", t, func() {
		document := bson.M{
			"value": "abc-(.*-def",
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("Decoding the document should return an error", func() {
			var w regexpWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON document containing an integer", t, func() {
		document := bson.M{
			"value": 12,
		}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("Decoding the document should return an error", func() {
			var w regexpWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldNotBeNil)
		})
	})

	Convey("Given an empty BSON document", t, func() {
		document := bson.M{}
		bsonDocument, err := bson.Marshal(document)
		So(err, ShouldBeNil)

		Convey("The document should be decoded without error", func() {
			var w regexpWrapper
			So(bson.Unmarshal(bsonDocument, &w), ShouldBeNil)

			Convey("The value should not be set", func() {
				So(w.Value.Set, ShouldBeFalse)
			})
		})
	})
}
