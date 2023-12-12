package ratelimit

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statistics"
	"github.com/rs/zerolog"
)

const writeChanTimeout = 100 * time.Millisecond

type StatsSender interface {
	Add(ts int64, pass bool)
}

func NewStatsSender(ch chan<- statistics.Message, logger zerolog.Logger) StatsSender {
	return &baseStatsSender{
		ch:     ch,
		logger: logger,
	}
}

type baseStatsSender struct {
	ch     chan<- statistics.Message
	logger zerolog.Logger
}

func (s *baseStatsSender) Add(ts int64, pass bool) {
	var dropped int64

	if !pass {
		dropped = 1
	}

	select {
	case s.ch <- statistics.Message{
		Timestamp: ts,
		Received:  1,
		Dropped:   dropped,
	}:
	case <-time.After(writeChanTimeout):
		s.logger.Debug().Msg("failed to write to channel by timeout")
	}

}
