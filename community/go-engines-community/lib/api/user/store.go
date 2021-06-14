package user

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	securitymodel "git.canopsis.net/canopsis/go-engines/lib/security/model"
	"git.canopsis.net/canopsis/go-engines/lib/security/password"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store interface {
	Find(ListRequest) (*AggregationResult, error)
	GetOneBy(string) (*User, error)
	Insert(EditRequest) (*User, error)
	Update(EditRequest) (*User, error)
	Delete(string) (bool, error)
}

func NewStore(dbClient mongo.DbClient, passwordEncoder password.Encoder) Store {
	return &store{
		dbClient:        dbClient,
		dbCollection:    dbClient.Collection(mongo.RightsMongoCollection),
		passwordEncoder: passwordEncoder,
	}
}

type store struct {
	dbClient        mongo.DbClient
	dbCollection    mongo.DbCollection
	passwordEncoder password.Encoder
}

func (s *store) Find(r ListRequest) (*AggregationResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := bson.M{}

	if r.Search != "" {
		searchRegexp := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", r.Search),
			Options: "i",
		}

		filter["$or"] = []bson.M{
			{"crecord_name": searchRegexp},
			{"firstname": searchRegexp},
			{"lastname": searchRegexp},
		}
	}

	sortBy := "name"
	if r.SortBy != "" {
		sortBy = r.SortBy
	}

	pipeline := []bson.M{
		{"$match": bson.M{"crecord_type": securitymodel.LineTypeSubject}},
		{"$match": filter},
	}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
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

func (s *store) GetOneBy(id string) (*User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pipeline := []bson.M{
		{"$match": bson.M{
			"_id":          id,
			"crecord_type": securitymodel.LineTypeSubject,
		}},
	}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
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

func (s *store) Insert(r EditRequest) (*User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := s.dbCollection.InsertOne(ctx, bson.M{
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

	return s.GetOneBy(r.Name)
}

func (s *store) Update(r EditRequest) (*User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	res, err := s.dbCollection.UpdateOne(ctx,
		bson.M{"_id": r.ID, "crecord_type": securitymodel.LineTypeSubject},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, err
	}

	if res.MatchedCount == 0 {
		return nil, nil
	}

	return s.GetOneBy(r.ID)
}

func (s *store) Delete(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	delCount, err := s.dbCollection.DeleteOne(ctx, bson.M{
		"_id":          id,
		"crecord_type": securitymodel.LineTypeSubject,
	})
	if err != nil {
		return false, err
	}

	return delCount > 0, nil
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
				"_id":  "$role._id",
				"name": "$role.crecord_name",
			},
		}},
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "defaultview",
			"foreignField": "_id",
			"as":           "defaultview",
		}},
		{"$unwind": bson.M{"path": "$defaultview", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"name":                      "$crecord_name",
			"email":                     "$mail",
			"ui_groups_navigation_type": "$groupsNavigationType",
			"ui_tours":                  "$tours",
		}},
	}
}
