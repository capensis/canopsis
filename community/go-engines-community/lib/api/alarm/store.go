package alarm

//go:generate mockgen -destination=../../../mocks/lib/api/alarm/alarm.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm Store

import (
	"context"
	"errors"
	"sort"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/perfdata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	InstructionExecutionStatusRunning = iota
	InstructionExecutionStatusPaused
	InstructionExecutionStatusCompleted
	InstructionExecutionStatusAborted
	InstructionExecutionStatusFailed
	InstructionExecutionStatusWaitResult
)

const (
	InstructionTypeManual = iota
	InstructionTypeAuto
	InstructionTypeSimplifiedManual
)

const InstructionStatusApproved = 0

type Store interface {
	Find(ctx context.Context, r ListRequestWithPagination, userId string) (*AggregationResult, error)
	GetAssignedInstructionsMap(ctx context.Context, alarmIds []string) (map[string][]AssignedInstruction, error)
	GetInstructionExecutionStatuses(ctx context.Context, alarmIDs []string, assignedInstructionsMap map[string][]AssignedInstruction) (map[string]ExecutionStatus, error)
	Count(ctx context.Context, r FilterRequest, userID string) (*Count, error)
	GetByID(ctx context.Context, id, userId string, onlyParents bool) (*Alarm, error)
	GetOpenByEntityID(ctx context.Context, id, userId string) (*Alarm, bool, error)
	FindByService(ctx context.Context, id string, r ListByServiceRequest, userId string) (*AggregationResult, error)
	FindByComponent(ctx context.Context, r ListByComponentRequest, userId string) (*AggregationResult, error)
	FindResolved(ctx context.Context, r ResolvedListRequest, userId string) (*AggregationResult, error)
	GetDetails(ctx context.Context, r DetailsRequest, userId string) (*Details, error)
	GetAssignedDeclareTicketsMap(ctx context.Context, alarmIds []string) (map[string][]AssignedDeclareTicketRule, error)
	Export(ctx context.Context, t export.Task) (export.DataCursor, error)
	GetLinks(ctx context.Context, ruleId string, alarmIds []string, userId string) ([]link.Link, bool, error)
}

type store struct {
	dbClient                         mongo.DbClient
	dbExportClient                   mongo.DbClient
	mainDbCollection                 mongo.DbCollection
	resolvedDbCollection             mongo.DbCollection
	dbInstructionCollection          mongo.DbCollection
	dbInstructionExecutionCollection mongo.DbCollection
	dbEntityCollection               mongo.DbCollection
	dbDeclareTicketCollection        mongo.DbCollection
	dbUserCollection                 mongo.DbCollection
	authorProvider                   author.Provider

	linkGenerator link.Generator

	timezoneConfigProvider config.TimezoneConfigProvider

	tplExecutor template.Executor

	decoder encoding.Decoder

	logger zerolog.Logger
}

func NewStore(
	dbClient,
	dbExportClient mongo.DbClient,
	linkGenerator link.Generator,
	timezoneConfigProvider config.TimezoneConfigProvider,
	authorProvider author.Provider,
	tplExecutor template.Executor,
	decoder encoding.Decoder,
	logger zerolog.Logger,
) Store {
	return &store{
		dbClient:                         dbClient,
		dbExportClient:                   dbExportClient,
		mainDbCollection:                 dbClient.Collection(mongo.AlarmMongoCollection),
		resolvedDbCollection:             dbClient.Collection(mongo.ResolvedAlarmMongoCollection),
		dbInstructionCollection:          dbClient.Collection(mongo.InstructionMongoCollection),
		dbInstructionExecutionCollection: dbClient.Collection(mongo.InstructionExecutionMongoCollection),
		dbEntityCollection:               dbClient.Collection(mongo.EntityMongoCollection),
		dbDeclareTicketCollection:        dbClient.Collection(mongo.DeclareTicketRuleMongoCollection),
		dbUserCollection:                 dbClient.Collection(mongo.UserCollection),
		authorProvider:                   authorProvider,

		linkGenerator: linkGenerator,

		timezoneConfigProvider: timezoneConfigProvider,

		tplExecutor: tplExecutor,

		decoder: decoder,

		logger: logger,
	}
}

func (s *store) Find(ctx context.Context, r ListRequestWithPagination, userId string) (*AggregationResult, error) {
	collection := s.mainDbCollection
	if r.GetOpenedFilter() == OnlyResolved {
		collection = s.resolvedDbCollection
	}

	now := datetime.NewCpsTime()
	pipeline, err := s.getQueryBuilder().CreateListAggregationPipeline(ctx, r, now, userId)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result AggregationResult
	for cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, s.postProcessResult(ctx, &result, r.WithDeclareTickets, r.WithInstructions, r.WithLinks, r.OnlyParents, userId)
}

func (s *store) GetByID(ctx context.Context, id, userId string, onlyParents bool) (*Alarm, error) {
	now := datetime.NewCpsTime()
	pipeline, err := s.getQueryBuilder().CreateGetAggregationPipeline(bson.M{"_id": id}, now, userId,
		OpenedAndRecentResolved, onlyParents)
	if err != nil {
		return nil, err
	}

	cursor, err := s.mainDbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := AggregationResult{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	if len(result.Data) == 0 {
		resolvedCursor, err := s.resolvedDbCollection.Aggregate(ctx, pipeline)
		if err != nil {
			return nil, err
		}
		defer resolvedCursor.Close(ctx)

		if resolvedCursor.Next(ctx) {
			err = resolvedCursor.Decode(&result)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(result.Data) == 0 {
		return nil, nil
	}

	return &result.Data[0], s.postProcessResult(ctx, &result, true, true, true, false, userId)
}

func (s *store) GetOpenByEntityID(ctx context.Context, entityID, userId string) (*Alarm, bool, error) {
	err := s.dbEntityCollection.FindOne(ctx, bson.M{
		"_id":     entityID,
		"enabled": true,
	}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, false, nil
		}
		return nil, false, err
	}

	now := datetime.NewCpsTime()
	pipeline, err := s.getQueryBuilder().CreateGetAggregationPipeline(bson.M{"d": entityID}, now, userId,
		OnlyOpened, false)
	if err != nil {
		return nil, false, err
	}

	cursor, err := s.mainDbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, false, err
	}
	defer cursor.Close(ctx)

	result := AggregationResult{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, false, err
		}
	}

	if len(result.Data) == 0 {
		return nil, true, nil
	}

	return &result.Data[0], true, s.postProcessResult(ctx, &result, true, true, true, false, userId)
}

func (s *store) FindByService(ctx context.Context, id string, r ListByServiceRequest, userId string) (*AggregationResult, error) {
	now := datetime.NewCpsTime()
	service := types.Entity{}
	err := s.dbEntityCollection.FindOne(ctx, bson.M{
		"_id":     id,
		"type":    types.EntityTypeService,
		"enabled": true,
	}).Decode(&service)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	entityMatch := bson.M{"entity.services": service.ID}
	if r.WithService {
		entityMatch = bson.M{"$or": []bson.M{
			{"entity._id": service.ID},
			entityMatch,
		}}
	}

	opened := true
	pipeline, err := s.getQueryBuilder().CreateAggregationPipelineByMatch(ctx,
		nil,
		entityMatch,
		r.Query,
		r.SortRequest,
		FilterRequest{BaseFilterRequest: BaseFilterRequest{
			Opened:   &opened,
			Category: r.Category,
			Search:   r.Search,
		}},
		now,
		userId,
	)
	if err != nil {
		return nil, err
	}

	cursor, err := s.mainDbCollection.Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result AggregationResult
	for cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, s.postProcessResult(ctx, &result, true, true, true, false, userId)
}

func (s *store) FindByComponent(ctx context.Context, r ListByComponentRequest, userId string) (*AggregationResult, error) {
	now := datetime.NewCpsTime()
	component := types.Entity{}
	err := s.dbEntityCollection.FindOne(ctx, bson.M{
		"_id":     r.ID,
		"type":    types.EntityTypeComponent,
		"enabled": true,
	}).Decode(&component)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	opened := true
	pipeline, err := s.getQueryBuilder().CreateAggregationPipelineByMatch(ctx,
		nil,
		bson.M{"entity.component": component.ID},
		r.Query,
		r.SortRequest,
		FilterRequest{BaseFilterRequest: BaseFilterRequest{
			Opened: &opened,
		}},
		now,
		userId,
	)
	if err != nil {
		return nil, err
	}

	cursor, err := s.mainDbCollection.Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result AggregationResult
	for cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, s.postProcessResult(ctx, &result, true, true, true, false, userId)
}

func (s *store) FindResolved(ctx context.Context, r ResolvedListRequest, userId string) (*AggregationResult, error) {
	now := datetime.NewCpsTime()

	err := s.dbEntityCollection.FindOne(ctx, bson.M{
		"_id":     r.ID,
		"enabled": true,
	}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	opened := false
	pipeline, err := s.getQueryBuilder().CreateAggregationPipelineByMatch(ctx,
		bson.M{"d": r.ID},
		nil,
		r.Query,
		r.SortRequest,
		FilterRequest{BaseFilterRequest: BaseFilterRequest{
			StartFrom: r.StartFrom,
			StartTo:   r.StartTo,
			Opened:    &opened,
		}},
		now,
		userId,
	)
	if err != nil {
		return nil, err
	}

	cursor, err := s.resolvedDbCollection.Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result AggregationResult
	for cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	err = s.fillLinks(ctx, &result, userId)
	if err != nil {
		s.logger.Err(err).Msg("cannot fill links")
	}

	return &result, nil
}

func (s *store) GetDetails(ctx context.Context, r DetailsRequest, userId string) (*Details, error) {
	now := datetime.NewCpsTime()
	match := bson.M{"_id": r.ID}
	collection := s.mainDbCollection
	switch r.GetOpenedFilter() {
	case OnlyOpened:
		match["v.resolved"] = nil
	case OnlyResolved:
		collection = s.resolvedDbCollection
	}

	const entityLookupName = "entity"
	pipeline := []bson.M{
		{"$match": match},
		{"$addFields": bson.M{
			"is_meta_alarm": bson.M{"$cond": bson.A{bson.M{"$not": bson.A{"$v.meta"}}, false, true}},
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           entityLookupName,
		}},
		{"$unwind": "$" + entityLookupName},
	}

	if r.WithDependencies {
		pipeline = append(pipeline, getImpactsCountPipeline()...)
	}

	if r.Steps != nil {
		stepMatch := bson.M{}
		if r.Steps.Type != "" {
			stepMatch["v.steps._t"] = r.Steps.Type
		}

		var stepsArray any = "$v.steps"
		if len(stepMatch) > 0 {
			pipeline = append(pipeline,
				bson.M{"$unwind": bson.M{
					"path":                       "$v.steps",
					"preserveNullAndEmptyArrays": true,
					"includeArrayIndex":          "step_index",
				}},
				bson.M{"$match": stepMatch},
			)
			if r.Steps.Reversed {
				pipeline = append(pipeline, bson.M{"$sort": bson.M{"step_index": -1}})
			} else {
				pipeline = append(pipeline, bson.M{"$sort": bson.M{"step_index": 1}})
			}
			pipeline = append(pipeline,
				bson.M{"$group": bson.M{
					"_id":   "$_id",
					"data":  bson.M{"$first": "$$ROOT"},
					"steps": bson.M{"$push": "$v.steps"},
				}},
				bson.M{"$replaceRoot": bson.M{"newRoot": bson.M{
					"$mergeObjects": bson.A{
						"$data",
						bson.M{"v": bson.M{
							"$mergeObjects": bson.A{
								"$v",
								bson.M{"steps": "$steps"},
							},
						}},
					},
				}}},
			)
		} else if r.Steps.Reversed {
			stepsArray = bson.M{"$reverseArray": "$v.steps"}
		}

		pipeline = append(pipeline, bson.M{"$addFields": bson.M{
			"steps.data": bson.M{"$slice": bson.A{
				stepsArray,
				(r.Steps.Page - 1) * r.Steps.Limit,
				r.Steps.Limit},
			},
			"steps_count": bson.M{"$size": "$v.steps"},
		}})
	}

	pipeline = append(pipeline, bson.M{"$project": bson.M{
		"v.steps": 0,
	}})

	var details Details
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&details)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	if r.Steps != nil {
		details.Steps.Meta, err = common.NewPaginatedMeta(r.Steps.Query, details.StepsCount)
		if err != nil {
			return nil, err
		}
	}

	if r.Children != nil {
		children := AggregationResult{
			Data: make([]Alarm, 0),
		}

		if details.IsMetaAlarm {
			childrenPipeline, err := s.getQueryBuilder().CreateChildrenAggregationPipeline(*r.Children,
				r.GetOpenedFilter(), details.Entity.ID, r.Search, userId, r.SearchBy, now)
			if err != nil {
				return nil, err
			}

			childrenCursor, err := collection.Aggregate(ctx, childrenPipeline)
			if err != nil {
				return nil, err
			}
			defer childrenCursor.Close(ctx)
			if childrenCursor.Next(ctx) {
				err = childrenCursor.Decode(&children)
				if err != nil {
					return nil, err
				}
			}

			err = s.postProcessResult(ctx, &children, r.WithDeclareTickets, r.WithInstructions, true, false, userId)
			if err != nil {
				return nil, err
			}
		}

		meta, err := common.NewPaginatedMeta(r.Children.Query, children.TotalCount)
		if err != nil {
			return nil, err
		}
		details.Children = &ChildrenDetails{
			Data: children.Data,
			Meta: meta,
		}
	}

	if len(r.PerfData) > 0 {
		perfDataRe := perfdata.Parse(r.PerfData)
		details.FilteredPerfData = perfdata.Filter(r.PerfData, perfDataRe, details.Entity.PerfData)
	}

	return &details, nil
}

func (s *store) Count(ctx context.Context, r FilterRequest, userID string) (*Count, error) {
	collection := s.mainDbCollection
	if r.GetOpenedFilter() == OnlyResolved {
		collection = s.resolvedDbCollection
	}

	pipeline, err := s.getQueryBuilder().CreateCountAggregationPipeline(ctx, r, userID, datetime.NewCpsTime())
	if err != nil {
		return nil, err
	}

	totalPipeline := []bson.M{{"$count": "count"}}

	totalActivePipeline := []bson.M{
		{"$match": bson.M{"$and": []bson.M{
			{"v.snooze": bson.M{"$exists": false}},
			{"$or": []bson.M{
				{"v.pbehavior_info": bson.M{"$exists": false}},
				{"v.pbehavior_info.canonical_type": bson.M{"$eq": pbehavior.TypeActive}},
			}},
		}}},
		{"$count": "count"},
	}

	totalSnoozePipeline := []bson.M{
		{"$match": bson.M{"v.snooze": bson.M{"$exists": true}}},
		{"$count": "count"},
	}

	totalAckPipeline := []bson.M{
		{"$match": bson.M{"v.ack": bson.M{"$exists": true}}},
		{"$count": "count"},
	}

	totalTicketPipeline := []bson.M{
		{"$match": bson.M{"v.ticket": bson.M{"$exists": true}}},
		{"$count": "count"},
	}

	totalPbehaviorPipeline := []bson.M{
		{"$match": bson.M{"$and": []bson.M{
			{"v.pbehavior_info": bson.M{"$exists": true}},
			{"v.pbehavior_info.canonical_type": bson.M{"$ne": pbehavior.TypeActive}},
		}}},
		{"$count": "count"},
	}

	aggregationPipeline := append(pipeline,
		bson.M{"$facet": bson.M{
			"total":           totalPipeline,
			"total_active":    totalActivePipeline,
			"total_snooze":    totalSnoozePipeline,
			"total_ack":       totalAckPipeline,
			"total_ticket":    totalTicketPipeline,
			"total_pbehavior": totalPbehaviorPipeline,
		}},
		bson.M{"$addFields": bson.M{
			"counts": bson.M{
				"$arrayToObject": bson.M{
					"$map": bson.M{
						"input": bson.M{"$objectToArray": "$$ROOT"},
						"as":    "each",
						"in": bson.M{
							"k": "$$each.k",
							"v": bson.M{"$sum": "$$each.v.count"},
						},
					},
				},
			},
		}},
		bson.M{"$replaceRoot": bson.M{"newRoot": "$counts"}},
	)

	cursor, err := collection.Aggregate(ctx, aggregationPipeline)
	if err != nil {
		return nil, err
	}

	var result Count

	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			cursor.Close(ctx)
			return nil, err
		}
	}

	err = cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) GetAssignedInstructionsMap(ctx context.Context, alarmIds []string) (map[string][]AssignedInstruction, error) {
	m, _, err := s.getAssignedInstructionsMap(ctx, alarmIds)
	return m, err
}

func (s *store) Export(ctx context.Context, t export.Task) (export.DataCursor, error) {
	r := ExportFetchParameters{}
	err := s.decoder.Decode([]byte(t.Parameters), &r)
	if err != nil {
		return nil, err
	}

	collectionName := mongo.AlarmMongoCollection
	if r.GetOpenedFilter() == OnlyResolved {
		collectionName = mongo.ResolvedAlarmMongoCollection
	}

	now := datetime.NewCpsTime()
	pipeline, err := s.getQueryBuilder().CreateOnlyListAggregationPipeline(ctx, ListRequest{
		FilterRequest: FilterRequest{
			BaseFilterRequest: r.BaseFilterRequest,
			SearchBy:          t.Fields.Fields(),
		},
	}, now, t.User)
	if err != nil {
		return nil, err
	}

	project := make(bson.M, len(t.Fields))
	withInstructions := false
	withLinks := false
	withModel := false
	for _, field := range t.Fields {
		if field.Name == "assigned_instructions" {
			withInstructions = true
			withModel = true
			continue
		}
		if field.Name == "links" || strings.HasPrefix(field.Name, "links.") {
			withLinks = true
			withModel = true
			continue
		}
		if field.Template != "" {
			withModel = true
			continue
		}

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

	if withModel {
		project["model"] = bson.M{
			"alarm": "$$ROOT",
			"entity": bson.M{"$mergeObjects": bson.A{
				"$entity",
				bson.M{"category": "$entity.category._id"},
			}},
		}
	}

	var instructions []Instruction
	if withInstructions {
		cursor, err := s.dbInstructionCollection.Find(ctx, bson.M{
			"type":   bson.M{"$in": bson.A{InstructionTypeManual, InstructionTypeSimplifiedManual}},
			"status": bson.M{"$in": bson.A{InstructionStatusApproved, nil}},
		}, options.Find().SetProjection(bson.M{"steps": 0}))
		if err != nil {
			return nil, err
		}
		err = cursor.All(ctx, &instructions)
		if err != nil {
			return nil, err
		}
	}

	pipeline = append(pipeline, bson.M{"$project": project})
	collection := s.dbExportClient.Collection(collectionName)
	cursor, err := collection.Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return nil, err
	}

	location := s.timezoneConfigProvider.Get().Location
	var linkGenerator link.Generator
	var user link.User
	if withLinks {
		linkGenerator = s.linkGenerator
		user, err = s.findUser(ctx, t.User)
		if err != nil {
			return nil, err
		}
	}
	exportCursor := newExportCursor(cursor, t.Fields, common.GetRealFormatTime(r.TimeFormat), location,
		instructions, linkGenerator, user, s.tplExecutor, withModel, s.logger)
	return exportCursor, nil
}

func (s *store) GetLinks(ctx context.Context, ruleId string, alarmIds []string, userId string) ([]link.Link, bool, error) {
	user, err := s.findUser(ctx, userId)
	if err != nil {
		return nil, false, err
	}

	links, err := s.linkGenerator.GenerateCombinedForAlarmsByRule(ctx, ruleId, alarmIds, user)
	if err != nil {
		if errors.Is(err, link.ErrNoRule) {
			return nil, false, nil
		}
		if errors.Is(err, link.ErrNotMatchedAlarm) {
			return nil, false, common.NewValidationError("ids", "Alarms aren't matched to rule.")
		}
		return nil, false, err
	}

	sort.Slice(links, func(i, j int) bool {
		return links[i].Label < links[j].Label
	})

	if links == nil {
		return []link.Link{}, true, nil
	}

	return links, true, nil
}

func (s *store) getAssignedInstructionsMap(ctx context.Context, alarmIds []string) (map[string][]AssignedInstruction, bson.M, error) {
	instructionCursor, err := s.dbInstructionCollection.Aggregate(
		ctx,
		[]bson.M{
			{"$match": bson.M{
				"type":    bson.M{"$in": bson.A{InstructionTypeManual, InstructionTypeSimplifiedManual}},
				"status":  bson.M{"$in": bson.A{InstructionStatusApproved, nil}},
				"enabled": true,
			}},
			{"$lookup": bson.M{
				"from":         mongo.InstructionExecutionMongoCollection,
				"localField":   "_id",
				"foreignField": "instruction",
				"as":           "executions",
			}},
			{"$addFields": bson.M{
				"executions": bson.M{"$filter": bson.M{
					"input": "$executions",
					"cond": bson.M{"$and": []bson.M{
						{"$in": bson.A{"$$this.status", []int{InstructionExecutionStatusRunning, InstructionExecutionStatusPaused}}},
						{"$in": bson.A{"$$this.alarm", alarmIds}},
					}},
				}},
			}},
			{"$addFields": bson.M{
				"executions": bson.M{
					"$map": bson.M{
						"input": "$executions",
						"in": bson.M{
							"_id":    "$$this._id",
							"alarm":  "$$this.alarm",
							"status": "$$this.status",
						},
					},
				},
			}},
			{"$project": bson.M{
				"steps": 0,
			}},
		},
	)
	if err != nil {
		return nil, nil, err
	}

	defer instructionCursor.Close(ctx)

	instructionMap := make(map[string]InstructionWithExecutions, canopsis.FacetLimit)
	instructionFiltersPipeline := make(bson.M, canopsis.FacetLimit)
	assignedInstructionsMap := make(map[string][]AssignedInstruction)
	allInstructionMatches := make([]bson.M, 0)

	for instructionCursor.Next(ctx) {
		var instruction InstructionWithExecutions
		err = instructionCursor.Decode(&instruction)
		if err != nil {
			return nil, nil, err
		}

		q, err := getInstructionQuery(instruction.Instruction)
		if err != nil {
			return nil, nil, err
		}

		if len(instructionFiltersPipeline) > canopsis.FacetLimit {
			err = s.processInstructionFiltersPipeline(ctx, alarmIds, instructionMap, instructionFiltersPipeline, assignedInstructionsMap)
			if err != nil {
				return nil, nil, err
			}

			instructionMap = make(map[string]InstructionWithExecutions, canopsis.FacetLimit)
			instructionFiltersPipeline = make(bson.M, canopsis.FacetLimit)
		}

		if q != nil {
			instructionMap[instruction.ID] = instruction
			instructionFiltersPipeline[instruction.ID] = []bson.M{{"$match": q}}
			allInstructionMatches = append(allInstructionMatches, q)
		}
	}

	anyInstructionMatch := bson.M{}
	if len(allInstructionMatches) > 0 {
		anyInstructionMatch = bson.M{"$or": allInstructionMatches}
	}

	if len(instructionFiltersPipeline) > 0 {
		err = s.processInstructionFiltersPipeline(ctx, alarmIds, instructionMap, instructionFiltersPipeline, assignedInstructionsMap)
	}

	return assignedInstructionsMap, anyInstructionMatch, err
}

func (s *store) processInstructionFiltersPipeline(
	ctx context.Context,
	alarmIds []string,
	instructionMap map[string]InstructionWithExecutions,
	instructionFiltersPipeline bson.M,
	assignedInstructionsMap map[string][]AssignedInstruction,
) error {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": alarmIds}}},
		{"$addFields": bson.M{
			"v.infos_array": bson.M{"$objectToArray": "$v.infos"},
			"v.duration": bson.M{"$subtract": bson.A{
				datetime.NewCpsTime(),
				"$v.creation_date",
			}},
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": bson.M{"path": "$entity", "preserveNullAndEmptyArrays": true}},
		{"$facet": instructionFiltersPipeline},
		{"$addFields": bson.M{
			"ids": bson.M{
				"$arrayToObject": bson.M{
					"$map": bson.M{
						"input": bson.M{"$objectToArray": "$$ROOT"},
						"as":    "each",
						"in": bson.M{
							"k": "$$each.k",
							"v": bson.M{"$map": bson.M{
								"input": "$$each.v",
								"as":    "e",
								"in":    "$$e._id",
							}},
						},
					},
				},
			},
		}},
		{"$unwind": "$ids"},
		{"$replaceRoot": bson.M{"newRoot": "$ids"}},
	}

	assignedInstructionsCursor, err := s.mainDbCollection.Aggregate(
		ctx,
		pipeline,
	)
	if err != nil {
		return err
	}

	defer assignedInstructionsCursor.Close(ctx)

	for assignedInstructionsCursor.Next(ctx) {
		assignedInstructions := make(map[string][]string)
		err = assignedInstructionsCursor.Decode(&assignedInstructions)
		if err != nil {
			return err
		}

		for instructionId, alarmIds := range assignedInstructions {
			for _, alarmId := range alarmIds {
				execution := instructionMap[instructionId].GetExecution(alarmId)
				assignedInstructionsMap[alarmId] = append(assignedInstructionsMap[alarmId], AssignedInstruction{
					ID:        instructionId,
					Name:      instructionMap[instructionId].Name,
					Type:      instructionMap[instructionId].Type,
					Execution: execution,
				})
			}
		}
	}

	return nil
}

func (s *store) GetInstructionExecutionStatuses(ctx context.Context, alarmIDs []string, assignedInstructionsMap map[string][]AssignedInstruction) (map[string]ExecutionStatus, error) {
	if len(alarmIDs) == 0 {
		return nil, nil
	}

	leftAlarms := make(map[string]struct{}, len(alarmIDs))
	for _, id := range alarmIDs {
		leftAlarms[id] = struct{}{}
	}

	cursor, err := s.dbInstructionExecutionCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"alarm": bson.M{"$in": alarmIDs},
			},
		},
		{
			"$sort": bson.M{
				"started_at": -1,
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"alarm":       "$alarm",
					"instruction": "$instruction",
				},
				"instruction_id":   bson.M{"$first": "$instruction"},
				"instruction_name": bson.M{"$first": "$name"},
				"instruction_type": bson.M{"$first": "$type"},
				"status":           bson.M{"$first": "$status"},
				"started_at":       bson.M{"$first": "$started_at"},
			},
		},
		{
			"$sort": bson.M{
				"started_at": -1,
			},
		},
		{
			"$group": bson.M{
				"_id": "$_id.alarm",
				"all_failed": bson.M{
					"$push": bson.M{
						"$cond": bson.M{
							"if": bson.M{
								"$eq": bson.A{"$status", InstructionExecutionStatusFailed},
							},
							"then": "$instruction_type",
							"else": "$$REMOVE",
						},
					},
				},
				"all_successful": bson.M{
					"$push": bson.M{
						"$cond": bson.M{
							"if": bson.M{
								"$eq": bson.A{"$status", InstructionExecutionStatusCompleted},
							},
							"then": "$instruction_type",
							"else": "$$REMOVE",
						},
					},
				},
				"running_manual_instructions": bson.M{
					"$push": bson.M{
						"$cond": bson.M{
							"if": bson.M{
								"$and": bson.A{
									bson.M{
										"$in": bson.A{"$status", bson.A{InstructionExecutionStatusRunning, InstructionExecutionStatusWaitResult}},
									},
									bson.M{
										"$in": bson.A{"$instruction_type", bson.A{InstructionTypeManual, InstructionTypeSimplifiedManual}},
									},
								},
							},
							"then": "$instruction_name",
							"else": "$$REMOVE",
						},
					},
				},
				"running_auto_instructions": bson.M{
					"$push": bson.M{
						"$cond": bson.M{
							"if": bson.M{
								"$and": bson.A{
									bson.M{
										"$in": bson.A{"$status", bson.A{InstructionExecutionStatusRunning, InstructionExecutionStatusWaitResult}},
									},
									bson.M{
										"$eq": bson.A{"$instruction_type", InstructionTypeAuto},
									},
								},
							},
							"then": "$instruction_name",
							"else": "$$REMOVE",
						},
					},
				},
				"failed_manual_instructions": bson.M{
					"$push": bson.M{
						"$cond": bson.M{
							"if": bson.M{
								"$and": bson.A{
									bson.M{
										"$eq": bson.A{"$status", InstructionExecutionStatusFailed},
									},
									bson.M{
										"$in": bson.A{"$instruction_type", bson.A{InstructionTypeManual, InstructionTypeSimplifiedManual}},
									},
								},
							},
							"then": "$instruction_name",
							"else": "$$REMOVE",
						},
					},
				},
				"failed_auto_instructions": bson.M{
					"$push": bson.M{
						"$cond": bson.M{
							"if": bson.M{
								"$and": bson.A{
									bson.M{
										"$eq": bson.A{"$status", InstructionExecutionStatusFailed},
									},
									bson.M{
										"$eq": bson.A{"$instruction_type", InstructionTypeAuto},
									},
								},
							},
							"then": "$instruction_name",
							"else": "$$REMOVE",
						},
					},
				},
				"successful_manual_instructions": bson.M{
					"$addToSet": bson.M{
						"$cond": bson.M{
							"if": bson.M{
								"$and": bson.A{
									bson.M{
										"$eq": bson.A{"$status", InstructionExecutionStatusCompleted},
									},
									bson.M{
										"$in": bson.A{"$instruction_type", bson.A{InstructionTypeManual, InstructionTypeSimplifiedManual}},
									},
								},
							},
							"then": "$instruction_name",
							"else": "$$REMOVE",
						},
					},
				},
				"successful_auto_instructions": bson.M{
					"$addToSet": bson.M{
						"$cond": bson.M{
							"if": bson.M{
								"$and": bson.A{
									bson.M{
										"$eq": bson.A{"$status", InstructionExecutionStatusCompleted},
									},
									bson.M{
										"$eq": bson.A{"$instruction_type", InstructionTypeAuto},
									},
								},
							},
							"then": "$instruction_name",
							"else": "$$REMOVE",
						},
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"last_failed":     bson.M{"$arrayElemAt": bson.A{"$all_failed", 0}},
				"last_successful": bson.M{"$arrayElemAt": bson.A{"$all_successful", 0}},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	var executionStatuses []ExecutionStatus
	err = cursor.All(ctx, &executionStatuses)
	if err != nil {
		return nil, err
	}
	statusesByAlarm := make(map[string]ExecutionStatus, len(executionStatuses))
	for _, v := range executionStatuses {
		delete(leftAlarms, v.ID)

		v.Icon = getInstructionExecutionIcon(v, assignedInstructionsMap)
		statusesByAlarm[v.ID] = v
	}

	for alarmID := range leftAlarms {
		if _, ok := assignedInstructionsMap[alarmID]; ok {
			statusesByAlarm[alarmID] = ExecutionStatus{
				Icon: IconManualAvailable,
			}
		}
	}

	return statusesByAlarm, nil
}

func getInstructionExecutionIcon(status ExecutionStatus, assignedInstructionsMap map[string][]AssignedInstruction) int {
	availableInstructionsMap := make(map[string]struct{}, len(assignedInstructionsMap[status.ID]))
	for _, instr := range assignedInstructionsMap[status.ID] {
		availableInstructionsMap[instr.Name] = struct{}{}
	}

	for _, name := range status.RunningManualInstructions {
		delete(availableInstructionsMap, name)
	}

	for _, name := range status.FailedManualInstructions {
		delete(availableInstructionsMap, name)
	}

	for _, name := range status.SuccessfulManualInstructions {
		delete(availableInstructionsMap, name)
	}

	runningManualInstruction := len(status.RunningManualInstructions) != 0
	runningAutoInstruction := len(status.RunningAutoInstructions) != 0
	failedManualInstruction := len(status.FailedManualInstructions) != 0
	failedAutoInstruction := len(status.FailedAutoInstructions) != 0
	successfulManualInstruction := len(status.SuccessfulManualInstructions) != 0
	successfulAutoInstruction := len(status.SuccessfulAutoInstructions) != 0
	availableInstructions := len(availableInstructionsMap) != 0
	lastFailed := status.LastFailed
	lastSuccessful := status.LastSuccessful

	if (failedManualInstruction || failedAutoInstruction) && lastFailed != nil {
		if *lastFailed == InstructionTypeAuto {
			if runningManualInstruction || runningAutoInstruction {
				return IconAutoFailedOtherInProgress
			} else if availableInstructions {
				return IconAutoFailedManualAvailable
			} else {
				return IconAutoFailed
			}
		} else {
			if runningManualInstruction || runningAutoInstruction {
				return IconManualFailedOtherInProgress
			} else if availableInstructions {
				return IconManualFailedManualAvailable
			} else {
				return IconManualFailed
			}
		}
	}

	if (successfulManualInstruction || successfulAutoInstruction) && lastSuccessful != nil {
		if *lastSuccessful == InstructionTypeAuto {
			if runningManualInstruction || runningAutoInstruction {
				return IconAutoSuccessfulOtherInProgress
			} else if availableInstructions {
				return IconAutoSuccessfulManualAvailable
			} else {
				return IconAutoSuccessful
			}
		} else {
			if runningManualInstruction || runningAutoInstruction {
				return IconManualSuccessfulOtherInProgress
			} else if availableInstructions {
				return IconManualSuccessfulManualAvailable
			} else {
				return IconManualSuccessful
			}
		}
	}

	if runningManualInstruction {
		return IconManualInProgress
	}

	if runningAutoInstruction {
		return IconAutoInProgress
	}

	if availableInstructions {
		return IconManualAvailable
	}

	return NoIcon
}

func (s *store) fillAssignedInstructions(ctx context.Context, result *AggregationResult) (map[string][]AssignedInstruction, bson.M, error) {
	var alarmIds []string
	for _, item := range result.Data {
		if item.Value.Resolved == nil {
			alarmIds = append(alarmIds, item.ID)
		}
	}

	if len(alarmIds) == 0 {
		return nil, nil, nil
	}

	assignedInstructionsMap, anyInstructionMatch, err := s.getAssignedInstructionsMap(ctx, alarmIds)
	if err != nil {
		return nil, nil, err
	}

	for i, alarmDocument := range result.Data {
		sort.Slice(assignedInstructionsMap[alarmDocument.ID], func(i, j int) bool {
			return assignedInstructionsMap[alarmDocument.ID][i].Name < assignedInstructionsMap[alarmDocument.ID][j].Name
		})

		assignedInstructions := assignedInstructionsMap[alarmDocument.ID]
		if len(assignedInstructions) == 0 {
			assignedInstructions = make([]AssignedInstruction, 0)
		}
		result.Data[i].AssignedInstructions = &assignedInstructions
	}

	return assignedInstructionsMap, anyInstructionMatch, nil
}

func (s *store) fillInstructionExecutionStatusesAndIcon(ctx context.Context, result *AggregationResult, assignedInstructions map[string][]AssignedInstruction) error {
	alarmIDs := make([]string, len(result.Data))
	for i, item := range result.Data {
		alarmIDs[i] = item.ID
	}
	if len(alarmIDs) == 0 {
		return nil
	}

	executionStatuses, err := s.GetInstructionExecutionStatuses(ctx, alarmIDs, assignedInstructions)
	if err != nil {
		return err
	}

	for i, v := range result.Data {
		result.Data[i].InstructionExecutionIcon = executionStatuses[v.ID].Icon
		result.Data[i].RunningManualInstructions = executionStatuses[v.ID].RunningManualInstructions
		result.Data[i].RunningAutoInstructions = executionStatuses[v.ID].RunningAutoInstructions
		result.Data[i].FailedManualInstructions = executionStatuses[v.ID].FailedManualInstructions
		result.Data[i].FailedAutoInstructions = executionStatuses[v.ID].FailedAutoInstructions
		result.Data[i].SuccessfulManualInstructions = executionStatuses[v.ID].SuccessfulManualInstructions
		result.Data[i].SuccessfulAutoInstructions = executionStatuses[v.ID].SuccessfulAutoInstructions
	}

	return nil
}

func (s *store) fillChildrenInstructionsFlag(ctx context.Context, result *AggregationResult, anyInstructionMatch bson.M) error {
	if len(anyInstructionMatch) == 0 {
		return nil
	}
	parentIds := make([]string, 0)
	requiredParents := make(map[string]bool)
	for _, v := range result.Data {
		if v.IsMetaAlarm != nil && *v.IsMetaAlarm && v.Value.Resolved == nil {
			parentIds = append(parentIds, v.Entity.ID)
			requiredParents[v.Entity.ID] = true
		}
	}
	if len(parentIds) == 0 {
		return nil
	}

	pipeline := []bson.M{
		{"$match": bson.M{
			"v.parents":  bson.M{"$in": parentIds},
			"v.resolved": nil,
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$match": anyInstructionMatch},
		{"$project": bson.M{
			"v.parents": 1,
		}},
	}
	cursor, err := s.mainDbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	hasInstructions := make(map[string]bool, len(parentIds))
	for cursor.Next(ctx) {
		alarm := Alarm{}
		err := cursor.Decode(&alarm)
		if err != nil {
			return err
		}
		for _, v := range alarm.Value.Parents {
			if !requiredParents[v] {
				continue
			}

			hasInstructions[v] = true
		}
	}

	for i := range result.Data {
		childrenInstructions := hasInstructions[result.Data[i].Entity.ID]
		result.Data[i].ChildrenInstructions = &childrenInstructions
	}

	return nil
}

// fillLinks sends a request to API v2 and fills result with links from a response.
func (s *store) fillLinks(ctx context.Context, result *AggregationResult, userId string) error {
	if result == nil || len(result.Data) == 0 {
		return nil
	}

	user, err := s.findUser(ctx, userId)
	if err != nil {
		return err
	}

	ids := make([]string, len(result.Data))
	for i, v := range result.Data {
		ids[i] = v.ID
	}

	linksByAlarmId, err := s.linkGenerator.GenerateForAlarms(ctx, ids, user)
	if err != nil || len(linksByAlarmId) == 0 {
		return err
	}

	for i, v := range result.Data {
		result.Data[i].Links = linksByAlarmId[v.ID]
		for _, links := range result.Data[i].Links {
			sort.Slice(links, func(i, j int) bool {
				return links[i].Label < links[j].Label
			})
		}
	}

	return nil
}

func (s *store) getQueryBuilder() *MongoQueryBuilder {
	return NewMongoQueryBuilder(s.dbClient, s.authorProvider)
}

func (s *store) fillAssignedDeclareTickets(ctx context.Context, result *AggregationResult) error {
	var alarmIDs []string
	for _, item := range result.Data {
		if item.Value.Resolved == nil {
			alarmIDs = append(alarmIDs, item.ID)
		}
	}

	if len(alarmIDs) == 0 {
		return nil
	}

	assignedRulesMap, err := s.GetAssignedDeclareTicketsMap(ctx, alarmIDs)
	if err != nil {
		return err
	}

	for idx, v := range result.Data {
		sort.Slice(assignedRulesMap[v.ID], func(i, j int) bool {
			return assignedRulesMap[v.ID][i].Name < assignedRulesMap[v.ID][j].Name
		})

		result.Data[idx].AssignedDeclareTicketRules = assignedRulesMap[v.ID]
	}

	return nil
}

func (s *store) GetAssignedDeclareTicketsMap(ctx context.Context, alarmIds []string) (map[string][]AssignedDeclareTicketRule, error) {
	declareTicketCursor, err := s.dbDeclareTicketCollection.Find(ctx, bson.M{"enabled": true})
	if err != nil {
		return nil, err
	}

	defer declareTicketCursor.Close(ctx)

	ruleMap := make(map[string]AssignedDeclareTicketRule, canopsis.FacetLimit)
	rulePipeline := make(bson.M, canopsis.FacetLimit)
	assignedRulesMap := make(map[string][]AssignedDeclareTicketRule)

	for declareTicketCursor.Next(ctx) {
		var rule DeclareTicketRule
		err = declareTicketCursor.Decode(&rule)
		if err != nil {
			return nil, err
		}

		q, err := rule.getDeclareTicketQuery()
		if err != nil {
			return nil, err
		}

		if len(rulePipeline) == canopsis.FacetLimit {
			err = s.processPipeline(ctx, alarmIds, ruleMap, rulePipeline, assignedRulesMap)
			if err != nil {
				return nil, err
			}

			ruleMap = make(map[string]AssignedDeclareTicketRule, canopsis.FacetLimit)
			rulePipeline = make(bson.M, canopsis.FacetLimit)
		}

		if q != nil {
			ruleMap[rule.ID] = AssignedDeclareTicketRule{ID: rule.ID, Name: rule.Name}
			rulePipeline[rule.ID] = []bson.M{{"$match": q}}
		}
	}

	if len(rulePipeline) > 0 {
		err = s.processPipeline(ctx, alarmIds, ruleMap, rulePipeline, assignedRulesMap)
	}

	return assignedRulesMap, err
}

func (s *store) findUser(ctx context.Context, id string) (link.User, error) {
	user := link.User{}
	cursor, err := s.dbUserCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$addFields": bson.M{"username": "$name"}},
	})
	if err != nil {
		return user, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&user)
		return user, err
	}

	return user, errors.New("user not found")
}

func (s *store) processPipeline(
	ctx context.Context,
	alarmIDs []string,
	ruleMap map[string]AssignedDeclareTicketRule,
	rulePipeline bson.M,
	assignedRulesMap map[string][]AssignedDeclareTicketRule,
) error {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": alarmIDs}}},
		{"$addFields": bson.M{
			"v.infos_array": bson.M{"$objectToArray": "$v.infos"},
			"v.duration": bson.M{"$subtract": bson.A{
				datetime.NewCpsTime(),
				"$v.creation_date",
			}},
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": bson.M{"path": "$entity", "preserveNullAndEmptyArrays": true}},
		{"$facet": rulePipeline},
		{"$addFields": bson.M{
			"ids": bson.M{
				"$arrayToObject": bson.M{
					"$map": bson.M{
						"input": bson.M{"$objectToArray": "$$ROOT"},
						"as":    "each",
						"in": bson.M{
							"k": "$$each.k",
							"v": bson.M{"$map": bson.M{
								"input": "$$each.v",
								"as":    "e",
								"in":    "$$e._id",
							}},
						},
					},
				},
			},
		}},
		{"$unwind": "$ids"},
		{"$replaceRoot": bson.M{"newRoot": "$ids"}},
	}

	assignedRulesCursor, err := s.mainDbCollection.Aggregate(
		ctx,
		pipeline,
	)
	if err != nil {
		return err
	}

	defer assignedRulesCursor.Close(ctx)

	for assignedRulesCursor.Next(ctx) {
		assignedRules := make(map[string][]string)
		err = assignedRulesCursor.Decode(&assignedRules)
		if err != nil {
			return err
		}

		for ruleID, alarmIds := range assignedRules {
			for _, alarmId := range alarmIds {
				assignedRulesMap[alarmId] = append(assignedRulesMap[alarmId], ruleMap[ruleID])
			}
		}
	}

	return nil
}

func (s *store) postProcessResult(
	ctx context.Context,
	result *AggregationResult,
	withDeclareTicket, withInstructions, withLinks, onlyParents bool,
	userId string,
) error {
	if withDeclareTicket {
		err := s.fillAssignedDeclareTickets(ctx, result)
		if err != nil {
			return err
		}
	}

	if withInstructions {
		assignedInstructionMap, anyInstructionMatch, err := s.fillAssignedInstructions(ctx, result)
		if err != nil {
			return err
		}
		err = s.fillInstructionExecutionStatusesAndIcon(ctx, result, assignedInstructionMap)
		if err != nil {
			return err
		}
		if onlyParents {
			err = s.fillChildrenInstructionsFlag(ctx, result, anyInstructionMatch)
			if err != nil {
				return err
			}
		}
	}

	if withLinks {
		err := s.fillLinks(ctx, result, userId)
		if err != nil {
			s.logger.Err(err).Msg("cannot fill links")
		}
	}

	return nil
}
