package websocket

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

// groupDelimiter separates the group name from the room ID in a composite room key.
// For example, room "alarms/abc123" belongs to group "alarms" with ID "abc123".
const groupDelimiter = "/"

var ErrRoomNotFound = errors.New("room not found")

type RoomRegistry interface {
	Register(room string, handlers RoomHandlers) error
	RegisterGroup(group string, handlers RoomHandlers) error
	Get(room string) (RoomHandlers, bool)
}

type Authorize func(ctx context.Context, userID string) (bool, error)

type RoomHandlers struct {
	OnJoin    func(ctx context.Context, opts JoinOptions) error
	OnLeave   func(ctx context.Context, opts LeaveOptions) error
	OnMessage func(ctx context.Context, opts MessageOptions) error
	Authorize Authorize
}

type JoinOptions struct {
	ConnID, UserID string
	RoomID         string
	Payload        any
}

type LeaveOptions struct {
	ConnID, UserID string
	RoomID         string
}

type MessageOptions struct {
	ConnID, UserID string
	RoomID         string
	Payload        any
}

type JoinError struct {
	Err     error
	Payload any
}

func (e *JoinError) Error() string {
	return e.Err.Error()
}

func NewJoinError(err error, payload any) *JoinError {
	return &JoinError{Err: err, Payload: payload}
}

func NewRoomRegistry() RoomRegistry {
	return &roomRegistry{
		rooms:  make(map[string]RoomHandlers),
		groups: make(map[string]RoomHandlers),
	}
}

type roomRegistry struct {
	mu     sync.RWMutex
	rooms  map[string]RoomHandlers
	groups map[string]RoomHandlers
}

func (r *roomRegistry) Register(room string, handlers RoomHandlers) error {
	if room == "" {
		return errors.New("room name cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.rooms[room]; exists {
		return fmt.Errorf("room %q already registered", room)
	}

	r.rooms[room] = handlers

	return nil
}

// RegisterGroup registers a handler for a family of rooms sharing a common prefix.
// Any room key of the form "<group>/<id>" will match, with opts.RoomID set to
// the id part when handlers are invoked.
func (r *roomRegistry) RegisterGroup(group string, handlers RoomHandlers) error {
	if group == "" {
		return errors.New("group name cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.groups[group]; exists {
		return fmt.Errorf("group %q already registered", group)
	}

	r.groups[group] = handlers

	return nil
}

func (r *roomRegistry) Get(room string) (RoomHandlers, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.rooms[room]
	if ok {
		return h, true
	}

	g, roomID, ok := strings.Cut(room, groupDelimiter)
	if ok {
		h, ok = r.groups[g]
		if ok {
			return r.wrapHandlers(roomID, h), true
		}
	}

	return h, ok
}

func (r *roomRegistry) wrapHandlers(roomID string, h RoomHandlers) RoomHandlers {
	if h.OnJoin != nil {
		onJoin := h.OnJoin
		h.OnJoin = func(ctx context.Context, opts JoinOptions) error {
			opts.RoomID = roomID

			return onJoin(ctx, opts)
		}
	}

	if h.OnLeave != nil {
		onLeave := h.OnLeave
		h.OnLeave = func(ctx context.Context, opts LeaveOptions) error {
			opts.RoomID = roomID

			return onLeave(ctx, opts)
		}
	}

	if h.OnMessage != nil {
		onMessage := h.OnMessage
		h.OnMessage = func(ctx context.Context, opts MessageOptions) error {
			opts.RoomID = roomID

			return onMessage(ctx, opts)
		}
	}

	return h
}
