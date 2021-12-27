package action

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"text/template"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
)

const (
	TaskNew = iota
	TaskNotMatched
	TaskCancelled
	TaskRpcError
)

type Task struct {
	Source         string
	Action         Action
	Alarm          types.Alarm
	Entity         types.Entity
	Step           int
	ExecutionID    string
	ScenarioID     string
	AckResources   bool
	Header         map[string]string
	Response       map[string]interface{}
	AdditionalData AdditionalData
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
	size                   int
	axeRpcClient           engine.RPCClient
	webhookRpcClient       engine.RPCClient
	alarmAdapter           alarm.Adapter
	encoder                encoding.Encoder
	logger                 zerolog.Logger
	timezoneConfigProvider *config.BaseTimezoneConfigProvider
}

func NewWorkerPool(
	size int,
	axeRpcClient engine.RPCClient,
	webhookRpcClient engine.RPCClient,
	alarmAdapter alarm.Adapter,
	encoder encoding.Encoder,
	logger zerolog.Logger,
	timezoneConfigProvider *config.BaseTimezoneConfigProvider,
) WorkerPool {
	return &pool{
		size:                   size,
		axeRpcClient:           axeRpcClient,
		webhookRpcClient:       webhookRpcClient,
		alarmAdapter:           alarmAdapter,
		encoder:                encoder,
		logger:                 logger,
		timezoneConfigProvider: timezoneConfigProvider,
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

						if !task.Action.AlarmPatterns.Matches(&task.Alarm) || !task.Action.EntityPatterns.Matches(&task.Entity) {
							resultChannel <- TaskResult{
								Source:      source,
								Alarm:       task.Alarm,
								Step:        task.Step,
								ExecutionID: task.ExecutionID,
								Status:      TaskNotMatched,
							}
						} else {
							err := s.call(ctx, task, id)
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
	params := make(map[string]interface{}, len(task.Action.Parameters))
	tplData := types.AlarmWithEntity{
		Alarm:  task.Alarm,
		Entity: task.Entity,
	}

	for k, v := range task.Action.Parameters {
		params[k] = v

		switch k {
		case "author", "output":
			if str, ok := v.(string); ok {
				var err error
				params[k], err = s.renderTemplate(str, tplData)
				if err != nil {
					return nil, fmt.Errorf("cannot render template scenario=%s: %w", task.ScenarioID, err)
				}
			}
		}
	}

	return &types.RPCAxeEvent{
		EventType:  task.Action.Type,
		Parameters: params,
		Alarm:      &task.Alarm,
		Entity:     &task.Entity,
	}, nil
}

func (s *pool) getRPCWebhookEvent(ctx context.Context, task Task) (*types.RPCWebhookEvent, error) {
	children := make([]types.Alarm, 0)
	if len(task.Alarm.Value.Children) > 0 {
		err := s.alarmAdapter.GetOpenedAlarmsByIDs(ctx, task.Alarm.Value.Children, &children)
		if err != nil {
			return nil, fmt.Errorf("cannot find children : %v", err)
		}
	}

	tplData := map[string]interface{}{
		"Alarm":          task.Alarm,
		"Entity":         task.Entity,
		"Children":       children,
		"Response":       task.Response,
		"Header":         task.Header,
		"AdditionalData": task.AdditionalData,
	}
	params := make(map[string]interface{}, len(task.Action.Parameters))
	for k, v := range task.Action.Parameters {
		params[k] = v

		switch k {
		case "request":
			if m, ok := v.(map[string]interface{}); ok {
				newRequest := make(map[string]interface{}, len(m))

				for requestKey, requestVal := range m {
					newRequest[requestKey] = requestVal

					switch requestKey {
					case "url", "payload":
						if str, ok := requestVal.(string); ok {
							var err error
							newRequest[requestKey], err = s.renderTemplate(str, tplData)
							if err != nil {
								return nil, fmt.Errorf("cannot render template scenario=%s: %w", task.ScenarioID, err)
							}
						}
					case "headers":
						if headers, ok := requestVal.(map[string]interface{}); ok {
							newHeaders := make(map[string]interface{}, len(headers))

							for headerKey, headerVal := range headers {
								newHeaders[headerKey] = headerVal

								if str, ok := headerVal.(string); ok {
									var err error
									newHeaders[headerKey], err = s.renderTemplate(str, tplData)
									if err != nil {
										return nil, fmt.Errorf("cannot render template scenario=%s: %w", task.ScenarioID, err)
									}
								}
							}

							newRequest[requestKey] = newHeaders
						}
					}
				}

				params[k] = newRequest
			}
		}
	}

	webhookParams := types.WebhookParameters{}
	err := mapstructure.Decode(params, &webhookParams)
	if err != nil {
		return nil, fmt.Errorf("cannot decode map struct scenario=%s : %v", task.ScenarioID, err)
	}

	return &types.RPCWebhookEvent{
		Parameters:   webhookParams,
		Alarm:        &task.Alarm,
		Entity:       &task.Entity,
		AckResources: task.AckResources,
		Header:       task.Header,
		Response:     task.Response,
		Message:      fmt.Sprintf("step %d of scenario %s", task.Step, task.ScenarioID),
	}, nil
}

func (s *pool) renderTemplate(templateStr string, data interface{}) (string, error) {
	var f template.FuncMap
	if s.timezoneConfigProvider != nil {
		timezone := s.timezoneConfigProvider.Get()
		f = types.GetTemplateFunc(&timezone)
	} else {
		f = types.GetTemplateFunc(nil)
	}

	t, err := template.New("template").Funcs(f).Parse(templateStr)
	if err != nil {
		return "", err
	}
	b := strings.Builder{}
	err = t.Execute(&b, data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
