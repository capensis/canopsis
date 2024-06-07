package pbehaviortype

import (
	"cmp"
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	libpriority "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/priority"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Store is an interface for pbehavior types storage
type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetByID(ctx context.Context, id string) (*Response, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
	GetNextPriority(ctx context.Context) (int64, error)
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	authorProvider        author.Provider
	defaultSearchByFields []string
	defaultSortBy         string
}

// NewStore instantiates pbehavior type store.
func NewStore(db mongo.DbClient, authorProvider author.Provider) Store {
	return &store{
		dbClient:              db,
		dbCollection:          db.Collection(mongo.PbehaviorTypeMongoCollection),
		authorProvider:        authorProvider,
		defaultSearchByFields: []string{"_id", "name", "description", "type"},
		defaultSortBy:         "name",
	}
}

// Find pbehavior types according to query.
func (s *store) Find(ctx context.Context, r ListRequest) (pbhResult *AggregationResult, err error) {
	prioritiesOfDefaultTypes, err := s.getPrioritiesOfDefaultTypes(ctx)
	if err != nil {
		return nil, err
	}

	match := bson.M{}
	if !r.WithHidden {
		match["hidden"] = bson.M{"$in": bson.A{false, nil}}
	}

	if r.OnlyDefault {
		match["priority"] = bson.M{"$in": prioritiesOfDefaultTypes}
	}

	if len(r.Types) > 0 {
		match["type"] = bson.M{"$in": r.Types}
	}

	pipeline := make([]bson.M, 0)
	if len(match) > 0 {
		pipeline = append(pipeline, bson.M{"$match": match})
	}

	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	var project []bson.M
	if r.WithFlags {
		project = getDefaultAndDeletablePipeline(prioritiesOfDefaultTypes)
	}
	cursor, err := s.dbCollection.Aggregate(
		ctx,
		pagination.CreateAggregationPipeline(
			r.Query,
			pipeline,
			common.GetSortQuery(cmp.Or(r.SortBy, s.defaultSortBy), r.Sort),
			project,
		),
		options.Aggregate().SetCollation(&options.Collation{Locale: "en"}),
	)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, errors.New("pbehavior types find return no data")
	}

	pbhResult = &AggregationResult{}
	return pbhResult, cursor.Decode(pbhResult)
}

// GetByID pbehavior type by id.
func (s *store) GetByID(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if cursor.Next(ctx) {
		var res Response

		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}

		return &res, nil
	}

	return nil, nil
}

// Insert creates new pbehavior type.
func (s *store) Insert(ctx context.Context, r CreateRequest) (*Response, error) {
	now := datetime.NewCpsTime()

	doc := transformRequestToDocument(r.EditRequest)
	doc.ID = cmp.Or(r.ID, utils.NewID())
	doc.IconName = r.IconName
	doc.Created = &now
	doc.Updated = &now

	prioritiesOfDefaultTypes, err := s.getPrioritiesOfDefaultTypes(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range prioritiesOfDefaultTypes {
		if p == doc.Priority {
			return nil, common.NewValidationError("priority", "Priority is taken by default type.")
		}
	}

	var res *Response

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = nil

		err := libpriority.UpdateFollowing(ctx, s.dbCollection, doc.ID, doc.Priority)
		if err != nil {
			return err
		}

		_, err = s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}

		res, err = s.GetByID(ctx, doc.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Update pbehavior type.
func (s *store) Update(ctx context.Context, r UpdateRequest) (*Response, error) {
	prioritiesOfDefaultTypes, err := s.getPrioritiesOfDefaultTypes(ctx)
	if err != nil {
		return nil, err
	}

	now := datetime.NewCpsTime()

	doc := transformRequestToDocument(r.EditRequest)
	doc.ID = r.ID
	doc.IconName = r.IconName
	doc.Updated = &now

	var res *Response

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = nil

		var oldType pbehavior.Type
		err = s.dbCollection.FindOne(ctx, bson.M{"_id": doc.ID},
			options.FindOne().SetProjection(bson.M{"priority": 1})).Decode(&oldType)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
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
			return common.NewValidationError("priority", "Priority is taken by default type.")
		}

		if doc.IconName == "" && (!isDefault || doc.Type != pbehavior.TypeActive) {
			return common.NewValidationError("icon_name", "IconName is missing.")
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
				"color":  doc.Color,
				"hidden": doc.Hidden,
			}})
			if err != nil {
				return err
			}

			if result.MatchedCount == 0 {
				return ErrDefaultType
			}
		} else {
			err = libpriority.UpdateFollowing(ctx, s.dbCollection, doc.ID, doc.Priority)
			if err != nil {
				return err
			}

			result, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": doc.ID}, bson.M{"$set": doc})
			if err != nil || result.MatchedCount == 0 {
				return err
			}
		}

		res, err = s.GetByID(ctx, doc.ID)
		return err
	})

	return res, err
}

// Delete pbehavior type by id
func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	isDefault, err := s.isDefault(ctx, id)
	if err != nil {
		return false, err
	}
	if isDefault {
		return false, ErrDefaultType
	}
	isLinkedToPbh, err := s.isLinkedToPbehavior(ctx, id)
	if err != nil {
		return false, err
	}
	if isLinkedToPbh {
		return false, ErrLinkedTypeToPbehavior
	}
	isLinkedToException, err := s.isLinkedToException(ctx, id)
	if err != nil {
		return false, err
	}
	if isLinkedToException {
		return false, ErrLinkedTypeToException
	}
	isLinkedToAction, err := s.isLinkedToAction(ctx, id)
	if err != nil {
		return false, err
	}
	if isLinkedToAction {
		return false, ErrLinkedToActionType
	}

	var deleted int64

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		deleted, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		return err
	})

	return deleted > 0, err
}

func (s *store) GetNextPriority(ctx context.Context) (int64, error) {
	cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
		{"$group": bson.M{
			"_id":      nil,
			"priority": bson.M{"$max": "$priority"},
		}},
	})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)
	res := struct {
		Priority int64 `bson:"priority"`
	}{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return 0, err
		}
	}

	return res.Priority + 1, nil
}

// isLinkedToPbehavior checks if there is pbehavior with linked type.
func (s *store) isLinkedToPbehavior(ctx context.Context, id string) (bool, error) {
	pbhCollection := s.dbClient.Collection(mongo.PbehaviorMongoCollection)
	res := pbhCollection.FindOne(ctx, bson.M{"$or": []bson.M{
		{"type_": id},
		{"exdates.type": id},
	}})
	if err := res.Err(); err == nil {
		return true, nil
	} else if !errors.Is(err, mongodriver.ErrNoDocuments) {
		return false, err
	}

	return false, nil
}

// isLinkedToException checks if there is exception with linked type.
func (s *store) isLinkedToException(ctx context.Context, id string) (bool, error) {
	exceptionCollection := s.dbClient.Collection(mongo.PbehaviorExceptionMongoCollection)
	res := exceptionCollection.FindOne(ctx, bson.M{"exdates.type": id})
	if err := res.Err(); err == nil {
		return true, nil
	} else if !errors.Is(err, mongodriver.ErrNoDocuments) {
		return false, err
	}

	return false, nil
}

// isLinkedToAction checks if there is action with linked type.
func (s *store) isLinkedToAction(ctx context.Context, id string) (bool, error) {
	actionCollection := s.dbClient.Collection(mongo.ScenarioMongoCollection)
	res := actionCollection.FindOne(ctx, bson.M{
		"actions": bson.M{
			"$elemMatch": bson.M{
				"type":            types.ActionTypePbehavior,
				"parameters.type": id,
			},
		},
	})
	if err := res.Err(); err == nil {
		return true, nil
	} else if !errors.Is(err, mongodriver.ErrNoDocuments) {
		return false, err
	}

	return false, nil
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

func getDefaultAndDeletablePipeline(prioritiesOfDefaultTypes []int64) []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from": mongo.PbehaviorMongoCollection,
			"let":  bson.M{"type": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$type", "$$type"}}}},
				{"$limit": 1},
				{"$project": bson.M{"_id": 1}},
			},
			"as": "pbhs",
		}},
		{"$lookup": bson.M{
			"from": mongo.ScenarioMongoCollection,
			"let":  bson.M{"type": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$actions.parameters.type", "$$type"}}}},
				{"$limit": 1},
				{"$project": bson.M{"_id": 1}},
			},
			"as": "actions",
		}},
		{"$addFields": bson.M{
			"default": bson.M{"$in": bson.A{"$priority", prioritiesOfDefaultTypes}},
			"deletable": bson.M{"$and": []bson.M{
				{"$not": bson.M{"$in": bson.A{"$priority", prioritiesOfDefaultTypes}}},
				{"$eq": bson.A{bson.M{"$size": "$pbhs"}, 0}},
				{"$eq": bson.A{bson.M{"$size": "$actions"}, 0}},
			}},
		}},
		{"$project": bson.M{
			"pbhs":    0,
			"actions": 0,
		}},
	}
}

func transformRequestToDocument(request EditRequest) *pbehavior.Type {
	return &pbehavior.Type{
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		Priority:    request.Priority,
		Color:       request.Color,
		Hidden:      request.Hidden,
		Author:      request.Author,
	}
}
