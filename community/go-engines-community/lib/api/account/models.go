package account

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/role"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/user"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	user.User   `bson:",inline"`
	Permissions []role.Permission `bson:"permissions" json:"permissions"`
	UITours     map[string]bool   `bson:"ui_tours" json:"ui_tours"`
}

type EditRequest struct {
	ID                     string          `json:"-"`
	Password               string          `json:"password"`
	UILanguage             string          `json:"ui_language"`
	UIGroupsNavigationType string          `json:"ui_groups_navigation_type"`
	UITheme                string          `json:"ui_theme"`
	DefaultView            string          `json:"defaultview"`
	UITours                map[string]bool `json:"ui_tours"`
}

func (r EditRequest) getUpdateBson(passwordEncoder password.Encoder) bson.M {
	bsonModel := bson.M{
		"ui_language":               r.UILanguage,
		"ui_groups_navigation_type": r.UIGroupsNavigationType,
		"ui_theme":                  r.UITheme,
		"defaultview":               r.DefaultView,
		"ui_tours":                  r.UITours,
	}
	if r.Password != "" {
		bsonModel["password"] = string(passwordEncoder.EncodePassword([]byte(r.Password)))
	}

	return bsonModel
}
