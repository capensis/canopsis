package user

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/colortheme"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy     string `form:"sort_by" binding:"oneoforempty=_id name enable source"`
	Permission string `form:"permission"`
}

type CreateRequest struct {
	EditRequest
	Source     string `json:"source" binding:"oneoforempty=ldap cas saml"`
	ExternalID string `json:"external_id"`
}

type UpdateRequest struct {
	ID string `json:"-"`
	EditRequest
}

type EditRequest struct {
	Password               string   `json:"password"`
	Name                   string   `json:"name" binding:"required,max=255"`
	Firstname              string   `json:"firstname" binding:"max=255"`
	Lastname               string   `json:"lastname" binding:"max=255"`
	Email                  string   `json:"email" binding:"required,email"`
	Roles                  []string `json:"roles" binding:"required,notblank"`
	UILanguage             string   `json:"ui_language" binding:"max=255"`
	UIGroupsNavigationType string   `json:"ui_groups_navigation_type" binding:"max=255"`
	UITheme                string   `json:"ui_theme" binding:"max=255"`
	IsEnabled              *bool    `json:"enable" binding:"required"`
	DefaultView            string   `json:"defaultview"`
}

func (r CreateRequest) getBson(passwordEncoder password.Encoder) bson.M {
	bsonModel := bson.M{
		"_id":                       r.Name,
		"name":                      r.Name,
		"lastname":                  r.Lastname,
		"firstname":                 r.Firstname,
		"email":                     r.Email,
		"roles":                     r.Roles,
		"password":                  string(passwordEncoder.EncodePassword([]byte(r.Password))),
		"ui_language":               r.UILanguage,
		"ui_theme":                  r.UITheme,
		"ui_groups_navigation_type": r.UIGroupsNavigationType,
		"enable":                    r.IsEnabled,
		"defaultview":               r.DefaultView,
		"authkey":                   utils.NewID(),
		"source":                    r.Source,
		"external_id":               r.ExternalID,
	}

	return bsonModel
}

func (r EditRequest) getBson(passwordEncoder password.Encoder) bson.M {
	bsonModel := bson.M{
		"name":                      r.Name,
		"lastname":                  r.Lastname,
		"firstname":                 r.Firstname,
		"email":                     r.Email,
		"roles":                     r.Roles,
		"ui_language":               r.UILanguage,
		"ui_theme":                  r.UITheme,
		"ui_groups_navigation_type": r.UIGroupsNavigationType,
		"enable":                    r.IsEnabled,
		"defaultview":               r.DefaultView,
	}
	if r.Password != "" {
		bsonModel["password"] = string(passwordEncoder.EncodePassword([]byte(r.Password)))
	}

	return bsonModel
}

type User struct {
	ID                     string           `bson:"_id" json:"_id"`
	Name                   string           `bson:"name" json:"name"`
	DisplayName            string           `bson:"display_name" json:"display_name"`
	Lastname               string           `bson:"lastname" json:"lastname"`
	Firstname              string           `bson:"firstname" json:"firstname"`
	Email                  string           `bson:"email" json:"email"`
	Roles                  []Role           `bson:"roles" json:"roles"`
	UILanguage             string           `bson:"ui_language" json:"ui_language"`
	UITheme                colortheme.Theme `bson:"ui_theme" json:"ui_theme"`
	UIGroupsNavigationType string           `bson:"ui_groups_navigation_type" json:"ui_groups_navigation_type"`
	Enabled                bool             `bson:"enable" json:"enable"`
	DefaultView            *View            `bson:"defaultview" json:"defaultview"`
	ExternalID             string           `bson:"external_id" json:"external_id"`
	Source                 string           `bson:"source" json:"source"`
	AuthApiKey             string           `bson:"authkey" json:"authkey"`

	ActiveConnects *int64 `bson:"-" json:"active_connects,omitempty"`
}

type Role struct {
	ID          string `bson:"_id" json:"_id"`
	Name        string `bson:"name" json:"name"`
	DefaultView *View  `bson:"defaultview" json:"defaultview"`
}

type View struct {
	ID    string `bson:"_id" json:"_id"`
	Title string `bson:"title" json:"title"`
}

type BulkUpdateRequestItem struct {
	ID string `json:"_id" binding:"required"`
	EditRequest
}

type BulkDeleteRequestItem struct {
	ID string `json:"_id" binding:"required"`
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
