package pbehaviortype

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
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
	Insert(ctx context.Context, model *Type) error
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Type, error)
	Update(ctx context.Context, id string, model *Type) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	db                    mongo.DbClient
	dbCollection          mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
}

// NewStore instantiates pbehavior type store.
func NewStore(db mongo.DbClient) Store {
	return &store{
		db:                    db,
		dbCollection:          db.Collection(mongo.PbehaviorTypeMongoCollection),
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

	sortBy := r.SortBy
	if sortBy == "" {
		sortBy = s.defaultSortBy
	}

	var project []bson.M
	if r.WithFlags {
		project = getDefaultAndDeletablePipeline(prioritiesOfDefaultTypes)
	}
	cursor, err := s.dbCollection.Aggregate(
		ctx,
		pagination.CreateAggregationPipeline(
			r.Query,
			pipeline,
			common.GetSortQuery(sortBy, r.Sort),
			project,
		),
		options.Aggregate().SetCollation(&options.Collation{Locale: "en"}),
	)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, fmt.Errorf("pbehavior types find return no data")
	}

	pbhResult = &AggregationResult{}
	return pbhResult, cursor.Decode(pbhResult)
}

// GetOneBy pbehavior type by id.
func (s *store) GetOneBy(ctx context.Context, id string) (*Type, error) {
	res := &Type{}

	if err := s.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(res); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return res, nil
}

// Insert creates new pbehavior type.
func (s *store) Insert(ctx context.Context, pt *Type) error {
	if pt.ID == "" {
		pt.ID = utils.NewID()
	}

	_, err := s.dbCollection.InsertOne(ctx, pt)
	if err != nil {
		if mongodriver.IsDuplicateKeyError(err) {
			return ErrorDuplicatePriority
		}

		return err
	}

	return nil
}

// Update pbehavior type.
func (s *store) Update(ctx context.Context, id string, pt *Type) (bool, error) {
	isDefault, err := s.IsDefault(ctx, id)
	if err != nil {
		return false, err
	}
	if pt.IconName == "" && (!isDefault || pt.Type != pbehavior.TypeActive) {
		return false, common.NewValidationError("icon_name", "IconName is missing.")
	}
	if isDefault {
		filter := bson.M{
			"_id":         id,
			"name":        pt.Name,
			"description": pt.Description,
			"type":        pt.Type,
			"priority":    pt.Priority,
			"icon_name":   pt.IconName,
		}
		if pt.IconName == "" {
			filter["icon_name"] = bson.M{"$in": bson.A{nil, ""}}
		}
		result, err := s.dbCollection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
			"color": pt.Color,
		}})
		if err != nil {
			return false, err
		}

		if result.MatchedCount == 0 {
			return false, ErrDefaultType
		}

		return true, nil
	}

	if pt.ID != id {
		pt.ID = id
	}
	result, err := s.dbCollection.ReplaceOne(ctx, bson.M{"_id": id}, pt)
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, nil
}

// Delete pbehavior type by id
func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	isDefault, err := s.IsDefault(ctx, id)
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

	r, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})

	return r > 0, err
}

// isLinkedToPbehavior checks if there is pbehavior with linked type.
func (s *store) isLinkedToPbehavior(ctx context.Context, id string) (bool, error) {
	pbhCollection := s.db.Collection(mongo.PbehaviorMongoCollection)
	res := pbhCollection.FindOne(ctx, bson.M{"$or": []bson.M{
		{"type_": id},
		{"exdates.type": id},
	}})
	if err := res.Err(); err == nil {
		return true, nil
	} else if err != mongodriver.ErrNoDocuments {
		return false, err
	}

	return false, nil
}

// isLinkedToException checks if there is execption with linked type.
func (s *store) isLinkedToException(ctx context.Context, id string) (bool, error) {
	exceptionCollection := s.db.Collection(mongo.PbehaviorExceptionMongoCollection)
	res := exceptionCollection.FindOne(ctx, bson.M{"exdates.type": id})
	if err := res.Err(); err == nil {
		return true, nil
	} else if err != mongodriver.ErrNoDocuments {
		return false, err
	}

	return false, nil
}

// isLinkedToAction checks if there is action with linked type.
func (s *store) isLinkedToAction(ctx context.Context, id string) (bool, error) {
	actionCollection := s.db.Collection(mongo.ScenarioMongoCollection)
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
	} else if err != mongodriver.ErrNoDocuments {
		return false, err
	}

	return false, nil
}

func (s *store) IsDefault(ctx context.Context, id string) (bool, error) {
	prioritiesOfDefaultTypes, err := s.getPrioritiesOfDefaultTypes(ctx)
	if err != nil {
		return false, err
	}

	res := s.dbCollection.FindOne(ctx, bson.M{"_id": id})
	if err := res.Err(); err == nil {
		var pbhType Type
		err = res.Decode(&pbhType)
		if err != nil {
			return false, err
		}

		for _, priority := range prioritiesOfDefaultTypes {
			if pbhType.Priority == priority {
				return true, nil
			}
		}
	} else if err != mongodriver.ErrNoDocuments {
		return false, err
	}

	return false, nil
}

func (s *store) getPrioritiesOfDefaultTypes(ctx context.Context) ([]int, error) {
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
		Type     string `bson:"_id"`
		Priority int    `bson:"priority"`
	}
	err = cursor.All(ctx, &doc)
	if err != nil {
		return nil, err
	}

	res := make([]int, len(doc))
	i := 0
	for _, d := range doc {
		res[i] = d.Priority
		i++
	}

	return res, nil
}

func getDefaultAndDeletablePipeline(prioritiesOfDefaultTypes []int) []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorMongoCollection,
			"localField":   "_id",
			"foreignField": "type",
			"as":           "pbhs",
		}},
		{"$lookup": bson.M{
			"from":         mongo.ScenarioMongoCollection,
			"localField":   "_id",
			"foreignField": "actions.parameters.type",
			"as":           "actions",
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
