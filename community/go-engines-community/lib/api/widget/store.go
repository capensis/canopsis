package widget

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	FindTabPrivacySettings(ctx context.Context, ids []string) (map[string]apisecurity.ViewTabPrivacySettings, error)
	FindViewIdByTab(ctx context.Context, tabId string) (string, error)
	GetOneBy(ctx context.Context, id string) (*Response, error)
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
	Copy(ctx context.Context, widgetID string, r CreateRequest) (*Response, error)
	CopyForTab(ctx context.Context, tabID, newTabID, author string, isPrivate bool) error
	UpdateGridPositions(ctx context.Context, items []EditGridPositionItemRequest) (bool, error)
}

func NewStore(dbClient mongo.DbClient, authorProvider author.Provider, enforcer security.Enforcer) Store {
	return &store{
		client:             dbClient,
		collection:         dbClient.Collection(mongo.WidgetMongoCollection),
		tabCollection:      dbClient.Collection(mongo.ViewTabMongoCollection),
		filterCollection:   dbClient.Collection(mongo.WidgetFiltersMongoCollection),
		userPrefCollection: dbClient.Collection(mongo.UserPreferencesMongoCollection),
		authorProvider:     authorProvider,
		enforcer:           enforcer,
	}
}

type store struct {
	client             mongo.DbClient
	collection         mongo.DbCollection
	tabCollection      mongo.DbCollection
	filterCollection   mongo.DbCollection
	userPrefCollection mongo.DbCollection
	authorProvider     author.Provider

	enforcer security.Enforcer
}

func (s *store) FindTabPrivacySettings(ctx context.Context, ids []string) (map[string]apisecurity.ViewTabPrivacySettings, error) {
	results := make([]struct {
		ID                                 string `bson:"_id"`
		apisecurity.ViewTabPrivacySettings `bson:"inline"`
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
			"view":       "$tab.view",
			"is_private": "$tab.is_private",
			"author":     "$tab.author",
		}},
	})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	tabInfos := make(map[string]apisecurity.ViewTabPrivacySettings)
	for _, result := range results {
		if result.View != "" {
			tabInfos[result.ID] = result.ViewTabPrivacySettings
		}
	}

	return tabInfos, nil
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
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$lookup": bson.M{
			"from": mongo.WidgetFiltersMongoCollection,
			"let":  bson.M{"widget": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{
					"$expr":              bson.M{"$eq": bson.A{"$widget", "$$widget"}},
					"is_user_preference": false,
				}},
			},
			"as": "filters",
		}},
		{"$unwind": bson.M{"path": "$filters", "preserveNullAndEmptyArrays": true}},
	}
	pipeline = append(pipeline, s.authorProvider.PipelineForField("filters.author")...)
	pipeline = append(pipeline,
		bson.M{"$sort": bson.M{"filters.position": 1}},
		bson.M{"$group": bson.M{
			"_id":  nil,
			"data": bson.M{"$first": "$$ROOT"},
			"filters": bson.M{"$push": bson.M{"$cond": bson.M{
				"if":   "$filters._id",
				"then": "$filters",
				"else": "$$REMOVE",
			}}},
		}},
		bson.M{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$data",
			bson.M{"filters": "$filters"},
		}}}},
	)
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
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

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Response, error) {
	tabInfo, err := s.getTabPrivacySettings(ctx, r.Tab)
	if err != nil {
		return nil, err
	}

	if tabInfo.IsPrivate && tabInfo.Author != r.Author {
		return nil, common.NewValidationError("tab", "Tab is private.")
	}

	if !tabInfo.IsPrivate {
		ok, err := s.enforcer.Enforce(r.Author, tabInfo.View, model.PermissionUpdate)
		if err != nil {
			panic(err)
		}

		if !ok {
			return nil, common.NewValidationError("tab", "Can't modify a tab.")
		}
	}

	now := datetime.NewCpsTime()
	widget := transformEditRequestToModel(r.EditRequest)
	widget.ID = utils.NewID()
	widget.Tab = r.Tab
	widget.Created = now
	widget.Updated = now
	widget.IsPrivate = tabInfo.IsPrivate

	filters := make([]interface{}, len(r.Filters))
	for i, filter := range r.Filters {
		doc := transformFilterRequestToModel(filter)
		doc.ID = utils.NewID()
		doc.Widget = widget.ID
		doc.Author = widget.Author
		doc.Position = int64(i)
		doc.Created = now
		doc.Updated = now
		doc.IsPrivate = tabInfo.IsPrivate

		if widget.Parameters.MainFilter == filter.ID {
			widget.Parameters.MainFilter = doc.ID
		}

		filters[i] = doc
	}

	var response *Response
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err := s.collection.InsertOne(ctx, widget)
		if err != nil {
			return err
		}

		if len(filters) > 0 {
			_, err := s.filterCollection.InsertMany(ctx, filters)
			if err != nil {
				return err
			}
		}

		response, err = s.GetOneBy(ctx, widget.ID)
		return err
	})

	return response, err
}

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Response, error) {
	oldWidget, err := s.GetOneBy(ctx, r.ID)
	if err != nil || oldWidget == nil {
		return nil, err
	}

	now := datetime.NewCpsTime()
	widget := transformEditRequestToModel(r.EditRequest)
	widget.ID = oldWidget.ID
	widget.Updated = now
	widget.IsPrivate = oldWidget.IsPrivate

	filters := make(map[string]view.WidgetFilter, len(r.Filters))
	for i, filter := range r.Filters {
		doc := transformFilterRequestToModel(filter)
		doc.Widget = widget.ID
		doc.Author = widget.Author
		doc.IsPrivate = widget.IsPrivate
		doc.Position = int64(i)
		doc.Updated = now

		filters[filter.ID] = doc
	}

	var response *Response
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		cursor, err := s.filterCollection.Find(ctx, bson.M{"widget": widget.ID, "is_user_preference": false})
		if err != nil {
			return err
		}
		defer cursor.Close(ctx)
		filterWriteModels := make([]mongodriver.WriteModel, 0, len(filters))
		updateFilterIds := make([]string, 0, len(filters))
		for cursor.Next(ctx) {
			idModel := struct {
				ID string `bson:"_id"`
			}{}
			err := cursor.Decode(&idModel)
			if err != nil {
				return err
			}
			if doc, ok := filters[idModel.ID]; ok {
				updateFilterIds = append(updateFilterIds, idModel.ID)
				filterWriteModels = append(filterWriteModels, mongodriver.NewUpdateOneModel().
					SetFilter(bson.M{"_id": idModel.ID}).
					SetUpdate(bson.M{"$set": doc}))
				delete(filters, idModel.ID)
			}
		}
		for id, doc := range filters {
			doc.ID = utils.NewID()
			doc.Created = now
			updateFilterIds = append(updateFilterIds, doc.ID)
			// Use id from request only to set main filter.
			if id == widget.Parameters.MainFilter {
				widget.Parameters.MainFilter = doc.ID
			}
			filterWriteModels = append(filterWriteModels, mongodriver.NewInsertOneModel().SetDocument(doc))
		}

		update := bson.M{"$set": widget}

		if oldWidget.Type == view.WidgetTypeJunit &&
			(widget.Type != oldWidget.Type ||
				widget.Parameters.IsAPI != oldWidget.Parameters.IsAPI ||
				widget.Parameters.Directory != oldWidget.Parameters.Directory ||
				widget.Parameters.ReportFileRegexp != oldWidget.Parameters.ReportFileRegexp) {
			update["$unset"] = bson.M{"internal_parameters": ""}
		}

		_, err = s.collection.UpdateOne(ctx, bson.M{"_id": widget.ID}, update)
		if err != nil {
			return err
		}

		if len(filterWriteModels) > 0 {
			_, err := s.filterCollection.BulkWrite(ctx, filterWriteModels)
			if err != nil {
				return err
			}
		}

		_, err = s.filterCollection.DeleteMany(ctx, bson.M{
			"widget":             widget.ID,
			"is_user_preference": false,
			"_id":                bson.M{"$nin": updateFilterIds},
		})
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

func (s *store) Copy(ctx context.Context, widgetID string, r CreateRequest) (*Response, error) {
	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		tabInfo, err := s.getTabPrivacySettings(ctx, r.Tab)
		if err != nil {
			return err
		}

		if tabInfo.IsPrivate && tabInfo.Author != r.Author {
			return common.NewValidationError("tab", "Tab is private.")
		}

		if !tabInfo.IsPrivate {
			ok, err := s.enforcer.Enforce(r.Author, tabInfo.View, model.PermissionUpdate)
			if err != nil {
				panic(err)
			}

			if !ok {
				return common.NewValidationError("tab", "Can't modify a tab.")
			}
		}

		response, err = s.copy(ctx, widgetID, tabInfo.IsPrivate, r)
		return err
	})

	return response, err
}

func (s *store) CopyForTab(ctx context.Context, tabID, newTabID, author string, isPrivate bool) error {
	cursor, err := s.collection.Find(ctx, bson.M{"tab": tabID}, options.Find().SetProjection(bson.M{"author": 0}))
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

		_, err = s.copy(ctx, w.ID, isPrivate, CreateRequest{
			Tab: newTabID,
			EditRequest: EditRequest{
				Title:          w.Title,
				Type:           w.Type,
				GridParameters: w.GridParameters,
				Parameters:     w.Parameters,
				Author:         author,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) copy(ctx context.Context, widgetID string, isPrivate bool, r CreateRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	newWidget := view.Widget{
		ID:             utils.NewID(),
		Tab:            r.Tab,
		Title:          r.Title,
		Type:           r.Type,
		GridParameters: r.GridParameters,
		Parameters:     r.Parameters,
		Author:         r.Author,
		IsPrivate:      isPrivate,
		Created:        now,
		Updated:        now,
	}

	cursor, err := s.filterCollection.Find(ctx, bson.M{
		"widget":             widgetID,
		"is_user_preference": false,
	}, options.Find().SetProjection(bson.M{"author": 0}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	mainFilter := ""
	filters := make([]interface{}, 0)
	for cursor.Next(ctx) {
		filter := view.WidgetFilter{}
		err := cursor.Decode(&filter)
		if err != nil {
			return nil, err
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
		filter.IsPrivate = isPrivate
		filters = append(filters, filter)
	}

	newWidget.Parameters.MainFilter = mainFilter
	_, err = s.collection.InsertOne(ctx, newWidget)
	if err != nil {
		return nil, err
	}

	if len(filters) > 0 {
		_, err := s.filterCollection.InsertMany(ctx, filters)
		if err != nil {
			return nil, err
		}
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
				return ValidationError{error: errors.New("widgets are related to different view tabs")}
			}
		}

		count, err := s.collection.CountDocuments(ctx, bson.M{"tab": tabId})
		if err != nil {
			return err
		}
		if count != int64(len(items)) {
			return ValidationError{error: errors.New("widgets are missing")}
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

func (s *store) getTabPrivacySettings(ctx context.Context, tabID string) (apisecurity.ViewTabPrivacySettings, error) {
	var tabInfo apisecurity.ViewTabPrivacySettings

	err := s.tabCollection.FindOne(ctx, bson.M{"_id": tabID}).Decode(&tabInfo)
	if err != nil && errors.Is(err, mongodriver.ErrNoDocuments) {
		return tabInfo, common.NewValidationError("tab", "Tab doesn't exist.")
	}

	return tabInfo, err
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

func transformEditRequestToModel(r EditRequest) view.Widget {
	return view.Widget{
		Title:          r.Title,
		Type:           r.Type,
		GridParameters: r.GridParameters,
		Parameters:     r.Parameters,
		Author:         r.Author,
	}
}

func transformFilterRequestToModel(r FilterRequest) view.WidgetFilter {
	return view.WidgetFilter{
		Title:                       r.Title,
		IsUserPreference:            false,
		AlarmPatternFields:          r.AlarmPatternFieldsRequest.ToModel(),
		EntityPatternFields:         r.EntityPatternFieldsRequest.ToModel(),
		PbehaviorPatternFields:      r.PbehaviorPatternFieldsRequest.ToModel(),
		WeatherServicePatternFields: r.WeatherServicePatternFieldsRequest.ToModel(),
	}
}
