package types_test

import (
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestComparison(t *testing.T) {
	Convey("Testing Ack", t, func() {
		var alarmChangeTypeAck = types.AlarmChange{
			Type:                types.AlarmChangeTypeAck,
			PreviousState:       types.AlarmStateOK,
			PreviousStateChange: datetime.NewCpsTime(),
		}
		So(len(alarmChangeTypeAck.GetTriggers()), ShouldEqual, 1)
		So(alarmChangeTypeAck.GetTriggers()[0], ShouldEqual, "ack")

		Convey("Testing default value", func() {
			var alarmChangeTypeNone = types.AlarmChange{
				Type:                types.AlarmChangeTypeNone,
				PreviousState:       types.AlarmStateOK,
				PreviousStateChange: datetime.NewCpsTime(),
			}
			So(len(alarmChangeTypeNone.GetTriggers()), ShouldEqual, 0)
		})
	})
}
