package link

import (
	"context"
	"sort"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
)

type Generator interface {
	Load(ctx context.Context) error
	GenerateForAlarm(ctx context.Context, alarm types.Alarm, entity types.Entity, user link.User) []Link
	GenerateForEntity(ctx context.Context, entity types.Entity, user link.User) []Link
}

type Link struct {
	Label    string `json:"label"`
	URL      string `json:"url"`
	Category string `json:"category"`
}

func NewGenerator(dbClient mongo.DbClient, externalDataGetter externaldata.Getter, tplExecutor template.Executor, logger zerolog.Logger) Generator {
	return &generator{
		externalDataGetter: externalDataGetter,
		tplExecutor:        tplExecutor,
		logger:             logger,
		ruleProvider:       link.NewRuleProvider(dbClient, tplExecutor, logger),
	}
}

type generator struct {
	externalDataGetter externaldata.Getter
	tplExecutor        template.Executor
	logger             zerolog.Logger
	ruleProvider       link.RuleProvider
}

func (g *generator) Load(ctx context.Context) error {
	return g.ruleProvider.Load(ctx)
}

func (g *generator) GenerateForAlarm(ctx context.Context, alarm types.Alarm, entity types.Entity, user link.User) (res []Link) {
	res = make([]Link, 0)
	var err error
	rules := g.ruleProvider.Rules()
	for _, rule := range rules {
		if rule.Type != link.TypeAlarm {
			continue
		}

		if !g.matchAlarm(rule, &alarm, &entity) {
			continue
		}

		item := struct {
			types.Alarm
			Entity       types.Entity
			ExternalData map[string]any
		}{
			Alarm:  alarm,
			Entity: entity,
		}
		item.ExternalData, err = g.fetchExternalData(ctx, rule.ExternalData, item)
		if err != nil {
			g.logger.Err(err).Str("rule", rule.ID).Msg("failed to fetch external data for link")
			continue
		}

		if rule.CodeExecutor != nil {
			// func generate(alarms, user, envVar1, envVar2) {}
			args := append([]any{[]any{item}, user}, g.tplVars()...)
			l, err := g.execCode(ctx, rule, args)
			if err != nil {
				g.logger.Err(err).Str("rule", rule.ID).Msg("failed to execute js code for link")
				continue
			}

			res = append(res, l...)
			continue
		}

		l, err := g.execTemplate(rule, map[string]any{
			"Alarms": []any{item},
			"User":   user,
		})
		if err != nil {
			g.logger.Err(err).Str("rule", rule.ID).Msg("failed to execute templates for link")
			continue
		}

		res = append(res, l...)
	}

	return res
}

func (g *generator) GenerateForEntity(ctx context.Context, entity types.Entity, user link.User) (res []Link) {
	res = make([]Link, 0)
	rules := g.ruleProvider.Rules()
	for _, rule := range rules {
		if rule.Type != link.TypeEntity {
			continue
		}

		if !g.matchEntity(rule, &entity) {
			continue
		}

		item := struct {
			types.Entity
			ExternalData map[string]any
		}{
			Entity: entity,
		}
		var err error
		item.ExternalData, err = g.fetchExternalData(ctx, rule.ExternalData, item)
		if err != nil {
			g.logger.Err(err).Str("rule", rule.ID).Msg("failed to fetch external data for link")
			continue
		}

		if rule.CodeExecutor != nil {
			// func generate(entities, user, envVar1, envVar2) {}
			args := append([]any{[]any{item}, user}, g.tplVars()...)
			l, err := g.execCode(ctx, rule, args)
			if err != nil {
				g.logger.Err(err).Str("rule", rule.ID).Msg("failed to execute js code for link")
				continue
			}

			res = append(res, l...)
			continue
		}

		l, err := g.execTemplate(rule, map[string]any{
			"Entities": []any{item},
			"User":     user,
		})
		if err != nil {
			g.logger.Err(err).Str("rule", rule.ID).Msg("failed to execute templates for link")
			continue
		}

		res = append(res, l...)
	}

	return res
}

func (g *generator) matchAlarm(rule link.ParsedRule, alarm *types.Alarm, entity *types.Entity) bool {
	ok, err := match.MatchAlarmPattern(rule.AlarmPattern, alarm)
	if err != nil {
		g.logger.Err(err).Str("rule", rule.ID).Msg("invalid alarm pattern in link rule")

		return false
	}

	if !ok {
		return false
	}

	return g.matchEntity(rule, entity)
}

func (g *generator) matchEntity(rule link.ParsedRule, entity *types.Entity) bool {
	ok, err := match.MatchEntityPattern(rule.EntityPattern, entity)
	if err != nil {
		g.logger.Err(err).Str("rule", rule.ID).Msg("invalid entity pattern in link rule")

		return false
	}

	return ok
}

func (g *generator) fetchExternalData(ctx context.Context, params []externaldata.ParsedRefParameters, data any) (map[string]any, error) {
	res, err := g.externalDataGetter.Get(ctx, externaldata.Rule{
		ExternalData: params,
	}, data)
	if err != nil {
		return nil, err
	}

	return res.ExternalData, nil
}

func (g *generator) execCode(ctx context.Context, rule link.ParsedRule, args []any) ([]Link, error) {
	gl, err := link.ExecCode(ctx, rule, args)
	if err != nil {
		return nil, err
	}

	return g.transformLinks(gl), nil
}

func (g *generator) execTemplate(rule link.ParsedRule, data any) ([]Link, error) {
	gl, err := link.ExecTemplate(rule, data, g.tplExecutor)
	if err != nil {
		return nil, err
	}

	return g.transformLinks(gl), nil
}

func (g *generator) transformLinks(gl []link.GeneratedLink) []Link {
	res := make([]Link, len(gl))
	for i, l := range gl {
		res[i] = Link{
			Label:    l.Label,
			URL:      l.Url,
			Category: l.Category,
		}
	}

	return res
}

func (g *generator) tplVars() []any {
	tplVars := g.tplExecutor.GetDefaultTplVars()
	keys := make([]string, 0, len(tplVars))
	for k := range tplVars {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	res := make([]any, len(keys))
	for i, k := range keys {
		res[i] = tplVars[k]
	}

	return res
}
