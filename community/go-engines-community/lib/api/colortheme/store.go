package colortheme

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	Canopsis       = "canopsis"
	CanopsisDark   = "canopsis_dark"
	ColorBlind     = "color_blind"
	ColorBlindDark = "color_blind_dark"
)

type Store interface {
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	GetByID(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
}

type store struct {
	dbClient          libmongo.DbClient
	dbColorCollection libmongo.DbCollection
	dbUserCollection  libmongo.DbCollection
	authorProvider    author.Provider

	defaultSearchByFields []string
	defaultThemeIDs       map[string]struct{}
}

func NewStore(
	dbClient libmongo.DbClient,
	authorProvider author.Provider,
) Store {
	return &store{
		dbClient:              dbClient,
		dbColorCollection:     dbClient.Collection(libmongo.ColorThemeCollection),
		dbUserCollection:      dbClient.Collection(libmongo.UserCollection),
		authorProvider:        authorProvider,
		defaultSearchByFields: []string{"_id", "name"},
		defaultThemeIDs: map[string]struct{}{
			Canopsis:       {},
			CanopsisDark:   {},
			ColorBlind:     {},
			ColorBlindDark: {},
		},
	}
}

func (s *store) Insert(ctx context.Context, r EditRequest) (*Response, error) {
	now := datetime.NewCpsTime()

	r.ID = utils.NewID()
	r.Created = now
	r.Updated = now
	r.Deletable = true

	var response *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		_, err := s.dbColorCollection.InsertOne(ctx, r)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return common.NewValidationError("name", "Name already exists.")
			}

			return err
		}

		response, err = s.GetByID(ctx, r.ID)
		return err
	})

	return response, err
}

func (s *store) GetByID(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.dbColorCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var res Response
		err = cursor.Decode(&res)
		if err != nil {
			return nil, err
		}

		return &res, nil
	}

	return nil, nil
}

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	var pipeline []bson.M

	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	sortBy := "name"
	if query.SortBy != "" {
		sortBy = query.SortBy
	}

	cursor, err := s.dbColorCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		common.GetSortQuery(sortBy, query.Sort),
	))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var result AggregationResult
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (s *store) Update(ctx context.Context, r EditRequest) (*Response, error) {
	if s.isDefaultTheme(r.ID) {
		return nil, ErrDefaultTheme
	}

	r.Updated = datetime.NewCpsTime()
	r.Deletable = true

	var response *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		res, err := s.dbColorCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": r})
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return common.NewValidationError("name", "Name already exists.")
			}

			return err
		}

		if res.MatchedCount == 0 {
			return nil
		}

		response, err = s.GetByID(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	if s.isDefaultTheme(id) {
		return false, ErrDefaultTheme
	}

	var deleted int64

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		res, err := s.dbColorCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		deleted, err = s.dbColorCollection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil || deleted == 0 {
			return err
		}

		_, err = s.dbUserCollection.UpdateMany(ctx, bson.M{"ui_theme": id},
			bson.M{"$set": bson.M{"ui_theme": Canopsis, "author": userID, "updated": datetime.NewCpsTime()}})
		return err
	})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s *store) isDefaultTheme(id string) bool {
	_, ok := s.defaultThemeIDs[id]
	return ok
}
