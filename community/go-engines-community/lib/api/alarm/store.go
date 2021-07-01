package alarm

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/expression/parser"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const InstructionStatusRunning = 0
const InstructionStatusPaused = 1

const InstructionsQueryAll = "all"

const linkFetchTimeout = 30 * time.Second

type Store interface {
	Find(context.Context, string, ListRequestWithPagination) (*AggregationResult, error)
	Count(context.Context, FilterRequest) (*Count, error)
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	fieldsAliases         map[string]string
	fieldsAliasesByRegex  map[string]string
	defaultSearchByFields []string
	defaultSortBy         string
	defaultSort           string
	links                 LinksFetcher
	// nested objects lookups depend on requested Filter:
	// these aggregation stages inserted at beginning of pipeline when Filter has some expression
	// or inserted at beginning of project stage otherwise
	deferredNestedObjects []bson.M
}

func NewStore(dbClient mongo.DbClient, legacyURL fmt.Stringer) Store {
	s := &store{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(alarm.AlarmCollectionName),
		fieldsAliases: map[string]string{
			"uid":            "_id",
			"connector":      "v.connector",
			"connector_name": "v.connector_name",
			"component":      "v.component",
			"resource":       "v.resource",
			"entity_id":      "d",
			"state":          "v.state.val",
			"status":         "v.status.val",
			"snooze":         "v.snooze",
			"ack":            "v.ack",
			"cancel":         "v.cancel",
			"ticket":         "v.ticket.val",
			"output":         "v.output",
			"opened":         "t",
			"resolved":       "v.resolved",
			"domain":         "v.extra.domain",
			"perimeter":      "v.extra.perimeter",
		},
		fieldsAliasesByRegex: map[string]string{
			"^infos\\.(.+)":           "entity.infos.$1",
			"^v\\.infos\\.\\*\\.(.+)": "v_infos_array.v.$1",
		},
		defaultSearchByFields: []string{
			"v.connector",
			"v.connector_name",
			"v.component",
			"v.resource",
		},
		defaultSortBy: "t",
		defaultSort:   common.SortDesc,
	}
	s.links = NewLinksFetcher(legacyURL, linkFetchTimeout)
	return s
}

// move $sort to the end of $project when sorting by added field v.duration
func (s *store) adjustSort(sort bson.M, project []bson.M) []bson.M {
	srt := sort["$sort"]
	if sortExpr, ok := srt.(bson.D); ok && len(sortExpr) > 1 && sortExpr[0].Key == "v.duration" {
		project = append(project, bson.M{"$sort": sortExpr})
		delete(sort, "$sort")
	} else if sortExpr, ok := srt.(bson.M); ok && sortExpr != nil && sortExpr["v.duration"] != nil {
		project = append(project, bson.M{"$sort": bson.M{"v.duration": sortExpr["v.duration"]}})
		delete(sort, "$sort")
	}
	return project
}

// insertDeferred iserts deferredNestedObjects at the beginning of pipeline or project stages depending on Filter
func (s *store) insertDeferred(r FilterRequest, pipeline *[]bson.M, project []bson.M) []bson.M {
	doLen := len(s.deferredNestedObjects)
	// deferred nested objects currently have only pbehavior collection lookup
	if strings.Contains(r.Filter, "pbehavior.") || strings.Contains(r.Filter, `"pbehavior"`) {
		p := make([]bson.M, doLen+len(*pipeline))
		copy(p, s.deferredNestedObjects)
		copy(p[doLen:], *pipeline)
		*pipeline = p
		return project
	}
	p := make([]bson.M, doLen+len(project))
	copy(p, s.deferredNestedObjects)
	copy(p[doLen:], project)
	return p
}

func (s *store) Find(ctx context.Context, apiKey string, r ListRequestWithPagination) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	s.addNestedObjects(r.FilterRequest, &pipeline)
	err := s.addFilter(ctx, r.FilterRequest, &pipeline)
	if err != nil {
		return nil, err
	}
	sort := s.getSort(r.ListRequest)
	entitiesToProject := s.resetEntities(r.ListRequest, &pipeline)
	project := s.getProject(r.ListRequest, entitiesToProject)
	project = s.insertDeferred(r.FilterRequest, &pipeline, project)
	project = s.adjustSort(sort, project)
	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		sort,
		project,
	), options.Aggregate().SetAllowDiskUse(true))

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

	err = s.fillChildren(ctx, r.ListRequest, &result)
	if err != nil {
		return nil, err
	}

	err = s.fillAssignedInstructions(ctx, &result)
	if err != nil {
		return nil, err
	}

	err = s.fillLinks(ctx, apiKey, &result)
	if err != nil {
		log.Printf("fillLinks error %s", err)
	}

	return &result, nil
}

func (s *store) Count(ctx context.Context, r FilterRequest) (*Count, error) {
	pipeline := make([]bson.M, 0)
	s.addNestedObjects(r, &pipeline)
	err := s.addFilter(ctx, r, &pipeline)
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

	cursor, err := s.dbCollection.Aggregate(ctx, aggregationPipeline)
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

func (s *store) fillChildren(ctx context.Context, r ListRequest, result *AggregationResult) error {
	childrenIds := make([]string, 0)
	for i := range result.Data {
		if result.Data[i].ChildrenIDs != nil {
			if len(result.Data[i].ChildrenIDs.Data) == 0 {
				result.Data[i].Children = &Children{
					Data:  nil,
					Total: result.Data[i].ChildrenIDs.Total,
				}
			} else {
				childrenIds = append(childrenIds, result.Data[i].ChildrenIDs.Data...)
			}
		}
	}

	if len(childrenIds) == 0 {
		return nil
	}

	pipeline := make([]bson.M, 0)
	pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": []bson.M{
		{"d": bson.M{"$in": childrenIds}},
		{"v.resolved": bson.M{"$exists": false}},
	}}})
	s.addNestedObjects(r.FilterRequest, &pipeline)
	pipeline = append(pipeline, s.getSort(r))
	pipeline = append(pipeline, s.getProject(r, false)...)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	var children []Alarm
	err = cursor.All(ctx, &children)
	if err != nil {
		return err
	}

	childrenByEntityID := make(map[string][]Alarm)
	for _, ch := range children {
		if _, ok := childrenByEntityID[ch.Entity.ID]; !ok {
			childrenByEntityID[ch.Entity.ID] = make([]Alarm, 0)
		}

		childrenByEntityID[ch.Entity.ID] = append(childrenByEntityID[ch.Entity.ID], ch)
	}

	for i := range result.Data {
		if result.Data[i].ChildrenIDs == nil {
			continue
		}
		for _, chID := range result.Data[i].ChildrenIDs.Data {
			if children, ok := childrenByEntityID[chID]; ok {
				if result.Data[i].Children == nil {
					result.Data[i].Children = &Children{
						Data:  make([]Alarm, 0),
						Total: result.Data[i].ChildrenIDs.Total,
					}
				}

				result.Data[i].Children.Data = append(result.Data[i].Children.Data, children...)
			}
		}
	}

	return nil
}

func (s *store) getAssignedInstructionsMap(ctx context.Context, alarmIds []string) (map[string][]InstructionWithAlarms, error) {
	instructionCursor, err := s.dbClient.Collection(mongo.InstructionMongoCollection).Aggregate(
		ctx,
		[]bson.M{
			{"$lookup": bson.M{
				"from": mongo.InstructionExecutionMongoCollection,
				"let":  bson.M{"instruction_ids": "$_id"},
				"pipeline": []bson.M{
					{"$match": bson.M{
						"$expr": bson.M{
							"$and": []bson.M{
								{"$eq": bson.A{"$instruction", "$$instruction_ids"}},
								{"$in": bson.A{"$status", []int{InstructionStatusRunning, InstructionStatusPaused}}},
								{"$in": bson.A{"$alarm", alarmIds}},
							},
						},
					}},
				},
				"as": "instruction_executions",
			}},
			{"$addFields": bson.M{
				"alarms_with_executions": bson.M{
					"$map": bson.M{
						"input": "$instruction_executions",
						"as":    "execution",
						"in": bson.M{
							"_id":    "$$execution._id",
							"alarm":  "$$execution.alarm",
							"status": "$$execution.status",
						},
					},
				},
			}},
			{"$project": bson.M{
				"_id":                    1,
				"name":                   1,
				"alarm_patterns":         1,
				"entity_patterns":        1,
				"alarms_with_executions": 1,
				"active_on_pbh":          1,
				"disabled_on_pbh":        1,
			}},
		},
	)
	if err != nil {
		return nil, err
	}

	defer instructionCursor.Close(ctx)

	instructionMap := make(map[string]InstructionWithAlarms)
	instructionFiltersPipeline := bson.M{}

	for instructionCursor.Next(ctx) {
		var instructionDocument InstructionWithAlarms
		err = instructionCursor.Decode(&instructionDocument)
		if err != nil {
			return nil, err
		}

		if instructionDocument.AlarmPatterns.IsSet() || instructionDocument.EntityPatterns.IsSet() {
			instructionMap[instructionDocument.ID] = instructionDocument

			and := []bson.M{
				instructionDocument.AlarmPatterns.AsMongoDriverQuery(),
				getEntityPatternsForEntity(instructionDocument.EntityPatterns.AsMongoDriverQuery()),
			}

			if len(instructionDocument.ActiveOnPbh) > 0 {
				and = append(and, bson.M{"v.pbehavior_info.type": bson.M{"$in": instructionDocument.ActiveOnPbh}})
			}

			if len(instructionDocument.DisabledOnPbh) > 0 {
				and = append(and, bson.M{"v.pbehavior_info.type": bson.M{"$nin": instructionDocument.DisabledOnPbh}})
			}

			instructionFiltersPipeline[instructionDocument.ID] = []bson.M{
				{
					"$match": bson.M{
						"$and": and,
					},
				},
			}
		}
	}

	if len(instructionMap) == 0 {
		return nil, nil
	}

	pipeline := []bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": alarmIds}}},
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

	assignedInstructionsCursor, err := s.dbClient.Collection(alarm.AlarmCollectionName).Aggregate(
		ctx,
		pipeline,
	)
	if err != nil {
		return nil, err
	}

	defer assignedInstructionsCursor.Close(ctx)

	assignedInstructionsMap := make(map[string][]InstructionWithAlarms)
	for assignedInstructionsCursor.Next(ctx) {
		assignedInstructions := make(map[string][]string)
		err = assignedInstructionsCursor.Decode(&assignedInstructions)
		if err != nil {
			return nil, err
		}

		for instructionId, alarmIds := range assignedInstructions {
			for _, alarmId := range alarmIds {
				execution := instructionMap[instructionId].GetExecution(alarmId)
				assignedInstructionsMap[alarmId] = append(assignedInstructionsMap[alarmId], InstructionWithAlarms{
					ID:            instructionId,
					Name:          instructionMap[instructionId].Name,
					Execution:     execution,
					ActiveOnPbh:   instructionMap[instructionId].ActiveOnPbh,
					DisabledOnPbh: instructionMap[instructionId].DisabledOnPbh,
				})
			}
		}
	}

	return assignedInstructionsMap, nil
}

func (s *store) fillAssignedInstructions(ctx context.Context, result *AggregationResult) error {
	var alarmIds []string
	for _, alarmDocument := range result.Data {
		alarmIds = append(alarmIds, alarmDocument.ID)
	}

	if len(alarmIds) == 0 {
		return nil
	}

	assignedInstructionsMap, err := s.getAssignedInstructionsMap(ctx, alarmIds)
	if err != nil {
		return err
	}

	for i, alarmDocument := range result.Data {
		sort.Slice(assignedInstructionsMap[alarmDocument.ID], func(i, j int) bool {
			return assignedInstructionsMap[alarmDocument.ID][i].Name < assignedInstructionsMap[alarmDocument.ID][j].Name
		})

		result.Data[i].AssignedInstructions = assignedInstructionsMap[alarmDocument.ID]
	}

	return nil
}

type LinksRequest struct {
	Entities []AlarmEntity `json:"entities"`
}

type AlarmEntity struct {
	AlarmID  string `json:"alarm"`
	EntityID string `json:"entity"`
}

type EntityLinks struct {
	AlarmEntity
	Links map[string]interface{} `json:"links"`
}

type LinksResponse struct {
	Data []EntityLinks
}

// request alarm links from API v2 and fill AggregationResult with links
func (s *store) fillLinks(ctx context.Context, apiKey string, result *AggregationResult) error {
	if result == nil || len(result.Data) == 0 {
		return nil
	}

	maxItems := len(result.Data)
	if maxItems > 100 {
		maxItems = 100
	}
	linksEntities := make([]AlarmEntity, 0, maxItems)
	alarms := make(map[string]int, maxItems)
	for i, al := range result.Data {
		linksEntities = append(linksEntities, AlarmEntity{
			AlarmID:  al.ID,
			EntityID: al.Entity.ID,
		})
		alarms[al.ID] = i
		if len(linksEntities) == maxItems {
			break
		}
	}

	res, err := s.links.Fetch(ctx, apiKey, linksEntities)
	if err != nil {
		return err
	}

	for _, rec := range res.Data {
		if i, ok := alarms[rec.AlarmID]; ok {
			result.Data[i].Links = make(map[string]interface{}, len(rec.Links))
			for category, link := range rec.Links {
				result.Data[i].Links[category] = link
			}
		}
	}

	return nil
}

func (s *store) addFilter(ctx context.Context, r FilterRequest, pipeline *[]bson.M) error {
	match := make([]bson.M, 0)
	s.addStartFromFilter(r, &match)
	s.addStartToFilter(r, &match)
	s.addOnlyOpenedFilter(r, &match)
	s.addOnlyResolvedFilter(r, &match)
	replacedKeys, err := s.addQueryFilter(r, &match)
	if err != nil {
		return err
	}
	s.addOnlyParentsFilter(r, &match)
	s.addCategoryFilter(r, &match)

	err = s.addOnlyInstructionsFilter(ctx, r, &match)
	if err != nil {
		return err
	}

	searchReplacedKeys := s.addSearchFilter(r, pipeline, &match)
	replacedKeys = append(replacedKeys, searchReplacedKeys...)
	err = s.addOnlyManualFilter(r, &match)
	if err != nil {
		return err
	}

	if len(match) > 0 {
		// Add auxiliary field to implement filtering by "v.infos" field.
		addField := false
		for _, key := range replacedKeys {
			if strings.HasPrefix(key, "v_infos_array") {
				addField = true
				break
			}
		}
		if addField {
			*pipeline = append(*pipeline, bson.M{"$addFields": bson.M{
				"v_infos_array": bson.M{"$objectToArray": "$v.infos"},
			}})
		}
		*pipeline = append(*pipeline, bson.M{"$match": bson.M{"$and": match}})
	}

	return nil
}

func (s *store) addStartFromFilter(r FilterRequest, match *[]bson.M) {
	if r.StartFrom == nil {
		return
	}

	*match = append(*match, bson.M{"t": bson.M{"$gte": r.StartFrom}})
}

func (s *store) addStartToFilter(r FilterRequest, match *[]bson.M) {
	if r.StartTo == nil {
		return
	}

	*match = append(*match, bson.M{"t": bson.M{"$lte": r.StartTo}})
}

func (s *store) addOnlyOpenedFilter(r FilterRequest, match *[]bson.M) {
	if !r.OnlyOpened || r.OnlyResolved {
		return
	}

	*match = append(*match, bson.M{"v.resolved": bson.M{"$exists": false}})
}

func (s *store) addOnlyResolvedFilter(r FilterRequest, match *[]bson.M) {
	if !r.OnlyResolved || r.OnlyOpened {
		return
	}

	*match = append(*match, bson.M{"v.resolved": bson.M{"$exists": true}})
}

func (s *store) addQueryFilter(r FilterRequest, match *[]bson.M) ([]string, error) {
	if r.Filter == "" {
		return nil, nil
	}

	var queryFilter bson.M
	err := json.Unmarshal([]byte(r.Filter), &queryFilter)
	if err != nil {
		return nil, err
	}

	resolvedFilter, replacedKeys := s.resolveAliasesInQuery(queryFilter)

	*match = append(*match, resolvedFilter.(bson.M))
	return replacedKeys, nil
}

func (s *store) addOnlyParentsFilter(r FilterRequest, match *[]bson.M) {
	if r.OnlyParents {
		*match = append(*match, bson.M{"$or": []bson.M{
			{"v.parents": bson.M{"$exists": false}},
			{"v.parents": bson.M{"$eq": []string{}}},
			{"v.meta": bson.M{"$exists": true}},
		}})
	} else {
		*match = append(*match, bson.M{"$or": []bson.M{
			{"v.meta": bson.M{"$exists": false}},
			{"v.meta": bson.M{"$eq": ""}},
		}})
	}
}

func (s *store) addOnlyInstructionsFilter(ctx context.Context, r FilterRequest, match *[]bson.M) error {
	var instructionFilters []bson.M
	var err error

	if r.WithoutInstructions != "" {
		if r.WithoutInstructions == InstructionsQueryAll {
			instructionFilters, err = s.getInstructionsFilters(ctx, bson.M{})
		} else {
			instructionFilters, err = s.getInstructionsFilters(ctx, bson.M{"name": bson.M{"$in": strings.Split(r.WithoutInstructions, ",")}})
		}

		if instructionFilters != nil {
			*match = append(*match, bson.M{"$nor": instructionFilters})
		}
	}

	if r.WithInstructions != "" {
		if r.WithInstructions == InstructionsQueryAll {
			instructionFilters, err = s.getInstructionsFilters(ctx, bson.M{})
		} else {
			instructionFilters, err = s.getInstructionsFilters(ctx, bson.M{"name": bson.M{"$in": strings.Split(r.WithInstructions, ",")}})
		}

		if instructionFilters != nil {
			*match = append(*match, bson.M{"$or": instructionFilters})
		} else {
			*match = append(*match, bson.M{"$nor": []bson.M{{}}})
		}
	}

	return err
}

func (s *store) getInstructionsFilters(ctx context.Context, filter bson.M) ([]bson.M, error) {
	instructionCursor, _ := s.dbClient.Collection(mongo.InstructionMongoCollection).Find(ctx, filter)
	defer instructionCursor.Close(ctx)

	var instructionFilters []bson.M

	for instructionCursor.Next(ctx) {
		var instructionDocument InstructionWithAlarms
		err := instructionCursor.Decode(&instructionDocument)
		if err != nil {
			return nil, err
		}

		and := []bson.M{
			instructionDocument.AlarmPatterns.AsMongoDriverQuery(),
			getEntityPatternsForEntity(instructionDocument.EntityPatterns.AsMongoDriverQuery()),
		}

		if len(instructionDocument.ActiveOnPbh) > 0 {
			and = append(and, bson.M{"v.pbehavior_info.type": bson.M{"$in": instructionDocument.ActiveOnPbh}})
		}

		if len(instructionDocument.DisabledOnPbh) > 0 {
			and = append(and, bson.M{"v.pbehavior_info.type": bson.M{"$nin": instructionDocument.DisabledOnPbh}})
		}

		if instructionDocument.AlarmPatterns.IsSet() || instructionDocument.EntityPatterns.IsSet() {
			instructionFilters = append(instructionFilters, bson.M{"$and": and})
		}
	}

	return instructionFilters, nil
}

func getEntityPatternsForEntity(patternBson bson.M) bson.M {
	newBson := make(bson.M)
	patternListInterface, ok := patternBson["$or"]
	if ok {
		patternList := patternListInterface.([]bson.M)
		newPatternsList := make([]bson.M, len(patternList))
		for i, pattern := range patternList {
			newPattern := make(bson.M)
			for k, vv := range pattern {
				newPattern["entity."+k] = vv
			}

			newPatternsList[i] = newPattern
		}

		// Just in case when an entity pattern's function AsMongoDriverQuery returns an empty bson array
		// that might happen if the pattern has bad format in mongo, but after unmarshalling it has isSet = true
		// since it's not possible for $or has empty array we just fill it with an empty bson
		// @todo: It's a temporary solution to avoid panic, the issue should be solved in UnmarshalBSONValue functions for patterns
		if len(newPatternsList) == 0 {
			newPatternsList = append(newPatternsList, bson.M{})
		}

		newBson["$or"] = newPatternsList
	} else {
		return patternBson
	}

	return newBson
}

func (s *store) addSearchFilter(r FilterRequest, pipeline *[]bson.M,
	match *[]bson.M) []string {
	if r.Search == "" {
		return nil
	}

	p := parser.NewParser()
	expr, err := p.Parse(r.Search)
	if err == nil {
		query := expr.Query()
		resolvedQuery, replacedKeys := s.resolveAliasesInQuery(query)
		*match = append(*match, resolvedQuery.(bson.M))

		return replacedKeys
	}

	searchRegexp := primitive.Regex{
		Pattern: fmt.Sprintf(".*%s.*", r.Search),
		Options: "i",
	}

	fields, replacedKeys := s.resolveAliases(r.SearchBy)
	if len(fields) == 0 {
		fields = s.defaultSearchByFields
	}

	searchMatch := make([]bson.M, len(fields))
	for i := range fields {
		searchMatch[i] = bson.M{fields[i]: searchRegexp}
	}

	if r.OnlyParents {
		*pipeline = append(*pipeline, bson.M{"$graphLookup": bson.M{
			"from":                    alarm.AlarmCollectionName,
			"startWith":               "$v.children",
			"connectFromField":        "v.children",
			"connectToField":          "d",
			"restrictSearchWithMatch": bson.M{"$or": searchMatch},
			"as":                      "filtered_children",
		}})
		*match = append(*match, bson.M{
			"$or": append(searchMatch, bson.M{"filtered_children": bson.M{"$ne": []string{}}}),
		})
	} else {
		*match = append(*match, bson.M{
			"$or": searchMatch,
		})
	}

	return replacedKeys
}

func (s *store) addOnlyManualFilter(r FilterRequest, match *[]bson.M) error {
	if r.OnlyManual {
		*match = append(*match, bson.M{"$expr": bson.M{"$eq": bson.A{"$meta_alarm_rule.type", metaalarm.RuleManualGroup}}})
	}

	return nil
}

func (s *store) addCategoryFilter(r FilterRequest, match *[]bson.M) {
	if r.Category == "" {
		return
	}

	*match = append(*match, bson.M{"entity.category._id": bson.M{"$eq": r.Category}})
}

func (s *store) addNestedObjects(r FilterRequest, pipeline *[]bson.M) {
	*pipeline = append(*pipeline,
		bson.M{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		bson.M{"$match": bson.M{"$or": []bson.M{{"entity.enabled": true}, {"entity": bson.M{"$exists": false}}}}},
		bson.M{"$unwind": bson.M{"path": "$entity", "preserveNullAndEmptyArrays": true}},
		bson.M{"$addFields": bson.M{
			"impact_state": bson.M{"$multiply": bson.A{"$v.state.val", "$entity.impact_level"}},
		}},
		bson.M{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "entity.category",
			"foreignField": "_id",
			"as":           "entity.category",
		}},
		bson.M{"$unwind": bson.M{"path": "$entity.category", "preserveNullAndEmptyArrays": true}},
	)
	s.deferredNestedObjects = []bson.M{
		{"$lookup": bson.M{
			"from":         pbehavior.PBehaviorCollectionName,
			"foreignField": "_id",
			"localField":   "v.pbehavior_info.id",
			"as":           "pbehavior",
		}},
		{"$unwind": bson.M{"path": "$pbehavior", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         pbehavior.TypeCollectionName,
			"foreignField": "_id",
			"localField":   "pbehavior.type_",
			"as":           "pbehavior.type",
		}},
		{"$unwind": bson.M{"path": "$pbehavior.type", "preserveNullAndEmptyArrays": true}},
	}
	if r.OnlyParents {
		*pipeline = append(*pipeline,
			bson.M{"$lookup": bson.M{
				"from":         metaalarm.RulesCollectionName,
				"localField":   "v.meta",
				"foreignField": "_id",
				"as":           "meta_alarm_rule",
			}},
			bson.M{"$unwind": bson.M{"path": "$meta_alarm_rule", "preserveNullAndEmptyArrays": true}},
		)
	}
}

func (s *store) resetEntities(r ListRequest, pipeline *[]bson.M) bool {
	if strings.HasPrefix(r.SortBy, "entity.") || strings.HasPrefix(r.SortBy, "infos.") {
		return false
	}
	*pipeline = append(*pipeline, bson.M{"$project": bson.M{"entity": 0}})
	return true
}

func (s *store) getProject(r ListRequest, entitiesToProject bool) []bson.M {
	now := time.Now().Unix()
	addFields := bson.M{
		"infos": "$v.infos",
		// outer field lastComment to make use it in $project
		"lastComment": bson.M{
			"$reduce": bson.M{
				"input": bson.M{
					"$slice": bson.A{
						bson.M{"$filter": bson.M{
							"input": bson.M{"$reverseArray": "$v.steps"},
							"as":    "steps",
							"cond": bson.M{
								"$eq": bson.A{"$$steps._t", "comment"},
							},
						}},
						1,
					},
				},
				"initialValue": bson.M{},
				"in":           bson.M{"$mergeObjects": bson.A{bson.M{}, "$$this"}},
			},
		},
		"v.duration": bson.M{"$subtract": bson.A{now, bson.M{"$cond": bson.M{
			"if":   "$v.activation_date",
			"then": "$v.activation_date",
			"else": "$v.creation_date",
		}}}},
		"v.current_state_duration": bson.M{"$subtract": bson.A{now, "$v.state.t"}},
	}

	project := bson.M{
		"filtered_children": 0,
		"v_infos_array":     0,
	}

	if !r.WithSteps {
		project["v.steps"] = 0
	}

	if r.OnlyParents {
		childrenPipeline := bson.M{"total": bson.M{"$size": "$v.children"}}
		if r.WithChildren {
			childrenPipeline["data"] = "$v.children"
		}

		addFields["is_meta_alarm"] = bson.M{"$cond": bson.A{bson.M{"$not": bson.A{"$v.meta"}}, false, true}}
		addFields["children_ids"] = bson.M{"$cond": bson.A{
			bson.M{"$not": bson.A{"$v.meta"}},
			nil,
			childrenPipeline,
		}}
		addFields["filtered_children_ids"] = "$filtered_children._id"
	}
	lastCommentNilWhenEmpty := bson.M{
		"$project": bson.M{
			"t": 1, "d": 1, "v": 1, "entity": 1, "infos": 1, "is_meta_alarm": 1,
			"children_ids": 1, "filtered_children_ids": 1,
			"pbehavior":       bson.M{"$cond": bson.M{"if": bson.M{"$eq": bson.A{bson.M{}, "$pbehavior"}}, "then": "$$REMOVE", "else": "$pbehavior"}},
			"meta_alarm_rule": 1, "causes": 1,
			"impact_state": 1,
			"lastComment": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$eq": bson.A{bson.M{}, "$lastComment"}},
					"then": nil,
					"else": "$lastComment",
				},
			},
		},
	}
	// remove outer lastComment
	project["lastComment"] = 0
	var pipeline []bson.M
	if r.OnlyParents {
		pipeline = s.getCausesPipeline()
	} else {
		pipeline = make([]bson.M, 0, 6)
	}
	if entitiesToProject {
		pipeline = append(pipeline,
			bson.M{"$lookup": bson.M{
				"from":         mongo.EntityMongoCollection,
				"foreignField": "_id",
				"localField":   "d",
				"as":           "entity",
			}},
			bson.M{"$unwind": bson.M{"path": "$entity", "preserveNullAndEmptyArrays": true}},
			bson.M{"$lookup": bson.M{
				"from":         mongo.EntityCategoryMongoCollection,
				"localField":   "entity.category",
				"foreignField": "_id",
				"as":           "entity.category",
			}},
			bson.M{"$unwind": bson.M{"path": "$entity.category", "preserveNullAndEmptyArrays": true}},
		)
	}
	pipeline = append(pipeline, []bson.M{
		{"$addFields": addFields},
		lastCommentNilWhenEmpty,
		{"$addFields": bson.M{"v.last_comment": "$lastComment"}},
		{"$project": project},
	}...)

	return pipeline
}

func (s *store) getSort(r ListRequest) bson.M {
	sortBy := s.resolveAlias(r.SortBy)
	sort := r.Sort

	if strings.HasSuffix(sortBy, ".") {
		sortBy = "v.last_update_date" // moved from reader.py
	}

	if sortBy == "" {
		sortBy = s.defaultSortBy
		if r.Sort == "" {
			sort = s.defaultSort
		}
	}

	return common.GetSortQuery(sortBy, sort)
}

func (s *store) getCausesPipeline() []bson.M {
	return []bson.M{
		{"$graphLookup": bson.M{
			"from":             alarm.AlarmCollectionName,
			"startWith":        "$v.parents",
			"connectFromField": "v.parents",
			"connectToField":   "d",
			"as":               "parents",
		}},
		{"$lookup": bson.M{
			"from":         metaalarm.RulesCollectionName,
			"localField":   "parents.v.meta",
			"foreignField": "_id",
			"as":           "causes_rules",
		}},
		{"$addFields": bson.M{
			"causes": bson.M{
				"$cond": []bson.M{
					{"$and": []bson.M{
						{"$ne": bson.A{"$v.parents", bson.A{}}},
						{"$ifNull": bson.A{"$v.parents", false}},
						{"v": bson.M{"parents": bson.M{"$type": "array"}}},
					}},
					{
						"rules": "$causes_rules",
						"total": bson.M{"$size": "$v.parents"},
					},
					nil,
				},
			},
		}},
	}
}

func (s *store) resolveAliases(v []string) (newV []string, replacedKeys []string) {
	res := make([]string, len(v))
	keys := make([]string, 0)

	for i, alias := range v {
		res[i] = s.resolveAlias(alias)
		if res[i] != alias {
			keys = append(keys, res[i])
		}
	}

	return res, keys
}

func (s *store) resolveAliasesInQuery(query interface{}) (newQuery interface{}, replacedKeys []string) {
	res := query
	val := reflect.ValueOf(res)
	keys := make([]string, 0)

	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			newVal, replaced := s.resolveAliasesInQuery(val.Index(i).Interface())
			keys = append(keys, replaced...)
			val.Index(i).Set(reflect.ValueOf(newVal))
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			newVal, replaced := s.resolveAliasesInQuery(val.MapIndex(key).Interface())
			keys = append(keys, replaced...)
			newKey := s.resolveAlias(key.String())
			if newKey != key.String() {
				keys = append(keys, newKey)
			}
			val.SetMapIndex(key, reflect.Value{})
			val.SetMapIndex(reflect.ValueOf(newKey), reflect.ValueOf(newVal))
		}
	}

	return res, keys
}

func (s *store) resolveAlias(v string) string {
	if v == "" {
		return v
	}

	for alias, field := range s.fieldsAliases {
		if alias == v {
			return field
		}
	}

	for expr, repl := range s.fieldsAliasesByRegex {
		r, err := regexp.Compile(expr)
		if err == nil {
			replace := r.ReplaceAllString(v, repl)
			if v != replace {
				return replace
			}
		}
	}

	return v
}
