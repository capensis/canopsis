package alarm

import (
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

var stateTitles = map[int]string{
	types.AlarmStateOK:       types.AlarmStateTitleOK,
	types.AlarmStateMinor:    types.AlarmStateTitleMinor,
	types.AlarmStateMajor:    types.AlarmStateTitleMajor,
	types.AlarmStateCritical: types.AlarmStateTitleCritical,
}
var statusTitles = map[int]string{
	types.AlarmStatusOff:       types.AlarmStatusTitleOff,
	types.AlarmStatusOngoing:   types.AlarmStatusTitleOngoing,
	types.AlarmStatusStealthy:  types.AlarmStatusTitleStealthy,
	types.AlarmStatusFlapping:  types.AlarmStatusTitleFlapping,
	types.AlarmStatusCancelled: types.AlarmStatusTitleCancelled,
}

func transformExportField(
	timeFormat string,
	location *time.Location,
) func(k string, v any) any {
	return func(k string, v any) any {
		switch k {
		case "v.state.val":
			if i, ok := getInt64(v); ok {
				return stateTitles[int(i)]
			}
		case "v.status.val":
			if i, ok := getInt64(v); ok {
				return statusTitles[int(i)]
			}
		case "t",
			"v.creation_date",
			"v.activation_date",
			"v.last_update_date",
			"v.last_event_date",
			"v.resolved":
			if i, ok := getInt64(v); ok {
				return types.NewCpsTime(i).In(location).Time.Format(timeFormat)
			}
		default:
			if strings.HasSuffix(k, ".t") {
				if i, ok := getInt64(v); ok {
					return types.NewCpsTime(i).In(location).Time.Format(timeFormat)
				}
			}
		}

		return v
	}
}

func getInt64(v any) (int64, bool) {
	switch i := v.(type) {
	case int64:
		return i, true
	case int32:
		return int64(i), true
	case int:
		return int64(i), true
	default:
		return 0, false
	}
}
