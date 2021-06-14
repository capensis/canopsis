package viewgroup

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Store interface {
	Find(r ListRequest, authorizedViewIds []string) (*AggregationResult, error)
	GetOneBy(string) (*ViewGroup, error)
	Insert([]EditRequest) ([]ViewGroup, error)
	Update([]BulkUpdateRequestItem) ([]ViewGroup, error)
	Delete([]string) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbCollection:     dbClient.Collection(mongo.ViewGroupMongoCollection),
		dbViewCollection: dbClient.Collection(mongo.ViewMongoCollection),
	}
}

type store struct {
	dbCollection     mongo.DbCollection
	dbViewCollection mongo.DbCollection
}

func (s *store) Find(r ListRequest, authorizedViewIds []string) (*AggregationResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := bson.M{}

	if r.Search != "" {
		searchRegexp := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", r.Search),
			Options: "i",
		}

		filter["$or"] = []bson.M{
			{"title": searchRegexp},
			{"description": searchRegexp},
		}
	}

	sort := common.GetSortQuery("position", common.SortAsc)
	pipeline := []bson.M{{"$match": filter}}
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
			bson.M{"$sort": bson.M{"views.position": 1}},
			bson.M{"$group": bson.M{
				"_id":       "$_id",
				"group":     bson.M{"$first": "$group"},
				"deletable": bson.M{"$first": "$deletable"},
				"views":     bson.M{"$push": "$views"},
			}},
			bson.M{"$addFields": bson.M{
				"views.group": "$group",
			}},
			bson.M{"$replaceRoot": bson.M{
				"newRoot": bson.M{"$mergeObjects": bson.A{
					"$group",
					bson.M{
						"deletable": "$deletable",
						"views": bson.M{"$filter": bson.M{
							"input": bson.M{"$cond": bson.M{
								"if":   "$views",
								"then": "$views",
								"else": bson.A{},
							}},
							"cond": bson.M{"$in": bson.A{"$$this._id", authorizedViewIds}},
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

func (s *store) GetOneBy(id string) (*ViewGroup, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func (s *store) Insert(r []EditRequest) ([]ViewGroup, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	count, err := s.dbCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	now := types.CpsTime{Time: time.Now()}
	models := make([]interface{}, len(r))
	groups := make([]ViewGroup, len(r))

	for i, item := range r {
		groups[i] = ViewGroup{
			ID:      utils.NewID(),
			Title:   item.Title,
			Author:  item.Author,
			Created: now,
			Updated: now,
		}
		models[i] = view.Group{
			ID:       groups[i].ID,
			Title:    groups[i].Title,
			Author:   groups[i].Author,
			Position: count + int64(i),
			Created:  groups[i].Created,
			Updated:  groups[i].Updated,
		}
	}

	_, err = s.dbCollection.InsertMany(ctx, models)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (s *store) Update(r []BulkUpdateRequestItem) ([]ViewGroup, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ids := make([]string, len(r))
	rByID := make(map[string]BulkUpdateRequestItem, len(r))
	for i, item := range r {
		ids[i] = item.ID
		rByID[item.ID] = item
	}

	groups, err := s.findByIDs(ctx, ids)
	if err != nil || len(groups) < len(ids) {
		return nil, err
	}

	models := make([]mongodriver.WriteModel, len(groups))
	now := types.CpsTime{Time: time.Now()}

	for i := range groups {
		groups[i].Title = rByID[groups[i].ID].Title
		groups[i].Author = rByID[groups[i].ID].Author
		groups[i].Updated = now

		models[i] = mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": groups[i].ID}).
			SetUpdate(bson.M{"$set": bson.M{
				"title":   groups[i].Title,
				"author":  groups[i].Author,
				"updated": groups[i].Updated,
			}})
	}

	_, err = s.dbCollection.BulkWrite(ctx, models)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (s *store) findByIDs(ctx context.Context, ids []string) ([]ViewGroup, error) {
	cursor, err := s.dbCollection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}},
		options.Find().SetSort(bson.D{
			{Key: "position", Value: 1},
			{Key: "_id", Value: 1},
		}))
	if err != nil {
		return nil, err
	}

	groups := make([]ViewGroup, 0)
	err = cursor.All(ctx, &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (s *store) Delete(ids []string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res := s.dbViewCollection.FindOne(ctx, bson.M{"group_id": bson.M{"$in": ids}})
	if err := res.Err(); err != nil {
		if err != mongodriver.ErrNoDocuments {
			return false, err
		}
	} else {
		return false, ErrLinkedToView
	}

	if len(ids) > 0 {
		count, err := s.dbCollection.CountDocuments(ctx, bson.M{"_id": bson.M{"$in": ids}})
		if err != nil {
			return false, err
		}

		if count < int64(len(ids)) {
			return false, nil
		}
	}

	delCount, err := s.dbCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return false, err
	}

	return delCount > 0, nil
}
