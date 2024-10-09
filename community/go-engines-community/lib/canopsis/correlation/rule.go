package correlation

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	RuleTypeRelation    = "relation"
	RuleTypeTimeBased   = "timebased"
	RuleTypeAttribute   = "attribute"
	RuleTypeComplex     = "complex"
	RuleTypeValueGroup  = "valuegroup"
	RuleTypeManualGroup = "manualgroup"
	RuleTypeCorel       = "corel"
)

type Rule struct {
	ID             string     `bson:"_id" json:"_id"`
	Type           string     `bson:"type" json:"type"`
	Name           string     `bson:"name" json:"name"`
	Author         string     `bson:"author" json:"author"`
	OutputTemplate string     `bson:"output_template" json:"output_template"`
	Config         RuleConfig `bson:"config" json:"config"`
	AutoResolve    bool       `bson:"auto_resolve" json:"auto_resolve"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
	TotalEntityPatternFields         `bson:",inline"`

	Created *types.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated *types.CpsTime `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`

	OldAlarmPatterns       oldpattern.AlarmPatternList  `bson:"old_alarm_patterns,omitempty" json:"old_alarm_patterns,omitempty"`
	OldEntityPatterns      oldpattern.EntityPatternList `bson:"old_entity_patterns,omitempty" json:"old_entity_patterns,omitempty"`
	OldTotalEntityPatterns oldpattern.EntityPatternList `bson:"old_total_entity_patterns,omitempty" json:"old_total_entity_patterns,omitempty"`
	OldEventPatterns       oldpattern.EventPatternList  `bson:"old_event_patterns,omitempty" json:"old_event_patterns,omitempty"`
}

func (r *Rule) Matches(event types.Event, alarmWithEntity types.AlarmWithEntity) (bool, error) {
	if !r.OldEventPatterns.IsSet() && !r.OldEntityPatterns.IsSet() && !r.OldAlarmPatterns.IsSet() &&
		len(r.EntityPattern) == 0 && len(r.AlarmPattern) == 0 {
		switch r.Type {
		case RuleTypeRelation,
			RuleTypeTimeBased,
			RuleTypeComplex,
			RuleTypeValueGroup:
			return true, nil
		}
	}

	if r.OldEventPatterns.IsSet() {
		if !r.OldEventPatterns.IsValid() {
			return false, pattern.ErrInvalidOldEventPattern
		}

		if !r.OldEventPatterns.Matches(event) {
			return false, nil
		}
	}

	return pattern.Match(alarmWithEntity.Entity, alarmWithEntity.Alarm, r.EntityPattern, r.AlarmPattern, r.OldEntityPatterns, r.OldAlarmPatterns)
}

func (r *Rule) IsManual() bool {
	return r.Type == RuleTypeManualGroup
}

func (r *Rule) GetStateID(group string) string {
	if group != "" {
		return r.ID + "&&" + group
	}

	return r.ID
}

type TotalEntityPatternFields struct {
	TotalEntityPattern pattern.Entity `bson:"total_entity_pattern" json:"total_entity_pattern,omitempty"`

	CorporateTotalEntityPattern      string `bson:"corporate_total_entity_pattern" json:"corporate_total_entity_pattern,omitempty"`
	CorporateTotalEntityPatternTitle string `bson:"corporate_total_entity_pattern_title" json:"corporate_total_entity_pattern_title,omitempty"`
}
