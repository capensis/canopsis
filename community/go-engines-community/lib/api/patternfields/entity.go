package patternfields

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
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

func ValidateEntityPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}

	p, ok := i.(pattern.Entity)

	return ok && match.ValidateEntityPattern(p, nil)
}

func GetForbiddenFieldsInEntityPattern(collection string) []string {
	switch collection {
	case mongo.StateSettingsMongoCollection:
		return []string{"last_event_date", "component", "component_infos"}
	case mongo.EntityMongoCollection:
		return []string{"last_event_date", "connector", "component_infos"}
	case mongo.PbehaviorMongoCollection,
		mongo.IdleRuleMongoCollection,
		mongo.DynamicInfosRulesMongoCollection,
		mongo.MetaAlarmRulesMongoCollection,
		mongo.FlappingRuleMongoCollection,
		mongo.ResolveRuleMongoCollection,
		mongo.ScenarioCollection,
		mongo.InstructionMongoCollection,
		mongo.KpiFilterMongoCollection,
		mongo.DeclareTicketRuleCollection,
		mongo.LinkRuleMongoCollection,
		mongo.AlarmTagCollection:
		return []string{"last_event_date"}
	default:
		return nil
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
