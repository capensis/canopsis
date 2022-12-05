package alarmtag

import (
	"context"
	"errors"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Updater interface {
	Add(tags map[string]string)
	Update(ctx context.Context) error
}

func NewUpdater(client mongo.DbClient) Updater {
	return &updater{
		client:           client,
		collection:       client.Collection(mongo.AlarmTagCollection),
		configCollection: client.Collection(mongo.ConfigurationMongoCollection),
		tags:             make(map[string]string),
	}
}

type updater struct {
	client           mongo.DbClient
	collection       mongo.DbCollection
	configCollection mongo.DbCollection

	tagsMx sync.Mutex
	tags   map[string]string
}

func (u *updater) Add(eventTags map[string]string) {
	u.tagsMx.Lock()
	defer u.tagsMx.Unlock()
	for k, v := range eventTags {
		t := types.TransformEventTag(k, v)
		if t != "" {
			u.tags[t] = k
		}
	}
}

func (u *updater) Update(ctx context.Context) error {
	u.tagsMx.Lock()
	tags := u.tags
	u.tags = make(map[string]string, len(u.tags))
	u.tagsMx.Unlock()

	return u.update(ctx, tags)
}

func (u *updater) update(ctx context.Context, tags map[string]string) error {
	if len(tags) == 0 {
		return nil
	}

	colors, err := u.getColors(ctx)
	if err != nil {
		return err
	}

	now := types.NewCpsTime()

	return u.client.WithTransaction(ctx, func(ctx context.Context) error {
		err := u.keepNewTags(ctx, tags)
		if err != nil || len(tags) == 0 {
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

		models := make([]interface{}, len(tags))
		i := 0
		k := 0
		for t, label := range tags {
			color := labelColors[label]
			if color == "" && len(colors) > 0 {
				colorIndex := (count + k) % len(colors)
				color = colors[colorIndex]
				k++
			}

			models[i] = types.AlarmTag{
				ID:      utils.NewID(),
				Value:   t,
				Label:   label,
				Color:   color,
				Created: now,
			}
			i++
		}

		_, err = u.collection.InsertMany(ctx, models)
		return err
	})
}

func (u *updater) keepNewTags(ctx context.Context, tags map[string]string) error {
	values := make([]string, len(tags))
	i := 0
	for t := range tags {
		values[i] = t
		i++
	}
	cursor, err := u.collection.Find(ctx, bson.M{"value": bson.M{"$in": values}})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		tag := struct {
			Value string `bson:"value"`
		}{}
		err = cursor.Decode(&tag)
		if err != nil {
			return err
		}

		delete(tags, tag.Value)
	}

	return nil
}

func (u *updater) getLabelColors(ctx context.Context, newTags map[string]string) (map[string]string, error) {
	newLabels := make([]string, len(newTags))
	i := 0
	for _, label := range newTags {
		newLabels[i] = label
		i++
	}

	cursor, err := u.collection.Find(ctx, bson.M{"label": bson.M{"$in": newLabels}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	colors := make(map[string]string, len(newLabels))
	for cursor.Next(ctx) {
		v := struct {
			Label string `bson:"label"`
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

func (u *updater) getLabelsCount(ctx context.Context) (int, error) {
	cursor, err := u.collection.Aggregate(ctx, []bson.M{
		{"$group": bson.M{
			"_id": "$label",
		}},
		{"$group": bson.M{
			"_id":   nil,
			"count": bson.M{"$sum": 1},
		}},
	})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		v := struct {
			Count int `bson:"count"`
		}{}
		err = cursor.Decode(&v)
		return v.Count, err
	}

	return 0, nil
}

func (u *updater) getColors(ctx context.Context) ([]string, error) {
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
