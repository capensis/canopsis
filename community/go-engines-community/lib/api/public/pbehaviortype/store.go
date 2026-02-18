package pbehaviortype

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/dbvalidation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	libpriority "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/priority"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store crud.Store[ListRequest, CreateRequest, UpdateRequest, Response]

type store struct {
	dbClient                 mongo.DbClient
	dbCollection             mongo.DbCollection
	dbPbhCollection          mongo.DbCollection
	dbPbhExceptionCollection mongo.DbCollection
	dbScenarioCollection     mongo.DbCollection
	authorProvider           author.Provider
	defaultSearchByFields    []string
	dupErrorParser           validation.DuplicateErrorParser
}

func NewStore(db mongo.DbClient, authorProvider author.Provider) Store {
	return &store{
		dbClient:                 db,
		dbCollection:             db.Collection(mongo.PbehaviorTypeMongoCollection),
		dbPbhCollection:          db.Collection(mongo.PbehaviorMongoCollection),
		dbPbhExceptionCollection: db.Collection(mongo.PbehaviorExceptionMongoCollection),
		dbScenarioCollection:     db.Collection(mongo.ScenarioCollection),
		authorProvider:           authorProvider,
		defaultSearchByFields:    []string{"_id", "name", "description", "type"},
		dupErrorParser:           validation.NewDuplicateErrorParser(),
	}
}

func (s *store) Find(ctx context.Context, r ListRequest) (pagination.ListResponse[Response], error) {
	res := pagination.ListResponse[Response]{}
	prioritiesOfDefaultTypes, err := s.getPrioritiesOfDefaultTypes(ctx)
	if err != nil {
		return res, err
	}

	match := bson.M{
		"hidden": bson.M{"$in": bson.A{false, nil}},
	}
	if r.OnlyDefault {
		match["priority"] = bson.M{"$in": prioritiesOfDefaultTypes}
	}

	if len(r.Types) > 0 {
		match["type"] = bson.M{"$in": r.Types}
	}

	beforeLimit := make([]bson.M, 0)
	if len(match) > 0 {
		beforeLimit = append(beforeLimit, bson.M{"$match": match})
	}

	filter := mongoquery.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		beforeLimit = append(beforeLimit, bson.M{"$match": filter})
	}

	afterLimit := []bson.M{
		{"$addFields": bson.M{
			"default": bson.M{"$in": bson.A{"$priority", prioritiesOfDefaultTypes}},
		}},
	}
	afterLimit = append(afterLimit, s.authorProvider.Pipeline()...)
	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		mongoquery.GetSortQuery(r.SortBy, r.Sort),
		afterLimit,
	)
	opts := options.Aggregate().SetCollation(&options.Collation{Locale: "en"})
	ar, err := db.AggregateOne[pagination.Data[Response]](ctx, s.dbCollection, pipeline, opts)
	if err != nil {
		return res, err
	}

	res = pagination.NewResponse[Response](r.Query, ar)

	return res, nil
}

func (s *store) GetByID(ctx context.Context, id string) (Response, error) {
	prioritiesOfDefaultTypes, err := s.getPrioritiesOfDefaultTypes(ctx)
	if err != nil {
		return Response{}, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$addFields": bson.M{
			"default": bson.M{"$in": bson.A{"$priority", prioritiesOfDefaultTypes}},
		}},
	}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	return db.AggregateOne[Response](ctx, s.dbCollection, pipeline)
}

func (s *store) Insert(ctx context.Context, r CreateRequest) (Response, error) {
	res := Response{}
	now := datetime.NewCpsTime()
	doc := s.transformRequestToDocument(r.EditRequest)
	doc.ID = utils.NewID()
	doc.IconName = r.IconName
	doc.Created = &now
	doc.Updated = &now
	prioritiesOfDefaultTypes, err := s.getPrioritiesOfDefaultTypes(ctx)
	if err != nil {
		return res, err
	}

	for _, p := range prioritiesOfDefaultTypes {
		if p == doc.Priority {
			return res, validation.NewSingleError("not_applicable", "Priority", "Priority", r)
		}
	}

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = Response{}
		err = libpriority.UpdateFollowing(ctx, s.dbCollection, doc.ID, doc.Priority)
		if err != nil {
			return err
		}

		_, err = s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, doc)
			}

			return err
		}

		res, err = s.GetByID(ctx, doc.ID)

		return err
	})

	return res, err
}

func (s *store) Update(ctx context.Context, id string, r UpdateRequest) (Response, error) {
	var res Response
	prioritiesOfDefaultTypes, err := s.getPrioritiesOfDefaultTypes(ctx)
	if err != nil {
		return res, err
	}

	now := datetime.NewCpsTime()
	doc := s.transformRequestToDocument(r.EditRequest)
	doc.ID = id
	doc.IconName = r.IconName
	doc.Updated = &now
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = Response{}

		var oldType pbehavior.Type
		err = s.dbCollection.
			FindOne(ctx, bson.M{"_id": doc.ID}, options.FindOne().SetProjection(bson.M{"priority": 1})).
			Decode(&oldType)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return httperror.ErrNotFound
			}

			return err
		}

		isDefault := false
		isTaken := false
		for _, priority := range prioritiesOfDefaultTypes {
			if oldType.Priority == priority {
				isDefault = true
				break
			}

			isTaken = isTaken || doc.Priority == priority
		}

		if isTaken && !isDefault {
			return validation.NewSingleError("not_applicable", "Priority", "Priority", r)
		}

		if doc.IconName == "" && (!isDefault || doc.Type != pbehavior.TypeActive) {
			return validation.NewSingleError("required", "IconName", "IconName", r)
		}

		if isDefault {
			filter := bson.M{
				"_id":         doc.ID,
				"name":        doc.Name,
				"description": doc.Description,
				"type":        doc.Type,
				"priority":    doc.Priority,
				"icon_name":   doc.IconName,
			}
			if doc.IconName == "" {
				filter["icon_name"] = bson.M{"$in": bson.A{nil, ""}}
			}

			result, err := s.dbCollection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
				"color": doc.Color,
			}})
			if err != nil {
				if mongodriver.IsDuplicateKeyError(err) {
					return s.dupErrorParser.Parse(err, doc)
				}

				return err
			}

			if result.MatchedCount == 0 {
				return httperror.NewForbiddenError("The default type cannot be modified.")
			}
		} else {
			err = libpriority.UpdateFollowing(ctx, s.dbCollection, doc.ID, doc.Priority)
			if err != nil {
				return err
			}

			result, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": doc.ID}, bson.M{"$set": doc})
			if err != nil || result.MatchedCount == 0 {
				if mongodriver.IsDuplicateKeyError(err) {
					return s.dupErrorParser.Parse(err, doc)
				}

				return err
			}
		}

		res, err = s.GetByID(ctx, doc.ID)

		return err
	})

	return res, err
}

func (s *store) Delete(ctx context.Context, id, userID string) error {
	return s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		err := s.validateDeleteRequest(ctx, id)
		if err != nil {
			return err
		}

		// required to get the author in action log listener.
		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil {
			return err
		}

		if res.MatchedCount == 0 {
			return httperror.ErrNotFound
		}

		deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			return err
		}

		if deleted == 0 {
			return httperror.ErrNotFound
		}

		return nil
	})
}

func (s *store) validateDeleteRequest(ctx context.Context, id string) error {
	isDefault, err := s.isDefault(ctx, id)
	if err != nil {
		return err
	}

	if isDefault {
		return httperror.NewForbiddenError("The default type cannot be deleted.")
	}

	err = dbvalidation.ValidateLinkedReference(ctx, s.dbPbhCollection, bson.M{
		"$or": []bson.M{
			{"type_": id},
			{"exdates.type": id},
		},
	}, "type", "a pbehavior")
	if err != nil {
		return err
	}

	err = dbvalidation.ValidateLinkedReference(ctx, s.dbPbhExceptionCollection, bson.M{
		"exdates.type": id,
	}, "type", "an exception")
	if err != nil {
		return err
	}

	return dbvalidation.ValidateLinkedReference(ctx, s.dbScenarioCollection, bson.M{
		"actions": bson.M{
			"$elemMatch": bson.M{
				"type":            types.ActionTypePbehavior,
				"parameters.type": id,
			},
		},
	}, "type", "an action scenario")
}

func (s *store) isDefault(ctx context.Context, id string) (bool, error) {
	prioritiesOfDefaultTypes, err := s.getPrioritiesOfDefaultTypes(ctx)
	if err != nil {
		return false, err
	}

	var pbhType pbehavior.Type
	err = s.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&pbhType)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	for _, priority := range prioritiesOfDefaultTypes {
		if pbhType.Priority == priority {
			return true, nil
		}
	}

	return false, nil
}

func (s *store) getPrioritiesOfDefaultTypes(ctx context.Context) ([]int64, error) {
	cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
		{"$group": bson.M{
			"_id":      "$type",
			"priority": bson.M{"$min": "$priority"},
		}},
	})
	if err != nil {
		return nil, err
	}

	var doc []struct {
		Priority int64 `bson:"priority"`
	}
	err = cursor.All(ctx, &doc)
	if err != nil {
		return nil, err
	}

	res := make([]int64, len(doc))
	i := 0
	for _, d := range doc {
		res[i] = d.Priority
		i++
	}

	return res, nil
}

func (s *store) transformRequestToDocument(request EditRequest) pbehavior.Type {
	return pbehavior.Type{
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		Priority:    request.Priority,
		Color:       request.Color,
		Author:      request.Author,
	}
}
