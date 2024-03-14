package link

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"
	"text/template"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/js"
	liblink "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
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
		dbClient:                client,
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		resolvedAlarmCollection: client.Collection(mongo.ResolvedAlarmMongoCollection),
		entityCollection:        client.Collection(mongo.EntityMongoCollection),
		linkCollection:          client.Collection(mongo.LinkRuleMongoCollection),
		tplExecutor:             tplExecutor,
		logger:                  logger,
	}
}

type generator struct {
	dbClient                mongo.DbClient
	alarmCollection         mongo.DbCollection
	resolvedAlarmCollection mongo.DbCollection
	entityCollection        mongo.DbCollection
	linkCollection          mongo.DbCollection
	tplExecutor             libtemplate.Executor
	logger                  zerolog.Logger

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

type entityWithAlarm struct {
	types.Entity `bson:",inline"`
	Alarm        *types.Alarm `bson:"alarm"`
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

type linkWithCategory struct {
	liblink.Link
	Category string
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

func (g *generator) GenerateForAlarm(ctx context.Context, alarm types.Alarm, entity types.Entity) (liblink.LinksByCategory, error) {
	res, err := g.runWorkers(ctx, func(ctx context.Context, rule parsedRule) (map[string][]linkWithCategory, error) {
		return g.generateLinksByAlarms(ctx, rule, []alarmWithData{
			{
				Alarm:  alarm,
				Entity: entity,
			},
		})
	})
	if err != nil {
		return nil, err
	}

	return res[alarm.ID], nil
}

func (g *generator) GenerateForAlarms(ctx context.Context, ids []string) (map[string]liblink.LinksByCategory, error) {
	alarms, err := g.getAlarms(ctx, ids)
	if err != nil || len(alarms) == 0 {
		return nil, err
	}

	return g.runWorkers(ctx, func(ctx context.Context, rule parsedRule) (map[string][]linkWithCategory, error) {
		return g.generateLinksByAlarms(ctx, rule, alarms)
	})
}

func (g *generator) GenerateForEntities(ctx context.Context, ids []string) (map[string]liblink.LinksByCategory, error) {
	entities, err := g.getEntities(ctx, ids)
	if err != nil || len(entities) == 0 {
		return nil, err
	}

	return g.runWorkers(ctx, func(ctx context.Context, rule parsedRule) (map[string][]linkWithCategory, error) {
		return g.generateLinksByEntities(ctx, rule, entities)
	})
}

func (g *generator) GenerateCombinedForAlarmsByRule(ctx context.Context, ruleId string, alarmIds []string) ([]liblink.Link, error) {
	rule := g.getRule(ruleId)
	if rule.ID == "" {
		return nil, liblink.ErrNoRule
	}

	alarms, err := g.getAlarms(ctx, alarmIds)
	if err != nil {
		return nil, err
	}
	if len(alarms) != len(alarmIds) {
		return nil, liblink.ErrNotMatchedAlarm
	}

	for _, alarm := range alarms {
		ok, err := rule.AlarmPattern.Match(alarm.Alarm)
		if err != nil {
			return nil, fmt.Errorf("invalid alarm pattern linkrule=%s: %w", rule.ID, err)
		}
		if !ok {
			return nil, liblink.ErrNotMatchedAlarm
		}

		ok, err = rule.EntityPattern.Match(alarm.Entity)
		if err != nil {
			return nil, fmt.Errorf("invalid entity pattern linkrule=%s: %w", rule.ID, err)
		}
		if !ok {
			return nil, liblink.ErrNotMatchedAlarm
		}
	}

	entities := make([]entityWithData, len(alarms))
	for i, alarm := range alarms {
		entities[i] = entityWithData{Entity: alarm.Entity}
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

		return g.getLinksByCode(ctx, rule.CodeExecutor, arg)
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

	return g.getLinksByTpl(rule.Links, rule.LinkTpls, data)
}

func (g *generator) runWorkers(
	ctx context.Context,
	f func(context.Context, parsedRule) (map[string][]linkWithCategory, error),
) (map[string]liblink.LinksByCategory, error) {
	eg, ctx := errgroup.WithContext(ctx)
	inCh := make(chan parsedRule)
	outCh := make(chan map[string][]linkWithCategory)

	go func() {
		defer close(inCh)

		g.rulesMx.RLock()
		defer g.rulesMx.RUnlock()
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
	for linksById := range outCh {
		for id, links := range linksById {
			if res[id] == nil {
				res[id] = make(map[string][]liblink.Link)
			}
			for _, link := range links {
				res[id][link.Category] = append(res[id][link.Category], link.Link)
			}
		}
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
			pr.CodeExecutor, err = js.Compile(rule.ID, rule.SourceCode)
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
	pipeline := []bson.M{
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
	}
	openPipeline := []bson.M{
		{"$match": bson.M{
			"_id":        bson.M{"$in": ids},
			"v.resolved": nil,
		}},
	}
	openPipeline = append(openPipeline, pipeline...)
	openCursor, err := g.alarmCollection.Aggregate(ctx, openPipeline)
	if err != nil {
		return nil, err
	}

	var openAlarms []alarmWithData
	err = openCursor.All(ctx, &openAlarms)
	if err != nil {
		return nil, err
	}

	if len(openAlarms) == len(ids) {
		return openAlarms, nil
	}

	resolvedPipeline := []bson.M{
		{"$match": bson.M{
			"_id": bson.M{"$in": ids},
		}},
	}
	resolvedPipeline = append(resolvedPipeline, pipeline...)
	resolvedCursor, err := g.resolvedAlarmCollection.Aggregate(ctx, resolvedPipeline)
	if err != nil {
		return nil, err
	}

	var resolvedAlarms []alarmWithData
	err = resolvedCursor.All(ctx, &resolvedAlarms)
	if err != nil {
		return nil, err
	}

	return append(openAlarms, resolvedAlarms...), nil
}

func (g *generator) getEntities(ctx context.Context, ids []string) ([]entityWithAlarm, error) {
	cursor, err := g.entityCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": ids}}},
		{"$lookup": bson.M{
			"from":         mongo.AlarmMongoCollection,
			"localField":   "_id",
			"foreignField": "d",
			"pipeline": []bson.M{
				{"$match": bson.M{"v.resolved": nil}},
			},
			"as": "alarm",
		}},
		{"$unwind": bson.M{"path": "$alarm", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"_id": 1}},
	})
	if err != nil {
		return nil, err
	}

	var entities []entityWithAlarm
	err = cursor.All(ctx, &entities)
	return entities, err
}

func (g *generator) generateLinksByAlarms(ctx context.Context, rule parsedRule, alarms []alarmWithData) (map[string][]linkWithCategory, error) {
	res := make(map[string][]linkWithCategory, len(alarms))
	for _, alarm := range alarms {
		ok, err := rule.AlarmPattern.Match(alarm.Alarm)
		if err != nil {
			return nil, fmt.Errorf("invalid alarm pattern linkrule=%s: %w", rule.ID, err)
		}
		if !ok {
			continue
		}

		ok, err = rule.EntityPattern.Match(alarm.Entity)
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
			res[alarm.ID], err = g.getLinksWithCategoryByCode(ctx, rule.ID, rule.CodeExecutor, arg)
			if err != nil {
				g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process alarm")
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

		res[alarm.ID], err = g.getLinksWithCategoryByTpl(rule.ID, rule.Links, rule.LinkTpls, data)
		if err != nil {
			g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process alarm")
		}
	}

	return res, nil
}

func (g *generator) generateLinksByEntities(ctx context.Context, rule parsedRule, entities []entityWithAlarm) (map[string][]linkWithCategory, error) {
	res := make(map[string][]linkWithCategory, len(entities))
	for _, entity := range entities {
		if entity.Alarm == nil && rule.Type == liblink.TypeAlarm {
			continue
		}

		if entity.Alarm != nil {
			ok, err := rule.AlarmPattern.Match(*entity.Alarm)
			if err != nil {
				return nil, fmt.Errorf("invalid alarm pattern linkrule=%s: %w", rule.ID, err)
			}
			if !ok {
				continue
			}
		}

		ok, err := rule.EntityPattern.Match(entity.Entity)
		if err != nil {
			return nil, fmt.Errorf("invalid entity pattern linkrule=%s: %w", rule.ID, err)
		}
		if !ok {
			continue
		}

		var arg any
		switch rule.Type {
		case liblink.TypeAlarm:
			v := []alarmWithData{{Entity: entity.Entity, Alarm: *entity.Alarm}}
			err := g.addExternalDataToAlarms(ctx, rule.ExternalData, rule.ExternalDataTpl, v)
			if err != nil {
				return nil, err
			}
			arg = v
		case liblink.TypeEntity:
			v := []entityWithData{{Entity: entity.Entity}}
			err := g.addExternalDataToEntities(ctx, rule.ExternalData, rule.ExternalDataTpl, v)
			if err != nil {
				return nil, err
			}
			arg = v
		}

		if rule.CodeExecutor != nil {
			res[entity.ID], err = g.getLinksWithCategoryByCode(ctx, rule.ID, rule.CodeExecutor, arg)
			if err != nil {
				g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process alarm")
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

		res[entity.ID], err = g.getLinksWithCategoryByTpl(rule.ID, rule.Links, rule.LinkTpls, data)
		if err != nil {
			g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process entity")
		}
	}

	return res, nil
}

func (g *generator) getRule(id string) parsedRule {
	g.rulesMx.RLock()
	defer g.rulesMx.RUnlock()
	for _, rule := range g.rules {
		if rule.ID == id {
			return rule
		}
	}

	return parsedRule{}
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

func (g *generator) getLinksWithCategoryByTpl(
	id string,
	linkTpls []liblink.Parameters,
	tpls []*template.Template,
	data map[string]any,
) ([]linkWithCategory, error) {
	res := make([]linkWithCategory, len(linkTpls))
	for i, linkTpl := range linkTpls {
		url, err := g.tplExecutor.ExecuteByTpl(tpls[i], data)
		if err != nil {
			return nil, err
		}

		res[i] = linkWithCategory{
			Category: linkTpl.Category,
			Link: liblink.Link{
				RuleID:   id,
				Label:    linkTpl.Label,
				IconName: linkTpl.IconName,
				Url:      url,
				Single:   linkTpl.Single,
				Action:   linkTpl.Action,
			},
		}
	}

	return res, nil
}

func (g *generator) getLinksByTpl(
	linkTpls []liblink.Parameters,
	tpls []*template.Template,
	data map[string]any,
) ([]liblink.Link, error) {
	res := make([]liblink.Link, len(linkTpls))
	for i, linkTpl := range linkTpls {
		url, err := g.tplExecutor.ExecuteByTpl(tpls[i], data)
		if err != nil {
			return nil, err
		}

		res[i] = liblink.Link{
			Label:    linkTpl.Label,
			IconName: linkTpl.IconName,
			Url:      url,
			Single:   linkTpl.Single,
			Action:   linkTpl.Action,
		}
	}

	return res, nil
}

func (g *generator) getLinksWithCategoryByCode(
	ctx context.Context,
	id string,
	codeExecutor js.Executor,
	arg any,
) ([]linkWithCategory, error) {
	r, err := codeExecutor.ExecuteFunc(ctx, jsFuncName, arg)
	if err != nil {
		return nil, err
	}

	s, ok := r.([]any)
	if !ok {
		return nil, fmt.Errorf("value is not slice")
	}

	res := make([]linkWithCategory, len(s))
	for i := 0; i < len(s); i++ {
		item, ok := s[i].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("value is not slice of map")
		}

		category, _ := item["category"].(string)
		label, _ := item["label"].(string)
		iconName, _ := item["icon_name"].(string)
		url, _ := item["url"].(string)
		action, _ := item["action"].(string)
		if url == "" {
			return nil, fmt.Errorf("url is missing")
		}
		if label == "" {
			return nil, fmt.Errorf("label is missing")
		}

		res[i] = linkWithCategory{
			Category: category,
			Link: liblink.Link{
				RuleID:   id,
				Label:    label,
				IconName: iconName,
				Url:      url,
				Action:   action,
			},
		}
		if single, ok := item["single"].(bool); ok {
			res[i].Link.Single = single
		}
	}

	return res, nil
}

func (g *generator) getLinksByCode(
	ctx context.Context,
	codeExecutor js.Executor,
	arg any,
) ([]liblink.Link, error) {
	r, err := codeExecutor.ExecuteFunc(ctx, jsFuncName, arg)
	if err != nil {
		return nil, err
	}

	s, ok := r.([]any)
	if !ok {
		return nil, fmt.Errorf("value is not slice")
	}

	res := make([]liblink.Link, len(s))
	for i := 0; i < len(s); i++ {
		item, ok := s[i].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("value is not slice of map")
		}

		label, _ := item["label"].(string)
		iconName, _ := item["icon_name"].(string)
		url, _ := item["url"].(string)
		action, _ := item["action"].(string)
		if url == "" {
			return nil, fmt.Errorf("url is missing")
		}
		if label == "" {
			return nil, fmt.Errorf("label is missing")
		}

		res[i] = liblink.Link{
			Label:    label,
			IconName: iconName,
			Url:      url,
			Action:   action,
		}
		if single, ok := item["single"].(bool); ok {
			res[i].Single = single
		}
	}

	return res, nil
}
