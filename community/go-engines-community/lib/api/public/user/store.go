package user

import (
	"cmp"
	"context"
	"errors"
	"slices"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/dbvalidation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/sortparam"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store interface {
	Find(ctx context.Context, r ListRequest) (pagination.ListResponse[Response], error)
	GetByID(ctx context.Context, id string) (Response, error)
	Insert(ctx context.Context, r CreateRequest, requestRoles []string) (Response, error)
	Update(ctx context.Context, id string, r EditRequest, requestRoles []string) (Response, error)
	UpdateByPatch(ctx context.Context, id string, r PatchRequest, requestRoles []string) (Response, error)
	Delete(ctx context.Context, id, userID string, requestRoles []string) error
}

func NewStore(dbClient mongo.DbClient, securityConfig security.Config) Store {
	return &store{
		dbClient:         dbClient,
		dbCollection:     dbClient.Collection(mongo.UserCollection),
		dbRoleCollection: dbClient.Collection(mongo.RoleCollection),
		securityConfig:   securityConfig,
		dupErrorParser:   validation.NewDuplicateErrorParser(),
	}
}

type store struct {
	dbClient         mongo.DbClient
	dbCollection     mongo.DbCollection
	dbRoleCollection mongo.DbCollection
	securityConfig   security.Config
	dupErrorParser   validation.DuplicateErrorParser
}

func (s *store) Find(ctx context.Context, r ListRequest) (res pagination.ListResponse[Response], err error) {
	sort, err := sortparam.Validate(cmp.Or(r.Sort, "name"), []string{"_id", "name", "lastname", "firstname", "email", "external_id", "created", "updated"})
	if err != nil {
		return res, err
	}

	filter := bson.M{}
	if r.Enabled != nil {
		filter["enabled"] = *r.Enabled
	}

	if len(r.Roles) > 0 {
		filter["roles"] = bson.M{"$in": r.Roles}
	}

	if len(r.Sources) > 0 {
		filter["source"] = bson.M{"$in": r.Sources}
	}

	if r.Author != "" {
		filter["author"] = r.Author
	}

	beforeLimit := make([]bson.M, 0)
	if len(filter) != 0 {
		beforeLimit = append(beforeLimit, bson.M{"$match": filter})
	}

	// The lookups don't keep the document order, so sort again once they are done.
	afterLimit := append(nestedObjectsPipeline(), sort)
	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		sort,
		afterLimit,
	)
	opts := options.Aggregate().SetCollation(&options.Collation{Locale: "en"})
	ar, err := db.AggregateOne[pagination.Data[Response]](ctx, s.dbCollection, pipeline, opts)
	if err != nil {
		return res, err
	}

	return pagination.NewResponse[Response](r.Query, ar), nil
}

func (s *store) GetByID(ctx context.Context, id string) (Response, error) {
	pipeline := append([]bson.M{
		{"$match": bson.M{"_id": id}},
	}, nestedObjectsPipeline()...)

	return db.AggregateOne[Response](ctx, s.dbCollection, pipeline)
}

func (s *store) Insert(ctx context.Context, r CreateRequest, requestRoles []string) (res Response, err error) {
	if slices.Contains(r.Roles, security.RoleAdmin) && !slices.Contains(requestRoles, security.RoleAdmin) {
		return res, httperror.NewForbiddenError("You cannot assign yourself the admin role.")
	}

	id := utils.NewID()
	now := datetime.NewCpsTime()
	doc := bson.M{
		"_id":         id,
		"name":        r.Name,
		"lastname":    r.Lastname,
		"firstname":   r.Firstname,
		"email":       r.Email,
		"roles":       r.Roles,
		"enabled":     r.Enabled,
		"authkey":     utils.NewID(),
		"source":      r.Source,
		"external_id": r.ExternalID,
		"author":      r.Author,
		"created":     now,
		"updated":     now,
	}

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		// Reset the result, the transaction may be retried.
		res = Response{}

		err := s.validateEditRequest(ctx, r.EditRequest, "")
		if err != nil {
			return err
		}

		_, err = s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Response{})
			}

			return err
		}

		res, err = s.GetByID(ctx, id)

		return err
	})

	return res, err
}

func (s *store) Update(ctx context.Context, id string, r EditRequest, requestRoles []string) (res Response, err error) {
	if id == r.Author && r.Enabled != nil && !*r.Enabled {
		return res, httperror.NewForbiddenError("You cannot disable your own account.")
	}

	isAdminRequest := slices.Contains(requestRoles, security.RoleAdmin)
	set := bson.M{
		"name":      r.Name,
		"lastname":  r.Lastname,
		"firstname": r.Firstname,
		"email":     r.Email,
		"roles":     r.Roles,
		"enabled":   r.Enabled,
		"author":    r.Author,
		"updated":   datetime.NewCpsTime(),
	}

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		// Reset the result, the transaction may be retried.
		res = Response{}

		err := s.validateEditRequest(ctx, r, id)
		if err != nil {
			return err
		}

		prevUser, err := s.checkUpdateAllowed(ctx, id, r.Roles, isAdminRequest)
		if err != nil {
			return err
		}

		err = s.checkIdpFields(r, prevUser)
		if err != nil {
			return err
		}

		res, err = s.updateAndGet(ctx, id, set)

		return err
	})

	return res, err
}

func (s *store) UpdateByPatch(ctx context.Context, id string, r PatchRequest, requestRoles []string) (res Response, err error) {
	if id == r.Author && r.Enabled != nil && !*r.Enabled {
		return res, httperror.NewForbiddenError("You cannot disable your own account.")
	}

	isAdminRequest := slices.Contains(requestRoles, security.RoleAdmin)
	set := bson.M{
		"author":  r.Author,
		"updated": datetime.NewCpsTime(),
	}

	for field, value := range map[string]*string{
		"name":      r.Name,
		"firstname": r.Firstname,
		"lastname":  r.Lastname,
		"email":     r.Email,
	} {
		if value != nil {
			set[field] = *value
		}
	}

	if len(r.Roles) > 0 {
		set["roles"] = r.Roles
	}

	if r.Enabled != nil {
		set["enabled"] = r.Enabled
	}

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		// Reset the result, the transaction may be retried.
		res = Response{}

		err := dbvalidation.ValidateExist(ctx, s.dbRoleCollection, r, "Roles", r.Roles)
		if err != nil {
			return err
		}

		if r.Name != nil {
			err = s.validateName(ctx, *r.Name, id, r)
			if err != nil {
				return err
			}
		}

		prevUser, err := s.checkUpdateAllowed(ctx, id, r.Roles, isAdminRequest)
		if err != nil {
			return err
		}

		// The IDP fields are checked on the resulting user, so fill the untouched ones with the stored values.
		er := EditRequest{
			Name:      valueOr(r.Name, prevUser.Name),
			Email:     valueOr(r.Email, prevUser.Email),
			Lastname:  valueOr(r.Lastname, prevUser.Lastname),
			Firstname: valueOr(r.Firstname, prevUser.Firstname),
			Roles:     r.Roles,
		}
		if len(r.Roles) == 0 {
			er.Roles = prevUser.Roles
		}

		err = s.checkIdpFields(er, prevUser)
		if err != nil {
			return err
		}

		res, err = s.updateAndGet(ctx, id, set)

		return err
	})

	return res, err
}

func (s *store) Delete(ctx context.Context, id, userID string, requestRoles []string) error {
	if id == userID {
		return httperror.NewForbiddenError("You cannot delete your own account.")
	}

	isAdminRequest := slices.Contains(requestRoles, security.RoleAdmin)
	var deleted int64
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		// Reset the result, the transaction may be retried.
		deleted = 0

		var prevUser security.User

		// required to get the author in action log listener.
		err := s.dbCollection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}}).Decode(&prevUser)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		if slices.Contains(prevUser.Roles, security.RoleAdmin) && !isAdminRequest {
			return httperror.NewForbiddenError("You do not have permission to delete an admin user.")
		}

		deleted, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})

		return err
	})
	if err != nil {
		return err
	}

	if deleted == 0 {
		return httperror.ErrNotFound
	}

	return s.deleteCorrespondingObjects(ctx, id)
}

// updateAndGet applies set to the user and returns the updated document, or a zero
// Response if the user has been removed in the meantime.
func (s *store) updateAndGet(ctx context.Context, id string, set bson.M) (res Response, err error) {
	updRes, err := s.dbCollection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": set},
	)
	if err != nil {
		if mongodriver.IsDuplicateKeyError(err) {
			return res, s.dupErrorParser.Parse(err, Response{})
		}

		return res, err
	}

	if updRes.MatchedCount == 0 {
		return res, nil
	}

	return s.GetByID(ctx, id)
}

// checkUpdateAllowed fetches the stored user and verifies that the request may update it with newRoles.
func (s *store) checkUpdateAllowed(ctx context.Context, id string, newRoles []string, isAdminRequest bool) (prevUser security.User, err error) {
	// The roles are replaced only if the request contains them, so the last admin can
	// lose the role only if newRoles is filled and doesn't contain it.
	if len(newRoles) > 0 && !slices.Contains(newRoles, security.RoleAdmin) {
		soleAdminID, err := s.soleAdminID(ctx)
		if err != nil {
			return prevUser, err
		}

		if soleAdminID == id {
			return prevUser, httperror.NewForbiddenError("You cannot remove the last admin role from yourself. At least one admin must remain.")
		}
	}

	err = s.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&prevUser)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return prevUser, httperror.ErrNotFound
		}

		return prevUser, err
	}

	if !isAdminRequest {
		if slices.Contains(prevUser.Roles, security.RoleAdmin) {
			return prevUser, httperror.NewForbiddenError("You do not have permission to update an admin user.")
		}

		if slices.Contains(newRoles, security.RoleAdmin) {
			return prevUser, httperror.NewForbiddenError("You cannot assign yourself the admin role.")
		}
	}

	return prevUser, nil
}

func (s *store) validateEditRequest(ctx context.Context, r EditRequest, id string) error {
	err := dbvalidation.ValidateExist(ctx, s.dbRoleCollection, r, "Roles", r.Roles)
	if err != nil {
		return err
	}

	return s.validateName(ctx, r.Name, id, r)
}

func (s *store) validateName(ctx context.Context, name, id string, r any) error {
	q := bson.M{"name": name}
	if id != "" {
		q["_id"] = bson.M{"$ne": id}
	}

	err := s.dbCollection.FindOne(ctx, q).Err()
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil
	}

	if err != nil {
		return err
	}

	return validation.NewSingleError("exist", "Name", "Name", r)
}

// soleAdminID returns the id of the only enabled admin user or an empty string if there are none or several of them.
func (s *store) soleAdminID(ctx context.Context) (string, error) {
	cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"enabled": true, "roles": security.RoleAdmin}},
		{"$group": bson.M{
			"_id":     nil,
			"count":   bson.M{"$sum": 1},
			"last_id": bson.M{"$first": "$_id"},
		}},
	})
	if err != nil {
		return "", err
	}

	defer cursor.Close(ctx)

	res := struct {
		Count  int64  `bson:"count"`
		LastID string `bson:"last_id"`
	}{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return "", err
		}
	}

	if res.Count != 1 {
		return "", nil
	}

	return res.LastID, nil
}

// checkIdpFields verifies that the request doesn't modify the fields and the roles which
// are provided by the identity provider of the user.
func (s *store) checkIdpFields(r EditRequest, prevUser security.User) error {
	idpFields, allowExtraRoles, ok := s.securityConfig.GetIdpFieldsExtraRolesAllowed(prevUser.Source)
	if !ok {
		return nil
	}

	var changedIdpFields []string
	for _, idpField := range idpFields {
		var changed bool
		switch idpField {
		case "name":
			changed = r.Name != prevUser.Name
		case "email":
			changed = r.Email != prevUser.Email
		case "lastname":
			changed = r.Lastname != prevUser.Lastname
		case "firstname":
			changed = r.Firstname != prevUser.Firstname
		}

		if changed {
			changedIdpFields = append(changedIdpFields, idpField)
		}
	}

	newRoles := make(map[string]bool, len(r.Roles))
	for _, role := range r.Roles {
		newRoles[role] = true
	}

	if !allowExtraRoles {
		for _, role := range r.Roles {
			if !slices.Contains(prevUser.IdPRoles, role) {
				changedIdpFields = append(changedIdpFields, "roles")
				break
			}
		}
	}

	if len(changedIdpFields) > 0 {
		slices.Sort(changedIdpFields)

		return httperror.NewForbiddenError("You cannot change IDP field: " + strings.Join(changedIdpFields, " ") + ".")
	}

	for _, role := range prevUser.IdPRoles {
		if !newRoles[role] {
			return httperror.NewForbiddenError("You remove IDP role: " + role + ".")
		}
	}

	return nil
}

func (s *store) deleteCorrespondingObjects(ctx context.Context, id string) error {
	deletes := []struct {
		collection string
		filter     bson.M
	}{
		{mongo.UserPreferencesMongoCollection, bson.M{"user": id}},
		{mongo.PatternMongoCollection, bson.M{"author": id, "is_corporate": false}},
		{mongo.ViewGroupMongoCollection, bson.M{"author": id, "is_private": true}},
		{mongo.ViewMongoCollection, bson.M{"author": id, "is_private": true}},
		{mongo.ViewTabMongoCollection, bson.M{"author": id, "is_private": true}},
		{mongo.WidgetMongoCollection, bson.M{"author": id, "is_private": true}},
		{mongo.WidgetFiltersMongoCollection, bson.M{
			"author": id,
			"$or": bson.A{
				bson.M{"is_user_preference": true},
				bson.M{"is_private": true},
			},
		}},
		{mongo.ShareTokenMongoCollection, bson.M{"user": id}},
		{mongo.UserNotificationCollection, bson.M{"user": id}},
	}

	for _, d := range deletes {
		_, err := s.dbClient.Collection(d.collection).DeleteMany(ctx, d.filter)
		if err != nil {
			return err
		}
	}

	_, err := s.dbClient.Collection(mongo.TemplateTestCollection).UpdateMany(ctx,
		bson.M{"data.user": id},
		bson.M{"$unset": bson.M{"data.user": ""}})

	return err
}

func nestedObjectsPipeline() []bson.M {
	pipeline := mongoquery.ArrayReferencePipeline("roles", mongo.RoleCollection)
	pipeline = append(pipeline, mongoquery.AuthorPipeline()...)

	return pipeline
}

func valueOr[T any](v *T, def T) T {
	if v == nil {
		return def
	}

	return *v
}
