package pbehaviortimespan

import (
	"context"
	"errors"
	"sort"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Service interface {
	GetTimespans(ctx context.Context, request TimespansRequest) ([]ItemResponse, error)
}

func NewService(dbClient mongo.DbClient, timezoneConfigProvider config.TimezoneConfigProvider) Service {
	return &service{
		exceptionCollection:    dbClient.Collection(mongo.PbehaviorExceptionMongoCollection),
		typeCollection:         dbClient.Collection(mongo.PbehaviorTypeMongoCollection),
		timezoneConfigProvider: timezoneConfigProvider,
	}
}

type service struct {
	exceptionCollection    mongo.DbCollection
	typeCollection         mongo.DbCollection
	timezoneConfigProvider config.TimezoneConfigProvider
}

func (s *service) GetTimespans(ctx context.Context, r TimespansRequest) ([]ItemResponse, error) {
	location := s.timezoneConfigProvider.Get().Location
	viewSpan := timespan.New(
		r.ViewFrom.Time.In(location),
		r.ViewTo.Time.In(location),
	)

	exdates, err := s.getExdates(ctx, r)
	if err != nil {
		return nil, err
	}

	pbhTypes, err := s.findPbhTypes(ctx)
	if err != nil {
		return nil, err
	}

	defaultTypes, err := pbehavior.ResolveDefaultTypes(pbhTypes)
	if err != nil {
		return nil, err
	}
	eventComputer := pbehavior.NewEventComputer(pbhTypes, defaultTypes)
	params := pbehavior.PbhEventParams{
		Start:   r.StartAt,
		End:     r.EndAt,
		RRule:   r.RRule,
		Type:    r.Type,
		Exdates: exdates,
	}
	computed, err := eventComputer.Compute(params, viewSpan)
	if err != nil {
		return nil, err
	}

	res := make([]ItemResponse, len(computed))
	for i, v := range computed {
		res[i] = ItemResponse{
			From: datetime.CpsTime{Time: v.Span.From()},
			To:   datetime.CpsTime{Time: v.Span.To()},
		}
		if spanType, ok := pbhTypes[v.ID]; ok {
			res[i].Type = spanType
		}
	}

	sort.Slice(res, sortResponse(res))

	return res, nil
}

func (s *service) getExdates(ctx context.Context, r TimespansRequest) ([]pbehavior.Exdate, error) {
	exdates := make([]pbehavior.Exdate, len(r.Exdates))
	for i, v := range r.Exdates {
		exdates[i] = pbehavior.Exdate{
			Exdate: types.Exdate{
				Begin: v.Begin,
				End:   v.End,
			},
			Type: v.Type,
		}
	}

	exceptions, err := s.findExceptions(ctx, r.Exceptions)
	if err != nil {
		return nil, err
	}

	for _, ex := range exceptions {
		for _, v := range ex.Exdates {
			exdates = append(exdates, pbehavior.Exdate{
				Exdate: types.Exdate{
					Begin: v.Begin,
					End:   v.End,
				},
				Type: v.Type,
			})
		}
	}

	return exdates, nil
}

func (s *service) findExceptions(ctx context.Context, ids []string) ([]pbehavior.Exception, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	cursor, err := s.exceptionCollection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var exceptions []pbehavior.Exception
	err = cursor.All(ctx, &exceptions)
	if err != nil {
		return nil, err
	}

	if len(exceptions) != len(ids) {
		return nil, errors.New("unknown exceptions")
	}

	return exceptions, nil
}

func (s *service) findPbhTypes(ctx context.Context) (map[string]pbehavior.Type, error) {
	cursor, err := s.typeCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	res := make(map[string]pbehavior.Type)

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var pbhType pbehavior.Type

		if err = cursor.Decode(&pbhType); err != nil {
			return nil, err
		}
		res[pbhType.ID] = pbhType
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func sortResponse(response []ItemResponse) func(i, j int) bool {
	return func(i, j int) bool {
		dateLeft := utils.DateOf(response[i].From.Time)
		dateRight := utils.DateOf(response[j].From.Time)

		if dateLeft.Before(dateRight) {
			return true
		}
		if dateLeft.After(dateRight) {
			return false
		}

		if response[i].Type.Priority > response[j].Type.Priority {
			return true
		}
		if response[i].Type.Priority < response[j].Type.Priority {
			return false
		}

		if response[i].From.Before(response[j].From) {
			return true
		}
		if response[i].From.After(response[j].From) {
			return false
		}

		return response[i].To.Before(response[j].To)
	}
}
