package bulk

// Bulk interface for InfluxDB
type Bulk interface {
	AddPoints(po ...PointOp) error
	Perform() error
}
