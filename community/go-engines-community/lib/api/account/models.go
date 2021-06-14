package account

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/role"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/user"
)

type User struct {
	user.User   `bson:",inline"`
	Permissions []role.Permission `bson:"permissions" json:"permissions"`
}
