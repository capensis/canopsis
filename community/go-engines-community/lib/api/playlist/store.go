package playlist

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const permissionPrefix = "Rights on playlist :"

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetById(ctx context.Context, id string) (*Playlist, error)
	Insert(ctx context.Context, userID string, r EditRequest) (*Playlist, error)
	Update(ctx context.Context, r EditRequest) (*Playlist, error)
	Delete(ctx context.Context, id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		collection:            dbClient.Collection(mongo.PlaylistMongoCollection),
		aclCollection:         dbClient.Collection(mongo.RightsMongoCollection),
		defaultSearchByFields: []string{"_id", "name", "author"},
		defaultSortBy:         "name",
	}
}

type store struct {
	collection            mongo.DbCollection
	aclCollection         mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)

	if len(r.Ids) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"_id": bson.M{"$in": r.Ids}}})
	}

	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := s.defaultSortBy
	if r.SortBy != "" {
		sortBy = r.SortBy
	}
	if sortBy == "interval" {
		sortBy = "interval.value"
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

func (s *store) GetById(ctx context.Context, id string) (*Playlist, error) {
	playlist := &Playlist{}

	if err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&playlist); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return playlist, nil
}

func (s *store) Insert(ctx context.Context, userID string, r EditRequest) (*Playlist, error) {
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

	_, err := s.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	err = s.createPermission(ctx, userID, id, r.Name)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (s *store) Update(ctx context.Context, r EditRequest) (*Playlist, error) {
	prevModel, err := s.GetById(ctx, r.ID)
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
	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"_id": r.ID},
		bson.M{"$set": model},
	)
	if err != nil {
		return nil, err
	}

	model.ID = r.ID

	err = s.updatePermission(ctx, r.ID, r.Name)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	if deleted == 0 {
		return false, nil
	}

	err = s.deletePermission(ctx, id)
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s *store) createPermission(ctx context.Context, userID, playlistID, playlistName string) error {
	_, err := s.aclCollection.InsertOne(ctx, bson.M{
		"_id":          playlistID,
		"crecord_name": playlistID,
		"crecord_type": securitymodel.LineTypeObject,
		"desc":         fmt.Sprintf("%s %s", permissionPrefix, playlistName),
		"type":         securitymodel.LineObjectTypeRW,
	})
	if err != nil {
		return err
	}

	res := s.aclCollection.FindOne(ctx, bson.M{
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

	_, err = s.aclCollection.UpdateMany(ctx,
		bson.M{
			"_id":          bson.M{"$in": bson.A{user.Role, security.RoleAdmin}},
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

func (s *store) updatePermission(ctx context.Context, playlistID, playlistName string) error {

	_, err := s.aclCollection.UpdateOne(ctx,
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

func (s *store) deletePermission(ctx context.Context, playlistID string) error {
	_, err := s.aclCollection.UpdateMany(ctx,
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

	_, err = s.aclCollection.DeleteOne(ctx, bson.M{
		"_id":          playlistID,
		"crecord_type": securitymodel.LineTypeObject,
	})

	return err
}
