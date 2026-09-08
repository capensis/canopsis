package link

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
)

func ExecTemplate(rule ParsedRule, data any, tplExecutor template.Executor) ([]GeneratedLink, error) {
	res := make([]GeneratedLink, len(rule.Links))
	for i, link := range rule.Links {
		url, err := tplExecutor.ExecuteByTpl(link.UrlTpl, data)
		if err != nil {
			return nil, err
		}

		label, err := tplExecutor.ExecuteByTpl(link.LabelTpl, data)
		if err != nil {
			return nil, err
		}

		res[i] = GeneratedLink{
			Category: link.Category,
			Link: Link{
				RuleID:     rule.ID,
				Label:      label,
				IconName:   link.IconName,
				Url:        url,
				Action:     link.Action,
				Single:     link.Single,
				HideInMenu: link.HideInMenu,
			},
		}
	}

	return res, nil
}
