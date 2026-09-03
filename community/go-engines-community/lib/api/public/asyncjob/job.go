package asyncjob

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

const (
	StatusPending = "pending"
	StatusOngoing = "ongoing"
	StatusFailed  = "failed"
	StatusDone    = "done"
)

// DirPerm is the permission of the directories holding the job files.
const DirPerm = 0o770

// FilePerm is the permission of the job files.
const FilePerm = 0o640

// Response is what the API answers for a job, whatever its type.
type Response struct {
	ID           string `json:"_id"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// Base holds the lifecycle fields shared by every job type.
type Base struct {
	ID        string            `bson:"_id" json:"_id"`
	Created   datetime.CpsTime  `bson:"created" json:"created"`
	LastPing  datetime.CpsTime  `bson:"last_ping" json:"last_ping"`
	Completed *datetime.CpsTime `bson:"completed,omitempty" json:"completed"`
	Status    string            `bson:"status" json:"status"`
	Attempts  int64             `bson:"attempts" json:"attempts"`
	// LeaseGeneration is incremented every time a worker claims the job. It fences a
	// worker which outlived its lease from updating the job after another worker claimed it.
	LeaseGeneration int64  `bson:"lease_generation" json:"-"`
	DurationMS      int64  `bson:"duration_ms,omitempty" json:"duration_ms"`
	ErrorMessage    string `bson:"error_message,omitempty"`
}

func NewBase(id string) Base {
	now := datetime.NewCpsTime()

	return Base{
		ID:       id,
		Created:  now,
		LastPing: now,
		Status:   StatusPending,
		Attempts: 1,
	}
}

func (b Base) JobID() string {
	return b.ID
}

func (b Base) JobStatus() string {
	return b.Status
}

func (b Base) Response() Response {
	return Response{
		ID:           b.ID,
		Status:       b.Status,
		ErrorMessage: b.ErrorMessage,
	}
}
