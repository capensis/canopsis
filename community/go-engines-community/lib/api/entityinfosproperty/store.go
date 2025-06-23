package entityinfosproperty

import (
	"cmp"
	"context"
	"errors"
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
}

func NewStore(
	dbClient libmongo.DbClient,
	authorProvider author.Provider,
) Store {
	return &store{
		dbClient:       dbClient,
		dbCollection:   dbClient.Collection(libmongo.EntityInfosPropertyCollection),
		authorProvider: authorProvider,

		defaultSearchByFields: []string{"key", "description", "alias"},
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
			dbClient.Collection(libmongo.DeclareTicketRuleMongoCollection),
			dbClient.Collection(libmongo.InstructionMongoCollection),
			dbClient.Collection(libmongo.DynamicInfosRulesMongoCollection),
			dbClient.Collection(libmongo.KpiFilterMongoCollection),
			dbClient.Collection(libmongo.MetaAlarmRulesMongoCollection),
			dbClient.Collection(libmongo.ScenarioMongoCollection),
		},
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
				return s.parseDupError(err)
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

	if query.Type != "" {
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

		err := s.dbCollection.FindOneAndUpdate(ctx, bson.M{"_id": r.ID}, bson.M{"$set": r}, options.FindOneAndUpdate().SetReturnDocument(options.Before)).Decode(&oldProp)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil
			}

			if mongo.IsDuplicateKeyError(err) {
				return s.parseDupError(err)
			}

			return err
		}

		query := []bson.M{
			{
				"$set": bson.M{
					"entity_pattern": getUpdateAliasQuery("entity_pattern", oldProp.Alias, r.Alias),
				},
			},
		}

		for _, collection := range s.linkedCollections {
			if collection.Name() == libmongo.MetaAlarmRulesMongoCollection {
				query[0]["$set"].(bson.M)["total_entity_pattern"] = getUpdateAliasQuery("total_entity_pattern", oldProp.Alias, r.Alias) //nolint:forcetypeassert
			}

			if collection.Name() == libmongo.ScenarioMongoCollection {
				query = []bson.M{
					{
						"$set": bson.M{
							"actions": bson.M{
								"$map": bson.M{
									"input": "$actions",
									"as":    "action",
									"in": bson.M{
										"$mergeObjects": []any{
											"$$action",
											bson.M{
												"entity_pattern": getUpdateAliasQuery("$action.entity_pattern", oldProp.Alias, r.Alias),
											},
										},
									},
								},
							},
						},
					},
				}
			}

			_, err = collection.UpdateMany(ctx, bson.M{"aliases": r.ID}, query)
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

		query := []bson.M{
			{
				"$set": bson.M{
					"aliases": bson.M{
						"$filter": bson.M{
							"input": "$aliases",
							"as":    "a",
							"cond":  bson.M{"$ne": bson.A{"$$a", id}},
						},
					},
					"entity_pattern": getUpdateAliasQuery("entity_pattern", oldProp.Alias, ""),
				},
			},
		}

		for _, collection := range s.linkedCollections {
			if collection.Name() == libmongo.MetaAlarmRulesMongoCollection {
				query[0]["$set"].(bson.M)["total_entity_pattern"] = getUpdateAliasQuery("total_entity_pattern", oldProp.Alias, "") //nolint:forcetypeassert
			}

			if collection.Name() == libmongo.ScenarioMongoCollection {
				query = []bson.M{
					{
						"$set": bson.M{
							"actions": bson.M{
								"$map": bson.M{
									"input": "$actions",
									"as":    "action",
									"in": bson.M{
										"$mergeObjects": []any{
											"$$action",
											bson.M{
												"entity_pattern": getUpdateAliasQuery("$action.entity_pattern", oldProp.Alias, ""),
											},
										},
									},
								},
							},
						},
					},
				}
			}

			_, err = collection.UpdateMany(ctx, bson.M{"aliases": id}, query)
			if err != nil {
				return err
			}
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
		case "key":
			return common.NewValidationError("key", "Key already exists.")
		case "alias":
			return common.NewValidationError("alias", "Alias already exists.")
		default:
			return common.NewValidationError(matchedStr, matchedStr+" already exists.")
		}
	}

	return fmt.Errorf("can't parse duplication error: %w", err)
}

func getUpdateAliasQuery(patternFieldName, oldAlias, newAlias string) bson.M {
	return bson.M{
		"$map": bson.M{
			"input": "$" + patternFieldName,
			"as":    "subarray",
			"in": bson.M{
				"$map": bson.M{
					"input": "$$subarray",
					"as":    "elem",
					"in": bson.M{
						"$mergeObjects": []any{
							"$$elem",
							bson.M{
								"alias": bson.M{
									"$cond": []any{
										bson.M{"$eq": []any{"$$elem.alias", oldAlias}},
										newAlias,
										"$$elem.alias",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
