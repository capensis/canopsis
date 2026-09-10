package entity

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	liblink "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
)

type Store interface {
	Find(ctx context.Context, r PaginatedListRequest, userID string) (pagination.ListResponse[Response], error)
	Search(ctx context.Context, r PaginatedSearchRequest, userID string) (pagination.ListResponse[Response], error)
}

type store struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
	transformer  patternfields.Transformer
	generator    link.Generator
	userProvider liblink.UserProvider
}

func NewStore(dbClient mongo.DbClient, transformer patternfields.Transformer, generator link.Generator, userProvider liblink.UserProvider) Store {
	return &store{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.EntityMongoCollection),
		transformer:  transformer,
		generator:    generator,
		userProvider: userProvider,
	}
}

func (s *store) Find(ctx context.Context, r PaginatedListRequest, userID string) (res pagination.ListResponse[Response], err error) {
	beforeLimit, afterLimit, sort, err := ListPipeline(r.ListRequest)
	if err != nil {
		return res, err
	}

	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		sort,
		afterLimit,
	)
	ar, err := db.AggregateAll[Response](ctx, s.dbCollection, pipeline)
	if err != nil {
		return res, err
	}

	user, err := s.userProvider.User(ctx, userID)
	if err != nil {
		return res, err
	}

	for i, item := range ar {
		ar[i].SetDefaults()
		ar[i].Links = s.generator.GenerateForEntity(ctx, item.Raw, user)
	}

	res = pagination.NewResponse[Response](r.Query, ar)

	return res, nil
}

func (s *store) Search(ctx context.Context, r PaginatedSearchRequest, userID string) (res pagination.ListResponse[Response], err error) {
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

	beforeLimit, afterLimit, sort, err := SearchPipeline(r.SearchRequest)
	if err != nil {
		return res, err
	}

	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		sort,
		afterLimit,
	)
	ar, err := db.AggregateAll[Response](ctx, s.dbCollection, pipeline)
	if err != nil {
		return res, err
	}

	user, err := s.userProvider.User(ctx, userID)
	if err != nil {
		return res, err
	}

	for i, item := range ar {
		ar[i].SetDefaults()
		ar[i].Links = s.generator.GenerateForEntity(ctx, item.Raw, user)
	}

	res = pagination.NewResponse[Response](r.Query, ar)

	return res, nil
}
