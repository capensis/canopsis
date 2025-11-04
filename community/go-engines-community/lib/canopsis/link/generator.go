package link

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"text/template"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/js"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/sync/errgroup"
)

const workers = 10
const jsFuncName = "generate"

var (
	ErrValueIsntSlice      = errors.New("value is not slice")
	ErrValueIsntSliceOrMap = errors.New("value is not slice or map")
	ErrLabelIsMissing      = errors.New("label is missing")
	ErrURLIsMissing        = errors.New("url is missing")
)

func NewGenerator(
	client mongo.DbClient,
	tplExecutor libtemplate.Executor,
	externalDataContainer *externaldata.GetterContainer,
	logger zerolog.Logger,
) Generator {
	return &generator{
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		resolvedAlarmCollection: client.Collection(mongo.ResolvedAlarmMongoCollection),
		entityCollection:        client.Collection(mongo.EntityMongoCollection),
		linkCollection:          client.Collection(mongo.LinkRuleMongoCollection),
		tplExecutor:             tplExecutor,
		externalDataContainer:   externalDataContainer,
		logger:                  logger,
	}
}

type generator struct {
	alarmCollection         mongo.DbCollection
	resolvedAlarmCollection mongo.DbCollection
	entityCollection        mongo.DbCollection
	linkCollection          mongo.DbCollection
	tplExecutor             libtemplate.Executor
	externalDataContainer   *externaldata.GetterContainer
	logger                  zerolog.Logger

	rulesMx sync.RWMutex
	rules   []parsedRule
}

type AlarmWithData struct {
	types.Alarm  `bson:",inline"`
	Entity       types.Entity   `bson:"entity"`
	ExternalData map[string]any `bson:"-"`
}

type EntityWithData struct {
	types.Entity `bson:",inline"`
	ExternalData map[string]any `bson:"-"`
}

type entityWithAlarm struct {
	types.Entity `bson:",inline"`
	Alarm        *types.Alarm `bson:"alarm"`
}

type parsedRule struct {
	ID            string
	Type          string
	AlarmPattern  pattern.Alarm
	EntityPattern pattern.Entity
	ExternalData  []externaldata.ParsedRefParameters
	Links         []Parameters
	UrlTpls       []*template.Template
	LabelTpls     []*template.Template
	CodeExecutor  js.Executor
}

type linkWithCategory struct {
	Link
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

func (g *generator) GenerateForAlarm(ctx context.Context, alarm types.Alarm, entity types.Entity, user User) (LinksByCategory, error) {
	res, err := g.runWorkers(ctx, func(ctx context.Context, rule parsedRule) (map[string][]linkWithCategory, error) {
		return g.generateLinksByAlarms(ctx, rule, []AlarmWithData{
			{
				Alarm:  alarm,
				Entity: entity,
			},
		}, user)
	})
	if err != nil {
		return nil, err
	}

	return res[alarm.ID], nil
}

func (g *generator) GenerateForAlarms(ctx context.Context, ids []string, user User) (map[string]LinksByCategory, error) {
	alarms, err := g.getAlarms(ctx, ids)
	if err != nil || len(alarms) == 0 {
		return nil, err
	}

	return g.runWorkers(ctx, func(ctx context.Context, rule parsedRule) (map[string][]linkWithCategory, error) {
		return g.generateLinksByAlarms(ctx, rule, alarms, user)
	})
}

func (g *generator) GenerateForEntities(ctx context.Context, ids []string, user User) (map[string]LinksByCategory, error) {
	entities, err := g.getEntities(ctx, ids)
	if err != nil || len(entities) == 0 {
		return nil, err
	}

	return g.runWorkers(ctx, func(ctx context.Context, rule parsedRule) (map[string][]linkWithCategory, error) {
		return g.generateLinksByEntities(ctx, rule, entities, user)
	})
}

func (g *generator) GenerateCombinedForAlarmsByRule(ctx context.Context, ruleId string, alarmIds []string, user User) ([]Link, error) {
	rule := g.getRule(ruleId)
	if rule.ID == "" {
		return nil, ErrNoRule
	}

	alarms, err := g.getAlarms(ctx, alarmIds)
	if err != nil {
		return nil, err
	}
	if len(alarms) != len(alarmIds) {
		return nil, ErrNotMatchedAlarm
	}

	for i := range alarms {
		ok, err := match.MatchAlarmPattern(rule.AlarmPattern, &alarms[i].Alarm)
		if err != nil {
			return nil, fmt.Errorf("invalid alarm pattern linkrule=%s: %w", rule.ID, err)
		}
		if !ok {
			return nil, ErrNotMatchedAlarm
		}

		ok, err = match.MatchEntityPattern(rule.EntityPattern, &alarms[i].Entity)
		if err != nil {
			return nil, fmt.Errorf("invalid entity pattern linkrule=%s: %w", rule.ID, err)
		}
		if !ok {
			return nil, ErrNotMatchedAlarm
		}
	}

	entities := make([]EntityWithData, len(alarms))
	for i, alarm := range alarms {
		entities[i] = EntityWithData{Entity: alarm.Entity}
	}

	err = g.addExternalData(ctx, rule, alarms, entities)
	if err != nil {
		return nil, err
	}

	if rule.CodeExecutor != nil {
		args := g.getCodeArgs(rule, alarms, entities, user)
		return g.getLinksByCode(ctx, rule.CodeExecutor, args)
	}

	data := g.getTplData(rule, alarms, entities, user)
	return g.getLinksByTpl(rule.Links, rule.UrlTpls, rule.LabelTpls, data)
}

func (g *generator) runWorkers(
	ctx context.Context,
	f func(context.Context, parsedRule) (map[string][]linkWithCategory, error),
) (map[string]LinksByCategory, error) {
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

	res := make(map[string]LinksByCategory)
	for linksById := range outCh {
		for id, links := range linksById {
			if res[id] == nil {
				res[id] = make(map[string][]Link)
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
		rule := Rule{}
		err = cursor.Decode(&rule)
		if err != nil {
			return nil, err
		}

		parsedExternalData := externaldata.ParseRefParameters(rule.ExternalData, g.tplExecutor)
		pr := parsedRule{
			ID:            rule.ID,
			Type:          rule.Type,
			AlarmPattern:  rule.AlarmPattern,
			EntityPattern: rule.EntityPattern,
			ExternalData:  parsedExternalData,
		}

		if rule.SourceCode != "" {
			pr.CodeExecutor, err = js.Compile(rule.ID, rule.SourceCode)
			if err != nil {
				g.logger.Err(err).Str("rule", rule.ID).Msg("invalid source code in link rule")
				continue
			}

			parsedRules = append(parsedRules, pr)
			continue
		}

		pr.Links = rule.Links
		pr.UrlTpls = make([]*template.Template, len(rule.Links))
		pr.LabelTpls = make([]*template.Template, len(rule.Links))
		for i, link := range rule.Links {
			if link.URL == "" {
				g.logger.Error().Str("rule", rule.ID).Msg("empty url template in link rule")
				break
			}

			parsedUrl := g.tplExecutor.Parse(link.URL)
			err = parsedUrl.Err
			if err != nil {
				g.logger.Err(err).Str("rule", rule.ID).Msg("invalid url template in link rule")
				break
			}

			if link.Label == "" {
				g.logger.Error().Str("rule", rule.ID).Msg("empty label template in link rule")
				break
			}

			parsedLabel := g.tplExecutor.Parse(link.Label)
			err = parsedLabel.Err
			if err != nil {
				g.logger.Err(err).Str("rule", rule.ID).Msg("invalid label template in link rule")
				break
			}

			pr.UrlTpls[i] = parsedUrl.Tpl
			pr.LabelTpls[i] = parsedLabel.Tpl
		}
		if err != nil {
			continue
		}

		parsedRules = append(parsedRules, pr)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return parsedRules, nil
}

func (g *generator) getAlarms(ctx context.Context, ids []string) ([]AlarmWithData, error) {
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

	var openAlarms []AlarmWithData
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

	var resolvedAlarms []AlarmWithData
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

func (g *generator) generateLinksByAlarms(ctx context.Context, rule parsedRule, alarms []AlarmWithData, user User) (map[string][]linkWithCategory, error) {
	res := make(map[string][]linkWithCategory, len(alarms))
	for i := range alarms {
		ok, err := match.MatchAlarmPattern(rule.AlarmPattern, &alarms[i].Alarm)
		if err != nil {
			g.logger.Err(err).Str("rule", rule.ID).Msg("invalid alarm pattern in link rule")
			continue
		}
		if !ok {
			continue
		}

		ok, err = match.MatchEntityPattern(rule.EntityPattern, &alarms[i].Entity)
		if err != nil {
			g.logger.Err(err).Str("rule", rule.ID).Msg("invalid entity pattern in link rule")
			continue
		}
		if !ok {
			continue
		}

		argAlarms := []AlarmWithData{alarms[i]}
		argEntities := []EntityWithData{{Entity: alarms[i].Entity}}
		err = g.addExternalData(ctx, rule, argAlarms, argEntities)
		if err != nil {
			g.logger.Err(err).Str("rule", rule.ID).Msg("cannot get external data by link rule")
			continue
		}

		if rule.CodeExecutor != nil {
			args := g.getCodeArgs(rule, argAlarms, argEntities, user)
			res[alarms[i].ID], err = g.getLinksWithCategoryByCode(ctx, rule.ID, rule.CodeExecutor, args)
			if err != nil {
				g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process alarm")
			}

			continue
		}

		data := g.getTplData(rule, argAlarms, argEntities, user)
		res[alarms[i].ID], err = g.getLinksWithCategoryByTpl(rule.ID, rule.Links, rule.UrlTpls, rule.LabelTpls, data)
		if err != nil {
			g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process alarm")
		}
	}

	return res, nil
}

func (g *generator) generateLinksByEntities(ctx context.Context, rule parsedRule, entities []entityWithAlarm, user User) (map[string][]linkWithCategory, error) {
	res := make(map[string][]linkWithCategory, len(entities))
	for i := range entities {
		if entities[i].Alarm == nil && rule.Type == TypeAlarm {
			continue
		}

		if entities[i].Alarm != nil {
			ok, err := match.MatchAlarmPattern(rule.AlarmPattern, entities[i].Alarm)
			if err != nil {
				g.logger.Err(err).Str("rule", rule.ID).Msg("invalid alarm pattern in link rule")
				continue
			}
			if !ok {
				continue
			}
		}

		ok, err := match.MatchEntityPattern(rule.EntityPattern, &entities[i].Entity)
		if err != nil {
			g.logger.Err(err).Str("rule", rule.ID).Msg("invalid entity pattern in link rule")
			continue
		}
		if !ok {
			continue
		}

		var argAlarms []AlarmWithData
		if entities[i].Alarm != nil {
			argAlarms = []AlarmWithData{{
				Alarm:  *entities[i].Alarm,
				Entity: entities[i].Entity,
			}}
		}

		argEntities := []EntityWithData{{Entity: entities[i].Entity}}
		err = g.addExternalData(ctx, rule, argAlarms, argEntities)
		if err != nil {
			g.logger.Err(err).Str("rule", rule.ID).Msg("cannot get external data by link rule")
			continue
		}

		if rule.CodeExecutor != nil {
			args := g.getCodeArgs(rule, argAlarms, argEntities, user)
			res[entities[i].ID], err = g.getLinksWithCategoryByCode(ctx, rule.ID, rule.CodeExecutor, args)
			if err != nil {
				g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process entity")
			}

			continue
		}

		data := g.getTplData(rule, argAlarms, argEntities, user)
		res[entities[i].ID], err = g.getLinksWithCategoryByTpl(rule.ID, rule.Links, rule.UrlTpls, rule.LabelTpls, data)
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

func (g *generator) addExternalData(
	ctx context.Context,
	rule parsedRule,
	alarms []AlarmWithData,
	entities []EntityWithData,
) error {
	switch rule.Type {
	case TypeAlarm:
		return g.addExternalDataToAlarms(ctx, rule.ExternalData, alarms)
	case TypeEntity:
		return g.addExternalDataToEntities(ctx, rule.ExternalData, entities)
	}

	return nil
}

func (g *generator) addExternalDataToAlarms(
	ctx context.Context,
	externalData []externaldata.ParsedRefParameters,
	data []AlarmWithData,
) error {
	if len(externalData) == 0 {
		return nil
	}

	var err error
	for i, item := range data {
		data[i].ExternalData = make(map[string]any, len(externalData))
		for _, params := range externalData {
			data[i].ExternalData[params.Reference], err = g.processExternalData(ctx, params, item)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *generator) addExternalDataToEntities(
	ctx context.Context,
	externalData []externaldata.ParsedRefParameters,
	data []EntityWithData,
) error {
	if len(externalData) == 0 {
		return nil
	}

	var err error
	for i, item := range data {
		data[i].ExternalData = make(map[string]any, len(externalData))
		for _, params := range externalData {
			data[i].ExternalData[params.Reference], err = g.processExternalData(ctx, params, item)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *generator) processExternalData(
	ctx context.Context,
	params externaldata.ParsedRefParameters,
	data any,
) (any, error) {
	getter, ok := g.externalDataContainer.Get(params.Type)
	if !ok {
		return nil, fmt.Errorf("cannot find external data getter by type %q", params.Type)
	}

	return getter.Get(ctx, params, data)
}

func (g *generator) getLinksWithCategoryByTpl(
	id string,
	linkParameters []Parameters,
	urlTpls []*template.Template,
	labelTpls []*template.Template,
	data map[string]any,
) ([]linkWithCategory, error) {
	res := make([]linkWithCategory, len(linkParameters))
	for i, linkTpl := range linkParameters {
		url, err := g.tplExecutor.ExecuteByTpl(urlTpls[i], data)
		if err != nil {
			return nil, err
		}

		label, err := g.tplExecutor.ExecuteByTpl(labelTpls[i], data)
		if err != nil {
			return nil, err
		}

		res[i] = linkWithCategory{
			Category: linkTpl.Category,
			Link: Link{
				RuleID:     id,
				Label:      label,
				IconName:   linkTpl.IconName,
				Url:        url,
				Action:     linkTpl.Action,
				Single:     linkTpl.Single,
				HideInMenu: linkTpl.HideInMenu,
			},
		}
	}

	return res, nil
}

func (g *generator) getLinksByTpl(
	linkParameters []Parameters,
	urlTpls []*template.Template,
	labelTpls []*template.Template,
	data map[string]any,
) ([]Link, error) {
	res := make([]Link, len(linkParameters))
	for i, linkTpl := range linkParameters {
		url, err := g.tplExecutor.ExecuteByTpl(urlTpls[i], data)
		if err != nil {
			return nil, err
		}

		label, err := g.tplExecutor.ExecuteByTpl(labelTpls[i], data)
		if err != nil {
			return nil, err
		}

		res[i] = Link{
			Label:      label,
			IconName:   linkTpl.IconName,
			Url:        url,
			Action:     linkTpl.Action,
			Single:     linkTpl.Single,
			HideInMenu: linkTpl.HideInMenu,
		}
	}

	return res, nil
}

func (g *generator) getLinksWithCategoryByCode(
	ctx context.Context,
	id string,
	codeExecutor js.Executor,
	args []any,
) ([]linkWithCategory, error) {
	r, err := codeExecutor.ExecuteFunc(ctx, jsFuncName, args...)
	if err != nil {
		return nil, err
	}

	s, ok := r.([]any)
	if !ok {
		return nil, ErrValueIsntSlice
	}

	res := make([]linkWithCategory, len(s))
	for i := 0; i < len(s); i++ {
		item, ok := s[i].(map[string]any)
		if !ok {
			return nil, ErrValueIsntSliceOrMap
		}

		category, _ := item["category"].(string)
		label, _ := item["label"].(string)
		iconName, _ := item["icon_name"].(string)
		url, _ := item["url"].(string)
		action, _ := item["action"].(string)
		if url == "" {
			return nil, ErrURLIsMissing
		}
		if label == "" {
			return nil, ErrLabelIsMissing
		}

		res[i] = linkWithCategory{
			Category: category,
			Link: Link{
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
		if hideInMenu, ok := item["hide_in_menu"].(bool); ok {
			res[i].Link.HideInMenu = hideInMenu
		}
	}

	return res, nil
}

func (g *generator) getLinksByCode(
	ctx context.Context,
	codeExecutor js.Executor,
	args []any,
) ([]Link, error) {
	r, err := codeExecutor.ExecuteFunc(ctx, jsFuncName, args...)
	if err != nil {
		return nil, err
	}

	s, ok := r.([]any)
	if !ok {
		return nil, ErrValueIsntSlice
	}

	res := make([]Link, len(s))
	for i := 0; i < len(s); i++ {
		item, ok := s[i].(map[string]any)
		if !ok {
			return nil, ErrValueIsntSliceOrMap
		}

		label, _ := item["label"].(string)
		iconName, _ := item["icon_name"].(string)
		url, _ := item["url"].(string)
		action, _ := item["action"].(string)
		if url == "" {
			return nil, ErrURLIsMissing
		}
		if label == "" {
			return nil, ErrLabelIsMissing
		}

		res[i] = Link{
			Label:    label,
			IconName: iconName,
			Url:      url,
			Action:   action,
		}
		if single, ok := item["single"].(bool); ok {
			res[i].Single = single
		}
		if hideInMenu, ok := item["hide_in_menu"].(bool); ok {
			res[i].HideInMenu = hideInMenu
		}
	}

	return res, nil
}

func (g *generator) getTplData(
	rule parsedRule,
	alarms []AlarmWithData,
	entities []EntityWithData,
	user User,
) map[string]any {
	var data map[string]any
	switch rule.Type {
	case TypeAlarm:
		data = map[string]any{
			"Alarms": alarms,
			"User":   user,
		}
	case TypeEntity:
		data = map[string]any{
			"Entities": entities,
			"User":     user,
		}
	}

	return data
}

func (g *generator) getCodeArgs(
	rule parsedRule,
	alarms []AlarmWithData,
	entities []EntityWithData,
	user User,
) []any {
	var items any
	switch rule.Type {
	case TypeAlarm:
		items = alarms
	case TypeEntity:
		items = entities
	}

	args := []any{items, user}
	tplVars := g.tplExecutor.GetDefaultTplVars()
	if len(tplVars) > 0 {
		keys := make([]string, 0, len(tplVars))
		for k := range tplVars {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		for _, k := range keys {
			args = append(args, tplVars[k])
		}
	}

	return args
}
