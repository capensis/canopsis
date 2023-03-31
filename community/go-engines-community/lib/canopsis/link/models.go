package link

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	TypeAlarm  = "alarm"
	TypeEntity = "entity"
)

type Rule struct {
	ID         string        `bson:"_id"`
	Name       string        `bson:"name"`
	Type       string        `bson:"type"`
	Enabled    bool          `bson:"enabled"`
	Links      []Parameters  `bson:"links"`
	SourceCode string        `bson:"source_code"`
	Author     string        `bson:"author"`
	Created    types.CpsTime `bson:"created,omitempty"`
	Updated    types.CpsTime `bson:"updated,omitempty"`

	ExternalData map[string]ExternalDataParameters `bson:"external_data"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
}

type Parameters struct {
	Label    string `bson:"label" json:"label" binding:"required,max=255"`
	Category string `bson:"category" json:"category" binding:"max=255"`
	IconName string `bson:"icon_name" json:"icon_name" binding:"required,max=255"`
	Url      string `bson:"url" json:"url" binding:"required,max=1000"`
	// Single to mark links unavailable to multiple selected alarms
	Single bool `bson:"single,omitempty" json:"single,omitempty"`
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
	RuleID   string `json:"rule_id,omitempty"`
	Label    string `json:"label"`
	IconName string `json:"icon_name"`
	Url      string `json:"url"`
	Single   bool   `json:"single,omitempty"`
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
