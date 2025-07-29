package scenario

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"strconv"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/priority"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/webhook"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Scenario, error)
	Find(ctx context.Context, q FilteredQuery) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Scenario, error)
	Update(ctx context.Context, r UpdateRequest) (*Scenario, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
	ValidateTemplates(ctx context.Context, request TemplateRequest) (map[string]template.ValidateResponse, error)
	GetTemplateVars() TemplateVarsResponse
}

type store struct {
	dbClient              mongo.DbClient
	collection            mongo.DbCollection
	alarmCollection       mongo.DbCollection
	tplDataCollection     mongo.DbCollection
	transformer           ModelTransformer
	authorProvider        author.Provider
	tplValidator          validator.Validator
	tplExecutor           libtemplate.Executor
	tplConfigProvider     config.TemplateConfigProvider
	defaultSearchByFields []string
	defaultSortBy         string
	outputTplVars         []template.VarResponse
	authorTplVars         []template.VarResponse
	whTplVars             []template.VarResponse
	firstWhTplVars        []template.VarResponse
	ticketTplVars         []template.VarResponse
}

func NewStore(
	db mongo.DbClient,
	authorProvider author.Provider,
	tplValidator validator.Validator,
	tplExecutor libtemplate.Executor,
	tplConfigProvider config.TemplateConfigProvider,
) Store {
	authorTplVars := template.GetAlarmVars("{{ ", " }}", ".Alarm", "", false)
	authorTplVars = append(authorTplVars, template.GetEntityVars("{{ ", " }}", ".Entity", "Entity", false)...)
	outputTplVars := make([]template.VarResponse, len(authorTplVars))
	copy(outputTplVars, authorTplVars)
	outputTplVars = append(outputTplVars,
		template.VarResponse{Name: "Trigger", Value: "{{ .AdditionalData.Trigger }}"},
		template.VarResponse{Name: "Author", Value: "{{ .AdditionalData.Author }}"},
		template.VarResponse{Name: "User id", Value: "{{ .AdditionalData.User }}"},
		template.VarResponse{Name: "Trigger event message", Value: "{{ .AdditionalData.Output }}"},
		template.VarResponse{Name: "Action initiator", Value: "{{ .AdditionalData.Initiator }}"},
		template.VarResponse{Name: "Rule name", Value: "{{ .AdditionalData.RuleName }}"},
	)
	ticketTplVars := []template.VarResponse{
		{Name: "Header field", Value: "{{ index .Header \"%field_name%\" }}"},
		{Name: "Response field", Value: "{{ index .Response \"%field_name%\" }}"},
		{Name: "Response field from the step #N", Value: "{{ index .ResponseMap \"%N%.%field_name%\" }}"},
	}
	firstWhTplVars := make([]template.VarResponse, len(outputTplVars))
	copy(firstWhTplVars, outputTplVars)
	firstWhTplVars = append(firstWhTplVars, template.GetAlarmVars("{{ range .Children }}{{ ", " }}{{ end }}", "", "Consequence alarm", true)...)
	firstWhTplVars = append(firstWhTplVars, template.GetEntityVars("{{ range .Children }}{{ ", " }}{{ end }}", ".Entity", "Consequence alarm entity", true)...)
	whTplVars := make([]template.VarResponse, 0, len(firstWhTplVars))
	whTplVars = append(whTplVars, firstWhTplVars...)
	whTplVars = append(whTplVars,
		template.VarResponse{Name: "Header field from the previous steps", Value: "{{ index .Header \"%field_name%\" }}"},
		template.VarResponse{Name: "Response field from the previous steps", Value: "{{ index .Response \"%field_name%\" }}"},
		template.VarResponse{Name: "Response field from the step #N", Value: "{{ index .ResponseMap \"%N%.%field_name%\" }}"},
	)

	return &store{
		dbClient:              db,
		collection:            db.Collection(mongo.ScenarioMongoCollection),
		alarmCollection:       db.Collection(mongo.AlarmMongoCollection),
		tplDataCollection:     db.Collection(mongo.TemplateDataCollection),
		transformer:           NewModelTransformer(),
		authorProvider:        authorProvider,
		tplValidator:          tplValidator,
		tplExecutor:           tplExecutor,
		tplConfigProvider:     tplConfigProvider,
		defaultSearchByFields: []string{"_id", "name", "author.name"},
		defaultSortBy:         "created",

		outputTplVars:  outputTplVars,
		authorTplVars:  authorTplVars,
		firstWhTplVars: firstWhTplVars,
		whTplVars:      whTplVars,
		ticketTplVars:  ticketTplVars,
	}
}

func (s *store) Find(ctx context.Context, r FilteredQuery) (*AggregationResult, error) {
	pipeline := s.authorProvider.Pipeline()
	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		s.getSort(r),
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	res := AggregationResult{}

	if cursor.Next(ctx) {
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (s *store) GetOneBy(ctx context.Context, id string) (*Scenario, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if cursor.Next(ctx) {
		res := &Scenario{}
		err := cursor.Decode(res)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Scenario, error) {
	now := datetime.NewCpsTime()
	model := s.transformer.TransformEditRequestToModel(r.EditRequest)

	if r.ID == "" {
		r.ID = utils.NewID()
	}

	model.ID = r.ID

	model.Created = now
	model.Updated = now
	var result *Scenario

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = nil
		_, err := s.collection.InsertOne(ctx, model)
		if err != nil {
			return err
		}

		err = priority.UpdateFollowing(ctx, s.collection, model.ID, model.Priority)
		if err != nil {
			return err
		}

		result, err = s.GetOneBy(ctx, model.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Scenario, error) {
	now := datetime.NewCpsTime()
	model := s.transformer.TransformEditRequestToModel(r.EditRequest)
	model.Updated = now

	var result *Scenario

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = nil
		res, err := s.collection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": model})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		err = priority.UpdateFollowing(ctx, s.collection, r.ID, r.Priority)
		if err != nil {
			return err
		}

		result, err = s.GetOneBy(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	var deleted int64

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		res, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		deleted, err = s.collection.DeleteOne(ctx, bson.M{"_id": id})
		return err
	})

	return deleted > 0, err
}

func (s *store) ValidateTemplates(ctx context.Context, r TemplateRequest) (map[string]template.ValidateResponse, error) {
	event, err := template.GetEventData(ctx, s.tplDataCollection, r.Data.Event)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, common.NewValidationError("data.event", "Event doesn't exist.")
	}

	var whTestData map[int]template.ResponseTestData
	if len(r.Data.Responses) > 0 {
		if len(r.Data.Responses) > len(r.Request.Actions) {
			return nil, common.NewValidationError("data.responses."+strconv.Itoa(len(r.Request.Actions)), "Response is redundant.")
		}

		whTestData, err = template.GetResponseData(ctx, s.tplDataCollection, r.Data.Responses)
		if err != nil {
			return nil, err
		}

		if len(whTestData) == 0 {
			return nil, common.NewValidationError("data.responses", "Responses don't exist.")
		}
	}

	alarm, err := s.findAlarm(ctx, event.GetEID())
	if err != nil {
		return nil, err
	}

	trigger := ""
	if len(r.Request.Triggers) > 0 {
		trigger = r.Request.Triggers[0].String()
	}

	additionalData := types.AdditionalData{
		Trigger:         trigger,
		AlarmChangeType: trigger,
		Initiator:       cmp.Or(event.Initiator, types.InitiatorExternal),
		Output:          event.Output,
		RuleName:        r.Request.Name,
	}

	return s.validateActionTpls(r, *event, alarm, additionalData, whTestData)
}

func (s *store) GetTemplateVars() TemplateVarsResponse {
	return TemplateVarsResponse{
		Output:       template.AddEnvVars(s.outputTplVars, s.tplConfigProvider),
		Author:       template.AddEnvVars(s.authorTplVars, s.tplConfigProvider),
		FirstWebhook: template.AddEnvVars(s.firstWhTplVars, s.tplConfigProvider),
		Webhook:      template.AddEnvVars(s.whTplVars, s.tplConfigProvider),
		Ticket:       template.AddEnvVars(s.ticketTplVars, s.tplConfigProvider),
	}
}

func (s *store) getSort(r FilteredQuery) bson.M {
	sortBy := cmp.Or(r.SortBy, s.defaultSortBy)
	if sortBy == "delay" {
		sortBy = "delay.value"
	}

	return common.GetSortQuery(sortBy, r.Sort)
}

func (s *store) findAlarm(ctx context.Context, entityID string) (webhook.TplAlarm, error) {
	if entityID == "" {
		return webhook.TplAlarm{}, common.NewValidationError("data.event", "Corresponding entity cannot be found.")
	}

	id := struct {
		ID string `bson:"_id"`
	}{}
	err := s.alarmCollection.FindOne(ctx, bson.M{"d": entityID}, options.FindOne().SetSort(bson.M{"t": -1}).SetProjection(bson.M{"_id": 1})).Decode(&id)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return webhook.TplAlarm{}, common.NewValidationError("data.event", "Corresponding alarm doesn't exist.")
		}

		return webhook.TplAlarm{}, err
	}

	cursor, err := s.alarmCollection.Aggregate(ctx, webhook.FetchAlarmsForTplPipeline([]string{id.ID}))
	if err != nil {
		return webhook.TplAlarm{}, fmt.Errorf("cannot fetch alarms: %w", err)
	}

	var alarms []webhook.TplAlarm
	err = cursor.All(ctx, &alarms)
	if err != nil {
		return webhook.TplAlarm{}, fmt.Errorf("cannot decode alarms: %w", err)
	}

	if len(alarms) == 0 {
		return webhook.TplAlarm{}, common.NewValidationError("data.alarm", "Alarm doesn't exist.")
	}

	return alarms[0], nil
}

func (s *store) validateActionTpls(
	r TemplateRequest,
	event types.Event,
	alarm webhook.TplAlarm,
	additionalData types.AdditionalData,
	whTestData map[int]template.ResponseTestData,
) (map[string]template.ValidateResponse, error) {
	response := make(map[string]template.ValidateResponse)
	whResponse := make(map[string]any)
	whResponseMap := make(map[string]any)
	whHeader := make(map[string]string)
	authorTplData := types.AlarmWithEntity{
		Alarm:  alarm.Alarm,
		Entity: alarm.Entity,
	}
	var err error
	for i, a := range r.Request.Actions {
		prefix := "actions." + strconv.Itoa(i) + ".parameters"
		var authorRes *template.ValidateResponse
		authorRes, additionalData, err = s.validateAuthorTpl(a, event, additionalData, authorTplData)
		if err != nil {
			return nil, err
		}

		if authorRes != nil {
			response[prefix+".author"] = *authorRes
		}

		switch a.Type {
		case types.ActionTypeWebhook:
			whTplData := webhook.NewTplData(false, []webhook.TplAlarm{alarm}, additionalData, whResponse, whResponseMap, whHeader)
			response[prefix+".request.url"], err = template.Validate(s.tplValidator, a.Parameters.Request.URL, whTplData)
			if err != nil {
				return nil, err
			}

			response[prefix+".request.payload"], err = template.Validate(s.tplValidator, a.Parameters.Request.Payload, whTplData)
			if err != nil {
				return nil, err
			}

			for k, v := range a.Parameters.Request.Headers {
				response[prefix+".request.headers."+k], err = template.Validate(s.tplValidator, v, whTplData)
				if err != nil {
					return nil, err
				}
			}

			iStr := strconv.Itoa(i)
			if td, ok := whTestData[i]; ok {
				whResponse = td.Body
				whHeader = td.Headers
				for k, v := range whResponse {
					whResponseMap[iStr+"."+k] = v
				}
			} else {
				return nil, common.NewValidationError("data.responses."+iStr, "Response is missing.")
			}

			if a.Parameters.DeclareTicket != nil {
				ticketTplData := map[string]any{
					"Response":    whResponse,
					"ResponseMap": whResponseMap,
					"Header":      whHeader,
				}

				if a.Parameters.DeclareTicket.TicketIDTpl != "" {
					response[prefix+".declare_ticket.ticket_id_tpl"], err = template.Validate(s.tplValidator, a.Parameters.DeclareTicket.TicketIDTpl, ticketTplData)
					if err != nil {
						return nil, err
					}
				}

				if a.Parameters.DeclareTicket.TicketURLTpl != "" {
					response[prefix+".declare_ticket.ticket_url_tpl"], err = template.Validate(s.tplValidator, a.Parameters.DeclareTicket.TicketURLTpl, ticketTplData)
					if err != nil {
						return nil, err
					}
				}
			}
		default:
			response[prefix+".output"], err = template.Validate(s.tplValidator, a.Parameters.Output, map[string]any{
				"Alarm":          alarm.Alarm,
				"Entity":         alarm.Entity,
				"AdditionalData": additionalData,
			})
			if err != nil {
				return nil, err
			}

			if _, ok := whTestData[i]; ok {
				return nil, common.NewValidationError("data.responses."+strconv.Itoa(i), "Response is redundant.")
			}
		}
	}

	return response, nil
}

func (s *store) validateAuthorTpl(
	a ActionRequest,
	event types.Event,
	additionalData types.AdditionalData,
	tplData any,
) (*template.ValidateResponse, types.AdditionalData, error) {
	if a.Parameters.ForwardAuthor != nil && *a.Parameters.ForwardAuthor {
		additionalData.Author = cmp.Or(event.Author, canopsis.DefaultEventAuthor)
		additionalData.User = event.UserID

		return nil, additionalData, nil
	}

	if a.Parameters.Author == "" {
		additionalData.Author = canopsis.DefaultEventAuthor
		additionalData.User = ""

		return &template.ValidateResponse{IsValid: true}, additionalData, nil
	}

	res, err := template.Validate(s.tplValidator, a.Parameters.Author, tplData)
	if err != nil {
		return nil, additionalData, err
	}

	if res.IsValid {
		additionalData.Author, err = s.tplExecutor.Execute(a.Parameters.Author, tplData)
		if err != nil {
			return nil, additionalData, err
		}

		additionalData.User = ""
	}

	return &res, additionalData, nil
}

func getNestedObjectsPipeline() []bson.M {
	return []bson.M{
		{"$unwind": bson.M{
			"path":                       "$actions",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "action_index",
		}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"localField":   "actions.parameters.type",
			"foreignField": "_id",
			"as":           "actions.parameters.type",
		}},
		{"$unwind": bson.M{"path": "$actions.parameters.type", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorReasonMongoCollection,
			"localField":   "actions.parameters.reason",
			"foreignField": "_id",
			"as":           "actions.parameters.reason",
		}},
		{"$unwind": bson.M{"path": "$actions.parameters.reason", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"action_index": 1}},
		{"$group": bson.M{
			"_id":     "$_id",
			"data":    bson.M{"$first": "$$ROOT"},
			"actions": bson.M{"$push": "$actions"},
		}},
		{"$replaceRoot": bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				"$data",
				bson.M{"actions": "$actions"},
			}},
		}},
	}
}
