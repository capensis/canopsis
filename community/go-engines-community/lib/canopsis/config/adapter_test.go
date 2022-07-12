package config_test

import (
	"context"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/testutils"
	"github.com/rs/zerolog"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConfWriteAndRead(t *testing.T) {
	Convey("Given a session, config DB adapter and a config", t, func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dbClient, err := mongo.NewClient(ctx, 0, 0, zerolog.Nop())
		if err != nil {
			panic(err)
		}

		c := config.NewAdapter(dbClient)

		source := testutils.GetTestConf()

		Convey("Init should work", func() {
			err := c.UpsertConfig(ctx, source)
			So(err, ShouldBeNil)
		})

		Convey("Readed conf should be good", func() {
			conf, err := c.GetConfig(ctx)
			So(err, ShouldBeNil)
			So(conf.Alarm.CancelAutosolveDelay, ShouldEqual, "1h")
		})
	})
}

func TestConfSave(t *testing.T) {
	Convey("Given a session, config DB adapter and a config", t, func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dbClient, err := mongo.NewClient(ctx, 0, 0, zerolog.Nop())
		if err != nil {
			panic(err)
		}

		c := config.NewAdapter(dbClient)
		source := testutils.GetTestConf()

		Convey("When we set a value", func() {
			err := c.UpsertConfig(ctx, source)
			So(err, ShouldBeNil)

			Convey("The the value is on the database", func() {
				_, err := c.GetConfig(ctx)
				So(err, ShouldBeNil)
			})
		})
	})
}
