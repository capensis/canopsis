package heartbeat

type Adapter interface {
	Get() ([]Heartbeat, error)
}
