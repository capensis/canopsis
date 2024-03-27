package canopsis

import (
	"time"
)

// Globals
const (
	AppName = "canopsis"

	ActionEngineName            = "engine-action"
	ActionQueueName             = "Engine_action"
	ActionAxeRPCClientQueueName = "Engine_action_axe_rpc_client"
	ActionConsumerName          = "action"
	ActionRPCConsumerName       = "action_rpc"
	ActionConnector             = "action"

	AxeEngineName                     = "engine-axe"
	AxeQueueName                      = "Engine_axe"
	AxePbehaviorRPCClientQueueName    = "Engine_axe_pbehavior_rpc_client"
	AxeDynamicInfosRPCClientQueueName = "Engine_axe_dynamic_infos_rpc_client"
	AxeRPCQueueServerName             = "Engine_axe_rpc_server"
	AxeConsumerName                   = "axe"
	AxeRPCConsumerName                = "axe_rpc"
	AxeConnector                      = "axe"

	CheExchangeName = ""
	CheEngineName   = "engine-che"
	CheQueueName    = "Engine_che"
	CheConsumerName = "che"
	CheConnector    = "che"

	DefaultBulkSize      = 1000
	DefaultBulkBytesSize = 16000000 // < MongoDB limit (16 megabytes)
	DefaultEventAuthor   = "system"

	DynamicInfosEngineName         = "engine-dynamic-infos"
	DynamicInfosQueueName          = "Engine_dynamic_infos"
	DynamicInfosConsumerName       = "dynamic-infos"
	DynamicInfosRPCConsumerName    = "dynamic-infos_rpc"
	DynamicInfosRPCQueueServerName = "Engine_dynamic_infos_rpc_server"

	PBehaviorEngineName         = "engine-pbehavior"
	PBehaviorRPCQueueServerName = "Engine_pbehavior_rpc_server"
	PBehaviorQueueRecomputeName = "Engine_pbehavior_recompute"
	PBehaviorRPCConsumerName    = "pbehavior_rpc"
	PBehaviorConsumerName       = "pbehavior"
	PBehaviorConnector          = "pbehavior"

	WebhookEngineName         = "engine-webhook"
	WebhookRPCQueueServerName = "Engine_webhook_rpc_server"
	WebhookRPCConsumerName    = "webhook_rpc"

	FIFOEngineName      = "engine-fifo"
	FIFOExchangeName    = ""
	FIFOQueueName       = "Engine_fifo"
	FIFOAckExchangeName = ""
	FIFOAckQueueName    = "FIFO_ack"
	FIFOConsumerName    = "fifo"
	FIFOAckConsumerName = "fifo_ack"

	CorrelationEngineName            = "engine-correlation"
	CorrelationQueueName             = "Engine_correlation"
	CorrelationAxeRPCClientQueueName = "Engine_correlation_axe_rpc_client"
	CorrelationConsumerName          = "correlation"
	CorrelationRPCConsumerName       = "correlation_rpc"
	CorrelationConnector             = "correlation"

	PeriodicalWaitTime     = time.Minute
	JsonContentType        = "application/json"
	CanopsisEventsExchange = "canopsis.events"

	RemediationEngineName            = "engine-remediation"
	RemediationConsumerName          = "remediation"
	RemediationRPCConsumerName       = "remediation_rpc"
	RemediationRPCQueueServerName    = "Engine_remediation_rpc_server"
	RemediationRPCQueueServerJobName = "Engine_remediation_rpc_server_job"
	RemediationConnector             = "remediation"

	TechMetricsFlushInterval = time.Second * 10

	DefaultFlushInterval = time.Second * 5

	FacetLimit = 1000

	ApiName      = "api"
	ApiConnector = "api"

	DefaultEventWorkers = 10

	DefaultSystemAlarmConnector = "system"
)
