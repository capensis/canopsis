package datastorage

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

func CanRun(
	lastExecuted datetime.CpsTime,
	scheduledTimes config.ScheduledTimes,
	location *time.Location,
) bool {
	if len(scheduledTimes) == 0 {
		return false
	}

	now := datetime.NewCpsTime().In(location)
	if !scheduledTimes.IsScheduledTime(now) {
		return false
	}

	return !isAlreadyExecuted(now, lastExecuted.In(location))
}

func isAlreadyExecuted(now, lastExecuted datetime.CpsTime) bool {
	return now.EqualDay(lastExecuted) && now.Hour() == lastExecuted.Hour()
}
