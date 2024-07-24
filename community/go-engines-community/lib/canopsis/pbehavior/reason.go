package pbehavior

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type Reason struct {
	ID          string           `bson:"_id,omitempty" json:"_id"`
	Name        string           `bson:"name" json:"name" binding:"required"`
	Description string           `bson:"description" json:"description" binding:"required"`
	Author      string           `bson:"author" json:"author"`
	Created     datetime.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated     datetime.CpsTime `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`

	// Hidden is used in API to hide documents from the list response
	Hidden bool `bson:"hidden" json:"hidden"`
}
