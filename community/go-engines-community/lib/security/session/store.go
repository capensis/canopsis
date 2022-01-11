// Package session contains implementation of http session.
package session

//go:generate mockgen -destination=../../../mocks/github.com/gorilla/sessions/store.go github.com/gorilla/sessions Store

import (
	"context"
	"errors"
	"github.com/gorilla/sessions"
	"time"
)

var ErrNoSession = errors.New("mongo: no session found")

// Store is an interface implemented by store that can clean expired sessions
// and count sessions.
type Store interface {
	sessions.Store
	// StartAutoClean starts a go routine that will every specified duration clean expired sessions.
	StartAutoClean(ctx context.Context, timeout time.Duration)
	// GetActiveSessionsCount returns count of active sessions.
	GetActiveSessionsCount(ctx context.Context) (int64, error)
	ExpireSessions(ctx context.Context, user string, provider string) error
}
