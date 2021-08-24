package websocket

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
)

type Authorizer interface {
	Auth(userId, room string) (bool, error)
}

func NewAuthorizer(enforcer security.Enforcer, roomPerms map[string][]string) Authorizer {
	return &authorizer{
		enforcer:  enforcer,
		roomPerms: roomPerms,
	}
}

type authorizer struct {
	enforcer  security.Enforcer
	roomPerms map[string][]string
}

func (a *authorizer) Auth(userId, room string) (bool, error) {
	perms := a.roomPerms[room]
	if len(perms) == 0 {
		return false, nil
	}

	vals := []interface{}{userId}
	for _, v := range perms {
		vals = append(vals, v)
	}

	return a.enforcer.Enforce(vals...)
}
