package pbehavior

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorexception"
	apireason "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorreason"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type ModelTransformer interface {
	TransformCreateRequestToModel(ctx context.Context, request CreateRequest) (*Response, error)
	TransformBulkCreateRequestToModels(ctx context.Context, requests []CreateRequest) ([]*Response, error)
	TransformUpdateRequestToModel(ctx context.Context, request UpdateRequest) (*Response, error)
	TransformBulkUpdateRequestToModels(ctx context.Context, requests []BulkUpdateRequestItem) ([]*Response, error)
	Patch(ctx context.Context, req PatchRequest, model *Response) error
}

type modelTransformer struct {
	dbClient             mongo.DbClient
	typeCollection       mongo.DbCollection
	reasonCollection     mongo.DbCollection
	reasonTransformer    apireason.ModelTransformer
	exceptionCollection  mongo.DbCollection
	exceptionTransformer pbehaviorexception.ModelTransformer
}

func NewModelTransformer(
	dbClient mongo.DbClient,
	reasonTransformer apireason.ModelTransformer,
	exceptionTransformer pbehaviorexception.ModelTransformer,
) ModelTransformer {
	return &modelTransformer{
		dbClient:             dbClient,
		typeCollection:       dbClient.Collection(pbehavior.TypeCollectionName),
		reasonCollection:     dbClient.Collection(pbehavior.ReasonCollectionName),
		reasonTransformer:    reasonTransformer,
		exceptionCollection:  dbClient.Collection(pbehavior.ExceptionCollectionName),
		exceptionTransformer: exceptionTransformer,
	}
}

func (t *modelTransformer) TransformCreateRequestToModel(ctx context.Context, request CreateRequest) (*Response, error) {
	reason, err := t.transformReasonToModel(ctx, request.Reason)
	if err != nil {
		return nil, err
	}

	pbhType, err := t.transformTypeToModel(ctx, request.Type)
	if err != nil {
		return nil, err
	}

	exdates, err := t.exceptionTransformer.TransformExdatesRequestToModel(ctx, request.Exdates)
	if err != nil {
		return nil, err
	}

	exceptions, err := t.transformExceptionsToModel(ctx, request.Exceptions)
	if err != nil {
		return nil, err
	}

	return &Response{
		ID:         request.ID,
		Author:     request.Author,
		Enabled:    *request.Enabled,
		Filter:     NewFilter(request.Filter),
		Name:       request.Name,
		Reason:     reason,
		RRule:      request.RRule,
		Start:      &request.Start,
		Stop:       request.Stop,
		Type:       pbhType,
		Exdates:    exdates,
		Exceptions: exceptions,
	}, nil
}

func (t *modelTransformer) TransformUpdateRequestToModel(ctx context.Context, request UpdateRequest) (*Response, error) {
	return t.TransformCreateRequestToModel(ctx, CreateRequest(request))
}

func (t *modelTransformer) TransformBulkCreateRequestToModels(ctx context.Context, requests []CreateRequest) ([]*Response, error) {
	models := make([]*Response, len(requests))
	var err error
	for i, request := range requests {
		models[i], err = t.TransformCreateRequestToModel(ctx, request)
		if err != nil {
			return nil, err
		}
	}

	return models, nil
}

func (t *modelTransformer) TransformBulkUpdateRequestToModels(ctx context.Context, requests []BulkUpdateRequestItem) ([]*Response, error) {
	models := make([]*Response, len(requests))
	var err error
	for i, request := range requests {
		models[i], err = t.TransformCreateRequestToModel(ctx, CreateRequest(request))
		if err != nil {
			return nil, err
		}
	}

	return models, nil
}

func (t *modelTransformer) transformReasonToModel(ctx context.Context, id string) (*apireason.Reason, error) {
	var reason apireason.Reason

	err := t.reasonCollection.
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&reason)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, ErrReasonNotExists
		} else {
			return nil, err
		}
	}

	return &reason, nil
}

func (t *modelTransformer) transformTypeToModel(ctx context.Context, id string) (*pbehavior.Type, error) {
	var pbhType pbehavior.Type

	err := t.typeCollection.
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&pbhType)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, pbehaviorexception.ErrTypeNotExists
		} else {
			return nil, err
		}
	}

	return &pbhType, nil
}

func (t *modelTransformer) transformExceptionsToModel(ctx context.Context, ids []string) ([]pbehaviorexception.Exception, error) {
	if len(ids) == 0 {
		return []pbehaviorexception.Exception{}, nil
	}

	cursor, err := t.exceptionCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": ids}}},
		{"$unwind": "$exdates"},
		{"$lookup": bson.M{
			"from":         pbehavior.TypeCollectionName,
			"localField":   "exdates.type",
			"foreignField": "_id",
			"as":           "exdates.type",
		}},
		{"$unwind": "$exdates.type"},
		{"$group": bson.M{
			"_id":     "$_id",
			"data":    bson.M{"$first": "$$ROOT"},
			"exdates": bson.M{"$push": "$exdates"},
		}},
		{"$replaceRoot": bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				"$data",
				bson.D{{Key: "exdates", Value: "$exdates"}}}},
		}},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	exceptions := make([]pbehaviorexception.Exception, len(ids))
	i := 0
	for ; cursor.Next(ctx); i++ {
		err := cursor.Decode(&exceptions[i])
		if err != nil {
			return nil, err
		}
	}

	if i < len(ids) {
		return nil, ErrExceptionNotExists
	}

	return exceptions, nil
}

func (t *modelTransformer) Patch(ctx context.Context, req PatchRequest, model *Response) error {
	var err error
	model.Author = req.Author
	if req.Enabled != nil {
		model.Enabled = *req.Enabled
	}
	if req.Filter != nil {
		model.Filter = NewFilter(req.Filter)
	}
	if req.Name != nil {
		model.Name = *req.Name
	}
	if req.RRule != nil {
		model.RRule = *req.RRule
	}
	if req.Start != nil {
		model.Start = req.Start
	}
	if req.Stop.isSet {
		model.Stop = req.Stop.CpsTime
	}
	if req.Type != nil {
		var pbhType *pbehavior.Type
		if pbhType, err = t.transformTypeToModel(ctx, *req.Type); err != nil {
			return err
		}
		model.Type = pbhType
	}
	if req.Reason != nil {
		var reason *apireason.Reason
		if reason, err = t.transformReasonToModel(ctx, *req.Reason); err != nil {
			return err
		}
		model.Reason = reason
	}
	if len(req.Exdates) > 0 {
		var exdates []pbehaviorexception.Exdate
		if exdates, err = t.exceptionTransformer.TransformExdatesRequestToModel(ctx, req.Exdates); err != nil {
			return err
		}
		model.Exdates = exdates
	}
	if len(req.Exceptions) > 0 {
		var exceptions []pbehaviorexception.Exception
		if exceptions, err = t.transformExceptionsToModel(ctx, req.Exceptions); err != nil {
			return err
		}
		model.Exceptions = exceptions
	}
	return err
}
