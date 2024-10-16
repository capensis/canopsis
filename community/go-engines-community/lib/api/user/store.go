package user

import (
	"cmp"
	"context"
	"slices"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Find(ctx context.Context, r ListRequest, userID string) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*User, error)
	Insert(ctx context.Context, r CreateRequest) (*User, error)
	Update(ctx context.Context, r UpdateRequest, userID string) (*User, error)
	Patch(ctx context.Context, r PatchRequest, userID string) (*User, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
}

func NewStore(
	dbClient mongo.DbClient,
	passwordEncoder password.Encoder,
	websocketStore websocket.Store,
	authorProvider author.Provider,
) Store {
	return &store{
		client:                 dbClient,
		collection:             dbClient.Collection(mongo.UserCollection),
		userPrefCollection:     dbClient.Collection(mongo.UserPreferencesMongoCollection),
		patternCollection:      dbClient.Collection(mongo.PatternMongoCollection),
		viewGroupsCollection:   dbClient.Collection(mongo.ViewGroupMongoCollection),
		viewCollection:         dbClient.Collection(mongo.ViewMongoCollection),
		viewTabCollection:      dbClient.Collection(mongo.ViewTabMongoCollection),
		widgetCollection:       dbClient.Collection(mongo.WidgetMongoCollection),
		widgetFilterCollection: dbClient.Collection(mongo.WidgetFiltersMongoCollection),
		shareTokenCollection:   dbClient.Collection(mongo.ShareTokenMongoCollection),

		passwordEncoder: passwordEncoder,
		websocketStore:  websocketStore,
		authorProvider:  authorProvider,

		defaultSearchByFields: []string{"_id", "name", "firstname", "lastname", "roles.name"},
		defaultSortBy:         "name",
	}
}

type store struct {
	client                 mongo.DbClient
	collection             mongo.DbCollection
	userPrefCollection     mongo.DbCollection
	patternCollection      mongo.DbCollection
	viewGroupsCollection   mongo.DbCollection
	viewCollection         mongo.DbCollection
	viewTabCollection      mongo.DbCollection
	widgetCollection       mongo.DbCollection
	widgetFilterCollection mongo.DbCollection
	shareTokenCollection   mongo.DbCollection

	passwordEncoder password.Encoder
	websocketStore  websocket.Store
	authorProvider  author.Provider

	defaultSearchByFields []string
	defaultSortBy         string
}

func (s *store) Find(ctx context.Context, r ListRequest, curUserID string) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	project := []bson.M{
		{"$addFields": bson.M{
			"username": "$name",
		}},
		{"$addFields": bson.M{
			"display_name": s.authorProvider.GetDisplayNameQuery(""),
		}},
	}
	project = append(project, s.authorProvider.Pipeline()...)

	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 || r.Permission != "" {
		pipeline = append(pipeline, getRolePipeline()...)
	} else {
		project = append(project, getRolePipeline()...)
	}

	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	if r.Permission != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"roles.permissions." + r.Permission: bson.M{"$exists": true}}})
	}

	project = append(project, getViewPipeline()...)
	project = append(project, getUiThemePipeline(s.authorProvider)...)

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		common.GetSortQuery(cmp.Or(r.SortBy, s.defaultSortBy), r.Sort),
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

	ids := make([]string, len(res.Data))
	for i, v := range res.Data {
		ids[i] = v.ID
	}

	conns, err := s.websocketStore.GetConnections(ctx, ids)
	if err != nil {
		return nil, err
	}

	var onlyOneAdmin bool
	var lastAdminID string
	if r.WithFlags {
		onlyOneAdmin, lastAdminID, _, err = s.checkLastAdmin(ctx)
		if err != nil {
			return nil, err
		}
	}

	for i := range res.Data {
		activeConnects := conns[res.Data[i].ID]
		res.Data[i].ActiveConnects = &activeConnects
		if r.WithFlags {
			deletable := res.Data[i].ID != curUserID && (!onlyOneAdmin || res.Data[i].ID != lastAdminID)
			res.Data[i].Deletable = &deletable
		}
	}

	return &res, nil
}

func (s *store) GetOneBy(ctx context.Context, id string) (*User, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
	}
	pipeline = append(pipeline, getNestedObjectsPipeline(s.authorProvider)...)
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

func (s *store) Insert(ctx context.Context, r CreateRequest) (*User, error) {
	insertDoc, err := r.getBson(s.passwordEncoder)
	if err != nil {
		return nil, err
	}

	var user *User
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		user = nil
		_, err = s.collection.InsertOne(ctx, insertDoc)
		if err != nil {
			return err
		}

		user, err = s.GetOneBy(ctx, r.Name)

		return err
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *store) Update(ctx context.Context, r UpdateRequest, curUserID string) (*User, error) {
	if r.ID == curUserID && r.IsEnabled != nil && !*r.IsEnabled {
		return nil, common.NewValidationError("enable", "user cannot disable itself")
	}

	updateDoc, err := r.getBson(s.passwordEncoder)
	if err != nil {
		return nil, err
	}

	var user *User
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		user = nil
		onlyOneAdmin, lastAdminID, adminRoleID, err := s.checkLastAdmin(ctx)
		if err != nil {
			return err
		}

		if onlyOneAdmin && lastAdminID == r.ID {
			if !slices.Contains(r.Roles, adminRoleID) {
				return common.NewValidationError("roles", "last admin cannot be edited")
			}

			if r.IsEnabled != nil && !*r.IsEnabled {
				return common.NewValidationError("enable", "last admin cannot be disabled")
			}
		}

		res, err := s.collection.UpdateOne(ctx,
			bson.M{"_id": r.ID},
			bson.M{"$set": updateDoc},
		)
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		user, err = s.GetOneBy(ctx, r.ID)

		return err
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *store) Patch(ctx context.Context, r PatchRequest, curUserID string) (*User, error) {
	if r.ID == curUserID && r.IsEnabled != nil && !*r.IsEnabled {
		return nil, common.NewValidationError("enable", "user cannot disable itself")
	}

	updateDoc, err := r.getBson(s.passwordEncoder)
	if err != nil {
		return nil, err
	}

	var user *User
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		user = nil
		onlyOneAdmin, lastAdminID, adminRoleID, err := s.checkLastAdmin(ctx)
		if err != nil {
			return err
		}

		if onlyOneAdmin && lastAdminID == r.ID {
			if len(r.Roles) > 0 && !slices.Contains(r.Roles, adminRoleID) {
				return common.NewValidationError("roles", "last admin cannot be edited")
			}

			if r.IsEnabled != nil && !*r.IsEnabled {
				return common.NewValidationError("enable", "last admin cannot be disabled")
			}
		}

		res, err := s.collection.UpdateOne(ctx,
			bson.M{"_id": r.ID},
			bson.M{"$set": updateDoc},
		)
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		user, err = s.GetOneBy(ctx, r.ID)

		return err
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	if id == userID {
		return false, common.NewValidationError("_id", "user cannot delete itself")
	}

	var deleted int64
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		onlyOneAdmin, lastAdminID, _, err := s.checkLastAdmin(ctx)
		if err != nil {
			return err
		}

		if onlyOneAdmin && id == lastAdminID {
			return common.NewValidationError("_id", "last admin cannot be deleted")
		}

		// required to get the author in action log listener.
		res, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		deleted, err = s.collection.DeleteOne(ctx, bson.M{"_id": id})

		return err
	})

	if err != nil || deleted == 0 {
		return false, err
	}

	err = s.deleteUserPreferences(ctx, id)
	if err != nil {
		return false, err
	}

	err = s.deletePatterns(ctx, id)
	if err != nil {
		return false, err
	}

	err = s.deleteViewPrivateObjects(ctx, id)
	if err != nil {
		return false, err
	}

	err = s.deleteShareTokens(ctx, id)
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

func (s *store) deletePatterns(ctx context.Context, id string) error {
	_, err := s.patternCollection.DeleteMany(ctx, bson.M{
		"author":       id,
		"is_corporate": false,
	})

	return err
}

func (s *store) deleteViewPrivateObjects(ctx context.Context, id string) error {
	_, err := s.viewGroupsCollection.DeleteMany(ctx, bson.M{
		"author":     id,
		"is_private": true,
	})
	if err != nil {
		return err
	}

	_, err = s.viewCollection.DeleteMany(ctx, bson.M{
		"author":     id,
		"is_private": true,
	})
	if err != nil {
		return err
	}

	_, err = s.viewTabCollection.DeleteMany(ctx, bson.M{
		"author":     id,
		"is_private": true,
	})
	if err != nil {
		return err
	}

	_, err = s.widgetCollection.DeleteMany(ctx, bson.M{
		"author":     id,
		"is_private": true,
	})
	if err != nil {
		return err
	}

	_, err = s.widgetFilterCollection.DeleteMany(ctx, bson.M{
		"author": id,
		"$or": bson.A{
			bson.M{"is_user_preference": true},
			bson.M{"is_private": true},
		},
	})

	return err
}

func (s *store) deleteShareTokens(ctx context.Context, id string) error {
	_, err := s.shareTokenCollection.DeleteMany(ctx, bson.M{
		"user": id,
	})

	return err
}

func (s *store) checkLastAdmin(ctx context.Context) (bool, string, string, error) {
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"enable": true}},
		{"$lookup": bson.M{
			"from":         mongo.RoleCollection,
			"localField":   "roles",
			"foreignField": "_id",
			"as":           "admin_role",
			"pipeline": []bson.M{
				{"$match": bson.M{"name": security.RoleAdmin}},
			},
		}},
		{"$unwind": "$admin_role"},
		{"$group": bson.M{
			"_id":           nil,
			"count":         bson.M{"$sum": 1},
			"last_id":       bson.M{"$first": "$_id"},
			"admin_role_id": bson.M{"$first": "$admin_role._id"},
		}},
	})
	if err != nil {
		return false, "", "", err
	}

	defer cursor.Close(ctx)
	res := struct {
		Count       int64  `bson:"count"`
		LastID      string `bson:"last_id"`
		AdminRoleID string `bson:"admin_role_id"`
	}{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return false, "", "", err
		}
	}

	return res.Count <= 1, res.LastID, res.AdminRoleID, nil
}

func getNestedObjectsPipeline(authorProvider author.Provider) []bson.M {
	pipeline := []bson.M{
		{"$addFields": bson.M{
			"username": "$name",
		}},
		{"$addFields": bson.M{
			"display_name": authorProvider.GetDisplayNameQuery(""),
		}},
	}
	pipeline = append(pipeline, getRolePipeline()...)
	pipeline = append(pipeline, getViewPipeline()...)
	pipeline = append(pipeline, getUiThemePipeline(authorProvider)...)
	pipeline = append(pipeline, authorProvider.Pipeline()...)

	return pipeline
}

func getRolePipeline() []bson.M {
	return []bson.M{
		{"$unwind": bson.M{"path": "$roles", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.RoleCollection,
			"localField":   "roles",
			"foreignField": "_id",
			"as":           "roles",
		}},
		{"$unwind": bson.M{"path": "$roles", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "roles.defaultview",
			"foreignField": "_id",
			"as":           "roles.defaultview",
		}},
		{"$unwind": bson.M{"path": "$roles.defaultview", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"role_index": 1}},
		{"$group": bson.M{
			"_id":  "$_id",
			"data": bson.M{"$first": "$$ROOT"},
			"roles": bson.M{"$push": bson.M{
				"$cond": bson.M{
					"if":   "$roles._id",
					"then": "$roles",
					"else": "$$REMOVE",
				},
			}},
		}},
		{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$data",
			bson.M{"roles": "$roles"},
		}}}},
	}
}

func getViewPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "defaultview",
			"foreignField": "_id",
			"as":           "defaultview",
		}},
		{"$unwind": bson.M{"path": "$defaultview", "preserveNullAndEmptyArrays": true}},
	}
}

func getUiThemePipeline(authorProvider author.Provider) []bson.M {
	pipeline := []bson.M{
		{
			"$addFields": bson.M{
				"ui_theme": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$or": bson.A{
								bson.M{"$eq": bson.A{"$ui_theme", ""}},
								bson.M{"$eq": bson.A{bson.M{"$ifNull": bson.A{"$ui_theme", ""}}, ""}},
							},
						},
						"then": "canopsis",
						"else": "$ui_theme",
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         mongo.ColorThemeCollection,
				"localField":   "ui_theme",
				"foreignField": "_id",
				"as":           "ui_theme",
			},
		},
		{
			"$unwind": bson.M{"path": "$ui_theme", "preserveNullAndEmptyArrays": true},
		},
	}

	return append(pipeline, authorProvider.PipelineForField("ui_theme.author")...)
}
