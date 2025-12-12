package patternfields

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
)

type EntityRequest struct {
	EntityPattern          pattern.Entity `json:"entity_pattern" binding:"entity_pattern"`
	CorporateEntityPattern string         `json:"corporate_entity_pattern"`

	CorporatePattern savedpattern.SavedPattern `json:"-"`
}

func (r EntityRequest) ToModel() savedpattern.EntityPatternFields {
	if r.CorporatePattern.ID == "" {
		return savedpattern.EntityPatternFields{
			EntityPattern: r.EntityPattern,
		}
	}

	return savedpattern.EntityPatternFields{
		EntityPattern:               r.CorporatePattern.EntityPattern,
		CorporateEntityPattern:      r.CorporatePattern.ID,
		CorporateEntityPatternTitle: r.CorporatePattern.Title,
	}
}

func (r EntityRequest) ToModelWithoutFields(collectionName string) savedpattern.EntityPatternFields {
	if r.CorporatePattern.ID == "" {
		return savedpattern.EntityPatternFields{
			EntityPattern: r.EntityPattern,
		}
	}

	forbiddenFields := GetForbiddenFieldsInEntityPattern(collectionName)

	return savedpattern.EntityPatternFields{
		EntityPattern:               r.CorporatePattern.EntityPattern.RemoveFields(forbiddenFields),
		CorporateEntityPattern:      r.CorporatePattern.ID,
		CorporateEntityPatternTitle: r.CorporatePattern.Title,
	}
}

func GetAliases(p pattern.Entity) []string {
	aliases := make([]string, 0)
	for _, g := range p {
		for _, c := range g {
			if c.Alias != "" {
				aliases = append(aliases, c.Alias)
			}
		}
	}

	return aliases
}
