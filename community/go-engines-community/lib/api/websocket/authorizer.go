package websocket

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
)

var ErrNotFoundRoom = errors.New("no room")
var ErrNotFoundRoomInGroup = errors.New("no group room")

// Authorizer is used to implement websocket room authentication and authorization.
type Authorizer interface {
	// Authenticate authenticates user by token.
	Authenticate(ctx context.Context, token string) (string, error)
	// Authorize checks if user has access to room.
	Authorize(ctx context.Context, userID, room string) (bool, error)
	// AddRoom adds room with permissions.
	AddRoom(room string, perms []string) error
	AddGroup(group string, perms []string, check GroupCheckExists) error
	GetGroupIds(group string) []string
	RemoveGroupRoom(group, id string) error
}

func NewAuthorizer(
	enforcer security.Enforcer,
	tokenProviders []security.TokenProvider,
) Authorizer {
	return &authorizer{
		enforcer:       enforcer,
		tokenProviders: tokenProviders,
		roomPerms:      make(map[string][]string),
		groupPerms:     make(map[string][]string),
		groupChecks:    make(map[string]GroupCheckExists),
		roomsByGroup:   make(map[string][]string),
	}
}

type authorizer struct {
	enforcer       security.Enforcer
	tokenProviders []security.TokenProvider
	mx             sync.RWMutex
	roomPerms      map[string][]string
	groupPerms     map[string][]string
	groupChecks    map[string]GroupCheckExists
	roomsByGroup   map[string][]string
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

func (a *authorizer) Authorize(ctx context.Context, userID, room string) (bool, error) {
	ok, err := a.authorizeRoom(userID, room)
	if err != nil {
		if errors.Is(err, ErrNotFoundRoom) {
			return a.authorizeGroupRoom(ctx, userID, room)
		}

		return false, err
	}

	return ok, nil
}

func (a *authorizer) AddRoom(room string, perms []string) error {
	a.mx.Lock()
	defer a.mx.Unlock()

	if _, ok := a.roomPerms[room]; ok {
		return fmt.Errorf("%q room already exists", room)
	}

	a.roomPerms[room] = perms
	return nil
}

func (a *authorizer) AddGroup(group string, perms []string, check GroupCheckExists) error {
	a.mx.Lock()
	defer a.mx.Unlock()

	if _, ok := a.groupPerms[group]; ok {
		return fmt.Errorf("%q group already exists", group)
	}

	a.groupPerms[group] = perms
	if check != nil {
		a.groupChecks[group] = check
	}
	return nil
}

func (a *authorizer) GetGroupIds(group string) []string {
	a.mx.RLock()
	defer a.mx.RUnlock()

	return a.roomsByGroup[group]
}

func (a *authorizer) RemoveGroupRoom(group, id string) error {
	a.mx.Lock()
	defer a.mx.Unlock()

	room := group + id
	if _, ok := a.roomPerms[room]; !ok {
		return nil
	}

	delete(a.roomPerms, room)

	k := 0
	for _, v := range a.roomsByGroup[group] {
		if id != v {
			a.roomsByGroup[group][k] = v
			k++
		}
	}

	a.roomsByGroup[group] = a.roomsByGroup[group][:k]
	return nil
}

func (a *authorizer) authorizeRoom(userID, room string) (bool, error) {
	a.mx.RLock()
	defer a.mx.RUnlock()

	if perms, ok := a.roomPerms[room]; ok {
		// Return authorized if room doesn't have permissions.
		if len(perms) == 0 {
			return true, nil
		}

		if userID == "" {
			return false, nil
		}

		vals := []any{userID}
		for _, v := range perms {
			vals = append(vals, v)
		}

		return a.enforcer.Enforce(vals...)
	}

	return false, ErrNotFoundRoom
}

func (a *authorizer) authorizeGroupRoom(ctx context.Context, userID, room string) (bool, error) {
	a.mx.Lock()
	defer a.mx.Unlock()

	for group, perms := range a.groupPerms {
		if !strings.HasPrefix(room, group) {
			continue
		}

		id := room[len(group):]
		if check, ok := a.groupChecks[group]; ok {
			ok, err := check(ctx, id)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, ErrNotFoundRoomInGroup
			}
		}

		// Return authorized if room doesn't have permissions.
		if len(perms) == 0 {
			a.roomPerms[room] = perms
			a.roomsByGroup[group] = append(a.roomsByGroup[group], id)
			return true, nil
		}

		if userID == "" {
			return false, nil
		}

		vals := []any{userID}
		for _, v := range perms {
			vals = append(vals, v)
		}

		ok, err := a.enforcer.Enforce(vals...)
		if err != nil || !ok {
			return false, err
		}

		a.roomPerms[room] = perms
		a.roomsByGroup[group] = append(a.roomsByGroup[group], id)

		return true, nil
	}

	return false, ErrNotFoundRoom
}
