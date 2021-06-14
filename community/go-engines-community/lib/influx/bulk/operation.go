package bulk

import "time"

// PointOp regroups all parameters for influx.NewPoint.
// All members are mendatory.
type PointOp struct {
	Name   string
	Tags   map[string]string
	Fields map[string]interface{}
	Time   time.Time
}
