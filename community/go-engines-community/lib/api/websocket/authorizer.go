package websocket

import (
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
)

var ErrUserNotFound = errors.New("user not found")

type Authorizer interface {
	Exists(userId string) (bool, error)
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

func (a *authorizer) Exists(userId string) (bool, error) {
	roles, err := a.enforcer.GetRolesForUser(userId)
	if err != nil {
		return false, err
	}

	return len(roles) > 0, nil
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

	roles, err := a.enforcer.GetRolesForUser(userId)
	if err != nil {
		return false, err
	}

	if len(roles) == 0 {
		return false, ErrUserNotFound
	}

	return a.enforcer.Enforce(vals...)
}
