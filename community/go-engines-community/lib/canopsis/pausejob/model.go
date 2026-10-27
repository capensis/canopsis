package pausejob

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// Job owns pause delivery until each recipient commits its state and continuation.
// Initialized distinguishes a job awaiting discovery from a completed empty job.
type Job struct {
	ID          string      `bson:"_id"`
	Parent      string      `bson:"parent"`
	ParentType  string      `bson:"parent_type"`
	Event       types.Event `bson:"event"`
	Initialized bool        `bson:"initialized"`
	Pending     []string    `bson:"pending"`
	RetryAt     time.Time   `bson:"retry_at"`
}
