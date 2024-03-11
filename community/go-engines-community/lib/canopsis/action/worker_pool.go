package action

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libwebhook "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/webhook"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TaskNew = iota
	TaskNotMatched
	TaskCancelled
	TaskRpcError
)

type Task struct {
	Source               string
	Action               Action
	Alarm                types.Alarm
	Entity               types.Entity
	Step                 int
	ExecutionCacheKey    string
	ExecutionID          string
	ScenarioID           string
	ScenarioName         string
	SkipForChild         bool
	IsMetaAlarmUpdated   bool
	SkipForInstruction   bool
	IsInstructionMatched bool
	Header               map[string]string
	Response             map[string]interface{}
	ResponseMap          map[string]interface{}
	AdditionalData       AdditionalData
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
	size                int
	axeRpcClient        engine.RPCClient
	webhookRpcClient    engine.RPCClient
	encoder             encoding.Encoder
	logger              zerolog.Logger
	templateExecutor    template.Executor
	alarmConfigProvider config.AlarmConfigProvider

	alarmCollection          mongo.DbCollection
	webhookHistoryCollection mongo.DbCollection
}

func NewWorkerPool(
	size int,
	dbClient mongo.DbClient,
	axeRpcClient engine.RPCClient,
	webhookRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
	templateExecutor template.Executor,
	alarmConfigProvider config.AlarmConfigProvider,
) WorkerPool {
	return &pool{
		size:                size,
		axeRpcClient:        axeRpcClient,
		webhookRpcClient:    webhookRpcClient,
		encoder:             encoder,
		logger:              logger,
		templateExecutor:    templateExecutor,
		alarmConfigProvider: alarmConfigProvider,

		alarmCollection:          dbClient.Collection(mongo.AlarmMongoCollection),
		webhookHistoryCollection: dbClient.Collection(mongo.WebhookHistoryMongoCollection),
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
							skip, err := s.call(ctx, task, id)
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
							if skip {
								resultChannel <- TaskResult{
									Source:            source,
									Alarm:             task.Alarm,
									Step:              task.Step,
									ExecutionCacheKey: task.ExecutionCacheKey,
									Status:            TaskNotMatched,
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

func (s *pool) call(ctx context.Context, task Task, workerId int) (bool, error) {
	var event interface{}
	var rpcClient engine.RPCClient
	var skip bool
	var err error
	switch task.Action.Type {
	case types.ActionTypeWebhook:
		rpcClient = s.webhookRpcClient
		event, skip, err = s.getRPCWebhookEvent(ctx, task)
	default:
		rpcClient = s.axeRpcClient
		event, err = s.getRPCAxeEvent(task)
	}

	if rpcClient == nil {
		return false, fmt.Errorf("cannot process action %s", task.Action.Type)
	}

	if err != nil || skip || event == nil {
		return skip, err
	}

	body, err := s.encoder.Encode(event)
	if err != nil {
		s.logger.Warn().Err(err).Msgf("Worker %d encode rpc for action '%s' failed", workerId, task.Action.Type)
		return false, err
	}

	err = rpcClient.Call(ctx, engine.RPCMessage{
		CorrelationID: fmt.Sprintf("%s&&%d", task.ExecutionCacheKey, task.Step),
		Body:          body,
	})
	if err != nil {
		s.logger.Warn().Err(err).Msgf("Worker %d send rpc for action '%s' failed", workerId, task.Action.Type)
		return false, err
	}

	return false, nil
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
		return nil, fmt.Errorf("cannot render output template scenario=%s: %w", task.ScenarioID, err)
	}

	additionalData, err := s.resolveAuthor(task)
	if err != nil {
		return nil, err
	}

	axeParams := rpc.AxeParameters{
		Author:         additionalData.Author,
		User:           additionalData.User,
		State:          params.State,
		Duration:       params.Duration,
		Name:           params.Name,
		Reason:         params.Reason,
		Type:           params.Type,
		RRule:          params.RRule,
		Tstart:         params.Tstart,
		Tstop:          params.Tstop,
		StartOnTrigger: params.StartOnTrigger,
	}

	if task.Action.Type == types.ActionTypeAssocTicket {
		axeParams.TicketInfo = types.TicketInfo{
			Ticket:           params.Ticket,
			TicketURL:        params.TicketURL,
			TicketSystemName: params.TicketSystemName,
			TicketRuleName:   types.RuleNameScenarioPrefix + task.ScenarioName,
			TicketRuleID:     task.ScenarioID,
			TicketData:       params.TicketData,
			TicketComment:    task.Action.Comment,
		}
		axeParams.Output = axeParams.TicketInfo.GetStepMessage()
	} else if params.Output != "" {
		outputBuilder := strings.Builder{}
		outputBuilder.WriteString(types.RuleNameScenarioPrefix)
		outputBuilder.WriteString(task.ScenarioName)
		alarmConfig := s.alarmConfigProvider.Get()
		outputBuilder.WriteString(". ")
		outputBuilder.WriteString(types.OutputCommentPrefix)
		outputBuilder.WriteString(utils.TruncateString(params.Output, alarmConfig.OutputLength))
		outputBuilder.WriteRune('.')
		axeParams.Output = outputBuilder.String()
	}

	return &rpc.AxeEvent{
		EventType:  task.Action.Type,
		Parameters: axeParams,
		Alarm:      &task.Alarm,
		Entity:     &task.Entity,
	}, nil
}

func (s *pool) getRPCWebhookEvent(ctx context.Context, task Task) (*rpc.WebhookEvent, bool, error) {
	children := make([]types.Alarm, 0)
	if len(task.Alarm.Value.Children) > 0 {
		cursor, err := s.alarmCollection.Find(ctx, bson.M{
			"d":          bson.M{"$in": task.Alarm.Value.Children},
			"v.resolved": nil,
		})
		if err != nil {
			return nil, false, fmt.Errorf("cannot find children: %w", err)
		}
		err = cursor.All(ctx, &children)
		if err != nil {
			return nil, false, fmt.Errorf("cannot decode children: %w", err)
		}
	}
	// Skip if instruction is in progress
	if task.SkipForInstruction && task.IsInstructionMatched {
		return nil, true, nil
	}
	// Skip webhooks for children
	if task.SkipForChild {
		if task.IsMetaAlarmUpdated {
			return nil, true, nil
		} else if len(task.Alarm.Value.Parents) > 0 {
			err := s.alarmCollection.FindOne(ctx, bson.M{
				"d":          bson.M{"$in": task.Alarm.Value.Parents},
				"v.resolved": nil,
			}, options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
			if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil, false, fmt.Errorf("cannot find parents: %w", err)
			}
			if err == nil {
				return nil, true, nil
			}
		}
	}

	additionalData, err := s.resolveAuthor(task)
	if err != nil {
		return nil, false, err
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
		return nil, false, fmt.Errorf("cannot render request url template scenario=%s: %w", task.ScenarioID, err)
	}
	request.Payload, err = s.templateExecutor.Execute(request.Payload, tplData)
	if err != nil {
		return nil, false, fmt.Errorf("cannot render request payload template scenario=%s: %w", task.ScenarioID, err)
	}

	headers := make(map[string]string, len(request.Headers))
	for k, v := range request.Headers {
		headers[k], err = s.templateExecutor.Execute(v, tplData)
		if err != nil {
			return nil, false, fmt.Errorf("cannot render request header %q template scenario=%s: %w", k, task.ScenarioID, err)
		}
	}
	request.Headers = headers

	history := libwebhook.History{
		ID:        utils.NewID(),
		Alarms:    []string{task.Alarm.ID},
		Scenario:  task.ScenarioID,
		Index:     int64(task.Step),
		Execution: task.ExecutionID,
		Name:      types.RuleNameScenarioPrefix + task.ScenarioName,

		SystemName:    task.Action.Parameters.TicketSystemName,
		Status:        libwebhook.StatusCreated,
		Comment:       task.Action.Comment,
		Request:       request,
		DeclareTicket: task.Action.Parameters.DeclareTicket,
		UserID:        additionalData.User,
		Username:      additionalData.Author,
		Initiator:     types.InitiatorSystem,
		CreatedAt:     datetime.NewCpsTime(),
	}

	err = s.webhookHistoryCollection.FindOneAndUpdate(ctx,
		bson.M{
			"execution": history.Execution,
			"scenario":  history.Scenario,
			"index":     history.Index,
		},
		bson.M{
			"$setOnInsert": history,
		},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&history)
	if err != nil {
		return nil, false, fmt.Errorf("cannot save webhook history scenario=%s: %w", task.ScenarioID, err)
	}

	return &rpc.WebhookEvent{
		Execution: history.ID,
	}, false, nil
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
		return additionalData, fmt.Errorf("cannot render author template scenario=%s: %w", task.ScenarioID, err)
	}

	return additionalData, nil
}
