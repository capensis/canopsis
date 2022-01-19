package viewtab

import (
	"context"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widget"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Store interface {
	Find(ctx context.Context, ids []string) ([]Tab, error)
	GetOneBy(ctx context.Context, id string) (*Tab, error)
	Insert(ctx context.Context, r EditRequest) (*Tab, error)
	Update(ctx context.Context, oldTab Tab, r EditRequest) (*Tab, error)
	Delete(ctx context.Context, id string) (bool, error)
	Copy(ctx context.Context, tab Tab, r CopyRequest) (*Tab, error)
	UpdatePositions(ctx context.Context, tabs []Tab) (bool, error)
}

func NewStore(dbClient mongo.DbClient, widgetStore widget.Store) Store {
	return &store{
		collection:         dbClient.Collection(mongo.ViewTabMongoCollection),
		widgetCollection:   dbClient.Collection(mongo.WidgetMongoCollection),
		playlistCollection: dbClient.Collection(mongo.PlaylistMongoCollection),
		userPrefCollection: dbClient.Collection(mongo.UserPreferencesMongoCollection),

		widgetStore: widgetStore,
	}
}

type store struct {
	collection         mongo.DbCollection
	widgetCollection   mongo.DbCollection
	playlistCollection mongo.DbCollection
	userPrefCollection mongo.DbCollection

	widgetStore widget.Store
}

func (s *store) Find(ctx context.Context, ids []string) ([]Tab, error) {
	tabs := make([]Tab, 0)
	cursor, err := s.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &tabs)
	if err != nil {
		return nil, err
	}

	tabsByID := make(map[string]Tab, len(tabs))
	for _, tab := range tabs {
		tabsByID[tab.ID] = tab
	}

	sortedTabs := make([]Tab, 0, len(tabs))
	for _, id := range ids {
		if tab, ok := tabsByID[id]; ok {
			sortedTabs = append(sortedTabs, tab)
		}
	}

	return sortedTabs, nil
}

func (s *store) GetOneBy(ctx context.Context, id string) (*Tab, error) {
	tabs := make([]Tab, 0)
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$lookup": bson.M{
			"from":         mongo.WidgetMongoCollection,
			"localField":   "_id",
			"foreignField": "tab",
			"as":           "widgets",
		}},
		{"$unwind": bson.M{"path": "$widgets", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.WidgetFiltersMongoCollection,
			"localField":   "widgets._id",
			"foreignField": "widget",
			"as":           "filters",
		}},
		{"$unwind": bson.M{"path": "$filters", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"filters.user": bson.M{"$cond": bson.M{
				"if":   "$filters.user",
				"then": "$filters.user",
				"else": "",
			}},
		}},
		{"$sort": bson.M{"filters.title": 1}},
		{"$group": bson.M{
			"_id": bson.M{
				"_id":     "$_id",
				"widgets": "$widgets._id",
			},
			"data":    bson.M{"$first": "$$ROOT"},
			"widgets": bson.M{"$first": "$widgets"},
			"filters": bson.M{"$push": "$filters"},
		}},
		{"$addFields": bson.M{
			"widgets.filters": bson.M{"$filter": bson.M{
				"input": bson.M{"$filter": bson.M{"input": "$filters", "cond": "$$this._id"}},
				"cond":  bson.M{"$eq": bson.A{"$$this.user", ""}},
			}},
		}},
		{"$sort": bson.D{{"widgets.grid_parameters.desktop.y", 1}, {"widgets.grid_parameters.desktop.x", 1}}},
		{"$group": bson.M{
			"_id":     "$_id._id",
			"data":    bson.M{"$first": "$data"},
			"widgets": bson.M{"$push": "$widgets"},
		}},
		{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$data",
			bson.M{"widgets": bson.M{"$filter": bson.M{
				"input": "$widgets",
				"cond":  "$$this._id",
			}}},
		}}}},
	})
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

func (s *store) Insert(ctx context.Context, r EditRequest) (*Tab, error) {
	count, err := s.collection.CountDocuments(ctx, bson.M{"view": r.View})
	if err != nil {
		return nil, err
	}
	now := types.CpsTime{Time: time.Now()}
	tab := view.Tab{
		ID:       utils.NewID(),
		Title:    r.Title,
		View:     r.View,
		Author:   r.Author,
		Position: count,
		Created:  now,
		Updated:  now,
	}

	_, err = s.collection.InsertOne(ctx, tab)
	if err != nil {
		return nil, err
	}

	return s.GetOneBy(ctx, tab.ID)
}

func (s *store) Update(ctx context.Context, oldTab Tab, r EditRequest) (*Tab, error) {
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

	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": tab.ID}, bson.M{"$set": tab})
	if err != nil {
		return nil, err
	}

	return s.GetOneBy(ctx, tab.ID)
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	isLinked, err := s.isLinked(ctx, id)
	if err != nil {
		return false, err
	}
	if isLinked {
		return false, ValidationErr{error: errors.New("view tab is linked to playlist")}
	}

	delCount, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	if delCount == 0 {
		return false, nil
	}

	err = s.deleteWidgets(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) Copy(ctx context.Context, tab Tab, r CopyRequest) (*Tab, error) {
	count, err := s.collection.CountDocuments(ctx, bson.M{"view": r.View})
	if err != nil {
		return nil, err
	}
	now := types.CpsTime{Time: time.Now()}
	newTab := view.Tab{
		ID:       utils.NewID(),
		Title:    tab.Title,
		View:     r.View,
		Author:   r.Author,
		Position: count,
		Created:  now,
		Updated:  now,
	}

	_, err = s.collection.InsertOne(ctx, newTab)
	if err != nil {
		return nil, err
	}

	cursor, err := s.widgetCollection.Find(ctx, bson.M{"tab": tab.ID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		w := widget.Widget{}
		err := cursor.Decode(&w)
		if err != nil {
			return nil, err
		}

		_, err = s.widgetStore.Copy(ctx, w, widget.CopyRequest{
			Tab:    newTab.ID,
			Author: r.Author,
		})
		if err != nil {
			return nil, err
		}
	}

	return s.GetOneBy(ctx, newTab.ID)
}

func (s *store) UpdatePositions(ctx context.Context, tabs []Tab) (bool, error) {
	viewId := ""
	for _, tab := range tabs {
		if viewId == "" {
			viewId = tab.View
		} else if viewId != tab.View {
			return false, ValidationErr{error: errors.New("view tabs are related to different views")}
		}
	}

	count, err := s.collection.CountDocuments(ctx, bson.M{"view": viewId})
	if err != nil {
		return false, err
	}
	if count != int64(len(tabs)) {
		return false, ValidationErr{error: errors.New("view tabs are missing")}
	}

	writeModels := make([]mongodriver.WriteModel, len(tabs))
	for i, tab := range tabs {
		writeModels[i] = mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": tab.ID}).
			SetUpdate(bson.M{"$set": bson.M{"position": i}})
	}

	res, err := s.collection.BulkWrite(ctx, writeModels)
	if err != nil {
		return false, err
	}

	return res.MatchedCount > 0, nil
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

	return nil
}
