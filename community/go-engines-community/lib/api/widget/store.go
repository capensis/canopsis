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
	CopyForTab(ctx context.Context, tabID, newTabID, author string) error
	UpdateGridPositions(ctx context.Context, items []EditGridPositionItemRequest) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		client:             dbClient,
		collection:         dbClient.Collection(mongo.WidgetMongoCollection),
		tabCollection:      dbClient.Collection(mongo.ViewTabMongoCollection),
		userPrefCollection: dbClient.Collection(mongo.UserPreferencesMongoCollection),
	}
}

type store struct {
	client             mongo.DbClient
	collection         mongo.DbCollection
	tabCollection      mongo.DbCollection
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
	widget := Response{}
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&widget)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &widget, nil
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
		if err != nil || response == nil {
			return err
		}

		err = s.updateUserPreferences(ctx, *response)
		if err != nil {
			return err
		}

		return nil
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

		res = true
		return nil
	})

	return res, err
}

func (s *store) Copy(ctx context.Context, widget Response, r EditRequest) (*Response, error) {
	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		var err error
		response, err = s.copy(ctx, r)
		return err
	})

	return response, err
}

func (s *store) CopyForTab(ctx context.Context, tabID, newTabID, author string) error {
	cursor, err := s.collection.Find(ctx, bson.M{"tab": tabID})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		w := Response{}
		err := cursor.Decode(&w)
		if err != nil {
			return err
		}

		_, err = s.copy(ctx, EditRequest{
			Tab:            newTabID,
			Title:          w.Title,
			Type:           w.Type,
			GridParameters: w.GridParameters,
			Parameters:     w.Parameters,
			Author:         author,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) copy(ctx context.Context, r EditRequest) (*Response, error) {
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

	_, err := s.collection.InsertOne(ctx, newWidget)
	if err != nil {
		return nil, err
	}

	return s.GetOneBy(ctx, newWidget.ID)
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

func (s *store) updateUserPreferences(ctx context.Context, response Response) error {
	if len(response.Parameters.ViewFilters) == 0 {
		return nil
	}

	cursor, err := s.userPrefCollection.Find(ctx, bson.M{
		"widget": response.ID,
	})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	writeModels := make([]mongodriver.WriteModel, 0)
	for cursor.Next(ctx) {
		userPref := struct {
			ID      string `bson:"_id"`
			Content struct {
				MainFilter interface{} `bson:"mainFilter"`
			} `bson:"content" json:"content"`
		}{}
		err = cursor.Decode(&userPref)
		if err != nil {
			return err
		}

		if userPref.Content.MainFilter == nil {
			continue
		}

		if filters, ok := userPref.Content.MainFilter.(bson.A); ok {
			newFilters := make([]interface{}, len(filters))
			updated := false
			for i, filter := range filters {
				newFilter := adjustUserPreferenceFilter(filter, response.Parameters.ViewFilters)
				if newFilter.Title == "" {
					newFilters[i] = filter
				} else {
					newFilters[i] = newFilter
					updated = true
				}
			}

			if updated {
				writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
					SetFilter(bson.M{"_id": userPref.ID}).
					SetUpdate(bson.M{"$set": bson.M{
						"content.mainFilter": newFilters,
					}}),
				)
			}
		} else {
			newFilter := adjustUserPreferenceFilter(userPref.Content.MainFilter, response.Parameters.ViewFilters)
			if newFilter.Title != "" {
				writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
					SetFilter(bson.M{"_id": userPref.ID}).
					SetUpdate(bson.M{"$set": bson.M{
						"content.mainFilter": newFilter,
					}}),
				)
			}
		}
	}

	if len(writeModels) > 0 {
		_, err = s.userPrefCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}
	}

	return nil
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

func adjustUserPreferenceFilter(filter interface{}, viewFilters []view.Filter) view.Filter {
	var title, value interface{}
	if d, ok := filter.(bson.D); ok {
		for _, e := range d {
			switch e.Key {
			case "title":
				title = e.Value
			case "filter":
				value = e.Value
			}
		}
	}

	if title != nil {
		for _, viewFilter := range viewFilters {
			if title == viewFilter.Title {
				if value != viewFilter.Filter {
					return viewFilter
				}

				break
			}
		}
	}

	return view.Filter{}
}
