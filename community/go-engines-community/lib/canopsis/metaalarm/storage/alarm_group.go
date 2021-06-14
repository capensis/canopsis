package storage

import "git.canopsis.net/canopsis/go-engines/lib/canopsis/types"

type AlarmGroup struct {
	OpenTime types.CpsTime
	Group    map[string]types.CpsTime
}

func (g AlarmGroup) GetAlarmIds() []string {
	ids := make([]string, 0, len(g.Group))

	for id := range g.Group {
		ids = append(ids, id)
	}

	return ids
}
