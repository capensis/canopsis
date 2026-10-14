package link

import (
	"context"
	"fmt"
	"sort"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/sync/errgroup"
)

const workers = 10

func NewGenerator(
	client mongo.DbClient,
	tplExecutor libtemplate.Executor,
	externalDataGetter externaldata.Getter,
	logger zerolog.Logger,
) Generator {
	return &generator{
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		resolvedAlarmCollection: client.Collection(mongo.ResolvedAlarmMongoCollection),
		entityCollection:        client.Collection(mongo.EntityMongoCollection),
		linkCollection:          client.Collection(mongo.LinkRuleMongoCollection),
		tplExecutor:             tplExecutor,
		externalDataGetter:      externalDataGetter,
		logger:                  logger,
		ruleProvider:            NewRuleProvider(client, tplExecutor, logger),
	}
}

type generator struct {
	alarmCollection         mongo.DbCollection
	resolvedAlarmCollection mongo.DbCollection
	entityCollection        mongo.DbCollection
	linkCollection          mongo.DbCollection
	tplExecutor             libtemplate.Executor
	externalDataGetter      externaldata.Getter
	logger                  zerolog.Logger
	ruleProvider            RuleProvider
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

func (g *generator) Load(ctx context.Context) error {
	return g.ruleProvider.Load(ctx)
}

func (g *generator) GenerateForAlarm(ctx context.Context, alarm types.Alarm, entity types.Entity, user User) (LinksByCategory, error) {
	res, err := g.runWorkers(ctx, func(ctx context.Context, rule ParsedRule) (map[string][]GeneratedLink, error) {
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

	return g.runWorkers(ctx, func(ctx context.Context, rule ParsedRule) (map[string][]GeneratedLink, error) {
		return g.generateLinksByAlarms(ctx, rule, alarms, user)
	})
}

func (g *generator) GenerateForEntities(ctx context.Context, ids []string, user User) (map[string]LinksByCategory, error) {
	entities, err := g.getEntities(ctx, ids)
	if err != nil || len(entities) == 0 {
		return nil, err
	}

	return g.runWorkers(ctx, func(ctx context.Context, rule ParsedRule) (map[string][]GeneratedLink, error) {
		return g.generateLinksByEntities(ctx, rule, entities, user)
	})
}

func (g *generator) GenerateCombinedForAlarmsByRule(ctx context.Context, ruleId string, alarmIds []string, user User) ([]Link, error) {
	rule, err := g.ruleProvider.Rule(ruleId)
	if err != nil {
		return nil, err
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
		links, err := ExecCode(ctx, rule, args)
		if err != nil {
			return nil, err
		}

		res := make([]Link, len(links))
		for i, link := range links {
			res[i] = link.Link
			res[i].RuleID = ""
		}

		return res, nil
	}

	data := g.getTplData(rule, alarms, entities, user)

	links, err := ExecTemplate(rule, data, g.tplExecutor)
	if err != nil {
		return nil, err
	}

	res := make([]Link, len(links))
	for i, link := range links {
		res[i] = link.Link
		res[i].RuleID = ""
	}

	return res, nil
}

func (g *generator) runWorkers(
	ctx context.Context,
	f func(context.Context, ParsedRule) (map[string][]GeneratedLink, error),
) (map[string]LinksByCategory, error) {
	eg, ctx := errgroup.WithContext(ctx)
	inCh := make(chan ParsedRule)
	outCh := make(chan map[string][]GeneratedLink)

	go func() {
		defer close(inCh)

		rules := g.ruleProvider.Rules()
		for _, rule := range rules {
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

	var openAlarms []AlarmWithData // nolint:prealloc
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

func (g *generator) generateLinksByAlarms(ctx context.Context, rule ParsedRule, alarms []AlarmWithData, user User) (map[string][]GeneratedLink, error) {
	res := make(map[string][]GeneratedLink, len(alarms))
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
			res[alarms[i].ID], err = ExecCode(ctx, rule, args)
			if err != nil {
				g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process alarm")
			}

			continue
		}

		data := g.getTplData(rule, argAlarms, argEntities, user)
		res[alarms[i].ID], err = ExecTemplate(rule, data, g.tplExecutor)
		if err != nil {
			g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process alarm")
		}
	}

	return res, nil
}

func (g *generator) generateLinksByEntities(ctx context.Context, rule ParsedRule, entities []entityWithAlarm, user User) (map[string][]GeneratedLink, error) {
	res := make(map[string][]GeneratedLink, len(entities))
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
			res[entities[i].ID], err = ExecCode(ctx, rule, args)
			if err != nil {
				g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process entity")
			}

			continue
		}

		data := g.getTplData(rule, argAlarms, argEntities, user)
		res[entities[i].ID], err = ExecTemplate(rule, data, g.tplExecutor)
		if err != nil {
			g.logger.Err(err).Str("linkrule", rule.ID).Msg("cannot process entity")
		}
	}

	return res, nil
}

func (g *generator) addExternalData(
	ctx context.Context,
	rule ParsedRule,
	alarms []AlarmWithData,
	entities []EntityWithData,
) error {
	if len(rule.ExternalData) == 0 {
		return nil
	}

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
		data[i].ExternalData, err = g.processExternalData(ctx, externalData, item)
		if err != nil {
			return err
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
		data[i].ExternalData, err = g.processExternalData(ctx, externalData, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) processExternalData(
	ctx context.Context,
	params []externaldata.ParsedRefParameters,
	data any,
) (map[string]any, error) {
	res, err := g.externalDataGetter.Get(ctx, externaldata.Rule{
		ExternalData: params,
	}, data)
	if err != nil {
		return nil, err
	}

	return res.ExternalData, nil
}

func (g *generator) getTplData(
	rule ParsedRule,
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
	rule ParsedRule,
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
