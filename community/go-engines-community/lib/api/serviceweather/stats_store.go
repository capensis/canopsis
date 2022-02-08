package serviceweather

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	StatsCollection = "entity_stats"
)

// StatsStore to read, save and reset entity stats
type StatsStore interface {
	GetStats(context.Context, string, *time.Location) (Stats, error)
	SetStats(context.Context, string, Stats, *time.Location) error
	ResetStats(context.Context, string) error
}

// NewStatsStore instantiates a StatsStore
func NewStatsStore(dbClient mongo.DbClient) StatsStore {
	return &statsStore{
		dbCollection: dbClient.Collection(StatsCollection),
	}
}

type statsStore struct {
	dbCollection mongo.DbCollection
}

// GetStats fetches entity stats by its ID. Returns zero stats when last event has before current date
func (s *statsStore) GetStats(ctx context.Context, eid string, location *time.Location) (result Stats, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	entityStats := s.dbCollection.FindOne(ctx, bson.M{"_id": eid})
	if err = entityStats.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return result, nil
		}
		return result, err
	}

	err = entityStats.Decode(&result)
	if err != nil || result.LastEvent.Time.IsZero() {
		return result, err
	}

	nowYear, nowMonth, nowDay := time.Now().In(location).Date()
	prevYear, prevMonth, prevDay := result.LastEvent.Time.In(location).Date()

	if !(nowYear == prevYear && nowMonth == prevMonth && nowDay == prevDay) {
		result.OKEventsCount = 0
		result.FailEventsCount = 0
	}
	return result, nil
}

// SetStats updates or innsert entity stats by Entity ID when last saved event has same date.
// The new day's stats starts from zero counters. It replaces stats when lastevevent's date was before the current date
func (s *statsStore) SetStats(ctx context.Context, eid string, st Stats, location *time.Location) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	lastEventTime := time.Time{}
	if st.LastEvent == nil {
		lastEventTime = time.Now()
		st.LastEvent = &types.CpsTime{Time: lastEventTime}
	}

	filter := bson.M{"_id": eid}
	updateFields := bson.M{"last_event": st.LastEvent}
	if st.FailEventsCount > 0 {
		updateFields["last_ko"] = st.LastFailEvent
	}
	update := bson.M{
		"$set": updateFields,
		"$inc": bson.M{"ok": st.OKEventsCount, "ko": st.FailEventsCount},
	}
	res := s.dbCollection.FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.Before))
	if err := res.Err(); err != nil && err != mongodriver.ErrNoDocuments {
		return err
	}

	prev := Stats{}
	if err := res.Decode(&prev); err != nil {
		return nil
	}

	if prev.LastEvent != nil {
		if lastEventTime.IsZero() {
			lastEventTime = st.LastEvent.Time
		}

		evtYear, evtMonth, evtDay := lastEventTime.In(location).Date()
		prevYear, prevMonth, prevDay := prev.LastEvent.Time.In(location).Date()

		if !(evtYear == prevYear && evtMonth == prevMonth && evtDay == prevDay) {
			updateFields["ko"] = st.FailEventsCount
			updateFields["ok"] = st.OKEventsCount
			if _, exist := updateFields["last_ko"]; !exist {
				updateFields["last_ko"] = prev.LastFailEvent
			}
			res = s.dbCollection.FindOneAndReplace(ctx, filter, updateFields)
			if err := res.Err(); err != nil && err != mongodriver.ErrNoDocuments {
				return err
			}
		}
	}

	return nil
}

// ResetStats saves empty stats
func (s *statsStore) ResetStats(ctx context.Context, eid string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	res := s.dbCollection.FindOneAndUpdate(ctx, bson.M{"_id": eid}, bson.M{
		"$set": bson.M{"ko": 0, "ok": 0},
	}, options.FindOneAndUpdate().SetUpsert(true))
	if err := res.Err(); err != nil && err != mongodriver.ErrNoDocuments {
		return err
	}
	return nil
}
