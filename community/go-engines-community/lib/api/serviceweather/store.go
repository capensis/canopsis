package serviceweather

import (
	"context"
	"errors"
	"sort"
	"time"

	alarmapi "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	libtypes "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Find(context.Context, ListRequest) (*AggregationResult, error)
	FindEntities(ctx context.Context, id string, query EntitiesListRequest, userId string) (*EntityAggregationResult, error)
}

func NewStore(
	dbClient mongo.DbClient,
	linkGenerator link.Generator,
	alarmStore alarmapi.Store,
	timezoneConfigProvider config.TimezoneConfigProvider,
	authorProvider author.Provider,
	logger zerolog.Logger,
) Store {
	return &store{
		dbClient:         dbClient,
		dbCollection:     dbClient.Collection(mongo.EntityMongoCollection),
		userDbCollection: dbClient.Collection(mongo.UserCollection),

		alarmStore:     alarmStore,
		linkGenerator:  linkGenerator,
		authorProvider: authorProvider,

		timezoneConfigProvider: timezoneConfigProvider,

		logger: logger,
	}
}

type store struct {
	dbClient         mongo.DbClient
	dbCollection     mongo.DbCollection
	userDbCollection mongo.DbCollection

	linkGenerator  link.Generator
	alarmStore     alarmapi.Store
	authorProvider author.Provider

	timezoneConfigProvider config.TimezoneConfigProvider

	logger zerolog.Logger
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline, err := s.getQueryBuilder().CreateListAggregationPipeline(ctx, r)
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

func (s *store) FindEntities(ctx context.Context, id string, r EntitiesListRequest, userId string) (*EntityAggregationResult, error) {
	err := s.dbCollection.FindOne(ctx, bson.M{
		"_id":     id,
		"type":    libtypes.EntityTypeService,
		"enabled": true,
	}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	location := s.timezoneConfigProvider.Get().Location
	now := datetime.CpsTime{Time: time.Now().In(location)}
	pipeline, err := s.getQueryBuilder().CreateListDependenciesAggregationPipeline(id, r, now)
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

	var alarmIds []string
	if r.WithDeclareTickets || r.WithInstructions {
		alarmIds = make([]string, 0, len(res.Data))
		for _, v := range res.Data {
			if v.AlarmID != "" {
				alarmIds = append(alarmIds, v.AlarmID)
			}
		}
	}

	if r.WithDeclareTickets {
		assignedDeclareTicketsMap, err := s.alarmStore.GetAssignedDeclareTicketsMap(ctx, alarmIds)
		if err != nil {
			return nil, err
		}

		for idx, v := range res.Data {
			sort.Slice(assignedDeclareTicketsMap[v.AlarmID], func(i, j int) bool {
				return assignedDeclareTicketsMap[v.AlarmID][i].Name < assignedDeclareTicketsMap[v.AlarmID][j].Name
			})

			res.Data[idx].AssignedDeclareTicketRules = assignedDeclareTicketsMap[v.AlarmID]
		}
	}

	if r.WithInstructions {
		assignedInstructionsMap, err := s.alarmStore.GetAssignedInstructionsMap(ctx, alarmIds)
		if err != nil {
			return nil, err
		}

		statusesByAlarm, err := s.alarmStore.GetInstructionExecutionStatuses(ctx, alarmIds, assignedInstructionsMap)
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
			res.Data[idx].InstructionExecutionIcon = statusesByAlarm[v.AlarmID].Icon
			res.Data[idx].RunningManualInstructions = statusesByAlarm[v.AlarmID].RunningManualInstructions
			res.Data[idx].RunningAutoInstructions = statusesByAlarm[v.AlarmID].RunningAutoInstructions
			res.Data[idx].FailedManualInstructions = statusesByAlarm[v.AlarmID].FailedManualInstructions
			res.Data[idx].FailedAutoInstructions = statusesByAlarm[v.AlarmID].FailedAutoInstructions
			res.Data[idx].SuccessfulManualInstructions = statusesByAlarm[v.AlarmID].SuccessfulManualInstructions
			res.Data[idx].SuccessfulAutoInstructions = statusesByAlarm[v.AlarmID].SuccessfulAutoInstructions
		}
	}

	err = s.fillLinks(ctx, &res, userId)
	if err != nil {
		s.logger.Err(err).Msg("cannot fill links")
	}

	return &res, nil
}

func (s *store) fillLinks(ctx context.Context, result *EntityAggregationResult, userId string) error {
	user, err := s.findUser(ctx, userId)
	if err != nil {
		return err
	}

	ids := make([]string, len(result.Data))
	for i, v := range result.Data {
		ids[i] = v.ID
	}

	linksByEntityId, err := s.linkGenerator.GenerateForEntities(ctx, ids, user)
	if err != nil || len(linksByEntityId) == 0 {
		return err
	}

	for i, v := range result.Data {
		result.Data[i].Links = linksByEntityId[v.ID]
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

func (s *store) findUser(ctx context.Context, id string) (link.User, error) {
	user := link.User{}
	cursor, err := s.userDbCollection.Aggregate(ctx, []bson.M{
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
