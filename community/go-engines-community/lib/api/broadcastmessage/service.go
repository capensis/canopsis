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

	messages, err := s.store.GetActive(ctx)
	if err != nil {
		s.logger.Err(err).Msg("cannot fetch messages")
	}
	previous := make(map[string]bool, len(messages))
	for _, v := range messages {
		previous[v.ID] = true
	}

	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-ch:
			if !ok {
				return
			}

			previous = s.check(ctx, previous, false)
		case <-ticker.C:
			previous = s.check(ctx, previous, true)
		}
	}
}

func (s *service) check(ctx context.Context, previous map[string]bool, onChange bool) map[string]bool {
	messages, err := s.store.GetActive(ctx)
	if err != nil {
		s.logger.Err(err).Msg("cannot fetch messages")
		return previous
	}

	if onChange {
		// Check if messages are changed.
		if len(previous) == len(messages) {
			equal := true
			for _, v := range messages {
				if !previous[v.ID] {
					equal = false
					break
				}
			}

			if equal {
				return previous
			}
		}
	}

	ids := make(map[string]bool, len(messages))
	for _, v := range messages {
		ids[v.ID] = true
	}
	s.websocketHub.Send(websocket.RoomBroadcastMessages, messages)

	return ids
}
