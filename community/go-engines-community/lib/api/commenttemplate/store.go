package commenttemplate

import (
	"cmp"
	"context"
	"fmt"
	"regexp"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Store interface {
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	GetByID(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
}

type store struct {
	dbClient         libmongo.DbClient
	dbCollection     libmongo.DbCollection
	widgetCollection libmongo.DbCollection
	authorProvider   author.Provider

	defaultSearchByFields []string
	dupErrorRegexp        *regexp.Regexp
}

func NewStore(
	dbClient libmongo.DbClient,
	authorProvider author.Provider,
) Store {
	return &store{
		dbClient:         dbClient,
		dbCollection:     dbClient.Collection(libmongo.CommentTemplateMongoCollection),
		widgetCollection: dbClient.Collection(libmongo.WidgetMongoCollection),
		authorProvider:   authorProvider,

		defaultSearchByFields: []string{"_id", "name"},
		dupErrorRegexp:        regexp.MustCompile(`{ ([^:]+)`),
	}
}

func (s *store) Insert(ctx context.Context, r EditRequest) (*Response, error) {
	now := datetime.NewCpsTime()

	doc := Document{
		ID:      utils.NewID(),
		Name:    r.Name,
		Fields:  r.Fields,
		Author:  r.Author,
		Created: now,
		Updated: now,
	}

	var res *Response

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = nil

		_, err := s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return s.parseDupError(err)
			}

			return err
		}

		res, err = s.GetByID(ctx, doc.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *store) GetByID(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		response := Response{}
		err := cursor.Decode(&response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	match := make([]bson.M, 0)
	pipeline := make([]bson.M, 0)

	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		match = append(match, filter)
	}

	if len(query.IDs) > 0 {
		match = append(match, bson.M{"_id": bson.M{"$in": query.IDs}})
	}

	if len(match) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": match}})
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		common.GetSortQuery(cmp.Or(query.SortBy, "_id"), query.Sort),
		s.authorProvider.Pipeline(),
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

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) Update(ctx context.Context, r EditRequest) (*Response, error) {
	doc := Document{
		Name:    r.Name,
		Fields:  r.Fields,
		Author:  r.Author,
		Updated: datetime.NewCpsTime(),
	}

	var res *Response

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = nil

		_, err := s.dbCollection.UpdateOne(
			ctx,
			bson.M{"_id": r.ID},
			bson.M{"$set": doc},
		)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return s.parseDupError(err)
			}

			return err
		}

		res, err = s.GetByID(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	var deleted int64

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		_, err = s.widgetCollection.UpdateMany(ctx, bson.M{"parameters.comment_templates": id},
			bson.M{"$pull": bson.M{"parameters.comment_templates": id}})
		if err != nil {
			return err
		}

		deleted, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		return err
	})

	return deleted > 0, err
}

func (s *store) parseDupError(err error) error {
	match := s.dupErrorRegexp.FindStringSubmatch(err.Error())
	if len(match) > 1 {
		matchedStr := match[1]

		switch matchedStr {
		case "name":
			return common.NewValidationError("name", "Name already exists.")
		default:
			return common.NewValidationError(matchedStr, matchedStr+" already exists.")
		}
	}

	return fmt.Errorf("can't parse duplication error: %w", err)
}
