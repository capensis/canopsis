package alarm

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Watcher interface {
	StartWatch(ctx context.Context, connId, userID, roomId string, data any) error
	StartWatchDetails(ctx context.Context, connId, userID, roomId string, data any) error
	StopWatch(connId, roomId string) error
}

func NewWatcher(
	client mongo.DbClient,
	hub websocket.Hub,
	store Store,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
	logger zerolog.Logger,
) Watcher {
	return &watcher{
		collection: client.Collection(mongo.AlarmMongoCollection),
		hub:        hub,
		store:      store,
		encoder:    encoder,
		decoder:    decoder,
		logger:     logger,
		streams:    make(map[string]map[string]streamData),
	}
}

// watcher subscribes to MongoDb to watch after changes of alarms.
type watcher struct {
	collection mongo.DbCollection
	hub        websocket.Hub
	store      Store
	encoder    encoding.Encoder
	decoder    encoding.Decoder
	logger     zerolog.Logger

	streamsMx sync.RWMutex
	streams   map[string]map[string]streamData
}

type streamData struct {
	connIdsByUserId map[string][]string
	cancel          context.CancelFunc
}

// StartWatch creates a new stream change or adds a connection to an existed one if there is already a stream change with the same request.
func (w *watcher) StartWatch(ctx context.Context, connId, userID, roomId string, data any) error {
	b, err := w.encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("unexpected data type: %w", err)
	}

	k := w.genKey(b)
	streamCtx, streamCancel := context.WithCancel(ctx)
	if !w.newStream(roomId, k, connId, userID, streamCancel) {
		return nil
	}

	var alarmIds []string
	err = w.decoder.Decode(b, &alarmIds)
	if err != nil {
		return fmt.Errorf("unexpected data type: %w", err)
	}

	stream, err := w.collection.Watch(ctx, []bson.M{
		{"$match": bson.M{
			"operationType":   "update",
			"documentKey._id": bson.M{"$in": alarmIds},
		}},
	})
	if err != nil {
		return fmt.Errorf("cannot watch collection: %w", err)
	}

	go func() {
		defer func() {
			_ = stream.Close(streamCtx)
			streamCancel()
		}()

		for stream.Next(streamCtx) {
			changeEvent := struct {
				DocumentKey struct {
					ID string `bson:"_id"`
				} `bson:"documentKey"`
			}{}
			err = stream.Decode(&changeEvent)
			if err != nil {
				w.logger.Err(err).Msgf("cannot decode alarm")
				continue
			}

			connIdsByUserId := w.getConnIds(roomId, k)
			for userID, connIds := range connIdsByUserId {
				res, err := w.store.GetByID(streamCtx, changeEvent.DocumentKey.ID, userID, true)
				if err != nil {
					w.logger.Err(err).Msgf("cannot get alarm")
					continue
				}
				if res == nil {
					w.logger.Error().Msgf("cannot find alarm")
					continue
				}

				w.hub.SendGroupRoomByConnections(connIds, websocket.RoomAlarmsGroup, roomId, res)
			}
		}
	}()

	return nil
}

// StartWatchDetails creates a new stream change or adds a connection to an existed one if there is already a stream change with the same request.
func (w *watcher) StartWatchDetails(ctx context.Context, connId, userID, roomId string, data any) error {
	b, err := w.encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("unexpected data type: %w", err)
	}

	k := w.genKey(b)
	streamCtx, streamCancel := context.WithCancel(ctx)
	if !w.newStream(roomId, k, connId, userID, streamCancel) {
		return nil
	}

	var requests []DetailsRequest
	err = w.decoder.Decode(b, &requests)
	if err != nil {
		return fmt.Errorf("unexpected data type: %w", err)
	}

	requestsById := make(map[string]DetailsRequest, len(requests))
	alarmIds := make([]string, len(requests))
	metaAlarmIds := make([]string, 0, len(requests))
	for i, request := range requests {
		request.Format()
		err = binding.Validator.ValidateStruct(request)
		if err != nil {
			return fmt.Errorf("invalid request %d: %w", i, err)
		}

		requestsById[request.ID] = request
		alarmIds[i] = request.ID
		if request.Children != nil && request.Children.Page > 0 {
			metaAlarmIds = append(metaAlarmIds, request.ID)
		}
	}

	metaAlarmEntityIds := make([]string, 0, len(metaAlarmIds))
	metaAlarmIdByEntityId := make(map[string]string, len(metaAlarmIds))
	if len(metaAlarmIds) > 0 {
		cursor, err := w.collection.Find(ctx, bson.M{
			"_id":        bson.M{"$in": metaAlarmIds},
			"v.resolved": nil,
		}, options.Find().SetProjection(bson.M{"d": 1, "v.meta": 1}))
		if err != nil {
			return fmt.Errorf("cannot find alarm: %w", err)
		}

		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			metaAlarm := types.Alarm{}
			err := cursor.Decode(&metaAlarm)
			if err != nil {
				return err
			}

			if metaAlarm.Value.Meta != "" {
				metaAlarmEntityIds = append(metaAlarmEntityIds, metaAlarm.EntityID)
				metaAlarmIdByEntityId[metaAlarm.EntityID] = metaAlarm.ID
			}
		}
	}

	var pipeline []bson.M
	opts := options.ChangeStream()
	if len(metaAlarmEntityIds) == 0 {
		pipeline = []bson.M{
			{"$match": bson.M{
				"operationType":   "update",
				"documentKey._id": bson.M{"$in": alarmIds},
			}},
		}
	} else {
		pipeline = []bson.M{
			{"$match": bson.M{"$or": []bson.M{
				{
					"operationType":   "update",
					"documentKey._id": bson.M{"$in": alarmIds},
				},
				{
					"operationType":          "update",
					"fullDocument.v.parents": bson.M{"$in": metaAlarmEntityIds},
				},
			}}},
		}
		opts = opts.SetFullDocument(options.UpdateLookup)
	}

	stream, err := w.collection.Watch(ctx, pipeline, opts)
	if err != nil {
		return fmt.Errorf("cannot watch collection: %w", err)
	}

	go func() {
		defer func() {
			_ = stream.Close(streamCtx)
			streamCancel()
		}()
		for stream.Next(streamCtx) {
			changeEvent := struct {
				DocumentKey struct {
					ID string `bson:"_id"`
				} `bson:"documentKey"`
				FullDocument types.Alarm `bson:"fullDocument"`
			}{}
			err = stream.Decode(&changeEvent)
			if err != nil {
				w.logger.Err(err).Msgf("cannot decode alarm")
				continue
			}

			connIdsByUserId := w.getConnIds(roomId, k)
			w.sendGroupRoomAlrmDetails(streamCtx, changeEvent.DocumentKey.ID, roomId, requestsById, connIdsByUserId)

			for _, parent := range changeEvent.FullDocument.Value.Parents {
				if metaAlarmId, ok := metaAlarmIdByEntityId[parent]; ok {
					w.sendGroupRoomAlrmDetails(streamCtx, metaAlarmId, roomId, requestsById, connIdsByUserId)
				}
			}
		}
	}()

	return nil
}

func (w *watcher) sendGroupRoomAlrmDetails(ctx context.Context, alarmId, roomId string, requestsById map[string]DetailsRequest, connIdsByUserId map[string][]string) {
	request, ok := requestsById[alarmId]
	if !ok {
		return
	}
	for userID, connIds := range connIdsByUserId {
		res, err := w.store.GetDetails(ctx, request, userID)
		if err != nil {
			w.logger.Err(err).Msgf("cannot get alarm")
			continue
		}
		if res != nil {
			res.ID = request.ID
			w.hub.SendGroupRoomByConnections(connIds, websocket.RoomAlarmDetailsGroup, roomId, res)
		}
	}
}

func (w *watcher) StopWatch(connId, roomId string) error {
	w.streamsMx.Lock()
	defer w.streamsMx.Unlock()

	for k, v := range w.streams[roomId] {
		for userID, connIds := range v.connIdsByUserId {
			index := -1

			for i, streamConnId := range connIds {
				if streamConnId == connId {
					index = i
					break
				}
			}

			if index < 0 {
				continue
			}

			w.streams[roomId][k].connIdsByUserId[userID] = append(connIds[:index], connIds[index+1:]...)
			if len(w.streams[roomId][k].connIdsByUserId[userID]) == 0 {
				delete(w.streams[roomId][k].connIdsByUserId, userID)

				if len(w.streams[roomId][k].connIdsByUserId) == 0 {
					delete(w.streams[roomId], k)
					v.cancel()
				}
			}

			return nil
		}
	}

	return nil
}

func (w *watcher) newStream(roomId, k, connId, userID string, streamCancel context.CancelFunc) bool {
	w.streamsMx.Lock()
	defer w.streamsMx.Unlock()

	if _, ok := w.streams[roomId]; !ok {
		w.streams[roomId] = map[string]streamData{k: {
			connIdsByUserId: map[string][]string{userID: {connId}},
			cancel:          streamCancel,
		}}

		return true
	}

	if _, ok := w.streams[roomId][k]; ok {
		w.streams[roomId][k].connIdsByUserId[userID] = append(w.streams[roomId][k].connIdsByUserId[userID], connId)
		return false
	}

	w.streams[roomId][k] = streamData{
		connIdsByUserId: map[string][]string{userID: {connId}},
		cancel:          streamCancel,
	}

	return true
}

func (w *watcher) getConnIds(roomId, k string) map[string][]string {
	w.streamsMx.RLock()
	defer w.streamsMx.RUnlock()

	return w.streams[roomId][k].connIdsByUserId
}

func (w *watcher) genKey(b []byte) string {
	cacheKey := sha256.Sum256(b)
	return hex.EncodeToString(cacheKey[:])
}
