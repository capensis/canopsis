package action

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

const (
	TaskNew = iota
	TaskNotMatched
	TaskCancelled
	TaskRpcError
)

type Task struct {
	Source            string
	Action            Action
	Alarm             types.Alarm
	Entity            types.Entity
	Step              int
	ExecutionCacheKey string
	ScenarioID        string
	AckResources      bool
	Header            map[string]string
	Response          map[string]interface{}
	ResponseMap       map[string]interface{}
	AdditionalData    AdditionalData
}

type TaskResult struct {
	Source            string
	Alarm             types.Alarm
	Step              int
	ExecutionCacheKey string
	AlarmChangeType   types.AlarmChangeType
	Status            int
	Header            map[string]string
	Response          map[string]interface{}
	Err               error
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
	templateExecutor *template.Executor
}

func NewWorkerPool(
	size int,
	axeRpcClient engine.RPCClient,
	webhookRpcClient engine.RPCClient,
	alarmAdapter alarm.Adapter,
	encoder encoding.Encoder,
	logger zerolog.Logger,
	timezoneConfigProvider config.TimezoneConfigProvider,
) WorkerPool {
	return &pool{
		size:             size,
		axeRpcClient:     axeRpcClient,
		webhookRpcClient: webhookRpcClient,
		alarmAdapter:     alarmAdapter,
		encoder:          encoder,
		logger:           logger,
		templateExecutor: template.NewExecutor(timezoneConfigProvider),
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

						s.logger.Debug().Interface("task", task).Msgf("Worker %d got task", id)

						match, err := task.Action.Match(task.Entity, task.Alarm)
						if err != nil {
							s.logger.Err(err).Msgf("match action %d from scenario %s returned error", task.Step, task.ScenarioID)
						}

						if !match {
							resultChannel <- TaskResult{
								Source:            source,
								Alarm:             task.Alarm,
								Step:              task.Step,
								ExecutionCacheKey: task.ExecutionCacheKey,
								Status:            TaskNotMatched,
							}
						} else {
							err := s.call(ctx, task, id)
							if err != nil {
								resultChannel <- TaskResult{
									Source:            source,
									Alarm:             task.Alarm,
									Step:              task.Step,
									ExecutionCacheKey: task.ExecutionCacheKey,
									Status:            TaskRpcError,
									Err:               err,
								}

								break
							}

							s.logger.Debug().Msgf("Worker %d send rpc for action '%s'", id, task.Action.Type)
						}

						s.logger.Debug().Interface("task", task).Msgf("Worker %d finished task", id)
					}
				}
			}(i)
		}

		wg.Wait()
	}()

	return resultChannel, nil
}

func (s *pool) call(ctx context.Context, task Task, workerId int) error {
	var event interface{}
	var rpcClient engine.RPCClient
	var err error
	switch task.Action.Type {
	case types.ActionTypeWebhook:
		rpcClient = s.webhookRpcClient
		event, err = s.getRPCWebhookEvent(ctx, task)
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

	err = rpcClient.Call(ctx, engine.RPCMessage{
		CorrelationID: fmt.Sprintf("%s&&%d", task.ExecutionCacheKey, task.Step),
		Body:          body,
	})
	if err != nil {
		s.logger.Warn().Err(err).Msgf("Worker %d send rpc for action '%s' failed", workerId, task.Action.Type)
		return err
	}

	return nil
}

func (s *pool) getRPCAxeEvent(task Task) (*rpc.AxeEvent, error) {
	params := task.Action.Parameters
	tplData := types.AlarmWithEntity{
		Alarm:  task.Alarm,
		Entity: task.Entity,
	}
	var err error
	params.Output, err = s.templateExecutor.Execute(params.Output, tplData)
	if err != nil {
		return nil, fmt.Errorf("cannot render template scenario=%s: %w", task.ScenarioID, err)
	}

	additionalData, err := s.resolveAuthor(task)
	if err != nil {
		return nil, err
	}

	axeParams := rpc.AxeParameters{
		Output: params.Output,
		Author: additionalData.Author,
		User:   additionalData.User,
		State:  params.State,
		TicketInfo: types.TicketInfo{
			Ticket:           params.Ticket,
			TicketURL:        params.TicketURL,
			TicketSystemName: params.TicketSystemName,
			TicketData:       params.TicketData,
		},
		Duration:       params.Duration,
		Name:           params.Name,
		Reason:         params.Reason,
		Type:           params.Type,
		RRule:          params.RRule,
		Tstart:         params.Tstart,
		Tstop:          params.Tstop,
		StartOnTrigger: params.StartOnTrigger,
	}

	return &rpc.AxeEvent{
		EventType:  task.Action.Type,
		Parameters: axeParams,
		Alarm:      &task.Alarm,
		Entity:     &task.Entity,
	}, nil
}

func (s *pool) getRPCWebhookEvent(ctx context.Context, task Task) (*rpc.WebhookEvent, error) {
	children := make([]types.Alarm, 0)
	if len(task.Alarm.Value.Children) > 0 {
		err := s.alarmAdapter.GetOpenedAlarmsByIDs(ctx, task.Alarm.Value.Children, &children)
		if err != nil {
			return nil, fmt.Errorf("cannot find children : %v", err)
		}
	}

	additionalData, err := s.resolveAuthor(task)
	if err != nil {
		return nil, err
	}

	tplData := map[string]interface{}{
		"Alarm":          task.Alarm,
		"Entity":         task.Entity,
		"Children":       children,
		"Response":       task.Response,
		"ResponseMap":    task.ResponseMap,
		"Header":         task.Header,
		"AdditionalData": additionalData,
	}

	request := *task.Action.Parameters.Request
	request.URL, err = s.templateExecutor.Execute(request.URL, tplData)
	if err != nil {
		return nil, fmt.Errorf("cannot render template scenario=%s: %w", task.ScenarioID, err)
	}
	request.Payload, err = s.templateExecutor.Execute(request.Payload, tplData)
	if err != nil {
		return nil, fmt.Errorf("cannot render template scenario=%s: %w", task.ScenarioID, err)
	}

	headers := make(map[string]string, len(request.Headers))
	for k, v := range request.Headers {
		headers[k], err = s.templateExecutor.Execute(v, tplData)
		if err != nil {
			return nil, fmt.Errorf("cannot render template scenario=%s: %w", task.ScenarioID, err)
		}
	}
	request.Headers = headers

	webhookParams := rpc.WebhookParameters{
		Request:       request,
		DeclareTicket: task.Action.Parameters.DeclareTicket,
		Scenario:      task.ScenarioID,
		Author:        additionalData.Author,
		User:          additionalData.User,
	}

	return &rpc.WebhookEvent{
		Parameters:   webhookParams,
		Alarm:        &task.Alarm,
		Entity:       &task.Entity,
		AckResources: task.AckResources,
		Header:       task.Header,
		Response:     task.Response,
		Message:      fmt.Sprintf("step %d of scenario %s", task.Step, task.ScenarioID),
	}, nil
}

func (s *pool) resolveAuthor(task Task) (AdditionalData, error) {
	additionalData := task.AdditionalData
	if task.Action.Parameters.ForwardAuthor != nil && *task.Action.Parameters.ForwardAuthor {
		return additionalData, nil
	}

	additionalData.User = ""

	if task.Action.Parameters.Author == "" {
		additionalData.Author = canopsis.DefaultEventAuthor
		return additionalData, nil
	}

	var err error
	additionalData.Author, err = s.templateExecutor.Execute(task.Action.Parameters.Author, types.AlarmWithEntity{
		Alarm:  task.Alarm,
		Entity: task.Entity,
	})
	if err != nil {
		return additionalData, fmt.Errorf("cannot render template scenario=%s: %w", task.ScenarioID, err)
	}

	return additionalData, nil
}
