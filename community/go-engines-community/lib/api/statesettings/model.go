package statesettings

type StateSettingRequest struct {
	ID              string           `json:"-" bson:"-"`
	Method          string           `json:"method" bson:"method" binding:"required"`
	Type            string           `json:"type" bson:"type" binding:"required"`
	JunitThresholds *JUnitThresholds `json:"junit_thresholds,omitempty" bson:"junit_thresholds,omitempty"`
}

type StateSetting struct {
	ID              string           `json:"_id" bson:"_id"`
	Method          string           `json:"method" bson:"method"`
	Type            string           `json:"type" bson:"type"`
	JunitThresholds *JUnitThresholds `json:"junit_thresholds,omitempty" bson:"junit_thresholds,omitempty"`
}

type StateThresholds struct {
	Minor    *float64 `json:"minor" bson:"minor" binding:"required,numeric,gte=0,lte=100,ltefield=Major,ltefield=Critical"`
	Major    *float64 `json:"major" bson:"major" binding:"required,numeric,gte=0,lte=100,ltefield=Critical"`
	Critical *float64 `json:"critical" bson:"critical" binding:"required,numeric,gte=0,lte=100"`
	Type     *int     `json:"type" bson:"type" binding:"required"`
}

type JUnitThresholds struct {
	Skipped  StateThresholds `json:"skipped" bson:"skipped" binding:"required"`
	Errors   StateThresholds `json:"errors" bson:"errors" binding:"required"`
	Failures StateThresholds `json:"failures" bson:"failures" binding:"required"`
}

type AggregationResult struct {
	Data       []StateSetting `bson:"data" json:"data"`
	TotalCount int64          `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
