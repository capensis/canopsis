package websocket

import (
	"context"
	"fmt"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
)

// Authorizer is used to implement websocket room authentication and authorization.
type Authorizer interface {
	// Authenticate authenticates user by token.
	Authenticate(ctx context.Context, token string) (string, error)
	// Authorize checks if user has access to room.
	Authorize(userId, room string) (bool, error)
	// AddRoom adds room with permissions.
	AddRoom(room string, perms []string) error
	// RemoveRoom removes room.
	RemoveRoom(room string) error
	HasRoom(room string) bool
}

func NewAuthorizer(
	enforcer security.Enforcer,
	tokenProviders []security.TokenProvider,
) Authorizer {
	return &authorizer{
		enforcer:       enforcer,
		tokenProviders: tokenProviders,
		roomPerms:      make(map[string][]string),
	}
}

type authorizer struct {
	enforcer       security.Enforcer
	tokenProviders []security.TokenProvider
	roomPermsMx    sync.RWMutex
	roomPerms      map[string][]string
}

func (a *authorizer) Authenticate(ctx context.Context, token string) (string, error) {
	for _, provider := range a.tokenProviders {
		user, err := provider.Auth(ctx, token)
		if err != nil {
			return "", err
		}

		if user != nil {
			return user.ID, nil
		}
	}

	return "", nil
}

func (a *authorizer) Authorize(userId, room string) (bool, error) {
	a.roomPermsMx.RLock()
	defer a.roomPermsMx.RUnlock()
	perms, ok := a.roomPerms[room]
	// Return unauthorized if room is missing.
	if !ok {
		return false, nil
	}

	// Return authorized if room doesn't have permissions.
	if len(perms) == 0 {
		return true, nil
	}

	if userId == "" {
		return false, nil
	}

	vals := []interface{}{userId}
	for _, v := range perms {
		vals = append(vals, v)
	}

	return a.enforcer.Enforce(vals...)
}

func (a *authorizer) AddRoom(room string, perms []string) error {
	a.roomPermsMx.Lock()
	defer a.roomPermsMx.Unlock()

	if _, ok := a.roomPerms[room]; ok {
		return fmt.Errorf("%q room already exists", room)
	}

	a.roomPerms[room] = perms
	return nil
}

func (a *authorizer) RemoveRoom(room string) error {
	a.roomPermsMx.Lock()
	defer a.roomPermsMx.Unlock()

	if _, ok := a.roomPerms[room]; !ok {
		return fmt.Errorf("%q room doesn't exists", room)
	}

	delete(a.roomPerms, room)
	return nil
}

func (a *authorizer) HasRoom(room string) bool {
	a.roomPermsMx.RLock()
	defer a.roomPermsMx.RUnlock()

	_, ok := a.roomPerms[room]
	return ok
}
