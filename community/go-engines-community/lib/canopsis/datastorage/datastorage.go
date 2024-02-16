package datastorage

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

func CanRun(
	lastExecuted datetime.CpsTime,
	scheduledTime *config.ScheduledTime,
	location *time.Location,
) bool {
	// Skip if schedule is not defined.
	if scheduledTime == nil {
		return false
	}
	// Check now = schedule.
	now := datetime.NewCpsTime().In(location)
	if now.Weekday() != scheduledTime.Weekday || now.Hour() != scheduledTime.Hour {
		return false
	}

	//Skip if already executed today.
	if lastExecuted.Unix() > 0 && lastExecuted.EqualDay(now) {
		return false
	}

	return true
}
