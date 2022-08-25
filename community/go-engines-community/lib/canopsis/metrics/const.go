package metrics

// criteria type
const (
	EntityCriteriaType = iota
	UserCriteriaType
)

// tech metrics
const (
	FIFOEvent         = "fifo_event"
	CheEvent          = "che_event"
	PBehaviorEvent    = "pbehavior_event"
	AxeEvent          = "axe_event"
	CorrelationEvent  = "correlation_event"
	ServiceEvent      = "service_event"
	DynamicInfosEvent = "dynamic_infos_event"
	ActionEvent       = "action_event"

	FIFOQueue           = "fifo_queue"
	AxePeriodical       = "axe_periodical"
	PBehaviorPeriodical = "pbehavior_periodical"
	CheInfos            = "che_infos"
	ApiRequests         = "api_requests"
)
