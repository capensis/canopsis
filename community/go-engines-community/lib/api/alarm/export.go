package alarm

import (
	"context"
	"strconv"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
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

func getDataFetcher(
	store Store,
	r ExportRequest,
	exportFields []string,
	location *time.Location,
) export.DataFetcher {
	return func(ctx context.Context, page, limit int64) ([]map[string]string, int64, error) {
		res, err := store.Find(ctx, ListRequestWithPagination{
			Query: pagination.Query{Page: page, Limit: limit, Paginate: true},
			ListRequest: ListRequest{
				FilterRequest: FilterRequest{
					BaseFilterRequest: r.BaseFilterRequest,
					SearchBy:          exportFields,
				},
			},
		})
		if err != nil {
			return nil, 0, err
		}

		data, err := export.ConvertToMap(res.Data, exportFields,
			common.GetRealFormatTime(r.TimeFormat), location)
		if err != nil {
			return nil, 0, err
		}

		for i := range data {
			if v, ok := data[i]["v.state.val"]; ok {
				key := stringToInt(v)
				if key >= 0 {
					data[i]["v.state.val"] = stateTitles[key]
				}
			}
			if v, ok := data[i]["v.status.val"]; ok {
				key := stringToInt(v)
				if key >= 0 {
					data[i]["v.status.val"] = statusTitles[key]
				}
			}
		}

		return data, res.TotalCount, nil
	}
}

func stringToInt(i string) int {
	key, err := strconv.Atoi(i)
	if err != nil {
		return -1
	}

	return key
}
