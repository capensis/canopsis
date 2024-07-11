package messageratestats

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"github.com/rs/zerolog"
)

type Watcher interface {
	StartWatch(ctx context.Context, connId, userID, roomId string, data any) error
	StopWatch(connId, roomId string) error
}

func NewWatcher(
	hub websocket.Hub,
	store Store,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
	tickDuration time.Duration,
	logger zerolog.Logger,
) Watcher {
	return &watcher{
		hub:          hub,
		store:        store,
		encoder:      encoder,
		decoder:      decoder,
		tickDuration: tickDuration,
		logger:       logger,
		streams:      make(map[string]streamData),
	}
}

type watcher struct {
	hub     websocket.Hub
	store   Store
	encoder encoding.Encoder
	decoder encoding.Decoder
	logger  zerolog.Logger

	tickDuration time.Duration

	streamsMx sync.RWMutex
	streams   map[string]streamData
}

type streamData struct {
	connIds []string
	cancel  context.CancelFunc
}

// StartWatch creates a new stream change or adds a connection to an existed one if there is already a stream change with the same request.
func (w *watcher) StartWatch(ctx context.Context, connId, _, _ string, data any) error {
	b, err := w.encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("unexpected data type: %w", err)
	}

	k := w.genKey(b)
	ctx, cancel := context.WithCancel(ctx)
	if !w.newStream(k, connId, cancel) {
		return nil
	}

	var searchRequest SearchRequest
	err = w.decoder.Decode(b, &searchRequest)
	if err != nil {
		return fmt.Errorf("unexpected data type: %w", err)
	}

	go func() {
		defer cancel()

		ticker := time.NewTicker(w.tickDuration)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				rates, err := w.store.Find(ctx, ListRequest{
					Interval:      IntervalMinute,
					SearchRequest: searchRequest,
				})
				if err != nil {
					w.logger.Err(err).Msg("cannot get minute message rates")
					continue
				}

				w.hub.SendGroupRoomByConnections(w.getConnIds(k), websocket.RoomMessageRates, "", rates)
			}
		}
	}()

	return nil
}

func (w *watcher) StopWatch(connId, _ string) error {
	w.streamsMx.Lock()
	defer w.streamsMx.Unlock()

	for k, v := range w.streams {
		index := slices.Index(v.connIds, connId)

		if index < 0 {
			continue
		}

		v.connIds = slices.Delete(v.connIds, index, index+1)
		if len(v.connIds) == 0 {
			delete(w.streams, k)
			v.cancel()
		} else {
			w.streams[k] = v
		}

		return nil
	}

	return nil
}

func (w *watcher) newStream(k, connId string, streamCancel context.CancelFunc) bool {
	w.streamsMx.Lock()
	defer w.streamsMx.Unlock()

	d, ok := w.streams[k]
	if !ok {
		w.streams[k] = streamData{
			connIds: []string{connId},
			cancel:  streamCancel,
		}

		return true
	}

	d.connIds = append(d.connIds, connId)
	w.streams[k] = d

	return false
}

func (w *watcher) getConnIds(k string) []string {
	w.streamsMx.RLock()
	defer w.streamsMx.RUnlock()

	return w.streams[k].connIds
}

func (w *watcher) genKey(b []byte) string {
	cacheKey := sha256.Sum256(b)
	return hex.EncodeToString(cacheKey[:])
}
