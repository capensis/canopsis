package viewgroup

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Find(ctx context.Context, r ListRequest, authorizedViewIds, ownedPrivateIds []string) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*ViewGroup, error)
	Insert(ctx context.Context, r EditRequest) (*ViewGroup, error)
	Update(ctx context.Context, r EditRequest) (*ViewGroup, error)
	Delete(ctx context.Context, id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient, authorProvider author.Provider) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.ViewGroupMongoCollection),
		dbViewCollection:      dbClient.Collection(mongo.ViewMongoCollection),
		authorProvider:        authorProvider,
		defaultSearchByFields: []string{"_id", "title"},
	}
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	dbViewCollection      mongo.DbCollection
	authorProvider        author.Provider
	defaultSearchByFields []string
}

func (s *store) Find(ctx context.Context, r ListRequest, authorizedViewIds, ownedPrivateIds []string) (*AggregationResult, error) {
	var pipeline []bson.M
	if r.WithPrivate {
		pipeline = []bson.M{{"$match": bson.M{"$or": bson.A{bson.M{"is_private": false}, bson.M{"author": r.UserID}}}}}
		if r.WithViews {
			authorizedViewIds = append(authorizedViewIds, ownedPrivateIds...)
		}
	} else {
		pipeline = []bson.M{{"$match": bson.M{"is_private": false}}}
	}

	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sort := bson.M{"$sort": bson.D{{Key: "is_private", Value: 1}, {Key: "position", Value: 1}, {Key: "title", Value: 1}}}
	project := s.authorProvider.Pipeline()

	if r.WithFlags || r.WithViews {
		project = append(project,
			bson.M{"$addFields": bson.M{
				"group": "$$ROOT",
			}},
			bson.M{"$lookup": bson.M{
				"from":         mongo.ViewMongoCollection,
				"localField":   "_id",
				"foreignField": "group_id",
				"as":           "views",
			}},
		)
	}

	if r.WithFlags {
		project = append(project,
			bson.M{"$addFields": bson.M{
				"deletable": bson.M{"$eq": bson.A{"$views", bson.A{}}},
			}},
		)
		if !r.WithViews {
			project = append(project, bson.M{"$project": bson.M{
				"views": 0,
			}})
		}
	}

	if r.WithViews {
		project = append(project,
			bson.M{"$unwind": bson.M{"path": "$views", "preserveNullAndEmptyArrays": true}},
			bson.M{"$addFields": bson.M{
				"views.group": "$group",
			}},
		)
		project = append(project, s.authorProvider.PipelineForField("views.author")...)

		if r.WithTabs {
			project = append(project,
				bson.M{"$lookup": bson.M{
					"from":         mongo.ViewTabMongoCollection,
					"localField":   "views._id",
					"foreignField": "view",
					"as":           "tabs",
				}},
				bson.M{"$unwind": bson.M{"path": "$tabs", "preserveNullAndEmptyArrays": true}},
			)
			project = append(project, s.authorProvider.PipelineForField("tabs.author")...)

			if r.WithWidgets {
				project = append(project,
					bson.M{"$lookup": bson.M{
						"from":         mongo.WidgetMongoCollection,
						"localField":   "tabs._id",
						"foreignField": "tab",
						"as":           "widgets",
					}},
					bson.M{"$unwind": bson.M{"path": "$widgets", "preserveNullAndEmptyArrays": true}},
				)
				project = append(project, s.authorProvider.PipelineForField("widgets.author")...)
				project = append(project,
					bson.M{"$lookup": bson.M{
						"from": mongo.WidgetFiltersMongoCollection,
						"let":  bson.M{"widget": "$widgets._id"},
						"pipeline": []bson.M{
							{"$match": bson.M{
								"$expr":              bson.M{"$eq": bson.A{"$widget", "$$widget"}},
								"is_user_preference": false,
							}},
						},
						"as": "filters",
					}},
					bson.M{"$unwind": bson.M{"path": "$filters", "preserveNullAndEmptyArrays": true}},
				)
				project = append(project, s.authorProvider.PipelineForField("filters.author")...)
				project = append(project,
					bson.M{"$sort": bson.M{"filters.position": 1}},
					bson.M{"$group": bson.M{
						"_id": bson.M{
							"_id":    "$_id",
							"view":   "$views._id",
							"tab":    "$tabs._id",
							"widget": "$widgets._id",
						},
						"group":     bson.M{"$first": "$group"},
						"deletable": bson.M{"$first": "$deletable"},
						"views":     bson.M{"$first": "$views"},
						"tabs":      bson.M{"$first": "$tabs"},
						"widgets":   bson.M{"$first": "$widgets"},
						"filters": bson.M{"$push": bson.M{"$cond": bson.M{
							"if":   "$filters._id",
							"then": "$filters",
							"else": "$$REMOVE",
						}}},
					}},
					bson.M{"$addFields": bson.M{
						"_id":             "$_id._id",
						"widgets.filters": "$filters",
					}},
					bson.M{"$sort": bson.D{
						{Key: "widgets.grid_parameters.desktop.y", Value: 1},
						{Key: "widgets.grid_parameters.desktop.x", Value: 1},
					}},
					bson.M{"$group": bson.M{
						"_id": bson.M{
							"_id":  "$_id",
							"view": "$views._id",
							"tab":  "$tabs._id",
						},
						"group":     bson.M{"$first": "$group"},
						"deletable": bson.M{"$first": "$deletable"},
						"views":     bson.M{"$first": "$views"},
						"tabs":      bson.M{"$first": "$tabs"},
						"widgets":   bson.M{"$push": "$widgets"},
					}},
					bson.M{"$addFields": bson.M{
						"_id": "$_id._id",
						"tabs.widgets": bson.M{"$filter": bson.M{
							"input": "$widgets",
							"cond":  "$$this._id",
						}},
					}},
				)
			}

			project = append(project,
				bson.M{"$sort": bson.M{"tabs.position": 1}},
				bson.M{"$group": bson.M{
					"_id": bson.M{
						"_id":  "$_id",
						"view": "$views._id",
					},
					"group":     bson.M{"$first": "$group"},
					"deletable": bson.M{"$first": "$deletable"},
					"views":     bson.M{"$first": "$views"},
					"tabs":      bson.M{"$push": "$tabs"},
				}},
				bson.M{"$addFields": bson.M{
					"_id": "$_id._id",
					"views.tabs": bson.M{"$filter": bson.M{
						"input": "$tabs",
						"cond":  "$$this._id",
					}},
				}},
			)
		}

		project = append(project,
			bson.M{"$sort": bson.M{"views.position": 1}},
			bson.M{"$group": bson.M{
				"_id":       "$_id",
				"group":     bson.M{"$first": "$group"},
				"deletable": bson.M{"$first": "$deletable"},
				"views":     bson.M{"$push": "$views"},
			}},
			bson.M{"$replaceRoot": bson.M{
				"newRoot": bson.M{"$mergeObjects": bson.A{
					"$group",
					bson.M{
						"deletable": "$deletable",
						"views": bson.M{"$filter": bson.M{
							"input": "$views",
							"cond":  bson.M{"$in": bson.A{"$$this._id", authorizedViewIds}},
						}},
					},
				}},
			}},
			sort,
		)
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		sort,
		project,
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	res := AggregationResult{}

	if cursor.Next(ctx) {
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (s *store) GetOneBy(ctx context.Context, id string) (*ViewGroup, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		viewgroup := &ViewGroup{}
		err := cursor.Decode(viewgroup)
		if err != nil {
			return nil, err
		}

		return viewgroup, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r EditRequest) (*ViewGroup, error) {
	now := datetime.NewCpsTime()
	doc := view.Group{
		ID:      utils.NewID(),
		Title:   r.Title,
		Author:  r.Author,
		Created: now,
		Updated: now,
	}
	var response *ViewGroup
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		position, err := s.getNextPosition(ctx)
		if err != nil {
			return err
		}

		doc.Position = position

		_, err = s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}

		response, err = s.GetOneBy(ctx, doc.ID)
		return err
	})

	return response, err
}

func (s *store) Update(ctx context.Context, r EditRequest) (*ViewGroup, error) {
	now := datetime.NewCpsTime()

	var response *ViewGroup
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": bson.M{
			"title":   r.Title,
			"author":  r.Author,
			"updated": now,
		}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		response, err = s.GetOneBy(ctx, r.ID)
		return err
	})

	return response, err
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	res := false
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = false
		err := s.dbViewCollection.FindOne(ctx, bson.M{"group_id": id}).Err()
		if err != nil {
			if !errors.Is(err, mongodriver.ErrNoDocuments) {
				return err
			}
		} else {
			return ErrLinkedToView
		}

		delCount, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			return err
		}

		res = delCount > 0
		return nil
	})

	return res, err
}

func (s *store) getNextPosition(ctx context.Context) (int64, error) {
	cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
		{"$group": bson.M{
			"_id":      nil,
			"position": bson.M{"$max": "$position"},
		}},
	})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		data := struct {
			Position int64 `bson:"position"`
		}{}
		err = cursor.Decode(&data)
		return data.Position + 1, err
	}

	return 0, nil
}
