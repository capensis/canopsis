package patternfields

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Transformer defines the behavior for enriching incoming request payloads
// (AlarmRequest, EntityRequest, PbehaviorRequest, WeatherServiceRequest)
// with corporate pattern data, user-scoped pattern data, and alias entity properties
// retrieved from the database.
//
// -----------------------------------------------------------------------------
// EXAMPLE 1: Enriching AlarmRequest with corporate patterns
// -----------------------------------------------------------------------------
//
//	type Request struct {
//	    patternfields.AlarmRequest
//	}
//
//	f := func(i Request) (savedpattern.AlarmPatternFields, error) {
//	    patterns, err := s.transformer.FetchCorporatePatterns(ctx, r.CorporateAlarmPattern)
//	    if err != nil {
//	        return savedpattern.AlarmPatternFields{}, err
//	    }
//
//	    var valErrs validator.ValidationErrors
//	    r.AlarmRequest, valErrs = s.transformer.ApplyAlarmCorporatePattern(r.AlarmRequest, patterns)
//	    if valErrs != nil {
//	        return savedpattern.AlarmPatternFields{}, validation.NewError(valErrs, r)
//	    }
//
//	    return r.AlarmRequest.ToModel(), nil
//	}
//
// -----------------------------------------------------------------------------
// EXAMPLE 2: Enriching EntityRequest with corporate patterns
// -----------------------------------------------------------------------------
//
//	type Request struct {
//	    patternfields.EntityRequest
//	}
//
//	f := func(i Request) (savedpattern.EntityPatternFields, []string, error) {
//	    patterns, err := s.transformer.FetchCorporatePatterns(ctx, r.CorporateEntityPattern)
//	    if err != nil {
//	        return savedpattern.EntityPatternFields{}, nil, err
//	    }
//
//	    aliases, err := s.transformer.FetchAliases(ctx, patternfields.GetAliases(fr.EntityPattern))
//	    if err != nil {
//	        return savedpattern.EntityPatternFields{}, nil, err
//	    }
//
//	    var aliasPropIDs []string
//	    var valErrs validator.ValidationErrors
//	    if r.CorporateEntityPattern != "" {
//	        r.EntityRequest, aliasPropIDs, valErrs = s.transformer.ApplyEntityCorporatePattern(r.EntityRequest, patterns)
//	    } else if r.EntityPattern != nil {
//	        r.EntityPattern, aliasPropIDs, valErrs = s.transformer.ApplyAliases(r.EntityPattern, aliases)
//	    }
//
//	    if valErrs != nil {
//	        return savedpattern.EntityPatternFields{}, nil, validation.NewError(valErrs, r)
//	    }
//
//	    return r.EntityRequest.ToModel(), aliasPropIDs, nil
//	}
//
// -----------------------------------------------------------------------------
// EXAMPLE 3: Enriching AlarmRequest and EntityRequest with corporate patterns
// -----------------------------------------------------------------------------
//
//	type Request struct {
//	    patternfields.AlarmRequest
//	    patternfields.EntityRequest
//	}
//
//	f := func(i Request) (apf savedpattern.AlarmPatternFields, epf savedpattern.EntityPatternFields, aliasPropIDs []string, err error) {
//	    return s.transformer.TransformAlarmAndEntityRequest(ctx, r.AlarmRequest, r.EntityRequest, r, s.collection.Name())
//	}
type Transformer interface {
	FetchCorporatePatterns(ctx context.Context, patternIDs ...string) (Patterns, error)
	FetchPatternsByUser(ctx context.Context, userID string, patternIDs ...string) (Patterns, error)
	FetchAliases(ctx context.Context, aliases []string) (Aliases, error)

	ApplyAlarmCorporatePattern(r AlarmRequest, patterns Patterns, fieldNs ...string) (AlarmRequest, validator.ValidationErrors)
	ApplyEntityCorporatePattern(r EntityRequest, patterns Patterns, fieldNs ...string) (EntityRequest, []string, validator.ValidationErrors)
	ApplyPbehaviorCorporatePattern(r PbehaviorRequest, patterns Patterns, fieldNs ...string) (PbehaviorRequest, validator.ValidationErrors)
	ApplyServiceWeatherCorporatePattern(r WeatherServiceRequest, patterns Patterns, fieldNs ...string) (WeatherServiceRequest, validator.ValidationErrors)

	ApplyAliases(p pattern.Entity, aliasProps Aliases, fieldNs ...string) (pattern.Entity, []string, validator.ValidationErrors)

	TransformEntityRequest(
		ctx context.Context,
		er EntityRequest,
		r any,
		collectionName string,
	) (epf savedpattern.EntityPatternFields, aliasPropIDs []string, err error)
	TransformAlarmAndEntityRequest(
		ctx context.Context,
		ar AlarmRequest,
		er EntityRequest,
		r any,
		collectionName string,
	) (apf savedpattern.AlarmPatternFields, epf savedpattern.EntityPatternFields, aliasPropIDs []string, err error)
	TransformAliases(ctx context.Context, p pattern.Entity, r any) (pattern.Entity, []string, error)
}

type Patterns map[string]savedpattern.SavedPattern

type Aliases map[string]Alias

func (a Aliases) GetPropIDs() []string {
	ids := make([]string, 0, len(a))
	for _, v := range a {
		ids = append(ids, v.ID)
	}

	return ids
}

type Alias struct {
	ID   string
	Name string
}

func NewTransformer(client mongo.DbClient) Transformer {
	return &baseTransformer{
		patternCollection:             client.Collection(mongo.PatternMongoCollection),
		entityInfosPropertyCollection: client.Collection(mongo.EntityInfosPropertyCollection),
	}
}

type baseTransformer struct {
	patternCollection             mongo.DbCollection
	entityInfosPropertyCollection mongo.DbCollection
}

func (t *baseTransformer) FetchCorporatePatterns(ctx context.Context, patternIDs ...string) (Patterns, error) {
	k := 0
	for _, d := range patternIDs {
		if d != "" {
			patternIDs[k] = d
			k++
		}
	}

	patternIDs = patternIDs[:k]
	if len(patternIDs) == 0 {
		return nil, nil
	}

	filter := bson.M{
		"_id":          bson.M{"$in": patternIDs},
		"is_corporate": true,
	}

	return t.fetchPatternsByFilter(ctx, filter)
}

func (t *baseTransformer) FetchPatternsByUser(ctx context.Context, userID string, patternIDs ...string) (Patterns, error) {
	k := 0
	for _, d := range patternIDs {
		if d != "" {
			patternIDs[k] = d
			k++
		}
	}

	patternIDs = patternIDs[:k]
	if len(patternIDs) == 0 {
		return nil, nil
	}

	filter := bson.M{
		"_id": bson.M{"$in": patternIDs},
		"$or": []bson.M{
			{"is_corporate": true},
			{"author": userID},
		},
	}

	return t.fetchPatternsByFilter(ctx, filter)
}

func (t *baseTransformer) FetchAliases(ctx context.Context, aliases []string) (Aliases, error) {
	cursor, err := t.entityInfosPropertyCollection.Find(
		ctx,
		bson.M{"alias": bson.M{"$in": aliases}},
		options.Find().SetProjection(bson.M{"alias": 1, "name": 1}),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot find entity aliases: %w", err)
	}

	defer cursor.Close(ctx)
	aliasProps := make(Aliases, len(aliases))
	for cursor.Next(ctx) {
		prop := struct {
			ID    string `bson:"_id"`
			Alias string `bson:"alias"`
			Name  string `bson:"name"`
		}{}
		err = cursor.Decode(&prop)
		if err != nil {
			return nil, fmt.Errorf("cannot decode entity alias: %w", err)
		}

		aliasProps[prop.Alias] = Alias{
			ID:   prop.ID,
			Name: prop.Name,
		}
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("cannot fetch entity aliases: %w", err)
	}

	return aliasProps, nil
}

func (t *baseTransformer) ApplyAlarmCorporatePattern(r AlarmRequest, patterns Patterns, fieldNs ...string) (AlarmRequest, validator.ValidationErrors) {
	if r.CorporateAlarmPattern != "" {
		p, ok := patterns[r.CorporateAlarmPattern]
		if !ok || p.Type != savedpattern.TypeAlarm {
			if len(fieldNs) == 0 {
				fieldNs = []string{"CorporateAlarmPattern"}
			}

			return r, validator.ValidationErrors{t.newNotExistErr(fieldNs)}
		}

		r.CorporatePattern = p
	}

	return r, nil
}

func (t *baseTransformer) ApplyEntityCorporatePattern(r EntityRequest, patterns Patterns, fieldNs ...string) (EntityRequest, []string, validator.ValidationErrors) {
	if r.CorporateEntityPattern == "" {
		return r, nil, nil
	}

	p, ok := patterns[r.CorporateEntityPattern]
	if !ok || p.Type != savedpattern.TypeEntity {
		if len(fieldNs) == 0 {
			fieldNs = []string{"CorporateEntityPattern"}
		}

		return r, nil, validator.ValidationErrors{t.newNotExistErr(fieldNs)}
	}

	r.CorporatePattern = p

	return r, p.Aliases, nil
}

func (t *baseTransformer) ApplyAliases(p pattern.Entity, aliasProps Aliases, fieldNs ...string) (pattern.Entity, []string, validator.ValidationErrors) {
	if p == nil {
		return p, nil, nil
	}

	if len(fieldNs) == 0 {
		fieldNs = []string{"EntityPattern"}
	}

	usedAliasProps := make(Aliases)
	valErrs := make(validator.ValidationErrors, 0)
	for gi, g := range p {
		for ci, c := range g {
			if c.Alias == "" {
				continue
			}

			n, ok := aliasProps[c.Alias]
			if !ok {
				aliasNs := make([]string, len(fieldNs), len(fieldNs)+3)
				copy(aliasNs, fieldNs)
				aliasNs = append(aliasNs, strconv.Itoa(gi), strconv.Itoa(ci), "Alias")
				valErrs = append(valErrs, t.newNotExistErr(aliasNs))
				continue
			}

			usedAliasProps[c.Alias] = n
			g[ci].Field = "infos." + n.Name
		}
	}

	return p, usedAliasProps.GetPropIDs(), valErrs
}

func (t *baseTransformer) ApplyPbehaviorCorporatePattern(r PbehaviorRequest, patterns Patterns, fieldNs ...string) (PbehaviorRequest, validator.ValidationErrors) {
	if r.CorporatePbehaviorPattern != "" {
		p, ok := patterns[r.CorporatePbehaviorPattern]
		if !ok || p.Type != savedpattern.TypePbehavior {
			if len(fieldNs) == 0 {
				fieldNs = []string{"CorporatePbehaviorPattern"}
			}

			return r, validator.ValidationErrors{t.newNotExistErr(fieldNs)}
		}

		r.CorporatePattern = p
	}

	return r, nil
}

func (t *baseTransformer) ApplyServiceWeatherCorporatePattern(r WeatherServiceRequest, patterns Patterns, fieldNs ...string) (WeatherServiceRequest, validator.ValidationErrors) {
	if r.CorporateWeatherServicePattern != "" {
		p, ok := patterns[r.CorporateWeatherServicePattern]
		if !ok || p.Type != savedpattern.TypeWeatherService {
			if len(fieldNs) == 0 {
				fieldNs = []string{"CorporateWeatherServicePattern"}
			}

			return r, validator.ValidationErrors{t.newNotExistErr(fieldNs)}
		}

		r.CorporatePattern = p
	}

	return r, nil
}

func (t *baseTransformer) TransformEntityRequest(
	ctx context.Context,
	er EntityRequest,
	r any,
	collectionName string,
) (epf savedpattern.EntityPatternFields, aliasPropIDs []string, err error) {
	var valErrs validator.ValidationErrors
	if er.CorporateEntityPattern != "" {
		var patterns Patterns
		patterns, err = t.FetchCorporatePatterns(ctx, er.CorporateEntityPattern)
		if err != nil {
			return epf, aliasPropIDs, err
		}

		er, aliasPropIDs, valErrs = t.ApplyEntityCorporatePattern(er, patterns)
	} else if er.EntityPattern != nil {
		var aliases Aliases
		aliases, err = t.FetchAliases(ctx, GetAliases(er.EntityPattern))
		if err != nil {
			return epf, aliasPropIDs, err
		}

		er.EntityPattern, aliasPropIDs, valErrs = t.ApplyAliases(er.EntityPattern, aliases)
	}

	if len(valErrs) > 0 {
		return epf, aliasPropIDs, validation.NewError(valErrs, r)
	}

	epf = er.ToModelWithoutFields(collectionName)

	return epf, aliasPropIDs, nil
}

func (t *baseTransformer) TransformAlarmAndEntityRequest(
	ctx context.Context,
	ar AlarmRequest,
	er EntityRequest,
	r any,
	collectionName string,
) (apf savedpattern.AlarmPatternFields, epf savedpattern.EntityPatternFields, aliasPropIDs []string, err error) {
	patterns, err := t.FetchCorporatePatterns(ctx,
		ar.CorporateAlarmPattern,
		er.CorporateEntityPattern,
	)
	if err != nil {
		return apf, epf, aliasPropIDs, err
	}

	var valErrs, applyErrs validator.ValidationErrors
	ar, applyErrs = t.ApplyAlarmCorporatePattern(ar, patterns)
	valErrs = append(valErrs, applyErrs...)
	if er.CorporateEntityPattern != "" {
		er, aliasPropIDs, applyErrs = t.ApplyEntityCorporatePattern(er, patterns)
	} else if er.EntityPattern != nil {
		var aliases Aliases
		aliases, err = t.FetchAliases(ctx, GetAliases(er.EntityPattern))
		if err != nil {
			return apf, epf, aliasPropIDs, err
		}

		er.EntityPattern, aliasPropIDs, applyErrs = t.ApplyAliases(er.EntityPattern, aliases)
	}

	valErrs = append(valErrs, applyErrs...)
	if len(valErrs) > 0 {
		return apf, epf, aliasPropIDs, validation.NewError(valErrs, r)
	}

	apf = ar.ToModelWithoutFields(collectionName)
	epf = er.ToModelWithoutFields(collectionName)

	return apf, epf, aliasPropIDs, nil
}

func (t *baseTransformer) TransformAliases(ctx context.Context, p pattern.Entity, r any) (pattern.Entity, []string, error) {
	if p == nil {
		return p, nil, nil
	}

	aliases, err := t.FetchAliases(ctx, GetAliases(p))
	if err != nil {
		return p, nil, err
	}

	var aliasPropIDs []string
	var valErrs validator.ValidationErrors
	p, aliasPropIDs, valErrs = t.ApplyAliases(p, aliases)
	if len(valErrs) > 0 {
		return p, nil, validation.NewError(valErrs, r)
	}

	return p, aliasPropIDs, nil
}

func (t *baseTransformer) fetchPatternsByFilter(ctx context.Context, filter bson.M) (Patterns, error) {
	cursor, err := t.patternCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("cannot find corporate patterns: %w", err)
	}

	defer cursor.Close(ctx)
	patterns := make(Patterns)
	for cursor.Next(ctx) {
		p := savedpattern.SavedPattern{}
		err = cursor.Decode(&p)
		if err != nil {
			return nil, fmt.Errorf("cannot decode corporate pattern: %w", err)
		}

		patterns[p.ID] = p
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("cannot fetch corporate patterns: %w", err)
	}

	return patterns, nil
}

func (t *baseTransformer) newNotExistErr(fieldNs []string) validator.FieldError {
	return validation.NewFieldError("not_exist", fieldNs[len(fieldNs)-1], strings.Join(fieldNs, "."))
}
