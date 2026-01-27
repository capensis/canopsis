package broadcastmessage

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"github.com/rs/zerolog"
)

// Service is used to implement websocket room for broadcast messages.
type Service interface {
	Start(ctx context.Context, ch <-chan bool)
}

func NewService(
	store Store,
	websocketHub websocket.Hub,
	interval time.Duration,
	logger zerolog.Logger,
) Service {
	return &service{
		store:        store,
		websocketHub: websocketHub,
		interval:     interval,
		logger:       logger,
	}
}

type service struct {
	store        Store
	websocketHub websocket.Hub
	interval     time.Duration
	logger       zerolog.Logger
}

func (s *service) Start(ctx context.Context, ch <-chan bool) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-ch:
			if !ok {
				return
			}

			s.check(ctx)
		case <-ticker.C:
			s.check(ctx)
		}
	}
}

func (s *service) check(ctx context.Context) {
	room := websocket.RoomBroadcastMessages
	conns := s.websocketHub.GetConnectionsByRoom(room)
	for _, c := range conns {
		messages, err := s.store.GetActive(ctx, c.UserID)
		if err != nil {
			s.logger.Err(err).Msg("cannot fetch messages")

			return
		}

		s.websocketHub.SendToConn(c.ID, room, messages)
	}
}
