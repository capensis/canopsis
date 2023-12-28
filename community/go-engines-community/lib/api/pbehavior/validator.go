package pbehavior

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorexception"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"github.com/teambition/rrule-go"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Validator struct {
	dbClient mongo.DbClient
}

func NewValidator(client mongo.DbClient) *Validator {
	return &Validator{dbClient: client}
}

func (v *Validator) ValidateCreateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)
	if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required", "")
	}
}

func (v *Validator) ValidateUpdateRequest(sl validator.StructLevel) {
	corporateEntityPattern := ""
	var entityPattern pattern.Entity
	switch r := sl.Current().Interface().(type) {
	case UpdateRequest:
		entityPattern = r.EntityPattern
		corporateEntityPattern = r.CorporateEntityPattern
	case BulkUpdateRequestItem:
		entityPattern = r.EntityPattern
		corporateEntityPattern = r.CorporateEntityPattern
	}

	if len(entityPattern) == 0 && corporateEntityPattern == "" {
		sl.ReportError(entityPattern, "EntityPattern", "EntityPattern", "required", "")
	}
}

func (v *Validator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	if r.CorporateEntityPattern == "" && len(r.EntityPattern) > 0 &&
		!match.ValidateEntityPattern(r.EntityPattern, common.GetForbiddenFieldsInEntityPattern(mongo.PbehaviorMongoCollection)) {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
	}

	if r.RRule != "" && !v.checkRrule(r.RRule) {
		sl.ReportError(r.RRule, "RRule", "RRule", "rrule", "")
	}

	var foundType *pbehavior.Type
	var err error
	if r.Type != "" {
		foundType, err = v.checkType(ctx, r.Type)
		if err != nil {
			panic(err)
		}
		if foundType == nil {
			sl.ReportError(r.Type, "Type", "Type", "not_exist", "")
		}
	}
	if r.Reason != "" {
		ok, err := v.checkReason(ctx, r.Reason)
		if err != nil {
			panic(err)
		}
		if !ok {
			sl.ReportError(r.Reason, "Reason", "Reason", "not_exist", "")
		}
	}
	ok, err := v.checkExdates(ctx, r.Exdates)
	if err != nil {
		panic(err)
	}
	if !ok {
		sl.ReportError(r.Exdates, "Exdates", "Exdates", "not_exist", "")
	}
	ok, err = v.checkExceptions(ctx, r.Exceptions)
	if err != nil {
		panic(err)
	}
	if !ok {
		sl.ReportError(r.Exceptions, "Exceptions", "Exceptions", "not_exist", "")
	}

	if r.Stop != nil && r.Start != nil && r.Stop.Before(*r.Start) {
		sl.ReportError(r.Stop, "Stop", "Stop", "gtfield", "Start")
	}

	if r.Stop == nil && foundType != nil && foundType.Type != pbehavior.TypePause {
		sl.ReportError(r.Stop, "Stop", "Stop", "required", "")
	}
}

func (v *Validator) ValidatePatchRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(PatchRequest)

	if r.RRule != nil && *r.RRule != "" && !v.checkRrule(*r.RRule) {
		sl.ReportError(r.RRule, "RRule", "RRule", "rrule", "")
	}

	if r.CorporateEntityPattern != nil && *r.CorporateEntityPattern == "" {
		sl.ReportError(r.CorporateEntityPattern, "CorporateEntityPattern", "CorporateEntityPattern", "required", "")
	}

	if r.CorporateEntityPattern == nil && r.EntityPattern != nil {
		if len(r.EntityPattern) == 0 {
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required", "")
		} else if !match.ValidateEntityPattern(r.EntityPattern, common.GetForbiddenFieldsInEntityPattern(mongo.PbehaviorMongoCollection)) {
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
		}
	}

	if r.Name != nil && *r.Name == "" {
		sl.ReportError(r.Name, "Name", "Name", "required", "")
	}

	var foundType *pbehavior.Type
	var err error
	if r.Type != nil {
		if *r.Type == "" {
			sl.ReportError(r.Type, "Type", "Type", "required", "")
		} else {
			foundType, err = v.checkType(ctx, *r.Type)
			if err != nil {
				panic(err)
			}
			if foundType == nil {
				sl.ReportError(r.Type, "Type", "Type", "not_exist", "")
			}
		}
	}
	if r.Reason != nil {
		if *r.Reason == "" {
			sl.ReportError(r.Reason, "Reason", "Reason", "required", "")
		} else {
			ok, err := v.checkReason(ctx, *r.Reason)
			if err != nil {
				panic(err)
			}
			if !ok {
				sl.ReportError(r.Reason, "Reason", "Reason", "not_exist", "")
			}
		}
	}
	if r.Exdates != nil {
		ok, err := v.checkExdates(ctx, r.Exdates)
		if err != nil {
			panic(err)
		}
		if !ok {
			sl.ReportError(r.Exdates, "Exdates", "Exdates", "not_exist", "")
		}
	}
	if r.Exceptions != nil {
		ok, err := v.checkExceptions(ctx, r.Exceptions)
		if err != nil {
			panic(err)
		}
		if !ok {
			sl.ReportError(r.Exceptions, "Exceptions", "Exceptions", "not_exist", "")
		}
	}
	if r.Color != nil && *r.Color != "" {
		err := validator.New().Var(*r.Color, "iscolor")
		if err != nil {
			sl.ReportError(r.Color, "Color", "Color", "iscolor", "")
		}
	}

	cursor, err := v.dbClient.Collection(mongo.PbehaviorMongoCollection).Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": r.ID}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"localField":   "type_",
			"foreignField": "_id",
			"as":           "type",
		}},
		{"$unwind": "$type"},
	})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)
	pbh := struct {
		pbehavior.PBehavior `bson:",inline"`
		pbehavior.Type      `bson:"type"`
	}{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&pbh)
		if err != nil {
			panic(err)
		}
	} else {
		return
	}

	if r.Start != nil && r.Stop.isSet {
		if r.Stop.val != nil && *r.Start >= *r.Stop.val {
			sl.ReportError(r.Stop, "Stop", "Stop", "gtfield", "Start")
		}
	} else if r.Start != nil && !r.Stop.isSet {
		if pbh.Stop != nil && *r.Start >= pbh.Stop.Unix() {
			sl.ReportError(r.Start, "Start", "Start", "ltfield", "Stop")
		}
	} else if r.Start == nil && r.Stop.isSet {
		if r.Stop.val != nil && pbh.Start.Unix() >= *r.Stop.val {
			sl.ReportError(r.Stop, "Stop", "Stop", "gtfield", "Start")
		}
	}

	if r.Stop.isSet {
		if r.Stop.val == nil {
			cannonicalType := ""
			if foundType != nil {
				cannonicalType = foundType.Type
			} else {
				cannonicalType = pbh.Type.Type
			}

			if cannonicalType != pbehavior.TypePause {
				sl.ReportError(r.Stop, "Stop", "Stop", "required", "")
			}
		}
	} else if foundType != nil && foundType.Type != pbehavior.TypePause && pbh.Stop == nil {
		sl.ReportError(r.Stop, "Stop", "Stop", "required", "")
	}
}

func (v *Validator) ValidateEntityCreateRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(BulkEntityCreateRequestItem)
	var foundType *pbehavior.Type
	var err error
	if r.Type != "" {
		foundType, err = v.checkType(ctx, r.Type)
		if err != nil {
			panic(err)
		}
		if foundType == nil {
			sl.ReportError(r.Type, "Type", "Type", "not_exist", "")
		}
	}
	if r.Reason != "" {
		ok, err := v.checkReason(ctx, r.Reason)
		if err != nil {
			panic(err)
		}
		if !ok {
			sl.ReportError(r.Reason, "Reason", "Reason", "not_exist", "")
		}
	}

	if r.Stop != nil && r.Start != nil && r.Stop.Before(*r.Start) {
		sl.ReportError(r.Stop, "Stop", "Stop", "gtfield", "Start")
	}

	if r.Stop == nil && foundType != nil && foundType.Type != pbehavior.TypePause {
		sl.ReportError(r.Stop, "Stop", "Stop", "required", "")
	}
}

func (v *Validator) ValidateConnectorCreateRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(BulkConnectorCreateRequestItem)

	if r.Type != "" {
		foundType, err := v.checkType(ctx, r.Type)
		if err != nil {
			panic(err)
		}

		if foundType == nil {
			sl.ReportError(r.Type, "Type", "Type", "not_exist", "")
		}
	}

	if r.Reason != "" {
		ok, err := v.checkReason(ctx, r.Reason)
		if err != nil {
			panic(err)
		}

		if !ok {
			sl.ReportError(r.Reason, "Reason", "Reason", "not_exist", "")
		}
	}

	if r.Stop != nil && r.Start != nil && r.Stop.Before(*r.Start) {
		sl.ReportError(r.Stop, "Stop", "Stop", "gtfield", "Start")
	}
}

func (v *Validator) checkRrule(r string) bool {
	_, err := rrule.StrToROption(r)
	return err == nil
}

func (v *Validator) checkType(ctx context.Context, id string) (*pbehavior.Type, error) {
	if id == "" {
		return nil, nil
	}
	foundType := pbehavior.Type{}
	err := v.dbClient.Collection(mongo.PbehaviorTypeMongoCollection).
		FindOne(ctx, bson.M{"_id": id}).Decode(&foundType)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &foundType, nil
}

func (v *Validator) checkReason(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return true, nil
	}
	err := v.dbClient.Collection(mongo.PbehaviorReasonMongoCollection).
		FindOne(ctx, bson.M{"_id": id}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (v *Validator) checkExdates(ctx context.Context, exdates []pbehaviorexception.ExdateRequest) (bool, error) {
	if len(exdates) == 0 {
		return true, nil
	}
	types := make([]string, 0, len(exdates))
	has := make(map[string]bool, len(exdates))
	for _, v := range exdates {
		if v.Type != "" && !has[v.Type] {
			has[v.Type] = true
			types = append(types, v.Type)
		}
	}

	count, err := v.dbClient.Collection(mongo.PbehaviorTypeMongoCollection).
		CountDocuments(ctx, bson.M{"_id": bson.M{"$in": types}})
	if err != nil {
		return false, err
	}

	return count == int64(len(types)), nil
}

func (v *Validator) checkExceptions(ctx context.Context, exceptions []string) (bool, error) {
	if len(exceptions) == 0 {
		return true, nil
	}
	count, err := v.dbClient.Collection(mongo.PbehaviorExceptionMongoCollection).
		CountDocuments(ctx, bson.M{"_id": bson.M{"$in": exceptions}})
	if err != nil {
		return false, err
	}

	return count == int64(len(exceptions)), nil
}

func (v *Validator) ValidateCalendarRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(CalendarByEntityIDRequest)
	if r.To.Unix() > 0 && r.From.Unix() > 0 && r.To.Before(r.From) {
		sl.ReportError(r.To, "To", "To", "gtfield", "From")
	}
}
