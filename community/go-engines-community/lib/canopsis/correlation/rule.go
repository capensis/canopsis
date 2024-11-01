package correlation

import (
	"cmp"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
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
	ID             string                      `bson:"_id,omitempty" json:"_id,omitempty"`
	Type           string                      `bson:"type" json:"type"`
	Name           string                      `bson:"name" json:"name"`
	Author         string                      `bson:"author" json:"author"`
	OutputTemplate string                      `bson:"output_template" json:"output_template"`
	Config         RuleConfig                  `bson:"config" json:"config"`
	Tags           types.CorrelationRuleTags   `bson:"tags" json:"tags"`
	Infos          []types.CorrelationRuleInfo `bson:"infos" json:"infos"`
	AutoResolve    bool                        `bson:"auto_resolve" json:"auto_resolve"`

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

func GetMetaAlarmComponentAndResource(
	extraInfosMeta EventExtraInfosMeta,
	templateExecutor template.Executor,
	logger zerolog.Logger,
) (string, string) {
	component := ""
	resource := ""
	var err error
	if extraInfosMeta.Rule.Config.ComponentTemplate != "" {
		component, err = templateExecutor.Execute(extraInfosMeta.Rule.Config.ComponentTemplate, extraInfosMeta)
		if err != nil {
			logger.Warn().Err(err).Str("rule", extraInfosMeta.Rule.ID).Msg("invalid component template")
		}
	}

	if extraInfosMeta.Rule.Config.ResourceTemplate != "" {
		resource, err = templateExecutor.Execute(extraInfosMeta.Rule.Config.ResourceTemplate, extraInfosMeta)
		if err != nil {
			logger.Warn().Err(err).Str("rule", extraInfosMeta.Rule.ID).Msg("invalid resource template")
		}

		resource += utils.NewID()
	}

	defaultComponent, defaultResource := GetMetaAlarmDefaultComponentAndResource()

	return cmp.Or(component, defaultComponent), cmp.Or(resource, defaultResource)
}

func GetMetaAlarmDefaultComponentAndResource() (string, string) {
	return DefaultMetaAlarmComponent, DefaultMetaAlarmEntityPrefix + utils.NewID()
}

type TotalEntityPatternFields struct {
	TotalEntityPattern pattern.Entity `bson:"total_entity_pattern" json:"total_entity_pattern,omitempty"`

	CorporateTotalEntityPattern      string `bson:"corporate_total_entity_pattern" json:"corporate_total_entity_pattern,omitempty"`
	CorporateTotalEntityPatternTitle string `bson:"corporate_total_entity_pattern_title" json:"corporate_total_entity_pattern_title,omitempty"`
}
