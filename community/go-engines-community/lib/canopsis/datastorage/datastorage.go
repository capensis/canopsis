package datastorage

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libtime "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/time"
)

func CanRun(
	lastExecuted libtime.CpsTime,
	scheduledTime *config.ScheduledTime,
	location *time.Location,
) bool {
	// Skip if schedule is not defined.
	if scheduledTime == nil {
		return false
	}
	// Check now = schedule.
	now := libtime.NewCpsTime().In(location)
	if now.Weekday() != scheduledTime.Weekday || now.Hour() != scheduledTime.Hour {
		return false
	}

	//Skip if already executed today.
	if lastExecuted.Unix() > 0 && lastExecuted.EqualDay(now) {
		return false
	}

	return true
}
