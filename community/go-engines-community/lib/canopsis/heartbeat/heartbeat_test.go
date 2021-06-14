package heartbeat_test

import (
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/heartbeat"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHeartBeatMappingsDup(t *testing.T) {
	li := heartbeat.NewItem(0)
	err1 := li.AddMapping("connector", "test")
	err2 := li.AddMapping("connector", "test")

	if err1 != nil {
		t.Fatalf("should not raise error: %v", err1)
	}

	if err2 == nil {
		t.Fatalf("should raise an error")
	}

	if len(li.Mappings) != 1 {
		t.Fatalf("bad mappings lenght: %d", len(li.Mappings))
	}
}

func TestHeartbeatToItem(t *testing.T) {
	Convey("Given a ConfHeartBeatItem", t, func() {
		cli := heartbeat.Heartbeat{
			Pattern:          map[string]string{"connector": "test"},
			ExpectedInterval: "10s",
		}

		Convey("I can translate to a HeartBeatItem", func() {
			li, err := cli.ToHeartBeatItem()
			So(err, ShouldBeNil)
			So(li.Mappings, ShouldContainKey, "connector")
			So(li.Mappings["connector"], ShouldEqual, "test")
		})

		Convey("I cannot translate a bad MaxDuration value", func() {
			cli.ExpectedInterval = "war string"
			_, err := cli.ToHeartBeatItem()
			So(err, ShouldNotBeNil)
		})
	})
}
