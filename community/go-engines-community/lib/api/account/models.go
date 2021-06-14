package account

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/role"
	"git.canopsis.net/canopsis/go-engines/lib/api/user"
)

type User struct {
	user.User   `bson:",inline"`
	Permissions []role.Permission `bson:"permissions" json:"permissions"`
}
