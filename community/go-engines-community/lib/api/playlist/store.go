package playlist

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	securitymodel "git.canopsis.net/canopsis/go-engines/lib/security/model"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"time"
)

const permissionPrefix = "Rights on playlist :"

type Store interface {
	Find(r ListRequest) (*AggregationResult, error)
	GetById(id string) (*Playlist, error)
	Insert(userID string, r EditRequest) (*Playlist, error)
	Update(r EditRequest) (*Playlist, error)
	Delete(id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbCollection:    dbClient.Collection(mongo.PlaylistMongoCollection),
		aclDbCollection: dbClient.Collection(mongo.RightsMongoCollection),
	}
}

type store struct {
	dbCollection    mongo.DbCollection
	aclDbCollection mongo.DbCollection
}

func (s *store) Find(r ListRequest) (*AggregationResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := bson.M{}

	if len(r.Ids) > 0 {
		filter["_id"] = bson.M{"$in": r.Ids}
	}

	if r.Search != "" {
		searchRegexp := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", r.Search),
			Options: "i",
		}

		filter["$or"] = []bson.M{
			{"name": searchRegexp},
		}
	}

	sortBy := "name"
	if r.SortBy != "" {
		sortBy = r.SortBy
	}
	if sortBy == "interval" {
		sortBy = "interval.seconds"
	}

	pipeline := []bson.M{{"$match": filter}}
	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
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

func (s *store) GetById(id string) (*Playlist, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playlist := &Playlist{}

	if err := s.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&playlist); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return playlist, nil
}

func (s *store) Insert(userID string, r EditRequest) (*Playlist, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := utils.NewID()
	now := types.CpsTime{Time: time.Now()}
	model := Playlist{
		ID:         id,
		Author:     r.Author,
		Name:       r.Name,
		Enabled:    *r.Enabled,
		Fullscreen: *r.Fullscreen,
		TabsList:   r.TabsList,
		Interval:   r.Interval,
		Created:    now,
		Updated:    now,
	}

	_, err := s.dbCollection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	err = s.createPermission(userID, id, r.Name)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (s *store) Update(r EditRequest) (*Playlist, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	prevModel, err := s.GetById(r.ID)
	if err != nil || prevModel == nil {
		return nil, err
	}

	now := types.CpsTime{Time: time.Now()}
	model := Playlist{
		Author:     r.Author,
		Name:       r.Name,
		Enabled:    *r.Enabled,
		Fullscreen: *r.Fullscreen,
		TabsList:   r.TabsList,
		Interval:   r.Interval,
		Created:    prevModel.Created,
		Updated:    now,
	}
	_, err = s.dbCollection.UpdateOne(
		ctx,
		bson.M{"_id": r.ID},
		bson.M{"$set": model},
	)
	if err != nil {
		return nil, err
	}

	model.ID = r.ID

	err = s.updatePermission(r.ID, r.Name)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (s *store) Delete(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	if deleted == 0 {
		return false, nil
	}

	err = s.deletePermission(id)
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s *store) createPermission(userID, playlistID, playlistName string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := s.aclDbCollection.InsertOne(ctx, bson.M{
		"_id":          playlistID,
		"crecord_name": playlistID,
		"crecord_type": securitymodel.LineTypeObject,
		"desc":         fmt.Sprintf("%s %s", permissionPrefix, playlistName),
		"type":         securitymodel.LineObjectTypeRW,
	})
	if err != nil {
		return err
	}

	res := s.aclDbCollection.FindOne(ctx, bson.M{
		"_id":          userID,
		"crecord_type": securitymodel.LineTypeSubject,
	})
	if err := res.Err(); err != nil {
		return err
	}

	user := struct {
		Role string `json:"role"`
	}{}
	err = res.Decode(&user)
	if err != nil {
		return err
	}

	_, err = s.aclDbCollection.UpdateOne(ctx,
		bson.M{
			"_id":          user.Role,
			"crecord_type": securitymodel.LineTypeRole,
		},
		bson.M{
			"$set": bson.M{
				"rights." + playlistID: bson.M{
					"checksum": securitymodel.PermissionBitmaskRead |
						securitymodel.PermissionBitmaskUpdate |
						securitymodel.PermissionBitmaskDelete,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return err
}

func (s *store) updatePermission(playlistID, playlistName string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := s.aclDbCollection.UpdateOne(ctx,
		bson.M{
			"_id":          playlistID,
			"crecord_type": securitymodel.LineTypeObject,
		},
		bson.M{
			"$set": bson.M{
				"desc": fmt.Sprintf("%s %s", permissionPrefix, playlistName),
			},
		},
	)

	return err
}

func (s *store) deletePermission(playlistID string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := s.aclDbCollection.UpdateMany(ctx,
		bson.M{
			"crecord_type":         securitymodel.LineTypeRole,
			"rights." + playlistID: bson.M{"$exists": true},
		},
		bson.M{
			"$unset": bson.M{"rights." + playlistID: ""},
		},
	)
	if err != nil {
		return err
	}

	_, err = s.aclDbCollection.DeleteOne(ctx, bson.M{
		"_id":          playlistID,
		"crecord_type": securitymodel.LineTypeObject,
	})

	return err
}
