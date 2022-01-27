package pbehaviortimespan

import (
	"context"
	"errors"
	"sort"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"github.com/teambition/rrule-go"
	"go.mongodb.org/mongo-driver/bson"
)

type Service interface {
	GetTimespans(ctx context.Context, request TimespansRequest) ([]timespansItemResponse, error)
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

func (s *service) GetTimespans(ctx context.Context, r TimespansRequest) ([]timespansItemResponse, error) {
	location := s.timezoneConfigProvider.Get().Location
	startAt := r.StartAt.Time.In(location)
	var endAt time.Time
	if r.EndAt == nil {
		endAt = r.ViewTo.Time
	} else {
		endAt = r.EndAt.Time
	}
	endAt = endAt.In(location)
	var event pbehavior.Event

	if r.RRule == "" {
		event = pbehavior.NewEvent(startAt, endAt)
	} else {
		rOption, err := rrule.StrToROption(r.RRule)
		if err != nil {
			panic(err)
		}

		event = pbehavior.NewRecEvent(startAt, endAt, rOption)
	}

	var spans []timespan.Span
	viewSpan := timespan.New(
		r.ViewFrom.Time.In(location),
		r.ViewTo.Time.In(location),
	)

	exdates, err := s.getExdates(ctx, r)
	if err != nil {
		return nil, err
	}
	var pbhTypes map[string]pbehavior.Type

	if r.ByDate {
		spans, err = pbehavior.GetDateSpans(event, viewSpan, exdates)
		if err != nil {
			return nil, err
		}

	} else {
		spans, err = pbehavior.GetTimeSpans(event, viewSpan, exdates)
		if err != nil {
			return nil, err
		}
		// append exception spans with types
		pbhTypeIDs := make([]string, 0, len(exdates))
		for _, exspan := range exdates {
			if exspan.Type() != "" {
				spans = append(spans, exspan)
				pbhTypeIDs = append(pbhTypeIDs, exspan.Type())
			}
		}
		pbhTypes, err = s.findPbhTypes(ctx, pbhTypeIDs)
		if err != nil {
			return nil, err
		}
	}

	sort.Sort(timespan.BySpans(spans))
	res := make([]timespansItemResponse, len(spans))
	for i, span := range spans {
		var spanType *pbehavior.Type
		if spanTypeId := span.Type(); spanTypeId != "" && pbhTypes != nil {
			if st, ok := pbhTypes[spanTypeId]; ok {
				spanType = &st
			}
		}
		res[i] = timespansItemResponse{
			From: types.NewCpsTime(span.From().Unix()),
			To:   types.NewCpsTime(span.To().Unix()),
			Type: spanType,
		}
	}

	return res, nil
}

func (s *service) getExdates(ctx context.Context, r TimespansRequest) ([]timespan.Span, error) {
	res := make([]timespan.Span, 0, len(r.Exdates))
	location := s.timezoneConfigProvider.Get().Location

	for _, v := range r.Exdates {
		res = append(res, timespan.TypedNew(
			v.Begin.Time.In(location),
			v.End.Time.In(location),
			v.Type,
		))
	}

	exceptions, err := s.findExceptions(ctx, r.Exceptions)
	if err != nil {
		return nil, err
	}

	for _, ex := range exceptions {
		for _, v := range ex.Exdates {
			res = append(res, timespan.TypedNew(
				v.Begin.Time.In(location),
				v.End.Time.In(location),
				v.Type,
			))
		}
	}

	return res, nil
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

func (s *service) findPbhTypes(ctx context.Context, ids []string) (map[string]pbehavior.Type, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	cursor, err := s.typeCollection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	res := make(map[string]pbehavior.Type, len(ids))

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
