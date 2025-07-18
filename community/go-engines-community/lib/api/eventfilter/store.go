package eventfilter

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apiexternaldata "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/externaldatatable"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/priority"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
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
	ValidateTemplates(ctx context.Context, request TemplateRequest) (map[string]template.Response, error)
	GetTemplateVars() TemplateVarsResponse
}

type store struct {
	dbClient                mongo.DbClient
	dbCollection            mongo.DbCollection
	dbFailureCollection     mongo.DbCollection
	dbExdataTableCollection mongo.DbCollection
	dbEntityCollection      mongo.DbCollection
	authorProvider          author.Provider
	tplValidator            validator.Validator
	tplExecutor             libtemplate.Executor
	tplConfigProvider       config.TemplateConfigProvider
	externalDataContainer   *externaldata.GetterContainer
	defaultSearchByFields   []string
	defaultSortBy           string
	tplVars                 []string
	tplVarExternalData      string
}

func NewStore(
	dbClient mongo.DbClient,
	authorProvider author.Provider,
	tplValidator validator.Validator,
	tplExecutor libtemplate.Executor,
	tplConfigProvider config.TemplateConfigProvider,
	externalDataContainer *externaldata.GetterContainer,
) Store {
	return &store{
		dbClient:                dbClient,
		dbCollection:            dbClient.Collection(mongo.EventFilterRuleCollection),
		dbFailureCollection:     dbClient.Collection(mongo.EventFilterFailureCollection),
		dbExdataTableCollection: dbClient.Collection(mongo.ExternalDataTableCollection),
		dbEntityCollection:      dbClient.Collection(mongo.EntityMongoCollection),
		authorProvider:          authorProvider,
		tplValidator:            tplValidator,
		tplConfigProvider:       tplConfigProvider,
		tplExecutor:             tplExecutor,
		externalDataContainer:   externalDataContainer,
		defaultSearchByFields:   []string{"_id", "author.name", "description", "type"},
		defaultSortBy:           "created",
		tplVars: []string{
			".Event.Connector",
			".Event.ConnectorName",
			".Event.Component",
			".Event.Resource",
			".Event.Output",
			"index .Event.ExtraInfos \"%infos_name%\"",
			".RegexMatch.%field%.%name%",
		},
		tplVarExternalData: ".ExternalData.%reference%",
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
		_, err := s.dbCollection.InsertOne(ctx, model)
		if err != nil {
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
	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
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

	update := bson.M{"$set": model}
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
		_, err := s.dbCollection.UpdateOne(
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

		response, err = s.GetByID(ctx, model.ID)
		return err
	})
	if err != nil || response == nil {
		return nil, err
	}

	_, err = s.dbFailureCollection.UpdateMany(ctx, bson.M{"rule": request.ID, "unread": true}, bson.M{
		"$unset": bson.M{
			"unread": "",
		},
	})
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

		// required to get the author in action log listener.
		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		deleted, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		return err
	})

	return deleted > 0, err
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

		return err
	})

	return ruleExists, err
}

func (s *store) ValidateTemplates(ctx context.Context, r TemplateRequest) (map[string]template.Response, error) {
	switch r.Request.Type {
	case eventfilter.RuleTypeEnrichment, eventfilter.RuleTypeChangeEntity:
	default:
		return nil, nil
	}

	data, err := s.getTplData(ctx, r)
	if err != nil {
		return nil, err
	}

	response := make(map[string]template.Response)
	externalData := make(map[string]any, len(r.Request.ExternalData))
	for i, d := range r.Request.ExternalData {
		switch d.Type {
		case externaldata.RefTypeTable:
			for k, v := range d.Regexp {
				response["external_data."+strconv.Itoa(i)+".regexp."+k], err = s.validateTpl(v, data)
				if err != nil {
					return nil, err
				}
			}

			for k, v := range d.Select {
				response["external_data."+strconv.Itoa(i)+".select."+k], err = s.validateTpl(v, data)
				if err != nil {
					return nil, err
				}
			}

			getter, ok := s.externalDataContainer.Get(d.Type)
			if !ok {
				return nil, fmt.Errorf("cannot find external data getter by type %q", d.Type)
			}

			params, err := apiexternaldata.TransformRefParameters(ctx, []externaldata.RefParameters{d}, s.dbExdataTableCollection)
			if err != nil {
				return nil, err
			}

			parsedParams := externaldata.ParseRefParameters(params, s.tplExecutor)
			if len(parsedParams) == 0 {
				return nil, errors.New("expected not empty array")
			}

			externalData[d.Reference], err = getter.Get(ctx, parsedParams[0], data)
			if err != nil {
				getterTplErr := &externaldata.GetterTplError{}
				getterErr := &externaldata.GetterError{}
				if errors.As(err, &getterTplErr) || errors.As(err, &getterErr) {
					return nil, common.NewValidationError("request.external_data."+strconv.Itoa(i), err.Error())
				}

				return nil, err
			}
		case externaldata.RefTypeAPI:
			if d.Request != nil {
				response["external_data."+strconv.Itoa(i)+".request.url"], err = s.validateTpl(d.Request.URL, data)
				if err != nil {
					return nil, err
				}

				response["external_data."+strconv.Itoa(i)+".request.payload"], err = s.validateTpl(d.Request.Payload, data)
				if err != nil {
					return nil, err
				}
			}

			var ok bool
			externalData[d.Reference], ok = r.Data.ExternalData[d.Reference]
			if !ok {
				return nil, common.NewValidationError("data.external_data."+d.Reference, d.Reference+" is missing.")
			}
		}
	}

	data.ExternalData = externalData
	switch r.Request.Type {
	case eventfilter.RuleTypeEnrichment:
		for i, action := range r.Request.Config.Actions {
			switch action.Type {
			case eventfilter.ActionSetFieldFromTemplate,
				eventfilter.ActionSetEntityInfoFromTemplate,
				eventfilter.ActionSetTagsFromTemplate:
				val, ok := action.Value.(string)
				if !ok {
					return nil, common.NewValidationError("request.config.actions."+strconv.Itoa(i)+".value", "Value must be string.")
				}

				response["config.actions."+strconv.Itoa(i)+".value"], err = s.validateTpl(val, data)
				if err != nil {
					return nil, err
				}
			}
		}
	case eventfilter.RuleTypeChangeEntity:
		if r.Request.Config.Resource != "" {
			response["config.resource"], err = s.validateTpl(r.Request.Config.Resource, data)
			if err != nil {
				return nil, err
			}
		}

		if r.Request.Config.Component != "" {
			response["config.component"], err = s.validateTpl(r.Request.Config.Component, data)
			if err != nil {
				return nil, err
			}
		}

		if r.Request.Config.Connector != "" {
			response["config.connector"], err = s.validateTpl(r.Request.Config.Connector, data)
			if err != nil {
				return nil, err
			}
		}

		if r.Request.Config.ConnectorName != "" {
			response["config.connector_name"], err = s.validateTpl(r.Request.Config.ConnectorName, data)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, nil
}

func (s *store) GetTemplateVars() TemplateVarsResponse {
	envVars := s.tplConfigProvider.Get().Vars
	envVarKeys := make([]string, len(envVars))
	i := 0
	for k := range envVars {
		envVarKeys[i] = "." + libtemplate.EnvVar + "." + k
		i++
	}

	sort.Strings(envVarKeys)
	tplVars := make([]string, len(s.tplVars)+len(envVars))
	copy(tplVars, s.tplVars)
	copy(tplVars[len(s.tplVars):], envVarKeys)
	confTplVars := make([]string, len(s.tplVars)+len(envVars)+1)
	copy(confTplVars, s.tplVars)
	confTplVars[len(s.tplVars)] = s.tplVarExternalData
	copy(confTplVars[len(s.tplVars)+1:], envVarKeys)

	return TemplateVarsResponse{
		ExternalData: tplVars,
		Config:       confTplVars,
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
		Author:              r.Author,
		Description:         r.Description,
		Type:                r.Type,
		Priority:            r.Priority,
		Enabled:             r.Enabled,
		Config:              r.Config,
		ExternalData:        externalData,
		EventPattern:        r.EventPattern,
		EntityPatternFields: r.EntityPatternFieldsRequest.ToModel(),
		RRule:               r.RRule,
		Start:               r.Start,
		Stop:                r.Stop,
		ResolvedStart:       r.Start,
		ResolvedStop:        r.Stop,
		Exdates:             exdates,
		Exceptions:          r.Exceptions,
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

func (s *store) getTplData(ctx context.Context, request TemplateRequest) (eventfilter.Template, error) {
	var res eventfilter.Template
	var entity types.Entity
	entityID := request.Data.Event.GetEID()
	if entityID != "" {
		err := s.dbEntityCollection.FindOne(ctx, bson.M{"_id": entityID, "soft_deleted": nil}).Decode(&entity)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return res, err
		}
	}

	var err error
	var matched bool
	var eventRegexMatches match.EventRegexMatches
	var entityRegexMatches match.EntityRegexMatches
	if len(request.Request.EventPattern) > 0 {
		matched, eventRegexMatches, err = match.MatchEventPatternWithRegexMatches(request.Request.EventPattern, &request.Data.Event.Event)
		if err != nil {
			return res, common.NewValidationError("request.event_pattern", "EventPattern is invalid event pattern.")
		}

		if !matched {
			return res, common.NewValidationError("data.event", "Event is not matched to event pattern.")
		}
	}

	var entityPattern pattern.Entity
	if len(request.Request.CorporatePattern.EntityPattern) > 0 {
		entityPattern = request.Request.CorporatePattern.EntityPattern
	} else {
		entityPattern = request.Request.EntityPattern
	}

	if len(entityPattern) > 0 {
		if entity.ID == "" {
			return res, common.NewValidationError("data.event", "Corresponding entity doesn't exist.")
		}

		matched, entityRegexMatches, err = match.MatchEntityPatternWithRegexMatches(entityPattern, &entity)
		if err != nil {
			return res, common.NewValidationError("request.entity_pattern", "EntityPattern is invalid entity pattern.")
		}

		if !matched {
			return res, common.NewValidationError("data.event", "Corresponding entity is not matched to entity pattern.")
		}
	}

	res = eventfilter.Template{
		Event: &request.Data.Event.Event,
		RegexMatch: eventfilter.RegexMatch{
			EventRegexMatches: eventRegexMatches,
			Entity:            entityRegexMatches,
		},
	}
	if entity.ID != "" {
		res.Event.Entity = &entity
	}

	return res, nil
}

func (s *store) validateTpl(str string, data eventfilter.Template) (template.Response, error) {
	isValid, errReport, err := s.tplValidator.Validate(str, data)
	if err != nil {
		return template.Response{}, err
	}

	if errReport != nil && strings.Contains(errReport.Message, tplValErrNoEntity) && data.Event.Entity == nil {
		return template.Response{}, common.NewValidationError("data.event", "Corresponding entity doesn't exist.")
	}

	return template.Response{
		IsValid: isValid,
		Err:     errReport,
	}, nil
}
