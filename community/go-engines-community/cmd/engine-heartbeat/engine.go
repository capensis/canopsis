package main

import (
	"context"
	"errors"
	"runtime/trace"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	cps "git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/heartbeat"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/streadway/amqp"
)

type EngineHeartBeat struct {
	cps.DefaultEngine
	References        References
	HeartBeatItems    heartbeat.SafeHeartbeatItems
	regheartBeatItems map[string]int
	SendAlarmFunc     func(string, int, string) error
}

// AddHeartBeatItem avoid to register duplicate HeartBeatItem-s in the Engine
func (e *EngineHeartBeat) AddHeartBeatItem(li heartbeat.Item) error {
	_, exists := e.regheartBeatItems[li.ID()]
	if !exists {
		e.HeartBeatItems.AddItem(li)
	} else {
		return errors.New("this HeartBeatItem already exists in the set")
	}
	return nil
}

func (e *EngineHeartBeat) ConsumerChan() (<-chan amqp.Delivery, error) {
	_, err := e.Sub.QueueDeclare(
		cps.HeartBeatQueueName,
		true,  // durable
		false, // autodelete
		false, // exclusive
		false, // nowait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}
	err = e.Sub.QueueBind(
		cps.HeartBeatQueueName,
		"#",
		"canopsis.events",
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return e.Sub.Consume(
		cps.HeartBeatQueueName,    // queue
		cps.HeartBeatConsumerName, // consumer
		false,                     // auto-ack
		false,                     // exclusive
		false,                     // no-local
		false,                     // no-wait
		nil,                       // args
	)
}

func (e *EngineHeartBeat) WorkerProcess(msg amqp.Delivery) {
	ctx, task := trace.NewTask(context.Background(), "heartbeat.WorkerProcess")
	defer task.End()

	trace.Logf(ctx, "event_size", "%d", len(msg.Body))

	var ievent types.GenericEvent
	err := ievent.JSONUnmarshal(msg.Body)
	if err != nil {
		e.processWorkerError(err, msg)
		return
	}

	_, err = e.workHeartBeat(ievent)
	if err != nil {
		e.processWorkerError(err, msg)
		return
	}
	e.DefaultEngine.AckMessage(msg)
}

func (e *EngineHeartBeat) sendalarm(alarmid string, state int, output string) error {
	event := types.Event{
		Connector:     "heartbeat",
		ConnectorName: "heartbeat",
		EventType:     types.EventTypeCheck,
		Component:     "heartbeats",
		Resource:      alarmid,
		SourceType:    types.SourceTypeResource,
		State:         types.CpsNumber(state),
		Output:        output,
	}

	bevent, err := e.References.JSONEncoder.Encode(event)
	if err != nil {
		return err
	}

	return e.References.ChannelPub.Publish(
		cps.HeartBeatExchangeName,
		"#",
		false,
		false,
		amqp.Publishing{
			Body:        bevent,
			ContentType: "application/json",
		},
	)
}

func (e *EngineHeartBeat) PeriodicalProcess() {
	e.Logger().Info().Msg("periodical process")
	err := e.LoadHeartbeatItems()
	if err != nil {
		e.Logger().Warn().Err(err).Msg("error on load heartbeat items")
		e.AskStop(canopsis.ExitEngine)
	}
	for _, heartBeatItem := range e.HeartBeatItems.Value() {
		id := heartBeatItem.ID()
		if id == "" {
			continue
		}
		result := e.References.Redis.Exists(id)
		err = result.Err()
		if result.Val() == 0 {
			err := e.SendAlarmFunc(id, types.AlarmStateCritical, heartBeatItem.Output)
			if err == nil {
				e.References.Redis.Set("alarm:"+id, 0, 0)
			}
		} else {
			_ = e.SendAlarmFunc(id, types.AlarmStateOK, heartBeatItem.Output)
		}
	}

}

func (e *EngineHeartBeat) warmup(li heartbeat.Item) error {
	res := e.References.Redis.Exists(li.ID())
	if res.Err() == nil && res.Val() == 0 {
		setr := e.References.Redis.Set(li.ID(), 0, li.MaxDuration)
		return setr.Err()
	}

	return nil
}

// LoadHeartbeatItems (re)loads of an engine's heartbeat items.
// Drop a heartbeats loaded from database and reloads again.
// Keep a heartbeats added by another way besides from the database.
func (e *EngineHeartBeat) LoadHeartbeatItems() error {
	heartbeatItems := make([]heartbeat.Item, 0)
	for _, hi := range e.HeartBeatItems.Value() {
		if hi.Source != heartbeat.HeartbeatSourceMongo {
			heartbeatItems = append(heartbeatItems, hi)
		}
	}

	if len(heartbeatItems) > 0 {
		e.Logger().Info().Msgf("keeped %d previously added heartbeat item(s) besides the database", len(heartbeatItems))
	}

	heartbeats, err := e.References.Adapter.Get()
	if err != nil {
		return err
	}

	if len(heartbeats) > 0 {
		e.Logger().Info().Msgf("loaded %d heartbeat item(s) from the database", len(heartbeats))
	}

	for _, hi := range heartbeats {
		li, err := hi.ToHeartBeatItem()
		if err != nil {
			return err
		}
		heartbeatItems = append(heartbeatItems, li)
	}

	e.HeartBeatItems = heartbeat.NewSafeHeartbeatItems()
	e.regheartBeatItems = make(map[string]int, 0)

	for _, li := range heartbeatItems {
		err = e.AddHeartBeatItem(li)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *EngineHeartBeat) Initialize() error {
	if err := e.DefaultEngine.Initialize(); err != nil {
		return err
	}

	if err := e.LoadHeartbeatItems(); err != nil {
		return err
	}

	for _, li := range e.HeartBeatItems.Value() {
		e.warmup(li)
	}

	e.Logger().Info().Msgf("loaded %d rule(s)", e.HeartBeatItems.Len())
	return nil
}

func (e *EngineHeartBeat) resolveAlarm(id, output string) error {
	return e.SendAlarmFunc(id, types.AlarmStateOK, output)
}

func (e *EngineHeartBeat) workHeartBeat(ievent types.GenericEvent) ([]int, error) {
	itemCounts := make([]int, 0)
	for itemCount, lm := range e.HeartBeatItems.Value() {
		id, err := ievent.PartialID(lm)
		if err != nil || id != lm.ID() {
			continue
		}

		// Set expiration date to lm.MaxDuration so when the
		// beat processing works, it only checks if the key
		// exists or no.
		cmdstatus := e.References.Redis.Set(id, nil, lm.MaxDuration)
		if cmdstatus.Err() != nil {
			e.Logger().Warn().Err(cmdstatus.Err()).Str("id", id).Msg("redis set error")
			continue
		}

		// Check for an alarm that we must solve.
		alarmID := "alarm:" + id
		cmdint := e.References.Redis.Exists(alarmID)
		if cmdint.Err() != nil {
			e.Logger().Warn().Err(cmdstatus.Err()).Str("alarmID", alarmID).Msg("redis exists error")
			continue
		}

		if cmdint.Val() != 0 {
			e.Logger().Info().Str("alarmID", alarmID).Msg("resolution")
			cmddel := e.References.Redis.Del(alarmID)
			if cmddel.Err() != nil {
				e.Logger().Warn().Err(cmddel.Err()).Str("alarmID", alarmID).Msg("redis del error")
			}
			err := e.resolveAlarm(id, lm.Output)
			if err != nil {
				e.Logger().Warn().Err(err).Str("id", id).Msg("error resolving alarm")
			}
		}

		itemCounts = append(itemCounts, itemCount)
	}

	return itemCounts, nil
}

func (e *EngineHeartBeat) processWorkerError(err error, msg amqp.Delivery) {
	e.Logger().Error().Err(err).Msgf("event processing error : %+v", string(msg.Body))

	e.DefaultEngine.ProcessWorkerError(err, msg)
}

func (e *EngineHeartBeat) GetRunInfo() engine.RunInfo {
	return engine.RunInfo{
		Name:         "engine-heartbeat",
		ConsumeQueue: cps.HeartBeatQueueName,
		PublishQueue: "",
	}
}
