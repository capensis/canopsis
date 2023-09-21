package trigger

type Trigger struct {
	Type       string     `json:"type"`
	Parameters Parameters `json:"parameters"`
}

func (t *Trigger) IsTriggered() bool {
	return false
}

type Parameters struct {
	EventsCountThreshold int `json:"events_count_threshold"`
}
