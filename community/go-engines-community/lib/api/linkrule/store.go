package linkrule

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apiexternaldata "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/externaldatatable"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/request"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

const defaultCategoriesLimit = 100

type Store interface {
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	GetByID(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
	GetCategories(ctx context.Context, r CategoriesRequest) (*CategoryResponse, error)
	ValidateTemplates(ctx context.Context, request TemplateRequest) (map[string]template.ValidateResponse, error)
	GetTemplateVars() TemplateVarsResponse
}

type store struct {
	client                mongo.DbClient
	collection            mongo.DbCollection
	exdataTableCollection mongo.DbCollection
	alarmCollection       mongo.DbCollection
	entityCollection      mongo.DbCollection
	userCollection        mongo.DbCollection
	tplTestCollection     mongo.DbCollection
	authorProvider        author.Provider
	transformer           common.PatternFieldsTransformer
	tplValidator          validator.Validator
	tplExecutor           libtemplate.Executor
	tplConfigProvider     config.TemplateConfigProvider
	externalDataContainer *externaldata.GetterContainer
	enforcer              security.Enforcer

	defaultSearchByFields []string
	defaultSortBy         string
	alarmTplVars          []template.VarResponse
	entityTplVars         []template.VarResponse
	alarmExdataTplVars    []template.VarResponse
	entityExdataTplVars   []template.VarResponse
	dupErrorParser        validation.DuplicateErrorParser
}

func NewStore(
	dbClient mongo.DbClient,
	authorProvider author.Provider,
	transformer common.PatternFieldsTransformer,
	tplValidator validator.Validator,
	tplExecutor libtemplate.Executor,
	tplConfigProvider config.TemplateConfigProvider,
	externalDataContainer *externaldata.GetterContainer,
	enforcer security.Enforcer,
) Store {
	userTplVars := []template.VarResponse{
		{Name: "email", Value: "{{ .User.Email }}"},
		{Name: "username", Value: "{{ .User.Username }}"},
		{Name: "firstname", Value: "{{ .User.Firstname }}"},
		{Name: "lastname", Value: "{{ .User.Lastname }}"},
		{Name: "externalID", Value: "{{ .User.ExternalID }}"},
		{Name: "externalSource", Value: "{{ .User.Source }}"},
		{Name: "roleIDs", Value: "{{ range .User.Roles }}{{ . }}{{ end }}"},
	}

	return &store{
		client:                dbClient,
		collection:            dbClient.Collection(mongo.LinkRuleMongoCollection),
		exdataTableCollection: dbClient.Collection(mongo.ExternalDataTableCollection),
		alarmCollection:       dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection:      dbClient.Collection(mongo.EntityMongoCollection),
		userCollection:        dbClient.Collection(mongo.UserCollection),
		tplTestCollection:     dbClient.Collection(mongo.TemplateTestCollection),
		authorProvider:        authorProvider,
		transformer:           transformer,
		tplValidator:          tplValidator,
		tplExecutor:           tplExecutor,
		tplConfigProvider:     tplConfigProvider,
		externalDataContainer: externalDataContainer,
		enforcer:              enforcer,

		defaultSearchByFields: []string{"_id", "author.name", "name"},
		defaultSortBy:         "created",
		alarmTplVars: []template.VarResponse{
			{
				Name:  "alarms",
				Value: template.GetAlarmVars("{{ range .Alarms }}{{ ", " }}{{ end }}", "", true),
			},
			{
				Name:  "entities",
				Value: template.GetEntityVars("{{ range .Alarms }}{{ ", " }}{{ end }}", ".Entity", true),
			},
			{
				Name:  "externalData",
				Value: "{{ range .Alarms }}{{ .ExternalData.%reference% }}{{ end }}",
			},
			{
				Name:  "user",
				Value: userTplVars,
			},
		},
		entityTplVars: []template.VarResponse{
			{
				Name:  "entities",
				Value: template.GetEntityVars("{{ range .Entities }}{{ ", " }}{{ end }}", "", true),
			},
			{
				Name:  "externalData",
				Value: "{{ range .Entities }}{{ .ExternalData.%reference% }}{{ end }}",
			},
			{
				Name:  "user",
				Value: userTplVars,
			},
		},
		alarmExdataTplVars: []template.VarResponse{
			{
				Name:  "alarm",
				Value: template.GetAlarmVars("{{ ", " }}", "", false),
			},
			{
				Name:  "entity",
				Value: template.GetEntityVars("{{ ", " }}", ".Entity", false),
			},
		},
		entityExdataTplVars: []template.VarResponse{
			{
				Name:  "entity",
				Value: template.GetEntityVars("{{ ", " }}", "", false),
			},
		},
		dupErrorParser: validation.NewDuplicateErrorParser(map[string]string{
			"name": "Name already exists.",
		}),
	}
}

func (s *store) Insert(ctx context.Context, request EditRequest) (*Response, error) {
	now := datetime.NewCpsTime()

	model, err := s.transformRequestToModel(ctx, request)
	if err != nil {
		return nil, err
	}

	model.ID = utils.NewID()
	model.Created = now
	model.Updated = now

	var response *Response
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		err = s.transformPatternRequestsToModel(ctx, request, &model)
		if err != nil {
			return err
		}

		_, err = s.collection.InsertOne(ctx, model)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err)
			}

			return err
		}

		response, err = s.GetByID(ctx, model.ID)

		return err
	})

	return response, err
}

func (s *store) GetByID(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)
	pipeline = append(pipeline, apiexternaldata.GetRefParametersLookups()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		response := Response{}
		err = cursor.Decode(&response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *store) Find(ctx context.Context, request ListRequest) (*AggregationResult, error) {
	pipeline := s.authorProvider.Pipeline()
	filter := common.GetSearchQuery(request.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sort := common.GetSortQuery(cmp.Or(request.SortBy, s.defaultSortBy), request.Sort)
	project := apiexternaldata.GetRefParametersLookups()
	project = append(project, sort)
	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		request.Query,
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

func (s *store) Update(ctx context.Context, request EditRequest) (*Response, error) {
	now := datetime.NewCpsTime()

	model, err := s.transformRequestToModel(ctx, request)
	if err != nil {
		return nil, err
	}

	model.ID = request.ID
	model.Updated = now

	var response *Response
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		err = s.transformPatternRequestsToModel(ctx, request, &model)
		if err != nil {
			return err
		}

		res, err := s.collection.UpdateOne(
			ctx,
			bson.M{"_id": request.ID},
			bson.M{"$set": model},
		)
		if err != nil || res.MatchedCount == 0 {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err)
			}

			return err
		}

		_, err = s.tplTestCollection.UpdateMany(ctx, bson.M{"rule._id": request.ID, "type": template.TypeTestLinkRule}, bson.M{
			"$set": bson.M{"rule.name": model.Name},
		})
		if err != nil {
			return err
		}

		response, err = s.GetByID(ctx, model.ID)

		return err
	})

	return response, err
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	var deleted int64
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		_, err := logger.DeleteByFilter(ctx, bson.M{"rule._id": id, "type": template.TypeTestLinkRule}, userID,
			s.tplTestCollection)
		if err != nil {
			return err
		}

		deleted, err = logger.DeleteOne(ctx, id, userID, s.collection)

		return err
	})

	return deleted > 0, err
}

// GetCategories returns list of distinct categories
func (s *store) GetCategories(ctx context.Context, r CategoriesRequest) (*CategoryResponse, error) {
	pipeline := make([]bson.M, 0)
	if r.Type != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"type": r.Type}})
	}
	queryLimit := r.Limit
	if queryLimit == 0 {
		queryLimit = defaultCategoriesLimit
	}
	pipeline = append(pipeline,
		bson.M{"$unwind": "$links"},
	)
	if r.Search != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{
			"links.category": bson.Regex{
				Pattern: ".*" + regexp.QuoteMeta(r.Search) + ".*",
				Options: "i",
			},
		}})
	}
	pipeline = append(pipeline,
		bson.M{"$group": bson.M{
			"_id": "$links.category",
		}},
		bson.M{"$sort": bson.M{"_id": 1}},
		bson.M{"$limit": queryLimit},
		bson.M{"$group": bson.M{
			"_id": nil,
			"categories": bson.M{
				"$push": "$_id",
			},
		}},
		bson.M{"$project": bson.M{"_id": 0}},
	)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	resp := CategoryResponse{}
	if !cursor.Next(ctx) {
		return &resp, nil
	}

	if err = cursor.Decode(&resp); err != nil {
		return nil, err
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *store) ValidateTemplates(ctx context.Context, r TemplateRequest) (map[string]template.ValidateResponse, error) {
	response := make(map[string]template.ValidateResponse)
	if r.Rule.SourceCode != "" {
		return response, nil
	}

	user, alarms, entities, err := s.getTplData(ctx, r)
	if err != nil {
		return nil, err
	}

	response, alarms, entities, err = s.validateExdataTpls(ctx, r, response, alarms, entities)
	if err != nil {
		return nil, err
	}

	var data map[string]any
	switch r.Rule.Type {
	case link.TypeAlarm:
		data = map[string]any{
			"Alarms": alarms,
			"User":   user,
		}
	case link.TypeEntity:
		data = map[string]any{
			"Entities": entities,
			"User":     user,
		}
	}

	for i, l := range r.Rule.Links {
		prefix := "links." + strconv.Itoa(i)
		response[prefix+".url"], err = template.Validate(s.tplValidator, l.URL, data)
		if err != nil {
			return nil, err
		}

		response[prefix+".label"], err = template.Validate(s.tplValidator, l.Label, data)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

func (s *store) GetTemplateVars() TemplateVarsResponse {
	alarmTplVars := template.AddEnvVars(s.alarmTplVars, s.tplConfigProvider)
	entityTplVars := template.AddEnvVars(s.entityTplVars, s.tplConfigProvider)
	res := TemplateVarsResponse{}
	res.Alarm.URL = alarmTplVars
	res.Alarm.Label = alarmTplVars
	res.Alarm.ExternalData = template.AddEnvVars(s.alarmExdataTplVars, s.tplConfigProvider)
	res.Entity.URL = entityTplVars
	res.Entity.Label = entityTplVars
	res.Entity.ExternalData = template.AddEnvVars(s.entityExdataTplVars, s.tplConfigProvider)

	return res
}

func (s *store) transformRequestToModel(ctx context.Context, r EditRequest) (link.Rule, error) {
	externalData, err := apiexternaldata.TransformRefParameters(ctx, r.ExternalData, s.exdataTableCollection)
	if err != nil {
		return link.Rule{}, err
	}

	rule := link.Rule{
		Name:         r.Name,
		Type:         r.Type,
		Enabled:      *r.Enabled,
		Links:        r.Links,
		SourceCode:   r.SourceCode,
		ExternalData: externalData,
		Author:       r.Author,
		EntityPatternFields: r.EntityPatternFieldsRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInEntityPattern(mongo.LinkRuleMongoCollection),
		),
	}
	if r.Type == link.TypeAlarm {
		rule.AlarmPatternFields = r.AlarmPatternFieldsRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInAlarmPattern(mongo.LinkRuleMongoCollection),
			common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.LinkRuleMongoCollection),
		)
	}

	return rule, nil
}

func (s *store) transformPatternRequestsToModel(ctx context.Context, req EditRequest, model *link.Rule) error {
	transformedEntityPatternRequest, err := s.transformer.TransformEntityPatternFieldsRequest(ctx, req.EntityPatternFieldsRequest)
	if err != nil {
		return err
	}

	if req.Type == link.TypeAlarm {
		transformedAlarmPatternRequest, err := s.transformer.TransformAlarmPatternFieldsRequest(ctx, req.AlarmPatternFieldsRequest)
		if err != nil {
			return err
		}

		model.AlarmPatternFields = transformedAlarmPatternRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInAlarmPattern(mongo.LinkRuleMongoCollection),
			common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.LinkRuleMongoCollection),
		)
	}

	model.Aliases = transformedEntityPatternRequest.Aliases
	model.EntityPatternFields = transformedEntityPatternRequest.ToModelWithoutFields(
		common.GetForbiddenFieldsInEntityPattern(mongo.LinkRuleMongoCollection),
	)

	return nil
}

func (s *store) getTplData(ctx context.Context, r TemplateRequest) (link.User, []link.AlarmWithData, []link.EntityWithData, error) {
	var user link.User
	var alarm link.AlarmWithData
	var entity link.EntityWithData
	userID := r.TestData.User
	if r.TestData.Test != "" {
		test := template.TestModel{}
		err := s.tplTestCollection.
			FindOne(ctx, bson.M{"_id": r.TestData.Test, "type": template.TypeTestLinkRule, "rule._id": r.Rule.ID}).
			Decode(&test)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return user, nil, nil, common.NewValidationError("testdata.test", "Test doesn't exist.")
			}

			return user, nil, nil, err
		}

		if userID == "" {
			userID = test.Data.User
		}

		switch r.Rule.Type {
		case link.TypeAlarm:
			if test.Data.Alarm != nil {
				alarm = link.AlarmWithData{
					Alarm:  test.Data.Alarm.Alarm,
					Entity: test.Data.Alarm.Entity,
				}
			}
		case link.TypeEntity:
			if test.Data.Entity != nil {
				entity = link.EntityWithData{
					Entity: *test.Data.Entity,
				}
			}
		}
	}

	if userID == "" {
		userID = r.Author
	} else {
		ok, err := s.enforcer.Enforce(r.Author, apisecurity.PermAcl, securitymodel.PermissionRead)
		if err != nil {
			return user, nil, nil, err
		}

		if !ok {
			return user, nil, nil, common.NewValidationError("testdata.user", "User is not accessible.")
		}
	}

	user, err := s.findUser(ctx, userID)
	if err != nil {
		return user, nil, nil, err
	}

	switch r.Rule.Type {
	case link.TypeAlarm:
		if r.TestData.Alarm == "" {
			if alarm.ID == "" {
				return user, nil, nil, common.NewValidationError("testdata.alarm", "Alarm is missing.")
			}
		} else if r.TestData.Alarm != alarm.ID { // keep snapshot from the test
			alarm, err = s.findAlarm(ctx, r.TestData.Alarm)
			if err != nil {
				return user, nil, nil, err
			}
		}
	case link.TypeEntity:
		if r.TestData.Entity == "" {
			if entity.ID == "" {
				return user, nil, nil, common.NewValidationError("testdata.entity", "Entity is missing.")
			}
		} else if r.TestData.Entity != entity.ID { // keep snapshot from the test
			entity, err = s.findEntity(ctx, r.TestData.Entity)
			if err != nil {
				return user, nil, nil, err
			}
		}
	}

	var alarms []link.AlarmWithData
	var entities []link.EntityWithData
	if alarm.ID != "" {
		alarms = []link.AlarmWithData{alarm}
	}

	if entity.ID != "" {
		entities = []link.EntityWithData{entity}
	}

	return user, alarms, entities, nil
}

// findAlarm fetches alarm with related entities.
func (s *store) findAlarm(ctx context.Context, alarmID string) (link.AlarmWithData, error) {
	var res link.AlarmWithData
	cursor, err := s.alarmCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"_id": alarmID,
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$project": bson.M{
			"v.steps": 0,
		}},
	})
	if err != nil {
		return res, fmt.Errorf("cannot find alarm: %w", err)
	}

	defer cursor.Close(ctx)
	if !cursor.Next(ctx) {
		return res, common.NewValidationError("testdata.alarm", "Alarm doesn't exist.")
	}

	err = cursor.Decode(&res)
	if err != nil {
		return res, fmt.Errorf("cannot decode alarm: %w", err)
	}

	if err = cursor.Err(); err != nil {
		return res, fmt.Errorf("cannot fetch alarm: %w", err)
	}

	return res, nil
}

func (s *store) findEntity(ctx context.Context, entityID string) (link.EntityWithData, error) {
	var res link.EntityWithData
	err := s.entityCollection.FindOne(ctx, bson.M{"_id": entityID}).Decode(&res)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return res, common.NewValidationError("testdata.entity", "Entity doesn't exist.")
		}

		return res, fmt.Errorf("cannot fetch entities: %w", err)
	}

	return res, nil
}

func (s *store) findUser(ctx context.Context, userID string) (link.User, error) {
	user := link.User{}
	cursor, err := s.userCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": userID}},
		{"$addFields": bson.M{"username": "$name"}},
	})
	if err != nil {
		return user, fmt.Errorf("cannot find user: %w", err)
	}

	defer cursor.Close(ctx)
	if !cursor.Next(ctx) {
		return user, common.NewValidationError("testdata.user", "User doesn't exist.")
	}

	err = cursor.Decode(&user)
	if err != nil {
		return user, err
	}

	if err = cursor.Err(); err != nil {
		return user, fmt.Errorf("cannot fetch user: %w", err)
	}

	return user, nil
}

func (s *store) validateExdataTpls(
	ctx context.Context,
	r TemplateRequest,
	response map[string]template.ValidateResponse,
	alarms []link.AlarmWithData,
	entities []link.EntityWithData,
) (map[string]template.ValidateResponse, []link.AlarmWithData, []link.EntityWithData, error) {
	for i, d := range r.Rule.ExternalData {
		prefix := "external_data." + strconv.Itoa(i)
		isValid := true
		for k, v := range d.Regexp {
			resKey := prefix + ".regexp." + k
			for _, alarm := range alarms {
				vr, err := template.Validate(s.tplValidator, v, alarm)
				if err != nil {
					return nil, nil, nil, err
				}

				if !vr.IsValid {
					response[resKey] = vr
					isValid = false
					break
				}

				if _, ok := response[resKey]; !ok {
					response[resKey] = vr
				}
			}

			for _, entity := range entities {
				vr, err := template.Validate(s.tplValidator, v, entity)
				if err != nil {
					return nil, nil, nil, err
				}

				if !vr.IsValid {
					response[resKey] = vr
					isValid = false
					break
				}

				if _, ok := response[resKey]; !ok {
					response[resKey] = vr
				}
			}
		}

		for k, v := range d.Select {
			resKey := prefix + ".select." + k
			for _, alarm := range alarms {
				vr, err := template.Validate(s.tplValidator, v, alarm)
				if err != nil {
					return nil, nil, nil, err
				}

				if !vr.IsValid {
					response[resKey] = vr
					isValid = false
					break
				}

				if _, ok := response[resKey]; !ok {
					response[resKey] = vr
				}
			}

			for _, entity := range entities {
				vr, err := template.Validate(s.tplValidator, v, entity)
				if err != nil {
					return nil, nil, nil, err
				}

				if !vr.IsValid {
					response[resKey] = vr
					isValid = false
					break
				}

				if _, ok := response[resKey]; !ok {
					response[resKey] = vr
				}
			}
		}

		if !isValid {
			continue
		}

		var err error
		for j := range alarms {
			if alarms[j].ExternalData == nil {
				alarms[j].ExternalData = make(map[string]map[string]any, len(r.Rule.ExternalData))
			}

			alarms[j].ExternalData[d.Reference], err = s.processTableExdata(ctx, d, alarms[j], "rule.external_data."+strconv.Itoa(i))
			if err != nil {
				return nil, nil, nil, err
			}
		}

		for j := range entities {
			if entities[j].ExternalData == nil {
				entities[j].ExternalData = make(map[string]map[string]any, len(r.Rule.ExternalData))
			}

			entities[j].ExternalData[d.Reference], err = s.processTableExdata(ctx, d, entities[j], "rule.external_data."+strconv.Itoa(i))
			if err != nil {
				return nil, nil, nil, err
			}
		}
	}

	return response, alarms, entities, nil
}

func (s *store) processTableExdata(
	ctx context.Context,
	d template.TemplateRefParameters,
	tplData any,
	field string,
) (map[string]any, error) {
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

	params, err := apiexternaldata.TransformRefParameters(ctx, []externaldata.RefParameters{refParam}, s.exdataTableCollection)
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
