package datastorage

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func CanRun(
	lastExecuted types.CpsTime,
	scheduledTime *config.ScheduledTime,
	location *time.Location,
) bool {
	// Skip if schedule is not defined.
	if scheduledTime == nil {
		return false
	}
	// Check now = schedule.
	now := time.Now().In(location)
	if now.Weekday() != scheduledTime.Weekday || now.Hour() != scheduledTime.Hour {
		return false
	}

	//Skip if already executed today.
	if lastExecuted.Unix() > 0 {
		dateFormat := "2006-01-02"
		if lastExecuted.Time.In(time.UTC).Format(dateFormat) == now.In(time.UTC).Format(dateFormat) {
			return false
		}
	}

	return true
}
