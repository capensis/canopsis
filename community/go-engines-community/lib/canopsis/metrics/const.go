package metrics

// criteria type
const (
	EntityCriteriaType = iota
	UserCriteriaType
)

// tech metrics
const (
	FIFOQueue = "fifo_queue"
	FIFOEvent = "fifo_event"
	CheEvent  = "che_event"
	AxeEvent  = "axe_event"
)
