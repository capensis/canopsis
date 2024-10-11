package action

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
)

type redisBasedManager struct {
	workerPool       WorkerPool
	taskChannel      chan Task
	outputChannel    chan ScenarioResult
	executionStorage ScenarioExecutionStorage
	scenarioStorage  ScenarioStorage
	logger           zerolog.Logger
}

func NewTaskManager(
	workerPool WorkerPool,
	executionStorage ScenarioExecutionStorage,
	scenarioStorage ScenarioStorage,
	logger zerolog.Logger,
) TaskManager {
	return &redisBasedManager{
		workerPool:       workerPool,
		executionStorage: executionStorage,
		scenarioStorage:  scenarioStorage,
		logger:           logger,
	}
}

func (e *redisBasedManager) Run(
	parentCtx context.Context,
	rpcResultChannel <-chan RpcResult,
	inputChannel <-chan ExecuteScenariosTask,
) (resCh <-chan ScenarioResult, resErr error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer func() {
		if resErr != nil {
			cancel()
		}
	}()

	e.logger.Debug().Msg("Task manager started")

	e.outputChannel = make(chan ScenarioResult)
	e.taskChannel = make(chan Task)

	taskResultChannel, err := e.workerPool.RunWorkers(ctx, e.taskChannel)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			cancel()
			close(e.outputChannel)
			close(e.taskChannel)
		}()

		wg := sync.WaitGroup{}
		wg.Add(3)

		go e.listenInputChannel(ctx, &wg, inputChannel)
		go e.listenRPCResultChannel(ctx, &wg, rpcResultChannel)
		go e.listenTaskResultChannel(ctx, &wg, taskResultChannel)

		wg.Wait()
	}()

	return e.outputChannel, nil
}

func (e *redisBasedManager) listenInputChannel(ctx context.Context, wg *sync.WaitGroup, inputChannel <-chan ExecuteScenariosTask) {
	e.logger.Debug().Msg("start listen scenario tasks")
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			e.logger.Debug().Msg("input listener cancelled")
			return
		case scenariosTask, ok := <-inputChannel:
			if !ok {
				e.logger.Debug().Msg("input channel closed")
				return
			}

			go func(ctx context.Context, task ExecuteScenariosTask) {
				if task.DelayedScenarioID != "" {
					scenario := e.scenarioStorage.GetScenario(task.DelayedScenarioID)
					if scenario == nil {
						e.logger.Error().Str("scenario", task.DelayedScenarioID).Msg("cannot find scenario")
						e.outputChannel <- ScenarioResult{
							Alarm:                task.Alarm,
							StartEventProcessing: task.Start,
							FifoAckEvent:         task.FifoAckEvent,
							Err:                  errors.New("scenario doesn't exist"),
						}
						return
					}
					_, err := e.executionStorage.Inc(ctx, task.Alarm.ID, 1, true)
					if err != nil {
						e.logger.Err(err).Str("scenario", task.DelayedScenarioID).Str("alarm", task.Alarm.ID).Msg("cannot run scenario")
						e.outputChannel <- ScenarioResult{
							Alarm:                task.Alarm,
							StartEventProcessing: task.Start,
							FifoAckEvent:         task.FifoAckEvent,
							Err:                  err,
						}
						return
					}

					e.startExecution(ctx, *scenario, task.Alarm, task.Entity, task.AdditionalData, task.FifoAckEvent,
						task.Start, task.IsMetaAlarmUpdated, task.IsInstructionMatched)
					return
				}

				if task.AbandonedExecutionCacheKey != "" {
					execution, err := e.executionStorage.Get(ctx, task.AbandonedExecutionCacheKey)
					if err != nil {
						e.logger.Err(err).Str("execution", task.AbandonedExecutionCacheKey).Msg("cannot find abandoned scenario")
						e.outputChannel <- ScenarioResult{
							Alarm:                task.Alarm,
							StartEventProcessing: task.Start,
							FifoAckEvent:         task.FifoAckEvent,
							Err:                  err,
						}
						return
					}

					step := 0
					for _, e := range execution.ActionExecutions {
						if e.Executed {
							step++
						} else {
							break
						}
					}

					action := execution.ActionExecutions[step].Action
					skipForChild := false
					if action.Parameters.SkipForChild != nil {
						skipForChild = *action.Parameters.SkipForChild
					}
					skipForInstruction := false
					if action.Parameters.SkipForInstruction != nil {
						skipForInstruction = *action.Parameters.SkipForInstruction
					}
					e.taskChannel <- Task{
						Source:               "input listener",
						Action:               action,
						Alarm:                task.Alarm,
						Entity:               task.Entity,
						Step:                 step,
						ExecutionID:          execution.ID,
						ExecutionCacheKey:    execution.GetCacheKey(),
						ScenarioID:           execution.ScenarioID,
						ScenarioName:         execution.ScenarioName,
						SkipForChild:         skipForChild,
						IsMetaAlarmUpdated:   execution.IsMetaAlarmUpdated,
						SkipForInstruction:   skipForInstruction,
						IsInstructionMatched: execution.IsInstructionMatched,
						Header:               execution.Header,
						Response:             execution.Response,
						ResponseMap:          execution.ResponseMap,
						AdditionalData:       task.AdditionalData,
					}

					return
				}

				ok, err := e.processTriggers(ctx, task)
				if err != nil {
					e.logger.Err(err).Str("alarm", task.Alarm.ID).Msg("cannot run scenarios")
					e.outputChannel <- ScenarioResult{
						Alarm:                task.Alarm,
						StartEventProcessing: task.Start,
						FifoAckEvent:         task.FifoAckEvent,
						Err:                  err,
					}
					return
				}

				if !ok {
					e.outputChannel <- ScenarioResult{
						Alarm:                task.Alarm,
						StartEventProcessing: task.Start,
						FifoAckEvent:         task.FifoAckEvent,
					}
				}
			}(ctx, scenariosTask)
		}
	}
}

func (e *redisBasedManager) finishExecution(
	ctx context.Context,
	alarm types.Alarm,
	execution ScenarioExecution,
	executionErr error,
) {
	if execution.Tries > 0 {
		err := e.executionStorage.Del(ctx, execution.GetCacheKey())
		if err != nil {
			e.logger.Err(err).Str("execution", execution.GetCacheKey()).Msg("cannot delete execution")
			return
		}

		return
	}

	count, err := e.executionStorage.Inc(ctx, alarm.ID, -1, false)
	if err != nil {
		e.logger.Err(err).Str("execution", execution.GetCacheKey()).Msg("cannot decrease counter")
		return
	}

	err = e.executionStorage.Del(ctx, execution.GetCacheKey())
	if err != nil {
		e.logger.Err(err).Str("execution", execution.GetCacheKey()).Msg("cannot delete execution")
		return
	}

	if count > 0 {
		return
	}

	select {
	case <-ctx.Done():
		return
	default:
		e.outputChannel <- ScenarioResult{
			Alarm:                alarm,
			StartEventProcessing: time.Unix(execution.StartEventProcessing, 0),
			Err:                  executionErr,
			ActionExecutions:     execution.ActionExecutions,
			FifoAckEvent:         execution.FifoAckEvent,
		}
	}
}

func (e *redisBasedManager) listenRPCResultChannel(ctx context.Context, wg *sync.WaitGroup, rpcResultChannel <-chan RpcResult) {
	e.logger.Debug().Msg("start listen engine-axe rpc responses")
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			e.logger.Debug().Msg("rpc listener cancelled")
			return
		case r, ok := <-rpcResultChannel:
			if !ok {
				e.logger.Debug().Msg("rpc channel closed")
				return
			}

			go func(ctx context.Context, result RpcResult) {
				taskRes := TaskResult{
					Source: "rpc listener",
				}
				execIdAndStep := strings.Split(result.CorrelationID, "&&")
				if len(execIdAndStep) < 2 {
					taskRes.Status = TaskRpcError
					taskRes.Err = fmt.Errorf("bad correlation id format: %s, impossible to get step and execution id", result.CorrelationID)

					e.processTaskResult(ctx, taskRes)

					return
				}

				executionCacheKey := execIdAndStep[0]
				step, err := strconv.Atoi(execIdAndStep[1])
				if err != nil {
					taskRes.Status = TaskRpcError
					taskRes.Err = fmt.Errorf("bad correlation id format: %s, failed to convert %s to int, error = %s", result.CorrelationID, execIdAndStep[1], err.Error())
					e.processTaskResult(ctx, taskRes)

					return
				}
				if result.Alarm == nil {
					taskRes.Status = TaskRpcError
					taskRes.Err = fmt.Errorf("missing alarm for correlation id: %s", result.CorrelationID)
					e.processTaskResult(ctx, taskRes)

					return
				}

				taskRes.Alarm = *r.Alarm
				taskRes.Step = step
				taskRes.AlarmChangeType = result.AlarmChangeType
				taskRes.ExecutionCacheKey = executionCacheKey
				taskRes.Header = result.WebhookHeader
				taskRes.Response = result.WebhookResponse

				if r.Error != nil {
					taskRes.Status = TaskRpcError
					taskRes.Err = r.Error
				}

				e.processTaskResult(ctx, taskRes)
			}(ctx, r)
		}
	}
}

func (e *redisBasedManager) listenTaskResultChannel(ctx context.Context, wg *sync.WaitGroup,
	taskResultChannel <-chan TaskResult) {
	e.logger.Debug().Msg("start listen workers results")
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			e.logger.Debug().Msg("task result listener cancelled")
			return
		case t, ok := <-taskResultChannel:
			if !ok {
				e.logger.Debug().Msg("task result channel closed")
				return
			}

			go e.processTaskResult(ctx, t)
		}
	}
}

func (e *redisBasedManager) processTaskResult(ctx context.Context, taskRes TaskResult) {
	if taskRes.ExecutionCacheKey == "" {
		e.logger.Error().Err(taskRes.Err).Msg("cannot get execution")
		return
	}

	scenarioExecution, err := e.executionStorage.Get(ctx, taskRes.ExecutionCacheKey)
	if err != nil || scenarioExecution == nil {
		e.logger.Error().Err(err).Str("execution", taskRes.ExecutionCacheKey).Msg("cannot get execution")
		return
	}

	if taskRes.Status == TaskCancelled {
		e.logger.Warn().Msgf("worker task was cancelled, error = %s", taskRes.Err.Error())

		if taskRes.ExecutionCacheKey != "" {
			e.logger.Debug().Str("source", taskRes.Source).
				Str("alarm", taskRes.Alarm.ID).
				Str("execution", taskRes.ExecutionCacheKey).
				Int("step", taskRes.Step).
				Msg("Worker returned error, drop scenario")
			e.finishExecution(ctx, taskRes.Alarm, *scenarioExecution, taskRes.Err)
		}

		return
	}

	if taskRes.Err != nil {
		e.logger.Err(taskRes.Err).
			Str("source", taskRes.Source).
			Str("alarm", taskRes.Alarm.ID).
			Str("execution", taskRes.ExecutionCacheKey).
			Int("step", taskRes.Step).Msg("Execution failed, drop scenario")
		e.finishExecution(ctx, taskRes.Alarm, *scenarioExecution, taskRes.Err)

		return
	}

	if taskRes.Status == TaskNotMatched && scenarioExecution.ActionExecutions[taskRes.Step].Action.DropScenarioIfNotMatched {
		e.logger.Debug().
			Str("source", taskRes.Source).
			Str("alarm", taskRes.Alarm.ID).
			Str("execution", taskRes.ExecutionCacheKey).Int("step", taskRes.Step).
			Msg("Action is not matched, drop scenario")
		e.finishExecution(ctx, taskRes.Alarm, *scenarioExecution, nil)
		return
	}

	scenarioExecution.ActionExecutions[taskRes.Step].Executed = true
	scenarioExecution.LastUpdate = time.Now().Unix()
	if scenarioExecution.Header == nil {
		scenarioExecution.Header = make(map[string]string)
	}
	if scenarioExecution.Response == nil {
		scenarioExecution.Response = make(map[string]interface{})
	}

	if scenarioExecution.ResponseMap == nil {
		scenarioExecution.ResponseMap = make(map[string]interface{})
	}

	for k, v := range taskRes.Header {
		scenarioExecution.Header[k] = v
	}

	responseCountStr := strconv.Itoa(scenarioExecution.ResponseCount)
	for k, v := range taskRes.Response {
		scenarioExecution.Response[k] = v
		scenarioExecution.ResponseMap[responseCountStr+"."+k] = v
	}

	if taskRes.Response != nil {
		scenarioExecution.ResponseCount++
	}

	err = e.executionStorage.Update(ctx, *scenarioExecution)
	if err != nil {
		e.logger.Err(err).Str("execution", scenarioExecution.GetCacheKey()).Msg("cannot save execution")
		e.finishExecution(ctx, taskRes.Alarm, *scenarioExecution, err)
		return
	}

	if scenarioExecution.ActionExecutions[taskRes.Step].Action.EmitTrigger &&
		taskRes.AlarmChangeType != types.AlarmChangeTypeNone {
		err := e.processEmittedTrigger(ctx, taskRes, *scenarioExecution)
		if err != nil {
			e.logger.Err(err).Str("execution", scenarioExecution.GetCacheKey()).Msg("cannot process emitted trigger")
			e.finishExecution(ctx, taskRes.Alarm, *scenarioExecution, err)
			return
		}
	}

	nextStep := taskRes.Step + 1
	if len(scenarioExecution.ActionExecutions) > nextStep {
		additionalData := scenarioExecution.AdditionalData
		additionalData.AlarmChangeType = taskRes.AlarmChangeType
		action := scenarioExecution.ActionExecutions[nextStep].Action
		skipForChild := false
		if action.Parameters.SkipForChild != nil {
			skipForChild = *action.Parameters.SkipForChild
		}
		skipForInstruction := false
		if action.Parameters.SkipForInstruction != nil {
			skipForInstruction = *action.Parameters.SkipForInstruction
		}
		nextTask := Task{
			Source:               "process task func",
			Action:               action,
			Alarm:                taskRes.Alarm,
			Entity:               scenarioExecution.Entity,
			Step:                 nextStep,
			ExecutionID:          scenarioExecution.ID,
			ExecutionCacheKey:    scenarioExecution.GetCacheKey(),
			ScenarioID:           scenarioExecution.ScenarioID,
			ScenarioName:         scenarioExecution.ScenarioName,
			SkipForChild:         skipForChild,
			IsMetaAlarmUpdated:   scenarioExecution.IsMetaAlarmUpdated,
			SkipForInstruction:   skipForInstruction,
			IsInstructionMatched: scenarioExecution.IsInstructionMatched,
			Header:               scenarioExecution.Header,
			Response:             scenarioExecution.Response,
			ResponseMap:          scenarioExecution.ResponseMap,
			AdditionalData:       additionalData,
		}

		select {
		case <-ctx.Done():
			return
		default:
			e.taskChannel <- nextTask
		}
	} else {
		e.logger.Debug().
			Str("source", taskRes.Source).
			Str("alarm", taskRes.Alarm.ID).
			Str("execution", taskRes.ExecutionCacheKey).
			Int("step", taskRes.Step).
			Msg("Scenario is finished")
		e.finishExecution(ctx, taskRes.Alarm, *scenarioExecution, nil)
	}
}

func (e *redisBasedManager) processTriggers(ctx context.Context, task ExecuteScenariosTask) (bool, error) {
	err := e.scenarioStorage.RunDelayedScenarios(ctx, task.Triggers, task.Alarm, task.Entity, task.AdditionalData)
	if err != nil {
		return false, err
	}

	scenariosByTrigger, err := e.scenarioStorage.GetTriggeredScenarios(task.Triggers, task.Alarm)
	if err != nil {
		return false, err
	}

	scenariosCount := 0
	for _, scenarios := range scenariosByTrigger {
		scenariosCount += len(scenarios)
	}

	if scenariosCount == 0 {
		return false, nil
	}

	_, err = e.executionStorage.Inc(ctx, task.Alarm.ID, int64(scenariosCount), true)
	if err != nil {
		return false, err
	}

	additionalData := task.AdditionalData
	for trigger, scenarios := range scenariosByTrigger {
		additionalData.Trigger = trigger
		for _, scenario := range scenarios {
			e.startExecution(ctx, scenario, task.Alarm, task.Entity, additionalData, task.FifoAckEvent,
				task.Start, task.IsMetaAlarmUpdated, task.IsInstructionMatched)
		}
	}

	return true, nil
}

func (e *redisBasedManager) processEmittedTrigger(
	ctx context.Context,
	prevTaskRes TaskResult,
	prevScenarioExecution ScenarioExecution,
) error {
	additionalData := prevScenarioExecution.AdditionalData
	additionalData.AlarmChangeType = prevTaskRes.AlarmChangeType

	alarmChange := types.AlarmChange{Type: prevTaskRes.AlarmChangeType}
	triggers := alarmChange.GetTriggers()

	if len(triggers) == 0 {
		return nil
	}
	err := e.scenarioStorage.RunDelayedScenarios(ctx, triggers, prevTaskRes.Alarm, prevScenarioExecution.Entity, additionalData)
	if err != nil {
		return err
	}

	scenariosByTrigger, err := e.scenarioStorage.GetTriggeredScenarios(triggers, prevTaskRes.Alarm)
	if err != nil {
		return err
	}

	scenariosCount := 0
	for _, scenarios := range scenariosByTrigger {
		scenariosCount += len(scenarios)
	}

	if scenariosCount == 0 {
		return nil
	}

	_, err = e.executionStorage.Inc(ctx, prevTaskRes.Alarm.ID, int64(scenariosCount), false)
	if err != nil {
		return err
	}

	for trigger, scenarios := range scenariosByTrigger {
		additionalData.Trigger = trigger
		for _, scenario := range scenarios {
			e.startExecution(ctx, scenario, prevTaskRes.Alarm, prevScenarioExecution.Entity, additionalData,
				prevScenarioExecution.FifoAckEvent, time.Unix(prevScenarioExecution.StartEventProcessing, 0), prevScenarioExecution.IsMetaAlarmUpdated, prevScenarioExecution.IsInstructionMatched)
		}
	}

	return nil
}

func (e *redisBasedManager) startExecution(
	ctx context.Context,
	scenario Scenario,
	alarm types.Alarm,
	entity types.Entity,
	data AdditionalData,
	fifoAckEvent types.Event,
	start time.Time,
	isMetaAlarmUpdated bool,
	isInstructionMatched bool,
) {
	e.logger.Debug().Msgf("Execute scenario = %s for alarm = %s", scenario.ID, alarm.ID)
	var executions []Execution
	for _, action := range scenario.Actions {
		executions = append(
			executions,
			Execution{
				Action:   action,
				Executed: false,
			},
		)
	}

	data.RuleName = scenario.Name

	execution := ScenarioExecution{
		ID:                   utils.NewID(),
		ScenarioID:           scenario.ID,
		ScenarioName:         scenario.Name,
		AlarmID:              alarm.ID,
		Entity:               entity,
		ActionExecutions:     executions,
		LastUpdate:           time.Now().Unix(),
		AdditionalData:       data,
		FifoAckEvent:         fifoAckEvent,
		IsMetaAlarmUpdated:   isMetaAlarmUpdated,
		IsInstructionMatched: isInstructionMatched,
		StartEventProcessing: start.Unix(),
	}
	ok, err := e.executionStorage.Create(ctx, execution)
	if err != nil {
		e.logger.Err(err).Str("scenario", scenario.ID).Str("alarm", alarm.ID).Msg("cannot save execution")
		return
	}
	if !ok {
		e.logger.Error().Str("scenario", scenario.ID).Str("alarm", alarm.ID).Msg("scenario is already executing")
		return
	}

	action := scenario.Actions[0]
	skipForChild := false
	if action.Parameters.SkipForChild != nil {
		skipForChild = *action.Parameters.SkipForChild
	}
	skipForInstruction := false
	if action.Parameters.SkipForInstruction != nil {
		skipForInstruction = *action.Parameters.SkipForInstruction
	}
	e.taskChannel <- Task{
		Source:               "input listener",
		Action:               action,
		Alarm:                alarm,
		Entity:               entity,
		Step:                 0,
		ExecutionID:          execution.ID,
		ExecutionCacheKey:    execution.GetCacheKey(),
		ScenarioID:           scenario.ID,
		ScenarioName:         scenario.Name,
		SkipForChild:         skipForChild,
		IsMetaAlarmUpdated:   execution.IsMetaAlarmUpdated,
		SkipForInstruction:   skipForInstruction,
		IsInstructionMatched: execution.IsInstructionMatched,
		AdditionalData:       data,
	}
}
