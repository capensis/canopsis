package eventfilter

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apiexternaldata "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/externaldatatable"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/priority"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/request"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/usernotification"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/http"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const tplValErrNoEntity = "Undefined key or field \".Event.Entity"

type Store interface {
	Insert(ctx context.Context, request CreateRequest) (*Response, error)
	GetByID(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, request UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
	FindFailures(ctx context.Context, id string, r FailureRequest) (*AggregationFailureResult, error)
	ReadFailures(ctx context.Context, id string) (bool, error)
	ValidateTemplates(ctx context.Context, request TemplateRequest) (map[string]template.ValidateResponse, error)
	GetTemplateVars() TemplateVarsResponse
	GetCopyVars() CopyVarsResponse
}

type store struct {
	dbClient                mongo.DbClient
	dbCollection            mongo.DbCollection
	dbFailureCollection     mongo.DbCollection
	dbExdataTableCollection mongo.DbCollection
	dbEntityCollection      mongo.DbCollection
	dbTplDataCollection     mongo.DbCollection
	dbTplTestCollection     mongo.DbCollection
	transformer             common.PatternFieldsTransformer
	authorProvider          author.Provider
	notificationStore       usernotification.Store
	tplValidator            validator.Validator
	tplExecutor             libtemplate.Executor
	tplConfigProvider       config.TemplateConfigProvider
	externalDataContainer   *externaldata.GetterContainer
	encoder                 encoding.Encoder
	decoder                 encoding.Decoder
	defaultSearchByFields   []string
	defaultSortBy           string
	exdataTplVars           []template.VarResponse
	configTplVars           []template.VarResponse
	configCopyVars          []template.VarResponse
	dupErrorParser          validation.DuplicateErrorParser
}

func NewStore(
	dbClient mongo.DbClient,
	authorProvider author.Provider,
	transformer common.PatternFieldsTransformer,
	notificationStore usernotification.Store,
	tplValidator validator.Validator,
	tplExecutor libtemplate.Executor,
	tplConfigProvider config.TemplateConfigProvider,
	externalDataContainer *externaldata.GetterContainer,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
) Store {
	exdataTplVars := template.GetEventVars("{{ ", " }}", ".Event", false)
	exdataTplVars = append(exdataTplVars,
		template.VarResponse{Name: "regexpMatch", Value: "{{ .RegexMatch.%field%.%name% }}"},
	)
	configTplVars := make([]template.VarResponse, 0, len(exdataTplVars)+1)
	configTplVars = append(configTplVars, exdataTplVars...)
	configTplVars = append(configTplVars, template.VarResponse{Name: "externalData", Value: "{{ .ExternalData.%reference% }}"})
	configCopyVars := template.GetEventCopyVars("Event")
	configCopyVars = append(configCopyVars,
		template.VarResponse{Name: "regexpMatch", Value: "RegexMatch.%field%.%name%"},
		template.VarResponse{Name: "externalData", Value: "ExternalData.%reference%"},
	)

	return &store{
		dbClient:                dbClient,
		dbCollection:            dbClient.Collection(mongo.EventFilterRuleCollection),
		dbFailureCollection:     dbClient.Collection(mongo.EventFilterFailureCollection),
		dbExdataTableCollection: dbClient.Collection(mongo.ExternalDataTableCollection),
		dbEntityCollection:      dbClient.Collection(mongo.EntityMongoCollection),
		dbTplDataCollection:     dbClient.Collection(mongo.TemplateTestDataCollection),
		dbTplTestCollection:     dbClient.Collection(mongo.TemplateTestCollection),
		transformer:             transformer,
		authorProvider:          authorProvider,
		notificationStore:       notificationStore,
		tplValidator:            tplValidator,
		tplConfigProvider:       tplConfigProvider,
		tplExecutor:             tplExecutor,
		externalDataContainer:   externalDataContainer,
		encoder:                 encoder,
		decoder:                 decoder,
		defaultSearchByFields:   []string{"_id", "author.name", "description", "type"},
		defaultSortBy:           "created",
		exdataTplVars:           exdataTplVars,
		configTplVars:           configTplVars,
		configCopyVars:          configCopyVars,
		dupErrorParser: validation.NewDuplicateErrorParser(map[string]string{
			"_id": "ID already exists.",
		}),
	}
}

func (s *store) Insert(ctx context.Context, request CreateRequest) (*Response, error) {
	model, err := s.transformRequestToDocument(ctx, request.EditRequest)
	if err != nil {
		return nil, err
	}

	model.ID = request.ID
	if model.ID == "" {
		model.ID = utils.NewID()
	}

	now := datetime.NewCpsTime()
	model.Created = &now
	model.Updated = &now

	var response *Response

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		err = s.transformEntityPatternRequestToModel(ctx, request.EntityPatternFieldsRequest, &model)
		if err != nil {
			return err
		}

		_, err = s.dbCollection.InsertOne(ctx, model)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err)
			}

			return err
		}

		err = priority.UpdateFollowing(ctx, s.dbCollection, model.ID, model.Priority)
		if err != nil {
			return err
		}

		response, err = s.GetByID(ctx, model.ID)
		return err
	})

	return response, err
}

func (s *store) GetByID(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
	}
	pipeline = append(pipeline, s.getResponseLookups()...)
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		var rule Response
		err = cursor.Decode(&rule)
		if err != nil {
			return nil, err
		}

		return &rule, nil
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	pipeline := s.authorProvider.Pipeline()
	andCond := make([]bson.M, 0)
	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		andCond = append(andCond, filter)
	}

	if query.OnlyUnreadFailure {
		andCond = append(andCond, bson.M{"unread_failures_count": bson.M{"$gt": 0}})
	}

	if len(andCond) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": andCond}})
	}

	sort := common.GetSortQuery(cmp.Or(query.SortBy, s.defaultSortBy), query.Sort)
	project := s.getResponseLookups()
	project = append(project, sort)
	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		sort,
		project,
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var result AggregationResult
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) Update(ctx context.Context, request UpdateRequest) (*Response, error) {
	updated := datetime.NewCpsTime()
	model, err := s.transformRequestToDocument(ctx, request.EditRequest)
	if err != nil {
		return nil, err
	}

	model.ID = request.ID
	model.Created = nil
	model.Updated = &updated

	update := make(bson.M)
	unset := bson.M{
		"events_count":          "",
		"unread_failures_count": "",
	}

	if model.Start == nil || model.Start.IsZero() || model.Stop == nil || model.Stop.IsZero() {
		unset["start"] = ""
		unset["stop"] = ""
		unset["resolved_start"] = ""
		unset["resolved_stop"] = ""
		unset["next_resolved_start"] = ""
		unset["next_resolved_stop"] = ""
	}

	if len(unset) != 0 {
		update["$unset"] = unset
	}

	var response *Response
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		err = s.transformEntityPatternRequestToModel(ctx, request.EntityPatternFieldsRequest, &model)
		if err != nil {
			return err
		}

		update["$set"] = model

		_, err = s.dbCollection.UpdateOne(
			ctx,
			bson.M{"_id": model.ID},
			update,
		)
		if err != nil {
			return err
		}

		err = priority.UpdateFollowing(ctx, s.dbCollection, model.ID, model.Priority)
		if err != nil {
			return err
		}

		_, err = s.dbTplTestCollection.UpdateMany(ctx, bson.M{"rule._id": model.ID, "type": template.TypeTestEventFilterRule}, bson.M{
			"$set": bson.M{"rule.name": model.Description},
		})
		if err != nil {
			return err
		}

		_, err = s.dbFailureCollection.UpdateMany(ctx, bson.M{"rule": request.ID, "unread": true}, bson.M{
			"$unset": bson.M{
				"unread": "",
			},
		})
		if err != nil {
			return err
		}

		response, err = s.GetByID(ctx, model.ID)

		return err
	})
	if err != nil || response == nil {
		return nil, err
	}

	err = s.notificationStore.DeleteForEventFilterFailure(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	_, err := s.dbFailureCollection.DeleteMany(ctx, bson.M{"rule": id})
	if err != nil {
		return false, err
	}

	var deleted int64

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		_, err = logger.DeleteByFilter(ctx, bson.M{"rule._id": id, "type": template.TypeTestEventFilterRule}, userID,
			s.dbTplTestCollection)
		if err != nil {
			return err
		}

		deleted, err = logger.DeleteOne(ctx, id, userID, s.dbCollection)

		return err
	})
	if err != nil || deleted == 0 {
		return false, err
	}

	err = s.notificationStore.DeleteForEventFilterFailure(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) FindFailures(ctx context.Context, id string, r FailureRequest) (*AggregationFailureResult, error) {
	err := s.dbCollection.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	match := bson.M{"rule": id}
	if r.Type != nil {
		match["type"] = r.Type
	}

	cursor, err := s.dbFailureCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		[]bson.M{{"$match": match}},
		common.GetSortQuery("t", common.SortDesc),
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var result AggregationFailureResult
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) ReadFailures(ctx context.Context, id string) (bool, error) {
	ruleExists := false
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		ruleExists = false
		err := s.dbCollection.
			FindOneAndUpdate(
				ctx,
				bson.M{"_id": id},
				bson.M{"$unset": bson.M{"unread_failures_count": ""}},
				options.FindOneAndUpdate().SetProjection(bson.M{"_id": 1}),
			).
			Err()
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		ruleExists = true
		_, err = s.dbFailureCollection.UpdateMany(ctx, bson.M{"rule": id, "unread": true}, bson.M{"$unset": bson.M{
			"unread": "",
		}})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil || !ruleExists {
		return false, err
	}

	err = s.notificationStore.DeleteForEventFilterFailure(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) ValidateTemplates(ctx context.Context, r TemplateRequest) (map[string]template.ValidateResponse, error) {
	switch r.Rule.Type {
	case eventfilter.RuleTypeEnrichment, eventfilter.RuleTypeChangeEntity:
	default:
		return nil, nil
	}

	var err error
	r.Rule.EntityPatternFieldsRequest, err = s.transformer.TransformEntityPatternFieldsRequest(ctx, r.Rule.EntityPatternFieldsRequest)
	if err != nil {
		return nil, err
	}

	tplData, exdataTestData, err := s.getTplData(ctx, r)
	if err != nil {
		return nil, err
	}

	response := make(map[string]template.ValidateResponse)
	response, tplData.ExternalData, err = s.validateExdataTpls(ctx, r, response, tplData, exdataTestData)
	if err != nil {
		return nil, err
	}

	response, err = s.validateConfigTpls(r, response, tplData)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *store) GetTemplateVars() TemplateVarsResponse {
	return TemplateVarsResponse{
		ExternalData: template.AddEnvVars(s.exdataTplVars, s.tplConfigProvider),
		Config:       template.AddEnvVars(s.configTplVars, s.tplConfigProvider),
	}
}

func (s *store) GetCopyVars() CopyVarsResponse {
	return CopyVarsResponse{
		Config: s.configCopyVars,
	}
}

func (s *store) transformRequestToDocument(ctx context.Context, r EditRequest) (eventfilter.Rule, error) {
	exdates := make([]types.Exdate, len(r.Exdates))
	for i := range r.Exdates {
		exdates[i].Begin = r.Exdates[i].Begin
		exdates[i].End = r.Exdates[i].End
	}

	externalData, err := apiexternaldata.TransformRefParameters(ctx, r.ExternalData, s.dbExdataTableCollection)
	if err != nil {
		return eventfilter.Rule{}, err
	}

	return eventfilter.Rule{
		Author:        r.Author,
		Description:   r.Description,
		Type:          r.Type,
		Priority:      r.Priority,
		Enabled:       r.Enabled,
		Config:        r.Config,
		ExternalData:  externalData,
		EventPattern:  r.EventPattern,
		RRule:         r.RRule,
		Start:         r.Start,
		Stop:          r.Stop,
		ResolvedStart: r.Start,
		ResolvedStop:  r.Stop,
		Exdates:       exdates,
		Exceptions:    r.Exceptions,
	}, nil
}

func (s *store) getResponseLookups() []bson.M {
	pipeline := apiexternaldata.GetRefParametersLookups()
	pipeline = append(pipeline, bson.M{
		"$lookup": bson.M{
			"from":         mongo.PbehaviorExceptionMongoCollection,
			"localField":   "exceptions",
			"foreignField": "_id",
			"as":           "exceptions",
		},
	})

	return pipeline
}

func (s *store) transformEntityPatternRequestToModel(ctx context.Context, r common.EntityPatternFieldsRequest, model *eventfilter.Rule) error {
	transformedEntityPatternRequest, err := s.transformer.TransformEntityPatternFieldsRequest(ctx, r)
	if err != nil {
		return err
	}

	model.Aliases = transformedEntityPatternRequest.Aliases
	model.EntityPatternFields = transformedEntityPatternRequest.ToModel()

	return nil
}

func (s *store) getTplData(ctx context.Context, r TemplateRequest) (eventfilter.Template, map[int]template.ResponseTestData, error) {
	var tplData eventfilter.Template
	eventDataID := r.TestData.Event
	exdataTestDataIDs := r.TestData.Responses
	if r.TestData.Test != "" {
		test := template.TestModel{}
		err := s.dbTplTestCollection.
			FindOne(ctx, bson.M{"_id": r.TestData.Test, "type": template.TypeTestEventFilterRule, "rule._id": r.Rule.ID}).
			Decode(&test)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return tplData, nil, common.NewValidationError("testdata.test", "Test doesn't exist.")
			}

			return tplData, nil, err
		}

		if eventDataID == "" {
			eventDataID = test.Data.Event
		}

		if exdataTestDataIDs == nil {
			exdataTestDataIDs = test.Data.Responses
		}
	}

	if eventDataID == "" {
		return tplData, nil, common.NewValidationError("testdata.event", "Event is missing.")
	}

	event, err := template.GetEventData(ctx, s.dbTplDataCollection, eventDataID, s.encoder, s.decoder)
	if err != nil {
		return tplData, nil, err
	}

	if event == nil {
		return tplData, nil, common.NewValidationError("testdata.event", "Event doesn't exist.")
	}

	var exdataTestData map[int]template.ResponseTestData
	if len(exdataTestDataIDs) > 0 {
		if len(exdataTestDataIDs) > len(r.Rule.ExternalData) {
			return tplData, nil, common.NewValidationError("testdata.responses."+strconv.Itoa(len(r.Rule.ExternalData)), "Response is redundant.")
		}

		exdataTestData, err = template.GetResponseData(ctx, s.dbTplDataCollection, exdataTestDataIDs)
		if err != nil {
			return tplData, nil, err
		}

		if len(exdataTestData) == 0 {
			return tplData, nil, common.NewValidationError("testdata.responses", "Response doesn't exist.")
		}
	}

	var entity types.Entity
	entityID := event.GetEID()
	if entityID != "" {
		err = s.dbEntityCollection.FindOne(ctx, bson.M{"_id": entityID, "soft_deleted": nil}).Decode(&entity)
		// ignore missing Entity in case only Event is used in templates
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return tplData, nil, err
		}
	}

	var matched bool
	var eventRegexMatches match.EventRegexMatches
	var entityRegexMatches match.EntityRegexMatches
	if len(r.Rule.EventPattern) > 0 {
		matched, eventRegexMatches, err = match.MatchEventPatternWithRegexMatches(r.Rule.EventPattern, event)
		if err != nil {
			return tplData, nil, common.NewValidationError("rule.event_pattern", "EventPattern is invalid event pattern.")
		}

		if !matched {
			return tplData, nil, common.NewValidationError("testdata.event", "Event is not matched to event pattern.")
		}
	}

	var entityPattern pattern.Entity
	if len(r.Rule.CorporatePattern.EntityPattern) > 0 {
		entityPattern = r.Rule.CorporatePattern.EntityPattern
	} else {
		entityPattern = r.Rule.EntityPattern
	}

	if len(entityPattern) > 0 {
		if entity.ID == "" {
			return tplData, nil, common.NewValidationError("testdata.event", "Corresponding entity doesn't exist.")
		}

		matched, entityRegexMatches, err = match.MatchEntityPatternWithRegexMatches(entityPattern, &entity)
		if err != nil {
			return tplData, nil, common.NewValidationError("rule.entity_pattern", "EntityPattern is invalid entity pattern.")
		}

		if !matched {
			return tplData, nil, common.NewValidationError("testdata.event", "Corresponding entity is not matched to entity pattern.")
		}
	}

	tplData = eventfilter.Template{
		Event: event,
		RegexMatch: eventfilter.RegexMatch{
			EventRegexMatches: eventRegexMatches,
			Entity:            entityRegexMatches,
		},
	}
	if entity.ID != "" {
		tplData.Event.Entity = &entity
	}

	return tplData, exdataTestData, nil
}

func (s *store) validateExdataTpls(
	ctx context.Context,
	r TemplateRequest,
	response map[string]template.ValidateResponse,
	tplData eventfilter.Template,
	exdataTestData map[int]template.ResponseTestData,
) (map[string]template.ValidateResponse, map[string]any, error) {
	externalData := make(map[string]any, len(r.Rule.ExternalData))
	var err error
	for i, d := range r.Rule.ExternalData {
		prefix := "external_data." + strconv.Itoa(i)
		switch d.Type {
		case externaldata.RefTypeTable:
			isValid := true
			for k, v := range d.Regexp {
				vr, err := s.validateTpl(v, tplData)
				if err != nil {
					return nil, nil, err
				}

				response[prefix+".regexp."+k] = vr
				if !vr.IsValid {
					isValid = false
				}
			}

			for k, v := range d.Select {
				vr, err := s.validateTpl(v, tplData)
				if err != nil {
					return nil, nil, err
				}

				response[prefix+".select."+k] = vr
				if !vr.IsValid {
					isValid = false
				}
			}

			if isValid {
				externalData[d.Reference], err = s.processTableExdata(ctx, d, tplData, "rule."+prefix)
				if err != nil {
					return nil, nil, err
				}
			}

			if _, ok := exdataTestData[i]; ok {
				return nil, nil, common.NewValidationError("testdata.responses."+strconv.Itoa(i), "Response is redundant.")
			}
		case externaldata.RefTypeAPI:
			if d.Request != nil {
				response[prefix+".request.url"], err = s.validateTpl(d.Request.URL, tplData)
				if err != nil {
					return nil, nil, err
				}

				response[prefix+".request.payload"], err = s.validateTpl(d.Request.Payload, tplData)
				if err != nil {
					return nil, nil, err
				}
			}

			if td, ok := exdataTestData[i]; ok {
				b, err := s.encoder.Encode(td.Body)
				if err != nil {
					return nil, nil, common.NewValidationError("testdata.responses."+strconv.Itoa(i), "Response is not JSON.")
				}

				flatten, basicRes, err := http.FlattenJSON(b)
				if err != nil {
					return nil, nil, common.NewValidationError("testdata.responses."+strconv.Itoa(i), "Response is not JSON.")
				}

				if flatten == nil {
					externalData[d.Reference] = basicRes
				} else {
					externalData[d.Reference] = flatten
				}
			} else {
				return nil, nil, common.NewValidationError("testdata.responses."+strconv.Itoa(i), "Response is missing.")
			}
		}
	}

	return response, externalData, nil
}

func (s *store) validateConfigTpls(
	r TemplateRequest,
	response map[string]template.ValidateResponse,
	tplData eventfilter.Template,
) (map[string]template.ValidateResponse, error) {
	var err error
	switch r.Rule.Type {
	case eventfilter.RuleTypeEnrichment:
		for i, action := range r.Rule.Config.Actions {
			switch action.Type {
			case eventfilter.ActionSetFieldFromTemplate,
				eventfilter.ActionSetEntityInfoFromTemplate,
				eventfilter.ActionSetTagsFromTemplate:
				val, ok := action.Value.(string)
				if !ok {
					return nil, common.NewValidationError("rule.config.actions."+strconv.Itoa(i)+".value", "Value must be string.")
				}

				response["config.actions."+strconv.Itoa(i)+".value"], err = s.validateTpl(val, tplData)
				if err != nil {
					return nil, err
				}
			}
		}
	case eventfilter.RuleTypeChangeEntity:
		if r.Rule.Config.Resource != "" {
			response["config.resource"], err = s.validateTpl(r.Rule.Config.Resource, tplData)
			if err != nil {
				return nil, err
			}
		}

		if r.Rule.Config.Component != "" {
			response["config.component"], err = s.validateTpl(r.Rule.Config.Component, tplData)
			if err != nil {
				return nil, err
			}
		}

		if r.Rule.Config.Connector != "" {
			response["config.connector"], err = s.validateTpl(r.Rule.Config.Connector, tplData)
			if err != nil {
				return nil, err
			}
		}

		if r.Rule.Config.ConnectorName != "" {
			response["config.connector_name"], err = s.validateTpl(r.Rule.Config.ConnectorName, tplData)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, nil
}

func (s *store) processTableExdata(
	ctx context.Context,
	d template.TemplateRefParameters,
	tplData eventfilter.Template,
	field string,
) (any, error) {
	getter, ok := s.externalDataContainer.Get(d.Type)
	if !ok {
		return nil, fmt.Errorf("cannot find external data getter by type %q", d.Type)
	}

	refParam := externaldata.RefParameters{
		Reference: d.Reference,
		Type:      d.Type,
		Table:     d.Table,
		Select:    d.Select,
		Regexp:    d.Regexp,
		SortBy:    d.SortBy,
		Sort:      d.Sort,
		Optional:  d.Optional,
	}
	if d.Request != nil {
		refParam.Request = &request.Parameters{
			URL:     d.Request.URL,
			Payload: d.Request.Payload,
			Headers: d.Request.Headers,
		}
	}

	params, err := apiexternaldata.TransformRefParameters(ctx, []externaldata.RefParameters{refParam}, s.dbExdataTableCollection)
	if err != nil {
		return nil, err
	}

	parsedParams := externaldata.ParseRefParameters(params, s.tplExecutor)
	if len(parsedParams) == 0 {
		return nil, errors.New("expected not empty array")
	}

	res, err := getter.Get(ctx, parsedParams[0], tplData)
	if err != nil {
		getterTplErr := &externaldata.GetterTplError{}
		getterErr := &externaldata.GetterError{}
		if errors.As(err, &getterTplErr) || errors.As(err, &getterErr) {
			return nil, common.NewValidationError(field, err.Error())
		}

		return nil, err
	}

	return res, nil
}

func (s *store) validateTpl(str string, data eventfilter.Template) (template.ValidateResponse, error) {
	response, err := template.Validate(s.tplValidator, str, data)
	if err != nil {
		return response, err
	}

	// if only Event is provided but Entity is used in templates
	if response.Err != nil && strings.Contains(response.Err.Message, tplValErrNoEntity) && data.Event.Entity == nil {
		return response, common.NewValidationError("testdata.event", "Corresponding entity doesn't exist.")
	}

	return response, nil
}
