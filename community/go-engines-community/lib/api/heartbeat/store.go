package heartbeat

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Store is an interface for heartbeats storage
type Store interface {
	Insert(ctx context.Context, model []*Response) error
	Find(ctx context.Context, r FilteredQuery) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Response, error)
	Update(ctx context.Context, model []*Response) error
	Delete(ctx context.Context, ids []string) error
}

type store struct {
	db                    mongo.DbClient
	collection            mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
}

// NewStore instantiates heartbeat store.
func NewStore(db mongo.DbClient) Store {
	return &store{
		db:                    db,
		collection:            db.Collection(mongo.HeartbeatMongoCollection),
		defaultSearchByFields: []string{"_id", "name", "description", "author"},
		defaultSortBy:         "created",
	}
}

// Find heartbeats according to query.
func (s *store) Find(ctx context.Context, r FilteredQuery) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := r.SortBy
	if sortBy == "" {
		sortBy = s.defaultSortBy
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		common.GetSortQuery(sortBy, r.Sort),
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

// GetOneBy heartbeat by id.
func (s *store) GetOneBy(ctx context.Context, id string) (*Response, error) {
	res := &Response{}
	if err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(res); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return res, nil
}

// Create new heartbeats.
func (s *store) Insert(ctx context.Context, heartbeats []*Response) error {
	now := types.CpsTime{Time: time.Now()}

	for i := range heartbeats {
		if heartbeats[i].ID == "" {
			heartbeats[i].ID = utils.NewID()
		}

		heartbeats[i].Created = &now
		heartbeats[i].Updated = &now
	}

	docs := make([]interface{}, len(heartbeats))
	for i := range docs {
		docs[i] = heartbeats[i]
	}

	_, err := s.collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	return nil
}

// Update heartbeats.
func (s *store) Update(ctx context.Context, heartbeats []*Response) error {
	ids := make([]string, len(heartbeats))
	for i := range heartbeats {
		ids[i] = heartbeats[i].ID
	}

	cursor, err := s.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}},
		options.Find().SetProjection(bson.M{"_id": 1, "created": 1}))
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	findHeartbeats := make(map[string]*Response)
	for cursor.Next(ctx) {
		h := Response{}
		err := cursor.Decode(&h)
		if err != nil {
			return err
		}

		findHeartbeats[h.ID] = &h
	}

	if len(findHeartbeats) < len(ids) {
		notFoundIDs := make([]string, 0)
		for _, id := range ids {
			if _, ok := findHeartbeats[id]; !ok {
				notFoundIDs = append(notFoundIDs, id)
			}
		}

		return NotFoundError(fmt.Errorf("not found %s", strings.Join(notFoundIDs, ",")))
	}

	now := types.CpsTime{Time: time.Now()}
	models := make([]mongodriver.WriteModel, len(heartbeats))
	for i := range models {
		heartbeats[i].Updated = &now
		models[i] = mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": heartbeats[i].ID}).
			SetUpdate(bson.M{"$set": heartbeats[i]})
	}

	_, err = s.collection.BulkWrite(ctx, models)
	if err != nil {
		return err
	}

	for i := range heartbeats {
		heartbeats[i].Created = findHeartbeats[heartbeats[i].ID].Created
	}

	return nil
}

// Delete heartbeats by id
func (s *store) Delete(ctx context.Context, ids []string) error {
	cursor, err := s.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}},
		options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	findIDs := make(map[string]interface{})
	for cursor.Next(ctx) {
		id := struct {
			ID string `bson:"_id"`
		}{}
		err := cursor.Decode(&id)
		if err != nil {
			return err
		}

		findIDs[id.ID] = nil
	}

	if len(findIDs) < len(ids) {
		notFoundIDs := make([]string, 0)
		for _, id := range ids {
			if _, ok := findIDs[id]; !ok {
				notFoundIDs = append(notFoundIDs, id)
			}
		}

		return NotFoundError(fmt.Errorf("not found %s", strings.Join(notFoundIDs, ",")))
	}

	_, err = s.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})

	return err
}
