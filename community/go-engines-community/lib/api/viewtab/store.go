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
	Find(ctx context.Context, ids []string) ([]view.Tab, error)
	GetOneBy(ctx context.Context, id string) (*view.Tab, error)
	Insert(ctx context.Context, r EditRequest) (*view.Tab, error)
	Update(ctx context.Context, oldTab view.Tab, r EditRequest) (*view.Tab, error)
	Delete(ctx context.Context, id string) (bool, error)
	Copy(ctx context.Context, tab view.Tab, r CopyRequest) (*view.Tab, error)
	UpdatePositions(ctx context.Context, tabs []view.Tab) (bool, error)
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

func (s *store) Find(ctx context.Context, ids []string) ([]view.Tab, error) {
	tabs := make([]view.Tab, 0)
	cursor, err := s.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &tabs)
	if err != nil {
		return nil, err
	}

	tabsByID := make(map[string]view.Tab, len(tabs))
	for _, tab := range tabs {
		tabsByID[tab.ID] = tab
	}

	sortedTabs := make([]view.Tab, 0, len(tabs))
	for _, id := range ids {
		if tab, ok := tabsByID[id]; ok {
			sortedTabs = append(sortedTabs, tab)
		}
	}

	return sortedTabs, nil
}

func (s *store) GetOneBy(ctx context.Context, id string) (*view.Tab, error) {
	model := view.Tab{}
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&model)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &model, nil
}

func (s *store) Insert(ctx context.Context, r EditRequest) (*view.Tab, error) {
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

	return &tab, nil
}

func (s *store) Update(ctx context.Context, oldTab view.Tab, r EditRequest) (*view.Tab, error) {
	now := types.CpsTime{Time: time.Now()}
	tab := view.Tab{
		ID:       oldTab.ID,
		Title:    r.Title,
		View:     r.View,
		Author:   r.Author,
		Position: oldTab.Position,
		Created:  oldTab.Created,
		Updated:  now,
	}

	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": tab.ID}, bson.M{"$set": tab})
	if err != nil {
		return nil, err
	}

	return &tab, nil
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

func (s *store) Copy(ctx context.Context, tab view.Tab, r CopyRequest) (*view.Tab, error) {
	count, err := s.collection.CountDocuments(ctx, bson.M{"view": r.View})
	if err != nil {
		return nil, err
	}
	id := tab.ID
	now := types.CpsTime{Time: time.Now()}
	tab.ID = utils.NewID()
	tab.View = r.View
	tab.Author = r.Author
	tab.Position = count
	tab.Created = now
	tab.Updated = now

	_, err = s.collection.InsertOne(ctx, tab)
	if err != nil {
		return nil, err
	}

	cursor, err := s.widgetCollection.Find(ctx, bson.M{"tab": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		w := view.Widget{}
		err := cursor.Decode(&w)
		if err != nil {
			return nil, err
		}

		_, err = s.widgetStore.Copy(ctx, w, widget.CopyRequest{
			Tab:    tab.ID,
			Author: r.Author,
		})
		if err != nil {
			return nil, err
		}
	}

	return &tab, nil
}

func (s *store) UpdatePositions(ctx context.Context, tabs []view.Tab) (bool, error) {
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
