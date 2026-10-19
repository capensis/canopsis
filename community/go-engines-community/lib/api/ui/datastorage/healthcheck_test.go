package datastorage

import (
	"testing"

	storage "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

func TestMergeHealthCheckHistory(t *testing.T) {
	t.Run("adds an observation without changing unrelated history", func(t *testing.T) {
		data := storage.DataStorage{History: map[string]storage.HistoryWithCount{
			"remediation": {Time: datetime.NewCpsTime(10), Deleted: 3},
		}}
		observation := storage.HistoryWithCount{Time: datetime.NewCpsTime(20)}

		(&store{}).mergeHealthCheckHistory(&data, &observation)

		if got := data.History[healthCheckHistoryKey]; !got.Time.Equal(observation.Time) || got.Deleted != 0 || got.Archived != 0 {
			t.Fatalf("unexpected health-check observation: %#v", got)
		}
		if got := data.History["remediation"]; !got.Time.Equal(datetime.NewCpsTime(10)) || got.Deleted != 3 {
			t.Fatalf("unrelated history changed: %#v", got)
		}
	})

	t.Run("retains a newer persisted observation and its fields", func(t *testing.T) {
		persisted := storage.HistoryWithCount{Time: datetime.NewCpsTime(30), Deleted: 4, Archived: 2}
		data := storage.DataStorage{History: map[string]storage.HistoryWithCount{healthCheckHistoryKey: persisted}}

		(&store{}).mergeHealthCheckHistory(&data, &storage.HistoryWithCount{Time: datetime.NewCpsTime(20)})

		if got := data.History[healthCheckHistoryKey]; !got.Time.Equal(persisted.Time) || got.Deleted != persisted.Deleted || got.Archived != persisted.Archived {
			t.Fatalf("newer persisted history regressed: %#v", got)
		}
	})
}
