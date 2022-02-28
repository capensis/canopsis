package viewgroup

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Store interface {
	Find(ctx context.Context, r ListRequest, authorizedViewIds []string) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*ViewGroup, error)
	Insert(ctx context.Context, r EditRequest) (*ViewGroup, error)
	Update(ctx context.Context, r EditRequest) (*ViewGroup, error)
	Delete(ctx context.Context, id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbCollection:          dbClient.Collection(mongo.ViewGroupMongoCollection),
		dbViewCollection:      dbClient.Collection(mongo.ViewMongoCollection),
		defaultSearchByFields: []string{"_id", "title", "author"},
		defaultSortBy:         "position",
	}
}

type store struct {
	dbCollection          mongo.DbCollection
	dbViewCollection      mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
}

func (s *store) Find(ctx context.Context, r ListRequest, authorizedViewIds []string) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sort := common.GetSortQuery(s.defaultSortBy, common.SortAsc)
	project := make([]bson.M, 0)

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

			if r.WithWidgets {
				project = append(project,
					bson.M{"$lookup": bson.M{
						"from":         mongo.WidgetMongoCollection,
						"localField":   "tabs._id",
						"foreignField": "tab",
						"as":           "widgets",
					}},
					bson.M{"$unwind": bson.M{"path": "$widgets", "preserveNullAndEmptyArrays": true}},
					bson.M{"$lookup": bson.M{
						"from":         mongo.WidgetFiltersMongoCollection,
						"localField":   "widgets._id",
						"foreignField": "widget",
						"as":           "filters",
					}},
					bson.M{"$unwind": bson.M{"path": "$filters", "preserveNullAndEmptyArrays": true}},
					bson.M{"$sort": bson.M{"filters.title": 1}},
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
						"filters":   bson.M{"$push": "$filters"},
					}},
					bson.M{"$addFields": bson.M{
						"_id": "$_id._id",
						"widgets.filters": bson.M{"$filter": bson.M{
							"input": "$filters",
							"cond":  bson.M{"$eq": bson.A{"$$this.is_private", false}},
						}},
					}},
					bson.M{"$sort": bson.D{{"widgets.grid_parameters.desktop.y", 1}, {"widgets.grid_parameters.desktop.x", 1}}},
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
	count, err := s.dbCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	now := types.CpsTime{Time: time.Now()}
	group := ViewGroup{
		ID:      utils.NewID(),
		Title:   r.Title,
		Author:  r.Author,
		Created: &now,
		Updated: &now,
	}

	_, err = s.dbCollection.InsertOne(ctx, view.Group{
		ID:       group.ID,
		Title:    group.Title,
		Author:   group.Author,
		Position: count,
		Created:  *group.Created,
		Updated:  *group.Updated,
	})
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (s *store) Update(ctx context.Context, r EditRequest) (*ViewGroup, error) {
	group, err := s.GetOneBy(ctx, r.ID)
	if err != nil || group == nil {
		return nil, err
	}

	now := types.CpsTime{Time: time.Now()}
	group.Title = r.Title
	group.Author = r.Author
	group.Updated = &now

	_, err = s.dbCollection.UpdateOne(ctx, bson.M{"_id": group.ID}, bson.M{"$set": bson.M{
		"title":   group.Title,
		"author":  group.Author,
		"updated": group.Updated,
	}})
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	res := s.dbViewCollection.FindOne(ctx, bson.M{"group_id": id})
	if err := res.Err(); err != nil {
		if err != mongodriver.ErrNoDocuments {
			return false, err
		}
	} else {
		return false, ErrLinkedToView
	}

	delCount, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return delCount > 0, nil
}
