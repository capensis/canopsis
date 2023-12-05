package link

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
)

const (
	TypeAlarm  = "alarm"
	TypeEntity = "entity"
)

const (
	ActionOpen = "open"
	ActionCopy = "copy"
)

type Rule struct {
	ID         string           `bson:"_id"`
	Name       string           `bson:"name"`
	Type       string           `bson:"type"`
	Enabled    bool             `bson:"enabled"`
	Links      []Parameters     `bson:"links"`
	SourceCode string           `bson:"source_code"`
	Author     string           `bson:"author"`
	Created    datetime.CpsTime `bson:"created,omitempty"`
	Updated    datetime.CpsTime `bson:"updated,omitempty"`

	ExternalData map[string]ExternalDataParameters `bson:"external_data"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
}

type Parameters struct {
	Label    string `bson:"label" json:"label" binding:"required,max=255"`
	Category string `bson:"category" json:"category" binding:"max=255"`
	IconName string `bson:"icon_name" json:"icon_name" binding:"required,max=255"`
	Url      string `bson:"url" json:"url" binding:"required,max=1000"`
	Action   string `bson:"action" json:"action" binding:"required,oneof=open copy"`
	// Single to mark links unavailable to multiple selected alarms
	Single     bool `bson:"single,omitempty" json:"single,omitempty"`
	HideInMenu bool `bson:"hide_in_menu,omitempty" json:"hide_in_menu,omitempty"`
}

type ExternalDataParameters struct {
	Type string `bson:"type" json:"type" binding:"required,oneof=mongo"`

	Collection string            `bson:"collection" json:"collection" binding:"required"`
	Select     map[string]string `bson:"select" json:"select"`
	Regexp     map[string]string `bson:"regexp" json:"regexp"`
	SortBy     string            `bson:"sort_by" json:"sort_by"`
	Sort       string            `bson:"sort" json:"sort" binding:"oneoforempty=asc desc"`
}

type Link struct {
	RuleID     string `json:"rule_id,omitempty"`
	Label      string `json:"label"`
	IconName   string `json:"icon_name"`
	Url        string `json:"url"`
	Single     bool   `json:"single,omitempty"`
	HideInMenu bool   `json:"hide_in_menu,omitempty"`
	Action     string `json:"action"`
}

type LinksByCategory map[string][]Link

func MergeLinks(left, right map[string]LinksByCategory) map[string]LinksByCategory {
	if left == nil {
		return right
	}

	for id, linksByCategory := range right {
		if left[id] == nil {
			left[id] = linksByCategory
		} else {
			for category, links := range linksByCategory {
				left[id][category] = append(left[id][category], links...)
			}
		}
	}

	return left
}
