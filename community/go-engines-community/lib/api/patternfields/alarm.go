package patternfields

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
)

type AlarmRequest struct {
	AlarmPattern          pattern.Alarm `json:"alarm_pattern" binding:"alarm_pattern"`
	CorporateAlarmPattern string        `json:"corporate_alarm_pattern"`

	CorporatePattern savedpattern.SavedPattern `json:"-"`
	IsPrivate        bool                      `json:"-"`
	User             string                    `json:"-"`
}

func (r AlarmRequest) ToModel() savedpattern.AlarmPatternFields {
	if r.CorporatePattern.ID == "" {
		return savedpattern.AlarmPatternFields{
			AlarmPattern: r.AlarmPattern,
		}
	}

	return savedpattern.AlarmPatternFields{
		AlarmPattern:               r.CorporatePattern.AlarmPattern,
		CorporateAlarmPattern:      r.CorporatePattern.ID,
		CorporateAlarmPatternTitle: r.CorporatePattern.Title,
	}
}

func (r AlarmRequest) ToModelWithoutFields(forbiddenFields, onlyTimeAbsoluteFields []string) savedpattern.AlarmPatternFields {
	if r.CorporatePattern.ID == "" {
		return savedpattern.AlarmPatternFields{
			AlarmPattern: r.AlarmPattern,
		}
	}

	return savedpattern.AlarmPatternFields{
		AlarmPattern:               r.CorporatePattern.AlarmPattern.RemoveFields(forbiddenFields, onlyTimeAbsoluteFields),
		CorporateAlarmPattern:      r.CorporatePattern.ID,
		CorporateAlarmPatternTitle: r.CorporatePattern.Title,
	}
}

func ValidateAlarmPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}

	p, ok := i.(pattern.Alarm)

	return ok && match.ValidateAlarmPattern(p, nil, nil)
}

func GetForbiddenFieldsInAlarmPattern(collection string) []string {
	switch collection {
	case mongo.IdleRuleMongoCollection,
		mongo.MetaAlarmRulesMongoCollection,
		mongo.FlappingRuleMongoCollection,
		mongo.ResolveRuleMongoCollection,
		mongo.ScenarioCollection,
		mongo.InstructionMongoCollection,
		mongo.DeclareTicketRuleCollection,
		mongo.LinkRuleMongoCollection:
		return []string{"v.last_event_date", "v.last_update_date", "v.resolved"}
	case mongo.DynamicInfosRulesMongoCollection:
		return []string{"v.last_event_date", "v.last_update_date", "v.resolved", "v.duration", "v.infos"}
	case mongo.AlarmTagCollection:
		return []string{"v.last_event_date", "v.last_update_date", "v.resolved", "v.duration", "tags"}
	default:
		return nil
	}
}

func GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(collection string) []string {
	switch collection {
	case mongo.IdleRuleMongoCollection,
		mongo.DynamicInfosRulesMongoCollection,
		mongo.MetaAlarmRulesMongoCollection,
		mongo.FlappingRuleMongoCollection,
		mongo.ResolveRuleMongoCollection,
		mongo.ScenarioCollection,
		mongo.InstructionMongoCollection,
		mongo.DeclareTicketRuleCollection,
		mongo.LinkRuleMongoCollection,
		mongo.AlarmTagCollection:
		return []string{"v.creation_date", "v.ack.t", "v.activation_date"}
	default:
		return nil
	}
}
