// stats contains implementation of http session statistics.
package stats

import (
	"context"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"time"
)

// Stats represents session statistics for http session by user.
type Stats struct {
	ID        string `bson:"_id" json:"_id"`
	SessionID string `bson:"session_id" json:"session_id"`
	// Deprecated
	OldSessionID string        `bson:"id_beaker_session" json:"id_beaker_session"`
	SessionStart types.CpsTime `bson:"start" json:"start"`
	UserID       string        `bson:"user_id" json:"user_id"`
	// Deprecated
	Username        string                 `bson:"username" json:"username"`
	LastPing        types.CpsTime          `bson:"last_ping" json:"last_ping"`
	LastVisiblePath interface{}            `bson:"last_visible_path,omitempty" json:"last_visible_path,omitempty"`
	LastVisiblePing types.CpsTime          `bson:"last_visible_ping" json:"last_visible_ping"`
	VisibleDuration int64                  `bson:"visible_duration" json:"visible_duration"`
	TabDuration     map[string]interface{} `bson:"tab_duration" json:"tab_duration"`
}

// SessionData represents required session info for stats.
type SessionData struct {
	SessionID, UserID string
}

// PathData represents required path info for stats.
type PathData struct {
	ViewID, TabID string
	Visible       bool
}

// Filter represents possible filters to find session stats.
type Filter struct {
	IsActive      *bool
	Usernames     []string
	StartedAfter  *types.CpsTime
	StoppedBefore *types.CpsTime
}

// Manager interfaces is used to implement stats store.
type Manager interface {
	// Ping updates stats if exists or creates new stats.
	Ping(SessionData, PathData) (*Stats, error)
	// Find returns stats list.
	Find(Filter) ([]Stats, error)
}

// manager saves stats to mongo db.
type manager struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
	frame        time.Duration
}

// NewManager creates new stats manager.
func NewManager(dbClient mongo.DbClient, frame time.Duration) Manager {
	return &manager{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.SessionStatsMongoCollection),
		frame:        frame,
	}
}

func (m *manager) Ping(data SessionData, pathData PathData) (*Stats, error) {
	if data.SessionID == "" || data.UserID == "" {
		return nil, errors.New("invalid session data")
	}

	now := types.CpsTime{Time: time.Now().In(time.UTC)}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cursor, err := m.dbCollection.Find(ctx, m.getActiveStatsFilter(data, now))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	s := Stats{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&s)
		if err != nil {
			return nil, err
		}
	} else {
		s = createStats(data, now)
	}

	if pathData.Visible {
		s.VisibleDuration += toSeconds(now.Sub(s.LastPing.Time))
		updateTabDuration(&s, pathData.ViewID, pathData.TabID)
		path := make([]string, 0, 2)
		if pathData.ViewID != "" {
			path = append(path, pathData.ViewID)
		}
		if pathData.TabID != "" {
			path = append(path, pathData.TabID)
		}
		s.LastVisiblePath = path
		s.LastVisiblePing = now
		s.LastPing = now
	} else {
		s.LastPing = now
	}

	_, err = m.dbCollection.UpdateOne(
		ctx,
		bson.M{"_id": s.ID},
		bson.M{"$set": s},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (m *manager) Find(f Filter) ([]Stats, error) {
	now := time.Now().In(time.UTC)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := make([]bson.M, 0)

	if f.IsActive != nil {
		lastPing := now.Unix() - toSeconds(m.frame)
		if *f.IsActive {
			filter = append(filter, bson.M{"last_ping": bson.M{"$gt": lastPing}})
		} else {
			filter = append(filter, bson.M{"last_ping": bson.M{"$lt": lastPing}})
		}
	}
	if f.StartedAfter != nil {
		filter = append(filter, bson.M{"start": bson.M{"$gt": f.StartedAfter.Unix()}})
	}
	if f.StoppedBefore != nil {
		filter = append(filter, bson.M{"last_ping": bson.M{"$lt": f.StoppedBefore.Unix()}})
	}

	pipeline := make([]bson.M, 0)

	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": filter}})
	}

	pipeline = append(
		pipeline,
		bson.M{"$lookup": bson.M{
			"from":         mongo.RightsMongoCollection,
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		}},
		bson.M{"$unwind": bson.M{"path": "$user", "preserveNullAndEmptyArrays": true}},
	)

	if len(f.Usernames) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$or": []bson.M{
			{"username": bson.M{"$in": f.Usernames}},
			{"user.crecord_name": bson.M{"$in": f.Usernames}},
		}}})
	}

	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"username": bson.M{"$cond": bson.M{"if": "$user.crecord_name", "then": "$user.crecord_name", "else": "$username"}},
	}})

	cursor, err := m.dbCollection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}

	var r []Stats
	err = cursor.All(ctx, &r)
	if err != nil {
		return nil, err
	}

	if r == nil {
		r = make([]Stats, 0)
	}

	return r, nil
}

// getActiveStatsFilter creates filter to find active stats for session.
func (m *manager) getActiveStatsFilter(data SessionData, now types.CpsTime) bson.M {
	lastPing := now.Unix() - toSeconds(m.frame)
	return bson.M{
		"session_id": data.SessionID,
		"last_ping":  bson.M{"$gt": lastPing},
	}
}

// createStats creates new stats with filled fields.
func createStats(data SessionData, now types.CpsTime) Stats {
	return Stats{
		ID:              utils.NewID(),
		SessionID:       data.SessionID,
		SessionStart:    now,
		UserID:          data.UserID,
		LastPing:        now,
		LastVisiblePath: []string{},
		LastVisiblePing: now,
		VisibleDuration: 0,
		TabDuration:     map[string]interface{}{},
	}
}

// updateTabDuration updates or adds path duration.
func updateTabDuration(s *Stats, viewID, tabID string) {
	if viewID == "" {
		return
	}

	now := time.Now()
	prevLastPing := s.LastPing
	prevLastVisiblePing := s.LastVisiblePing
	lastVisiblePath, _ := s.LastVisiblePath.([]string)
	var prevViewID, prevTabID string
	if len(lastVisiblePath) > 0 {
		prevViewID = lastVisiblePath[0]
	}
	if len(lastVisiblePath) > 1 {
		prevTabID = lastVisiblePath[1]
	}

	duration := toSeconds(now.Sub(prevLastPing.Time))
	if viewID == prevViewID && tabID == prevTabID {
		duration = toSeconds(now.Sub(prevLastVisiblePing.Time))
	}

	if s.TabDuration == nil {
		s.TabDuration = make(map[string]interface{})
	}

	if v, ok := s.TabDuration[viewID]; ok {
		if m, ok := v.(map[string]interface{}); ok && tabID != "" {
			if pd, ok := m[tabID]; ok {
				m[tabID] = toInt(pd) + duration
			} else {
				m[tabID] = duration
			}
		} else {
			s.TabDuration[viewID] = toInt(s.TabDuration[viewID]) + duration
		}
	} else if tabID == "" {
		s.TabDuration[viewID] = duration
	} else {
		s.TabDuration[viewID] = map[string]interface{}{tabID: duration}
	}
}

// toSeconds returns duration value in seconds.
func toSeconds(d time.Duration) int64 {
	return int64(math.Ceil(d.Seconds()))
}

// toInt casts var to int.
func toInt(v interface{}) int64 {
	if i, ok := v.(int); ok {
		return int64(i)
	}
	if i, ok := v.(int32); ok {
		return int64(i)
	}
	if i, ok := v.(int64); ok {
		return i
	}

	return 0
}
