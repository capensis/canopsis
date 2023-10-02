package wrapper

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewGenerator(generators ...link.Generator) link.Generator {
	return &generator{
		generators: generators,
	}
}

type generator struct {
	generators []link.Generator
}

func (g *generator) Load(ctx context.Context) error {
	for _, v := range g.generators {
		err := v.Load(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) GenerateForAlarm(ctx context.Context, alarm types.Alarm, entity types.Entity, user link.User) (link.LinksByCategory, error) {
	var res link.LinksByCategory
	for _, v := range g.generators {
		linksByCategory, err := v.GenerateForAlarm(ctx, alarm, entity, user)
		if err != nil {
			return nil, err
		}

		if res == nil {
			res = linksByCategory
		} else {
			for category, links := range linksByCategory {
				res[category] = append(res[category], links...)
			}
		}
	}

	return res, nil
}

func (g *generator) GenerateForAlarms(ctx context.Context, ids []string, user link.User) (map[string]link.LinksByCategory, error) {
	var res map[string]link.LinksByCategory
	for _, v := range g.generators {
		linksById, err := v.GenerateForAlarms(ctx, ids, user)
		if err != nil {
			return nil, err
		}

		res = link.MergeLinks(res, linksById)
	}

	return res, nil
}

func (g *generator) GenerateForEntities(ctx context.Context, ids []string, user link.User) (map[string]link.LinksByCategory, error) {
	var res map[string]link.LinksByCategory
	for _, v := range g.generators {
		linksById, err := v.GenerateForEntities(ctx, ids, user)
		if err != nil {
			return nil, err
		}

		res = link.MergeLinks(res, linksById)
	}

	return res, nil
}

func (g *generator) GenerateCombinedForAlarmsByRule(ctx context.Context, ruleId string, alarmIds []string, user link.User) ([]link.Link, error) {
	var res []link.Link
	for _, v := range g.generators {
		links, err := v.GenerateCombinedForAlarmsByRule(ctx, ruleId, alarmIds, user)
		if err != nil {
			return nil, err
		}
		res = append(res, links...)
	}

	return res, nil
}
