package alarmtag

import (
	"context"
	"errors"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExternalUpdater interface {
	Add(tags map[string]string)
	Update(ctx context.Context) error
}

func NewExternalUpdater(client mongo.DbClient) ExternalUpdater {
	return &externalUpdater{
		client:                  client,
		alarmTagCollection:      client.Collection(mongo.AlarmTagCollection),
		alarmTagColorCollection: client.Collection(mongo.AlarmTagColorCollection),
		configCollection:        client.Collection(mongo.ConfigurationMongoCollection),
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		tags:                    make(map[string]string),
	}
}

type externalUpdater struct {
	client                  mongo.DbClient
	alarmTagColorCollection mongo.DbCollection
	alarmTagCollection      mongo.DbCollection
	configCollection        mongo.DbCollection
	alarmCollection         mongo.DbCollection

	tagsMx sync.Mutex
	tags   map[string]string
}

func (u *externalUpdater) Add(eventTags map[string]string) {
	u.tagsMx.Lock()
	defer u.tagsMx.Unlock()
	for k, v := range eventTags {
		t := types.TransformEventTag(k, v)
		if t != "" {
			u.tags[t] = k
		}
	}
}

func (u *externalUpdater) Update(ctx context.Context) error {
	u.tagsMx.Lock()
	tags := u.tags
	u.tags = make(map[string]string, len(u.tags))
	u.tagsMx.Unlock()

	return u.update(ctx, tags)
}

func (u *externalUpdater) update(ctx context.Context, tags map[string]string) error {
	if len(tags) == 0 {
		return nil
	}

	colors, err := u.getColors(ctx)
	if err != nil {
		return err
	}

	now := datetime.NewCpsTime()

	return u.client.WithTransaction(ctx, func(ctx context.Context) error {
		internalTagsToRemove, err := u.keepNewTags(ctx, tags)
		if err != nil || len(tags) == 0 {
			return err
		}

		err = u.removeInternalTags(ctx, internalTagsToRemove)
		if err != nil {
			return err
		}

		labelColors, err := u.getLabelColors(ctx, tags)
		if err != nil {
			return err
		}

		count, err := u.getLabelsCount(ctx)
		if err != nil {
			return err
		}

		models := make([]interface{}, 0, len(tags))
		k := 0
		for t, label := range tags {
			color := labelColors[label]
			if color == "" && len(colors) > 0 {
				colorIndex := (count + k) % len(colors)
				color = colors[colorIndex]

				v := struct {
					Color string `bson:"color"`
				}{}

				err = u.alarmTagColorCollection.FindOneAndUpdate(
					ctx,
					bson.M{"_id": label},
					bson.M{
						"$setOnInsert": bson.M{
							"color": color,
						},
					},
					options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
				).Decode(&v)
				if err != nil {
					return err
				}

				labelColors[label] = v.Color
				color = v.Color

				k++
			}

			models = append(models, AlarmTag{
				ID:      utils.NewID(),
				Type:    TypeExternal,
				Value:   t,
				Label:   label,
				Color:   color,
				Created: now,
				Updated: now,
			})
		}

		_, err = u.alarmTagCollection.InsertMany(ctx, models)
		return err
	})
}

func (u *externalUpdater) keepNewTags(ctx context.Context, tags map[string]string) ([]string, error) {
	values := make([]string, 0, len(tags))
	for t := range tags {
		values = append(values, t)
	}

	internalTagsToRemove := make([]string, 0)
	cursor, err := u.alarmTagCollection.Find(ctx, bson.M{"value": bson.M{"$in": values}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		tag := struct {
			Type  int64  `bson:"type"`
			Value string `bson:"value"`
		}{}

		err = cursor.Decode(&tag)
		if err != nil {
			return nil, err
		}

		switch tag.Type {
		case TypeExternal:
			delete(tags, tag.Value)
		case TypeInternal:
			internalTagsToRemove = append(internalTagsToRemove, tag.Value)
		}
	}

	return internalTagsToRemove, nil
}

func (u *externalUpdater) getLabelColors(ctx context.Context, tags map[string]string) (map[string]string, error) {
	labels := make([]string, 0, len(tags))
	for _, label := range tags {
		labels = append(labels, label)
	}

	cursor, err := u.alarmTagColorCollection.Find(ctx, bson.M{"_id": bson.M{"$in": labels}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	colors := make(map[string]string, len(labels))
	for cursor.Next(ctx) {
		v := struct {
			Label string `bson:"_id"`
			Color string `bson:"color"`
		}{}

		err = cursor.Decode(&v)
		if err != nil {
			return nil, err
		}

		colors[v.Label] = v.Color
	}

	return colors, nil
}

func (u *externalUpdater) getLabelsCount(ctx context.Context) (int, error) {
	c, err := u.alarmTagColorCollection.CountDocuments(ctx, bson.M{})

	return int(c), err
}

func (u *externalUpdater) getColors(ctx context.Context) ([]string, error) {
	v := struct {
		Colors []string `bson:"colors"`
	}{}
	err := u.configCollection.FindOne(ctx, bson.M{"_id": config.AlarmTagColorKeyName}).Decode(&v)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return v.Colors, nil
}

func (u *externalUpdater) removeInternalTags(ctx context.Context, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	_, err := u.alarmTagCollection.DeleteMany(ctx, bson.M{
		"type":  TypeInternal,
		"value": bson.M{"$in": tags},
	})
	if err != nil {
		return err
	}

	match := make([]bson.M, len(tags))
	unset := bson.M{}
	for i, tag := range tags {
		match[i] = bson.M{"itags." + tag: bson.M{"$ne": nil}}
		unset["itags."+tag] = ""
	}
	cursor, err := u.alarmCollection.Find(ctx, bson.M{
		"v.resolved": nil,
		"$or":        match,
	}, options.Find().SetProjection(bson.M{
		"_id":   1,
		"etags": 1,
	}))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	writeModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)
	for cursor.Next(ctx) {
		var alarm types.Alarm
		err = cursor.Decode(&alarm)
		if err != nil {
			return err
		}

		internalTags := make([]string, 0, len(tags))
		for _, tag := range tags {
			found := false
			for _, externalTag := range alarm.ExternalTags {
				if tag == externalTag {
					found = true
					break
				}
			}
			if !found {
				internalTags = append(internalTags, tag)
			}
		}

		writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{
				"_id":        alarm.ID,
				"v.resolved": nil,
			}).
			SetUpdate(bson.M{
				"$pullAll": bson.M{"tags": internalTags},
				"$unset":   unset,
			}))
		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err = u.alarmCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
		}
	}

	if len(writeModels) > 0 {
		_, err = u.alarmCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}
	}

	return nil
}
