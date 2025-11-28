package user

import (
	"cmp"
	"context"
	"errors"
	"slices"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

type Store interface {
	Find(ctx context.Context, r ListRequest, userID string, requestRoles []string) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*User, error)
	Insert(ctx context.Context, r CreateRequest, requestRoles []string) (*User, error)
	Update(ctx context.Context, r BulkUpdateRequestItem, userID string, requestRoles []string) (*User, error)
	Patch(ctx context.Context, r BulkPatchRequestItem, userID string, requestRoles []string) (*User, error)
	Delete(ctx context.Context, r BulkDeleteRequestItem, userID string, requestRoles []string) (bool, error)
}

func NewStore(
	dbClient mongo.DbClient,
	passwordEncoder password.Encoder,
	websocketStore websocket.Store,
	authorProvider author.Provider,
	securityConfig security.Config,
) Store {
	return &store{
		client:                 dbClient,
		userCollection:         dbClient.Collection(mongo.UserCollection),
		userPrefCollection:     dbClient.Collection(mongo.UserPreferencesMongoCollection),
		patternCollection:      dbClient.Collection(mongo.PatternMongoCollection),
		viewGroupsCollection:   dbClient.Collection(mongo.ViewGroupMongoCollection),
		viewCollection:         dbClient.Collection(mongo.ViewMongoCollection),
		viewTabCollection:      dbClient.Collection(mongo.ViewTabMongoCollection),
		widgetCollection:       dbClient.Collection(mongo.WidgetMongoCollection),
		widgetFilterCollection: dbClient.Collection(mongo.WidgetFiltersMongoCollection),
		shareTokenCollection:   dbClient.Collection(mongo.ShareTokenMongoCollection),
		notificationCollection: dbClient.Collection(mongo.UserNotificationCollection),
		tplTestCollection:      dbClient.Collection(mongo.TemplateTestCollection),
		colorThemeCollection:   dbClient.Collection(mongo.ColorThemeCollection),
		roleCollection:         dbClient.Collection(mongo.RoleCollection),

		passwordEncoder: passwordEncoder,
		websocketStore:  websocketStore,
		authorProvider:  authorProvider,
		securityConfig:  securityConfig,

		defaultSearchByFields: []string{"_id", "name", "firstname", "lastname", "roles.name"},
		defaultSortBy:         "name",
		dupErrorParser:        validation.NewDuplicateErrorParser(),
	}
}

type store struct {
	client                 mongo.DbClient
	userCollection         mongo.DbCollection
	userPrefCollection     mongo.DbCollection
	patternCollection      mongo.DbCollection
	viewGroupsCollection   mongo.DbCollection
	viewCollection         mongo.DbCollection
	viewTabCollection      mongo.DbCollection
	widgetCollection       mongo.DbCollection
	widgetFilterCollection mongo.DbCollection
	shareTokenCollection   mongo.DbCollection
	notificationCollection mongo.DbCollection
	tplTestCollection      mongo.DbCollection
	colorThemeCollection   mongo.DbCollection
	roleCollection         mongo.DbCollection

	passwordEncoder password.Encoder
	websocketStore  websocket.Store
	authorProvider  author.Provider
	securityConfig  security.Config

	defaultSearchByFields []string
	defaultSortBy         string
	dupErrorParser        validation.DuplicateErrorParser
}

func (s *store) Find(ctx context.Context, r ListRequest, curUserID string, requestRoles []string) (*AggregationResult, error) {
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

	filter := mongoquery.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 || r.Permission != "" {
		pipeline = append(pipeline, getRolePipeline()...)
	} else {
		project = append(project, getRolePipeline()...)
	}

	if len(r.IDs) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"_id": bson.M{"$in": r.IDs}}})
	}

	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	if r.Permission != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"roles.permissions." + r.Permission: bson.M{"$exists": true}}})
	}

	project = append(project, getViewPipeline()...)

	sort := mongoquery.GetSortQuery(cmp.Or(r.SortBy, s.defaultSortBy), r.Sort)
	project = append(project, sort)
	cursor, err := s.userCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
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

	ids := make([]string, len(res.Data))
	for i, v := range res.Data {
		ids[i] = v.ID
	}

	conns, err := s.websocketStore.GetConnections(ctx, ids)
	if err != nil {
		return nil, err
	}

	isAdminRequest := slices.Contains(requestRoles, security.RoleAdmin)

	for i := range res.Data {
		activeConnects := conns[res.Data[i].ID]
		res.Data[i].ActiveConnects = &activeConnects
		if r.WithFlags {
			resAdminUser := false
			for _, role := range res.Data[i].Roles {
				if role.ID == security.RoleAdmin {
					resAdminUser = true
					break
				}
			}

			deletable := res.Data[i].ID != curUserID && (!resAdminUser || isAdminRequest)
			res.Data[i].Deletable = &deletable
		}

		idpFields, _, ok := s.securityConfig.GetIdpFieldsExtraRolesAllowed(res.Data[i].Source)
		if ok {
			res.Data[i].IdPFields = idpFields
		}
	}

	return &res, nil
}

func (s *store) GetOneBy(ctx context.Context, id string) (*User, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
	}

	pipeline = append(pipeline, s.getNestedObjectsPipeline(s.authorProvider)...)
	cursor, err := s.userCollection.Aggregate(ctx, pipeline)
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

		idpFields, _, ok := s.securityConfig.GetIdpFieldsExtraRolesAllowed(user.Source)
		if ok {
			user.IdPFields = idpFields
		}

		return user, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r CreateRequest, requestRoles []string) (*User, error) {
	insertDoc, err := r.getBson(s.passwordEncoder)
	if err != nil {
		return nil, err
	}

	if slices.Contains(r.Roles, security.RoleAdmin) && !slices.Contains(requestRoles, security.RoleAdmin) {
		return nil, httperror.NewForbiddenError("You cannot assign yourself the admin role.")
	}

	var user *User
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		user = nil

		err = s.validateEditRequest(ctx, r.EditRequest, "")
		if err != nil {
			return err
		}

		_, err = s.userCollection.InsertOne(ctx, insertDoc)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Role{})
			}

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

func (s *store) Update(ctx context.Context, r BulkUpdateRequestItem, curUserID string, requestRoles []string) (*User, error) {
	if r.ID == curUserID && r.IsEnabled != nil && !*r.IsEnabled {
		return nil, httperror.NewForbiddenError("You cannot disable your own account.")
	}

	updateDoc, err := r.getBson(s.passwordEncoder)
	if err != nil {
		return nil, err
	}

	isAdminRequest := slices.Contains(requestRoles, security.RoleAdmin)

	var user *User
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		user = nil

		err = s.validateEditRequest(ctx, r.EditRequest, r.ID)
		if err != nil {
			return err
		}

		onlyOneAdmin, lastAdminID, err := s.checkLastAdmin(ctx)
		if err != nil {
			return err
		}

		if onlyOneAdmin && lastAdminID == r.ID && !slices.Contains(r.Roles, security.RoleAdmin) {
			return httperror.NewForbiddenError("You cannot remove the last admin role from yourself. At least one admin must remain.")
		}

		var prevUser security.User
		err = s.userCollection.FindOne(ctx, bson.M{"_id": r.ID}).Decode(&prevUser)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		if !isAdminRequest {
			if slices.Contains(prevUser.Roles, security.RoleAdmin) {
				return httperror.NewForbiddenError("You do not have permission to update an admin user.")
			}

			if slices.Contains(r.Roles, security.RoleAdmin) {
				return httperror.NewForbiddenError("You cannot assign yourself the admin role.")
			}
		}

		s.filterIdpFields(updateDoc, prevUser.Source, r.Roles, prevUser.IdPRoles)

		res, err := s.userCollection.UpdateOne(ctx,
			bson.M{"_id": r.ID},
			bson.M{"$set": updateDoc},
		)
		if err != nil || res.MatchedCount == 0 {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Role{})
			}

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

func (s *store) Patch(ctx context.Context, r BulkPatchRequestItem, curUserID string, requestRoles []string) (*User, error) {
	if r.ID == curUserID && r.IsEnabled != nil && !*r.IsEnabled {
		return nil, httperror.NewForbiddenError("You cannot disable your own account.")
	}

	updateDoc, err := r.getBson(s.passwordEncoder)
	if err != nil {
		return nil, err
	}

	isAdminRequest := slices.Contains(requestRoles, security.RoleAdmin)

	var user *User
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		user = nil

		err = validation.ValidateExist(ctx, s.viewCollection, r, "DefaultView", r.DefaultView)
		if err != nil {
			return err
		}

		err = validation.ValidateExist(ctx, s.colorThemeCollection, r, "UITheme", r.UITheme)
		if err != nil {
			return err
		}

		err = validation.ValidateExist(ctx, s.roleCollection, r, "Roles", r.Roles)
		if err != nil {
			return err
		}

		if r.Name != nil {
			err = s.validateName(ctx, *r.Name, r.ID, r)
			if err != nil {
				return err
			}
		}

		onlyOneAdmin, lastAdminID, err := s.checkLastAdmin(ctx)
		if err != nil {
			return err
		}

		if onlyOneAdmin && lastAdminID == r.ID && !slices.Contains(r.Roles, security.RoleAdmin) {
			return httperror.NewForbiddenError("You cannot remove the last admin role from yourself. At least one admin must remain.")
		}

		var prevUser security.User

		err = s.userCollection.FindOne(ctx, bson.M{"_id": r.ID}).Decode(&prevUser)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		if !isAdminRequest {
			if slices.Contains(prevUser.Roles, security.RoleAdmin) {
				return httperror.NewForbiddenError("You do not have permission to update an admin user.")
			}

			if slices.Contains(r.Roles, security.RoleAdmin) {
				return httperror.NewForbiddenError("You cannot assign yourself the admin role.")
			}
		}

		s.filterIdpFields(updateDoc, prevUser.Source, r.Roles, prevUser.IdPRoles)

		res, err := s.userCollection.UpdateOne(ctx,
			bson.M{"_id": r.ID},
			bson.M{"$set": updateDoc},
		)
		if err != nil || res.MatchedCount == 0 {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Role{})
			}

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

func (s *store) filterIdpFields(updateDoc bson.M, source string, requestRoles, idpRoles []string) {
	idpFields, allowExtraRoles, ok := s.securityConfig.GetIdpFieldsExtraRolesAllowed(source)
	if !ok {
		return
	}

	for _, idpField := range idpFields {
		delete(updateDoc, idpField)
	}

	if !allowExtraRoles {
		updateDoc["roles"] = idpRoles
		return
	}

	if len(requestRoles) != 0 {
		mergedRoles := make([]string, len(idpRoles), len(requestRoles)+len(idpRoles))
		copy(mergedRoles, idpRoles)

		idpRolesMap := make(map[string]bool, len(idpRoles))
		for _, role := range idpRoles {
			idpRolesMap[role] = true
		}

		for _, role := range requestRoles {
			if !idpRolesMap[role] {
				mergedRoles = append(mergedRoles, role)
			}
		}

		updateDoc["roles"] = mergedRoles
	}
}

func (s *store) Delete(ctx context.Context, r BulkDeleteRequestItem, userID string, requestRoles []string) (bool, error) {
	id := r.ID
	if id == userID {
		return false, httperror.NewForbiddenError("You cannot delete your own account.")
	}

	isAdminRequest := slices.Contains(requestRoles, security.RoleAdmin)

	var deleted int64
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		var prevUser security.User

		// required to get the author in action log listener.
		err := s.userCollection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}}).Decode(&prevUser)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		if slices.Contains(prevUser.Roles, security.RoleAdmin) && !isAdminRequest {
			return httperror.NewForbiddenError("You do not have permission to delete an admin user.")
		}

		deleted, err = s.userCollection.DeleteOne(ctx, bson.M{"_id": id})

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

	err = s.deleteNotifications(ctx, id)
	if err != nil {
		return false, err
	}

	err = s.clearTplTests(ctx, id)
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

func (s *store) deleteNotifications(ctx context.Context, id string) error {
	_, err := s.notificationCollection.DeleteMany(ctx, bson.M{
		"user": id,
	})

	return err
}

func (s *store) clearTplTests(ctx context.Context, id string) error {
	_, err := s.tplTestCollection.UpdateMany(ctx,
		bson.M{"data.user": id},
		bson.M{"$unset": bson.M{"data.user": ""}})

	return err
}

func (s *store) checkLastAdmin(ctx context.Context) (bool, string, error) {
	cursor, err := s.userCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"enable": true, "roles": security.RoleAdmin}},
		{"$group": bson.M{
			"_id":     nil,
			"count":   bson.M{"$sum": 1},
			"last_id": bson.M{"$first": "$_id"},
		}},
	})
	if err != nil {
		return false, "", err
	}

	defer cursor.Close(ctx)

	res := struct {
		Count  int64  `bson:"count"`
		LastID string `bson:"last_id"`
	}{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return false, "", err
		}
	}

	return res.Count <= 1, res.LastID, nil
}

func (s *store) getNestedObjectsPipeline(authorProvider author.Provider) []bson.M {
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
	pipeline = append(pipeline, authorProvider.Pipeline()...)

	return pipeline
}

func (s *store) validateEditRequest(ctx context.Context, r EditRequest, id string) error {
	err := validation.ValidateExist(ctx, s.viewCollection, r, "DefaultView", r.DefaultView)
	if err != nil {
		return err
	}

	err = validation.ValidateExist(ctx, s.colorThemeCollection, r, "UITheme", r.UITheme)
	if err != nil {
		return err
	}

	err = validation.ValidateExist(ctx, s.roleCollection, r, "Roles", r.Roles)
	if err != nil {
		return err
	}

	return s.validateName(ctx, r.Name, id, r)
}

func (s *store) validateName(ctx context.Context, name, id string, r any) error {
	q := bson.M{"$or": []bson.M{
		{"name": name},
		{"_id": name},
	}}
	if id != "" {
		q["_id"] = bson.M{"$ne": id}
	}

	err := s.userCollection.FindOne(ctx, q).Err()
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return err
	}

	if err == nil {
		return validation.NewError(
			validator.ValidationErrors{validation.NewFieldError("exist", "Name", "Name")},
			r,
		)
	}

	return nil
}

func getRolePipeline() []bson.M {
	return []bson.M{
		{"$unwind": bson.M{
			"path":                       "$roles",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "role_index",
		}},
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
