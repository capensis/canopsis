package pbehaviorics

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	pbehaviorapi "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	GetOneBy(ctx context.Context, id string) (*pbehaviorapi.Response, error)
	FindMaxPriority(ctx context.Context) (int64, error)
	FindMinPriority(ctx context.Context) (int64, error)
}

func NewStore(dbClient mongo.DbClient, authorProvider author.Provider) Store {
	return &store{
		dbClient:       dbClient,
		authorProvider: authorProvider,
	}
}

type store struct {
	dbClient       mongo.DbClient
	authorProvider author.Provider
}

func (s *store) GetOneBy(ctx context.Context, id string) (*pbehaviorapi.Response, error) {
	return pbehaviorapi.NewStore(s.dbClient, nil, nil, nil, s.authorProvider).GetOneBy(ctx, id)
}

func (s *store) FindMaxPriority(ctx context.Context) (int64, error) {
	return s.findPriority(ctx, true)
}

func (s *store) FindMinPriority(ctx context.Context) (int64, error) {
	return s.findPriority(ctx, false)
}

func (s *store) findPriority(ctx context.Context, desc bool) (int64, error) {
	coll := s.dbClient.Collection(mongo.PbehaviorTypeMongoCollection)
	sort := 1
	if desc {
		sort = -1
	}
	cursor, err := coll.Find(ctx, bson.M{},
		options.Find().SetSort(bson.M{"priority": sort}).SetLimit(1))
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		var t pbehavior.Type
		err = cursor.Decode(&t)
		if err != nil {
			return 0, err
		}

		return t.Priority, nil
	}

	return 0, nil
}
