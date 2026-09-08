package link

import (
	"context"
	"errors"
)

const jsFuncName = "generate"

var (
	ErrValueIsntSlice      = errors.New("value is not slice")
	ErrValueIsntSliceOrMap = errors.New("value is not slice or map")
	ErrLabelIsMissing      = errors.New("label is missing")
	ErrURLIsMissing        = errors.New("url is missing")
)

type GeneratedLink struct {
	Link
	Category string
}

func ExecCode(ctx context.Context, rule ParsedRule, args []any) ([]GeneratedLink, error) {
	r, err := rule.CodeExecutor.ExecuteFunc(ctx, jsFuncName, args...)
	if err != nil {
		return nil, err
	}

	s, ok := r.([]any)
	if !ok {
		return nil, ErrValueIsntSlice
	}

	res := make([]GeneratedLink, len(s))
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

		res[i] = GeneratedLink{
			Category: category,
			Link: Link{
				RuleID:   rule.ID,
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
