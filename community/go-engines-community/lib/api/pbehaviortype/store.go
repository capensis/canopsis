package pbehaviortype

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"

	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Store is an interface for pbhavior types storage
type Store interface {
	Insert(model *Type) error
	Find(r ListRequest) (*AggregationResult, error)
	GetOneBy(id string) (*Type, error)
	Update(id string, model *Type) (bool, error)
	Delete(id string) (bool, error)
}

type store struct {
	db            mongo.DbClient
	defaultSortBy string
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
		db:            db,
		defaultSortBy: "name",
	}
}

func (s *store) getCollection() mongo.DbCollection {
	return s.db.Collection(pbehavior.TypeCollectionName)
}

// Find pbhavior types according to query.
func (s *store) Find(r ListRequest) (pbhResult *AggregationResult, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	prioritiesOfDefaultTypes, err := s.getPrioritiesOfDefaultTypes()
	if err != nil {
		return nil, err
	}
	collection := s.getCollection()
	filter := bson.M{}

	if r.Search != "" {
		searchRegexp := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", r.Search),
			Options: "i",
		}

		filter["$or"] = []bson.M{
			{"name": searchRegexp},
			{"description": searchRegexp},
		}
	}
	if r.OnlyDefault {
		filter["priority"] = bson.M{"$in": prioritiesOfDefaultTypes}
	}

	sortBy := r.SortBy
	if sortBy == "" {
		sortBy = s.defaultSortBy
	}

	var project []bson.M
	if r.WithFlags {
		project = getEditableAndDeletablePipeline(prioritiesOfDefaultTypes)
	}
	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		[]bson.M{{"$match": filter}},
		common.GetSortQuery(sortBy, r.Sort),
		project,
	)
	cursor, err := collection.Aggregate(ctx, pipeline,
		options.Aggregate().SetCollation(&options.Collation{Locale: "en"}))

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
func (s *store) GetOneBy(id string) (*Type, error) {
	res := &Type{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
func (s *store) Insert(pt *Type) error {
	const errorDuplicateKey = 11000
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if pt.ID == "" {
		pt.ID = utils.NewID()
	}

	_, err := s.getCollection().InsertOne(ctx, pt)
	if err != nil {
		if mwe, ok := err.(mongodriver.WriteException); ok {
			for _, we := range mwe.WriteErrors {
				if we.Code == errorDuplicateKey {
					return ErrorDuplicatePriority
				}
			}
		}
	}

	return err
}

// Update pbehavior type.
func (s *store) Update(id string, pt *Type) (bool, error) {
	isDefault, err := s.IsDefault(id)
	if err != nil {
		return false, err
	}
	if isDefault {
		return false, ErrDefaultType
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
func (s *store) Delete(id string) (bool, error) {
	ToPbehavior, err := s.isLinkedToPbehavior(id)
	if err != nil {
		return false, err
	}
	if ToPbehavior {
		return false, ErrLinkedTypeToPbehavior
	}
	isLinkedToException, err := s.isLinkedToException(id)
	if err != nil {
		return false, err
	}
	if isLinkedToException {
		return false, ErrLinkedTypeToException
	}
	isLinkedToAction, err := s.isLinkedToAction(id)
	if err != nil {
		return false, err
	}
	if isLinkedToAction {
		return false, ErrLinkedToActionType
	}

	isDefault, err := s.IsDefault(id)
	if err != nil {
		return false, err
	}
	if isDefault {
		return false, ErrDefaultType
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, err := s.getCollection().DeleteOne(ctx, bson.M{"_id": id})

	return r > 0, err
}

// isLinkedToPbehavior checks if there is pbehavior with linked type.
func (s *store) isLinkedToPbehavior(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
func (s *store) isLinkedToException(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
func (s *store) isLinkedToAction(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

func (s *store) IsDefault(id string) (bool, error) {
	prioritiesOfDefaultTypes, err := s.getPrioritiesOfDefaultTypes()
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

func (s *store) getPrioritiesOfDefaultTypes() ([]int, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
