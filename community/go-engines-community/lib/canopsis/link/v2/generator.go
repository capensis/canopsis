package link

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sync"
	"text/template"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/js"
	liblink "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libreflect "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/reflect"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

const workers = 10
const jsFuncName = "generate"

func NewGenerator(
	client mongo.DbClient,
	tplExecutor libtemplate.Executor,
	logger zerolog.Logger,
) liblink.Generator {
	return &generator{
		dbClient:         client,
		alarmCollection:  client.Collection(mongo.AlarmMongoCollection),
		entityCollection: client.Collection(mongo.EntityMongoCollection),
		linkCollection:   client.Collection(mongo.LinkRuleMongoCollection),
		tplExecutor:      tplExecutor,
		logger:           logger,
	}
}

type generator struct {
	dbClient         mongo.DbClient
	alarmCollection  mongo.DbCollection
	entityCollection mongo.DbCollection
	linkCollection   mongo.DbCollection
	tplExecutor      libtemplate.Executor
	logger           zerolog.Logger

	rulesMx sync.RWMutex
	rules   []parsedRule
}

type alarmWithData struct {
	types.Alarm  `bson:",inline"`
	Entity       types.Entity              `bson:"entity"`
	ExternalData map[string]map[string]any `bson:"-"`
}

type entityWithData struct {
	types.Entity `bson:",inline"`
	ExternalData map[string]map[string]any `bson:"-"`
}

type parsedRule struct {
	ID              string
	Type            string
	AlarmPattern    pattern.Alarm
	EntityPattern   pattern.Entity
	ExternalData    map[string]liblink.ExternalDataParameters
	ExternalDataTpl map[string]map[string]map[string]*template.Template
	Links           []liblink.Parameters
	LinkTpls        []*template.Template
	CodeExecutor    js.Executor
}

func (g *generator) Load(ctx context.Context) error {
	rules, err := g.getRules(ctx)
	if err != nil {
		return err
	}

	g.rulesMx.Lock()
	defer g.rulesMx.Unlock()
	g.rules = rules
	return nil
}

func (g *generator) GenerateForAlarms(ctx context.Context, ids []string) (map[string]liblink.LinksByCategory, error) {
	alarms, err := g.getAlarms(ctx, ids)
	if err != nil || len(alarms) == 0 {
		return nil, err
	}

	return g.runWorkers(ctx, func(ctx context.Context, rule parsedRule) (map[string]liblink.LinksByCategory, error) {
		return g.generateLinksByAlarms(ctx, rule, alarms)
	})
}

func (g *generator) GenerateForEntities(ctx context.Context, ids []string) (map[string]liblink.LinksByCategory, error) {
	entities, err := g.getEntities(ctx, ids)
	if err != nil || len(entities) == 0 {
		return nil, err
	}

	return g.runWorkers(ctx, func(ctx context.Context, rule parsedRule) (map[string]liblink.LinksByCategory, error) {
		return g.generateLinksByEntities(ctx, rule, entities)
	})
}

func (g *generator) GenerateForAllAlarms(ctx context.Context, ids []string) ([]liblink.Link, error) {
	alarms, err := g.getAlarms(ctx, ids)
	if err != nil || len(alarms) != len(ids) {
		return nil, err
	}
	entities := make([]entityWithData, len(alarms))
	for i, alarm := range alarms {
		entities[i] = entityWithData{Entity: alarm.Entity}
	}

	linksMap, err := g.runWorkers(ctx, func(ctx context.Context, rule parsedRule) (map[string]liblink.LinksByCategory, error) {
		return g.generateLinksByAllAlarms(ctx, rule, alarms, entities)
	})
	if err != nil {
		return nil, err
	}

	res := make([]liblink.Link, 0)
	for _, linksByCategory := range linksMap {
		for _, links := range linksByCategory {
			res = append(res, links...)
		}
	}

	return res, nil
}

func (g *generator) runWorkers(
	ctx context.Context,
	f func(ctx context.Context, rule parsedRule) (map[string]liblink.LinksByCategory, error),
) (map[string]liblink.LinksByCategory, error) {
	g.rulesMx.RLock()
	defer g.rulesMx.RUnlock()

	if len(g.rules) == 0 {
		return nil, nil
	}

	eg, ctx := errgroup.WithContext(ctx)
	inCh := make(chan parsedRule)
	outCh := make(chan map[string]liblink.LinksByCategory)

	go func() {
		defer close(inCh)
		for _, rule := range g.rules {
			select {
			case <-ctx.Done():
				return
			case inCh <- rule:
			}
		}
	}()

	for i := 0; i < workers; i++ {
		eg.Go(func() error {
			for rule := range inCh {
				res, err := f(ctx, rule)
				if err != nil {
					g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process link rule")
					continue
				}

				outCh <- res
			}

			return nil
		})
	}

	go func() {
		_ = eg.Wait()
		close(outCh)
	}()

	res := make(map[string]liblink.LinksByCategory)
	for v := range outCh {
		liblink.MergeLinks(res, v)
	}

	err := eg.Wait()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (g *generator) getRules(ctx context.Context) ([]parsedRule, error) {
	cursor, err := g.linkCollection.Find(ctx, bson.M{"enabled": true})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	parsedRules := make([]parsedRule, 0)
	for cursor.Next(ctx) {
		rule := liblink.Rule{}
		err = cursor.Decode(&rule)
		if err != nil {
			return nil, err
		}

		externalDataTpl := make(map[string]map[string]map[string]*template.Template)
		for ref, params := range rule.ExternalData {
			externalDataTpl[ref] = map[string]map[string]*template.Template{
				"select": make(map[string]*template.Template, len(params.Select)),
				"regexp": make(map[string]*template.Template, len(params.Regexp)),
			}
			for k, v := range params.Select {
				externalDataTpl[ref]["select"][k], err = g.tplExecutor.Parse(v)
				if err != nil {
					return nil, fmt.Errorf("invalid template linkrule=%s: %w", rule.ID, err)
				}
			}
			for k, v := range params.Regexp {
				externalDataTpl[ref]["regexp"][k], err = g.tplExecutor.Parse(v)
				if err != nil {
					return nil, fmt.Errorf("invalid template linkrule=%s: %w", rule.ID, err)
				}
			}
		}

		pr := parsedRule{
			ID:              rule.ID,
			Type:            rule.Type,
			AlarmPattern:    rule.AlarmPattern,
			EntityPattern:   rule.EntityPattern,
			ExternalData:    rule.ExternalData,
			ExternalDataTpl: externalDataTpl,
		}

		if rule.SourceCode != "" {
			pr.CodeExecutor, err = js.Compile(rule.ID, rule.SourceCode, jsFuncName)
			if err != nil {
				return nil, fmt.Errorf("invalid source code linkrule=%s: %w", rule.ID, err)
			}

			parsedRules = append(parsedRules, pr)
			continue
		}

		pr.Links = rule.Links
		pr.LinkTpls = make([]*template.Template, len(rule.Links))
		for i, link := range rule.Links {
			pr.LinkTpls[i], err = g.tplExecutor.Parse(link.Url)
			if err != nil {
				return nil, fmt.Errorf("invalid template linkrule=%s: %w", rule.ID, err)
			}
		}

		parsedRules = append(parsedRules, pr)
	}

	return parsedRules, nil
}

func (g *generator) getAlarms(ctx context.Context, ids []string) ([]alarmWithData, error) {
	cursor, err := g.alarmCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"_id":        bson.M{"$in": ids},
			"v.resolved": nil,
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$project": bson.M{
			"v.steps": 0,
		}},
		{"$sort": bson.M{"_id": 1}},
	})
	if err != nil {
		return nil, err
	}

	var alarms []alarmWithData
	err = cursor.All(ctx, &alarms)
	return alarms, err
}

func (g *generator) getEntities(ctx context.Context, ids []string) ([]entityWithData, error) {
	cursor, err := g.entityCollection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}}, options.Find().SetSort(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}

	var entities []entityWithData
	err = cursor.All(ctx, &entities)
	return entities, err
}

func (g *generator) generateLinksByAlarms(ctx context.Context, rule parsedRule, alarms []alarmWithData) (map[string]liblink.LinksByCategory, error) {
	res := make(map[string]liblink.LinksByCategory, len(alarms))
	for _, alarm := range alarms {
		ok, err := rule.AlarmPattern.Match(alarm.Alarm)
		if err != nil {
			return nil, fmt.Errorf("invalid alarm pattern linkrule=%s: %w", rule.ID, err)
		}
		if !ok {
			continue
		}

		ok, _, err = rule.EntityPattern.Match(alarm.Entity)
		if err != nil {
			return nil, fmt.Errorf("invalid entity pattern linkrule=%s: %w", rule.ID, err)
		}
		if !ok {
			continue
		}

		var arg any
		switch rule.Type {
		case liblink.TypeAlarm:
			v := []alarmWithData{alarm}
			err := g.addExternalDataToAlarms(ctx, rule.ExternalData, rule.ExternalDataTpl, v)
			if err != nil {
				return nil, err
			}
			arg = v
		case liblink.TypeEntity:
			v := []entityWithData{{Entity: alarm.Entity}}
			err := g.addExternalDataToEntities(ctx, rule.ExternalData, rule.ExternalDataTpl, v)
			if err != nil {
				return nil, err
			}
			arg = v
		}

		if rule.CodeExecutor != nil {
			res[alarm.ID], err = g.processCode(ctx, rule.CodeExecutor, arg)
			if err != nil {
				g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process alarm")
				continue
			}

			continue
		}

		var data map[string]any
		switch rule.Type {
		case liblink.TypeAlarm:
			data = map[string]any{
				"Alarms": arg,
			}
		case liblink.TypeEntity:
			data = map[string]any{
				"Entities": arg,
			}
		}

		res[alarm.ID], err = g.processLinks(rule.Links, rule.LinkTpls, data)
		if err != nil {
			g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process alarm")
			continue
		}
	}

	return res, nil
}

func (g *generator) generateLinksByEntities(ctx context.Context, rule parsedRule, entities []entityWithData) (map[string]liblink.LinksByCategory, error) {
	if rule.Type != liblink.TypeEntity {
		return nil, nil
	}

	res := make(map[string]liblink.LinksByCategory, len(entities))
	for _, entity := range entities {
		ok, _, err := rule.EntityPattern.Match(entity.Entity)
		if err != nil {
			return nil, fmt.Errorf("invalid entity pattern linkrule=%s: %w", rule.ID, err)
		}
		if !ok {
			continue
		}

		arg := []entityWithData{entity}
		err = g.addExternalDataToEntities(ctx, rule.ExternalData, rule.ExternalDataTpl, arg)
		if err != nil {
			return nil, err
		}

		if rule.CodeExecutor != nil {
			res[entity.ID], err = g.processCode(ctx, rule.CodeExecutor, arg)
			if err != nil {
				g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process entity")
				continue
			}

			continue
		}

		res[entity.ID], err = g.processLinks(rule.Links, rule.LinkTpls, map[string]any{
			"Entities": arg,
		})
		if err != nil {
			g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process entity")
			continue
		}
	}

	return res, nil
}

func (g *generator) generateLinksByAllAlarms(
	ctx context.Context,
	rule parsedRule,
	alarms []alarmWithData,
	entities []entityWithData,
) (map[string]liblink.LinksByCategory, error) {
	for _, alarm := range alarms {
		ok, err := rule.AlarmPattern.Match(alarm.Alarm)
		if err != nil {
			return nil, fmt.Errorf("invalid alarm pattern linkrule=%s: %w", rule.ID, err)
		}
		if !ok {
			return nil, nil
		}

		ok, _, err = rule.EntityPattern.Match(alarm.Entity)
		if err != nil {
			return nil, fmt.Errorf("invalid entity pattern linkrule=%s: %w", rule.ID, err)
		}
		if !ok {
			return nil, nil
		}
	}

	switch rule.Type {
	case liblink.TypeAlarm:
		err := g.addExternalDataToAlarms(ctx, rule.ExternalData, rule.ExternalDataTpl, alarms)
		if err != nil {
			return nil, err
		}
	case liblink.TypeEntity:
		err := g.addExternalDataToEntities(ctx, rule.ExternalData, rule.ExternalDataTpl, entities)
		if err != nil {
			return nil, err
		}
	}

	if rule.CodeExecutor != nil {
		var arg any
		switch rule.Type {
		case liblink.TypeAlarm:
			arg = alarms
		case liblink.TypeEntity:
			arg = entities
		}

		linksByCategory, err := g.processCode(ctx, rule.CodeExecutor, arg)
		if err != nil {
			return nil, err
		}

		return map[string]liblink.LinksByCategory{
			"": linksByCategory,
		}, nil
	}

	var data map[string]any
	switch rule.Type {
	case liblink.TypeAlarm:
		data = map[string]any{
			"Alarms": alarms,
		}
	case liblink.TypeEntity:
		data = map[string]any{
			"Entities": entities,
		}
	}

	linksByCategory, err := g.processLinks(rule.Links, rule.LinkTpls, data)
	if err != nil {
		return nil, err
	}

	return map[string]liblink.LinksByCategory{
		"": linksByCategory,
	}, nil
}

func (g *generator) addExternalDataToAlarms(
	ctx context.Context,
	externalData map[string]liblink.ExternalDataParameters,
	externalDataTpl map[string]map[string]map[string]*template.Template,
	data []alarmWithData,
) error {
	if len(externalData) == 0 {
		return nil
	}

	var err error
	for i, item := range data {
		data[i].ExternalData = make(map[string]map[string]any, len(externalData))
		for ref, params := range externalData {
			data[i].ExternalData[ref], err = g.processExternalData(ctx, params, externalDataTpl[ref], item)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *generator) addExternalDataToEntities(
	ctx context.Context,
	externalData map[string]liblink.ExternalDataParameters,
	externalDataTpl map[string]map[string]map[string]*template.Template,
	data []entityWithData,
) error {
	if len(externalData) == 0 {
		return nil
	}

	var err error
	for i, item := range data {
		data[i].ExternalData = make(map[string]map[string]any, len(externalData))
		for ref, params := range externalData {
			data[i].ExternalData[ref], err = g.processExternalData(ctx, params, externalDataTpl[ref], item)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *generator) processExternalData(
	ctx context.Context,
	params liblink.ExternalDataParameters,
	tpls map[string]map[string]*template.Template,
	data any,
) (map[string]any, error) {
	collection := g.dbClient.Collection(params.Collection)
	sort := mongo.GetSort(params.SortBy, params.Sort)
	query := bson.M{}
	var err error
	for k := range params.Select {
		tpl := tpls["select"][k]
		query[k], err = g.tplExecutor.ExecuteByTpl(tpl, data)
		if err != nil {
			return nil, fmt.Errorf("cannot execute select template %q: %w", k, err)
		}
	}

	if len(params.Regexp) == 0 {
		res := make(map[string]any)
		err = collection.
			FindOne(ctx, query, options.FindOne().SetSort(sort)).
			Decode(&res)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil, nil
			}

			return nil, err
		}

		return res, nil
	}

	regexpMap := make(map[string]string, len(params.Regexp))
	for k := range params.Regexp {
		tpl := tpls["regexp"][k]
		regexpMap[k], err = g.tplExecutor.ExecuteByTpl(tpl, data)
		if err != nil {
			return nil, fmt.Errorf("cannot execute regexp template %q: %w", k, err)
		}
	}

	cursor, err := collection.Find(ctx, query, options.Find().SetSort(sort))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var row map[string]any
		err := cursor.Decode(&row)
		if err != nil {
			return nil, fmt.Errorf("cannot decode data: %w", err)
		}

		matched := false
		for field, v := range regexpMap {
			regexpStr, ok := row[field].(string)
			if !ok {
				matched = false
				break
			}

			re, err := regexp.Compile(regexpStr)
			if err != nil {
				return nil, fmt.Errorf("cannot compile %q regexp %q: %w", field, regexpStr, err)
			}

			matched = re.Match([]byte(v))
			if !matched {
				break
			}
		}

		if matched {
			return row, nil
		}
	}

	return nil, nil
}

func (g *generator) processLinks(
	links []liblink.Parameters,
	tpls []*template.Template,
	data map[string]any,
) (liblink.LinksByCategory, error) {
	res := make(liblink.LinksByCategory)
	for i, link := range links {
		url, err := g.tplExecutor.ExecuteByTpl(tpls[i], data)
		if err != nil {
			return nil, err
		}

		res[link.Category] = append(res[link.Category], liblink.Link{
			Label:    link.Label,
			IconName: link.IconName,
			Url:      url,
		})
	}

	return res, nil
}

func (g *generator) processCode(
	ctx context.Context,
	codeExecutor js.Executor,
	arg any,
) (liblink.LinksByCategory, error) {
	r, err := codeExecutor.Execute(ctx, arg)
	if err != nil {
		return nil, err
	}

	v := libreflect.UnwrapPointer(reflect.ValueOf(r))
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("generate returns not slice")
	}

	res := make(liblink.LinksByCategory)
	for i := 0; i < v.Len(); i++ {
		item := libreflect.UnwrapPointer(v.Index(i))
		if item.Kind() != reflect.Map {
			return nil, fmt.Errorf("generate returns not slice of map")
		}

		category := g.getMapStringItem(item, "category")
		label := g.getMapStringItem(item, "label")
		iconName := g.getMapStringItem(item, "icon_name")
		url := g.getMapStringItem(item, "url")
		if url == "" {
			return nil, fmt.Errorf("generate returns no url")
		}
		if label == "" {
			return nil, fmt.Errorf("generate returns no label")
		}

		res[category] = append(res[category], liblink.Link{
			Label:    label,
			IconName: iconName,
			Url:      url,
		})
	}

	return res, nil
}

func (g *generator) getMapStringItem(v reflect.Value, k string) string {
	val := libreflect.UnwrapPointer(v.MapIndex(reflect.ValueOf(k)))
	if val.Kind() == reflect.String {
		return val.String()
	}

	return ""
}
