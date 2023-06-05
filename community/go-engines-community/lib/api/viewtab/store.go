package viewtab

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widget"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Find(ctx context.Context, ids []string) ([]Response, error)
	GetOneBy(ctx context.Context, id string) (*Response, error)
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	Update(ctx context.Context, oldTab Response, r EditRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
	Copy(ctx context.Context, tab Response, r EditRequest) (*Response, error)
	CopyForView(ctx context.Context, viewID, newViewID, author string) error
	UpdatePositions(ctx context.Context, tabs []Response) (bool, error)
}

func NewStore(dbClient mongo.DbClient, widgetStore widget.Store) Store {
	return &store{
		client:             dbClient,
		collection:         dbClient.Collection(mongo.ViewTabMongoCollection),
		widgetCollection:   dbClient.Collection(mongo.WidgetMongoCollection),
		filterCollection:   dbClient.Collection(mongo.WidgetFiltersMongoCollection),
		playlistCollection: dbClient.Collection(mongo.PlaylistMongoCollection),
		userPrefCollection: dbClient.Collection(mongo.UserPreferencesMongoCollection),

		widgetStore: widgetStore,
	}
}

type store struct {
	client             mongo.DbClient
	collection         mongo.DbCollection
	widgetCollection   mongo.DbCollection
	filterCollection   mongo.DbCollection
	playlistCollection mongo.DbCollection
	userPrefCollection mongo.DbCollection

	widgetStore widget.Store
}

func (s *store) Find(ctx context.Context, ids []string) ([]Response, error) {
	tabs := make([]Response, 0)
	pipeline := []bson.M{{"$match": bson.M{"_id": bson.M{"$in": ids}}}}
	pipeline = append(pipeline, author.Pipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &tabs)
	if err != nil {
		return nil, err
	}

	tabsByID := make(map[string]Response, len(tabs))
	for _, tab := range tabs {
		tabsByID[tab.ID] = tab
	}

	sortedTabs := make([]Response, 0, len(tabs))
	for _, id := range ids {
		if tab, ok := tabsByID[id]; ok {
			sortedTabs = append(sortedTabs, tab)
		}
	}

	return sortedTabs, nil
}

func (s *store) GetOneBy(ctx context.Context, id string) (*Response, error) {
	tabs := make([]Response, 0)
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$lookup": bson.M{
			"from":         mongo.WidgetMongoCollection,
			"localField":   "_id",
			"foreignField": "tab",
			"as":           "widgets",
		}},
		{"$unwind": bson.M{"path": "$widgets", "preserveNullAndEmptyArrays": true}},
	}
	pipeline = append(pipeline, author.PipelineForField("widgets.author")...)
	pipeline = append(pipeline,
		bson.M{"$lookup": bson.M{
			"from":         mongo.WidgetFiltersMongoCollection,
			"localField":   "widgets._id",
			"foreignField": "widget",
			"as":           "filters",
		}},
		bson.M{"$unwind": bson.M{"path": "$filters", "preserveNullAndEmptyArrays": true}},
	)
	pipeline = append(pipeline, author.PipelineForField("filters.author")...)
	pipeline = append(pipeline,
		bson.M{"$sort": bson.M{"filters.position": 1}},
		bson.M{"$group": bson.M{
			"_id": bson.M{
				"_id":    "_id",
				"widget": "$widgets._id",
			},
			"data":    bson.M{"$first": "$$ROOT"},
			"widgets": bson.M{"$first": "$widgets"},
			"filters": bson.M{"$push": "$filters"},
		}},
		bson.M{"$addFields": bson.M{
			"_id": "$_id._id",
			"widgets.filters": bson.M{"$filter": bson.M{
				"input": "$filters",
				"cond":  bson.M{"$eq": bson.A{"$$this.is_private", false}},
			}},
		}},
		bson.M{"$sort": bson.D{
			{Key: "widgets.grid_parameters.desktop.y", Value: 1},
			{Key: "widgets.grid_parameters.desktop.x", Value: 1},
		}},
		bson.M{"$group": bson.M{
			"_id":     "$_id",
			"data":    bson.M{"$first": "$data"},
			"widgets": bson.M{"$push": "$widgets"},
		}},
		bson.M{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$data",
			bson.M{"widgets": bson.M{"$filter": bson.M{
				"input": "$widgets",
				"cond":  "$$this._id",
			}}},
		}}}},
	)
	pipeline = append(pipeline, author.Pipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &tabs)
	if err != nil {
		return nil, err
	}

	if len(tabs) > 0 {
		return &tabs[0], nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r EditRequest) (*Response, error) {
	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		position, err := s.getNextPosition(ctx, r.View)
		if err != nil {
			return err
		}
		now := types.CpsTime{Time: time.Now()}
		tab := view.Tab{
			ID:       utils.NewID(),
			Title:    r.Title,
			View:     r.View,
			Author:   r.Author,
			Position: position,
			Created:  now,
			Updated:  now,
		}

		_, err = s.collection.InsertOne(ctx, tab)
		if err != nil {
			return err
		}

		response, err = s.GetOneBy(ctx, tab.ID)
		return err
	})

	return response, err
}

func (s *store) Update(ctx context.Context, oldTab Response, r EditRequest) (*Response, error) {
	now := types.CpsTime{Time: time.Now()}
	tab := view.Tab{
		ID:       oldTab.ID,
		Title:    r.Title,
		View:     r.View,
		Author:   r.Author,
		Position: oldTab.Position,
		Created:  *oldTab.Created,
		Updated:  now,
	}

	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err := s.collection.UpdateOne(ctx, bson.M{"_id": tab.ID}, bson.M{"$set": tab})
		if err != nil {
			return err
		}

		response, err = s.GetOneBy(ctx, tab.ID)
		return err
	})

	return response, err
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	res := false
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = false
		isLinked, err := s.isLinked(ctx, id)
		if err != nil {
			return err
		}
		if isLinked {
			return ValidationErr{error: errors.New("view tab is linked to playlist")}
		}

		delCount, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil || delCount == 0 {
			return err
		}

		err = s.deleteWidgets(ctx, id)
		if err != nil {
			return err
		}

		res = true
		return nil
	})

	return res, err
}

func (s *store) Copy(ctx context.Context, tab Response, r EditRequest) (*Response, error) {
	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		var err error
		response, err = s.copy(ctx, tab.ID, r)
		return err
	})

	return response, err
}

func (s *store) CopyForView(ctx context.Context, viewID, newViewID, author string) error {
	cursor, err := s.collection.Find(ctx, bson.M{"view": viewID}, options.Find().SetProjection(bson.M{"author": 0}))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		t := Response{}
		err := cursor.Decode(&t)
		if err != nil {
			return err
		}

		_, err = s.copy(ctx, t.ID, EditRequest{
			Title:  t.Title,
			View:   newViewID,
			Author: author,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) copy(ctx context.Context, tabID string, r EditRequest) (*Response, error) {
	position, err := s.getNextPosition(ctx, r.View)
	if err != nil {
		return nil, err
	}
	now := types.NewCpsTime()
	newTab := view.Tab{
		ID:       utils.NewID(),
		Title:    r.Title,
		View:     r.View,
		Author:   r.Author,
		Position: position,
		Created:  now,
		Updated:  now,
	}

	_, err = s.collection.InsertOne(ctx, newTab)
	if err != nil {
		return nil, err
	}

	err = s.widgetStore.CopyForTab(ctx, tabID, newTab.ID, r.Author)
	if err != nil {
		return nil, err
	}

	return s.GetOneBy(ctx, newTab.ID)
}

func (s *store) UpdatePositions(ctx context.Context, tabs []Response) (bool, error) {
	viewId := ""
	for _, tab := range tabs {
		if viewId == "" {
			viewId = tab.View
		} else if viewId != tab.View {
			return false, ValidationErr{error: errors.New("view tabs are related to different views")}
		}
	}

	res := false
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = false
		count, err := s.collection.CountDocuments(ctx, bson.M{"view": viewId})
		if err != nil {
			return err
		}
		if count != int64(len(tabs)) {
			return ValidationErr{error: errors.New("view tabs are missing")}
		}

		writeModels := make([]mongodriver.WriteModel, len(tabs))
		for i, tab := range tabs {
			writeModels[i] = mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": tab.ID}).
				SetUpdate(bson.M{"$set": bson.M{"position": i}})
		}

		writeRes, err := s.collection.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}

		res = writeRes.MatchedCount > 0
		return nil
	})

	return res, err
}

func (s *store) isLinked(ctx context.Context, id string) (bool, error) {
	err := s.playlistCollection.FindOne(ctx, bson.M{"tabs_list": id}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *store) deleteWidgets(ctx context.Context, id string) error {
	widgetCursor, err := s.widgetCollection.Find(ctx, bson.M{"tab": id})
	if err != nil {
		return err
	}
	widgets := make([]view.Widget, 0)
	err = widgetCursor.All(ctx, &widgets)
	if err != nil || len(widgets) == 0 {
		return err
	}

	widgetIds := make([]string, len(widgets))
	for i, w := range widgets {
		widgetIds[i] = w.ID
	}

	_, err = s.widgetCollection.DeleteMany(ctx, bson.M{"tab": id})
	if err != nil {
		return err
	}

	_, err = s.userPrefCollection.DeleteMany(ctx, bson.M{"widget": bson.M{"$in": widgetIds}})
	if err != nil {
		return err
	}

	_, err = s.filterCollection.DeleteMany(ctx, bson.M{"widget": bson.M{"$in": widgetIds}})
	if err != nil {
		return err
	}

	return nil
}

func (s *store) getNextPosition(ctx context.Context, view string) (int64, error) {
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"view": view}},
		{"$group": bson.M{
			"_id":      nil,
			"position": bson.M{"$max": "$position"},
		}},
	})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		data := struct {
			Position int64 `bson:"position"`
		}{}
		err = cursor.Decode(&data)
		return data.Position + 1, err
	}

	return 0, nil
}
