package widget

import (
	"context"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	FindViewIds(ctx context.Context, ids []string) (map[string]string, error)
	FindViewIdByTab(ctx context.Context, tabId string) (string, error)
	GetOneBy(ctx context.Context, id string) (*Response, error)
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
	Copy(ctx context.Context, widget Response, r EditRequest) (*Response, error)
	UpdateGridPositions(ctx context.Context, items []EditGridPositionItemRequest) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		client:             dbClient,
		collection:         dbClient.Collection(mongo.WidgetMongoCollection),
		tabCollection:      dbClient.Collection(mongo.ViewTabMongoCollection),
		filterCollection:   dbClient.Collection(mongo.WidgetFiltersMongoCollection),
		userPrefCollection: dbClient.Collection(mongo.UserPreferencesMongoCollection),
	}
}

type store struct {
	client             mongo.DbClient
	collection         mongo.DbCollection
	tabCollection      mongo.DbCollection
	filterCollection   mongo.DbCollection
	userPrefCollection mongo.DbCollection
}

func (s *store) FindViewIds(ctx context.Context, ids []string) (map[string]string, error) {
	results := make([]struct {
		ID   string `bson:"_id"`
		View string `bson:"view"`
	}, 0)
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": ids}}},
		{"$lookup": bson.M{
			"from":         mongo.ViewTabMongoCollection,
			"localField":   "tab",
			"foreignField": "_id",
			"as":           "tab",
		}},
		{"$unwind": bson.M{"path": "$tab"}},
		{"$project": bson.M{
			"view": "$tab.view",
		}},
	})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	viewIds := make(map[string]string)
	for _, result := range results {
		if result.View != "" {
			viewIds[result.ID] = result.View
		}
	}

	return viewIds, nil
}

func (s *store) FindViewIdByTab(ctx context.Context, tabId string) (string, error) {
	result := struct {
		View string `bson:"view"`
	}{}
	err := s.tabCollection.FindOne(ctx, bson.M{"_id": tabId}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return "", nil
		}
		return "", err
	}

	return result.View, nil
}

func (s *store) GetOneBy(ctx context.Context, id string) (*Response, error) {
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$lookup": bson.M{
			"from":         mongo.WidgetFiltersMongoCollection,
			"localField":   "_id",
			"foreignField": "widget",
			"as":           "filters",
		}},
		{"$unwind": bson.M{"path": "$filters", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"filters.title": 1}},
		{"$group": bson.M{
			"_id":     nil,
			"data":    bson.M{"$first": "$$ROOT"},
			"filters": bson.M{"$push": "$filters"},
		}},
		{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$data",
			bson.M{"filters": bson.M{"$filter": bson.M{
				"input": "$filters",
				"cond":  bson.M{"$eq": bson.A{"$$this.is_private", false}},
			}}},
		}}}},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		widget := Response{}
		err = cursor.Decode(&widget)
		if err != nil {
			return nil, err
		}

		return &widget, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r EditRequest) (*Response, error) {
	now := types.NewCpsTime()
	widget := transformEditRequestToModel(r)
	widget.ID = utils.NewID()
	widget.Created = now
	widget.Updated = now

	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err := s.collection.InsertOne(ctx, widget)
		if err != nil {
			return err
		}

		response, err = s.GetOneBy(ctx, widget.ID)
		return err
	})

	return response, err
}

func (s *store) Update(ctx context.Context, r EditRequest) (*Response, error) {
	oldWidget, err := s.GetOneBy(ctx, r.ID)
	if err != nil || oldWidget == nil {
		return nil, err
	}

	now := types.NewCpsTime()
	widget := transformEditRequestToModel(r)
	widget.ID = oldWidget.ID
	widget.Updated = now
	// Empty InternalParameters to remove from update query.
	widget.InternalParameters = view.InternalParameters{}
	update := bson.M{"$set": widget}

	if oldWidget.Type == view.WidgetTypeJunit &&
		(widget.Type != oldWidget.Type ||
			widget.Parameters.IsAPI != oldWidget.Parameters.IsAPI ||
			widget.Parameters.Directory != oldWidget.Parameters.Directory ||
			widget.Parameters.ReportFileRegexp != oldWidget.Parameters.ReportFileRegexp) {
		update["$unset"] = bson.M{"internal_parameters": ""}
	}

	var response *Response
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err = s.collection.UpdateOne(ctx, bson.M{"_id": widget.ID}, update)
		if err != nil {
			return err
		}

		response, err = s.GetOneBy(ctx, widget.ID)
		return err
	})

	return response, err
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	res := false
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = false
		delCount, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil || delCount == 0 {
			return err
		}

		err = s.deleteUserPreferences(ctx, id)
		if err != nil {
			return err
		}

		err = s.deleteFilters(ctx, id)
		if err != nil {
			return err
		}

		res = true
		return nil
	})

	return res, err
}

func (s *store) Copy(ctx context.Context, widget Response, r EditRequest) (*Response, error) {
	now := types.NewCpsTime()
	newWidget := view.Widget{
		ID:             utils.NewID(),
		Tab:            r.Tab,
		Title:          r.Title,
		Type:           r.Type,
		GridParameters: r.GridParameters,
		Parameters:     r.Parameters,
		Author:         r.Author,
		Created:        now,
		Updated:        now,
	}

	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		cursor, err := s.filterCollection.Find(ctx, bson.M{
			"widget":          widget.ID,
			"is_private":      false,
			"old_mongo_query": nil, //do not copy old filters
		})
		if err != nil {
			return err
		}
		defer cursor.Close(ctx)

		mainFilter := ""
		filters := make([]interface{}, 0)
		for cursor.Next(ctx) {
			filter := view.WidgetFilter{}
			err := cursor.Decode(&filter)
			if err != nil {
				return err
			}

			newId := utils.NewID()
			// Main filter can be old filter so keep main filter in this case.
			if newWidget.Parameters.MainFilter == filter.ID {
				mainFilter = newId
			}

			filter.ID = newId
			filter.Widget = newWidget.ID
			filter.Author = r.Author
			filter.Created = now
			filter.Updated = now
			filters = append(filters, filter)
		}

		newWidget.Parameters.MainFilter = mainFilter
		_, err = s.collection.InsertOne(ctx, newWidget)
		if err != nil {
			return err
		}

		if len(filters) > 0 {
			_, err := s.filterCollection.InsertMany(ctx, filters)
			if err != nil {
				return err
			}
		}

		response, err = s.GetOneBy(ctx, newWidget.ID)
		return err
	})

	return response, err
}

func (s *store) UpdateGridPositions(ctx context.Context, items []EditGridPositionItemRequest) (bool, error) {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}

	res := false
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = false
		widgets := make([]view.Widget, 0, len(items))
		cursor, err := s.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
		if err != nil {
			return err
		}

		err = cursor.All(ctx, &widgets)
		if err != nil || len(widgets) != len(items) {
			return err
		}

		tabId := ""
		for _, w := range widgets {
			if tabId == "" {
				tabId = w.Tab
			} else if tabId != w.Tab {
				return ValidationErr{error: errors.New("widgets are related to different view tabs")}
			}
		}

		count, err := s.collection.CountDocuments(ctx, bson.M{"tab": tabId})
		if err != nil {
			return err
		}
		if count != int64(len(items)) {
			return ValidationErr{error: errors.New("widgets are missing")}
		}

		writeModels := make([]mongodriver.WriteModel, len(widgets))
		for i, item := range items {
			writeModels[i] = mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": item.ID}).
				SetUpdate(bson.M{"$set": bson.M{"grid_parameters": item.GridParameters}})
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

func (s *store) deleteUserPreferences(ctx context.Context, widgetID string) error {
	_, err := s.userPrefCollection.DeleteMany(ctx, bson.M{
		"widget": widgetID,
	})

	return err
}

func (s *store) deleteFilters(ctx context.Context, widgetID string) error {
	_, err := s.filterCollection.DeleteMany(ctx, bson.M{
		"widget": widgetID,
	})

	return err
}

func transformEditRequestToModel(request EditRequest) view.Widget {
	return view.Widget{
		Tab:            request.Tab,
		Title:          request.Title,
		Type:           request.Type,
		GridParameters: request.GridParameters,
		Parameters:     request.Parameters,
		Author:         request.Author,
	}
}
