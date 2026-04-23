package alarm

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Watcher interface {
	StartWatch(ctx context.Context, opts websocket.JoinOptions) error
	StartWatchDetails(ctx context.Context, opts websocket.JoinOptions) error
	StopWatch(ctx context.Context, opts websocket.LeaveOptions) error
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
func (w *watcher) StartWatch(ctx context.Context, opts websocket.JoinOptions) error {
	b, err := w.encoder.Encode(opts.Payload)
	if err != nil {
		return fmt.Errorf("unexpected data type: %w", err)
	}

	k := w.genKey(b)
	streamCtx, streamCancel := context.WithCancel(ctx)
	roomID := opts.RoomID
	if !w.newStream(roomID, k, opts.ConnID, opts.UserID, streamCancel) {
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

			connIdsByUserId := w.getConnIds(roomID, k)
			for userID, connIds := range connIdsByUserId {
				res, err := w.store.GetByID(streamCtx, changeEvent.DocumentKey.ID, userID)
				if err != nil {
					w.logger.Err(err).Msgf("cannot get alarm")
					continue
				}
				if res == nil {
					w.logger.Error().Msgf("cannot find alarm")
					continue
				}

				w.hub.SendMessage(ctx, res, websocket.ToConnection(websocket.GroupRoom(websocket.RoomAlarmsGroup, roomID), connIds...))
			}
		}
	}()

	return nil
}

// StartWatchDetails creates a new stream change or adds a connection to an existed one if there is already a stream change with the same request.
func (w *watcher) StartWatchDetails(ctx context.Context, opts websocket.JoinOptions) error {
	b, err := w.encoder.Encode(opts.Payload)
	if err != nil {
		return fmt.Errorf("unexpected data type: %w", err)
	}

	k := w.genKey(b)
	streamCtx, streamCancel := context.WithCancel(ctx)
	roomID := opts.RoomID
	if !w.newStream(roomID, k, opts.ConnID, opts.UserID, streamCancel) {
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
	csOts := options.ChangeStream()
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
		csOts = csOts.SetFullDocument(options.UpdateLookup)
	}

	stream, err := w.collection.Watch(ctx, pipeline, csOts)
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

			connIdsByUserId := w.getConnIds(roomID, k)
			w.sendGroupRoomAlrmDetails(streamCtx, changeEvent.DocumentKey.ID, roomID, requestsById, connIdsByUserId)

			for _, parent := range changeEvent.FullDocument.Value.Parents {
				if metaAlarmId, ok := metaAlarmIdByEntityId[parent]; ok {
					w.sendGroupRoomAlrmDetails(streamCtx, metaAlarmId, roomID, requestsById, connIdsByUserId)
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
			w.hub.SendMessage(ctx, res, websocket.ToConnection(websocket.GroupRoom(websocket.RoomAlarmDetailsGroup, roomId), connIds...))
		}
	}
}

func (w *watcher) StopWatch(_ context.Context, opts websocket.LeaveOptions) error {
	w.streamsMx.Lock()
	defer w.streamsMx.Unlock()

	roomID := opts.RoomID
	connID := opts.ConnID
	for k, v := range w.streams[roomID] {
		for userID, connIds := range v.connIdsByUserId {
			index := slices.Index(connIds, connID)
			if index < 0 {
				continue
			}

			w.streams[roomID][k].connIdsByUserId[userID] = slices.Delete(connIds, index, index+1)
			if len(w.streams[roomID][k].connIdsByUserId[userID]) == 0 {
				delete(w.streams[roomID][k].connIdsByUserId, userID)

				if len(w.streams[roomID][k].connIdsByUserId) == 0 {
					delete(w.streams[roomID], k)
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
