package pbehavior

import (
	libtime "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/time"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Exdate struct {
	types.Exdate `bson:"inline"`
	Type         string `bson:"type" json:"type"`
}

type Exception struct {
	ID          string           `bson:"_id,omitempty" json:"_id"`
	Name        string           `bson:"name" json:"name"`
	Description string           `bson:"description" json:"description"`
	Exdates     []Exdate         `bson:"exdates" json:"exdates"`
	Created     *libtime.CpsTime `bson:"created,omitempty" json:"created"`
}
