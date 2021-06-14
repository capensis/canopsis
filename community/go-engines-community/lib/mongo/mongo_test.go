package mongo_test

import (
	"os"
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewMongoSession(t *testing.T) {
	Convey("Wanna check some good MongoDB session?", t, func() {
		Convey("Bad url", func() {
			ou := os.Getenv(mongo.EnvURL)
			os.Setenv(mongo.EnvURL, "-... .- -.. / ..- .-. .-..")
			_, err := mongo.NewSession(time.Second)
			So(err, ShouldNotBeNil)
			os.Setenv(mongo.EnvURL, ou)
		})

		Convey("bad host", func() {
			ou := os.Getenv(mongo.EnvURL)
			os.Setenv(mongo.EnvURL, "mongodb://I-TOLD-YOU-I-AM-A-GHOST:27017")
			_, err := mongo.NewSession(time.Second)
			So(err, ShouldNotBeNil)
			os.Setenv(mongo.EnvURL, ou)
		})

		Convey("all good", func() {
			_, err := mongo.NewSession(time.Second)
			So(err, ShouldBeNil)
		})
	})
}
