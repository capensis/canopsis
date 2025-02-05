package canopsis

import (
	"time"
)

// Workers
const (
	DefaultEventWorkers         = 10
	DefaultUserEventWorkers     = 2
	DefaultSystemEventWorkers   = 4
	DefaultExternalEventWorkers = 4
	DefaultRpcWorkers           = 4
)

// App names
const (
	AppName = "canopsis"

	ActionEngineName       = "engine-action"
	ApiName                = "api"
	AxeEngineName          = "engine-axe"
	CheEngineName          = "engine-che"
	CorrelationEngineName  = "engine-correlation"
	DynamicInfosEngineName = "engine-dynamic-infos"
	FIFOEngineName         = "engine-fifo"
	PBehaviorEngineName    = "engine-pbehavior"
	RemediationEngineName  = "engine-remediation"
	SnmpEngineName         = "engine-snmp"
	WebhookEngineName      = "engine-webhook"
)

// Internal connectors
const (
	ActionConnector      = "action"
	ApiConnector         = "api"
	AxeConnector         = "axe"
	CheConnector         = "che"
	CorrelationConnector = "correlation"
	PBehaviorConnector   = "pbehavior"
	RemediationConnector = "remediation"

	DefaultSystemAlarmConnector = "system"
)

// Exchanges
const (
	DefaultExchangeName = ""
	EngineExchangeName  = "canopsis.engine"
	EventsExchangeName  = "canopsis.events"
)

// Queues
const (
	ActionQueuePrefix           = "Engine_action"
	ActionExternalQueueName     = ActionQueuePrefix + "_external"
	ActionSystemQueueName       = ActionQueuePrefix + "_system"
	ActionUserQueueName         = ActionQueuePrefix + "_user"
	ActionAxeRPCClientQueueName = ActionQueuePrefix + "_axe_rpc_client"

	AxeQueuePrefix                    = "Engine_axe"
	AxeExternalQueueName              = AxeQueuePrefix + "_external"
	AxeSystemQueueName                = AxeQueuePrefix + "_system"
	AxeUserQueueName                  = AxeQueuePrefix + "_user"
	AxePbehaviorRPCClientQueueName    = AxeQueuePrefix + "_pbehavior_rpc_client"
	AxeDynamicInfosRPCClientQueueName = AxeQueuePrefix + "_dynamic_infos_rpc_client"
	AxeRPCQueueServerName             = AxeQueuePrefix + "_rpc_server"

	CheQueuePrefix       = "Engine_che"
	CheExternalQueueName = CheQueuePrefix + "_external"
	CheSystemQueueName   = CheQueuePrefix + "_system"
	CheUserQueueName     = CheQueuePrefix + "_user"

	CorrelationQueuePrefix           = "Engine_correlation"
	CorrelationExternalQueueName     = CorrelationQueuePrefix + "_external"
	CorrelationSystemQueueName       = CorrelationQueuePrefix + "_system"
	CorrelationUserQueueName         = CorrelationQueuePrefix + "_user"
	CorrelationAxeRPCClientQueueName = CorrelationQueuePrefix + "_axe_rpc_client"

	DynamicInfosQueuePrefix        = "Engine_dynamic_infos"
	DynamicInfosExternalQueueName  = DynamicInfosQueuePrefix + "_external"
	DynamicInfosSystemQueueName    = DynamicInfosQueuePrefix + "_system"
	DynamicInfosUserQueueName      = DynamicInfosQueuePrefix + "_user"
	DynamicInfosRPCQueueServerName = DynamicInfosQueuePrefix + "_rpc_server"

	FIFOQueueName    = "Engine_fifo"
	FIFOAckQueueName = "FIFO_ack"

	PBehaviorRPCQueueServerName = "Engine_pbehavior_rpc_server"
	PBehaviorQueueRecomputeName = "Engine_pbehavior_recompute"

	RemediationRPCQueueServerName    = "Engine_remediation_rpc_server"
	RemediationRPCQueueServerJobName = "Engine_remediation_rpc_server_job"

	WebhookRPCQueueServerName = "Engine_webhook_rpc_server"
)

// Consumers
const (
	ActionExternalConsumerName = "action_external"
	ActionSystemConsumerName   = "action_system"
	ActionUserConsumerName     = "action_user"
	ActionRPCConsumerName      = "action_rpc"

	AxeExternalConsumerName = "axe_external"
	AxeSystemConsumerName   = "axe_system"
	AxeUserConsumerName     = "axe_user"
	AxeRPCConsumerName      = "axe_rpc"

	CheExternalConsumerName = "che_external"
	CheSystemConsumerName   = "che_system"
	CheUserConsumerName     = "che_user"

	CorrelationExternalConsumerName = "correlation_external"
	CorrelationSystemConsumerName   = "correlation_system"
	CorrelationUserConsumerName     = "correlation_user"
	CorrelationRPCConsumerName      = "correlation_rpc"

	DynamicInfosExternalConsumerName = "dynamic_infos_external"
	DynamicInfosSystemConsumerName   = "dynamic_infos_system"
	DynamicInfosUserConsumerName     = "dynamic_infos_user"
	DynamicInfosRPCConsumerName      = "dynamic_infos_rpc"

	FIFOConsumerName    = "fifo"
	FIFOAckConsumerName = "fifo_ack"

	PBehaviorRPCConsumerName = "pbehavior_rpc"
	PBehaviorConsumerName    = "pbehavior"

	RemediationRPCConsumerName = "remediation_rpc"

	WebhookRPCConsumerName = "webhook_rpc"
)

// Other
const (
	DefaultBulkSize      = 1000
	DefaultBulkBytesSize = 16000000 // < MongoDB limit (16 megabytes)

	DefaultEventAuthor = "system"

	JsonContentType = "application/json"

	RecorderConsumerName    = "recorder"
	RecorderRPCConsumerName = "recorder_rpc"

	PeriodicalWaitTime       = time.Minute
	TechMetricsFlushInterval = time.Second * 10
	DefaultFlushInterval     = time.Second * 5

	FacetLimit = 1000

	SubDirUpload   = "upload-files"
	SubDirIcons    = "icons"
	SubDirExport   = "export-files"
	SubDirImport   = "import-files"
	SubDirJunit    = "junit-files"
	SubDirJunitAPI = "junit-api-files"
)
