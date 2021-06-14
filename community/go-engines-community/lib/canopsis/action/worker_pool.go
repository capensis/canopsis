package action

import (
	"context"
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/rs/zerolog"
	"sync"
)

const (
	TaskNew = iota
	TaskNotMatched
	TaskCancelled
	TaskRpcError
)

type Task struct {
	Source       string
	Action       Action
	Alarm        types.Alarm
	Entity       types.Entity
	Step         int
	ExecutionID  string
	ScenarioID   string
	AckResources bool
	Header       map[string]string
	Response     map[string]interface{}
}

type TaskResult struct {
	Source          string
	Alarm           types.Alarm
	Step            int
	ExecutionID     string
	AlarmChangeType types.AlarmChangeType
	Status          int
	Header          map[string]string
	Response        map[string]interface{}
	Err             error
}

type WorkerPool interface {
	RunWorkers(ctx context.Context, taskChannel <-chan Task) (<-chan TaskResult, error)
}

type pool struct {
	size             int
	axeRpcClient     engine.RPCClient
	webhookRpcClient engine.RPCClient
	alarmAdapter     alarm.Adapter
	encoder          encoding.Encoder
	logger           zerolog.Logger
}

func NewWorkerPool(
	size int,
	axeRpcClient engine.RPCClient,
	webhookRpcClient engine.RPCClient,
	alarmAdapter alarm.Adapter,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) WorkerPool {
	return &pool{
		size:             size,
		axeRpcClient:     axeRpcClient,
		webhookRpcClient: webhookRpcClient,
		alarmAdapter:     alarmAdapter,
		encoder:          encoder,
		logger:           logger,
	}
}

func (s *pool) RunWorkers(ctx context.Context, taskChannel <-chan Task) (<-chan TaskResult, error) {
	s.logger.Info().Msg("Worker pool started")

	resultChannel := make(chan TaskResult)

	if s.size < 1 {
		return nil, errors.New("action worker pool error: size is less than 1")
	}

	go func() {
		wg := sync.WaitGroup{}
		defer close(resultChannel)

		for i := 0; i < s.size; i++ {
			wg.Add(1)
			go func(id int) {
				s.logger.Info().Msgf("Worker %d started", id)
				defer wg.Done()

				workerCtx, cancel := context.WithCancel(ctx)
				defer cancel()

				source := fmt.Sprintf("worker %d", id)

				for {
					select {
					case <-workerCtx.Done():
						s.logger.Debug().Msgf("Worker %d cancelled", id)

						resultChannel <- TaskResult{
							Source: source,
							Status: TaskCancelled,
							Err:    workerCtx.Err(),
						}

						return
					case task, ok := <-taskChannel:
						if !ok {
							return
						}

						s.logger.Debug().Msgf("Worker %d got task - %v", id, task)

						if !task.Action.AlarmPatterns.Matches(&task.Alarm) || !task.Action.EntityPatterns.Matches(&task.Entity) {
							resultChannel <- TaskResult{
								Source:      source,
								Alarm:       task.Alarm,
								Step:        task.Step,
								ExecutionID: task.ExecutionID,
								Status:      TaskNotMatched,
							}
						} else {
							err := s.call(task, id)
							if err != nil {
								resultChannel <- TaskResult{
									Source:      source,
									Alarm:       task.Alarm,
									Step:        task.Step,
									ExecutionID: task.ExecutionID,
									Status:      TaskRpcError,
									Err:         err,
								}
								break
							}

							s.logger.Debug().Msgf("Worker %d send rpc for action '%s'", id, task.Action.Type)
						}

						s.logger.Debug().Msgf("Worker %d finished task - %v", id, task)
					}
				}
			}(i)
		}

		wg.Wait()
	}()

	return resultChannel, nil
}

func (s *pool) call(task Task, workerId int) error {
	var event interface{}
	var rpcClient engine.RPCClient
	var err error
	switch task.Action.Type {
	case types.ActionTypeWebhook:
		rpcClient = s.webhookRpcClient
		event, err = s.getRPCWebhookEvent(task)
	default:
		rpcClient = s.axeRpcClient
		event, err = s.getRPCAxeEvent(task)
	}

	if rpcClient == nil {
		return fmt.Errorf("cannot process action %s", task.Action.Type)
	}

	if err != nil {
		return err
	}

	body, err := s.encoder.Encode(event)
	if err != nil {
		s.logger.Warn().Err(err).Msgf("Worker %d encode rpc for action '%s' failed", workerId, task.Action.Type)
		return err
	}

	err = rpcClient.Call(engine.RPCMessage{
		CorrelationID: fmt.Sprintf("%s&&%d", task.ExecutionID, task.Step),
		Body:          body,
	})
	if err != nil {
		s.logger.Warn().Err(err).Msgf("Worker %d send rpc for action '%s' failed", workerId, task.Action.Type)
		return err
	}

	return nil
}

func (s *pool) getRPCAxeEvent(task Task) (*types.RPCAxeEvent, error) {
	if t, ok := task.Action.Parameters.(types.Templater); ok {
		err := t.Template(types.AlarmWithEntity{
			Alarm:  task.Alarm,
			Entity: task.Entity,
		})

		if err != nil {
			return nil, err
		}
	}

	return &types.RPCAxeEvent{
		EventType:  task.Action.Type,
		Parameters: task.Action.Parameters,
		Alarm:      &task.Alarm,
		Entity:     &task.Entity,
	}, nil
}

func (s *pool) getRPCWebhookEvent(task Task) (*types.RPCWebhookEvent, error) {
	params, ok := task.Action.Parameters.(*types.WebhookParameters)
	if !ok {
		return nil, errors.New("invalid parameters")
	}

	children := make([]types.Alarm, 0)
	if len(task.Alarm.Value.Children) > 0 {
		err := s.alarmAdapter.GetOpenedAlarmsByAlarmIDs(task.Alarm.Value.Children, &children)
		if err != nil {
			return nil, fmt.Errorf("cannot find children : %v", err)
		}
	}

	err := params.Template(map[string]interface{}{
		"Alarm":    task.Alarm,
		"Entity":   task.Entity,
		"Children": children,
		"Response": task.Response,
		"Header":   task.Header,
	})
	if err != nil {
		return nil, err
	}

	return &types.RPCWebhookEvent{
		Parameters:   *params,
		Alarm:        &task.Alarm,
		AckResources: task.AckResources,
		Header:       task.Header,
		Response:     task.Response,
		Message:      fmt.Sprintf("step %d of scenario %s", task.Step, task.ScenarioID),
	}, nil
}
