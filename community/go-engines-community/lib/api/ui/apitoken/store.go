package apitoken

import (
	"cmp"
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	libtoken "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/security/token"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Store interface {
	Find(ctx context.Context, r ListRequest, userID string) (AggregationResult, error)
	Insert(ctx context.Context, r EditRequest) (Response, error)
	Delete(ctx context.Context, id, userID string) error
}

func NewStore(dbClient mongo.DbClient, authorProvider author.Provider) Store {
	return &store{
		dbClient:       dbClient,
		dbCollection:   dbClient.Collection(mongo.PublicAPITokenCollection),
		defaultSortBy:  "name",
		authorProvider: authorProvider,
	}
}

type store struct {
	dbClient       mongo.DbClient
	dbCollection   mongo.DbCollection
	defaultSortBy  string
	authorProvider author.Provider
}

func (s *store) Find(ctx context.Context, r ListRequest, userID string) (AggregationResult, error) {
	beforeLimit := []bson.M{
		{"$match": bson.M{
			"user": userID,
		}},
		{"$project": bson.M{
			"hash": 0,
		}},
	}
	sort := mongoquery.GetSortQuery(cmp.Or(r.SortBy, s.defaultSortBy), r.Sort)
	afterLimit := s.userPipeline()
	afterLimit = append(afterLimit, sort)

	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		sort,
		afterLimit,
	)

	return db.AggregateOne[AggregationResult](ctx, s.dbCollection, pipeline)
}

// Insert creates a token and returns it along with its generated value.
// The value is returned only here: the collection keeps its hash and a masked copy,
// so it cannot be read again afterwards.
func (s *store) Insert(ctx context.Context, r EditRequest) (Response, error) {
	token, err := libtoken.Generate()
	if err != nil {
		return Response{}, err
	}

	now := datetime.NewCpsTime()
	doc := Token{
		ID:       utils.NewID(),
		Name:     r.Name,
		Value:    libtoken.Mask(token),
		Hash:     libtoken.Hash(token),
		User:     r.Author,
		Created:  &now,
		Accessed: nil,
	}
	if r.Duration != nil {
		doc.Expired = new(r.Duration.AddTo(now))
	}

	pipeline := append(
		[]bson.M{{"$match": bson.M{"_id": doc.ID}}, {"$project": bson.M{"hash": 0}}},
		s.userPipeline()...,
	)

	var res Response
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = Response{}
		_, err := s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}

		res, err = db.AggregateOne[Response](ctx, s.dbCollection, pipeline)

		return err
	})
	if err != nil {
		return Response{}, err
	}

	// The aggregation read back the masked value, replace it by the real one.
	res.Value = token

	return res, nil
}

func (s *store) Delete(ctx context.Context, id, userID string) error {
	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id, "user": userID})
	if err != nil {
		return err
	}

	if deleted == 0 {
		return httperror.ErrNotFound
	}

	return nil
}

func (s *store) userPipeline() []bson.M {
	author := s.authorProvider.PipelineForField("user")
	roles := rolesPipeline()

	pipeline := make([]bson.M, 0, len(author)+len(roles))
	pipeline = append(pipeline, author...)

	return append(pipeline, roles...)
}

func rolesPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.UserCollection,
			"localField":   "user._id",
			"foreignField": "_id",
			"as":           "user_with_roles",
		}},
		{"$unwind": bson.M{"path": "$user_with_roles", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"user.roles": "$user_with_roles.roles",
		}},
		{"$project": bson.M{"user_with_roles": 0}},
		{"$unwind": bson.M{
			"path":                       "$user.roles",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "role_index",
		}},
		{"$lookup": bson.M{
			"from":         mongo.RoleCollection,
			"localField":   "user.roles",
			"foreignField": "_id",
			"as":           "role",
		}},
		{"$unwind": bson.M{"path": "$role", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"role_index": 1}},
		{"$group": bson.M{
			"_id":  "$_id",
			"data": bson.M{"$first": "$$ROOT"},
			"roles": bson.M{"$push": bson.M{
				"$cond": bson.M{
					"if":   "$role._id",
					"then": "$role",
					"else": "$$REMOVE",
				},
			}},
		}},
		{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$data",
			bson.M{"roles": "$roles"},
		}}}},
	}
}
