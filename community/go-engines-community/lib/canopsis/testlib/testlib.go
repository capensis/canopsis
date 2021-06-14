package testlib

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// GetCheckEvent generate a check event
func GetCheckEvent() types.Event {
	now := types.CpsTime{Time: time.Now().Truncate(time.Second)}
	return types.Event{
		Connector:     "rts",
		ConnectorName: "RA",
		Component:     "Soviet",
		Resource:      "Torun",
		EventType:     types.EventTypeCheck,
		SourceType:    types.SourceTypeResource,
		State:         types.AlarmStateMinor,
		Timestamp:     now,
	}
}
