package action

//go:generate mockgen -destination=../../../mocks/lib/canopsis/action/action.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action Adapter,DelayedScenarioManager,DelayedScenarioStorage,ScenarioExecutionStorage,ScenarioStorage,WorkerPool

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Adapter interface {
	GetEnabled(ctx context.Context) ([]Scenario, error)
	GetEnabledById(ctx context.Context, id string) (Scenario, error)
	GetEnabledByIDs(ctx context.Context, ids []string) ([]Scenario, error)
}

// Service allows you to manipulate actions.
type Service interface {
	// Process parse an event to see if an action is suitable.
	Process(ctx context.Context, event *types.Event) error

	// ListenScenarioFinish receives message when all scenarios for event are finished
	// and acknowledges fifo.
	ListenScenarioFinish(ctx context.Context, channel <-chan ScenarioResult)

	// ProcessAbandonedExecutions checks execution storage and processes executions which
	// weren't updated for a long time
	ProcessAbandonedExecutions(ctx context.Context) error
}

// ScenarioStorage is used to provide scenarios.
type ScenarioStorage interface {
	// ReloadScenarios trigger a refresh on scenarios cache from DB
	ReloadScenarios(ctx context.Context) error

	// GetTriggeredScenarios returns scenarios which are triggered by triggers.
	GetTriggeredScenarios(
		triggers []string,
		alarm types.Alarm,
	) (triggered []Scenario, err error)

	// RunDelayedScenarios starts delay timeout for scenarios which are triggered by triggers.
	RunDelayedScenarios(
		ctx context.Context,
		triggers []string,
		alarm types.Alarm,
		entity types.Entity,
		additionalData AdditionalData,
	) error

	// GetScenario returns scenario.
	GetScenario(id string) *Scenario
}

// TaskManager is used to execute scenarios.
type TaskManager interface {
	Run(ctx context.Context, rpcResultChannel <-chan RpcResult,
		inputChannel <-chan ExecuteScenariosTask) (<-chan ScenarioResult, error)
}

type ScenarioExecution struct {
	ID                 string                 `json:"-"`
	ScenarioID         string                 `json:"-"`
	AlarmID            string                 `json:"-"`
	Entity             types.Entity           `json:"e"`
	ActionExecutions   []Execution            `json:"ae"`
	LastUpdate         int64                  `json:"u"`
	AckResources       bool                   `json:"ar"`
	Tries              int64                  `json:"t"`
	Header             map[string]string      `json:"h,omitempty"`
	Response           map[string]interface{} `json:"r,omitempty"`
	ResponseMap        map[string]interface{} `json:"rm,omitempty"`
	ResponseCount      int                    `json:"rc"`
	AdditionalData     AdditionalData         `json:"ad"`
	FifoAckEvent       types.Event            `json:"fev"`
	IsMetaAlarmUpdated bool                   `json:"mau,omitempty"`
}

type ScenarioResult struct {
	Alarm            types.Alarm
	Err              error
	ActionExecutions []Execution
	FifoAckEvent     types.Event
}

type ExecuteScenariosTask struct {
	Triggers             []string
	DelayedScenarioID    string
	AbandonedExecutionID string
	Entity               types.Entity
	Alarm                types.Alarm
	AckResources         bool
	AdditionalData       AdditionalData
	FifoAckEvent         types.Event
	IsMetaAlarmUpdated   bool
}

type AdditionalData struct {
	AlarmChangeType types.AlarmChangeType `json:"alarm_change_type"`
	Author          string                `json:"author"`
	User            string                `json:"user"`
	Initiator       string                `json:"initiator"`
	Output          string                `json:"event_output"`
}

type Execution struct {
	Action   Action `json:"a"`
	Executed bool   `json:"e"`
}

type RpcResult struct {
	CorrelationID   string
	Alarm           *types.Alarm
	AlarmChangeType types.AlarmChangeType
	Header          map[string]string
	Response        map[string]interface{}
	Error           error
}
