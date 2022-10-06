package alarm

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

const manualMetaAlarmsLimit = 100

type Store interface {
	Find(ctx context.Context, apiKey string, r ListRequestWithPagination) (*AggregationResult, error)
	GetAssignedInstructionsMap(ctx context.Context, alarmIds []string) (map[string][]AssignedInstruction, error)
	GetInstructionExecutionStatuses(ctx context.Context, alarmIDs []string) (map[string]ExecutionStatus, error)
	Count(ctx context.Context, r FilterRequest) (*Count, error)
	GetByID(ctx context.Context, id, apiKey string) (*Alarm, error)
	GetOpenByEntityID(ctx context.Context, id, apiKey string) (*Alarm, bool, error)
	FindManual(ctx context.Context, search string) ([]ManualResponse, error)
	FindByService(ctx context.Context, id, apiKey string, r ListByServiceRequest) (*AggregationResult, error)
	FindByComponent(ctx context.Context, r ListByComponentRequest, apiKey string) (*AggregationResult, error)
	FindResolved(ctx context.Context, r ResolvedListRequest, apiKey string) (*AggregationResult, error)
	GetDetails(ctx context.Context, apiKey string, r DetailsRequest) (*Details, error)
}

type store struct {
	dbClient                         mongo.DbClient
	mainDbCollection                 mongo.DbCollection
	resolvedDbCollection             mongo.DbCollection
	dbInstructionCollection          mongo.DbCollection
	dbInstructionExecutionCollection mongo.DbCollection
	dbEntityCollection               mongo.DbCollection

	queryBuilder *MongoQueryBuilder

	linksFetcher common.LinksFetcher

	logger zerolog.Logger
}

func NewStore(dbClient mongo.DbClient, linksFetcher common.LinksFetcher, logger zerolog.Logger) Store {
	return &store{
		dbClient:                         dbClient,
		mainDbCollection:                 dbClient.Collection(mongo.AlarmMongoCollection),
		resolvedDbCollection:             dbClient.Collection(mongo.ResolvedAlarmMongoCollection),
		dbInstructionCollection:          dbClient.Collection(mongo.InstructionMongoCollection),
		dbInstructionExecutionCollection: dbClient.Collection(mongo.InstructionExecutionMongoCollection),
		dbEntityCollection:               dbClient.Collection(mongo.EntityMongoCollection),

		queryBuilder: NewMongoQueryBuilder(dbClient),

		linksFetcher: linksFetcher,

		logger: logger,
	}
}

func (s *store) Find(ctx context.Context, apiKey string, r ListRequestWithPagination) (*AggregationResult, error) {
	collection := s.mainDbCollection
	if r.GetOpenedFilter() == OnlyResolved {
		collection = s.resolvedDbCollection
	}

	now := types.NewCpsTime()
	pipeline, err := s.queryBuilder.CreateListAggregationPipeline(ctx, r, now)
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

	if r.WithInstructions {
		anyInstructionMatch, err := s.fillAssignedInstructions(ctx, &result, now)
		if err != nil {
			return nil, err
		}
		err = s.fillInstructionFlags(ctx, &result)
		if err != nil {
			return nil, err
		}
		if r.OnlyParents {
			err = s.fillChildrenInstructionsFlag(ctx, &result, anyInstructionMatch)
			if err != nil {
				return nil, err
			}
		}
	}

	if r.WithLinks {
		err = s.fillLinks(ctx, apiKey, &result)
		if err != nil {
			s.logger.Err(err).Msg("cannot fill links")
		}
	}

	return &result, nil
}

func (s *store) GetByID(ctx context.Context, id, apiKey string) (*Alarm, error) {
	pipeline, err := s.queryBuilder.CreateGetAggregationPipeline(bson.M{"_id": id}, types.NewCpsTime())
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

	_, err = s.fillAssignedInstructions(ctx, &result, types.NewCpsTime())
	if err != nil {
		return nil, err
	}
	err = s.fillInstructionFlags(ctx, &result)
	if err != nil {
		return nil, err
	}

	err = s.fillLinks(ctx, apiKey, &result)
	if err != nil {
		s.logger.Err(err).Msg("cannot fill links")
	}

	return &result.Data[0], nil
}

func (s *store) GetOpenByEntityID(ctx context.Context, entityID, apiKey string) (*Alarm, bool, error) {
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

	pipeline, err := s.queryBuilder.CreateGetAggregationPipeline(bson.M{
		"d":          entityID,
		"v.resolved": nil,
	}, types.NewCpsTime())
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

	_, err = s.fillAssignedInstructions(ctx, &result, types.NewCpsTime())
	if err != nil {
		return nil, false, err
	}
	err = s.fillInstructionFlags(ctx, &result)
	if err != nil {
		return nil, false, err
	}

	err = s.fillLinks(ctx, apiKey, &result)
	if err != nil {
		s.logger.Err(err).Msg("cannot fill links")
	}

	return &result.Data[0], true, nil
}

func (s *store) FindManual(ctx context.Context, search string) ([]ManualResponse, error) {
	pipeline := []bson.M{
		{"$match": bson.M{
			"v.meta":     bson.M{"$ne": nil},
			"v.resolved": nil,
		}},
	}
	pipeline = append(pipeline, getMetaAlarmRuleLookup()...)
	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"meta_alarm_rule.type": correlation.RuleManualGroup,
	}})
	if search != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{
			"v.display_name": primitive.Regex{
				Pattern: fmt.Sprintf(".*%s.*", search),
				Options: "i",
			},
		}})
	}
	pipeline = append(pipeline,
		common.GetSortQuery("v.display_name", common.SortAsc),
		bson.M{"$limit": manualMetaAlarmsLimit},
		bson.M{"$project": bson.M{
			"name": "$v.display_name",
		}},
	)

	cursor, err := s.mainDbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var res []ManualResponse
	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, err
	}

	if res == nil {
		res = make([]ManualResponse, 0)
	}

	return res, nil
}

func (s *store) FindByService(ctx context.Context, id, apiKey string, r ListByServiceRequest) (*AggregationResult, error) {
	now := types.NewCpsTime()
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

	ids := service.Depends
	if r.WithService {
		ids = append(ids, id)
	}

	pipeline, err := s.queryBuilder.CreateAggregationPipelineByMatch(ctx, bson.M{
		"d":          bson.M{"$in": ids},
		"v.resolved": nil,
	}, r.Query, r.SortRequest, FilterRequest{BaseFilterRequest: BaseFilterRequest{
		Category: r.Category,
		Search:   r.Search,
	}}, now)
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

	_, err = s.fillAssignedInstructions(ctx, &result, now)
	if err != nil {
		return nil, err
	}
	err = s.fillInstructionFlags(ctx, &result)
	if err != nil {
		return nil, err
	}

	err = s.fillLinks(ctx, apiKey, &result)
	if err != nil {
		s.logger.Err(err).Msg("cannot fill links")
	}

	return &result, nil
}

func (s *store) FindByComponent(ctx context.Context, r ListByComponentRequest, apiKey string) (*AggregationResult, error) {
	now := types.NewCpsTime()
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

	pipeline, err := s.queryBuilder.CreateAggregationPipelineByMatch(ctx, bson.M{
		"d":          bson.M{"$in": component.Depends},
		"v.resolved": nil,
	}, r.Query, r.SortRequest, FilterRequest{}, now)
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

	_, err = s.fillAssignedInstructions(ctx, &result, now)
	if err != nil {
		return nil, err
	}
	err = s.fillInstructionFlags(ctx, &result)
	if err != nil {
		return nil, err
	}

	err = s.fillLinks(ctx, apiKey, &result)
	if err != nil {
		s.logger.Err(err).Msg("cannot fill links")
	}

	return &result, nil
}

func (s *store) FindResolved(ctx context.Context, r ResolvedListRequest, apiKey string) (*AggregationResult, error) {
	now := types.NewCpsTime()

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

	match := bson.M{"d": r.ID}
	opened := false
	pipeline, err := s.queryBuilder.CreateAggregationPipelineByMatch(ctx, match, r.Query, r.SortRequest, FilterRequest{BaseFilterRequest: BaseFilterRequest{
		StartFrom: r.StartFrom,
		StartTo:   r.StartTo,
		Opened:    &opened,
	}}, now)
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

	err = s.fillLinks(ctx, apiKey, &result)
	if err != nil {
		s.logger.Err(err).Msg("cannot fill links")
	}

	return &result, nil
}

func (s *store) GetDetails(ctx context.Context, apiKey string, r DetailsRequest) (*Details, error) {
	now := types.NewCpsTime()
	match := bson.M{"_id": r.ID}
	collection := s.mainDbCollection
	switch r.GetOpenedFilter() {
	case OnlyOpened:
		match["v.resolved"] = nil
	case OnlyResolved:
		collection = s.resolvedDbCollection
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$addFields": bson.M{
			"is_meta_alarm": bson.M{"$cond": bson.A{bson.M{"$not": bson.A{"$v.meta"}}, false, true}},
		}},
	}

	if r.Steps != nil {
		pipeline = append(pipeline, bson.M{"$addFields": bson.M{
			"steps.data":  bson.M{"$slice": bson.A{"$v.steps", (r.Steps.Page - 1) * r.Steps.Limit, r.Steps.Limit}},
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
		details.Steps.Meta, err = common.NewPaginatedMeta(*r.Steps, details.StepsCount)
		if err != nil {
			return nil, err
		}
	}

	if r.Children != nil {
		children := AggregationResult{
			Data: make([]Alarm, 0),
		}

		if details.IsMetaAlarm {
			childrenPipeline, err := s.queryBuilder.CreateChildrenAggregationPipeline(*r.Children, r.GetOpenedFilter(), details.EntityID, now)
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

			if r.WithInstructions {
				_, err = s.fillAssignedInstructions(ctx, &children, now)
				if err != nil {
					return nil, err
				}
				err = s.fillInstructionFlags(ctx, &children)
				if err != nil {
					return nil, err
				}
			}

			err = s.fillLinks(ctx, apiKey, &children)
			if err != nil {
				s.logger.Err(err).Msg("cannot fill links")
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

	return &details, nil
}

func (s *store) Count(ctx context.Context, r FilterRequest) (*Count, error) {
	collection := s.mainDbCollection
	if r.GetOpenedFilter() == OnlyResolved {
		collection = s.resolvedDbCollection
	}

	pipeline, err := s.queryBuilder.CreateCountAggregationPipeline(ctx, r, types.NewCpsTime())
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
	m, _, err := s.getAssignedInstructionsMap(ctx, alarmIds, types.NewCpsTime())
	return m, err
}

func (s *store) getAssignedInstructionsMap(ctx context.Context, alarmIds []string, now types.CpsTime) (map[string][]AssignedInstruction, bson.M, error) {
	instructionCursor, err := s.dbInstructionCollection.Aggregate(
		ctx,
		[]bson.M{
			{"$match": bson.M{
				"type":   bson.M{"$in": bson.A{InstructionTypeManual, InstructionTypeSimplifiedManual}},
				"status": bson.M{"$in": bson.A{InstructionStatusApproved, nil}},
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

	instructionMap := make(map[string]InstructionWithExecutions)
	instructionFiltersPipeline := bson.M{}
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

	if len(instructionMap) == 0 {
		return nil, anyInstructionMatch, nil
	}

	pipeline := []bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": alarmIds}}},
		{"$addFields": bson.M{
			"v.infos_array": bson.M{"$objectToArray": "$v.infos"},
			"v.duration": bson.M{"$subtract": bson.A{
				bson.M{"$cond": bson.M{
					"if":   "$v.resolved",
					"then": "$v.resolved",
					"else": now,
				}},
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
		return nil, nil, err
	}

	defer assignedInstructionsCursor.Close(ctx)

	assignedInstructionsMap := make(map[string][]AssignedInstruction)
	for assignedInstructionsCursor.Next(ctx) {
		assignedInstructions := make(map[string][]string)
		err = assignedInstructionsCursor.Decode(&assignedInstructions)
		if err != nil {
			return nil, nil, err
		}

		for instructionId, alarmIds := range assignedInstructions {
			for _, alarmId := range alarmIds {
				execution := instructionMap[instructionId].GetExecution(alarmId)
				assignedInstructionsMap[alarmId] = append(assignedInstructionsMap[alarmId], AssignedInstruction{
					ID:        instructionId,
					Name:      instructionMap[instructionId].Name,
					Execution: execution,
				})
			}
		}
	}

	return assignedInstructionsMap, anyInstructionMatch, nil
}

func (s *store) GetInstructionExecutionStatuses(ctx context.Context, alarmIDs []string) (map[string]ExecutionStatus, error) {
	if len(alarmIDs) == 0 {
		return nil, nil
	}
	cursor, err := s.dbInstructionExecutionCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"alarm": bson.M{"$in": alarmIDs},
		}},
		{"$lookup": bson.M{
			"from":         mongo.InstructionMongoCollection,
			"localField":   "instruction",
			"foreignField": "_id",
			"as":           "instruction",
		}},
		{"$unwind": "$instruction"},
		{"$group": bson.M{
			"_id": "$alarm",
			"auto_statuses": bson.M{"$addToSet": bson.M{"$cond": bson.M{
				"if":   bson.M{"$eq": bson.A{"$instruction.type", InstructionTypeAuto}},
				"then": "$status",
				"else": "$$REMOVE",
			}}},
			"manual_statuses": bson.M{"$addToSet": bson.M{"$cond": bson.M{
				"if":   bson.M{"$eq": bson.A{"$instruction.type", InstructionTypeManual}},
				"then": "$status",
				"else": "$$REMOVE",
			}}},
		}},
		{"$addFields": bson.M{
			"auto_running": bson.M{"$gt": bson.A{
				bson.M{"$size": bson.M{"$filter": bson.M{
					"input": "$auto_statuses",
					"cond":  bson.M{"$in": bson.A{"$$this", []int{InstructionExecutionStatusRunning, InstructionExecutionStatusWaitResult}}},
				}}},
				0,
			}},
			"auto_failed": bson.M{"$gt": bson.A{
				bson.M{"$size": bson.M{"$filter": bson.M{
					"input": "$auto_statuses",
					"cond":  bson.M{"$in": bson.A{"$$this", []int{InstructionExecutionStatusAborted, InstructionExecutionStatusFailed}}},
				}}},
				0,
			}},
			"manual_running": bson.M{"$gt": bson.A{
				bson.M{"$size": bson.M{"$filter": bson.M{
					"input": "$manual_statuses",
					"cond":  bson.M{"$eq": bson.A{"$$this", InstructionExecutionStatusRunning}},
				}}},
				0,
			}},
			"manual_waiting_result": bson.M{"$gt": bson.A{
				bson.M{"$size": bson.M{"$filter": bson.M{
					"input": "$manual_statuses",
					"cond":  bson.M{"$eq": bson.A{"$$this", InstructionExecutionStatusWaitResult}},
				}}},
				0,
			}},
		}},
		{"$addFields": bson.M{
			"auto_all_completed": bson.M{"$and": bson.A{
				bson.M{"$not": "$auto_running"},
				bson.M{"$gt": bson.A{bson.M{"$size": "$auto_statuses"}, 0}},
			}},
		}},
	})
	if err != nil {
		return nil, err
	}

	executionStatuses := make([]ExecutionStatus, 0)
	err = cursor.All(ctx, &executionStatuses)
	if err != nil {
		return nil, err
	}
	statusesByAlarm := make(map[string]ExecutionStatus, len(executionStatuses))
	for _, v := range executionStatuses {
		statusesByAlarm[v.ID] = v
	}

	return statusesByAlarm, nil
}

func (s *store) fillAssignedInstructions(ctx context.Context, result *AggregationResult, now types.CpsTime) (bson.M, error) {
	var alarmIds []string
	for _, item := range result.Data {
		if item.Value.Resolved == nil {
			alarmIds = append(alarmIds, item.ID)
		}
	}

	if len(alarmIds) == 0 {
		return nil, nil
	}

	assignedInstructionsMap, anyInstructionMatch, err := s.getAssignedInstructionsMap(ctx, alarmIds, now)
	if err != nil {
		return nil, err
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

	return anyInstructionMatch, nil
}

func (s *store) fillInstructionFlags(ctx context.Context, result *AggregationResult) error {
	alarmIDs := make([]string, len(result.Data))
	for i, item := range result.Data {
		if item.Value.Resolved == nil {
			alarmIDs[i] = item.ID
		}
	}
	if len(alarmIDs) == 0 {
		return nil
	}

	statusesByAlarm, err := s.GetInstructionExecutionStatuses(ctx, alarmIDs)
	if err != nil {
		return err
	}

	for i, v := range result.Data {
		result.Data[i].IsAutoInstructionRunning = statusesByAlarm[v.ID].AutoRunning
		result.Data[i].IsAllAutoInstructionsCompleted = statusesByAlarm[v.ID].AutoAllCompleted
		result.Data[i].IsAutoInstructionFailed = statusesByAlarm[v.ID].AutoFailed
		result.Data[i].IsManualInstructionRunning = statusesByAlarm[v.ID].ManualRunning
		result.Data[i].IsManualInstructionWaitingResult = statusesByAlarm[v.ID].ManualWaitingResult
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
func (s *store) fillLinks(ctx context.Context, apiKey string, result *AggregationResult) error {
	if result == nil || len(result.Data) == 0 {
		return nil
	}

	reqEntities := make([]common.FetchLinksRequestItem, len(result.Data))
	alarmIndexes := make(map[string]int, len(result.Data))

	for i, item := range result.Data {
		reqEntities[i] = common.FetchLinksRequestItem{
			AlarmID:  item.ID,
			EntityID: item.Entity.ID,
		}
		alarmIndexes[item.ID] = i
	}

	res, err := s.linksFetcher.Fetch(ctx, apiKey, common.FetchLinksRequest{Entities: reqEntities})
	if err != nil || res == nil {
		return err
	}

	for _, rec := range res.Data {
		if i, ok := alarmIndexes[rec.AlarmID]; ok && len(rec.Links) > 0 {
			links := make(map[string]interface{}, len(rec.Links))
			for category, link := range rec.Links {
				links[category] = link
			}

			result.Data[i].Links = links
		}
	}

	return nil
}
