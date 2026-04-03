package statesettings

import (
	"cmp"
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/priority"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const StickySortField = "on_top"

type Store interface {
	GetByID(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
}

type store struct {
	dbClient                 mongo.DbClient
	dbCollection             mongo.DbCollection
	notifyDbCollection       mongo.DbCollection
	authorProvider           author.Provider
	stateSettingsUpdatesChan chan statesetting.RuleUpdatedMessage
	defaultSearchByFields    []string
	dupErrorParser           validation.DuplicateErrorParser
	transformer              patternfields.Transformer
}

func NewStore(
	dbClient mongo.DbClient,
	stateSettingsUpdatesChan chan statesetting.RuleUpdatedMessage,
	authorProvider author.Provider,
	transformer patternfields.Transformer,
) Store {
	return &store{
		dbClient:                 dbClient,
		dbCollection:             dbClient.Collection(mongo.StateSettingsMongoCollection),
		notifyDbCollection:       dbClient.Collection(mongo.EngineNotificationCollection),
		authorProvider:           authorProvider,
		stateSettingsUpdatesChan: stateSettingsUpdatesChan,
		defaultSearchByFields:    []string{"_id", "title"},
		dupErrorParser:           validation.NewDuplicateErrorParser(),
		transformer:              transformer,
	}
}

func (s *store) GetByID(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}, addEditableAndDeletableFields()}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var res Response

		err = cursor.Decode(&res)
		if err != nil {
			return nil, err
		}

		return &res, nil
	}

	return nil, nil
}

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	var pipeline []bson.M

	filter := mongoquery.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	pipeline = append(pipeline, addEditableAndDeletableFields())
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		s.getSortQuery(cmp.Or(query.SortBy, "title"), query.Sort),
	))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var result AggregationResult
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (s *store) Insert(ctx context.Context, r EditRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	var response *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		model, err := s.transformRequestToModel(ctx, r)
		if err != nil {
			return err
		}

		model.ID = utils.NewID()
		model.Created = &now
		model.Updated = &now
		_, err = s.dbCollection.InsertOne(ctx, model)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Response{})
			}

			return err
		}

		err = priority.UpdateFollowing(ctx, s.dbCollection, model.ID, model.Priority)
		if err != nil {
			return err
		}

		if model.Method == statesetting.MethodDependencies || model.Method == statesetting.MethodInherited {
			err = s.updateNotify(ctx)
			if err != nil {
				return err
			}
		}

		response, err = s.GetByID(ctx, model.ID)

		return err
	})
	if err != nil {
		return nil, err
	}

	if response != nil && (r.Method == statesetting.MethodDependencies || r.Method == statesetting.MethodInherited) {
		s.stateSettingsUpdatesChan <- statesetting.RuleUpdatedMessage{
			ID:         response.ID,
			NewPattern: response.EntityPattern,
			NewType:    *response.Type,
			Updated:    datetime.NewCpsTime(),
		}
	}

	return response, nil
}

func (s *store) Update(ctx context.Context, r EditRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	var response *Response
	var oldVersion statesetting.StateSetting
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		model, err := s.transformRequestToModel(ctx, r)
		if err != nil {
			return err
		}

		model.Updated = &now
		err = s.dbCollection.FindOneAndUpdate(
			ctx,
			bson.M{"_id": r.ID},
			bson.M{"$set": model},
			options.FindOneAndUpdate().SetReturnDocument(options.Before),
		).Decode(&oldVersion)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Response{})
			}

			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		err = priority.UpdateFollowing(ctx, s.dbCollection, r.ID, model.Priority)
		if err != nil {
			return err
		}

		if model.Method == statesetting.MethodDependencies || model.Method == statesetting.MethodInherited {
			err = s.updateNotify(ctx)
			if err != nil {
				return err
			}
		}

		response, err = s.GetByID(ctx, r.ID)

		return err
	})
	if err != nil {
		return nil, err
	}

	if response != nil && (response.Method == statesetting.MethodDependencies || response.Method == statesetting.MethodInherited) {
		s.stateSettingsUpdatesChan <- statesetting.RuleUpdatedMessage{
			ID:         response.ID,
			NewPattern: response.EntityPattern,
			NewType:    *response.Type,
			OldPattern: oldVersion.EntityPattern,
			OldType:    oldVersion.Type,
			Updated:    datetime.NewCpsTime(),
		}
	}

	return response, nil
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	if id == statesetting.ServiceID {
		return false, httperror.NewForbiddenError("The default service rule cannot be deleted.")
	}

	if id == statesetting.JUnitID {
		return false, httperror.NewForbiddenError("The default junit rule cannot be deleted.")
	}

	var oldVersion statesetting.StateSetting

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		oldVersion = statesetting.StateSetting{}

		// required to get the author in action log listener.
		err := s.dbCollection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}}).Decode(&oldVersion)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil || deleted == 0 {
			return err
		}

		return s.updateNotify(ctx)
	})
	if err != nil {
		return false, err
	}

	if oldVersion.Method == statesetting.MethodDependencies || oldVersion.Method == statesetting.MethodInherited {
		s.stateSettingsUpdatesChan <- statesetting.RuleUpdatedMessage{
			ID:         oldVersion.ID,
			OldPattern: oldVersion.EntityPattern,
			OldType:    oldVersion.Type,
			Updated:    datetime.NewCpsTime(),
		}
	}

	return oldVersion.ID != "", nil
}

func (s *store) getSortQuery(sortBy, sort string) bson.M {
	sortDir := 1
	if sort == pagination.SortDesc {
		sortDir = -1
	}

	q := bson.D{{Key: StickySortField, Value: -1}, {Key: sortBy, Value: sortDir}}
	if sortBy != "_id" {
		q = append(q, bson.E{Key: "_id", Value: 1})
	}

	return bson.M{"$sort": q}
}

// updateNotify updates a single document to trigger engine-che to update state settings rules
func (s *store) updateNotify(ctx context.Context) error {
	_, err := s.notifyDbCollection.UpdateOne(
		ctx,
		bson.M{"_id": statesetting.StateSettingsNotificationID},
		bson.M{"$set": bson.M{"time": time.Now()}},
		options.UpdateOne().SetUpsert(true),
	)

	return err
}

func (s *store) transformRequestToModel(ctx context.Context, r EditRequest) (statesetting.StateSetting, error) {
	model := statesetting.StateSetting{
		Method:          r.Method,
		Enabled:         r.Enabled,
		Priority:        r.Priority,
		StateThresholds: r.StateThresholds.ToModel(),
		JUnitThresholds: r.JUnitThresholds.ToModel(),
	}
	if r.Title != nil {
		model.Title = *r.Title
	}

	if r.Type != nil {
		model.Type = *r.Type
	}

	patterns, err := s.transformer.FetchCorporatePatterns(ctx,
		r.CorporateEntityPattern,
		r.CorporateInheritedEntityPattern,
	)
	if err != nil {
		return model, err
	}

	patternAliases := append(
		patternfields.GetAliases(r.EntityPattern),
		patternfields.GetAliases(r.InheritedEntityPattern)...,
	)
	aliases, err := s.transformer.FetchAliases(ctx, patternAliases)
	if err != nil {
		return model, err
	}

	aliasPropMap := make(map[string]bool)
	var applyAliasPropIDs []string
	var valErrs, applyErrs validator.ValidationErrors
	if r.CorporateEntityPattern != "" {
		r.EntityRequest, applyAliasPropIDs, applyErrs = s.transformer.ApplyEntityCorporatePattern(r.EntityRequest, patterns)
	} else if r.EntityPattern != nil && len(patternAliases) != 0 {
		r.EntityPattern, applyAliasPropIDs, applyErrs = s.transformer.ApplyAliases(r.EntityPattern, aliases)
	}

	valErrs = append(valErrs, applyErrs...)
	for _, id := range applyAliasPropIDs {
		aliasPropMap[id] = true
	}

	inheritedEntityRequest := patternfields.EntityRequest{
		EntityPattern:          r.InheritedEntityPattern,
		CorporateEntityPattern: r.CorporateInheritedEntityPattern,
	}
	if r.CorporateInheritedEntityPattern != "" {
		inheritedEntityRequest, applyAliasPropIDs, applyErrs = s.transformer.ApplyEntityCorporatePattern(inheritedEntityRequest, patterns, "CorporateInheritedEntityPattern")
	} else if r.InheritedEntityPattern != nil && len(patternAliases) != 0 {
		inheritedEntityRequest.EntityPattern, applyAliasPropIDs, applyErrs = s.transformer.ApplyAliases(inheritedEntityRequest.EntityPattern, aliases, "InheritedEntityPattern")
	}

	valErrs = append(valErrs, applyErrs...)
	for _, id := range applyAliasPropIDs {
		aliasPropMap[id] = true
	}

	if len(valErrs) > 0 {
		return model, validation.NewError(valErrs, r)
	}

	r.InheritedEntityPatternRequest = InheritedEntityPatternRequest{
		InheritedEntityPattern:          inheritedEntityRequest.EntityPattern,
		CorporateInheritedEntityPattern: inheritedEntityRequest.CorporateEntityPattern,
		CorporatePattern:                inheritedEntityRequest.CorporatePattern,
	}

	aliasPropIDs := make([]string, 0, len(aliasPropMap))
	for id := range aliasPropMap {
		aliasPropIDs = append(aliasPropIDs, id)
	}

	model.Aliases = aliasPropIDs
	model.EntityPatternFields = r.EntityRequest.ToModelWithoutFields(s.dbCollection.Name())
	model.InheritedEntityPatternFields = r.InheritedEntityPatternRequest.ToModelWithoutFields(s.dbCollection.Name())

	return model, nil
}

func addEditableAndDeletableFields() bson.M {
	return bson.M{
		"$addFields": bson.M{
			"editable": bson.M{"$cond": bson.M{
				"if":   bson.M{"$eq": bson.A{"$_id", statesetting.ServiceID}},
				"then": false,
				"else": true,
			}},
			"deletable": bson.M{"$cond": bson.M{
				"if":   bson.M{"$in": bson.A{"$_id", bson.A{statesetting.ServiceID, statesetting.JUnitID}}},
				"then": false,
				"else": true,
			}},
		},
	}
}
