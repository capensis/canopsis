package axe

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type internalTagPeriodicalWorker struct {
	PeriodicalInterval time.Duration
	Logger             zerolog.Logger
	TagCollection      mongo.DbCollection
	AlarmCollection    mongo.DbCollection
}

func (w *internalTagPeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *internalTagPeriodicalWorker) Work(ctx context.Context) {
	tags, err := w.loadTags(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot load tags")
		return
	}

	secNow := datetime.NewCpsTime()
	microNow := datetime.NewMicroTime()
	existedTags := make([]string, len(tags))
	for i, tag := range tags {
		existedTags[i] = tag.Value
		err = w.addTag(ctx, tag, tags, secNow, microNow)
		if err != nil {
			w.Logger.Err(err).Str("tag", tag.Value).Msg("cannot update alarms")
			return
		}
		err = w.removeTag(ctx, tag, tags, secNow, microNow)
		if err != nil {
			w.Logger.Err(err).Str("tag", tag.Value).Msg("cannot update alarms")
			return
		}
	}

	err = w.removeTags(ctx, existedTags, tags, secNow, microNow)
	if err != nil {
		w.Logger.Err(err).Msg("cannot update alarms")
		return
	}
}

func (w *internalTagPeriodicalWorker) loadTags(ctx context.Context) ([]alarmtag.AlarmTag, error) {
	cursor, err := w.TagCollection.Find(ctx, bson.M{"type": alarmtag.TypeInternal})
	if err != nil {
		return nil, err
	}

	var tags []alarmtag.AlarmTag
	err = cursor.All(ctx, &tags)
	return tags, err
}

func (w *internalTagPeriodicalWorker) addTag(
	ctx context.Context,
	tag alarmtag.AlarmTag,
	tags []alarmtag.AlarmTag,
	secNow datetime.CpsTime,
	microNow datetime.MicroTime,
) error {
	alarmMatch := bson.M{
		"t": bson.M{"$lt": secNow},
		"$or": []bson.M{
			{"itags_upd": bson.M{"$lt": microNow}},
			{"itags_upd": nil},
		},
		"itags": bson.M{"$nin": bson.A{tag.Value}},
	}
	if len(tag.AlarmPattern) > 0 {
		q, err := db.AlarmPatternToMongoQuery(tag.AlarmPattern, "")
		if err != nil {
			return err
		}
		alarmMatch = bson.M{"$and": []bson.M{alarmMatch, q}}
	}
	pipeline := []bson.M{
		{"$match": alarmMatch},
		{"$replaceRoot": bson.M{"newRoot": bson.M{"alarm": "$$ROOT"}}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "alarm.d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
	}
	if len(tag.EntityPattern) > 0 {
		q, err := db.EntityPatternToMongoQuery(tag.EntityPattern, "entity")
		if err != nil {
			return err
		}
		pipeline = append(pipeline, bson.M{"$match": q})
	}
	pipeline = append(pipeline, bson.M{"$project": bson.M{
		"alarm.v.steps": 0,
	}})

	cursor, err := w.AlarmCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	return w.updateByCursor(ctx, cursor, tags, microNow)
}

func (w *internalTagPeriodicalWorker) removeTag(
	ctx context.Context,
	tag alarmtag.AlarmTag,
	tags []alarmtag.AlarmTag,
	secNow datetime.CpsTime,
	microNow datetime.MicroTime,
) error {
	alarmMatch := bson.M{
		"t":         bson.M{"$lt": secNow},
		"itags_upd": bson.M{"$lt": microNow},
		"itags":     tag.Value,
	}
	if len(tag.AlarmPattern) > 0 && len(tag.EntityPattern) == 0 {
		q, err := db.AlarmPatternToNegativeMongoQuery(tag.AlarmPattern, "")
		if err != nil {
			return err
		}
		alarmMatch = bson.M{"$and": []bson.M{alarmMatch, q}}
	}
	pipeline := []bson.M{
		{"$match": alarmMatch},
		{"$replaceRoot": bson.M{"newRoot": bson.M{"alarm": "$$ROOT"}}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "alarm.d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
	}
	if len(tag.EntityPattern) > 0 {
		entityQuery, err := db.EntityPatternToNegativeMongoQuery(tag.EntityPattern, "entity")
		if err != nil {
			return err
		}
		var match bson.M
		if len(tag.AlarmPattern) > 0 {
			alarmQuery, err := db.AlarmPatternToNegativeMongoQuery(tag.AlarmPattern, "")
			if err != nil {
				return err
			}
			match = bson.M{"$or": []bson.M{
				alarmQuery,
				entityQuery,
			}}
		} else {
			match = entityQuery
		}
		pipeline = append(pipeline, bson.M{"$match": match})
	}
	pipeline = append(pipeline, bson.M{"$project": bson.M{
		"alarm.v.steps": 0,
	}})

	cursor, err := w.AlarmCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	return w.updateByCursor(ctx, cursor, tags, microNow)
}

func (w *internalTagPeriodicalWorker) removeTags(
	ctx context.Context,
	existedTags []string,
	tags []alarmtag.AlarmTag,
	secNow datetime.CpsTime,
	microNow datetime.MicroTime,
) error {
	cursor, err := w.AlarmCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"t":         bson.M{"$lt": secNow},
			"itags_upd": bson.M{"$lt": microNow},
			"itags": bson.M{"$elemMatch": bson.M{
				"$nin": existedTags,
			}},
		}},
		{"$replaceRoot": bson.M{"newRoot": bson.M{"alarm": "$$ROOT"}}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "alarm.d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$project": bson.M{
			"alarm.v.steps": 0,
		}},
	})
	if err != nil {
		return err
	}

	return w.updateByCursor(ctx, cursor, tags, microNow)
}

func (w *internalTagPeriodicalWorker) updateByCursor(
	ctx context.Context,
	cursor mongo.Cursor,
	tags []alarmtag.AlarmTag,
	microNow datetime.MicroTime,
) error {
	defer cursor.Close(ctx)

	writeModels := make([]mongodriver.WriteModel, 0)
	for cursor.Next(ctx) {
		alarm := types.AlarmWithEntity{}
		err := cursor.Decode(&alarm)
		if err != nil {
			return err
		}

		matchedTags := w.matchAlarm(alarm.Entity, alarm.Alarm, tags)
		writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{
				"_id": alarm.Alarm.ID,
				"$or": []bson.M{
					{"itags_upd": bson.M{"$lt": microNow}},
					{"itags_upd": nil},
				},
			}).
			SetUpdate([]bson.M{
				{"$set": bson.M{
					"itags_upd": microNow,
					"itags":     matchedTags,
					"tags": bson.M{"$concatArrays": bson.A{
						bson.M{"$cond": bson.M{
							"if":   "$etags",
							"then": "$etags",
							"else": bson.A{},
						}},
						matchedTags,
					}},
				}},
			}),
		)

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err = w.AlarmCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
		}
	}

	if len(writeModels) > 0 {
		_, err := w.AlarmCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *internalTagPeriodicalWorker) matchAlarm(entity types.Entity, alarm types.Alarm, tags []alarmtag.AlarmTag) []string {
	matchedTags := make([]string, 0)
	for _, tag := range tags {
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
