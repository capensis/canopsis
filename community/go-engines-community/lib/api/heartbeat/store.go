package heartbeat

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Store is an interface for heartbeats storage
type Store interface {
	Insert(model []*Heartbeat) error
	Find(r FilteredQuery) (*AggregationResult, error)
	GetOneBy(id string) (*Heartbeat, error)
	Update(model []*Heartbeat) error
	Delete(ids []string) error
}

type store struct {
	db            mongo.DbClient
	collection    mongo.DbCollection
	defaultSortBy string
}

// NewStore instantiates heartbeat store.
func NewStore(db mongo.DbClient) Store {
	return &store{
		db:            db,
		collection:    db.Collection(mongo.HeartbeatMongoCollection),
		defaultSortBy: "created",
	}
}

// Find heartbeats according to query.
func (s *store) Find(r FilteredQuery) (*AggregationResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := bson.M{}

	if r.Search != "" {
		searchRegexp := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", r.Search),
			Options: "i",
		}

		filter["$or"] = []bson.M{
			{"name": searchRegexp},
			{"description": searchRegexp},
			{"author": searchRegexp},
		}
	}

	sortBy := r.SortBy
	if sortBy == "" {
		sortBy = s.defaultSortBy
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		[]bson.M{{"$match": filter}},
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
func (s *store) GetOneBy(id string) (*Heartbeat, error) {
	res := &Heartbeat{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(res); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return res, nil
}

// Create new heartbeats.
func (s *store) Insert(heartbeats []*Heartbeat) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
func (s *store) Update(heartbeats []*Heartbeat) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	findHeartbeats := make(map[string]*Heartbeat)
	for cursor.Next(ctx) {
		h := Heartbeat{}
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
func (s *store) Delete(ids []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
