package wrapper

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
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

func (g *generator) GenerateForAlarms(ctx context.Context, ids []string) (map[string]link.LinksByCategory, error) {
	var res map[string]link.LinksByCategory
	for _, v := range g.generators {
		linksById, err := v.GenerateForAlarms(ctx, ids)
		if err != nil {
			return nil, err
		}

		res = link.MergeLinks(res, linksById)
	}

	return res, nil
}

func (g *generator) GenerateForEntities(ctx context.Context, ids []string) (map[string]link.LinksByCategory, error) {
	var res map[string]link.LinksByCategory
	for _, v := range g.generators {
		linksById, err := v.GenerateForEntities(ctx, ids)
		if err != nil {
			return nil, err
		}

		res = link.MergeLinks(res, linksById)
	}

	return res, nil
}

func (g *generator) GenerateCombinedForAlarmsByRule(ctx context.Context, ruleId string, alarmIds []string) ([]link.Link, error) {
	var res []link.Link
	for _, v := range g.generators {
		links, err := v.GenerateCombinedForAlarmsByRule(ctx, ruleId, alarmIds)
		if err != nil {
			return nil, err
		}
		res = append(res, links...)
	}

	return res, nil
}
