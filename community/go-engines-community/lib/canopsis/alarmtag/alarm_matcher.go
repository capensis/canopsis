package alarmtag

import (
	"context"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type InternalTagAlarmMatcher interface {
	Load(ctx context.Context) error
	Match(entity types.Entity, alarm types.Alarm) []string
}

func NewInternalTagAlarmMatcher(client mongo.DbClient) InternalTagAlarmMatcher {
	return &alarmMatcher{
		collection: client.Collection(mongo.AlarmTagCollection),
	}
}

type alarmMatcher struct {
	collection mongo.DbCollection

	tagsMx sync.RWMutex
	tags   []AlarmTag
}

func (m *alarmMatcher) Load(ctx context.Context) error {
	cursor, err := m.collection.Find(ctx, bson.M{"type": TypeInternal})
	if err != nil {
		return err
	}

	var tags []AlarmTag
	err = cursor.All(ctx, &tags)
	if err != nil {
		return err
	}

	m.tagsMx.Lock()
	defer m.tagsMx.Unlock()

	m.tags = tags
	return nil
}

func (m *alarmMatcher) Match(entity types.Entity, alarm types.Alarm) []string {
	m.tagsMx.RLock()
	defer m.tagsMx.RUnlock()

	matchedTags := make([]string, 0)
	for _, tag := range m.tags {
		ok, err := match.MatchEntityPattern(tag.EntityPattern, &entity)
		if err != nil || !ok {
			continue
		}

		ok, err = match.MatchAlarmPattern(tag.AlarmPattern, &alarm)
		if err != nil || !ok {
			continue
		}

		matchedTags = append(matchedTags, tag.Value)
	}

	return matchedTags
}
