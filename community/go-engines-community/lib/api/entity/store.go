package entity

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/statesettings"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/perfdata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Find(ctx context.Context, r ListRequestWithPagination) (*AggregationResult, error)
	Toggle(ctx context.Context, id string, enabled bool) (bool, SimplifiedEntity, error)
	GetContextGraph(ctx context.Context, id string) (*ContextGraphResponse, error)
	Export(ctx context.Context, t export.Task) (export.DataCursor, error)
	CheckStateSetting(ctx context.Context, r CheckStateSettingRequest) (StateSettingResponse, error)
	GetStateSetting(ctx context.Context, id string) (StateSettingResponse, error)
}

type store struct {
	db                      mongo.DbClient
	dbExport                mongo.DbClient
	mainCollection          mongo.DbCollection
	archivedCollection      mongo.DbCollection
	stateSettingsCollection mongo.DbCollection
	timezoneConfigProvider  config.TimezoneConfigProvider
	decoder                 encoding.Decoder
}

func NewStore(db, dbExport mongo.DbClient, timezoneConfigProvider config.TimezoneConfigProvider, decoder encoding.Decoder) Store {
	return &store{
		db:                      db,
		dbExport:                dbExport,
		mainCollection:          db.Collection(mongo.EntityMongoCollection),
		archivedCollection:      db.Collection(mongo.ArchivedEntitiesMongoCollection),
		stateSettingsCollection: db.Collection(mongo.StateSettingsMongoCollection),
		timezoneConfigProvider:  timezoneConfigProvider,
		decoder:                 decoder,
	}
}

func (s *store) Find(ctx context.Context, r ListRequestWithPagination) (*AggregationResult, error) {
	location := s.timezoneConfigProvider.Get().Location
	now := datetime.CpsTime{Time: time.Now().In(location)}

	pipeline, err := s.getQueryBuilder().CreateListAggregationPipeline(ctx, r, now)
	if err != nil {
		return nil, err
	}

	cursor, err := s.mainCollection.Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))

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

	s.fillConnectorType(&res)
	s.fillPerfData(&res, r.PerfData)

	return &res, nil
}

func (s *store) Toggle(ctx context.Context, id string, enabled bool) (bool, SimplifiedEntity, error) {
	var isToggled bool
	var oldSimplifiedEntity SimplifiedEntity

	err := s.db.WithTransaction(ctx, func(ctx context.Context) error {
		isToggled = false
		oldSimplifiedEntity = SimplifiedEntity{}

		cursor, err := s.mainCollection.Aggregate(ctx, []bson.M{
			{"$match": bson.M{"_id": id}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$_id",
				"connectFromField":        "_id",
				"connectToField":          "component",
				"as":                      "resources",
				"restrictSearchWithMatch": bson.M{"type": types.EntityTypeResource},
				"maxDepth":                0,
			}},
			{"$project": bson.M{
				"_id":       1,
				"name":      1,
				"component": 1,
				"enabled":   1,
				"type":      1,
				"resources": bson.M{"$map": bson.M{"input": "$resources", "in": "$$this._id"}},
			}},
		})
		if err != nil {
			return err
		}
		if cursor.Next(ctx) {
			err = cursor.Decode(&oldSimplifiedEntity)
			if err != nil {
				return err
			}
		} else {
			return nil
		}

		_, err = s.mainCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"enabled": enabled}})
		if err != nil {
			return err
		}

		isToggled = oldSimplifiedEntity.Enabled != enabled
		return nil
	})

	if oldSimplifiedEntity.ID == "" {
		return false, SimplifiedEntity{}, nil
	}

	if isToggled && !enabled && oldSimplifiedEntity.Type == types.EntityTypeComponent {
		depLen := len(oldSimplifiedEntity.Resources)
		from := 0

		for to := canopsis.DefaultBulkSize; to <= depLen; to += canopsis.DefaultBulkSize {
			_, err = s.mainCollection.UpdateMany(
				ctx,
				bson.M{"_id": bson.M{"$in": oldSimplifiedEntity.Resources[from:to]}},
				bson.M{"$set": bson.M{"enabled": enabled}},
			)
			if err != nil {
				return isToggled, oldSimplifiedEntity, err
			}

			from = to
		}

		if from < depLen {
			_, err = s.mainCollection.UpdateMany(
				ctx,
				bson.M{"_id": bson.M{"$in": oldSimplifiedEntity.Resources[from:depLen]}},
				bson.M{"$set": bson.M{"enabled": enabled}},
			)
			if err != nil {
				return isToggled, oldSimplifiedEntity, err
			}
		}
	}

	return isToggled, oldSimplifiedEntity, err
}

func (s *store) GetContextGraph(ctx context.Context, id string) (*ContextGraphResponse, error) {
	entity := Entity{}
	err := s.mainCollection.
		FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"type": 1})).
		Decode(&entity)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	switch entity.Type {
	case types.EntityTypeResource:
		pipeline = append(pipeline, []bson.M{
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$component",
				"connectFromField":        "component",
				"connectToField":          "_id",
				"as":                      "component",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$connector",
				"connectFromField":        "connector",
				"connectToField":          "_id",
				"as":                      "connector",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$services",
				"connectFromField":        "services",
				"connectToField":          "_id",
				"as":                      "services",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$addFields": bson.M{
				"impact": bson.M{"$concatArrays": bson.A{
					bson.M{"$map": bson.M{"input": "$component", "in": "$$this._id"}},
					bson.M{"$map": bson.M{"input": "$services", "in": "$$this._id"}},
				}},
				"depends": bson.M{"$map": bson.M{"input": "$connector", "in": "$$this._id"}},
			}},
		}...)
	case types.EntityTypeComponent:
		pipeline = append(pipeline, []bson.M{
			{"$graphLookup": bson.M{
				"from":             mongo.EntityMongoCollection,
				"startWith":        "$_id",
				"connectFromField": "_id",
				"connectToField":   "component",
				"as":               "resources",
				"restrictSearchWithMatch": bson.M{
					"type":         types.EntityTypeResource,
					"soft_deleted": bson.M{"$exists": false},
				},
				"maxDepth": 0,
			}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$connector",
				"connectFromField":        "connector",
				"connectToField":          "_id",
				"as":                      "connector",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$services",
				"connectFromField":        "services",
				"connectToField":          "_id",
				"as":                      "services",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$addFields": bson.M{
				"impact": bson.M{"$concatArrays": bson.A{
					bson.M{"$map": bson.M{"input": "$connector", "in": "$$this._id"}},
					bson.M{"$map": bson.M{"input": "$services", "in": "$$this._id"}},
				}},
				"depends": bson.M{"$map": bson.M{"input": "$resources", "in": "$$this._id"}},
			}},
		}...)
	case types.EntityTypeConnector:
		pipeline = append(pipeline, []bson.M{
			{"$graphLookup": bson.M{
				"from":             mongo.EntityMongoCollection,
				"startWith":        "$_id",
				"connectFromField": "_id",
				"connectToField":   "connector",
				"as":               "resources",
				"restrictSearchWithMatch": bson.M{
					"type":         types.EntityTypeResource,
					"soft_deleted": bson.M{"$exists": false},
				},
				"maxDepth": 0,
			}},
			{"$graphLookup": bson.M{
				"from":             mongo.EntityMongoCollection,
				"startWith":        "$_id",
				"connectFromField": "_id",
				"connectToField":   "connector",
				"as":               "components",
				"restrictSearchWithMatch": bson.M{
					"type":         types.EntityTypeComponent,
					"soft_deleted": bson.M{"$exists": false},
				},
				"maxDepth": 0,
			}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$services",
				"connectFromField":        "services",
				"connectToField":          "_id",
				"as":                      "services",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$addFields": bson.M{
				"impact": bson.M{"$concatArrays": bson.A{
					bson.M{"$map": bson.M{"input": "$resources", "in": "$$this._id"}},
					bson.M{"$map": bson.M{"input": "$services", "in": "$$this._id"}},
				}},
				"depends": bson.M{"$map": bson.M{"input": "$components", "in": "$$this._id"}},
			}},
		}...)
	case types.EntityTypeService:
		pipeline = append(pipeline, []bson.M{
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$_id",
				"connectFromField":        "_id",
				"connectToField":          "services",
				"as":                      "depends",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$services",
				"connectFromField":        "services",
				"connectToField":          "_id",
				"as":                      "impact",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$addFields": bson.M{
				"impact":  bson.M{"$map": bson.M{"input": "$impact", "in": "$$this._id"}},
				"depends": bson.M{"$map": bson.M{"input": "$depends", "in": "$$this._id"}},
			}},
		}...)
	}

	pipeline = append(pipeline, bson.M{"$project": bson.M{
		"impact":  1,
		"depends": 1,
	}})
	cursor, err := s.mainCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		res := ContextGraphResponse{}
		err = cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
		return &res, nil
	}

	return nil, nil
}

func (s *store) Export(ctx context.Context, t export.Task) (export.DataCursor, error) {
	r := BaseFilterRequest{}
	err := s.decoder.Decode([]byte(t.Parameters), &r)
	if err != nil {
		return nil, err
	}

	now := datetime.NewCpsTime()
	pipeline, err := s.getQueryBuilder().CreateOnlyListAggregationPipeline(ctx, ListRequest{
		BaseFilterRequest: r,
		SearchBy:          t.Fields.Fields(),
	}, now)
	if err != nil {
		return nil, err
	}

	project := make(bson.M, len(t.Fields))
	for _, field := range t.Fields {
		found := false
		for anotherField := range project {
			if strings.HasPrefix(field.Name, anotherField+".") {
				found = true
				break
			} else if strings.HasPrefix(anotherField, field.Name+".") {
				delete(project, anotherField)
				break
			}
		}
		if !found {
			project[field.Name] = 1
		}
	}
	pipeline = append(pipeline, bson.M{"$project": project})
	collection := s.dbExport.Collection(mongo.EntityMongoCollection)
	cursor, err := collection.Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return nil, err
	}

	return export.NewMongoCursor(cursor, t.Fields.Fields(), nil), nil
}

func (s *store) CheckStateSetting(ctx context.Context, r CheckStateSettingRequest) (StateSettingResponse, error) {
	response := StateSettingResponse{}
	cursor, err := s.stateSettingsCollection.Find(
		ctx,
		bson.M{
			"enabled": true,
			"type":    r.Type,
			"method": bson.M{
				"$in": []string{statesetting.MethodInherited, statesetting.MethodDependencies},
			},
		},
		options.Find().SetSort(bson.M{"priority": 1}).SetProjection(bson.M{
			"title": 1, "method": 1, "entity_pattern": 1, "type": 1,
			"inherited_entity_pattern": 1, "state_thresholds": 1,
		}),
	)
	if err != nil {
		return response, err
	}

	defer cursor.Close(ctx)

	ent := types.Entity{
		ID:          r.ID,
		Name:        r.Name,
		Connector:   r.Connector,
		Type:        r.Type,
		Infos:       TransformInfosRequest(r.Infos),
		Category:    r.Category,
		ImpactLevel: r.ImpactLevel,
	}

	for cursor.Next(ctx) {
		var stateSetting statesetting.StateSetting
		err = cursor.Decode(&stateSetting)
		if err != nil {
			return response, err
		}

		matched, err := match.MatchEntityPattern(*stateSetting.EntityPattern, &ent)
		if err != nil {
			return response, err
		}

		if matched {
			return getStateSettingResponse(stateSetting), nil
		}
	}

	if r.Type == types.EntityTypeService {
		return s.getDefaultStateSettingForService(ctx)
	}

	return response, nil
}

func (s *store) GetStateSetting(ctx context.Context, id string) (StateSettingResponse, error) {
	var response StateSettingResponse
	cursor, err := s.mainCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$lookup": bson.M{
			"from":         mongo.StateSettingsMongoCollection,
			"localField":   "state_info._id",
			"foreignField": "_id",
			"as":           "state_setting",
		}},
		{"$unwind": bson.M{"path": "$state_setting", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.EntityCountersCollection,
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "counters",
		}},
		{"$unwind": bson.M{"path": "$counters", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.AlarmMongoCollection,
			"localField":   "_id",
			"foreignField": "d",
			"as":           "alarm",
			"pipeline": []bson.M{
				{"$match": bson.M{"v.resolved": nil}},
			},
		}},
		{"$unwind": bson.M{"path": "$alarm", "preserveNullAndEmptyArrays": true}},
		{"$project": bson.M{
			"type":           1,
			"state_setting":  1,
			"state_counters": "$counters.state",
			"state":          "$alarm.v.state.val",
		}},
	})
	if err != nil {
		return response, err
	}

	if !cursor.Next(ctx) {
		return response, ErrNoFound
	}

	defer cursor.Close(ctx)
	res := struct {
		Type          string                       `bson:"type"`
		StateSetting  statesetting.StateSetting    `bson:"state_setting"`
		StateCounters entitycounters.StateCounters `bson:"state_counters"`
		State         int                          `bson:"state"`
	}{}
	err = cursor.Decode(&res)
	if err != nil {
		return response, err
	}

	response = getStateSettingResponse(res.StateSetting)
	if response.ID == "" && res.Type == types.EntityTypeService {
		return s.getDefaultStateSettingForService(ctx)
	}

	if res.StateSetting.Method == statesetting.MethodDependencies && res.StateSetting.StateThresholds != nil {
		var stateThreshold statesetting.StateThreshold
		switch res.State {
		case types.AlarmStateCritical:
			if res.StateSetting.StateThresholds.Critical != nil {
				stateThreshold = *res.StateSetting.StateThresholds.Critical
			}
		case types.AlarmStateMajor:
			if res.StateSetting.StateThresholds.Major != nil {
				stateThreshold = *res.StateSetting.StateThresholds.Major
			}
		case types.AlarmStateMinor:
			if res.StateSetting.StateThresholds.Minor != nil {
				stateThreshold = *res.StateSetting.StateThresholds.Minor
			}
		case types.AlarmStateOK:
			if res.StateSetting.StateThresholds.OK != nil {
				stateThreshold = *res.StateSetting.StateThresholds.OK
			}
		default:
			return response, fmt.Errorf("unknown state %d of entity %q", res.State, id)
		}

		var thresholdStateDependsCount int
		switch stateThreshold.State {
		case types.AlarmStateTitleCritical:
			thresholdStateDependsCount = res.StateCounters.Critical
		case types.AlarmStateTitleMajor:
			thresholdStateDependsCount = res.StateCounters.Major
		case types.AlarmStateTitleMinor:
			thresholdStateDependsCount = res.StateCounters.Minor
		case types.AlarmStateTitleOK:
			thresholdStateDependsCount = res.StateCounters.Ok
		}

		if thresholdStateDependsCount > 0 {
			response.DependsCount = res.StateCounters.Critical + res.StateCounters.Major + res.StateCounters.Minor + res.StateCounters.Ok
			response.ThresholdState = stateThreshold.State
			response.ThresholdStateDependsCount = thresholdStateDependsCount
		}
	}

	return response, nil
}

func (s *store) fillConnectorType(result *AggregationResult) {
	if result == nil {
		return
	}
	for i := range result.Data {
		result.Data[i].fillConnectorType()
	}
}

func (s *store) getQueryBuilder() *MongoQueryBuilder {
	return NewMongoQueryBuilder(s.db)
}

func (s *store) fillPerfData(result *AggregationResult, perfData []string) {
	if len(perfData) == 0 {
		return
	}

	perfDataRe := perfdata.Parse(perfData)
	for i, entity := range result.Data {
		result.Data[i].FilteredPerfData = perfdata.Filter(perfData, perfDataRe, entity.PerfData)
	}
}

func (s *store) getDefaultStateSettingForService(ctx context.Context) (StateSettingResponse, error) {
	var response StateSettingResponse
	var stateSetting statesetting.StateSetting
	err := s.stateSettingsCollection.FindOne(ctx, bson.M{"_id": statesetting.ServiceID}).Decode(&stateSetting)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return response, nil
		}

		return response, err
	}

	return getStateSettingResponse(stateSetting), nil
}

func getStateSettingResponse(stateSetting statesetting.StateSetting) StateSettingResponse {
	response := StateSettingResponse{}
	response.ID = stateSetting.ID
	response.Title = stateSetting.Title
	response.Type = stateSetting.Type
	response.Method = stateSetting.Method
	response.InheritedEntityPattern = stateSetting.InheritedEntityPattern
	if stateSetting.StateThresholds != nil {
		response.StateThresholds = &statesettings.StateThresholds{}
		response.StateThresholds.Critical = convertStateThreshold(stateSetting.StateThresholds.Critical)
		response.StateThresholds.Major = convertStateThreshold(stateSetting.StateThresholds.Major)
		response.StateThresholds.Minor = convertStateThreshold(stateSetting.StateThresholds.Minor)
		response.StateThresholds.OK = convertStateThreshold(stateSetting.StateThresholds.OK)
	}

	return response
}

func convertStateThreshold(src *statesetting.StateThreshold) *statesettings.StateThreshold {
	if src == nil {
		return nil
	}
	return &statesettings.StateThreshold{
		Method: src.Method,
		State:  src.State,
		Cond:   src.Cond,
		Value:  src.Value,
	}
}
