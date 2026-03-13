package widget

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	tplvalidator "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	entityTplVarsIndex = 1
)

type Store interface {
	FindTabPrivacySettings(ctx context.Context, ids []string) (map[string]apisecurity.ViewTabPrivacySettings, error)
	FindViewIdByTab(ctx context.Context, tabId string) (string, error)
	GetOneBy(ctx context.Context, id string) (*Response, error)
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
	Copy(ctx context.Context, widgetID string, r CreateRequest) (*Response, error)
	CopyForTab(ctx context.Context, tabID, newTabID, author string, isPrivate bool) error
	UpdateGridPositions(ctx context.Context, items []EditGridPositionItemRequest) (bool, error)
	ValidateTemplates(ctx context.Context, request TemplateRequest) (map[string]template.ValidateResponse, error)
	GetTemplateVars(ctx context.Context) (TemplateVarsResponse, error)
}

func NewStore(
	dbClient mongo.DbClient,
	authorProvider author.Provider,
	enforcer security.Enforcer,
	transformer patternfields.Transformer,
	tplValidator tplvalidator.Validator,
	tplConfigProvider config.TemplateConfigProvider,
) Store {
	return &store{
		client:                    dbClient,
		collection:                dbClient.Collection(mongo.WidgetMongoCollection),
		tabCollection:             dbClient.Collection(mongo.ViewTabMongoCollection),
		filterCollection:          dbClient.Collection(mongo.WidgetFiltersMongoCollection),
		userPrefCollection:        dbClient.Collection(mongo.UserPreferencesMongoCollection),
		widgetTemplateCollection:  dbClient.Collection(mongo.WidgetTemplateMongoCollection),
		alarmCollection:           dbClient.Collection(mongo.AlarmMongoCollection),
		tplTestCollection:         dbClient.Collection(mongo.TemplateTestCollection),
		entityInfosPropCollection: dbClient.Collection(mongo.EntityInfosPropertyCollection),
		commentTemplateCollection: dbClient.Collection(mongo.CommentTemplateMongoCollection),
		authorProvider:            authorProvider,
		transformer:               transformer,
		enforcer:                  enforcer,
		tplValidator:              tplValidator,
		tplConfigProvider:         tplConfigProvider,
		tplVars: []template.VarResponse{
			{
				Name:  "alarm",
				Value: template.GetAlarmVars("{{ ", " }}", ".Alarm", false),
			},
			{
				Name:  "entity",
				Value: template.GetEntityVars("{{ ", " }}", ".Entity", false),
			},
		},
	}
}

type store struct {
	client                    mongo.DbClient
	collection                mongo.DbCollection
	tabCollection             mongo.DbCollection
	filterCollection          mongo.DbCollection
	userPrefCollection        mongo.DbCollection
	widgetTemplateCollection  mongo.DbCollection
	alarmCollection           mongo.DbCollection
	tplTestCollection         mongo.DbCollection
	entityInfosPropCollection mongo.DbCollection
	commentTemplateCollection mongo.DbCollection
	authorProvider            author.Provider
	transformer               patternfields.Transformer
	enforcer                  security.Enforcer
	tplValidator              tplvalidator.Validator
	tplConfigProvider         config.TemplateConfigProvider
	tplVars                   []template.VarResponse
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
		{"$lookup": bson.M{
			"from":         mongo.CommentTemplateMongoCollection,
			"localField":   "parameters.comment_templates",
			"foreignField": "_id",
			"as":           "comment_templates",
		}},
		{"$unwind": bson.M{"path": "$filters", "preserveNullAndEmptyArrays": true}},
		{"$unwind": bson.M{"path": "$comment_templates", "preserveNullAndEmptyArrays": true}},
	}
	pipeline = append(pipeline, s.authorProvider.PipelineForField("filters.author")...)
	pipeline = append(pipeline, s.authorProvider.PipelineForField("comment_templates.author")...)
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
			"comment_templates": bson.M{"$push": bson.M{"$cond": bson.M{
				"if":   "$comment_templates._id",
				"then": "$comment_templates",
				"else": "$$REMOVE",
			}}},
		}},
		bson.M{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$data",
			bson.M{"filters": "$filters"},
			bson.M{"comment_templates": "$comment_templates"},
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

	if tabInfo.ID == "" {
		return nil, validation.NewSingleError("not_exist", "Tab", "Tab", r)
	}

	if tabInfo.IsPrivate && tabInfo.Author != r.Author {
		return nil, validation.NewSingleError("not_exist", "Tab", "Tab", r)
	}

	if !tabInfo.IsPrivate {
		ok, err := s.enforcer.Enforce(r.Author, tabInfo.View, model.PermissionUpdate)
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, httperror.NewForbiddenError("")
		}
	}

	// todo: put inside transaction too
	err = s.transformTemplateFields(ctx, &r.EditRequest)
	if err != nil {
		return nil, err
	}

	now := datetime.NewCpsTime()
	widget := transformEditRequestToModel(r.EditRequest)
	widget.ID = utils.NewID()
	widget.Tab = r.Tab
	widget.Created = now
	widget.Updated = now
	widget.IsPrivate = tabInfo.IsPrivate

	var response *Response
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		patterns, aliases, err := s.fetchPatterns(ctx, r.Filters)
		if err != nil {
			return err
		}

		var valErrs validator.ValidationErrors

		for i := range r.Parameters.CommentTemplates {
			err := s.commentTemplateCollection.FindOne(ctx, bson.M{"_id": r.Parameters.CommentTemplates[i]}).Err()
			if err != nil {
				if errors.Is(err, mongodriver.ErrNoDocuments) {
					valErrs = append(valErrs, validation.NewFieldError("not_exist", strconv.Itoa(i), "Parameters.CommentTemplates."+strconv.Itoa(i)))
					continue
				}

				return err
			}
		}

		filters := make([]view.WidgetFilter, len(r.Filters))
		for i, filterRequest := range r.Filters {
			doc := view.WidgetFilter{
				ID:               utils.NewID(),
				Title:            filterRequest.Title,
				IsUserPreference: false,
				Widget:           widget.ID,
				Author:           widget.Author,
				Position:         int64(i),
				Created:          now,
				Updated:          now,
				IsPrivate:        tabInfo.IsPrivate,
			}

			if widget.Parameters.MainFilter == filterRequest.ID {
				widget.Parameters.MainFilter = doc.ID
			}

			fErrs := s.transformPatternRequestsToModel(filterRequest, i, &doc, patterns, aliases)
			if fErrs != nil {
				valErrs = append(valErrs, fErrs...)
				continue
			}

			filters[i] = doc
		}

		if len(valErrs) > 0 {
			return validation.NewError(valErrs, r)
		}

		_, err = s.collection.InsertOne(ctx, widget)
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

	// todo: put inside transaction too
	err = s.transformTemplateFields(ctx, &r.EditRequest)
	if err != nil {
		return nil, err
	}

	now := datetime.NewCpsTime()
	widget := transformEditRequestToModel(r.EditRequest)
	widget.ID = oldWidget.ID
	widget.Updated = now
	widget.IsPrivate = oldWidget.IsPrivate

	var response *Response
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		patterns, aliases, err := s.fetchPatterns(ctx, r.Filters)
		if err != nil {
			return err
		}

		var valErrs validator.ValidationErrors

		for i := range r.Parameters.CommentTemplates {
			err := s.commentTemplateCollection.FindOne(ctx, bson.M{"_id": r.Parameters.CommentTemplates[i]}).Err()
			if err != nil {
				if errors.Is(err, mongodriver.ErrNoDocuments) {
					valErrs = append(valErrs, validation.NewFieldError("not_exist", strconv.Itoa(i), "Parameters.CommentTemplates."+strconv.Itoa(i)))
					continue
				}

				return err
			}
		}

		filters := make(map[string]view.WidgetFilter, len(r.Filters))
		for i, filterRequest := range r.Filters {
			doc := view.WidgetFilter{
				Title:            filterRequest.Title,
				IsUserPreference: false,
				Widget:           widget.ID,
				Author:           widget.Author,
				IsPrivate:        widget.IsPrivate,
				Position:         int64(i),
				Updated:          now,
			}

			fErrs := s.transformPatternRequestsToModel(filterRequest, i, &doc, patterns, aliases)
			if fErrs != nil {
				valErrs = append(valErrs, fErrs...)
				continue
			}

			filters[filterRequest.ID] = doc
		}

		if len(valErrs) > 0 {
			return validation.NewError(valErrs, r)
		}

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

		// required to get the author in action log listener.
		_, err = s.filterCollection.UpdateMany(
			ctx,
			bson.M{
				"widget":             widget.ID,
				"is_user_preference": false,
				"_id":                bson.M{"$nin": updateFilterIds},
			},
			bson.M{
				"$set": bson.M{"author": r.Author},
			},
		)
		if err != nil {
			return err
		}

		_, err = s.filterCollection.DeleteMany(ctx, bson.M{
			"widget":             widget.ID,
			"is_user_preference": false,
			"_id":                bson.M{"$nin": updateFilterIds},
		})
		if err != nil {
			return err
		}

		_, err = s.tplTestCollection.UpdateMany(ctx, bson.M{"rule._id": widget.ID, "type": template.TypeTestWidget}, bson.M{
			"$set": bson.M{"rule.name": widget.Title},
		})
		if err != nil {
			return err
		}

		response, err = s.GetOneBy(ctx, widget.ID)

		return err
	})

	return response, err
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	res := false
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = false

		// required to get the author in action log listener.
		result, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || result.MatchedCount == 0 {
			return err
		}

		delCount, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil || delCount == 0 {
			return err
		}

		err = s.deleteUserPreferences(ctx, id)
		if err != nil {
			return err
		}

		err = s.deleteFilters(ctx, id, userID)
		if err != nil {
			return err
		}

		_, err = logger.DeleteByFilter(ctx, bson.M{"rule._id": id, "type": template.TypeTestWidget}, userID,
			s.tplTestCollection)
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

		if tabInfo.ID == "" {
			return validation.NewSingleError("not_exist", "Tab", "Tab", r)
		}

		if tabInfo.IsPrivate && tabInfo.Author != r.Author {
			return validation.NewSingleError("not_exist", "Tab", "Tab", r)
		}

		if !tabInfo.IsPrivate {
			ok, err := s.enforcer.Enforce(r.Author, tabInfo.View, model.PermissionUpdate)
			if err != nil {
				return err
			}

			if !ok {
				return httperror.NewForbiddenError("")
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
				return validation.NewSingleError("not_applicable", "items", "items", nil)
			}
		}

		count, err := s.collection.CountDocuments(ctx, bson.M{"tab": tabId})
		if err != nil {
			return err
		}
		if count != int64(len(items)) {
			return validation.NewSingleErrorWithParam("slicelen", "items", "items", strconv.FormatInt(count, 10), nil)
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

func (s *store) ValidateTemplates(ctx context.Context, r TemplateRequest) (map[string]template.ValidateResponse, error) {
	var alarm types.AlarmWithEntity
	var err error
	if r.TestData.Test != "" {
		alarm, err = template.GetAlarmDataFromTest(ctx, s.tplTestCollection, r.TestData.Test, template.TypeTestWidget, r.Rule.ID)
		if err != nil {
			return nil, err
		}
	}

	if r.TestData.Alarm == "" {
		if alarm.Alarm.ID == "" {
			return nil, validation.NewSingleError("required", "Alarm", "TestData.Alarm", r)
		}
	} else if r.TestData.Alarm != alarm.Alarm.ID { // keep snapshot from the test
		alarm, err = s.findAlarm(ctx, r.TestData.Alarm)
		if err != nil {
			return nil, err
		}

		if alarm.Alarm.ID == "" {
			return nil, validation.NewSingleError("not_exist", "Alarm", "TestData.Alarm", r)
		}
	}

	response := make(map[string]template.ValidateResponse)
	for i, column := range r.Rule.Columns {
		if column.Template != "" {
			response["columns."+strconv.Itoa(i)+".template"], err = template.Validate(s.tplValidator, column.Template, alarm)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, nil
}

func (s *store) GetTemplateVars(ctx context.Context) (TemplateVarsResponse, error) {
	columnTplVars := template.AddEnvVars(s.tplVars, s.tplConfigProvider)

	aliases, err := template.GetAliases(ctx, s.entityInfosPropCollection)
	if err != nil {
		return TemplateVarsResponse{}, fmt.Errorf("failed to get aliases: %w", err)
	}

	err = template.AddAliasesVars(columnTplVars, aliases, entityTplVarsIndex, "{{ (index .Entity.Infos \"", "\").Value }}")
	if err != nil {
		return TemplateVarsResponse{}, fmt.Errorf("failed to add aliases to columnTplVars: %w", err)
	}

	return TemplateVarsResponse{
		Column: columnTplVars,
	}, nil
}

func (s *store) getTabPrivacySettings(ctx context.Context, tabID string) (apisecurity.ViewTabPrivacySettings, error) {
	var tabInfo apisecurity.ViewTabPrivacySettings

	err := s.tabCollection.FindOne(ctx, bson.M{"_id": tabID}).Decode(&tabInfo)
	if err != nil && errors.Is(err, mongodriver.ErrNoDocuments) {
		return tabInfo, nil
	}

	return tabInfo, err
}

func (s *store) deleteUserPreferences(ctx context.Context, widgetID string) error {
	_, err := s.userPrefCollection.DeleteMany(ctx, bson.M{
		"widget": widgetID,
	})

	return err
}

func (s *store) deleteFilters(ctx context.Context, widgetID, userID string) error {
	// required to get the author in action log listener.
	_, err := s.filterCollection.UpdateMany(ctx, bson.M{"widget": widgetID}, bson.M{"$set": bson.M{"author": userID}})
	if err != nil {
		return err
	}

	_, err = s.filterCollection.DeleteMany(ctx, bson.M{
		"widget": widgetID,
	})

	return err
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
	filters := make([]any, 0)
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

func (s *store) transformPatternRequestsToModel(
	r FilterRequest,
	idx int,
	model *view.WidgetFilter,
	patterns patternfields.Patterns,
	aliases patternfields.Aliases,
) validator.ValidationErrors {
	filtersFieldName := "Filters"
	sIdx := strconv.Itoa(idx)
	var valErrs, applyErrs validator.ValidationErrors
	r.AlarmRequest, applyErrs = s.transformer.ApplyAlarmCorporatePattern(r.AlarmRequest, patterns, filtersFieldName, sIdx, "CorporateAlarmPattern")
	valErrs = append(valErrs, applyErrs...)
	if r.CorporateEntityPattern != "" {
		r.EntityRequest, model.Aliases, applyErrs = s.transformer.ApplyEntityCorporatePattern(r.EntityRequest, patterns, filtersFieldName, sIdx, "CorporateEntityPattern")
	} else if r.EntityPattern != nil {
		r.EntityPattern, model.Aliases, applyErrs = s.transformer.ApplyAliases(r.EntityPattern, aliases, filtersFieldName, sIdx, "EntityPattern")
	}

	valErrs = append(valErrs, applyErrs...)
	r.PbehaviorRequest, applyErrs = s.transformer.ApplyPbehaviorCorporatePattern(r.PbehaviorRequest, patterns, filtersFieldName, sIdx, "CorporatePbehaviorPattern")
	valErrs = append(valErrs, applyErrs...)
	r.WeatherServiceRequest, applyErrs = s.transformer.ApplyServiceWeatherCorporatePattern(r.WeatherServiceRequest, patterns, filtersFieldName, sIdx, "CorporateServiceWeatherPattern")
	valErrs = append(valErrs, applyErrs...)
	if len(valErrs) > 0 {
		return valErrs
	}

	model.AlarmPatternFields = r.AlarmRequest.ToModel()
	model.EntityPatternFields = r.EntityRequest.ToModel()
	model.PbehaviorPatternFields = r.PbehaviorRequest.ToModel()
	model.WeatherServicePatternFields = r.WeatherServiceRequest.ToModel()

	return nil
}

func (s *store) transformTemplateFields(ctx context.Context, r *EditRequest) error {
	widgetParametersByType := view.GetWidgetTemplateParameters()[r.Type]
	for tplType, widgetParameters := range widgetParametersByType {
		for _, parameter := range widgetParameters {
			parameters := r.Parameters.RemainParameters
			key := parameter
			parts := strings.Split(parameter, ".")
			if len(parts) > 1 {
				key = parts[len(parts)-1]
				var ok bool
				for i := 0; i < len(parts)-1; i++ {
					parameters, ok = parameters[parts[i]].(map[string]any)
					if !ok {
						break
					}
				}
				if !ok {
					continue
				}
			}

			tplId, ok := parameters[key+"Template"].(string)
			if !ok || tplId == "" {
				continue
			}
			tpl := view.WidgetTemplate{}
			err := s.widgetTemplateCollection.FindOne(ctx, bson.M{
				"_id":  tplId,
				"type": tplType,
			}).Decode(&tpl)
			if err != nil {
				if errors.Is(err, mongodriver.ErrNoDocuments) {
					return validation.NewSingleError("not_exist", parameter+"Template", "parameters."+parameter+"Template", nil)
				}

				return err
			}

			parameters[key+"TemplateTitle"] = tpl.Title
			switch tpl.Type {
			case view.WidgetTemplateTypeAlarmColumns,
				view.WidgetTemplateTypeEntityColumns:
				parameters[key] = tpl.Columns
			case view.WidgetTemplateTypeAlarmMoreInfos,
				view.WidgetTemplateTypeAlarmExportToPDF,
				view.WidgetTemplateTypeServiceWeatherItem,
				view.WidgetTemplateTypeServiceWeatherModal,
				view.WidgetTemplateTypeServiceWeatherEntity:
				parameters[key] = tpl.Content
			case view.WidgetTemplateTypeAlarmQuickActions,
				view.WidgetTemplateTypeAlarmQuickMassActions:
				parameters[key] = tpl.Actions
			case view.WidgetTemplateTypeAlarmSortColumns:
				parameters[key] = tpl.SortColumns
			}
		}
	}

	return nil
}

func (s *store) findAlarm(ctx context.Context, alarmID string) (types.AlarmWithEntity, error) {
	cursor, err := s.alarmCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": alarmID}},
		{"$replaceRoot": bson.M{"newRoot": bson.M{
			"alarm": "$$ROOT",
		}}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "alarm.d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": bson.M{"path": "$entity", "preserveNullAndEmptyArrays": true}},
	})
	if err != nil {
		return types.AlarmWithEntity{}, fmt.Errorf("cannot find alarm: %w", err)
	}

	defer cursor.Close(ctx)
	if !cursor.Next(ctx) {
		return types.AlarmWithEntity{}, nil
	}

	var alarm types.AlarmWithEntity
	err = cursor.Decode(&alarm)
	if err != nil {
		return types.AlarmWithEntity{}, fmt.Errorf("cannot decode alarm: %w", err)
	}

	if err = cursor.Err(); err != nil {
		return types.AlarmWithEntity{}, fmt.Errorf("cannot find alarm: %w", err)
	}

	return alarm, nil
}

func (s *store) fetchPatterns(ctx context.Context, filters []FilterRequest) (patternfields.Patterns, patternfields.Aliases, error) {
	patternIDs := make([]string, 0, len(filters)*4)
	aliases := make([]string, 0, len(filters))
	for _, fr := range filters {
		patternIDs = append(patternIDs,
			fr.CorporateEntityPattern,
			fr.CorporateAlarmPattern,
			fr.CorporatePbehaviorPattern,
			fr.CorporateWeatherServicePattern,
		)
		aliases = append(aliases, patternfields.GetAliases(fr.EntityPattern)...)
	}

	patterns, err := s.transformer.FetchCorporatePatterns(ctx, patternIDs...)
	if err != nil {
		return nil, nil, err
	}

	aliasProps, err := s.transformer.FetchAliases(ctx, aliases)
	if err != nil {
		return nil, nil, err
	}

	return patterns, aliasProps, nil
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
