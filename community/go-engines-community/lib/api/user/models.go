package user

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" binding:"oneoforempty=_id name role enable"`
}

type EditRequest struct {
	ID                     string          `json:"-"`
	Password               string          `json:"password"`
	Name                   string          `json:"name" binding:"required,max=255"`
	Firstname              string          `json:"firstname" binding:"max=255"`
	Lastname               string          `json:"lastname" binding:"max=255"`
	Email                  string          `json:"email" binding:"required,email"`
	Role                   string          `json:"role" binding:"required"`
	UILanguage             string          `json:"ui_language"`
	UIGroupsNavigationType string          `json:"ui_groups_navigation_type"`
	IsEnabled              *bool           `json:"enable" binding:"required"`
	DefaultView            string          `json:"defaultview"`
	UITours                map[string]bool `json:"ui_tours"`
}

type User struct {
	ID                     string          `bson:"_id" json:"_id"`
	Name                   string          `bson:"name" json:"name"`
	Lastname               string          `bson:"lastname" json:"lastname"`
	Firstname              string          `bson:"firstname" json:"firstname"`
	Email                  string          `bson:"email" json:"email"`
	Role                   Role            `bson:"role" json:"role"`
	UILanguage             string          `bson:"ui_language" json:"ui_language"`
	UITours                map[string]bool `bson:"ui_tours" json:"ui_tours"`
	UIGroupsNavigationType string          `bson:"ui_groups_navigation_type" json:"ui_groups_navigation_type"`
	IsEnabled              bool            `bson:"enable" json:"enable"`
	DefaultView            *View           `bson:"defaultview" json:"defaultview"`
	ExternalID             string          `bson:"external_id" json:"external_id"`
	Source                 string          `bson:"source" json:"source"`
	AuthApiKey             string          `bson:"authkey" json:"authkey"`
}

type Role struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
}

type View struct {
	ID    string `bson:"_id" json:"_id"`
	Title string `bson:"title" json:"title"`
}

type AggregationResult struct {
	Data       []User `bson:"data" json:"data"`
	TotalCount int64  `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
