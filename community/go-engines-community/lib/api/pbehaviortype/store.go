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

// Store is an interface for pbhavior types storage
type Store interface {
	Insert(ctx context.Context, model *Type) error
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Type, error)
	Update(ctx context.Context, id string, model *Type) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	db                    mongo.DbClient
	defaultSearchByFields []string
	defaultSortBy         string
}

// NewStore instantiates pbehavior type store.
func NewStore(db mongo.DbClient) Store {
	// temporarily until feat/#2344/mongo-indexes not merged
	keys := bson.M{"priority": 1}
	indexOptions := options.Index().SetBackground(true).SetUnique(true)
	_, err := db.Collection(pbehavior.TypeCollectionName).Indexes().CreateOne(
		context.Background(), mongodriver.IndexModel{
			Keys:    &keys,
			Options: indexOptions,
		})
	if err != nil {
		panic(err)
	}

	return &store{
		db:                    db,
		defaultSearchByFields: []string{"_id", "name", "description", "type"},
		defaultSortBy:         "name",
	}
}

func (s *store) getCollection() mongo.DbCollection {
	return s.db.Collection(pbehavior.TypeCollectionName)
}

// Find pbhavior types according to query.
func (s *store) Find(ctx context.Context, r ListRequest) (pbhResult *AggregationResult, err error) {
	prioritiesOfDefaultTypes, err := s.getPrioritiesOfDefaultTypes(ctx)
	if err != nil {
		return nil, err
	}
	collection := s.getCollection()

	pipeline := make([]bson.M, 0)
	if r.OnlyDefault {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"priority": bson.M{"$in": prioritiesOfDefaultTypes}}})
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
		project = getEditableAndDeletablePipeline(prioritiesOfDefaultTypes)
	}
	cursor, err := collection.Aggregate(
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
	collection := s.getCollection()

	if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(res); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return res, nil
}

// Create new pbehavior type.
func (s *store) Insert(ctx context.Context, pt *Type) error {

	if pt.ID == "" {
		pt.ID = utils.NewID()
	}

	_, err := s.getCollection().InsertOne(ctx, pt)
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
	if isDefault {
		return false, ErrDefaultType
	}

	if pt.ID != id {
		pt.ID = id
	}
	result, err := s.getCollection().ReplaceOne(ctx, bson.M{"_id": id}, pt)
	if err != nil {
		return false, err
	}
	return result.ModifiedCount > 0 || result.MatchedCount > 0, nil
}

// Delete pbehavior type by id
func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	ToPbehavior, err := s.isLinkedToPbehavior(ctx, id)
	if err != nil {
		return false, err
	}
	if ToPbehavior {
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

	isDefault, err := s.IsDefault(ctx, id)
	if err != nil {
		return false, err
	}
	if isDefault {
		return false, ErrDefaultType
	}

	r, err := s.getCollection().DeleteOne(ctx, bson.M{"_id": id})

	return r > 0, err
}

// isLinkedToPbehavior checks if there is pbehavior with linked type.
func (s *store) isLinkedToPbehavior(ctx context.Context, id string) (bool, error) {
	pbhCollection := s.db.Collection(pbehavior.PBehaviorCollectionName)
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
	exceptionCollection := s.db.Collection(pbehavior.ExceptionCollectionName)
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

	res := s.getCollection().FindOne(ctx, bson.M{"_id": id})
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
	collection := s.getCollection()
	cursor, err := collection.Aggregate(ctx, []bson.M{
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

func getEditableAndDeletablePipeline(prioritiesOfDefaultTypes []int) []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         pbehavior.PBehaviorCollectionName,
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
			"editable": bson.M{"$not": bson.M{"$in": bson.A{"$priority", prioritiesOfDefaultTypes}}},
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
