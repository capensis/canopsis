package patternfields

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
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
	Aliases          []string                  `json:"-"`
	IsPrivate        bool                      `json:"-"`
	User             string                    `json:"-"`
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

func (r EntityRequest) ToModelWithoutFields(forbiddenFields []string) savedpattern.EntityPatternFields {
	if r.CorporatePattern.ID == "" {
		return savedpattern.EntityPatternFields{
			EntityPattern: r.EntityPattern,
		}
	}

	return savedpattern.EntityPatternFields{
		EntityPattern:               r.CorporatePattern.EntityPattern.RemoveFields(forbiddenFields),
		CorporateEntityPattern:      r.CorporatePattern.ID,
		CorporateEntityPatternTitle: r.CorporatePattern.Title,
	}
}

type TotalEntityRequest struct {
	TotalEntityPattern          pattern.Entity `json:"total_entity_pattern" binding:"entity_pattern"`
	CorporateTotalEntityPattern string         `json:"corporate_total_entity_pattern"`

	CorporatePattern savedpattern.SavedPattern `json:"-"`
	Aliases          []string                  `json:"-"`
}

func (r TotalEntityRequest) ToModel() correlation.TotalEntityPatternFields {
	if r.CorporatePattern.ID == "" {
		return correlation.TotalEntityPatternFields{
			TotalEntityPattern: r.TotalEntityPattern,
		}
	}

	return correlation.TotalEntityPatternFields{
		TotalEntityPattern:               r.CorporatePattern.EntityPattern,
		CorporateTotalEntityPattern:      r.CorporatePattern.ID,
		CorporateTotalEntityPatternTitle: r.CorporatePattern.Title,
	}
}

func (r TotalEntityRequest) ToModelWithoutFields(forbiddenFields []string) correlation.TotalEntityPatternFields {
	if r.CorporatePattern.ID == "" {
		return correlation.TotalEntityPatternFields{
			TotalEntityPattern: r.TotalEntityPattern,
		}
	}

	return correlation.TotalEntityPatternFields{
		TotalEntityPattern:               r.CorporatePattern.EntityPattern.RemoveFields(forbiddenFields),
		CorporateTotalEntityPattern:      r.CorporatePattern.ID,
		CorporateTotalEntityPatternTitle: r.CorporatePattern.Title,
	}
}

func ValidateEntityPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}
	p, ok := i.(pattern.Entity)
	if !ok {
		return false
	}

	return match.ValidateEntityPattern(p, nil)
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
