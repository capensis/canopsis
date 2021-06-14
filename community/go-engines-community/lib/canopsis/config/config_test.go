package config_test

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	"git.canopsis.net/canopsis/go-engines/lib/testutils"
	"os"
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConfWriteAndRead(t *testing.T) {
	Convey("Given a session, config DB adapter and a config", t, func() {
		dbClient, err := mongo.NewClient(0, 0)
		if err != nil {
			panic(err)
		}

		c := config.NewAdapter(dbClient)

		source := testutils.GetTestConf()

		Convey("Init should work", func() {
			err := c.UpsertConfig(source)
			So(err, ShouldBeNil)
		})

		Convey("Readed conf should be good", func() {
			conf, err := c.GetConfig()
			So(err, ShouldBeNil)
			So(conf.Alarm.FlappingFreqLimit, ShouldEqual, 1)
			So(conf.Alarm.CancelAutosolveDelay, ShouldEqual, "1h")
		})
	})
}

func TestConfSave(t *testing.T) {
	Convey("Given a session, config DB adapter and a config", t, func() {
		dbClient, err := mongo.NewClient(0, 0)
		if err != nil {
			panic(err)
		}

		c := config.NewAdapter(dbClient)
		source := testutils.GetTestConf()

		Convey("When we set a value", func() {
			source.Alarm.FlappingInterval = 666
			err := c.UpsertConfig(source)
			So(err, ShouldBeNil)

			Convey("The the value is on the database", func() {
				conf, err := c.GetConfig()
				So(err, ShouldBeNil)
				So(conf.Alarm.FlappingInterval, ShouldEqual, 666)
			})
		})
	})
}

func TestGetConfWrongMongo(t *testing.T) {
	Convey("Given a bad mongo url", t, func() {
		ou := os.Getenv(mongo.EnvURL)
		os.Setenv(mongo.EnvURL, "howmanytimeshaveitoldyaidontexist?")
		_, err := mongo.NewClient(0, 0)
		So(err, ShouldNotBeNil)
		os.Setenv(mongo.EnvURL, ou)
	})
}
