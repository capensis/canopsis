package user

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*User, error)
	Insert(ctx context.Context, r EditRequest) (*User, error)
	Update(ctx context.Context, r EditRequest) (*User, error)
	Delete(ctx context.Context, id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient, passwordEncoder password.Encoder) Store {
	return &store{
		collection:            dbClient.Collection(mongo.RightsMongoCollection),
		userPrefCollection:    dbClient.Collection(mongo.UserPreferencesMongoCollection),
		passwordEncoder:       passwordEncoder,
		defaultSearchByFields: []string{"_id", "crecord_name", "firstname", "lastname"},
		defaultSortBy:         "name",
	}
}

type store struct {
	collection            mongo.DbCollection
	userPrefCollection    mongo.DbCollection
	passwordEncoder       password.Encoder
	defaultSearchByFields []string
	defaultSortBy         string
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"crecord_type": securitymodel.LineTypeSubject}},
		{"$addFields": bson.M{
			"name":  "$crecord_name",
			"email": "$mail",
		}},
	}

	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	if r.Permission != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{fmt.Sprintf("role.rights.%s", r.Permission): bson.M{"$exists": true}}})
	}

	sortBy := s.defaultSortBy
	if r.SortBy != "" {
		sortBy = r.SortBy
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

func (s *store) GetOneBy(ctx context.Context, id string) (*User, error) {
	pipeline := []bson.M{
		{"$match": bson.M{
			"_id":          id,
			"crecord_type": securitymodel.LineTypeSubject,
		}},
	}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		user := &User{}
		err := cursor.Decode(user)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r EditRequest) (*User, error) {
	_, err := s.collection.InsertOne(ctx, bson.M{
		"_id":                  r.Name,
		"crecord_name":         r.Name,
		"crecord_type":         securitymodel.LineTypeSubject,
		"lastname":             r.Lastname,
		"firstname":            r.Firstname,
		"mail":                 r.Email,
		"role":                 r.Role,
		"shadowpasswd":         string(s.passwordEncoder.EncodePassword([]byte(r.Password))),
		"ui_language":          r.UILanguage,
		"groupsNavigationType": r.UIGroupsNavigationType,
		"enable":               r.IsEnabled,
		"defaultview":          r.DefaultView,
		"authkey":              utils.NewID(),
	})
	if err != nil {
		return nil, err
	}

	return s.GetOneBy(ctx, r.Name)
}

func (s *store) Update(ctx context.Context, r EditRequest) (*User, error) {
	update := bson.M{
		"crecord_name":         r.Name,
		"lastname":             r.Lastname,
		"firstname":            r.Firstname,
		"mail":                 r.Email,
		"role":                 r.Role,
		"ui_language":          r.UILanguage,
		"groupsNavigationType": r.UIGroupsNavigationType,
		"enable":               r.IsEnabled,
		"defaultview":          r.DefaultView,
		"tours":                r.UITours,
	}
	if r.Password != "" {
		update["shadowpasswd"] = string(s.passwordEncoder.EncodePassword([]byte(r.Password)))
	}

	res, err := s.collection.UpdateOne(ctx,
		bson.M{"_id": r.ID, "crecord_type": securitymodel.LineTypeSubject},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, err
	}

	if res.MatchedCount == 0 {
		return nil, nil
	}

	return s.GetOneBy(ctx, r.ID)
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	delCount, err := s.collection.DeleteOne(ctx, bson.M{
		"_id":          id,
		"crecord_type": securitymodel.LineTypeSubject,
	})
	if err != nil {
		return false, err
	}

	if delCount == 0 {
		return false, nil
	}

	err = s.deleteUserPreferences(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) deleteUserPreferences(ctx context.Context, id string) error {
	_, err := s.userPrefCollection.DeleteMany(ctx, bson.M{
		"user": id,
	})

	return err
}

func getNestedObjectsPipeline() []bson.M {
	return []bson.M{
		{"$graphLookup": bson.M{
			"from":             mongo.RightsMongoCollection,
			"startWith":        "$role",
			"connectFromField": "role",
			"connectToField":   "_id",
			"as":               "role",
		}},
		{"$unwind": bson.M{"path": "$role", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"role": bson.M{
				"_id":         "$role._id",
				"name":        "$role.crecord_name",
				"rights":      "$role.rights",
				"defaultview": "$role.defaultview",
			},
		}},
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "defaultview",
			"foreignField": "_id",
			"as":           "defaultview",
		}},
		{"$unwind": bson.M{"path": "$defaultview", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "role.defaultview",
			"foreignField": "_id",
			"as":           "role.defaultview",
		}},
		{"$unwind": bson.M{"path": "$role.defaultview", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"name":                      "$crecord_name",
			"email":                     "$mail",
			"ui_groups_navigation_type": "$groupsNavigationType",
			"ui_tours":                  "$tours",
		}},
	}
}
