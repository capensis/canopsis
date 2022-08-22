package serviceweather

import (
	"context"
	"errors"
	"sort"
	"time"

	alarmapi "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libtypes "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Find(context.Context, ListRequest) (*AggregationResult, error)
	FindEntities(ctx context.Context, id, apiKey string, query EntitiesListRequest) (*EntityAggregationResult, error)
}

func NewStore(
	dbClient mongo.DbClient,
	linksFetcher common.LinksFetcher,
	alarmStore alarmapi.Store,
	timezoneConfigProvider config.TimezoneConfigProvider,
	logger zerolog.Logger,
) Store {
	return &store{
		dbCollection: dbClient.Collection(mongo.EntityMongoCollection),

		alarmStore:   alarmStore,
		linksFetcher: linksFetcher,

		timezoneConfigProvider: timezoneConfigProvider,

		queryBuilder: NewMongoQueryBuilder(dbClient),

		logger: logger,
	}
}

type store struct {
	dbCollection mongo.DbCollection

	linksFetcher common.LinksFetcher
	alarmStore   alarmapi.Store

	timezoneConfigProvider config.TimezoneConfigProvider

	queryBuilder *MongoQueryBuilder

	logger zerolog.Logger
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline, err := s.queryBuilder.CreateListAggregationPipeline(ctx, r)
	if err != nil {
		return nil, err
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var res AggregationResult
	if cursor.Next(ctx) {
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (s *store) FindEntities(ctx context.Context, id, apiKey string, r EntitiesListRequest) (*EntityAggregationResult, error) {
	var service libtypes.Entity
	err := s.dbCollection.FindOne(ctx, bson.M{
		"_id":     id,
		"type":    libtypes.EntityTypeService,
		"enabled": true,
	}).Decode(&service)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	location := s.timezoneConfigProvider.Get().Location
	now := libtypes.CpsTime{Time: time.Now().In(location)}
	pipeline, err := s.queryBuilder.CreateListDependenciesAggregationPipeline(service.Depends, r, now)
	if err != nil {
		return nil, err
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var res EntityAggregationResult
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	if !service.PbehaviorInfo.IsActive() {
		for i := range res.Data {
			res.Data[i].IsGrey = true
		}
	}

	if r.WithInstructions {
		alarmIds := make([]string, 0, len(res.Data))
		for _, v := range res.Data {
			if v.AlarmID != "" {
				alarmIds = append(alarmIds, v.AlarmID)
			}
		}

		assignedInstructionsMap, err := s.alarmStore.GetAssignedInstructionsMap(ctx, alarmIds)
		if err != nil {
			return nil, err
		}

		statusesByAlarm, err := s.alarmStore.GetInstructionExecutionStatuses(ctx, alarmIds)
		if err != nil {
			return nil, err
		}

		for idx, v := range res.Data {
			sort.Slice(assignedInstructionsMap[v.AlarmID], func(i, j int) bool {
				return assignedInstructionsMap[v.AlarmID][i].Name < assignedInstructionsMap[v.AlarmID][j].Name
			})

			assignedInstructions := assignedInstructionsMap[v.AlarmID]
			if assignedInstructions == nil {
				assignedInstructions = make([]alarmapi.AssignedInstruction, 0)
			}
			res.Data[idx].AssignedInstructions = &assignedInstructions
			res.Data[idx].IsAutoInstructionRunning = statusesByAlarm[v.AlarmID].AutoRunning
			res.Data[idx].IsAllAutoInstructionsCompleted = statusesByAlarm[v.AlarmID].AutoAllCompleted
			res.Data[idx].IsAutoInstructionFailed = statusesByAlarm[v.AlarmID].AutoFailed
			res.Data[idx].IsManualInstructionRunning = statusesByAlarm[v.AlarmID].ManualRunning
			res.Data[idx].IsManualInstructionWaitingResult = statusesByAlarm[v.AlarmID].ManualWaitingResult
		}
	}

	err = s.fillLinks(ctx, apiKey, &res)
	if err != nil {
		s.logger.Err(err).Msg("cannot fill links")
	}

	return &res, nil
}

func (s *store) fillLinks(ctx context.Context, apiKey string, result *EntityAggregationResult) error {
	reqEntities := make([]common.FetchLinksRequestItem, 0, len(result.Data))
	entities := make(map[string][]int, len(result.Data))
	for i, entity := range result.Data {
		if _, ok := entities[entity.ID]; !ok {
			reqEntities = append(reqEntities, common.FetchLinksRequestItem{
				EntityID: entity.ID,
			})
			entities[entity.ID] = make([]int, 0, 1)
		}
		// map entity ID with record number in result.Data list
		entities[entity.ID] = append(entities[entity.ID], i)
	}
	res, err := s.linksFetcher.Fetch(ctx, apiKey, common.FetchLinksRequest{Entities: reqEntities})
	if err != nil || res == nil {
		return err
	}

	for _, rec := range res.Data {
		if l, ok := entities[rec.EntityID]; ok {
			for _, i := range l {
				result.Data[i].Links = make([]WeatherLink, 0, len(rec.Links))
				for category, link := range rec.Links {
					result.Data[i].Links = append(result.Data[i].Links, WeatherLink{
						Category: category,
						Links:    link,
					})
				}
			}
		}
	}

	return nil
}
