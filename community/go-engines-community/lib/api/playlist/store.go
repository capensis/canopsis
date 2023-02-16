package playlist

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

const permissionPrefix = "Rights on playlist :"

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetById(ctx context.Context, id string) (*Response, error)
	Insert(ctx context.Context, userID string, r EditRequest) (*Response, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		client:                dbClient,
		collection:            dbClient.Collection(mongo.PlaylistMongoCollection),
		aclCollection:         dbClient.Collection(mongo.RightsMongoCollection),
		defaultSearchByFields: []string{"_id", "name"},
		defaultSortBy:         "name",
	}
}

type store struct {
	client                mongo.DbClient
	collection            mongo.DbCollection
	aclCollection         mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline := author.Pipeline()

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

func (s *store) GetById(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, author.Pipeline()...)

	cursor, err := s.collection.Aggregate(ctx, pipeline)
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

func (s *store) Insert(ctx context.Context, userID string, r EditRequest) (*Response, error) {
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

	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err := s.collection.InsertOne(ctx, model)
		if err != nil {
			return err
		}

		err = s.createPermission(ctx, userID, id, r.Name)
		if err != nil {
			return err
		}

		response, err = s.GetById(ctx, id)
		return err
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *store) Update(ctx context.Context, r EditRequest) (*Response, error) {
	now := types.CpsTime{Time: time.Now()}
	model := Playlist{
		Author:     r.Author,
		Name:       r.Name,
		Enabled:    *r.Enabled,
		Fullscreen: *r.Fullscreen,
		TabsList:   r.TabsList,
		Interval:   r.Interval,
		Updated:    now,
	}

	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		res, err := s.collection.UpdateOne(
			ctx,
			bson.M{"_id": r.ID},
			bson.M{"$set": model},
		)
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		model.ID = r.ID
		err = s.updatePermission(ctx, r.ID, r.Name)
		if err != nil {
			return err
		}

		response, err = s.GetById(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted := false
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = false
		d, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil || d == 0 {
			return err
		}

		err = s.deletePermission(ctx, id)
		if err != nil {
			return err
		}

		deleted = true
		return nil
	})

	return deleted, err
}

func (s *store) createPermission(ctx context.Context, userID, playlistID, playlistName string) error {
	_, err := s.aclCollection.InsertOne(ctx, bson.M{
		"_id":          playlistID,
		"crecord_name": playlistID,
		"crecord_type": securitymodel.LineTypeObject,
		"description":  fmt.Sprintf("%s %s", permissionPrefix, playlistName),
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
				"description": fmt.Sprintf("%s %s", permissionPrefix, playlistName),
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
