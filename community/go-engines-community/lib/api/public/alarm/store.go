package alarm

import (
	"cmp"
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/sortparam"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	liblink "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Store interface {
	Find(ctx context.Context, r PaginatedListRequest, userID string) (pagination.ListResponse[Response], error)
	Search(ctx context.Context, r PaginatedSearchRequest, userID string) (pagination.ListResponse[Response], error)
	FindSteps(ctx context.Context, id string, r StepRequest) (pagination.ListResponse[Step], error)
}

type store struct {
	dbClient             mongo.DbClient
	dbCollection         mongo.DbCollection
	dbResolvedCollection mongo.DbCollection
	transformer          patternfields.Transformer
	generator            link.Generator
	userProvider         liblink.UserProvider
}

func NewStore(dbClient mongo.DbClient, transformer patternfields.Transformer, generator link.Generator, userProvider liblink.UserProvider) Store {
	return &store{
		dbClient:             dbClient,
		dbCollection:         dbClient.Collection(mongo.AlarmMongoCollection),
		dbResolvedCollection: dbClient.Collection(mongo.ResolvedAlarmMongoCollection),
		transformer:          transformer,
		generator:            generator,
		userProvider:         userProvider,
	}
}

func (s *store) Find(ctx context.Context, r PaginatedListRequest, userID string) (res pagination.ListResponse[Response], err error) {
	r.MetaAlarms, err = ResolveMetaAlarms(ctx, r.MetaAlarms, s.dbCollection, s.dbResolvedCollection)
	if err != nil {
		return res, err
	}

	beforeLimit, afterLimit, sort, err := ListPipeline(r.ListRequest, datetime.NewCpsTime())
	if err != nil {
		return res, err
	}

	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		sort,
		afterLimit,
	)
	ar, err := db.AggregateAll[Response](ctx, s.collection(r.Resolution), pipeline)
	if err != nil {
		return res, err
	}

	user, err := s.userProvider.User(ctx, userID)
	if err != nil {
		return res, err
	}

	for i, item := range ar {
		ar[i].setDefaults()
		ar[i].Entity.Links = s.generator.GenerateForEntity(ctx, item.RawEntity, user)
		ar[i].Links = s.generator.GenerateForAlarm(ctx, item.RawAlarm, item.RawEntity, user)
	}

	res = pagination.NewResponse[Response](r.Query, ar)

	return res, nil
}

func (s *store) Search(ctx context.Context, r PaginatedSearchRequest, userID string) (res pagination.ListResponse[Response], err error) {
	r.MetaAlarms, err = ResolveMetaAlarms(ctx, r.MetaAlarms, s.dbCollection, s.dbResolvedCollection)
	if err != nil {
		return res, err
	}

	if len(r.EntityPattern) > 0 {
		patternAliases := patternfields.GetAliases(r.EntityPattern)
		if len(patternAliases) != 0 {
			aliases, err := s.transformer.FetchAliases(ctx, patternAliases)
			if err != nil {
				return res, err
			}

			var valErrs validator.ValidationErrors
			r.EntityPattern, _, valErrs = s.transformer.ApplyAliases(r.EntityPattern, aliases, "EntityPattern")
			if len(valErrs) > 0 {
				return res, validation.NewError(valErrs, r)
			}
		}
	}

	beforeLimit, afterLimit, sort, err := SearchPipeline(r.SearchRequest, datetime.NewCpsTime())
	if err != nil {
		return res, err
	}

	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		sort,
		afterLimit,
	)
	ar, err := db.AggregateAll[Response](ctx, s.collection(r.Resolution), pipeline)
	if err != nil {
		return res, err
	}

	user, err := s.userProvider.User(ctx, userID)
	if err != nil {
		return res, err
	}

	for i, item := range ar {
		ar[i].setDefaults()
		ar[i].Entity.Links = s.generator.GenerateForEntity(ctx, item.RawEntity, user)
		ar[i].Links = s.generator.GenerateForAlarm(ctx, item.RawAlarm, item.RawEntity, user)
	}

	res = pagination.NewResponse[Response](r.Query, ar)

	return res, nil
}

func (s *store) FindSteps(ctx context.Context, id string, r StepRequest) (res pagination.ListResponse[Step], err error) {
	sort, err := sortparam.Validate(cmp.Or(r.Sort, "t"), stepSortableFields)
	if err != nil {
		return res, err
	}

	filter := bson.M{}
	if len(r.Types) > 0 {
		filter["_t"] = bson.M{"$in": r.Types}
	} else {
		filter["_t"] = bson.M{"$in": types.GetAlarmStepTypes()}
	}

	if len(r.Initiators) > 0 {
		filter["initiator"] = bson.M{"$in": r.Initiators}
	}

	if r.From != nil || r.To != nil {
		q := bson.M{}
		if r.From != nil {
			q["$gte"] = r.From
		}

		if r.To != nil {
			q["$lte"] = r.To
		}

		filter["t"] = q
	}

	beforeLimit := []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$unwind": bson.M{"path": "$v.steps", "includeArrayIndex": "_id"}},
		{"$replaceRoot": bson.M{"newRoot": "$v.steps"}},
	}
	if len(filter) > 0 {
		beforeLimit = append(beforeLimit, bson.M{"$match": filter})
	}

	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		sort,
	)
	ar, err := db.AggregateAll[Step](ctx, s.dbCollection, pipeline)
	if err != nil {
		return res, err
	}

	if len(ar) == 0 {
		ar, err = db.AggregateAll[Step](ctx, s.dbResolvedCollection, pipeline)
		if err != nil {
			return res, err
		}

		if len(ar) == 0 {
			return res, httperror.ErrNotFound
		}
	}

	res = pagination.NewResponse[Step](r.Query, ar)

	return res, nil
}

func (s *store) collection(resolution string) mongo.DbCollection {
	if resolution == ResolutionResolved {
		return s.dbResolvedCollection
	}

	return s.dbCollection
}

func ResolveMetaAlarms(ctx context.Context, metaAlarmIDs []string, dbCollection, dbResolvedCollection mongo.DbCollection) (entityIDs []string, err error) {
	if len(metaAlarmIDs) == 0 {
		return entityIDs, nil
	}

	pipeline := []bson.M{
		{"$match": bson.M{
			"_id":    bson.M{"$in": metaAlarmIDs},
			"v.meta": bson.M{"$nin": bson.A{nil, ""}},
		}},
		{"$project": bson.M{
			"_id": "$d",
		}},
	}
	entityIDs, err = db.AggregateIDs(ctx, dbCollection, pipeline)
	if err != nil {
		return entityIDs, err
	}

	resolvedEntityIDs, err := db.AggregateIDs(ctx, dbResolvedCollection, pipeline)
	if err != nil {
		return entityIDs, err
	}

	entityIDs = append(entityIDs, resolvedEntityIDs...)

	return entityIDs, nil
}
