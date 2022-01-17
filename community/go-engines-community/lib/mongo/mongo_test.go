package mongo_test

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestNewMongoSession(t *testing.T) {
	Convey("Wanna check some good MongoDB session?", t, func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		Convey("Bad url", func() {
			ou := os.Getenv(mongo.EnvURL)
			os.Setenv(mongo.EnvURL, "-... .- -.. / ..- .-. .-..")
			_, err := mongo.NewClient(ctx, 0, 0, zerolog.Nop())
			So(err, ShouldNotBeNil)
			os.Setenv(mongo.EnvURL, ou)
		})

		Convey("bad host", func() {
			ou := os.Getenv(mongo.EnvURL)
			os.Setenv(mongo.EnvURL, "mongodb://I-TOLD-YOU-I-AM-A-GHOST:27017")
			_, err := mongo.NewClient(ctx, 0, 0, zerolog.Nop())
			So(err, ShouldNotBeNil)
			os.Setenv(mongo.EnvURL, ou)
		})

		Convey("all good", func() {
			_, err := mongo.NewClient(ctx, 0, 0, zerolog.Nop())
			So(err, ShouldBeNil)
		})
	})
}
