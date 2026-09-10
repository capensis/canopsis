package pbehavior

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/dbvalidation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/sortparam"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/go-playground/validator/v10"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store interface {
	Find(ctx context.Context, r ListRequest) (pagination.ListResponse[Response], error)
	GetByID(ctx context.Context, id string) (Response, error)
	Insert(ctx context.Context, r EditRequest) (Response, error)
	Update(ctx context.Context, id string, r EditRequest) (Response, bool, error)
	UpdateByPatch(ctx context.Context, id string, r PatchRequest) (res Response, recomputeInherited bool, err error)
	Delete(ctx context.Context, id, userID string) (bool, error)
	ConnectorCreate(ctx context.Context, r BulkConnectorCreateRequestItem) (Response, error)
	ConnectorDelete(ctx context.Context, r BulkConnectorDeleteRequestItem) (string, error)
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	typeDbCollection      mongo.DbCollection
	reasonDbCollection    mongo.DbCollection
	exceptionDbCollection mongo.DbCollection

	timezoneConfigProvider config.TimezoneConfigProvider

	transformer    patternfields.Transformer
	dupErrorParser validation.DuplicateErrorParser
}

func NewStore(
	dbClient mongo.DbClient,
	timezoneConfigProvider config.TimezoneConfigProvider,
	transformer patternfields.Transformer,
) Store {
	return &store{
		dbClient:               dbClient,
		dbCollection:           dbClient.Collection(mongo.PbehaviorMongoCollection),
		typeDbCollection:       dbClient.Collection(mongo.PbehaviorTypeMongoCollection),
		reasonDbCollection:     dbClient.Collection(mongo.PbehaviorReasonMongoCollection),
		exceptionDbCollection:  dbClient.Collection(mongo.PbehaviorExceptionMongoCollection),
		timezoneConfigProvider: timezoneConfigProvider,
		transformer:            transformer,
		dupErrorParser:         validation.NewDuplicateErrorParser(),
	}
}

func (s *store) Find(ctx context.Context, r ListRequest) (res pagination.ListResponse[Response], err error) {
	sort, err := sortparam.Validate(cmp.Or(r.Sort, "name"), []string{"_id", "name", "tstart", "tstop", "rrule", "created", "updated"})
	if err != nil {
		return res, err
	}

	filter := bson.M{}
	if r.Enabled != nil {
		filter["enabled"] = *r.Enabled
	}

	if r.Inherited != nil {
		filter["inherited"] = *r.Inherited
	}

	if len(r.Types) > 0 {
		filter["type"] = bson.M{"$in": r.Types}
	}

	if len(r.Reasons) > 0 {
		filter["reason"] = bson.M{"$in": r.Reasons}
	}

	if len(r.Origins) > 0 {
		filter["origin"] = bson.M{"$in": r.Origins}
	}

	if r.Author != "" {
		filter["author"] = r.Author
	}

	beforeLimit := make([]bson.M, 0)
	if len(filter) != 0 {
		beforeLimit = append(beforeLimit, bson.M{"$match": filter})
	}

	afterLimit := nestedObjectsPipeline()
	afterLimit = append(afterLimit, sort)
	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		sort,
		afterLimit,
	)
	opts := options.Aggregate().SetCollation(&options.Collation{Locale: "en"})
	ar, err := db.AggregateAll[Response](ctx, s.dbCollection, pipeline, opts)
	if err != nil {
		return res, err
	}

	res = pagination.NewResponse[Response](r.Query, ar)

	return res, nil
}

func (s *store) GetByID(ctx context.Context, id string) (Response, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
	}
	pipeline = append(pipeline, nestedObjectsPipeline()...)

	res, err := db.AggregateOne[Response](ctx, s.dbCollection, pipeline)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *store) Insert(ctx context.Context, r EditRequest) (Response, error) {
	var res Response
	now := datetime.NewCpsTime()
	doc := s.transformRequestToDocument(r)
	doc.ID = utils.NewID()

	loc := s.timezoneConfigProvider.Get().Location
	var err error
	if doc.Timezone != "" {
		loc, err = time.LoadLocation(doc.Timezone)
		if err != nil {
			return res, fmt.Errorf("invalid timezone %q: %w", doc.Timezone, err)
		}
	}

	rruleEnd, err := pbehavior.GetRruleEnd(*r.Start, r.RRule, loc)
	if err != nil {
		return res, err
	}

	doc.Created = &now
	doc.Updated = &now
	doc.Comments = make([]pbehavior.Comment, 0)
	doc.RRuleEnd = rruleEnd

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = Response{}

		err = s.validateEditRequest(ctx, r)
		if err != nil {
			return err
		}

		doc.EntityPatternFields, doc.Aliases, err = s.transformer.TransformEntityRequest(ctx, patternfields.EntityRequest{
			EntityPattern: r.EntityPattern,
		}, r, s.dbCollection.Name())
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

		res, err = s.GetByID(ctx, doc.ID)

		return err
	})

	return res, err
}

func (s *store) Update(ctx context.Context, id string, r EditRequest) (res Response, recomputeInherited bool, err error) {
	now := datetime.NewCpsTime()
	doc := s.transformRequestToDocument(r)
	doc.Updated = &now

	unset := bson.M{
		"rrule_cstart": "",
	}
	if r.Stop == nil {
		unset["tstop"] = ""
	}

	loc := s.timezoneConfigProvider.Get().Location
	if doc.Timezone != "" {
		loc, err = time.LoadLocation(doc.Timezone)
		if err != nil {
			return res, false, fmt.Errorf("invalid timezone %q: %w", doc.Timezone, err)
		}
	}

	rruleEnd, err := pbehavior.GetRruleEnd(*r.Start, r.RRule, loc)
	if err != nil {
		return res, false, err
	}

	if rruleEnd == nil {
		unset["rrule_end"] = ""
	} else {
		doc.RRuleEnd = rruleEnd
	}

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = Response{}
		recomputeInherited = false

		err = s.validateEditRequest(ctx, r)
		if err != nil {
			return err
		}

		var prevPbh pbehavior.PBehavior
		err = s.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&prevPbh)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return httperror.ErrNotFound
			}

			return err
		}

		err = s.validateUpdateRequest(r, prevPbh)
		if err != nil {
			return err
		}

		doc.EntityPatternFields, doc.Aliases, err = s.transformer.TransformEntityRequest(ctx, patternfields.EntityRequest{
			EntityPattern: r.EntityPattern,
		}, r, s.dbCollection.Name())
		if err != nil {
			return err
		}

		_, err = s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
			"$set":   doc,
			"$unset": unset,
		})
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Response{})
			}

			return err
		}

		res, err = s.GetByID(ctx, id)
		if err != nil {
			return err
		}

		recomputeInherited = prevPbh.Inherited || res.Inherited

		return nil
	})

	return res, recomputeInherited, err
}

func (s *store) UpdateByPatch(ctx context.Context, id string, r PatchRequest) (res Response, recomputeInherited bool, err error) {
	set := bson.M{
		"author":  r.Author,
		"updated": datetime.NewCpsTime(),
	}

	unset := bson.M{}
	rruleUpdated := false
	if r.Name != nil {
		set["name"] = *r.Name
	}
	if r.Enabled != nil {
		set["enabled"] = *r.Enabled
	}
	if r.Reason != nil {
		set["reason"] = *r.Reason
	}
	if r.Type != nil {
		set["type"] = *r.Type
	}
	if r.RRule != nil {
		set["rrule"] = *r.RRule
		rruleUpdated = true
	}
	if r.Timezone != nil {
		set["timezone"] = *r.Timezone
		rruleUpdated = true
	}
	if r.Start != nil {
		set["tstart"] = *r.Start
		rruleUpdated = true
	}
	if r.Stop.isSet {
		if r.Stop.val == nil {
			unset["tstop"] = ""
		} else {
			set["tstop"] = *r.Stop.val
		}
	}
	if r.Exdates != nil {
		exdates := make([]pbehavior.Exdate, len(r.Exdates))
		for i := range r.Exdates {
			exdates[i].Type = r.Exdates[i].Type
			exdates[i].Begin = r.Exdates[i].Begin
			exdates[i].End = r.Exdates[i].End
		}

		set["exdates"] = exdates
	}
	if r.Color != nil {
		set["color"] = *r.Color
	}

	if r.Inherited != nil {
		set["inherited"] = *r.Inherited
	}

	if rruleUpdated {
		unset["rrule_cstart"] = ""
	}

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = Response{}
		recomputeInherited = false

		prevPbh := struct {
			pbehavior.PBehavior `bson:",inline"`
			pbehavior.Type      `bson:"type"`
		}{}
		cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
			{"$match": bson.M{"_id": id}},
			{"$lookup": bson.M{
				"from":         mongo.PbehaviorTypeMongoCollection,
				"localField":   "type",
				"foreignField": "_id",
				"as":           "type",
			}},
			{"$unwind": "$type"},
		})
		if err != nil {
			return err
		}

		defer cursor.Close(ctx)
		if !cursor.Next(ctx) {
			return httperror.ErrNotFound
		}

		err = cursor.Decode(&prevPbh)
		if err != nil {
			return err
		}

		if err = cursor.Err(); err != nil {
			return err
		}

		err = s.validatePatchRequest(ctx, r, prevPbh.PBehavior, prevPbh.Type)
		if err != nil {
			return err
		}

		if r.EntityPattern != nil {
			epf, aliases, err := s.transformer.TransformEntityRequest(ctx, patternfields.EntityRequest{
				EntityPattern: r.EntityPattern,
			}, r, s.dbCollection.Name())
			if err != nil {
				return err
			}

			set["aliases"] = aliases
			set["entity_pattern"] = epf.EntityPattern
			set["corporate_entity_pattern"] = ""
			set["corporate_entity_pattern_title"] = ""
		}

		update := bson.M{"$set": set}
		if len(unset) > 0 {
			update["$unset"] = unset
		}

		_, err = s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Response{})
			}

			return err
		}

		res, err = s.GetByID(ctx, id)
		if err != nil {
			return err
		}

		if rruleUpdated {
			loc := s.timezoneConfigProvider.Get().Location
			if res.Timezone != "" {
				loc, err = time.LoadLocation(res.Timezone)
				if err != nil {
					return fmt.Errorf("invalid timezone %q: %w", res.Timezone, err)
				}
			}

			rruleEnd, err := pbehavior.GetRruleEnd(*res.Start, res.RRule, loc)
			if err != nil {
				return err
			}

			if rruleEnd == nil {
				_, err = s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$unset": bson.M{"rrule_end": ""}})
			} else {
				_, err = s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"rrule_end": rruleEnd}})
			}
			if err != nil {
				return err
			}
		}

		recomputeInherited = prevPbh.Inherited || res.Inherited

		return nil
	})

	return res, recomputeInherited, err
}

func (s *store) Delete(ctx context.Context, id, userID string) (recomputeInherited bool, err error) {
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		recomputeInherited = false
		// required to get the author in action log listener.
		var prevPbh pbehavior.PBehavior
		err := s.dbCollection.
			FindOneAndUpdate(ctx,
				bson.M{"_id": id},
				bson.M{"$set": bson.M{"author": userID}},
				options.FindOneAndUpdate().SetReturnDocument(options.Before).SetProjection(bson.M{"inherited": 1}),
			).
			Decode(&prevPbh)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return httperror.ErrNotFound
			}

			return err
		}

		deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			return err
		}

		if deleted == 0 {
			return httperror.ErrNotFound
		}

		recomputeInherited = prevPbh.Inherited

		return nil
	})

	return recomputeInherited, err
}

func (s *store) ConnectorCreate(ctx context.Context, r BulkConnectorCreateRequestItem) (Response, error) {
	entities := make([]string, 0, len(r.Entities))
	exists := make(map[string]struct{}, len(r.Entities))
	for _, entity := range r.Entities {
		if _, ok := exists[entity]; !ok {
			exists[entity] = struct{}{}
			entities = append(entities, entity)
		}
	}

	now := datetime.NewCpsTime()
	doc := pbehavior.PBehavior{
		ID:       utils.NewID(),
		Author:   r.Author,
		Enabled:  true,
		Name:     r.Name,
		Reason:   r.Reason,
		Start:    r.Start,
		Stop:     r.Stop,
		Type:     r.Type,
		Created:  &now,
		Updated:  &now,
		Color:    r.Color,
		Origin:   r.Origin,
		Entities: entities,
		EntityPatternFields: savedpattern.EntityPatternFields{
			EntityPattern: pattern.Entity{
				{
					{
						Field: "_id",
						Condition: pattern.Condition{
							Type:  pattern.ConditionIsOneOf,
							Value: entities,
						},
					},
				},
			},
		},
		Comments: []pbehavior.Comment{
			{
				ID:        utils.NewID(),
				Author:    r.Author,
				Origin:    r.Origin,
				Timestamp: &now,
				Message:   r.Comment,
			},
		},
	}

	var res Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = Response{}

		err := dbvalidation.ValidateExist(ctx, s.typeDbCollection, r, "Type", r.Type)
		if err != nil {
			return err
		}

		err = dbvalidation.ValidateExist(ctx, s.reasonDbCollection, r, "Reason", r.Reason)
		if err != nil {
			return err
		}

		var findDoc pbehavior.PBehavior
		err = s.dbCollection.FindOneAndUpdate(ctx,
			bson.M{
				"origin":           r.Origin,
				"tstart":           r.Start,
				"tstop":            r.Stop,
				"comments.message": r.Comment,
				"entities":         bson.M{"$ne": nil},
			},
			[]bson.M{
				{"$set": bson.M{
					"entities": bson.M{"$cond": bson.M{
						"if": "$_id",
						"then": bson.M{"$setUnion": bson.A{
							"$entities",
							r.Entities,
						}},
						"else": nil,
					}},
				}},
				{"$set": bson.M{
					"entity_pattern": bson.M{"$cond": bson.M{
						"if": "$_id",
						"then": bson.A{bson.A{bson.M{
							"field": "_id",
							"cond": bson.M{
								"type":  pattern.ConditionIsOneOf,
								"value": "$entities",
							},
						}}},
						"else": nil,
					}},
				}},
				{
					"$replaceRoot": bson.M{"newRoot": bson.M{
						"$cond": bson.M{
							"if": "$_id",
							"then": bson.M{"$mergeObjects": bson.A{
								bson.M{
									"entity_pattern": "$entity_pattern",
									"entities":       "$entities",
									"author":         r.Author,
									"updated":        now,
								},
								"$$ROOT",
							}},
							"else": bson.M{"$literal": doc},
						},
					}},
				},
			},
			options.FindOneAndUpdate().
				SetUpsert(true).
				SetReturnDocument(options.After).
				SetProjection(bson.M{"_id": 1}),
		).Decode(&findDoc)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Response{})
			}

			return err
		}

		res, err = s.GetByID(ctx, findDoc.ID)

		return err
	})

	return res, err
}

func (s *store) ConnectorDelete(ctx context.Context, r BulkConnectorDeleteRequestItem) (string, error) {
	now := datetime.NewCpsTime()
	id := ""
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		id = ""
		var pbh pbehavior.PBehavior
		err := s.dbCollection.FindOneAndUpdate(ctx,
			bson.M{
				"origin":           r.Origin,
				"tstart":           r.Start,
				"tstop":            r.Stop,
				"comments.message": r.Comment,
				"entities":         bson.M{"$ne": nil},
			},
			bson.M{
				"$pullAll": bson.M{
					"entity_pattern.0.0.cond.value": r.Entities,
					"entities":                      r.Entities,
				},
				"$set": bson.M{
					"author":  r.Author,
					"updated": now,
				},
			},
			options.FindOneAndUpdate().SetProjection(bson.M{"_id": 1}),
		).Decode(&pbh)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return httperror.ErrNotFound
			}

			return err
		}

		id = pbh.ID
		// Delete if entities are empty
		_, err = s.dbCollection.DeleteOne(ctx, bson.M{
			"_id":      id,
			"entities": bson.A{},
		})

		return err
	})

	return id, err
}

func (s *store) transformRequestToDocument(r EditRequest) pbehavior.PBehavior {
	exdates := make([]pbehavior.Exdate, len(r.Exdates))
	for i := range r.Exdates {
		exdates[i].Type = r.Exdates[i].Type
		exdates[i].Begin = r.Exdates[i].Begin
		exdates[i].End = r.Exdates[i].End
	}

	return pbehavior.PBehavior{
		Author:    r.Author,
		Enabled:   *r.Enabled,
		Name:      r.Name,
		Reason:    r.Reason,
		RRule:     r.RRule,
		Timezone:  r.Timezone,
		Start:     r.Start,
		Stop:      r.Stop,
		Type:      r.Type,
		Exdates:   exdates,
		Color:     r.Color,
		Inherited: r.Inherited,
		EntityPatternFields: savedpattern.EntityPatternFields{
			EntityPattern: r.EntityPattern,
		},
	}
}

func (s *store) validateEditRequest(ctx context.Context, r EditRequest) error {
	t, err := s.findType(ctx, r.Type)
	if err != nil {
		return err
	}

	if t.ID == "" {
		return validation.NewSingleError("not_exist", "Type", "Type", r)
	}

	if r.Stop == nil && t.Type != pbehavior.TypePause {
		return validation.NewSingleError("required", "Stop", "Stop", r)
	}

	err = dbvalidation.ValidateExist(ctx, s.reasonDbCollection, r, "Reason", r.Reason)
	if err != nil {
		return err
	}

	valErrs, err := s.validateExdates(ctx, r.Exdates)
	if err != nil {
		return err
	}

	if len(valErrs) > 0 {
		return validation.NewError(valErrs, r)
	}

	return nil
}

func (s *store) validateUpdateRequest(r EditRequest, prevPbh pbehavior.PBehavior) error {
	if prevPbh.Origin == "" {
		return nil
	}

	if len(prevPbh.Entities) > 0 {
		return httperror.NewForbiddenError("The external pbehavior cannot be modified.")
	}

	var fieldErrs validator.ValidationErrors
	check := func(field string, cond bool) {
		if cond {
			fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", field, field))
		}
	}

	check("Enabled", !*r.Enabled)
	check("RRule", r.RRule != "")
	check("Timezone", r.Timezone != "")
	check("Exdates", len(r.Exdates) > 0)
	if diff := pretty.Compare(prevPbh.EntityPattern, r.EntityPattern); diff != "" {
		check("EntityPattern", true)
	}

	if len(fieldErrs) > 0 {
		return validation.NewError(fieldErrs, r)
	}

	return nil
}

func (s *store) validatePatchRequest(ctx context.Context, r PatchRequest, prevPbh pbehavior.PBehavior, prevType pbehavior.Type) error {
	var newType pbehavior.Type
	var err error
	if r.Type != nil {
		newType, err = s.findType(ctx, *r.Type)
		if err != nil {
			return err
		}

		if newType.ID == "" {
			return validation.NewSingleError("not_exist", "Type", "Type", r)
		}
	}

	err = dbvalidation.ValidateExist(ctx, s.reasonDbCollection, r, "Reason", r.Reason)
	if err != nil {
		return err
	}

	valErrs, err := s.validateExdates(ctx, r.Exdates)
	if err != nil {
		return err
	}

	if len(valErrs) > 0 {
		return validation.NewError(valErrs, r)
	}

	stopValErr := validation.NewSingleError("required", "Stop", "Stop", r)
	if r.Stop.isSet {
		if r.Stop.val == nil {
			cannonicalType := ""
			if newType.ID != "" {
				cannonicalType = newType.Type
			} else {
				cannonicalType = prevType.Type
			}

			if cannonicalType != pbehavior.TypePause {
				return stopValErr
			}
		}
	} else if newType.ID != "" && newType.Type != pbehavior.TypePause && prevPbh.Stop == nil {
		return stopValErr
	}

	if r.Start != nil && !r.Stop.isSet && prevPbh.Stop != nil && *r.Start >= prevPbh.Stop.Unix() {
		return validation.NewSingleErrorWithParam("ltfield", "Start", "Start", "Stop", r)
	} else if r.Start == nil && r.Stop.isSet && r.Stop.val != nil && prevPbh.Start.Unix() >= *r.Stop.val {
		return validation.NewSingleErrorWithParam("gtfield", "Stop", "Stop", "Start", r)
	}

	if prevPbh.Origin != "" {
		if len(prevPbh.Entities) > 0 {
			return httperror.NewForbiddenError("The external pbehavior cannot be modified.")
		}

		var fieldErrs validator.ValidationErrors
		check := func(field string, cond bool) {
			if cond {
				fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", field, field))
			}
		}

		check("Enabled", r.Enabled != nil && !*r.Enabled)
		check("RRule", r.RRule != nil && *r.RRule != "")
		check("Timezone", r.Timezone != nil && *r.Timezone != "")
		check("Exdates", len(r.Exdates) > 0)
		if r.EntityPattern != nil {
			if diff := pretty.Compare(prevPbh.EntityPattern, r.EntityPattern); diff != "" {
				check("EntityPattern", true)
			}
		}

		if len(fieldErrs) > 0 {
			return validation.NewError(fieldErrs, r)
		}
	}

	return nil
}

func (s *store) findType(ctx context.Context, id string) (pbehavior.Type, error) {
	t := pbehavior.Type{}
	err := s.typeDbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return t, nil
		}

		return t, err
	}

	return t, nil
}

func (s *store) validateExdates(ctx context.Context, exdates []ExdateRequest) (validator.ValidationErrors, error) {
	if len(exdates) == 0 {
		return nil, nil
	}

	ids := make([]string, len(exdates))
	for i, v := range exdates {
		ids[i] = v.Type
	}

	cursor, err := s.typeDbCollection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, fmt.Errorf("cannot find types: %w", err)
	}

	defer cursor.Close(ctx)
	found := make(map[string]bool, len(ids))
	for cursor.Next(ctx) {
		t := struct {
			ID string `bson:"_id"`
		}{}
		err = cursor.Decode(&t)
		if err != nil {
			return nil, fmt.Errorf("cannot decode type: %w", err)
		}

		found[t.ID] = true
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("cannot fetch types: %w", err)
	}

	var fieldErrs validator.ValidationErrors
	for i, v := range exdates {
		if !found[v.Type] {
			fieldErrs = append(fieldErrs, validation.NewFieldError("not_exist", "Type", "Exdates."+strconv.Itoa(i)+".Type"))
		}
	}

	return fieldErrs, nil
}

func nestedObjectsPipeline() []bson.M {
	pipeline := mongoquery.ReferencePipeline("type", mongo.PbehaviorTypeMongoCollection)
	pipeline = append(pipeline, mongoquery.ReferencePipeline("reason", mongo.PbehaviorReasonMongoCollection)...)
	pipeline = append(pipeline, exdatesPipeline()...)
	pipeline = append(pipeline, mongoquery.ArrayReferencePipeline("exceptions", mongo.PbehaviorExceptionMongoCollection)...)
	pipeline = append(pipeline, mongoquery.AuthorPipeline()...)
	pipeline = append(pipeline, originCommentPipeline()...)

	return pipeline
}

func exdatesPipeline() []bson.M {
	pipeline := []bson.M{
		{"$unwind": bson.M{
			"path":                       "$exdates",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "exdate_index",
		}},
	}

	pipeline = append(pipeline, mongoquery.ReferencePipeline("exdates.type", mongo.PbehaviorTypeMongoCollection)...)
	pipeline = append(pipeline,
		bson.M{"$sort": bson.M{"exdate_index": 1}},
		bson.M{"$group": bson.M{
			"_id":  "$_id",
			"data": bson.M{"$first": "$$ROOT"},
			"exdates": bson.M{"$push": bson.M{"$cond": bson.M{
				"if":   "$exdates.begin",
				"then": "$exdates",
				"else": "$$REMOVE",
			}}},
		}},
		bson.M{"$replaceRoot": bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				"$data",
				bson.M{"exdates": "$exdates"},
			}},
		}},
	)

	return pipeline
}

func originCommentPipeline() []bson.M {
	return []bson.M{
		{"$addFields": bson.M{
			"origin_comment": bson.M{"$cond": bson.M{
				"if": "$origin",
				"then": bson.M{"$first": bson.M{
					"$map": bson.M{
						"input": "$comments",
						"in":    "$$this.message",
					},
				}},
				"else": "",
			}},
		}},
	}
}
