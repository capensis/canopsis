package entityinfosproperty

import (
	"cmp"
	"context"
	"errors"
	"regexp"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	GetByID(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
}

type store struct {
	dbClient       libmongo.DbClient
	dbCollection   libmongo.DbCollection
	authorProvider author.Provider

	defaultSearchByFields []string
	dupErrorRegexp        *regexp.Regexp

	linkedCollections []libmongo.DbCollection
	dupErrorParser    validation.DuplicateErrorParser
}

func NewStore(
	dbClient libmongo.DbClient,
	authorProvider author.Provider,
) Store {
	return &store{
		dbClient:       dbClient,
		dbCollection:   dbClient.Collection(libmongo.EntityInfosPropertyCollection),
		authorProvider: authorProvider,

		defaultSearchByFields: []string{"name", "description", "alias"},
		dupErrorRegexp:        regexp.MustCompile(`{ ([^:]+)`),

		linkedCollections: []libmongo.DbCollection{
			dbClient.Collection(libmongo.AlarmTagCollection),
			dbClient.Collection(libmongo.EntityMongoCollection),
			dbClient.Collection(libmongo.EventFilterRuleCollection),
			dbClient.Collection(libmongo.FlappingRuleMongoCollection),
			dbClient.Collection(libmongo.IdleRuleMongoCollection),
			dbClient.Collection(libmongo.LinkRuleMongoCollection),
			dbClient.Collection(libmongo.PatternMongoCollection),
			dbClient.Collection(libmongo.PbehaviorMongoCollection),
			dbClient.Collection(libmongo.ResolveRuleMongoCollection),
			dbClient.Collection(libmongo.WidgetFiltersMongoCollection),
			dbClient.Collection(libmongo.DeclareTicketRuleCollection),
			dbClient.Collection(libmongo.InstructionMongoCollection),
			dbClient.Collection(libmongo.DynamicInfosRulesMongoCollection),
			dbClient.Collection(libmongo.KpiFilterMongoCollection),
			dbClient.Collection(libmongo.MetaAlarmRulesMongoCollection),
			dbClient.Collection(libmongo.ScenarioCollection),
		},
		dupErrorParser: validation.NewDuplicateErrorParser(map[string]string{
			"name":  "Name already exists.",
			"alias": "Alias already exists.",
		}),
	}
}

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Response, error) {
	now := datetime.NewCpsTime()

	r.ID = utils.NewID()
	r.Created = now
	r.Updated = now

	var response *Response

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		_, err := s.dbCollection.InsertOne(ctx, r)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err)
			}

			return err
		}

		response, err = s.GetByID(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return response, nil
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

	return nil, nil
}

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	pipeline := s.authorProvider.Pipeline()

	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if filter == nil {
		filter = bson.M{}
	}

	if query.Type != nil {
		filter["type"] = query.Type
	}

	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		common.GetSortQuery(cmp.Or(query.SortBy, "created"), query.Sort),
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

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Response, error) {
	r.Updated = datetime.NewCpsTime()

	var response *Response

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		var oldProp InfoProperty

		unset := bson.M{}
		if r.Alias == "" {
			unset["alias"] = ""
		}

		err := s.dbCollection.FindOneAndUpdate(ctx, bson.M{"_id": r.ID}, bson.M{"$set": r, "$unset": unset}, options.FindOneAndUpdate().SetReturnDocument(options.Before)).Decode(&oldProp)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil
			}

			if mongo.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err)
			}

			return err
		}

		if oldProp.Alias != "" {
			if r.Alias == "" {
				err = s.removeAliasFromLinkedCollections(ctx, r.ID, oldProp.Alias)
			} else if r.Alias != oldProp.Alias {
				err = s.updateAliasInLinkedCollections(ctx, r.ID, oldProp.Alias, r.Alias)
			}

			if err != nil {
				return err
			}
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
	var deleted int64

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		var oldProp InfoProperty

		// required to get the author in action log listener.
		err := s.dbCollection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}}, options.FindOneAndUpdate().SetReturnDocument(options.Before)).Decode(&oldProp)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil
			}
		}

		if oldProp.Alias != "" {
			err = s.removeAliasFromLinkedCollections(ctx, id, oldProp.Alias)
			if err != nil {
				return err
			}
		}

		deleted, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		return err
	})

	return deleted > 0, err
}

func (s *store) updateAliasInLinkedCollections(ctx context.Context, id, oldAlias, newAlias string) error {
	for _, collection := range s.linkedCollections {
		var update bson.M
		switch collection.Name() {
		case libmongo.MetaAlarmRulesMongoCollection:
			update = bson.M{
				"entity_pattern.$[].$[i].alias":       newAlias,
				"total_entity_pattern.$[].$[i].alias": newAlias,
			}
		case libmongo.ScenarioCollection:
			update = bson.M{
				"actions.$[].entity_pattern.$[].$[i].alias": newAlias,
			}
		default:
			update = bson.M{
				"entity_pattern.$[].$[i].alias": newAlias,
			}
		}

		_, err := collection.UpdateMany(
			ctx,
			bson.M{"aliases": id},
			bson.M{"$set": update},
			options.UpdateMany().SetArrayFilters([]any{bson.M{"i.alias": oldAlias}}),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) removeAliasFromLinkedCollections(ctx context.Context, id, oldAlias string) error {
	for _, collection := range s.linkedCollections {
		var unset bson.M

		switch collection.Name() {
		case libmongo.MetaAlarmRulesMongoCollection:
			unset = bson.M{
				"entity_pattern.$[].$[i].alias":       "",
				"total_entity_pattern.$[].$[i].alias": "",
			}
		case libmongo.ScenarioCollection:
			unset = bson.M{
				"actions.$[].entity_pattern.$[].$[i].alias": "",
			}
		default:
			unset = bson.M{
				"entity_pattern.$[].$[i].alias": "",
			}
		}

		_, err := collection.UpdateMany(
			ctx,
			bson.M{"aliases": id},
			bson.M{"$unset": unset, "$pull": bson.M{"aliases": id}},
			options.UpdateMany().SetArrayFilters([]any{bson.M{"i.alias": oldAlias}}),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
