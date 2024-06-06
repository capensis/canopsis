package correlation

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
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

	Created *datetime.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated *datetime.CpsTime `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
}

func (r *Rule) Matches(alarmWithEntity types.AlarmWithEntity) (bool, error) {
	if len(r.EntityPattern) == 0 && len(r.AlarmPattern) == 0 {
		switch r.Type {
		case RuleTypeRelation,
			RuleTypeTimeBased,
			RuleTypeComplex,
			RuleTypeValueGroup:
			return true, nil
		}
	}

	return match.Match(&alarmWithEntity.Entity, &alarmWithEntity.Alarm, r.EntityPattern, r.AlarmPattern)
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
